package controllers

import (
	"html/template"
	"manage/models"
	"time"
)

type TagController struct {
	BaseController
}

var tagNameMaxLength = 10

// Get 查看标签列表
func (c *TagController) Get() {
	c.Layout = "layouts/master.html"
	c.TplName = "tag/index.html"
	c.LayoutSections = map[string]string{
		"Style":  "tag/style.html",
		"Script": "tag/script.html",
	}

	AddLog(c.Ctx, "查看标签列表", "", "PAGE", "SUCCESS")

	tags, _ := models.GetTags(nil)

	c.Data = map[interface{}]interface{}{
		"xsrfdata": template.HTML(c.XSRFFormHTML()),
		"tags":     tags,
	}
}

// Post 添加标签
func (c *TagController) Post() {
	c.EnableRender = false

	name := c.GetString("name")
	if len(name) == 0 || len([]rune(name)) > tagNameMaxLength {
		AddLog(c.Ctx, "添加标签 "+name, "名称的长度为0-10", "{\"code\": 400000, \"msg\": \"名称的长度为0-10\"}", "FAIL")
		c.Data["json"] = &JSONResponse{Code: 400000, Msg: "名称的长度为0-10"}
		c.ServeJSON()
		return
	}

	if models.IsTagExists(map[string]interface{}{"name": name}) {
		AddLog(c.Ctx, "添加标签 "+name, "标签已存在", "{\"code\": 400001, \"msg\": \"标签已存在\"}", "FAIL")
		c.Data["json"] = &JSONResponse{Code: 400001, Msg: "标签已存在"}
		c.ServeJSON()
		return
	}

	tagId, err := models.AddTag(models.Tag{Name: name})
	if err != nil {
		AddLog(c.Ctx, "添加标签 "+name, err.Error(), "{\"code\": 400002, \"msg\": \"标签添加失败\"}", "FAIL")
		c.Data["json"] = &JSONResponse{Code: 400002, Msg: "标签添加失败"}
		c.ServeJSON()
		return
	}

	AddLog(c.Ctx, "添加标签 "+name, "", "{\"code\": 200, \"msg\": \"OK\"}", "SUCCESS")
	c.Data["json"] = &JSONResponse{Code: 200, Msg: "OK", Data: tagId}
	c.ServeJSON()
}

// DeleteTag 删除标签
func (c *TagController) DeleteTag() {
	c.EnableRender = false

	id, _ := c.GetInt("id")
	tag, _ := models.GetOneTag(map[string]interface{}{"id": id})
	if tag.Id == 0 {
		AddLog(c.Ctx, "删除标签 "+tag.Name, "", "{\"code\": 400000, \"msg\": \"标签不存在\"}", "FAIL")
		c.Data["json"] = &JSONResponse{Code: 400000, Msg: "标签不存在"}
		c.ServeJSON()
		return
	}

	if tag.ArticleNum != 0 {
		AddLog(c.Ctx, "删除标签 "+tag.Name, "", "{\"code\": 400001, \"msg\": \"标签下有文章，不能删除\"}", "FAIL")
		c.Data["json"] = &JSONResponse{Code: 400001, Msg: "标签下有文章，不能删除"}
		c.ServeJSON()
		return
	}

	if _, err := models.DeleteTag(map[string]interface{}{"id": id}); err != nil {
		AddLog(c.Ctx, "删除标签 "+tag.Name, err.Error(), "{\"code\": 400002, \"msg\": \"标签删除失败\"}", "FAIL")
		c.Data["json"] = &JSONResponse{Code: 400002, Msg: "标签删除失败"}
		c.ServeJSON()
		return
	}

	AddLog(c.Ctx, "删除标签 "+tag.Name, "", "{\"code\": 200, \"msg\": \"OK\"}", "SUCCESS")
	c.Data["json"] = &JSONResponse{Code: 200, Msg: "OK"}
	c.ServeJSON()
}

// UpdateTag 编辑标签
func (c *TagController) UpdateTag() {
	c.EnableRender = false

	id, _ := c.GetInt("id")
	name := c.GetString("name")

	tag, _ := models.GetOneTag(map[string]interface{}{"id": id})
	if tag.Id == 0 {
		AddLog(c.Ctx, "修改标签名称", "标签不存在", "{\"code\": 400000, \"msg\": \"标签不存在\"}", "FAIL")
		c.Data["json"] = &JSONResponse{Code: 400000, Msg: "标签不存在"}
		c.ServeJSON()
		return
	}

	if len(name) == 0 || len([]rune(name)) > tagNameMaxLength {
		AddLog(c.Ctx, "修改标签 "+tag.Name+" 的名称为 "+name, "名称的长度为0-10", "{\"code\": 400000, \"msg\": \"名称的长度为0-10\"}", "FAIL")
		c.Data["json"] = &JSONResponse{Code: 400001, Msg: "名称的长度为0-10"}
		c.ServeJSON()
		return
	}

	// 如果名字没有更改，直接返回
	if name == tag.Name {
		AddLog(c.Ctx, "修改标签 "+tag.Name+" 的名称为 "+name, "", "{\"code\": 200, \"msg\": \"OK\"}", "SUCCESS")
		c.Data["json"] = &JSONResponse{Code: 200, Msg: "OK"}
		c.ServeJSON()
		return
	}

	// 判断是否存在重名标签
	if models.IsTagExists(map[string]interface{}{"name": name}) {
		AddLog(c.Ctx, "修改标签 "+tag.Name+" 的名称为 "+name, "标签已存在", "{\"code\": 400002, \"msg\": \"标签已存在\"}", "FAIL")
		c.Data["json"] = &JSONResponse{Code: 400002, Msg: "标签已存在"}
		c.ServeJSON()
		return
	}

	if _, err := models.UpdateTag(map[string]interface{}{"id": id}, map[string]interface{}{
		"name":       name,
		"updated_at": time.Now().Format("2006-01-02 15:04:05"),
	}); err != nil {
		AddLog(c.Ctx, "修改标签 "+tag.Name+" 的名称为 "+name, err.Error(), "{\"code\": 400003, \"msg\": \"标签编辑失败\"}", "FAIL")
		c.Data["json"] = &JSONResponse{Code: 400003, Msg: "标签编辑失败"}
		c.ServeJSON()
		return
	}

	AddLog(c.Ctx, "修改标签 "+tag.Name+" 的名称为 "+name, "", "{\"code\": 200, \"msg\": \"OK\"}", "SUCCESS")
	c.Data["json"] = &JSONResponse{Code: 200, Msg: "OK"}
	c.ServeJSON()
}
