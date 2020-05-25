package controllers

import (
	"html/template"
	"manage/models"
	"time"
)

type CategoryController struct {
	BaseController
}

var categoryNameMaxLength = 10

func (c *CategoryController) Get() {
	c.Layout = "layouts/master.html"
	c.TplName = "category/index.html"
	c.LayoutSections = map[string]string{
		"Style":  "category/style.html",
		"Script": "category/script.html",
	}

	AddLog(c.Ctx, "查看栏目列表", "", "PAGE", "SUCCESS")

	categories, _ := models.GetCategories(nil)

	c.Data = map[interface{}]interface{}{
		"xsrfdata":   template.HTML(c.XSRFFormHTML()),
		"categories": categories,
	}
}

// Post 添加栏目
func (c *CategoryController) Post() {
	c.EnableRender = false

	name := c.GetString("name")
	if len(name) == 0 || len([]rune(name)) > categoryNameMaxLength {
		AddLog(c.Ctx, "添加栏目 "+name, "名称的长度为0-10", "{\"code\": 400000, \"msg\": \"名称的长度为0-10\"}", "FAIL")
		c.Data["json"] = &JSONResponse{Code: 400000, Msg: "名称的长度为0-10"}
		c.ServeJSON()
		return
	}

	if models.IsCategoryExists(map[string]interface{}{"name": name}) {
		AddLog(c.Ctx, "添加栏目 "+name, "栏目已存在", "{\"code\": 400001, \"msg\": \"栏目已存在\"}", "FAIL")
		c.Data["json"] = &JSONResponse{Code: 400001, Msg: "栏目已存在"}
		c.ServeJSON()
		return
	}

	if _, err := models.AddCategory(models.Category{Name: name}); err != nil {
		AddLog(c.Ctx, "添加栏目 "+name, err.Error(), "{\"code\": 400002, \"msg\": \"栏目添加失败\"}", "FAIL")
		c.Data["json"] = &JSONResponse{Code: 400002, Msg: "栏目添加失败"}
		c.ServeJSON()
		return
	}

	AddLog(c.Ctx, "添加栏目 "+name, "", "{\"code\": 200, \"msg\": \"OK\"}", "SUCCESS")
	c.Data["json"] = &JSONResponse{Code: 200, Msg: "OK"}
	c.ServeJSON()
}

// DeleteCategory 删除栏目
func (c *CategoryController) DeleteCategory() {
	c.EnableRender = false

	id, _ := c.GetInt("id")
	category, _ := models.GetCategory(map[string]interface{}{"id": id})
	if category.Id == 0 {
		AddLog(c.Ctx, "删除栏目 "+category.Name, "", "{\"code\": 400000, \"msg\": \"栏目不存在\"}", "FAIL")
		c.Data["json"] = &JSONResponse{Code: 400000, Msg: "栏目不存在"}
		c.ServeJSON()
		return
	}

	if category.ArticleNum != 0 {
		AddLog(c.Ctx, "删除栏目 "+category.Name, "", "{\"code\": 400001, \"msg\": \"栏目下有文章，不能删除\"}", "FAIL")
		c.Data["json"] = &JSONResponse{Code: 400001, Msg: "栏目下有文章，不能删除"}
		c.ServeJSON()
		return
	}

	if _, err := models.DeleteCategory(map[string]interface{}{"id": id}); err != nil {
		AddLog(c.Ctx, "删除栏目 "+category.Name, err.Error(), "{\"code\": 400002, \"msg\": \"栏目删除失败\"}", "FAIL")
		c.Data["json"] = &JSONResponse{Code: 400002, Msg: "栏目删除失败"}
		c.ServeJSON()
		return
	}

	AddLog(c.Ctx, "删除栏目 "+category.Name, "", "{\"code\": 200, \"msg\": \"OK\"}", "SUCCESS")
	c.Data["json"] = &JSONResponse{Code: 200, Msg: "OK"}
	c.ServeJSON()
}

// UpdateCategory 编辑栏目
func (c *CategoryController) UpdateCategory() {
	c.EnableRender = false

	id, _ := c.GetInt("id")
	name := c.GetString("name")

	category, _ := models.GetCategory(map[string]interface{}{"id": id})
	if category.Id == 0 {
		AddLog(c.Ctx, "修改栏目名称", "栏目不存在", "{\"code\": 400000, \"msg\": \"栏目不存在\"}", "FAIL")
		c.Data["json"] = &JSONResponse{Code: 400000, Msg: "栏目不存在"}
		c.ServeJSON()
		return
	}

	if len(name) == 0 || len([]rune(name)) > categoryNameMaxLength {
		AddLog(c.Ctx, "修改栏目 "+category.Name+" 的名称为 "+name, "名称的长度为0-10", "{\"code\": 400000, \"msg\": \"名称的长度为0-10\"}", "FAIL")
		c.Data["json"] = &JSONResponse{Code: 400001, Msg: "名称的长度为0-10"}
		c.ServeJSON()
		return
	}

	// 如果名字没有更改，直接返回
	if name == category.Name {
		AddLog(c.Ctx, "修改栏目 "+category.Name+" 的名称为 "+name, "", "{\"code\": 200, \"msg\": \"OK\"}", "SUCCESS")
		c.Data["json"] = &JSONResponse{Code: 200, Msg: "OK"}
		c.ServeJSON()
		return
	}

	// 判断是否存在重名栏目
	if models.IsCategoryExists(map[string]interface{}{"name": name}) {
		AddLog(c.Ctx, "修改栏目 "+category.Name+" 的名称为 "+name, "栏目已存在", "{\"code\": 400002, \"msg\": \"栏目已存在\"}", "FAIL")
		c.Data["json"] = &JSONResponse{Code: 400002, Msg: "栏目已存在"}
		c.ServeJSON()
		return
	}

	if _, err := models.UpdateCategory(map[string]interface{}{"id": id}, map[string]interface{}{
		"name":       name,
		"updated_at": time.Now().Format("2006-01-02 15:04:05"),
	}); err != nil {
		AddLog(c.Ctx, "修改栏目 "+category.Name+" 的名称为 "+name, err.Error(), "{\"code\": 400003, \"msg\": \"栏目编辑失败\"}", "FAIL")
		c.Data["json"] = &JSONResponse{Code: 400003, Msg: "栏目编辑失败"}
		c.ServeJSON()
		return
	}

	AddLog(c.Ctx, "修改栏目 "+category.Name+" 的名称为 "+name, "", "{\"code\": 200, \"msg\": \"OK\"}", "SUCCESS")
	c.Data["json"] = &JSONResponse{Code: 200, Msg: "OK"}
	c.ServeJSON()
}
