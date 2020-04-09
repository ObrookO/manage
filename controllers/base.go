package controllers

import (
	"encoding/json"
	"manage/models"
	"manage/utils"

	"github.com/mojocn/base64Captcha"

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

// 获取验证码
func (c *BaseController) GetCaptcha() {
	c.EnableRender = false

	captcha := base64Captcha.NewCaptcha(base64Captcha.DefaultDriverDigit, base64Captcha.DefaultMemStore)
	id, bs64, err := captcha.Generate()

	if err != nil {
		AddLog(c.Ctx, "生成验证码", err.Error(), "{\"code\": 400000, \"msg\": \"获取验证码失败\"}", "FAIL")
		c.Data["json"] = &JSONResponse{Code: 400000, Msg: "获取验证码失败", Data: map[string]string{"id": utils.RandomStr(20), "captcha": ""}}
		c.ServeJSON()
		return
	}

	AddLog(c.Ctx, "生成验证码", "", "{\"code\": 200, \"msg\": \"OK\"}", "SUCCESS")
	c.Data["json"] = &JSONResponse{Code: 200, Msg: "OK", Data: map[string]string{"id": id, "captcha": bs64}}
	c.ServeJSON()
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
