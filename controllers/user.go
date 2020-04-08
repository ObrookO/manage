package controllers

import (
	"html/template"
	"manage/models"
	"manage/utils"
	"strings"
	"time"
)

type UserController struct {
	BaseController
}

func (c *UserController) Get() {
	c.Layout = "layouts/master.html"
	c.TplName = "user/index.html"
	c.LayoutSections = map[string]string{
		"Style":  "user/style.html",
		"Script": "user/script.html",
	}

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

	users, err := models.GetUsers(filter)
	if err != nil {
		AddLog(c.Ctx, "获取用户列表", err.Error(), "page")
	} else {
		AddLog(c.Ctx, "获取用户列表", "", "page")
	}

	c.Data = map[interface{}]interface{}{
		"xsrfdata":      template.HTML(c.XSRFFormHTML()),
		"commentStatus": commentStatus,
		"status":        status,
		"keyword":       keyword,
		"users":         users,
	}
}

// ChangeCommentStatus 修改用户的评论状态
func (c *UserController) ChangeCommentStatus() {
	c.EnableRender = false

	userId, _ := c.GetInt("id")
	status, _ := c.GetInt("status")

	// 判断用户是否存在
	if !models.IsUserExists(map[string]interface{}{"id": userId}) {
		c.Data["json"] = &JSONResponse{Code: 400000, Msg: "用户不存在"}
		c.ServeJSON()
		AddLog(c.Ctx, "修改用户评论权限", "用户不存在", "{\"code\": 400000, \"msg\": \"用户不存在\"}")
		return
	}

	// 判断status是否合法
	if !utils.ObjInIntSlice(status, []int{0, 1}) {
		c.Data["json"] = &JSONResponse{Code: 400001, Msg: "无效的status"}
		c.ServeJSON()
		AddLog(c.Ctx, "修改用户评论权限", "无效的status", "{\"code\": 400001, \"msg\": \"无效的status\"}")
		return
	}

	// 修改用户信息
	if _, err := models.UpdateUser(map[string]interface{}{"id": userId}, map[string]interface{}{
		"allow_comment": status,
		"updated_at":    time.Now().Format("2006-01-02 15:04:05"),
	}); err != nil {
		c.Data["json"] = &JSONResponse{Code: 400002, Msg: "操作失败"}
		c.ServeJSON()
		AddLog(c.Ctx, "修改用户评论权限", err.Error(), "{\"code\": 400002, \"msg\": \"操作失败\"}")
		return
	}

	c.Data["json"] = &JSONResponse{Code: 200, Msg: "OK"}
	c.ServeJSON()
	AddLog(c.Ctx, "修改用户评论权限", "", "{\"code\": 200, \"msg\": \"OK\"}")
}

// ChangeStatus 修改用户的状态
func (c *UserController) ChangeStatus() {
	c.EnableRender = false

	userId, _ := c.GetInt("id")
	status, _ := c.GetInt("status")

	// 判断用户是否存在
	if !models.IsUserExists(map[string]interface{}{"id": userId}) {
		c.Data["json"] = &JSONResponse{Code: 400000, Msg: "用户不存在"}
		c.ServeJSON()
		AddLog(c.Ctx, "修改用户状态", "", "{\"code\": 400000, \"msg\": \"用户不存在\"}")
		return
	}

	// 判断status是否合法
	if !utils.ObjInIntSlice(status, []int{0, 1}) {
		c.Data["json"] = &JSONResponse{Code: 400001, Msg: "无效的status"}
		c.ServeJSON()
		AddLog(c.Ctx, "修改用户状态", "", "{\"code\": 400001, \"msg\": \"无效的status\"}")
		return
	}

	// 修改用户信息
	if _, err := models.UpdateUser(map[string]interface{}{"id": userId}, map[string]interface{}{
		"status":     status,
		"updated_at": time.Now().Format("2006-01-02 15:04:05"),
	}); err != nil {
		c.Data["json"] = &JSONResponse{Code: 400002, Msg: "操作失败"}
		c.ServeJSON()
		AddLog(c.Ctx, "修改用户状态", "", "{\"code\": 400002, \"msg\": \"操作失败\"}")
		return
	}

	c.Data["json"] = &JSONResponse{Code: 200, Msg: "OK"}
	c.ServeJSON()
	AddLog(c.Ctx, "修改用户状态", "", "{\"code\": 200, \"msg\": \"OK\"}")
}
