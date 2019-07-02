package router

import (
	"console/controller"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/up", &controller.FileController{}, "POST:Upload")
}
