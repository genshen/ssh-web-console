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

func Register() {
	// serve static files
	// In dev mode, resource files (for example /static/*) and views(fro example /index.html) are served separately.
	// In production mode, resource files and views are served by statikFS (for example /*).
	if utils.Config.Site.RunMode == RunModeDev {
		if utils.Config.Dev.StaticPrefix == utils.Config.Dev.ViewsPrefix {
			log.Fatal(`static prefix and views prefix can not be the same, check your config.`)
			return
		}
		// server resource files
		if utils.Config.Dev.StaticRedirect == "" {
			// serve locally
			localFile := justFilesFilesystem{http.Dir(utils.Config.Dev.StaticDir)}
			http.Handle(utils.Config.Dev.StaticPrefix, http.StripPrefix(utils.Config.Dev.StaticPrefix, http.FileServer(localFile)))
		} else {
			// serve by redirection
			http.HandleFunc(utils.Config.Dev.StaticPrefix, func(writer http.ResponseWriter, req *http.Request) {
				http.Redirect(writer, req, utils.Config.Dev.StaticRedirect+req.URL.Path, http.StatusMovedPermanently)
			})
		}
		// serve views files.
		utils.MemStatic(utils.Config.Dev.ViewsDir)
		http.HandleFunc(utils.Config.Dev.ViewsPrefix, func(w http.ResponseWriter, r *http.Request) {
			utils.ServeHTTP(w, r) // server soft static files.
		})
	} else {
		statikFS, err := fs.New()
		if err != nil {
			log.Fatal(err)
		}
		http.Handle(utils.Config.Prod.StaticPrefix, http.StripPrefix(utils.Config.Prod.StaticPrefix, http.FileServer(statikFS)))
	}

	bct := utils.Config.SSH.BufferCheckerCycleTime
	// api
	http.HandleFunc("/api/signin", controllers.SignIn)
	http.HandleFunc("/api/sftp/upload", controllers.AuthPreChecker(files.FileUpload{}))
	http.HandleFunc("/api/sftp/ls", controllers.AuthPreChecker(files.List{}))
	http.HandleFunc("/api/sftp/dl", controllers.AuthPreChecker(files.Download{}))
	http.HandleFunc("/ws/ssh", controllers.AuthPreChecker(controllers.NewSSHWSHandle(bct)))
	http.HandleFunc("/ws/sftp", controllers.AuthPreChecker(files.SftpEstablish{}))
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
