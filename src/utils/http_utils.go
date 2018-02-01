package utils

import (
	"log"
	"strconv"
	"net/http"
	"encoding/json"
	"io"
	"fmt"
	"io/ioutil"
	"strings"
	"time"
	"compress/gzip"
	"path/filepath"
	"bytes"
	"os"
	"mime"
	"crypto/sha256"
	"encoding/hex"
)

func ServeJSON(w http.ResponseWriter, j interface{}) {
	if j == nil {
		http.Error(w, "empty response data", 400)
		return
	}
	// w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(j)
	// for request: json.NewDecoder(res.Body).Decode(&body)
}

func Abort(w http.ResponseWriter, message string, code int) {
	http.Error(w, message, code)
}

func GetQueryInt(r *http.Request, key string, defaultValue int) int {
	value, err := strconv.Atoi(r.URL.Query().Get(key))
	if err != nil {
		log.Println("Error: get params cols error:", err)
		return defaultValue
	}
	return value
}

func GetQueryInt32(r *http.Request, key string, defaultValue uint32) uint32 {
	value, err := strconv.Atoi(r.URL.Query().Get(key))
	if err != nil {
		log.Println("Error: get params cols error:", err)
		return defaultValue
	}
	return uint32(value)
}

// serve all views files from memory storage.
// basic idea: https://github.com/bouk/staticfiles
type staticFilesFile struct {
	data  []byte
	mime  string
	mtime time.Time
	// size is the size before compression. If 0, it means the data is uncompressed
	size int64
	// hash is a sha256 hash of the file contents. Used for the Etag, and useful for caching
	hash string
}

var staticFiles = make(map[string]*staticFilesFile)

// NotFound is called when no asset is found.
// It defaults to http.NotFound but can be overwritten
var NotFound = http.NotFound

// read all files in views directory and map to "staticFiles"
func initHttpUtils() {
	files := processDir(Config.Site.ViewsDir, "")
	for _, file := range files {
		var b bytes.Buffer
		var b2 bytes.Buffer
		hash := sha256.New()

		f, err := os.Open(filepath.Join(Config.Site.ViewsDir,file))
		if err != nil {
			log.Fatal(err)
		}
		stat, err := f.Stat()
		if err != nil {
			log.Fatal(err)
		}
		if _, err := b.ReadFrom(f); err != nil {
			log.Fatal(err)
		}
		f.Close()

		compressedWriter, _ := gzip.NewWriterLevel(&b2, gzip.BestCompression)
		writer := io.MultiWriter(compressedWriter, hash)
		if _, err := writer.Write(b.Bytes()); err != nil {
			log.Fatal(err)
		}
		compressedWriter.Close()
		file = strings.Replace(file, "\\", "/", -1)
		if b2.Len() < b.Len() {
			staticFiles[file] = &staticFilesFile{
				data:  b2.Bytes(),
				mime:  mime.TypeByExtension(filepath.Ext(file)),
				mtime: time.Unix(stat.ModTime().Unix(), 0),
				size:  stat.Size(),
				hash:  hex.EncodeToString(hash.Sum(nil)),
			}
		} else {
			staticFiles[file] = &staticFilesFile{
				data:  b.Bytes(),
				mime:  mime.TypeByExtension(filepath.Ext(file)),
				mtime: time.Unix(stat.ModTime().Unix(), 0),
				hash:  hex.EncodeToString(hash.Sum(nil)),
			}
		}
		b.Reset()
		b2.Reset()
		hash.Reset()
	}
}

// todo large memory!!
func processDir(prefix, dir string) (fileSlice []string) {
	files, err := ioutil.ReadDir(filepath.Join(prefix, dir))
	var allFiles []string
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		if strings.HasPrefix(file.Name(), ".") {
			continue
		}

		dir := filepath.Join(dir, file.Name())
		//if skipFile(path.Join(id...), excludeSlice) {
		//	continue
		//}

		if file.IsDir() {
			for _, v := range processDir(prefix, dir) {
				allFiles = append(allFiles, v)
			}
		} else {
			allFiles = append(allFiles, dir)
		}
	}
	return allFiles
}

// ServeHTTP serves a request, attempting to reply with an embedded file.
func ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	filename := strings.TrimPrefix(req.URL.Path, "/")
	ServeHTTPByName(rw, req, filename)
}

// ServeHTTPByName serves a request by the key(param filename) in map.
func ServeHTTPByName(rw http.ResponseWriter, req *http.Request, filename string) {
	if f, ok := staticFiles[filename]; !ok {
		NotFound(rw, req)
		return
	} else {
		header := rw.Header()
		if f.hash != "" {
			if hash := req.Header.Get("If-None-Match"); hash == f.hash {
				rw.WriteHeader(http.StatusNotModified)
				return
			}
			header.Set("ETag", f.hash)
		}
		if !f.mtime.IsZero() {
			if t, err := time.Parse(http.TimeFormat, req.Header.Get("If-Modified-Since")); err == nil && f.mtime.Before(t.Add(1*time.Second)) {
				rw.WriteHeader(http.StatusNotModified)
				return
			}
			header.Set("Last-Modified", f.mtime.UTC().Format(http.TimeFormat))
		}
		header.Set("Content-Type", f.mime)

		// Check if the asset is compressed in the binary
		if f.size == 0 { // not compressed
			header.Set("Content-Length", strconv.Itoa(len(f.data)))
			rw.Write(f.data)
		} else {
			if header.Get("Content-Encoding") == "" && strings.Contains(req.Header.Get("Accept-Encoding"), "gzip") {
				header.Set("Content-Encoding", "gzip")
				header.Set("Content-Length", strconv.Itoa(len(f.data)))
				rw.Write(f.data)
			} else {
				header.Set("Content-Length", strconv.Itoa(int(f.size)))
				reader, _ := gzip.NewReader(bytes.NewReader(f.data))
				io.Copy(rw, reader)
				reader.Close()
			}
		}
	}
}

// Server is simply ServeHTTP but wrapped in http.HandlerFunc so it can be passed into net/http functions directly.
var Server http.Handler = http.HandlerFunc(ServeHTTP)

// Open allows you to read an embedded file directly. It will return a decompressing Reader if the file is embedded in compressed format.
// You should close the Reader after you're done with it.
func Open(name string) (io.ReadCloser, error) {
	f, ok := staticFiles[name]
	if !ok {
		return nil, fmt.Errorf("Asset %s not found", name)
	}

	if f.size == 0 {
		return ioutil.NopCloser(bytes.NewReader(f.data)), nil
	}
	return gzip.NewReader(bytes.NewReader(f.data))
}

// ModTime returns the modification time of the original file.
// Useful for caching purposes
// Returns zero time if the file is not in the bundle
func ModTime(file string) (t time.Time) {
	if f, ok := staticFiles[file]; ok {
		t = f.mtime
	}
	return
}

// Hash returns the hex-encoded SHA256 hash of the original file
// Used for the Etag, and useful for caching
// Returns an empty string if the file is not in the bundle
func Hash(file string) (s string) {
	if f, ok := staticFiles[file]; ok {
		s = f.hash
	}
	return
}
