package controller

type RegisterController struct {
	BaseController
}

func (it *RegisterController) QueryRegisterPrefix() {
	defer DeferFunc(it.Ctx.Output)
	rt := []string{"192.168.31.188/test"}
	Return(it.Ctx.Output, rt, nil)
}
