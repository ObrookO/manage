package main

import (
	_ "manage/routers"

	"github.com/astaxie/beego/orm"

	_ "github.com/astaxie/beego/session/redis"
	_ "github.com/go-sql-driver/mysql"

	"github.com/astaxie/beego"
)

func init() {
	orm.RegisterDriver("mysql", orm.DRMySQL)
	// 注册默认数据库
	orm.RegisterDataBase("default", "mysql", beego.AppConfig.String("db_user")+":"+beego.AppConfig.String(
		"db_password")+"@tcp(127.0.0.1:3306)/blog?charset=utf8&loc=Asia%2FShanghai")
}

func main() {
	beego.SetStaticPath("/acss", "static/css")
	beego.SetStaticPath("/aimg", "static/img")
	beego.SetStaticPath("/ajs", "static/js")
	beego.SetStaticPath("/ai", "static/i")
	beego.SetStaticPath("/afonts", "static/fonts")

	// 注册自定义函数
	beego.AddFuncMap("getIndex", getIndex)

	beego.Run()
}

// 获取索引
func getIndex(i, step int) int {
	return i + step
}
