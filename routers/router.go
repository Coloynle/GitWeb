package routers

import (
	"GitWeb/controllers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func init() {
	List := &controllers.List{}
	Menu := &controllers.Menu{}

	ns :=
		beego.NewNamespace("/v1",
			beego.NSCond(func(ctx *context.Context) bool {
				if ua := ctx.Input.UserAgent(); ua != "" {
					return true
				}
				return true
			}),
			beego.NSBefore(),
			beego.NSNamespace("/list",
				beego.NSRouter("/get", List, "get:Get"),
				beego.NSRouter("/reset", List, "get:ResetDir"),
				beego.NSRouter("/test", List, "get:Test"),
			),beego.NSNamespace("/git",
				beego.NSRouter("/update", List, "post:ResetGit"),
				beego.NSRouter("/branch/show", List, "post:GetBranchList"),
				beego.NSRouter("/branch/change", List, "post:ChangeBranch"),
				// beego.NSRouter("/test", List, "get:Test"),
			),beego.NSNamespace("/ignore",
				beego.NSRouter("/get", List, "get:IgnoreList"),
				beego.NSRouter("/set", List, "post:SetIgnore"),
				beego.NSRouter("/statusChange", List, "post:UpdateIgnoreStatus"),
				beego.NSRouter("/delete", List, "post:DeleteIgnore"),
				// beego.NSRouter("/branch/change", List, "post:ChangeBranch"),
				// beego.NSRouter("/test", List, "get:Test"),
			),beego.NSNamespace("/setting",
				beego.NSRouter("/workPath", Menu, "post:SetWorkPath"),
				// beego.NSRouter("/test", List, "get:Test"),
			),
		)
	beego.AddNamespace(ns)
}
