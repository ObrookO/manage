package controllers

import "manage/models"

type ResourceController struct {
	BaseController
}

// Get 干货收藏列表
func (c *ResourceController) Get() {
	c.Layout = "layouts/master.html"
	c.TplName = "resource/index.html"
	c.LayoutSections = map[string]string{
		"Style":  "resource/index_style.html",
		"Script": "resource/index_script.html",
	}

	addLog(c.Ctx, "查看干货收藏列表", "", "PAGE")

	list, _ := models.GetAllResource(nil)
	c.Data["list"] = list
}

// Post 添加干货收藏
func (c *ResourceController) Post() {
	var resource models.Resource

	if err := c.ParseForm(&resource); err != nil {
		addLog(c.Ctx, "添加干货收藏", err.Error(), "{\"code\": 400000, \"msg\": \"操作失败\"}")
		c.Data["json"] = &JSONResponse{Code: 400000, Msg: "操作失败"}
		c.ServeJSON()
		return
	}

	// 表单验证
	if err := validData(resource); err != nil {
		addLog(c.Ctx, "添加干货收藏", err.Error(), "{\"code\": 400001, \"msg\": \""+err.Error()+"\"}")
		c.Data["json"] = &JSONResponse{Code: 400001, Msg: err.Error()}
		c.ServeJSON()
		return
	}

	if _, err := models.AddResource(resource); err != nil {
		addLog(c.Ctx, "添加干货收藏", err.Error(), "{\"code\": 400002, \"msg\": \"操作失败\"}")
		c.Data["json"] = &JSONResponse{Code: 400002, Msg: "操作失败"}
		c.ServeJSON()
		return
	}

	addLog(c.Ctx, "添加干货收藏", "", "{\"code\": 200, \"msg\": \"OK\"}")
	c.Data["json"] = &JSONResponse{Code: 200, Msg: "OK"}
	c.ServeJSON()
	return
}

// DeleteResource 删除干货收藏
func (c *ResourceController) DeleteResource() {
	id, _ := c.GetInt("id")

	resource, _ := models.GetOneResource(map[string]interface{}{"id": id})
	if resource.Id == 0 {
		addLog(c.Ctx, "删除干货收藏", "干货收藏不存在", "{\"code\": 400000, \"msg\": \"干货收藏不存在\"}")
		c.Data["json"] = &JSONResponse{Code: 400000, Msg: "干货收藏不存在"}
		c.ServeJSON()
		return
	}

	logContent := "删除干货资源，标题：" + resource.Title

	if _, err := models.DeleteResource(map[string]interface{}{"id": id}); err != nil {
		addLog(c.Ctx, logContent, err.Error(), "{\"code\": 400001, \"msg\": \"操作失败\"}")
		c.Data["json"] = &JSONResponse{Code: 400001, Msg: "操作失败"}
		c.ServeJSON()
		return
	}

	addLog(c.Ctx, logContent, "", "{\"code\": 200, \"msg\": \"OK\"}")
	c.Data["json"] = &JSONResponse{Code: 200, Msg: "OK"}
	c.ServeJSON()
	return

}
