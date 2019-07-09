package router

import (
	"console/controller"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/plugins/cors"
)

func init() {
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		AllowCredentials: true,
	}))
	beego.Router("/docker/image/query", &controller.DockerController{}, "GET:QueryImage")
	beego.Router("/docker/image/list", &controller.DockerController{}, "GET:ListImage")
	beego.Router("/docker/image/build", &controller.DockerController{}, "POST:BuildImage")
	beego.Router("/docker/image/push", &controller.DockerController{}, "POST:PushImage")
	beego.Router("/docker/image/delete", &controller.DockerController{}, "DELETE:DeleteImage")
}
