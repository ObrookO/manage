package controllers

import (
	"encoding/json"
	"html/template"
	"manage/models"
	"strings"

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
	IsAdmin     bool            // 是否是超级管理员
	ManagerInfo *models.Manager // 管理员信息
)

func (c *BaseController) Prepare() {
	// 不需要验证登录的路由
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

		// 判断是否登录
		if l != nil && m != nil {
			ManagerInfo = m.(*models.Manager)
			c.Data["nickname"] = ManagerInfo.Nickname
			if ManagerInfo.IsAdmin == 1 {
				IsAdmin = true
			}

			if !c.IsAjax() {
				c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML()) // 定义全局xsrf
				c.Data["isAdmin"] = IsAdmin                          // 定义是否是管理员
			} else {
				c.EnableRender = false // ajax请求不加载模板
			}

			// 判断权限
			if !IsAdmin {
				if strings.Contains(path, "/logs") || strings.Contains(path, "/accounts") || strings.Contains(path, "/managers") {
					if c.IsAjax() {
						c.Data["json"] = &JSONResponse{Code: 500, Msg: "非法访问"}
						c.ServeJSON()
					} else {
						c.Abort("401")
					}
				}
			}
		} else {
			c.Redirect(c.URLFor("AuthController.Login"), 302)
		}
	}
}

// 记录日志
func AddLog(ctx *context.Context, content, reason, response string) {
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

	result := "SUCCESS"
	if reason != "" {
		result = "FAIL"
	}

	log := models.AdminLog{
		Content:  content,
		Ip:       ctx.Input.IP(),
		Url:      path,
		Method:   method,
		Query:    ctx.Request.URL.RawQuery,
		Headers:  string(h),
		Body:     body,
		Response: response,
		Result:   result,
		Reason:   reason,
	}

	log.Manager = ManagerInfo

	models.AddAdminLog(log)
}
