package controllers

import (
	"html/template"
	"manage/models"
	"strings"
	"time"

	utils "github.com/ObrookO/go-utils"
)

type AccountController struct {
	BaseController
}

func (c *AccountController) Get() {
	c.Layout = "layouts/master.html"
	c.TplName = "account/index.html"
	c.LayoutSections = map[string]string{
		"Style":  "account/style.html",
		"Script": "account/script.html",
	}

	AddLog(c.Ctx, "查看账号列表", "", "PAGE", "SUCCESS")

	filter := map[string]interface{}{}
	// 根据评论权限查询
	commentStatus, _ := c.GetInt("cs", -1)
	if utils.ObjInIntSlice(commentStatus, []int{0, 1}) {
		filter["allow_comment"] = commentStatus
	}

	// 根据账号权限查询
	status, _ := c.GetInt("s", -1)
	if utils.ObjInIntSlice(status, []int{0, 1}) {
		filter["status"] = status
	}

	// 根据关键字查询
	keyword := c.GetString("key")
	if len(keyword) > 0 {
		if strings.Contains(keyword, "@") {
			strSlice := strings.Split(keyword, "@")
			filter["email__istartswith"] = strSlice[0]
		} else {
			filter["username__istartswith"] = keyword
		}
	}

	accounts, _ := models.GetAccounts(filter)

	c.Data = map[interface{}]interface{}{
		"xsrfdata":      template.HTML(c.XSRFFormHTML()),
		"commentStatus": commentStatus,
		"status":        status,
		"keyword":       keyword,
		"accounts":      accounts,
	}
}

// ChangeCommentStatus 修改账号的评论权限
func (c *AccountController) ChangeCommentStatus() {
	c.EnableRender = false

	accountId, _ := c.GetInt("id")
	status, _ := c.GetInt("status")
	statusSlice := []string{"禁用", "启用"}

	// 判断账号是否存在
	account, _ := models.GetOneAccount(map[string]interface{}{"id": accountId})
	if account.Id == 0 {
		AddLog(c.Ctx, statusSlice[status]+"账号 "+account.Username+" 评论权限", "账号不存在", "{\"code\": 400000, \"msg\": \"账号不存在\"}", "FAIL")
		c.Data["json"] = &JSONResponse{Code: 400000, Msg: "账号不存在"}
		c.ServeJSON()
		return
	}

	// 判断status是否合法
	if !utils.ObjInIntSlice(status, []int{0, 1}) {
		AddLog(c.Ctx, "修改账号 "+account.Username+" 评论权限", "无效的status", "{\"code\": 400001, \"msg\": \"无效的status\"}", "FAIL")
		c.Data["json"] = &JSONResponse{Code: 400001, Msg: "无效的status"}
		c.ServeJSON()
		return
	}

	// 修改账号信息
	if _, err := models.UpdateAccount(map[string]interface{}{"id": accountId}, map[string]interface{}{
		"allow_comment": status,
		"updated_at":    time.Now().Format("2006-01-02 15:04:05"),
	}); err != nil {
		AddLog(c.Ctx, statusSlice[status]+"账号 "+account.Username+" 评论权限", err.Error(), "{\"code\": 400002, \"msg\": \"操作失败\"}", "FAIL")
		c.Data["json"] = &JSONResponse{Code: 400002, Msg: "操作失败"}
		c.ServeJSON()
		return
	}

	AddLog(c.Ctx, statusSlice[status]+"账号 "+account.Username+" 评论权限", "", "{\"code\": 200, \"msg\": \"OK\"}", "SUCCESS")
	c.Data["json"] = &JSONResponse{Code: 200, Msg: "OK"}
	c.ServeJSON()
}

// ChangeStatus 修改账号的状态
func (c *AccountController) ChangeStatus() {
	c.EnableRender = false

	accountId, _ := c.GetInt("id")
	status, _ := c.GetInt("status")
	statusSlice := []string{"禁用", "启用"}

	// 判断账号是否存在
	account, _ := models.GetOneAccount(map[string]interface{}{"id": accountId})
	if account.Id == 0 {
		AddLog(c.Ctx, statusSlice[status]+"账号 "+account.Username, "账号不存在", "{\"code\": 400000, \"msg\": \"账号不存在\"}", "FAIL")
		c.Data["json"] = &JSONResponse{Code: 400000, Msg: "账号不存在"}
		c.ServeJSON()
		return
	}

	// 判断status是否合法
	if !utils.ObjInIntSlice(status, []int{0, 1}) {
		AddLog(c.Ctx, "修改账号 "+account.Username+" 状态", "无效的status", "{\"code\": 400001, \"msg\": \"无效的status\"}", "FAIL")
		c.Data["json"] = &JSONResponse{Code: 400001, Msg: "无效的status"}
		c.ServeJSON()
		return
	}

	// 修改账号信息
	if _, err := models.UpdateAccount(map[string]interface{}{"id": accountId}, map[string]interface{}{
		"status":     status,
		"updated_at": time.Now().Format("2006-01-02 15:04:05"),
	}); err != nil {
		AddLog(c.Ctx, statusSlice[status]+"账号 "+account.Username, "", "{\"code\": 400002, \"msg\": \"操作失败\"}", "FAIL")
		c.Data["json"] = &JSONResponse{Code: 400002, Msg: "操作失败"}
		c.ServeJSON()
		return
	}

	AddLog(c.Ctx, statusSlice[status]+"账号 "+account.Username, "", "{\"code\": 200, \"msg\": \"OK\"}", "SUCCESS")
	c.Data["json"] = &JSONResponse{Code: 200, Msg: "OK"}
	c.ServeJSON()
}
