package files

import (
	"github.com/genshen/ssh-web-console/src/utils"
	"io"
	"log"
	"net/http"
	"path"
)

type Download struct{}

func (d Download) ShouldClearSessionAfterExec() bool {
	return false
}

func (d Download) ServeAfterAuthenticated(w http.ResponseWriter, r *http.Request, claims *utils.Claims, session utils.Session) {
	cid := r.URL.Query().Get("cid") // get connection id.
	if client := utils.ForkSftpClient(cid); client == nil {
		utils.Abort(w, "error: lost sftp connection.", 400)
		log.Println("Error: lost sftp connection.")
		return
	} else {
		if wd, err := client.Getwd(); err == nil {
			relativePath := r.URL.Query().Get("path") // get path.
			fullPath := path.Join(wd, relativePath)
			if fileInfo, err := client.Stat(fullPath); err == nil && !fileInfo.IsDir() {
				if file, err := client.Open(fullPath); err == nil {
					defer file.Close()
					w.Header().Add("Content-Disposition", "attachment;filename="+fileInfo.Name())
					w.Header().Add("Content-Type", "application/octet-stream")
					io.Copy(w, file)
					return
				}
			}
		}
		utils.Abort(w, "no such file", 400)
		return
	}
}
