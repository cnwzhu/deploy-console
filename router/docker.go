package router

import (
	"console/controller"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/docker", &controller.DockerController{}, "POST:QueryImage")
}
