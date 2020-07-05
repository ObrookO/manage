package controllers

import (
	"fmt"
	"manage/models"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/utils"

	utils2 "github.com/ObrookO/go-utils"
)

type ArticleController struct {
	BaseController
}

var (
	filenamePrefixMap = map[string]string{"cover": "C", "content": "D"} // 图片的前缀
	allowImageExt     = []string{".png", ".jpg", ".jpeg", ".gif"}       // 允许的图片类型
	allowImageSize    = 1 << 21                                         // 2M // 允许的图片大小
)

// Get 文章列表
func (c *ArticleController) Get() {
	c.Layout = "layouts/master.html"
	c.TplName = "article/index.html"
	c.LayoutSections = map[string]string{
		"Script": "article/index_script.html",
	}

	// 查询条件
	filter := map[string]interface{}{}
	if !IsAdmin {
		filter["manager_id"] = ManagerInfo.Id
	}

	category, _ := c.GetInt("c")
	if category > 0 {
		filter["category_id"] = category
	}
	isScroll, _ := c.GetInt("s", -1)
	if isScroll > -1 && utils2.ObjInIntSlice(isScroll, []int{0, 1}) {
		filter["is_scroll"] = isScroll
	}
	allowComment, _ := c.GetInt("ac", -1)
	if allowComment > -1 && utils2.ObjInIntSlice(allowComment, []int{0, 1}) {
		filter["allow_comment"] = allowComment
	}
	isRecommend, _ := c.GetInt("r", -1)
	if isRecommend > -1 && utils2.ObjInIntSlice(isRecommend, []int{0, 1}) {
		filter["is_recommend"] = isRecommend
	}
	status, _ := c.GetInt("st", -1)
	if status > -1 && utils2.ObjInIntSlice(status, []int{0, 1}) {
		filter["status"] = status
	}
	keyword := c.GetString("k")
	if len(keyword) > 0 {
		filter["title__icontains"] = keyword
	}

	AddLog(c.Ctx, "查看文章列表", "", "PAGE")

	categories, _ := models.GetAllCategories(nil)
	articles, _ := models.GetAllArticles(filter)

	c.Data["c"] = category
	c.Data["s"] = isScroll
	c.Data["ac"] = allowComment
	c.Data["r"] = isRecommend
	c.Data["st"] = status
	c.Data["keyword"] = keyword
	c.Data["categories"] = categories
	c.Data["articles"] = articles
}

// ChangeStatus 修改文章状态
func (c *ArticleController) ChangeStatus() {
	id, _ := c.GetInt("id")
	status, _ := c.GetInt("status", -1)

	allowStatusSlice := []int{0, 1}
	statusMap := map[int]string{0: "草稿", 1: "正常"}

	// 判断状态
	if !utils2.ObjInIntSlice(status, allowStatusSlice) {
		AddLog(c.Ctx, "修改文章状态", "无效的status", "{\"code\": 400000, \"msg\": \"参数错误\"}")
		c.Data["json"] = &JSONResponse{Code: 400000, Msg: "参数错误"}
		c.ServeJSON()
		return
	}

	// 判断文章是否存在
	article, _ := models.GetOneArticle(map[string]interface{}{"id": id})
	if article.Id == 0 {
		AddLog(c.Ctx, "修改文章状态", "文章不存在", "{\"code\": 400001, \"msg\": \"文章不存在\"}")
		c.Data["json"] = &JSONResponse{Code: 400001, Msg: "文章不存在"}
		c.ServeJSON()
		return
	}

	logContent := "修改文章状态为：" + statusMap[status] + "，文章作者：" + article.Manager.Nickname + "，文章标题：" + article.Title

	// 判断权限
	if !IsAdmin {
		if ManagerInfo.Id != article.Manager.Id {
			AddLog(c.Ctx, logContent, "非法操作", "{\"code\": 500, \"msg\": \"非法操作\"}")
			c.Data["json"] = &JSONResponse{Code: 500, Msg: "非法操作"}
			c.ServeJSON()
			return
		}
	}

	// 更新文章状态
	if _, err := models.UpdateArticleWithFilter(map[string]interface{}{"id": id}, map[string]interface{}{
		"status":     status,
		"updated_at": time.Now().Format("2006-01-02 15:04:05"),
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

// Delete 删除文章
func (c *ArticleController) Delete() {
	id, _ := c.GetInt("id")
	article, _ := models.GetOneArticle(map[string]interface{}{"id": id})
	if article.Id == 0 {
		AddLog(c.Ctx, "删除文章", "文章不存在", "{\"code\": 400000, \"msg\": \"文章不存在\"}")
		c.Data["json"] = &JSONResponse{Code: 400000, Msg: "文章不存在"}
		c.ServeJSON()
		return
	}

	logContent := "删除文章，文章作者：" + article.Manager.Nickname + "，文章标题： " + article.Title

	// 判断权限
	if !IsAdmin {
		if ManagerInfo.Id != article.Manager.Id {
			AddLog(c.Ctx, logContent, "非法操作", "{\"code\": 500, \"msg\": \"非法操作\"}")
			c.Data["json"] = &JSONResponse{Code: 500, Msg: "非法操作"}
			c.ServeJSON()
			return
		}
	}

	if _, err := models.DeleteArticle(map[string]interface{}{"id": id}); err != nil {
		AddLog(c.Ctx, logContent, err.Error(), "{\"code\": 400001, \"msg\": \"文章删除失败\"}")
		c.Data["json"] = &JSONResponse{Code: 400001, Msg: "文章删除失败"}
		c.ServeJSON()
		return
	}

	// 删除封面图
	os.Remove("static/upload/" + article.Cover)
	AddLog(c.Ctx, logContent, "", "{\"code\": 200, \"msg\": \"OK\"}")
	c.Data["json"] = &JSONResponse{Code: 200, Msg: "OK"}
	c.ServeJSON()
}

// UploadImage 上传图片
func (c *ArticleController) UploadImage() {
	fileKey := "file"
	file, header, err := c.GetFile(fileKey)
	if err != nil {
		AddLog(c.Ctx, "上传图片", err.Error(), "{\"code\": 400000, \"msg\": \"图片上传失败\"}")
		c.Data["json"] = &JSONResponse{Code: 400000, Msg: "图片上传失败"}
		c.ServeJSON()
		return
	}

	defer file.Close()

	logContent := "上传图片，原始名称：" + header.Filename
	fileType := c.GetString("type")
	allowFileType := []string{"cover", "content"}

	if !utils.InSlice(fileType, allowFileType) {
		AddLog(c.Ctx, logContent, "无效的type", "{\"code\": 400001, \"msg\":\"参数错误\"}")
		c.Data["json"] = &JSONResponse{Code: 400001, Msg: "参数错误"}
		c.ServeJSON()
		return
	}

	// 判断文件后缀
	ext := strings.ToLower(path.Ext(header.Filename))
	if !utils.InSlice(ext, allowImageExt) {
		AddLog(c.Ctx, logContent, "只允许上传png，gif，jpg，jpeg格式的图片", "{\"code\": 400002, \"msg\":\"只允许上传png，gif，jpg，jpeg格式的图片\"}")
		c.Data["json"] = &JSONResponse{Code: 400002, Msg: "只允许上传png，gif，jpg，jpeg格式的图片"}
		c.ServeJSON()
		return
	}

	// 判断文件大小
	size := header.Size
	if int(size) > allowImageSize {
		AddLog(c.Ctx, logContent, "只允许上传2M以下的图片", "{\"code\": 400003, \"msg\":\"只允许上传2M以下的图片\"}")
		c.Data["json"] = &JSONResponse{Code: 400003, Msg: "只允许上传2M以下的图片"}
		c.ServeJSON()
		return
	}

	filename := filenamePrefixMap[fileType] + time.Now().Format("20060102150405") + utils2.RandomStr(10) + ext
	if err := c.SaveToFile(fileKey, "static/upload/"+filename); err != nil {
		AddLog(c.Ctx, logContent, err.Error(), "{\"code\": 400004, \"msg\":\"图片上传失败\"}")
		c.Data["json"] = &JSONResponse{Code: 400004, Msg: "图片上传失败"}
		c.ServeJSON()
		return
	}

	AddLog(c.Ctx, logContent+"，保存名称："+filename, "", "{\"code\": 200, \"msg\":\"OK\"}")
	c.Data["json"] = &JSONResponse{Code: 200, Msg: "OK", Data: filename}
	c.ServeJSON()
}

// Add 添加文章页面
func (c *ArticleController) Add() {
	c.Layout = "layouts/master.html"
	c.TplName = "article/add.html"
	c.LayoutSections = map[string]string{
		"Style":  "article/add_style.html",
		"Script": "article/add_script.html",
	}

	categories, _ := models.GetAllCategories(nil)
	tags, _ := models.GetTags(nil)

	c.Data["categories"] = categories
	c.Data["tags"] = tags
}

// Post 添加文章
func (c *ArticleController) Post() {
	article := models.Article{}
	if err := c.ParseForm(&article); err != nil {
		AddLog(c.Ctx, "添加文章", err.Error(), "{\"code\": 400000, \"msg\": \"操作失败\"}")
		c.Data["json"] = &JSONResponse{Code: 400000, Msg: "操作失败"}
		c.ServeJSON()
		return
	}

	logContent := "添加文章，文章标题：" + article.Title

	// 设置栏目
	categoryId, _ := c.GetInt("categoryId")
	article.Category = &models.Category{
		Id: categoryId,
	}

	// 设置标签
	tags := c.GetString("tags")
	for _, t := range strings.Split(tags, ",") {
		if models.IsTagExists(map[string]interface{}{"id": t}) {
			tagId, _ := strconv.Atoi(t)
			article.Tags = append(article.Tags, &models.Tag{Id: tagId})
		}
	}

	// 表单验证
	if err := validData(article); err != nil {
		AddLog(c.Ctx, logContent, err.Error(), "{\"code\": 400001, \"msg\": \""+err.Error()+"\"}")
		c.Data["json"] = &JSONResponse{Code: 400001, Msg: err.Error()}
		c.ServeJSON()
		return
	}

	// 验证栏目是否存在
	if !models.IsCategoryExists(map[string]interface{}{"id": article.Category.Id}) {
		AddLog(c.Ctx, logContent, "栏目不存在", "{\"code\": 400002, \"msg\": \"栏目不存在\"}")
		c.Data["json"] = &JSONResponse{Code: 400002, Msg: "栏目不存在"}
		c.ServeJSON()
		return
	}

	// 验证封面是否存在
	if _, err := os.Stat("static/upload/" + article.Cover); err != nil {
		if os.IsNotExist(err) {
			AddLog(c.Ctx, logContent, "封面不存在", "{\"code\": 400003, \"msg\": \"封面不存在\"}")
			c.Data["json"] = &JSONResponse{Code: 400003, Msg: "封面不存在"}
			c.ServeJSON()
			return
		}
	}

	// 判断内容是否为空
	if len(strings.TrimLeft(article.Content, "")) == 0 {
		AddLog(c.Ctx, logContent, "内容不能为空", "{\"code\": 400004, \"msg\": \"内容不能为空\"}")
		c.Data["json"] = &JSONResponse{Code: 400004, Msg: "内容不能为空"}
		c.ServeJSON()
		return
	}

	// 判断是否是草稿
	if isDraft, _ := c.GetInt("isDraft"); isDraft != 1 {
		article.Status = 1
	}

	// 管理员信息
	article.Manager = &models.Manager{
		Id: ManagerInfo.Id,
	}

	id, err := models.AddArticle(article)
	if err != nil {
		AddLog(c.Ctx, logContent, err.Error(), "{\"code\": 400005, \"msg\":\"操作失败\"}")
		c.Data["json"] = &JSONResponse{Code: 400005, Msg: "操作失败"}
		c.ServeJSON()
		return
	}

	AddLog(c.Ctx, logContent+fmt.Sprintf("文章id：%v，文章标题：%v", id, article.Title), "", "{\"code\": 200, \"msg\":\"OK\"}")
	c.Data["json"] = &JSONResponse{Code: 200, Msg: "OK"}
	c.ServeJSON()
}

// Edit 编辑文章页面
func (c *ArticleController) Edit() {
	c.Layout = "layouts/master.html"
	c.TplName = "article/edit.html"
	c.LayoutSections = map[string]string{
		"Style":  "article/edit_style.html",
		"Script": "article/edit_script.html",
	}

	id, _ := c.GetInt(":id")
	article, _ := models.GetOneArticle(map[string]interface{}{"id": id})
	if article.Id == 0 {
		c.Abort("404")
	}

	categories, _ := models.GetAllCategories(nil)
	tags, _ := models.GetTags(nil)

	c.Data["categories"] = categories
	c.Data["id"] = id
	c.Data["tags"] = tags
	c.Data["article"] = article
}

// Update 更新文章
func (c *ArticleController) Update() {
	article := models.Article{}
	if err := c.ParseForm(&article); err != nil {
		AddLog(c.Ctx, "更新文章", err.Error(), "{\"code\": 400000, \"msg\":\"更新文章失败\"}")
		c.Data["json"] = &JSONResponse{Code: 400000, Msg: "更新文章失败"}
		c.ServeJSON()
		return
	}

	// 判断文章是否存在
	if !models.IsArticleExists(map[string]interface{}{"id": article.Id}) {
		AddLog(c.Ctx, "更新文章", "文章不存在", "{\"code\": 400001, \"msg\":\"文章不存在\"}")
		c.Data["json"] = &JSONResponse{Code: 400001, Msg: "文章不存在"}
		c.ServeJSON()
		return
	}

	logContent := "更新文章，文章作者：" + article.Manager.Nickname + "，文章标题：" + article.Title

	// 判断权限
	if !IsAdmin {
		if ManagerInfo.Id != article.Manager.Id {
			AddLog(c.Ctx, logContent, "非法操作", "{\"code\": 500, \"msg\":\"非法操作\"}")
			c.Data["json"] = &JSONResponse{Code: 500, Msg: "非法操作"}
			c.ServeJSON()
			return
		}
	}

	// 设置栏目
	categoryId, _ := c.GetInt("categoryId")
	article.Category = &models.Category{
		Id: categoryId,
	}

	// 设置标签
	tags := c.GetString("tags")
	for _, t := range strings.Split(tags, ",") {
		if models.IsTagExists(map[string]interface{}{"id": t}) {
			tagId, _ := strconv.Atoi(t)
			article.Tags = append(article.Tags, &models.Tag{Id: tagId})
		}
	}

	// 表单验证
	if err := validData(article); err != nil {
		AddLog(c.Ctx, logContent, err.Error(), "{\"code\": 400002, \"msg\":\""+err.Error()+"\"}")
		c.Data["json"] = &JSONResponse{Code: 400002, Msg: err.Error()}
		c.ServeJSON()
		return
	}

	// 验证栏目是否存在
	if !models.IsCategoryExists(map[string]interface{}{"id": article.Category.Id}) {
		AddLog(c.Ctx, logContent, "栏目不存在", "{\"code\": 400003, \"msg\":\"栏目不存在\"}")
		c.Data["json"] = &JSONResponse{Code: 400003, Msg: "栏目不存在"}
		c.ServeJSON()
		return
	}

	// 验证封面是否存在
	if _, err := os.Stat("static/upload/" + article.Cover); err != nil {
		if os.IsNotExist(err) {
			AddLog(c.Ctx, logContent, "封面不存在", "{\"code\": 400004, \"msg\":\"封面不存在\"}")
			c.Data["json"] = &JSONResponse{Code: 400004, Msg: "封面不存在"}
			c.ServeJSON()
			return
		}
	}

	// 判断内容是否为空
	if len(strings.TrimLeft(article.Content, "")) == 0 {
		AddLog(c.Ctx, logContent, "内容不能为空", "{\"code\": 400005, \"msg\":\"内容不能为空\"}")
		c.Data["json"] = &JSONResponse{Code: 400005, Msg: "内容不能为空"}
		c.ServeJSON()
		return
	}

	if _, err := models.UpdateArticle(article, "title", "keyword", "category_id", "description", "cover", "content", "is_scroll", "is_recommend",
		"allow_comment"); err != nil {
		AddLog(c.Ctx, logContent, err.Error(), "{\"code\": 400006, \"msg\":\"操作失败\"}")
		c.Data["json"] = &JSONResponse{Code: 400006, Msg: "操作失败"}
		c.ServeJSON()
		return
	}

	AddLog(c.Ctx, logContent, "", "{\"code\": 200, \"msg\":\"OK\"}")
	c.Data["json"] = &JSONResponse{Code: 200, Msg: "OK"}
	c.ServeJSON()
	return
}
