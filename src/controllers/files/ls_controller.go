package files

import (
	"net/http"
	"github.com/genshen/webConsole/src/utils"
	"log"
	"path"
	"github.com/genshen/webConsole/src/models"
)

type List struct{}
type Ls struct {
	Name  string `json:"name"`
	Path  string `json:"path"` // including Name
	IsDir bool   `json:"is_dir"`
}

func (f List) ShouldClearSessionAfterExec() bool {
	return false
}

func (f List) ServeAfterAuthenticated(w http.ResponseWriter, r *http.Request, claims *utils.Claims, session utils.Session) {
	response := models.JsonResponse{HasError: true}
	cid := r.URL.Query().Get("cid") // get connection id.
	if client := ForkSftpClient(cid); client == nil {
		utils.Abort(w, "error: lost sftp connection.", 400)
		log.Println("Error: lost sftp connection.")
		return
	} else {
		if wd, err := client.Getwd(); err == nil {
			relativePath := r.URL.Query().Get("path") // get path.
			fullPath := path.Join(wd, relativePath)
			if files, err := client.ReadDir(fullPath); err != nil {
				response.Addition = "no such path"
			} else {
				response.HasError = false
				fileList := make([]Ls,0) // this will not be converted to null if slice is empty.
				for _, file := range files {
					fileList = append(fileList, Ls{Name: file.Name(), IsDir: file.IsDir(), Path:path.Join(relativePath, file.Name())})
				}
				response.Message = fileList
			}
		} else {
			response.Addition = "no such path"
		}
	}
	utils.ServeJSON(w, response)
}
