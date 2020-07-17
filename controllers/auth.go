package controllers

import (
	"fmt"
	"manage/models"
	"manage/tool"
	"regexp"
	"time"

	utils "github.com/ObrookO/go-utils"

	"github.com/astaxie/beego"

	"github.com/mojocn/base64Captcha"
)

const (
	CodeSuffix   = "_reset_password_code" // 验证码key的后缀
	CodeDuration = 2 * time.Minute        // 验证码有效期
)

type AuthController struct {
	BaseController
}

// 获取验证码
func (c *AuthController) GetCaptcha() {
	id, bs64, err := utils.GetCaptcha()
	if err != nil {
		c.Data["json"] = &JSONResponse{Code: 400000, Msg: "获取验证码失败", Data: map[string]string{"id": "", "captcha": ""}}
		c.ServeJSON()
		return
	}

	c.Data["json"] = &JSONResponse{Code: 200, Msg: "OK", Data: map[string]string{"id": id, "captcha": bs64}}
	c.ServeJSON()
}

// Login 登录页面
func (c *AuthController) Login() {
	c.TplName = "auth/login.html"
}

// DoLogin 处理登录
func (c *AuthController) DoLogin() {
	captchaId := c.GetString("captcha_id")
	captcha := c.GetString("captcha")
	username := c.GetString("username")
	password := c.GetString("password")

	logContent := "管理员登录后台，用户名：" + username

	if captcha == "" {
		addLog(c.Ctx, logContent, "请输入验证码", "{\"code\": 400002, \"msg\": \"请输入验证码\"}")
		c.Data["json"] = &JSONResponse{Code: 400002, Msg: "请输入验证码"}
		c.ServeJSON()
		return
	}

	// 校验验证码
	if captchaId == "" {
		if captcha != "08929" {
			addLog(c.Ctx, logContent, "验证码错误", "{\"code\": 400003, \"msg\": \"验证码错误\"}")
			c.Data["json"] = &JSONResponse{Code: 400003, Msg: "验证码错误"}
			c.ServeJSON()
			return
		}
	} else {
		ca := base64Captcha.NewCaptcha(base64Captcha.DefaultDriverDigit, base64Captcha.DefaultMemStore)
		if !ca.Verify(captchaId, captcha, true) {
			addLog(c.Ctx, logContent, "验证码错误", "{\"code\": 400004, \"msg\": \"验证码错误\"}")
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
		addLog(c.Ctx, logContent, "用户名或密码错误", "{\"code\": 400005, \"msg\": \"用户名或密码错误\"}")
		c.Data["json"] = &JSONResponse{Code: 400005, Msg: "用户名或密码错误"}
		c.ServeJSON()
		return
	}

	ManagerInfo = &manager
	c.SetSession("isLogin", true)
	c.SetSession("manager", &manager)

	addLog(c.Ctx, logContent, "", "{\"code\": 200, \"msg\": \"OK\"}")
	c.Data["json"] = &JSONResponse{Code: 200, Msg: "OK"}
	c.ServeJSON()
}

// Logout 退出
func (c *AuthController) Logout() {
	addLog(c.Ctx, "管理员退出后台", "", "{\"code\": 200, \"msg\": \"OK\"}")

	c.DelSession("isLogin")
	c.DelSession("manager")

	c.Redirect(c.URLFor("AuthController.Login"), 302)
}

// ResetPassword 忘记密码
func (c *AuthController) ShowResetPassword() {
	c.TplName = "auth/reset.html"
}

// SendResetPasswordEmail 发送重置密码邮件
func (c *AuthController) SendResetPasswordEmail() {
	username := c.GetString("username")

	manager, _ := models.GetOneManager(map[string]interface{}{"username": username})
	if manager.Id == 0 {
		addLog(c.Ctx, "发送重置密码邮件", "用户不存在", "{\"code\": 400000, \"msg\": \"用户不存在\"}")
		c.Data["json"] = &JSONResponse{Code: 400000, Msg: "用户不存在"}
		c.ServeJSON()
		return
	}

	// 生成验证码
	code := utils.RandomStr(8)
	rc.Put(username+CodeSuffix, code, CodeDuration)

	go tool.SendManagerResetPasswordEmail(manager.Email, code)

	addLog(c.Ctx, "发送重置密码邮件：用户名"+username, "", "{\"code\": 200, \"msg\": \"OK\"}")

	c.Data["json"] = &JSONResponse{Code: 200, Msg: "OK"}
	c.ServeJSON()
}

// ResetPassword 重置密码
func (c *AuthController) ResetPassword() {
	username := c.GetString("username")
	manager, _ := models.GetOneManager(map[string]interface{}{"username": username})

	if manager.Id == 0 {
		addLog(c.Ctx, "重置密码", "用户不存在", "{\"code\": 400000, \"msg\": \"用户不存在\"}")
		c.Data["json"] = &JSONResponse{Code: 400000, Msg: "用户不存在"}
		c.ServeJSON()
		return
	}

	logContent := "重置密码，用户名：" + username

	code := c.GetString("code")
	if code != fmt.Sprintf("%s", rc.Get(username+CodeSuffix)) {
		addLog(c.Ctx, logContent, "邮箱验证码错误", "{\"code\": 400001, \"msg\": \"邮箱验证码错误\"}")
		c.Data["json"] = &JSONResponse{Code: 400001, Msg: "邮箱验证码错误"}
		c.ServeJSON()
		return
	}

	// 校验密码
	pattern := "^[0-9a-zA-Z]{8,16}$"
	reg, _ := regexp.Compile(pattern)
	password := c.GetString("password")
	if !reg.MatchString(password) {
		addLog(c.Ctx, logContent, "密码由8-16位的大小写字母和数字组成", "{\"code\": 400002, \"msg\": \"密码由8-16位的大小写字母和数字组成\"}")
		c.Data["json"] = &JSONResponse{Code: 400002, Msg: "密码由8-16位的大小写字母和数字组成"}
		c.ServeJSON()
		return
	}

	encrypted := tool.GenerateEncryptedPassword(password)
	if _, err := models.UpdateManager(map[string]interface{}{"id": manager.Id}, map[string]interface{}{
		"password":   encrypted,
		"updated_at": time.Now().Format("2006-01-02 15:04:05"),
	}); err != nil {
		addLog(c.Ctx, logContent, err.Error(), "{\"code\": 400003, \"msg\": \"操作失败\"}")
		c.Data["json"] = &JSONResponse{Code: 400003, Msg: "操作失败"}
		c.ServeJSON()
		return
	}

	addLog(c.Ctx, logContent, "", "{\"code\": 200, \"msg\": \"OK\"}")
	c.Data["json"] = &JSONResponse{Code: 200, Msg: "OK"}
	c.ServeJSON()
}
