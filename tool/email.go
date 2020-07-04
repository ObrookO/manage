package tool

import (
	"manage/models"
	"strings"

	"github.com/astaxie/beego"
	"gopkg.in/gomail.v2"
)

const (
	newManager    = iota // 添加用户
	resetPassword        // 重置密码
)

var appURL = beego.AppConfig.String("appurl")                                            // 博客首页地址
var emailFooter = "登录地址：<a href=\"" + appURL + "\" target=\"_blank\">" + appURL + "</a>" // 邮件页脚

// SendNewManagerEmail 给新创建的用户发送邮件
func SendNewManagerEmail(toAddress, account, rawPassword string) {
	subject := "账号创建成功通知"
	contentType := "text/html"
	content := []string{
		"登录账号：" + account,
		"初始密码：" + rawPassword,
		strings.Repeat("<br>", 3) + emailFooter,
	}

	sendEmail(newManager, toAddress, subject, contentType, strings.Join(content, "<br>"))
}

// sendEmail 发送邮件
func sendEmail(emailType int, toAddress, subject, contentType, msg string) {
	m := gomail.NewMessage()

	host := beego.AppConfig.String("email_host")
	port, _ := beego.AppConfig.Int("email_port")
	fromAddress := beego.AppConfig.String("email_from_address")
	fromName := beego.AppConfig.String("email_from_name")
	password := beego.AppConfig.String("email_password")

	m.SetHeader("From", fromAddress)
	// 设置发送人别名
	m.SetAddressHeader("From", fromAddress, fromName)
	m.SetHeader("To", toAddress)
	m.SetHeader("Subject", subject)
	m.SetBody(contentType, msg)

	var reason string
	var result = "SUCCESS"

	d := gomail.NewDialer(host, port, fromAddress, password)
	if err := d.DialAndSend(m); err != nil {
		result = "FAIL"
		reason = err.Error()
	}

	models.AddEmailLog(models.EmailLog{
		EmailType: emailType,
		Address:   toAddress,
		Content:   msg,
		Result:    result,
		Reason:    reason,
	})
}
