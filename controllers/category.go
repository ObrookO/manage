package controllers

import (
	"html/template"
	"manage/models"
)

type CategoryController struct {
	BaseController
}

func (c *CategoryController) Get() {
	c.Layout = "layouts/master.html"
	c.TplName = "category/index.html"
	c.LayoutSections = map[string]string{
		"Style":  "category/style.html",
		"Script": "category/script.html",
	}

	filter := map[string]interface{}{}
	categories, err := models.GetCategories(filter)

	if err != nil {
		AddLog(c.Ctx, "获取栏目列表", err.Error(), "page")
	} else {
		AddLog(c.Ctx, "获取栏目列表", "", "page")
	}

	c.Data = map[interface{}]interface{}{
		"xsrfdata":   template.HTML(c.XSRFFormHTML()),
		"categories": categories,
	}
}

func (c *CategoryController) Post() {
	c.EnableRender = false

	name := c.GetString("name")
	shortName := c.GetString("short")

	if len(name) == 0 || len(name) > 10 {
		c.Data["json"] = &JSONResponse{Code: 400000, Msg: "名称的长度为0-10"}
		c.ServeJSON()
		AddLog(c.Ctx, "添加栏目"+name, "名称的长度为0-10", "{\"code\": 400000, \"msg\": \"名称的长度为0-10\"}")
		return
	}

	if len(shortName) == 0 || len(shortName) > 10 {
		c.Data["json"] = &JSONResponse{Code: 400001, Msg: "简称的长度为0-10"}
		c.ServeJSON()
		AddLog(c.Ctx, "添加栏目"+name, "简称的长度为0-10", "{\"code\": 400001, \"msg\": \"简称的长度为0-10\"}")
		return
	}

	if models.IsCategoryExists(map[string]interface{}{"name": name}) {
		c.Data["json"] = &JSONResponse{Code: 400002, Msg: "栏目已存在"}
		c.ServeJSON()
		AddLog(c.Ctx, "添加栏目"+name, "栏目已存在", "{\"code\": 400002, \"msg\": \"栏目已存在\"}")
		return
	}

	if _, err := models.AddCategory(&models.Category{Name: name, ShortName: shortName}); err != nil {
		c.Data["json"] = &JSONResponse{Code: 400003, Msg: "栏目添加失败"}
		c.ServeJSON()
		AddLog(c.Ctx, "添加栏目"+name, err.Error(), "{\"code\": 400003, \"msg\": \"栏目添加失败\"}")
		return
	}

	c.Data["json"] = &JSONResponse{Code: 200, Msg: "OK"}
	c.ServeJSON()
	AddLog(c.Ctx, "添加栏目"+name, "", "{\"code\": 200, \"msg\": \"OK\"}")
}

func (c *CategoryController) DeleteCategory() {
	c.EnableRender = false

	id, _ := c.GetInt("id")
	category, _ := models.GetCategory(map[string]interface{}{"id": id})
	if category.Id == 0 {
		c.Data["json"] = &JSONResponse{Code: 400000, Msg: "栏目不存在"}
		c.ServeJSON()
		AddLog(c.Ctx, "删除栏目"+category.Name, "", "{\"code\": 400000, \"msg\": \"栏目不存在\"}")
		return
	}

	if category.ArticleNum != 0 {
		c.Data["json"] = &JSONResponse{Code: 400001, Msg: "栏目下有文章，不能删除"}
		c.ServeJSON()
		AddLog(c.Ctx, "删除栏目"+category.Name, "", "{\"code\": 400001, \"msg\": \"栏目下有文章，不能删除\"}")
		return
	}

	if _, err := models.DeleteCategory(map[string]interface{}{"id": id}); err != nil {
		c.Data["json"] = &JSONResponse{Code: 400002, Msg: "栏目删除失败"}
		c.ServeJSON()
		AddLog(c.Ctx, "删除栏目"+category.Name, err.Error(), "{\"code\": 400002, \"msg\": \"栏目删除失败\"}")
		return
	}

	c.Data["json"] = &JSONResponse{Code: 200, Msg: "OK"}
	c.ServeJSON()
	AddLog(c.Ctx, "删除栏目"+category.Name, "", "{\"code\": 200, \"msg\": \"OK\"}")
}
