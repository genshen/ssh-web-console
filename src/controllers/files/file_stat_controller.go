package files

import (
	"net/http"
	"github.com/genshen/webConsole/src/utils"
)

type FileStat struct{}

func (f FileStat) ShouldClearSessionAfterExec() bool {
	return false
}

func (f FileStat) ServeAfterAuthenticated(w http.ResponseWriter, r *http.Request, claims *utils.Claims, session *utils.Session) {

}
