package controller

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
)

type ReturnMsg struct {
	Code  int32       `json:"code"`
	Data  interface{} `json:"data"`
	State string      `json:"state"`
}

//docker控制器
type BaseController struct {
	beego.Controller
}

func DeferFunc(o *context.BeegoOutput) {
	if e := recover(); e != nil {
		logs.Error(" 错误 %s\r\n", e)
		Return(o, nil, &e)
	}
}

func Return(o *context.BeegoOutput, data interface{}, err *interface{}) {
	o.ContentType("application/json")
	if err != nil {
		bytes, e := json.Marshal(ReturnMsg{
			Data:  data,
			Code:  500,
			State: "fail",
		})
		if e != nil {
			panic(e)
		}
		_ = o.Body(bytes)
	} else {
		bytes, e := json.Marshal(ReturnMsg{
			Data:  data,
			Code:  200,
			State: "ok",
		})
		if e != nil {
			panic(e)
		}
		_ = o.Body(bytes)
	}
}
