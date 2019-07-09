package router

import (
	"console/controller"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/plugins/cors"
)

func init() {
	//添加跨域允许
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		AllowCredentials: true,
	}))
	//查询
	beego.Router("/docker/image/query", &controller.DockerController{}, "GET:QueryImage")
	//websocket注册
	beego.Router("/docker/image/ws/register", &controller.DockerController{}, "GET:ImageBuildWebsocketRegister")
	//查询所有镜像
	beego.Router("/docker/image/list", &controller.DockerController{}, "GET:ListImage")
	//构建镜像
	beego.Router("/docker/image/build", &controller.DockerController{}, "POST:BuildImage")
	//推送镜像
	beego.Router("/docker/image/push", &controller.DockerController{}, "POST:PushImage")
	//删除镜像
	beego.Router("/docker/image/delete", &controller.DockerController{}, "DELETE:DeleteImage")
}
