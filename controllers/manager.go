package controllers

import (
	"manage/models"
	"manage/tool"
)

type ManagerController struct {
	BaseController
}

// Get 用户列表
func (c *ManagerController) Get() {
	c.Layout = "layouts/master.html"
	c.TplName = "manager/index.html"
	c.LayoutSections = map[string]string{
		"Style":  "manager/index_style.html",
		"Script": "manager/index_script.html",
	}

	addLog(c.Ctx, "查看用户列表", "", "PAGE")

	managers, _ := models.GetAllManagers(nil)

	c.Data["managers"] = managers
}

// Post 添加用户
func (c *ManagerController) Post() {
	manager := models.Manager{}

	if err := c.ParseForm(&manager); err != nil {
		addLog(c.Ctx, "添加用户", err.Error(), "{\"code\":400000,\"msg\": \"操作失败\"}")
		c.Data["json"] = &JSONResponse{Code: 400000, Msg: "操作失败"}
		c.ServeJSON()
		return
	}

	logContent := "添加用户，用户名：" + manager.Username

	// 表单验证
	if err := validData(manager); err != nil {
		addLog(c.Ctx, logContent, err.Error(), "{\"code\": 400001, \"msg\": \""+err.Error()+"\"}")
		c.Data["json"] = &JSONResponse{Code: 400001, Msg: err.Error()}
		c.ServeJSON()
		return
	}

	// 判断用户名是否重复
	if models.IsManagerExists(map[string]interface{}{"username": manager.Username}) {
		addLog(c.Ctx, logContent, "用户已存在", "{\"code\": 400002, \"msg\": \"用户已存在\"}")
		c.Data["json"] = &JSONResponse{Code: 400002, Msg: "用户已存在"}
		c.ServeJSON()
		return
	}

	// 判断邮箱是否重复
	if models.IsManagerExists(map[string]interface{}{"email": manager.Email}) {
		addLog(c.Ctx, logContent, "邮箱已被占用", "{\"code\": 400003, \"msg\": \"邮箱已被占用\"}")
		c.Data["json"] = &JSONResponse{Code: 400003, Msg: "邮箱已被占用"}
		c.ServeJSON()
		return
	}

	rawPassword := tool.GenerateRawPassword()
	manager.Password = tool.GenerateEncryptedPassword(rawPassword)
	if _, err := models.AddManager(manager); err != nil {
		addLog(c.Ctx, logContent, err.Error(), "{\"code\": 400004, \"msg\": \"操作失败\"}")
		c.Data["json"] = &JSONResponse{Code: 400004, Msg: "操作失败"}
		c.ServeJSON()
		return
	}

	// 发送用户名和密码
	go tool.SendNewManagerEmail(manager.Email, manager.Username, rawPassword)

	addLog(c.Ctx, logContent, "", "{\"code\": 200, \"msg\": \"OK\"}")
	c.Data["json"] = &JSONResponse{Code: 200, Msg: "OK"}
	c.ServeJSON()
}

// GetInfo 获取用户信息
func (c *ManagerController) GetInfo() {
	id, _ := c.GetInt("id")
	manager, _ := models.GetOneManager(map[string]interface{}{"id": id}, "id", "username", "nickname", "email", "is_admin")
	logContent := "获取用户信息，用户名：" + manager.Username

	if manager.Id == 0 {
		addLog(c.Ctx, logContent, "用户不存在", "{\"code\": 400000, \"msg\": \"用户不存在\"}")
		c.Data["json"] = &JSONResponse{Code: 400000, Msg: "用户不存在"}
		c.ServeJSON()
		return
	}

	addLog(c.Ctx, logContent, "", "{\"code\": 200, \"msg\": \"OK\"}")
	c.Data["json"] = &JSONResponse{Code: 200, Msg: "OK", Data: manager}
	c.ServeJSON()
}

// Update 编辑用户
func (c *ManagerController) Update() {
	id, _ := c.GetInt("id")
	manager := models.Manager{}

	if err := c.ParseForm(&manager); err != nil {
		addLog(c.Ctx, "编辑用户", err.Error(), "{\"code\":400000,\"msg\": \"操作失败\"}")
		c.Data["json"] = &JSONResponse{Code: 400000, Msg: "操作失败"}
		c.ServeJSON()
		return
	}

	rawManager, _ := models.GetOneManager(map[string]interface{}{"id": id})
	if rawManager.Id == 0 {
		addLog(c.Ctx, "编辑用户", "用户不存在", "{\"code\": 400001, \"msg\": \"用户不存在\"}")
		c.Data["json"] = &JSONResponse{Code: 400001, Msg: "用户不存在"}
		c.ServeJSON()
		return
	}

	logContent := "编辑用户，用户名：" + rawManager.Username

	// 表单验证
	if err := validData(manager); err != nil {
		addLog(c.Ctx, logContent, err.Error(), "{\"code\": 400002, \"msg\": \""+err.Error()+"\"}")
		c.Data["json"] = &JSONResponse{Code: 400002, Msg: err.Error()}
		c.ServeJSON()
		return
	}

	// 判断邮箱是否重复
	if models.IsManagerExists(map[string]interface{}{"email": manager.Email, "id__gt": id}) {
		addLog(c.Ctx, logContent, "邮箱已被占用", "{\"code\": 400003, \"msg\": \"邮箱已被占用\"}")
		c.Data["json"] = &JSONResponse{Code: 400003, Msg: "邮箱已被占用"}
		c.ServeJSON()
		return
	}

	if models.IsManagerExists(map[string]interface{}{"email": manager.Email, "id__lt": id}) {
		addLog(c.Ctx, logContent, "邮箱已被占用", "{\"code\": 400003, \"msg\": \"邮箱已被占用\"}")
		c.Data["json"] = &JSONResponse{Code: 400003, Msg: "邮箱已被占用"}
		c.ServeJSON()
		return
	}

	if _, err := models.UpdateManager(map[string]interface{}{"id": id}, map[string]interface{}{
		"nickname": manager.Nickname,
		"email":    manager.Email,
		"is_admin": manager.IsAdmin,
	}); err != nil {
		addLog(c.Ctx, logContent, err.Error(), "{\"code\": 400004, \"msg\": \"操作失败\"}")
		c.Data["json"] = &JSONResponse{Code: 400004, Msg: "操作失败"}
		c.ServeJSON()
		return
	}

	addLog(c.Ctx, logContent, "", "{\"code\": 200, \"msg\": \"OK\"}")
	c.Data["json"] = &JSONResponse{Code: 200, Msg: "OK"}
	c.ServeJSON()
}

// Delete 删除用户
func (c *ManagerController) Delete() {
	id, _ := c.GetInt("id")

	manager, _ := models.GetOneManager(map[string]interface{}{"id": id})
	logContent := "删除用户，用户名：" + manager.Username

	if manager.Id == 0 {
		addLog(c.Ctx, logContent, "用户不存在", "{\"code\": 400000, \"msg\": \"用户不存在\"}")
		c.Data["json"] = &JSONResponse{Code: 400000, Msg: "用户不存在"}
		c.ServeJSON()
		return
	}

	if _, err := models.DeleteManager(map[string]interface{}{"id": id}); err != nil {
		addLog(c.Ctx, logContent, err.Error(), "{\"code\": 400001, \"msg\": \"操作失败\"}")
		c.Data["json"] = &JSONResponse{Code: 400001, Msg: "操作失败"}
		c.ServeJSON()
		return
	}

	addLog(c.Ctx, logContent, "", "{\"code\": 200, \"msg\": \"OK\"}")
	c.Data["json"] = &JSONResponse{Code: 200, Msg: "OK"}
	c.ServeJSON()
}
