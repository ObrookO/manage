package controllers

import (
	"manage/models"
	"strings"
)

type LogController struct {
	BaseController
}

var logPageLimit = 20

// HomeLog 前台日志列表
func (c *LogController) HomeLog() {
	c.Layout = "layouts/master.html"
	c.TplName = "log/home.html"
	c.LayoutSections = map[string]string{
		"Style":  "log/style.html",
		"Script": "log/script.html",
	}

	AddLog(c.Ctx, "查看前台日志列表", "", "PAGE", "SUCCESS")

	filter := map[string]interface{}{}
	logs, _ := models.GetHomeLogs(filter, 0, logPageLimit)

	c.Data = map[interface{}]interface{}{
		"logs": logs,
	}
}

// AdminLog 后台日志列表
func (c *LogController) AdminLog() {
	c.Layout = "layouts/master.html"
	c.TplName = "log/admin.html"
	c.LayoutSections = map[string]string{
		"Style":  "log/style.html",
		"Script": "log/script.html",
	}

	AddLog(c.Ctx, "查看后台日志列表", "", "PAGE", "SUCCESS")

	offset := 0
	filter := map[string]interface{}{}

	page, _ := c.GetInt("p")
	if page > 0 {
		offset = (page - 1) * logPageLimit
	} else {
		page = 1
	}

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

	// 根据关键字查询
	keyword := c.GetString("key")
	if len(keyword) > 0 {
		if strings.Contains(keyword, "/") {
			filter["url__istartswith"] = keyword
		} else {
			filter["content__istartswith"] = keyword
		}
	}

	logs, _ := models.GetAdminLogs(filter, offset, logPageLimit)

	c.Data = map[interface{}]interface{}{
		"page":    page,
		"method":  method,
		"result":  result,
		"keyword": keyword,
		"logs":    logs,
		"num":     len(logs),
		"limit":   logPageLimit,
	}
}
