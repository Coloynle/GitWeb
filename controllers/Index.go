package controllers

import (
	"github.com/astaxie/beego"
)

type Index struct {
	beego.Controller
}

func (this *Index) Prepare() {
	// 当前请求URL
	this.Data["url"] = this.Ctx.Input.GetData("uri")
}

func (this *Index) Get(){

	this.TplName = "index/index.html"
}

