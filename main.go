package main

import (
	_ "manage/routers"

	"github.com/astaxie/beego/orm"

	_ "github.com/astaxie/beego/session/redis"
	_ "github.com/go-sql-driver/mysql"

	"github.com/astaxie/beego"
)

func init() {
	host := beego.AppConfig.String("db_host")
	db := beego.AppConfig.String("db_name")
	user := beego.AppConfig.String("db_user")
	pass := beego.AppConfig.String("db_password")

	orm.RegisterDriver("mysql", orm.DRMySQL)
	// 注册默认数据库
	orm.RegisterDataBase("default", "mysql", user+":"+pass+"@tcp("+host+")/"+db+"?charset=utf8&loc=Asia%2FShanghai")
}

func main() {
	if beego.BConfig.RunMode == "dev" {
		//orm.Debug = true
	}

	beego.SetStaticPath("/css", "static/css")
	beego.SetStaticPath("/img", "static/img")
	beego.SetStaticPath("/js", "static/js")
	beego.SetStaticPath("/i", "static/i")
	beego.SetStaticPath("/fonts", "static/fonts")

	// 注册自定义函数
	beego.AddFuncMap("add", add)
	beego.AddFuncMap("sub", sub)

	beego.Run()
}

//  加法计算
func add(i, step int) int {
	return i + step
}

// 减法计算
func sub(i, step int) int {
	return i - step
}
