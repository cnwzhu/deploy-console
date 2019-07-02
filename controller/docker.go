package controller

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"io/ioutil"
	"log"
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
