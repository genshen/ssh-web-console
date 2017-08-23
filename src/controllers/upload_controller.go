package controllers

import (
	"github.com/astaxie/beego"
	"github.com/genshen/webConsole/src/models"
	"github.com/genshen/webConsole/src/utils"
)

type UploadController struct {
	BaseController
}

func (this *UploadController) UploadFile() {
	file, header, err := this.GetFile("file")
	if err != nil {
		beego.Error("getfile err ", err)
		this.Abort("503")
	} else {
		v := this.GetSession("userinfo")
		if v == nil {
			beego.Error("Cannot get Session data:", err)
			this.Abort("503")
		} else {
			user := v.(models.UserInfo)
			if err := utils.UploadFile(user, file, header); err != nil {
				beego.Error("sftp error:", err)
				this.Abort("503")
			} else {
				this.Ctx.WriteString("sussess")
			}
		}
	}
}
