package routers

import (
	"os"
	"net/http"
	"github.com/genshen/webConsole/src/controllers"
	"github.com/genshen/webConsole/src/utils"
)

func init() {
	if utils.Config.Site.RunMode == "dev" {
		http.HandleFunc("/static/", func(writer http.ResponseWriter, req *http.Request) {
			http.Redirect(writer, req, "localhost:8080"+req.URL.Path, http.StatusMovedPermanently)
		})
	} else {
		fs := justFilesFilesystem{http.Dir(utils.Config.Site.StaticDir)}
		//http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(fs)))
		http.Handle("/static/", http.FileServer(fs))
	}

	http.HandleFunc("/", controllers.Get)
	http.HandleFunc("/signin", controllers.SignIn)
	http.HandleFunc("/ssh/uploadfile", controllers.AuthPreChecker(controllers.FileUpload{}))
	http.HandleFunc("/ws/ssh", controllers.AuthPreChecker(controllers.SSHWebSocketHandle{}))
}

func Run() {
	http.ListenAndServe(utils.Config.Site.ListenAddr,nil)
}

/*
* disable directory index, code from https://groups.google.com/forum/#!topic/golang-nuts/bStLPdIVM6w
 */
type justFilesFilesystem struct {
	fs http.FileSystem
}

func (fs justFilesFilesystem) Open(name string) (http.File, error) {
	f, err := fs.fs.Open(name)
	if err != nil {
		return nil, err
	}
	return neuteredReaddirFile{f}, nil
}

type neuteredReaddirFile struct {
	http.File
}

func (f neuteredReaddirFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
}
