package controller

import (
	"console/service"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"io"
	"io/ioutil"
	"strconv"
)

//docker控制器
type DockerController struct {
	beego.Controller
}

//查询镜像信息
type QueryImageInfo struct {
	Q string `json:"q"`
}

//对象转string方法
func (it *QueryImageInfo) String() string {
	return "q: " + it.Q
}

//根据参数查询镜像信息
func (it *DockerController) QueryImage() {

}

//查询所有镜像
func (it *DockerController) ListImage() {
	defer func() {
		if e := recover(); e != nil {
			logs.Error(" 错误 %s\r\n", e)
		}
	}()
	list := service.DockerImageList()
	e := it.Ctx.Output.Body(list)
	if e != nil {
		panic(e)
	}
}

//删除镜像
func (it *DockerController) DeleteImage() {
	defer func() {
		if e := recover(); e != nil {
			logs.Error(" 错误 %s\r\n", e)
		}
	}()
	body := it.Ctx.Request.Body
	defer body.Close()
	readAll, e := ioutil.ReadAll(body)
	if e != nil {
		panic(e)
	}
	var deleteImage = &struct {
		Id string
	}{}
	e = json.Unmarshal(readAll, deleteImage)
	if e != nil {
		panic(e)
	}
	list := service.DockerImageDelete(deleteImage.Id)
	e = it.Ctx.Output.Body(list)
	if e != nil {
		panic(e)
	}
}

//推送镜像
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

//build镜像
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
	logs.Info("full name: " + info.Prefix + "/" + info.Name + ":" + info.Version)
	go service.DockerImagePush(info.Prefix+"/"+info.Name+":"+info.Version, flag)
	<-flag
	logs.Info("push结束")
	e = it.Ctx.Output.Body([]byte(`{"msg":"ok"}`))
	if e != nil {
		panic(e)
	}
}
