package main

import (
	"GitWeb/controllers"
	"GitWeb/models"
	_ "GitWeb/routers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func init()  {
	mysqlConn := controllers.MySqlConn{}
	orm.RegisterDriver("mysql",orm.DRMySQL)
	orm.RegisterDataBase("default","mysql",mysqlConn.GetDataSource())
	// orm.RegisterDataBase("default","mysql","root:@/golang?charset=utf8")
	orm.RegisterModel(new(models.Ignore))
	orm.RunSyncdb("default",false,true)
}

func main() {
	// beego.LoadAppConfig("ini", "conf/const.conf")
	// beego.LoadAppConfig("ini", "conf/config.conf")

	cfg,_ := config.NewConfig("ini","conf/app.conf")
	cfg.SaveConfigFile("conf/user.conf")

	var Uri = func(ctx *context.Context) {
		ctx.Input.SetData("uri",ctx.Request.RequestURI)
	}

	beego.InsertFilter("/*",beego.BeforeExec,Uri)

	beego.Run()
}

