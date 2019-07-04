package service

import (
	"encoding/json"
	"testing"
)

func TestNginxBuild(t *testing.T) {

}

func TestSeria(t *testing.T) {
	temp, _ := json.Marshal(struct {
		ImageList interface{}
	}{[]string{"42", "fsd"}})
	t.Log(temp)
}

func TestSelf(t *testing.T) {
	var deleteImage = &struct {
		Id string
	}{}
	s := `{"Id":"fgsffsgsgfffsgsfg"}`
	e := json.Unmarshal([]byte(s), deleteImage)
	if e != nil {
		t.Log(e)
	}
}
