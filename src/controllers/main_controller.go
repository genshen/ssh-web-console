package controllers

import (
	"github.com/genshen/webConsole/src/models"
	"github.com/genshen/webConsole/src/utils"
)

type MainController struct {
	BaseController
}

func (this *MainController) Get() {
	this.TplName = "index.html"
}

const (
	SIGN_IN_FORM_TYPE_ERROR_VALID    = iota
	SIGN_IN_FORM_TYPE_ERROR_PASSWORD
	SIGN_IN_FORM_TYPE_ERROR_TEST
)

func (this *MainController) SignIn() {
	var err error;
	userinfo := models.UserInfo{}
	userinfo.Host = this.GetString("host")
	userinfo.Port, err = this.GetInt("port", 22)
	userinfo.Username = this.GetString("username")
	userinfo.Password = this.GetString("passwd")
	if err == nil && userinfo.Host != "" && userinfo.Username != "" {
		//try to login ssh account
		ssh := utils.SSH{}
		ssh.Node.Host = userinfo.Host
		ssh.Node.Port = userinfo.Port
		_, err := ssh.Connect(userinfo.Username, userinfo.Password)
		if err != nil {
			errUnmarshal := models.SignInFormValid{HasError: true, Message: SIGN_IN_FORM_TYPE_ERROR_PASSWORD}
			this.Data["json"] = &errUnmarshal
		} else {
			defer ssh.Close()
			// create session
			if session, err := ssh.Client.NewSession(); err == nil {
				if err := session.Run("whoami"); err == nil {
					this.SetSession("userinfo", userinfo)
					errUnmarshal := models.SignInFormValid{HasError: false}
					this.Data["json"] = &errUnmarshal
					this.ServeJSON()
					return
				}
			}
			errUnmarshal := models.SignInFormValid{HasError: true, Message: SIGN_IN_FORM_TYPE_ERROR_TEST}
			this.Data["json"] = &errUnmarshal
		}
	} else {
		errUnmarshal := models.SignInFormValid{HasError: true, Message: SIGN_IN_FORM_TYPE_ERROR_VALID}
		this.Data["json"] = &errUnmarshal
	}
	this.ServeJSON()
}
