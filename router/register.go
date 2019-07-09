package router

import (
	"console/controller"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/docker/image/register/prefix/query", &controller.RegisterController{}, "GET:QueryRegisterPrefix")
}
