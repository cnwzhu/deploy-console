package controller

import (
	"console/service"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"io"
	"strconv"
)

type DockerController struct {
	beego.Controller
}

type QueryImageInfo struct {
	Q string `json:"q"`
}

func (it *QueryImageInfo) String() string {
	return "q: " + it.Q
}

func (it *DockerController) QueryImage() {

}

func (it *DockerController) PushImage() {
	defer func() {
		if e := recover(); e != nil {
			logs.Error(" 错误 %s\r\n", e)
		}
	}()
	flag := make(chan struct{})
	defer close(flag)
	service.DockerImagePush("", flag)
	<-flag
}

func (it *DockerController) BuildImage() {
	defer func() {
		if e := recover(); e != nil {
			logs.Error(" 错误 %s\r\n", e)
		}
	}()
	file, header, e := it.GetFile("image_file")
	if e != nil {
		panic(e)
	}
	i, e := strconv.Atoi(it.Ctx.Request.Form.Get("type"))
	info := &service.ImageSimpleBuildInfo{
		Name:    it.Ctx.Request.Form.Get("image_name"),
		Version: it.Ctx.Request.Form.Get("image_version"),
		Prefix:  it.Ctx.Request.Form.Get("image_prefix"),
		Type:    i,
	}
	defer file.Close()
	r := file.(io.Reader)
	flag := make(chan struct{})
	defer close(flag)
	go service.DockerImageBuild(&r, info, flag, header)
	logs.Info("开始build")
	<-flag
	logs.Info("build结束")
	logs.Info("开始push")
	logs.Info("full name: "+info.Prefix+"/"+info.Name+":"+info.Version)
	go service.DockerImagePush(info.Prefix+"/"+info.Name+":"+info.Version, flag)
	<-flag
	logs.Info("push结束")
	e = it.Ctx.Output.Body([]byte(`{"msg":"ok"}`))
	if e != nil {
		panic(e)
	}
}
