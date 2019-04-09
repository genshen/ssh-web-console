package routers

import (
	"github.com/genshen/ssh-web-console/src/controllers"
	"github.com/genshen/ssh-web-console/src/controllers/files"
	"github.com/genshen/ssh-web-console/src/utils"
	_ "github.com/genshen/ssh-web-console/statik"
	"github.com/rakyll/statik/fs"
	"log"
	"net/http"
	"os"
)

const (
	RunModeDev  = "dev"
	RunModeProd = "prod"
)

func init() {
	// static
	if utils.Config.Site.RunMode == RunModeDev {
		http.HandleFunc("/static/", func(writer http.ResponseWriter, req *http.Request) {
			http.Redirect(writer, req, "localhost:8080"+req.URL.Path, http.StatusMovedPermanently)
		})
	} else {
		//fs := justFilesFilesystem{http.Dir(utils.Config.Site.HardStaticDir)}
		//http.Handle(utils.Config.Site.StaticPrefix, http.StripPrefix(utils.Config.Site.StaticPrefix, http.FileServer(fs)))
		statikFS, err := fs.New()
		if err != nil {
			log.Fatal(err)
		}
		http.Handle(utils.Config.Site.StaticPrefix, http.StripPrefix(utils.Config.Site.StaticPrefix, http.FileServer(statikFS)))
	}

	// api
	http.HandleFunc("/api/signin", controllers.SignIn)
	http.HandleFunc("/api/sftp/upload", controllers.AuthPreChecker(files.FileUpload{}))
	http.HandleFunc("/api/sftp/ls", controllers.AuthPreChecker(files.List{}))
	http.HandleFunc("/api/sftp/dl", controllers.AuthPreChecker(files.Download{}))
	http.HandleFunc("/ws/ssh", controllers.AuthPreChecker(controllers.NewSSHWSHandle()))
	http.HandleFunc("/ws/sftp", controllers.AuthPreChecker(files.SftpEstablish{}))
}

func Run() {
	http.ListenAndServe(utils.Config.Site.ListenAddr, nil)
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
