package main

import (
	"manage/initialization"
	_ "manage/routers"

	_ "github.com/astaxie/beego/session/redis"

	"github.com/astaxie/beego"
)

func main() {
	// 设置资源路径
	beego.SetStaticPath("/css", "static/css")
	beego.SetStaticPath("/img", "static/img")
	beego.SetStaticPath("/js", "static/js")
	beego.SetStaticPath("/i", "static/i")
	beego.SetStaticPath("/fonts", "static/fonts")
	beego.SetStaticPath("/uploads", "static/upload")

	// 注册自定义函数
	beego.AddFuncMap("add", add)
	beego.AddFuncMap("sub", sub)

	// 初始化
	initialization.InitializeManager()

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
