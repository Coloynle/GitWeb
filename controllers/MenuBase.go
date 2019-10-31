package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/config"
	"regexp"
	"strconv"
	"strings"
)

// base MenuController
type Menu struct {
	beego.Controller
	WorkPath string     // 工作目录
	dirG     DirList    // Git项目数组
	ignore   IgnoreList // 忽略目录
	Limit    int        // 每页的大小
	Count    int
}

func (this *Menu) Prepare() {
	this.dirG = make(DirList)
	this.ignore = make(IgnoreList)
}

type Dir struct {
	id        int      `json:"id"`
	Path      string   `json:"path"`
	Name      string   `json:"name"`
	Branch    []string `json:"branch"`
	NowBranch string   `json:"now_branch"`
}
type DirList map[int]Dir

type Ignore struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Path   string `json:"path"`
	Status int    `json:"status"`
}
type IgnoreList map[int]Ignore

type returnStatus struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (d IgnoreList) Len() int { return len(d) }

func (d DirList) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }
func (d DirList) Len() int           { return len(d) }
func (d DirList) Less(i, j int) bool { return d[i].Path < d[j].Name }

type MySqlConn struct {
	NetworkType  string `json:"network_type"`
	Host         string `json:"host"`
	Port         string `json:"port"`
	Name         string `json:"name"`
	Password     string `json:"password"`
	DatabaseName string `json:"database_name"`
}

func (m *MySqlConn) GetDataSource() string {
	m.Name = beego.AppConfig.String("mysql::Name")
	m.Password = beego.AppConfig.String("mysql::Password")
	m.NetworkType = beego.AppConfig.String("mysql::NetworkType")
	m.Host = beego.AppConfig.String("mysql::Host")
	m.Port = beego.AppConfig.String("mysql::Port")
	m.DatabaseName = beego.AppConfig.String("mysql::DatabaseName")
	m.Name = beego.AppConfig.String("mysql::Name")
	return m.Name + ":" + m.Password + "@" + m.NetworkType + "(" + m.Host + ":" + m.Port + ")/" + m.DatabaseName + "?charset=utf8&loc=Asia%2FShanghai&parseTime=true"
}

func (this *Menu) SetWorkPath() {
	str := this.GetString("path")

	matched, err, str := VerifyPath(str)
	if matched != true || err != nil {
		this.Data["json"] = this.returnJson(0, "路径错误，请修改路径为正确的格式")
		this.ServeJSON()
		return
	}
	this.WorkPath = str
	setConfig("list::workpath", this.WorkPath)
	result := this.returnJson(1, this.WorkPath)
	this.Data["json"] = result
	this.ServeJSON()
	return
}

// 对输入路径进行格式验证并且格式化
func VerifyPath(str string) (matched bool, err error, path string) {
	// 对输入路径进行格式验证并且格式化
	str = strings.Replace(str, "\\", "/", -1)
	regstr := "\\/{2,}"
	reg, _ := regexp.Compile(regstr)
	str = reg.ReplaceAllString(str, "/")
	last := str[len(str)-1:]
	if last == "/" {
		str = str[0 : len(str)-1]
	}
	matched, err = regexp.MatchString("^[A-Za-z]:(\\/[\u4E00-\u9FA5A-Za-z0-9]+)+$", str)
	return matched, err, str
}

// 修改config即刻生效并且存至user.conf
func setConfig(key string,value string){
	beego.AppConfig.Set(key,value)
	cfg,_ := config.NewConfig("ini","conf/user.conf")
	cfg.Set(key,value)
	cfg.SaveConfigFile("conf/user.conf")
	// 重载配置文件
	beego.LoadAppConfig("ini","conf/app.conf")
}

func (this *Menu) returnJson(code int, message string) *returnStatus {
	result := &returnStatus{
		Code:    code,
		Message: message,
	}
	return result
}

func (this *Menu) stringToInt64(str string,sep string) []int64{
	strs := strings.Split(str,sep)
	arrInt64 := make([]int64,len(strs))
	for i:=0;i< len(strs); i++{
		arrInt64[i],_ = strconv.ParseInt(strs[i], 10, 64)
	}
	return arrInt64
}