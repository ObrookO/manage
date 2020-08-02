package controllers

import (
	"fmt"
	"manage/models"
	"time"
)

type TagController struct {
	BaseController
}

var tagNameMaxLength = 30

// Get 查看标签列表
func (c *TagController) Get() {
	c.Layout = "layouts/master.html"
	c.TplName = "tag/index.html"
	c.LayoutSections = map[string]string{
		"Style":  "tag/index_style.html",
		"Script": "tag/index_script.html",
	}

	addLog(c.Ctx, "查看标签列表", "", "PAGE")

	tags, _ := models.GetTags(nil)

	c.Data["tags"] = tags
}

// Post 添加标签
func (c *TagController) Post() {
	name := c.GetString("name")
	logContent := "添加标签 " + name

	if len(name) == 0 || len([]rune(name)) > tagNameMaxLength {
		addLog(c.Ctx, logContent, fmt.Sprintf("名称的长度为0-%v", tagNameMaxLength), fmt.Sprintf("{\"code\": 400000, \"msg\": \"名称的长度为0-%v\"}", tagNameMaxLength))
		c.Data["json"] = &JSONResponse{Code: 400000, Msg: fmt.Sprintf("名称的长度为0-%v", tagNameMaxLength)}
		c.ServeJSON()
		return
	}

	if models.IsTagExists(map[string]interface{}{"name": name}) {
		addLog(c.Ctx, logContent, "标签已存在", "{\"code\": 400001, \"msg\": \"标签已存在\"}")
		c.Data["json"] = &JSONResponse{Code: 400001, Msg: "标签已存在"}
		c.ServeJSON()
		return
	}

	tagId, err := models.AddTag(models.Tag{Name: name})
	if err != nil {
		addLog(c.Ctx, logContent, err.Error(), "{\"code\": 400002, \"msg\": \"操作失败\"}")
		c.Data["json"] = &JSONResponse{Code: 400002, Msg: "操作失败"}
		c.ServeJSON()
		return
	}

	addLog(c.Ctx, logContent, "", "{\"code\": 200, \"msg\": \"OK\"}")
	c.Data["json"] = &JSONResponse{Code: 200, Msg: "OK", Data: tagId}
	c.ServeJSON()
}

// DeleteTag 删除标签
func (c *TagController) DeleteTag() {
	id, _ := c.GetInt("id")
	tag, _ := models.GetOneTag(map[string]interface{}{"id": id})
	logContent := "删除标签 " + tag.Name
	if tag.Id == 0 {
		addLog(c.Ctx, logContent, "", "{\"code\": 400000, \"msg\": \"标签不存在\"}")
		c.Data["json"] = &JSONResponse{Code: 400000, Msg: "标签不存在"}
		c.ServeJSON()
		return
	}

	if models.GetArticleNumOfTag(map[string]interface{}{"tag_id": tag.Id}) != 0 {
		addLog(c.Ctx, logContent, "", "{\"code\": 400001, \"msg\": \"标签下有文章，不能删除\"}")
		c.Data["json"] = &JSONResponse{Code: 400001, Msg: "标签下有文章，不能删除"}
		c.ServeJSON()
		return
	}

	if _, err := models.DeleteTag(map[string]interface{}{"id": id}); err != nil {
		addLog(c.Ctx, logContent, err.Error(), "{\"code\": 400002, \"msg\": \"操作失败\"}")
		c.Data["json"] = &JSONResponse{Code: 400002, Msg: "操作失败"}
		c.ServeJSON()
		return
	}

	addLog(c.Ctx, logContent, "", "{\"code\": 200, \"msg\": \"OK\"}")
	c.Data["json"] = &JSONResponse{Code: 200, Msg: "OK"}
	c.ServeJSON()
}

// UpdateTag 编辑标签
func (c *TagController) UpdateTag() {
	id, _ := c.GetInt("id")
	name := c.GetString("name")

	tag, _ := models.GetOneTag(map[string]interface{}{"id": id})
	if tag.Id == 0 {
		addLog(c.Ctx, "修改标签名称", "标签不存在", "{\"code\": 400000, \"msg\": \"标签不存在\"}")
		c.Data["json"] = &JSONResponse{Code: 400000, Msg: "标签不存在"}
		c.ServeJSON()
		return
	}

	logContent := "修改标签 " + tag.Name + " 的名称为 " + name

	if len(name) == 0 || len([]rune(name)) > tagNameMaxLength {
		addLog(c.Ctx, logContent, "名称的长度为0-10", "{\"code\": 400000, \"msg\": \"名称的长度为0-10\"}")
		c.Data["json"] = &JSONResponse{Code: 400001, Msg: "名称的长度为0-10"}
		c.ServeJSON()
		return
	}

	// 如果名字没有更改，直接返回
	if name == tag.Name {
		addLog(c.Ctx, logContent, "", "{\"code\": 200, \"msg\": \"OK\"}")
		c.Data["json"] = &JSONResponse{Code: 200, Msg: "OK"}
		c.ServeJSON()
		return
	}

	// 判断是否存在重名标签
	if models.IsTagExists(map[string]interface{}{"name": name}) {
		addLog(c.Ctx, logContent, "标签已存在", "{\"code\": 400002, \"msg\": \"标签已存在\"}")
		c.Data["json"] = &JSONResponse{Code: 400002, Msg: "标签已存在"}
		c.ServeJSON()
		return
	}

	if _, err := models.UpdateTagWithFilter(map[string]interface{}{"id": id}, map[string]interface{}{
		"name":       name,
		"updated_at": time.Now().Format("2006-01-02 15:04:05"),
	}); err != nil {
		addLog(c.Ctx, logContent, err.Error(), "{\"code\": 400003, \"msg\": \"标签编辑失败\"}")
		c.Data["json"] = &JSONResponse{Code: 400003, Msg: "标签编辑失败"}
		c.ServeJSON()
		return
	}

	addLog(c.Ctx, logContent, "", "{\"code\": 200, \"msg\": \"OK\"}")
	c.Data["json"] = &JSONResponse{Code: 200, Msg: "OK"}
	c.ServeJSON()
}
