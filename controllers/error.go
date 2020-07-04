package controllers

import "github.com/astaxie/beego"

type ErrorController struct {
	beego.Controller
}

// Error401 401页面
func (c *ErrorController) Error401() {
	c.TplName = "error/401.html"
}

// 404页面
func (c *ErrorController) Error404() {
	c.TplName = "error/404.html"
}
