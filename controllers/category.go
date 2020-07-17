package controllers

import (
	"fmt"
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
		"Style":  "category/index_style.html",
		"Script": "category/index_script.html",
	}

	addLog(c.Ctx, "查看栏目列表", "", "PAGE")

	categories, _ := models.GetAllCategories(nil)

	c.Data["categories"] = categories
}

// Post 添加栏目
func (c *CategoryController) Post() {
	name := c.GetString("name")

	if len(name) == 0 || len([]rune(name)) > categoryNameMaxLength {
		addLog(c.Ctx, "添加栏目", fmt.Sprintf("名称的长度为0-%v", categoryNameMaxLength), fmt.Sprintf("{\"code\": 400000, \"msg\": \"名称的长度为0-%v\"}",
			categoryNameMaxLength))
		c.Data["json"] = &JSONResponse{Code: 400000, Msg: fmt.Sprintf("名称的长度为0-%v", categoryNameMaxLength)}
		c.ServeJSON()
		return
	}

	logContent := "添加栏目，栏目名称：" + name

	if models.IsCategoryExists(map[string]interface{}{"name": name}) {
		addLog(c.Ctx, logContent, "栏目已存在", "{\"code\": 400001, \"msg\": \"栏目已存在\"}")
		c.Data["json"] = &JSONResponse{Code: 400001, Msg: "栏目已存在"}
		c.ServeJSON()
		return
	}

	categoryId, err := models.AddCategory(models.Category{Name: name})
	if err != nil {
		addLog(c.Ctx, logContent, err.Error(), "{\"code\": 400002, \"msg\": \"操作失败\"}")
		c.Data["json"] = &JSONResponse{Code: 400002, Msg: "操作失败"}
		c.ServeJSON()
		return
	}

	addLog(c.Ctx, logContent, "", fmt.Sprintf("{\"code\": 200, \"msg\": \"OK\", \"data\": %v }", categoryId))
	c.Data["json"] = &JSONResponse{Code: 200, Msg: "OK", Data: categoryId}
	c.ServeJSON()
}

// DeleteCategory 删除栏目
func (c *CategoryController) DeleteCategory() {
	id, _ := c.GetInt("id")
	category, _ := models.GetCategory(map[string]interface{}{"id": id})

	if category.Id == 0 {
		addLog(c.Ctx, "删除栏目", "", "{\"code\": 400000, \"msg\": \"栏目不存在\"}")
		c.Data["json"] = &JSONResponse{Code: 400000, Msg: "栏目不存在"}
		c.ServeJSON()
		return
	}

	logContent := "删除栏目，栏目名称：" + category.Name

	if models.GetArticleNumOfCategory(map[string]interface{}{"category_id": category.Id}) != 0 {
		addLog(c.Ctx, logContent, "", "{\"code\": 400001, \"msg\": \"栏目下有文章，不能删除\"}")
		c.Data["json"] = &JSONResponse{Code: 400001, Msg: "栏目下有文章，不能删除"}
		c.ServeJSON()
		return
	}

	if _, err := models.DeleteCategory(map[string]interface{}{"id": id}); err != nil {
		addLog(c.Ctx, logContent, err.Error(), "{\"code\": 400002, \"msg\": \"操作失败\"}")
		c.Data["json"] = &JSONResponse{Code: 400002, Msg: "操作失败"}
		c.ServeJSON()
		return
	}

	addLog(c.Ctx, logContent, "", "{\"code\": 200, \"msg\": \"OK\"}")
	c.Data["json"] = &JSONResponse{Code: 200, Msg: "OK"}
	c.ServeJSON()
}

// UpdateCategory 编辑栏目
func (c *CategoryController) UpdateCategory() {
	id, _ := c.GetInt("id")
	name := c.GetString("name")

	category, _ := models.GetCategory(map[string]interface{}{"id": id})
	if category.Id == 0 {
		addLog(c.Ctx, "修改栏目名称", "栏目不存在", "{\"code\": 400000, \"msg\": \"栏目不存在\"}")
		c.Data["json"] = &JSONResponse{Code: 400000, Msg: "栏目不存在"}
		c.ServeJSON()
		return
	}

	logContent := "修改栏目名称，原始名称：" + category.Name + "，新名称：" + name
	if len(name) == 0 || len([]rune(name)) > categoryNameMaxLength {
		addLog(c.Ctx, logContent, "名称的长度为0-10", "{\"code\": 400000, \"msg\": \"名称的长度为0-10\"}")
		c.Data["json"] = &JSONResponse{Code: 400001, Msg: "名称的长度为0-10"}
		c.ServeJSON()
		return
	}

	// 如果名字没有更改，直接返回
	if name == category.Name {
		addLog(c.Ctx, logContent, "", "{\"code\": 200, \"msg\": \"OK\"}")
		c.Data["json"] = &JSONResponse{Code: 200, Msg: "OK"}
		c.ServeJSON()
		return
	}

	// 判断是否存在重名栏目
	if models.IsCategoryExists(map[string]interface{}{"name": name}) {
		addLog(c.Ctx, logContent, "栏目已存在", "{\"code\": 400002, \"msg\": \"栏目已存在\"}")
		c.Data["json"] = &JSONResponse{Code: 400002, Msg: "栏目已存在"}
		c.ServeJSON()
		return
	}

	if _, err := models.UpdateCategoryWithFilter(map[string]interface{}{"id": id}, map[string]interface{}{
		"name":       name,
		"updated_at": time.Now().Format("2006-01-02 15:04:05"),
	}); err != nil {
		addLog(c.Ctx, logContent, err.Error(), "{\"code\": 400003, \"msg\": \"操作失败\"}")
		c.Data["json"] = &JSONResponse{Code: 400003, Msg: "操作失败"}
		c.ServeJSON()
		return
	}

	addLog(c.Ctx, logContent, "", "{\"code\": 200, \"msg\": \"OK\"}")
	c.Data["json"] = &JSONResponse{Code: 200, Msg: "OK"}
	c.ServeJSON()
}
