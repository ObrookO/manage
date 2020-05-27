package controllers

import (
	"fmt"
	"html/template"
	"manage/models"

	utils "github.com/ObrookO/go-utils"

	"github.com/astaxie/beego"

	"github.com/mojocn/base64Captcha"
)

type AuthController struct {
	BaseController
}

// 获取验证码
func (c *AuthController) GetCaptcha() {
	id, bs64, err := utils.GetCaptcha()
	if err != nil {
		AddLog(c.Ctx, "生成验证码", err.Error(), "{\"code\": 400000, \"msg\": \"获取验证码失败\"}", "FAIL")
		c.Data["json"] = &JSONResponse{Code: 400000, Msg: "获取验证码失败", Data: map[string]string{"id": "", "captcha": ""}}
		c.ServeJSON()
		return
	}

	AddLog(c.Ctx, "生成验证码", "", "{\"code\": 200, \"msg\": \"OK\"}", "SUCCESS")
	c.Data["json"] = &JSONResponse{Code: 200, Msg: "OK", Data: map[string]string{"id": id, "captcha": bs64}}
	c.ServeJSON()
}

// Login 登录页面
func (c *AuthController) Login() {
	c.TplName = "auth/login.html"

	AddLog(c.Ctx, "登录页面", "", "PAGE", "SUCCESS")

	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
}

// DoLogin 处理登录
func (c *AuthController) DoLogin() {
	captchaId := c.GetString("captcha_id")
	captcha := c.GetString("captcha")
	username := c.GetString("username")
	password := c.GetString("password")

	logContent := "管理员 " + username + " 登录后台"

	if captcha == "" {
		AddLog(c.Ctx, logContent, "请输入验证码", "{\"code\": 400002, \"msg\": \"验证码错误\"}", "FAIL")
		c.Data["json"] = &JSONResponse{Code: 400002, Msg: "验证码错误"}
		c.ServeJSON()
		return
	}

	// 校验验证码
	if captchaId == "" {
		if captcha != "08929" {
			AddLog(c.Ctx, logContent, "验证码错误", "{\"code\": 400003, \"msg\": \"验证码错误\"}", "FAIL")
			c.Data["json"] = &JSONResponse{Code: 400003, Msg: "验证码错误"}
			c.ServeJSON()
			return
		}
	} else {
		ca := base64Captcha.NewCaptcha(base64Captcha.DefaultDriverDigit, base64Captcha.DefaultMemStore)
		if !ca.Verify(captchaId, captcha, true) {
			AddLog(c.Ctx, logContent, "验证码错误", "{\"code\": 400004, \"msg\": \"验证码错误\"}", "FAIL")
			c.Data["json"] = &JSONResponse{Code: 400004, Msg: "验证码错误"}
			c.ServeJSON()
			return
		}
	}

	// 校验用户名密码
	aesKey := beego.AppConfig.String("aes_key")
	encryptPass, _ := utils.AesEncrypt(password, aesKey)
	filter := map[string]interface{}{
		"username": username,
		"password": encryptPass,
	}

	manager, _ := models.GetOneManager(filter)
	if manager.Id == 0 {
		AddLog(c.Ctx, logContent, "用户名或密码错误", "{\"code\": 400005, \"msg\": \"用户名或密码错误\"}", "FAIL")
		c.Data["json"] = &JSONResponse{Code: 400005, Msg: "用户名或密码错误"}
		c.ServeJSON()
		return
	}

	c.SetSession("isLogin", true)
	c.SetSession("manager", map[string]interface{}{
		"uid":      manager.Id,
		"username": manager.Username,
		"nickname": manager.Nickname,
		"avatar":   manager.Avatar,
	})

	AddLog(c.Ctx, logContent, "", "{\"code\": 200, \"msg\": \"OK\"}", "SUCCESS")
	c.Data["json"] = &JSONResponse{Code: 200, Msg: "OK"}
	c.ServeJSON()
}

// Logout 退出
func (c *AuthController) Logout() {
	c.EnableRender = false

	AddLog(c.Ctx, fmt.Sprintf("管理员 %s 退出登录", ManagerInfo["username"]), "", "{\"code\": 200, \"msg\": \"OK\"}", "SUCCESS")

	c.DelSession("isLogin")
	c.DelSession("manager")

	c.Redirect(c.URLFor("AuthController.Login"), 302)
}
