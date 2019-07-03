package router

import (
	"console/controller"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/docker/image/query", &controller.DockerController{}, "GET:QueryImage")
	beego.Router("/docker/image/build", &controller.DockerController{}, "POST:BuildImage")
	beego.Router("/docker/image/push", &controller.DockerController{}, "POST:PushImage")
}
