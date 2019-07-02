package controller

import (
	"console/service"
	"fmt"
	"github.com/astaxie/beego"
	"io"
	"strconv"
)

type FileController struct {
	beego.Controller
}

func (it *FileController) Upload() {
	defer func() {
		if e := recover(); e != nil {
			fmt.Printf(" 错误 %s\r\n", e)
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
	<-flag
	e = it.Ctx.Output.Body([]byte(`{"msg":"ok"}`))
	if e != nil {
		panic(e)
	}
}
