package controllers

import (
	"github.com/genshen/ssh-web-console/src/models"
	"github.com/genshen/ssh-web-console/src/utils"
	"golang.org/x/crypto/ssh"
	"net/http"
	"strconv"
)

func SignIn(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method.", 405)
	} else {
		var err error
		var errUnmarshal models.JsonResponse
		err = r.ParseForm()
		if err != nil {
			panic(err)
		}
		userinfo := models.UserInfo{}
		userinfo.Host = r.Form.Get("host")
		port := r.Form.Get("port")
		userinfo.Username = r.Form.Get("username")
		userinfo.Password = r.Form.Get("passwd")

		userinfo.Port, err = strconv.Atoi(port)
		if err != nil {
			userinfo.Port = 22
		}

		if userinfo.Host != "" && userinfo.Username != "" {
			//try to login session account
			session := utils.SSHShellSession{}
			session.Node = utils.NewSSHNode(userinfo.Host, userinfo.Port)
			err := session.Connect(userinfo.Username, ssh.Password(userinfo.Password))
			if err != nil {
				errUnmarshal = models.JsonResponse{HasError: true, Message: models.SIGN_IN_FORM_TYPE_ERROR_PASSWORD}
			} else {
				defer session.Close()
				// create session
				client, err := session.GetClient()
				if err != nil {
					// bad connection.
					return
				}
				if session, err := client.NewSession(); err == nil {
					if err := session.Run("whoami"); err == nil {
						if token, expireUnix, err := utils.JwtNewToken(userinfo.JwtConnection, utils.Config.Jwt.Issuer); err == nil {
							errUnmarshal = models.JsonResponse{HasError: false, Addition: token}
							utils.ServeJSON(w, errUnmarshal)
							utils.SessionStorage.Put(token, expireUnix, userinfo)
							return
						}
					}
				}
				errUnmarshal = models.JsonResponse{HasError: true, Message: models.SIGN_IN_FORM_TYPE_ERROR_TEST}
			}
		} else {
			errUnmarshal = models.JsonResponse{HasError: true, Message: models.SIGN_IN_FORM_TYPE_ERROR_VALID}
		}
		utils.ServeJSON(w, errUnmarshal)
	}
}
