package controllers

import (
	"manage/models"
	"strings"
	"time"

	utils "github.com/ObrookO/go-utils"
)

type AccountController struct {
	BaseController
}

// Get 账号列表
func (c *AccountController) Get() {
	c.Layout = "layouts/master.html"
	c.TplName = "account/index.html"
	c.LayoutSections = map[string]string{
		"Script": "account/index_script.html",
	}

	AddLog(c.Ctx, "查看账号列表", "", "PAGE")

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

	accounts, _ := models.GetAllAccounts(filter)

	c.Data["commentStatus"] = commentStatus
	c.Data["status"] = status
	c.Data["keyword"] = keyword
	c.Data["accounts"] = accounts
}

// ChangeCommentStatus 修改账号的评论权限
func (c *AccountController) ChangeCommentStatus() {
	accountId, _ := c.GetInt("id")
	status, _ := c.GetInt("status")
	statusSlice := []string{"禁用", "启用"}

	// 判断status是否合法
	if !utils.ObjInIntSlice(status, []int{0, 1}) {
		AddLog(c.Ctx, "修改账号评论权限", "无效的status", "{\"code\": 400000, \"msg\": \"参数错误\"}")
		c.Data["json"] = &JSONResponse{Code: 400000, Msg: "参数错误"}
		c.ServeJSON()
		return
	}

	// 判断账号是否存在
	account, _ := models.GetOneAccount(map[string]interface{}{"id": accountId})
	if account.Id == 0 {
		AddLog(c.Ctx, "修改账号评论权限", "账号不存在", "{\"code\": 400001, \"msg\": \"账号不存在\"}")
		c.Data["json"] = &JSONResponse{Code: 400001, Msg: "账号不存在"}
		c.ServeJSON()
		return
	}

	logContent := statusSlice[status] + "账号评论权限，账号：" + account.Username

	// 修改账号信息
	if _, err := models.UpdateAccountWithFilter(map[string]interface{}{"id": accountId}, map[string]interface{}{
		"allow_comment": status,
		"updated_at":    time.Now().Format("2006-01-02 15:04:05"),
	}); err != nil {
		AddLog(c.Ctx, logContent, err.Error(), "{\"code\": 400002, \"msg\": \"操作失败\"}")
		c.Data["json"] = &JSONResponse{Code: 400002, Msg: "操作失败"}
		c.ServeJSON()
		return
	}

	AddLog(c.Ctx, logContent, "", "{\"code\": 200, \"msg\": \"OK\"}")
	c.Data["json"] = &JSONResponse{Code: 200, Msg: "OK"}
	c.ServeJSON()
}

// ChangeStatus 修改账号的状态
func (c *AccountController) ChangeStatus() {
	accountId, _ := c.GetInt("id")
	status, _ := c.GetInt("status")
	statusSlice := []string{"禁用", "启用"}

	// 判断status是否合法
	if !utils.ObjInIntSlice(status, []int{0, 1}) {
		AddLog(c.Ctx, "修改账号状态", "无效的status", "{\"code\": 400000, \"msg\": \"参数错误\"}")
		c.Data["json"] = &JSONResponse{Code: 400000, Msg: "参数错误"}
		c.ServeJSON()
		return
	}
	// 判断账号是否存在
	account, _ := models.GetOneAccount(map[string]interface{}{"id": accountId})
	if account.Id == 0 {
		AddLog(c.Ctx, "修改账号状态", "账号不存在", "{\"code\": 400001, \"msg\": \"账号不存在\"}")
		c.Data["json"] = &JSONResponse{Code: 400001, Msg: "账号不存在"}
		c.ServeJSON()
		return
	}

	logContent := statusSlice[status] + "账号，账号：" + account.Username

	// 修改账号信息
	if _, err := models.UpdateAccountWithFilter(map[string]interface{}{"id": accountId}, map[string]interface{}{
		"status":     status,
		"updated_at": time.Now().Format("2006-01-02 15:04:05"),
	}); err != nil {
		AddLog(c.Ctx, logContent, "", "{\"code\": 400002, \"msg\": \"操作失败\"}")
		c.Data["json"] = &JSONResponse{Code: 400002, Msg: "操作失败"}
		c.ServeJSON()
		return
	}

	AddLog(c.Ctx, logContent, "", "{\"code\": 200, \"msg\": \"OK\"}")
	c.Data["json"] = &JSONResponse{Code: 200, Msg: "OK"}
	c.ServeJSON()
}
