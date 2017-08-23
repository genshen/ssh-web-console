package routers

import (
	"github.com/genshen/webConsole/src/controllers"
	"github.com/astaxie/beego"
)

func init() {
	intiFilter()

    beego.Router("/", &controllers.MainController{})
    beego.Router("/signin", &controllers.MainController{},"post:SignIn")
    beego.Router("/ssh/uploadfile", &controllers.UploadController{},"post:UploadFile")

	beego.Router("/ws/ssh", &controllers.WebSocketController{},"get:SSHWebSocketHandle")
}
