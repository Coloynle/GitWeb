package controllers

import (
	"GitWeb/models"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/orm"
	"io/ioutil"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type List struct {
	Menu
	FlushDir bool
}

func (this *List) Prepare() {
	this.dirG = DirList{}
	this.ignore = IgnoreList{}

	dirGSession := this.GetSession(beego.AppConfig.String("session::GitProjectDirectory"))
	ignoreSession := this.GetSession(beego.AppConfig.String("session::IgnoreDirectory"))
	countSession := this.GetSession(beego.AppConfig.String("session::GitProjectCount"))
	flagSession := this.GetSession(beego.AppConfig.String("session::flushflag"))

	if dirGSession != nil {
		// 初始化git项目目录map
		json.Unmarshal(dirGSession.([]byte), &this.dirG)
	}
	if ignoreSession != nil {
		// 初始化忽略目录map
		json.Unmarshal(ignoreSession.([]byte), &this.ignore)
	}
	if countSession != nil {
		// 初始化list总数
		this.Count = countSession.(int)
	}
	if flagSession != nil {
		this.FlushDir = flagSession.(bool)
	} else {
		this.FlushDir = true
	}

	this.WorkPath = beego.AppConfig.String("list::workpath")
	this.Limit, _ = beego.AppConfig.Int("list::limit")

	// 当前请求URL
	this.Data["url"] = this.Ctx.Input.GetData("uri")
}

func (this *List) Test() {
	cfg, _ := config.NewConfig("ini", "conf/config.conf")
	cfg.Set("list::limit", "10")
	cfg.SaveConfigFile("conf/user.conf")

	this.Ctx.WriteString("123")
}

// 项目列表主页 list/get
func (this *List) Get() {
	if 0 == this.Count || this.FlushDir == true {
		this.Count = 0
		this.setDirG()
	}
	this.getIgnore()
	dir := this.getDirPage()
	pageCount := len(dir)                                   // 缓存总页数
	count := (pageCount-1)*this.Limit + len(dir[pageCount]) // 缓存总数

	if count < 0 {
		count = 0
	}

	// 项目列表
	this.Data["dir"] = dir
	// 项目条数
	this.Data["count"] = count
	// 每页条数
	this.Data["limit"] = this.Limit
	// 工作路径
	this.Data["workPath"] = this.WorkPath

	this.TplName = "index/list.html"
}

// 设置Git项目数组
func (this *List) setDirG() bool {
	// 初始化 i.dirG
	this.dirG = DirList{}
	this.getDir(this.WorkPath)
	if len(this.dirG) == 0 {
		return false
	} else {
		this.FlushDir = false
		jsonStr, _ := json.Marshal(this.dirG) // 列表Json化
		this.SetSession(beego.AppConfig.String("session::GitProjectDirectory"), jsonStr)
		this.SetSession(beego.AppConfig.String("session::GitProjectCount"), this.Count)
		this.SetSession(beego.AppConfig.String("session::flushFlag"), this.FlushDir)
		return true
	}
}

// 获取路径下所有Git项目 (递归)
func (this *List) getDir(path string) bool {
	files, _ := ioutil.ReadDir(path)
	for _, f := range files {
		if f.IsDir() {
			if f.Name() == ".git" {
				return true
			}
		}
	}
	for _, f := range files {
		if f.IsDir() {
			str := f.Name()
			// 过滤.开头的文件
			matched, err := regexp.MatchString("^\\.\\S*", str)
			if err == nil && !matched {
				dirPath := path + "/" + f.Name()
				GitTF := this.getDir(dirPath)
				if GitTF == true {
					var command CommandController
					BranchList := command.GetBranchList(dirPath)
					NowBranch := command.GetNowBranch(dirPath)
					this.Count += 1
					this.dirG[this.Count] = Dir{this.Count, dirPath, f.Name(), BranchList, NowBranch}
				}
			}
		}
	}
	return false
}

// 返回真实留下的项目（分页）
func (this *List) getDirPage() map[int]DirList {
	var realDir map[int]DirList
	var temp DirList
	realDir = make(map[int]DirList)
	temp = make(DirList)
	dirList := this.dirG
	page := 1
	count := 0
	max := dirList.Len()
	fmt.Println(max)
	for index := 1; index <= max; index++ {
		fmt.Println("test12321")
		if len(this.ignore) == 0 {
			temp[index] = dirList[index]
			count++
		} else {
			flag := 1
			for _, ignore := range this.ignore {
				// 地址末尾加斜线 避免 ad -> admin 被匹配
				if !(!strings.Contains(dirList[index].Path+"/", ignore.Path+"/") || ignore.Status != 1) {
					flag = 0
				}
			}
			if flag == 1 {
				temp[index] = dirList[index]
				count++
			}
		}
		if count >= this.Limit || index == max {
			realDir[page] = temp
			temp = DirList{}
			page++
			count = 0
		}
	}
	// 判断当只有一页 且第一页为空的时候 直接置空
	if len(realDir) == 1 && len(realDir[1]) == 0 {
		realDir = make(map[int]DirList)
	}
	return realDir
}

// 更新某个路径下的git项目
func (this *List) ResetGit() {
	var result map[string]map[string][]string
	result = make(map[string]map[string][]string)

	// 获取IDS
	getIds := this.GetString("ids")
	idsString := strings.Split(getIds, ",")
	ids := make([]int, len(idsString))
	for i := 0; i < len(idsString); i++ {
		ids[i], _ = strconv.Atoi(idsString[i])
	}
	for _, id := range ids {
		path := this.dirG[id].Path
		comm := CommandController{WorkDir: path}
		result[path] = make(map[string][]string)
		result[path] = comm.updateProject()
	}
	this.Data["json"] = result
	this.ServeJSON()
}

// 返回某个项目分支列表
func (this *List) GetBranchList() {
	id, _ := this.GetInt("id")
	this.Data["json"] = this.dirG[id].Branch
	this.ServeJSON()
}

// 切换分支AJAX
func (this *List) ChangeBranch() {
	id, _ := this.GetInt("id")
	branch := this.GetString("branch")
	path := this.dirG[id].Path
	var command CommandController
	result := command.CheckoutBranch(path, branch)

	// 修改当前项目分支信息
	dir := this.dirG[id]
	dir.NowBranch = command.GetNowBranch(path)
	this.dirG[id] = dir
	if dir.NowBranch == branch {
		result = []string{"1"}
		this.SetSession(beego.AppConfig.String("session::flushFlag"), true) // 更新列表
	}

	this.Data["json"] = result
	this.ServeJSON()
}

// 设置git项目目录
func (this *List) ResetDir() {
	var code int
	var message string
	// 初始化总条数
	this.Count = 0
	if this.setDirG() {
		code = 1
		message = "获取成功"
	} else {
		code = 0
		message = "获取失败或目录下无git项目，请重新设置工作目录"
	}
	this.Data["json"] = this.returnJson(code, message)
	this.ServeJSON()
}

// 忽略列表主页
func (this *List) IgnoreList() {
	// 获取项目列表
	if 0 == this.ignore.Len() {
		this.getIgnore()
	}
	dir := this.getIgnorePage()
	pageCount := len(dir)
	count := (pageCount-1)*this.Limit + len(dir[pageCount])
	if count < 0 {
		count = 0
	}
	// 项目列表
	this.Data["dir"] = dir
	// 项目条数
	this.Data["count"] = count
	// 每页条数
	this.Data["limit"] = this.Limit
	// 工作路径
	this.Data["workPath"] = this.WorkPath

	//test
	this.Data["ignore"] = this.ignore

	this.TplName = "index/ignore.html"
}

// 从数据库获取忽略列表
func (this *List) getIgnore() {
	ignore := models.Ignore{}
	ignores := ignore.GetAll(1000)
	this.ignore = make(IgnoreList)
	for _, value := range ignores {
		this.ignore[value.Id] = Ignore{value.Id, value.Name, value.Path, value.Status}
	}
}

// 返回忽略路径（分页）
func (this *List) getIgnorePage() map[int]IgnoreList {
	var ignore map[int]IgnoreList
	var temp IgnoreList
	var value Ignore
	ignore = make(map[int]IgnoreList)
	temp = make(IgnoreList)
	ignoreList := this.ignore
	page := 1
	count := 0
	num := 0
	max := ignoreList.Len()

	// map 有序化
	var ids []int
	for id := range ignoreList{
		ids = append(ids,id)
	}
	sort.Ints(ids)

	for _,id := range ids{
		 value = ignoreList[id]
		temp[value.Id] = value
		count++
		num++
		if count >= this.Limit || num == max {
			ignore[page] = temp
			temp = IgnoreList{}
			page++
			count = 0
		}
	}

	return ignore
}

// 设置忽略目录
func (this *List) SetIgnore() {
	var name string
	var path string
	var result []int64

	// 获取IDS
	getIds := this.GetString("ids")
	status, _ := this.GetInt("status")
	idsString := strings.Split(getIds, ",")
	ids := make([]int, len(idsString))
	for i := 0; i < len(idsString); i++ {
		ids[i], _ = strconv.Atoi(idsString[i])
	}
	for _, id := range ids {
		path = this.dirG[id].Path
		name = this.dirG[id].Name
		id := this.setIgnore(name, path, status)
		result = append(result, id)
	}
	this.Data["json"] = result
	this.ServeJSON()
}

// 设置忽略列表状态（添加方法）
func (this *List) setIgnore(name string, path string, status int) int64 {
	ignore := models.Ignore{}
	ignore.Path = path
	ignore.Status = status
	ignore.Name = name
	id, _ := ignore.Insert(ignore)
	return id
}

// 更新忽略状态
func (this *List) UpdateIgnoreStatus() {
	ids := this.GetString("ids")
	status, _ := this.GetInt("status")
	bool := this.updateIgnore(ids, status)
	var code int
	var message string
	if bool {
		code = 1
		message = "更新成功"
	} else {
		code = 0
		message = "更新失败"
	}
	result := this.returnJson(code, message)
	this.Data["json"] = result
	this.ServeJSON()
}

// 设置忽略列表状态（改变状态方法）
func (this *List) updateIgnore(id string, status int) bool {
	idsString := strings.Split(id, ",")
	ids := make([]int64, len(idsString))
	for i := 0; i < len(idsString); i++ {
		ids[i], _ = strconv.ParseInt(idsString[i], 10, 64)
	}
	ignore := models.Ignore{}
	ret, _ := ignore.Update(orm.Params{"id__in": ids}, orm.Params{"status": status})
	if ret != -1 {
		return true
	} else {
		return false
	}
}

// 删除忽略条目
func (this *List) DeleteIgnore() {
	ids := this.GetString("ids")
	bool := this.deleteIgnore(ids)
	var code int
	var message string
	if bool {
		code = 1
		message = "更新成功"
	} else {
		code = 0
		message = "更新失败"
	}
	result := this.returnJson(code, message)
	this.Data["json"] = result
	this.ServeJSON()
}

// 设置忽略列表状态（改变状态方法）
func (this *List) deleteIgnore(id string) bool {
	ids := this.stringToInt64(id, ",")
	ignore := models.Ignore{}
	ret, _ := ignore.Delete(orm.Params{"id__in": ids})
	if ret != -1 {
		return true
	} else {
		return false
	}
}
