package main

import (
	"encoding/gob"
	"manage/controllers"
	"manage/initialization"
	"manage/models"
	_ "manage/routers"

	_ "github.com/astaxie/beego/session/redis"

	"github.com/astaxie/beego"
)

func init() {
	gob.Register(&models.Manager{})
}

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
	beego.AddFuncMap("getArticleNumByCategoryId", getArticleNumByCategoryId)
	beego.AddFuncMap("getArticleNumByTagId", getArticleNumByTagId)

	beego.ErrorController(&controllers.ErrorController{})

	// 初始化
	initialization.InitializeManager()

	beego.Run()
}

// add 加法计算
func add(i, step int) int {
	return i + step
}

// sub 减法计算
func sub(i, step int) int {
	return i - step
}

// getArticleNumByCategoryId 获取栏目下的文章数量
func getArticleNumByCategoryId(categoryId int) int64 {
	return models.GetArticleNumOfCategory(map[string]interface{}{"category_id": categoryId})
}

// getArticleNumByTagId 获取栏目下的文章数量
func getArticleNumByTagId(tagId int) int64 {
	return models.GetArticleNumOfTag(map[string]interface{}{"tag_id": tagId})
}
