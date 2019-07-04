package service

import "github.com/astaxie/beego/config/env"

const NULL = ""

func local() {
	env.Get("user", NULL)
	env.Get("email", NULL)
	env.Get("image_prefix", NULL)
	env.Get("", NULL)
}

func file() {

}
