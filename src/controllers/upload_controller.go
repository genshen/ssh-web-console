package controllers

import (
	"log"
	"net/http"
	"github.com/genshen/webConsole/src/utils"
	"github.com/genshen/webConsole/src/models"
)

type FileUpload struct{}


func (c FileUpload) ShouldClearSessionAfterExec() bool{
	return false
}

func (f FileUpload) ServeAfterAuthenticated(w http.ResponseWriter, r *http.Request, claims *utils.Claims, session *utils.Session) {
	//file, header, err := this.GetFile("file")

	r.ParseMultipartForm(32 << 20)
	file, header, err := r.FormFile("file")
	if err != nil {
		log.Println("Error: getfile err ", err)
		utils.Abort(w, "error", 503)
		return
	}
	defer file.Close()

	user := session.Value.(models.UserInfo)
	if err := utils.UploadFile(utils.SftpNode{Host: user.Host, Port: user.Port}, user.Username, user.Password, file, header); err != nil {
		log.Println("Error: sftp error:", err)
		utils.Abort(w, "message", 503)
	} else {
		w.Write([]byte("sussess"))
	}
}
