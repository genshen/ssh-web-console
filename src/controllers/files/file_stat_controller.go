package files

import (
	"github.com/genshen/ssh-web-console/src/utils"
	"net/http"
)

type FileStat struct{}

func (f FileStat) ShouldClearSessionAfterExec() bool {
	return false
}

func (f FileStat) ServeAfterAuthenticated(w http.ResponseWriter, r *http.Request, claims *utils.Claims, session *utils.Session) {

}
