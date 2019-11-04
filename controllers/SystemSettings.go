package controllers

import (
	"github.com/astaxie/beego"
)

type System struct {
	beego.Controller
	configKey []string
	// conf      map[int]Config
	conf []Config
}

type Config struct {
	Key   string
	Value string
}

func (this *System) initConfigKeys() {
	this.configKey = []string{
		"appname",
		"session::GitProjectDirectory",
		"session::IgnoreDirectory",
		"session::GitProjectCount",
		"session::flushflag",
		"list::workpath",
		"list::limit",
		"mysql::name",
		"mysql::password",
		"mysql::databasename",
		"mysql::networktype",
		"mysql::host",
		"mysql::port",
	}
}

func (this *System) Prepare() {
	// 初始化配置列表
	this.initConfigKeys()
	// 当前请求URL
	this.Data["url"] = this.Ctx.Input.GetData("uri")
}

func (this *System) Index() {
	this.initConfigKeys()
	for _, v := range this.configKey {
		value := beego.AppConfig.String(v)
		this.conf = append(this.conf, Config{Key: v, Value: value})
	}
	this.Data["setting"] = this.conf
	this.TplName = "system/index.html"
}

func (this *System) SaveConfig() {
	for _, v := range this.configKey {
		value := this.GetString(v)
		this.conf = append(this.conf, Config{Key: v, Value: value})
	}
	for _, val := range this.conf {
		setConfig(val.Key, val.Value)
	}
	redirectUrl := this.URLFor("System.Index")
	this.Ctx.Redirect(302, redirectUrl)
}
