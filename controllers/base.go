package controllers

import (
	"encoding/json"
	"manage/models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

type JSONResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type BaseController struct {
	beego.Controller
}

var (
	ManagerInfo map[string]interface{} // 管理员信息
)

func (c *BaseController) Prepare() {
	l := c.GetSession("isLogin")
	m := c.GetSession("manager")

	if l != nil && m != nil {
		ManagerInfo = m.(map[string]interface{})
		c.Data = map[interface{}]interface{}{
			"username": ManagerInfo["username"],
		}
	} else {
		c.Redirect(c.URLFor("AuthController.Login"), 302)
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

	log := models.AdminLog{
		ManagerId: 0,
		Content:   content,
		Ip:        ctx.Input.IP(),
		Url:       ctx.Request.URL.Path,
		Method:    ctx.Request.Method,
		Query:     ctx.Request.URL.RawQuery,
		Headers:   string(h),
		Body:      string(b),
		Response:  response,
		Result:    result,
		Reason:    reason,
	}

	models.AddAdminLog(log)
}
