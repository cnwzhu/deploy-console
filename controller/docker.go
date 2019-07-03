package controller

import (
	"console/service"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"io"
	"io/ioutil"
	"log"
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
	var info QueryImageInfo
	b := it.Ctx.Request.Body
	defer b.Close()
	all, _ := ioutil.ReadAll(b)
	log.Println(string(all))
	json.Unmarshal(all, &info)
	log.Println(info)
	bytes, _ := json.Marshal(info)
	it.Ctx.Output.Body(bytes)
}

func (it *DockerController) BuildImage()  {
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
	<-flag
	e = it.Ctx.Output.Body([]byte(`{"msg":"ok"}`))
	if e != nil {
		panic(e)
	}
}
