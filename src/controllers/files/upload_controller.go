package files

import (
	"github.com/genshen/ssh-web-console/src/utils"
	"github.com/pkg/sftp"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path"
)

type FileUpload struct{}

func (f FileUpload) ShouldClearSessionAfterExec() bool {
	return false
}

func (f FileUpload) ServeAfterAuthenticated(w http.ResponseWriter, r *http.Request, claims *utils.Claims, session utils.Session) {
	cid := r.URL.Query().Get("cid") // get connection id.
	if sftpClient := ForkSftpClient(cid); sftpClient == nil {
		utils.Abort(w, "error: lost sftp connection.", 400)
		log.Println("Error: lost sftp connection.")
		return
	} else {
		//file, header, err := this.GetFile("file")
		r.ParseMultipartForm(32 << 20)
		file, header, err := r.FormFile("file")
		if err != nil {
			log.Println("Error: getfile err ", err)
			utils.Abort(w, "error", 503)
			return
		}
		defer file.Close()

		if err := UploadFile(sftpClient, file, header); err != nil {
			log.Println("Error: sftp error:", err)
			utils.Abort(w, "message", 503)
		} else {
			w.Write([]byte("success"))
		}
	}
}

// upload file to server via sftp.
func UploadFile(client *sftp.Client, srcFile multipart.File, header *multipart.FileHeader) error {
	var fullPath string
	if wd, err := client.Getwd(); err == nil {
		fullPath = path.Join(wd, "/tmp/")
		if _, err := client.Stat(fullPath); err != nil {
			if os.IsNotExist(err) {
				if err := client.Mkdir(fullPath); err != nil {
					return err
				}
			} else {
				return err
			}
		}
	} else {
		return err
	}

	dstFile, err := client.Create(path.Join(fullPath, header.Filename))
	if err != nil {
		return err
	}
	defer srcFile.Close()
	defer dstFile.Close()

	_, err = dstFile.ReadFrom(srcFile)
	if err != nil {
		return err
	}
	return nil
}
