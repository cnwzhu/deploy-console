package domain

import "fmt"

type ReturnMsg struct {
	Code  int32       `json:"code"`
	Data  interface{} `json:"data"`
	State string      `json:"state"`
}

func (it *ReturnMsg) String() string {
	return fmt.Sprintf("code: %s,  state: %s,  data: %s", it.Code, it.State, it.Data)
}

