package controller

import (
	"console/service"
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"github.com/gorilla/websocket"
	"io"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

func init() {
	go handleMessages()
}

//docker控制器
type DockerController struct {
	BaseController
}

//查询镜像信息
type QueryImageInfo struct {
	Q string `json:"q"`
}

//socket消息
type Message struct {
	Message string `json:"message"`
}

var (
	upgrader  = websocket.Upgrader{}
	clients   = make(map[*websocket.Conn]bool)
	broadcast = make(chan Message)
)

//对象转string方法
func (it *QueryImageInfo) String() string {
	return "q: " + it.Q
}

//根据参数查询镜像信息
func (it *DockerController) QueryImage() {

}

//查询所有镜像
func (it *DockerController) ListImage() {
	defer DeferFunc(it.Ctx.Output)
	list := service.DockerImageList()
	Return(it.Ctx.Output, list, nil)
}

//删除镜像
func (it *DockerController) DeleteImage() {
	defer DeferFunc(it.Ctx.Output)
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
	if strings.TrimSpace(deleteImage.Id) == "" {
		panic("删除镜像不能为空")
	}
	list := service.DockerImageDelete(deleteImage.Id)
	Return(it.Ctx.Output, list, nil)
}

//推送镜像
func (it *DockerController) PushImage() {
	defer DeferFunc(it.Ctx.Output)
	image := it.Ctx.Request.Form.Get("image")
	if strings.TrimSpace(image) == "" {
		panic("镜像不能为空")
	}
	flag := make(chan struct{})
	defer close(flag)
	go service.DockerImagePush(image, flag)
	<-flag
	Return(it.Ctx.Output, nil, nil)
}

//build镜像
func (it *DockerController) BuildImage() {
	defer DeferFunc(it.Ctx.Output)
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
	Return(it.Ctx.Output, nil, nil)
}

func (it *DockerController) ImageBuildWebsocketRegister() {
	ws, err := upgrader.Upgrade(it.Ctx.ResponseWriter, it.Ctx.Request, nil)
	if err != nil {
		log.Fatal(err)
	}
	clients[ws] = true
}

func handleMessages() {
	for {
		msg := <-broadcast
		logs.Info("clients len ", len(clients))
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("client.WriteJSON error: %v", err)
				e := client.Close()
				log.Printf("client.WriteJSON error: %v", e)
				delete(clients, client)
			}
		}
	}
}
