package controllers

import (
	"html/template"
)

type AuthController struct {
	BaseController
}

// Login 登录页面
func (c *AuthController) Login() {
	c.TplName = "auth/login.html"

	c.Data = map[interface{}]interface{}{
		"xsrfdata": template.HTML(c.XSRFFormHTML()),
	}
}

// 处理登录
func (c *AuthController) DoLogin() {
	c.EnableRender = false
}
