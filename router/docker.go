package router

import (
	"console/controller"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/docker/query", &controller.DockerController{}, "GET:QueryImage")
	beego.Router("/docker/build", &controller.DockerController{}, "POST:BuildImage")
}
