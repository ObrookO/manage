package controllers

import (
	"encoding/json"
	"html/template"
	"manage/models"

	"github.com/astaxie/beego/utils"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

type BaseController struct {
	beego.Controller
}

type JSONResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

var (
	ManagerInfo = map[string]interface{}{} // 管理员信息
)

func (c *BaseController) Prepare() {
	if !c.IsAjax() {
		c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML()) // 定义全局xsrf
	} else {
		c.EnableRender = false // ajax请求不加载模板
	}

	noCheckUrl := []string{
		"GET" + beego.URLFor("AuthController.GetCaptcha"),
		"GET" + beego.URLFor("AuthController.Login"),
		"POST" + beego.URLFor("AuthController.DoLogin"),
		"GET" + beego.URLFor("AuthController.Logout"),
	}

	method := c.Ctx.Request.Method
	path := c.Ctx.Request.URL.Path
	if !utils.InSlice(method+path, noCheckUrl) {
		l := c.GetSession("isLogin")
		m := c.GetSession("manager")

		if l != nil && m != nil {
			ManagerInfo = m.(map[string]interface{})
			c.Data["username"] = ManagerInfo["username"]
		} else {
			c.Redirect(c.URLFor("AuthController.Login"), 302)
		}
	}
}

// 记录日志
func AddLog(ctx *context.Context, content, reason, response, result string) {
	// 请求头转json
	h, _ := json.Marshal(ctx.Request.Header)

	// body转json
	b := []byte{}
	if len(ctx.Request.PostForm) > 0 {
		b, _ = json.Marshal(ctx.Request.PostForm)
	}

	// 不记录请求体的地址
	noSaveBodyUrl := []string{
		"POST" + beego.URLFor("AuthController.DoLogin"),
	}

	body := ""
	method := ctx.Request.Method
	path := ctx.Request.URL.Path
	if !utils.InSlice(method+path, noSaveBodyUrl) {
		body = string(b)
	}

	log := models.AdminLog{
		ManagerId: 0,
		Content:   content,
		Ip:        ctx.Input.IP(),
		Url:       path,
		Method:    method,
		Query:     ctx.Request.URL.RawQuery,
		Headers:   string(h),
		Body:      body,
		Response:  response,
		Result:    result,
		Reason:    reason,
	}

	models.AddAdminLog(log)
}
