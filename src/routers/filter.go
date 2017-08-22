package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func intiFilter() {
	if beego.BConfig.RunMode == beego.DEV{
		beego.InsertFilter("/dist/*",beego.BeforeRouter, func(ctx *context.Context){  //todo for debug
			ctx.Redirect(302,"http://localhost:8080"+ctx.Request.RequestURI)
		})
	}
}
