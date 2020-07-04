package controllers

import (
	"manage/models"
	"strings"
)

type LogController struct {
	BaseController
}

// EmailLog 邮件日志列表
func (c *LogController) EmailLog() {
	c.Layout = "layouts/master.html"
	c.TplName = "log/email.html"

	c.LayoutSections = map[string]string{
		"Style":  "log/email_style.html",
		"Script": "log/email_script.html",
	}

	filter := map[string]interface{}{}

	et, _ := c.GetInt("et", -1)
	if et > 0 {
		filter["email_type"] = et
	}
	// 根据发送结果
	result := c.GetString("r")
	if len(result) > 0 {
		filter["result"] = result
	}
	// 根据收件人邮箱查询
	keyword := c.GetString("key")
	if len(keyword) > 0 {
		filter["address__istartswith"] = keyword
	}

	emailType := map[int64]string{0: "初始密码邮件", 1: "重置密码邮件"}
	logs, _ := models.GetAllEmailLogs(filter)

	c.Data["emailType"] = emailType
	c.Data["et"] = et
	c.Data["result"] = result
	c.Data["keyword"] = keyword
	c.Data["logs"] = logs
}

// HomeLog 前台日志列表
func (c *LogController) HomeLog() {
	c.Layout = "layouts/master.html"
	c.TplName = "log/home.html"
	//c.LayoutSections = map[string]string{
	//	"Style":  "log/index_style.html",
	//	"Script": "log/index_script.html",
	//}

	filter := map[string]interface{}{}
	logs, _ := models.GetAllHomeLogs(filter)

	c.Data = map[interface{}]interface{}{
		"logs": logs,
	}
}

// AdminLog 后台日志列表
func (c *LogController) AdminLog() {
	c.Layout = "layouts/master.html"
	c.TplName = "log/admin.html"
	c.LayoutSections = map[string]string{
		"Style":  "log/admin_style.html",
		"Script": "log/admin_script.html",
	}

	filter := map[string]interface{}{}

	// 根据请求方式查询
	method := c.GetString("m")
	if len(method) > 0 {
		filter["method"] = method
	}
	// 根据请求结果查询
	result := c.GetString("r")
	if len(result) > 0 {
		filter["result"] = result
	}
	// 根据操作人查询
	managerId, _ := c.GetInt("o")
	if managerId > 0 {
		filter["manager_id"] = managerId
	}
	// 根据关键字查询
	keyword := c.GetString("key")
	if len(keyword) > 0 {
		if strings.Contains(keyword, "/") {
			filter["url__istartswith"] = keyword
		} else {
			filter["content__istartswith"] = keyword
		}
	}

	managers, _ := models.GetAllManagers(nil)
	logs, _ := models.GetAllAdminLogs(filter)

	c.Data["managers"] = managers
	c.Data["method"] = method
	c.Data["result"] = result
	c.Data["managerId"] = managerId
	c.Data["keyword"] = keyword
	c.Data["logs"] = logs
}
