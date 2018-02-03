package controllers

import (
	"log"
	"strings"
	"net/http"
	"github.com/genshen/webConsole/src/utils"
)

type AfterAuthenticated interface {
	// make sure token and session is not nil.
	ServeAfterAuthenticated(w http.ResponseWriter, r *http.Request, token *utils.Claims, session utils.Session)
	ShouldClearSessionAfterExec() bool
}

func AuthPreChecker(i AfterAuthenticated) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var authHead = r.Header.Get("Authorization")
		var token string
		if authHead != "" {
			lIndex := strings.LastIndex(authHead, " ")
			if lIndex < 0 || lIndex+1 >= len(authHead) {
				utils.Abort(w, "invalid token", 400)
				log.Println("Error: invalid token", 400)
				return
			} else {
				token = authHead[lIndex+1:]
			}
		} else {
			if token = r.URL.Query().Get(utils.Config.Jwt.QueryTokenKey); token == "" {
				utils.Abort(w, "invalid token", 400)
				log.Println("Error: invalid token", 400)
				return
			} // else token != "", then passed and go on running
		}

		if claims, err := utils.JwtVerify(token); err != nil {
			http.Error(w, "invalid token", 400)
			log.Println("Error: Cannot setup WebSocket connection:", err)
		} else { // check passed.
			// check session.
			if session, ok := utils.SessionStorage.Get(token); !ok { // make a session copy.
				utils.Abort(w, "Error: Cannot get Session data:", 400)
				log.Println("Error: Cannot get Session data for token", token)
			} else {
				if i.ShouldClearSessionAfterExec() {
					defer utils.SessionStorage.Delete(token)
					i.ServeAfterAuthenticated(w, r, claims, session)
				}else{
					i.ServeAfterAuthenticated(w, r, claims, session)
				}
			}
		}
	}
}
