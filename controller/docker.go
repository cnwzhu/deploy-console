package controller

import (
	"console/service"
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"github.com/gorilla/websocket"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
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
	ug = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	buildWsClients   = make(map[*websocket.Conn]bool)
	buildWsBroadcast = make(chan Message)
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
	go service.DockerImagePush(image)
	buildWsBroadcast <- Message{"push完成"}
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
	go inner(buildWsBroadcast, &file, info, header)
	Return(it.Ctx.Output, nil, nil)
}

func inner(ch chan Message, file *multipart.File, info *service.ImageSimpleBuildInfo, header *multipart.FileHeader) {
	r := (*file).(io.Reader)
	service.DockerImageBuild(&r, info, header)
	ch <- Message{"构建完成"}
	logs.Info("full name: " + info.Prefix + "/" + info.Name + ":" + info.Version)
	service.DockerImagePush(info.Prefix + "/" + info.Name + ":" + info.Version)
	ch <- Message{"push完成"}
	defer (*file).Close()
}

func (it *DockerController) ImageBuildWebsocketRegister() {
	defer DeferFunc(it.Ctx.Output)
	ws, err := ug.Upgrade(it.Ctx.ResponseWriter, it.Ctx.Request, nil)
	if err != nil {
		log.Fatal(err)
	}
	buildWsClients[ws] = true
}

func handleMessages() {
	for {
		msg := <-buildWsBroadcast
		logs.Info("clients len ", len(buildWsClients))
		logs.Info(msg.Message)
		for client := range buildWsClients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("client.WriteJSON error: %v", err)
				e := client.Close()
				log.Printf("client.WriteJSON error: %v", e)
				delete(buildWsClients, client)
			}
		}
	}
}

func (it *DockerController) Test() {
	buildWsBroadcast <- Message{"test ok"}
}
