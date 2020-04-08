package controllers

import (
	"encoding/json"
	"manage/models"

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

func AddLog(ctx *context.Context, content, reason, response string) {
	// 请求头转json
	h, _ := json.Marshal(ctx.Request.Header)

	b := []byte{}
	// body转json
	if len(ctx.Request.PostForm) > 0 {
		b, _ = json.Marshal(ctx.Request.PostForm)
	}

	adminLog := models.AdminLog{
		ManagerId: 0,
		Ip:        ctx.Input.IP(),
		Url:       ctx.Request.URL.Path,
		Method:    ctx.Request.Method,
		Query:     ctx.Request.URL.RawQuery,
		Headers:   string(h),
		Body:      string(b),
		Response:  response,
		Content:   content,
		Reason:    reason,
	}

	models.AddAdminLog(adminLog)
}
