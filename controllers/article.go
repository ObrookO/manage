package controllers

import (
	"errors"
	"fmt"
	"html/template"
	"manage/models"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/validation"

	utils2 "github.com/ObrookO/go-utils"

	"github.com/astaxie/beego/utils"
)

type ArticleController struct {
	BaseController
}

var (
	articlePageLimit  = 20                                              // 每页显示的文章数
	allowFileType     = []string{"cover", "content"}                    // 允许的文件的key
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
	category, _ := c.GetInt("c")
	if category > 0 {
		filter["category_id"] = category
	}
	isScroll, _ := c.GetInt("s", -1)
	if isScroll > -1 {
		filter["is_scroll"] = isScroll
	}
	allowComment, _ := c.GetInt("ac", -1)
	if allowComment > -1 && utils2.ObjInIntSlice(allowComment, []int{0, 1}) {
		filter["allow_comment"] = allowComment
	}
	isRecommend, _ := c.GetInt("r", -1)
	if isRecommend > -1 {
		filter["is_recommend"] = isRecommend
	}
	keyword := c.GetString("k")
	if len(keyword) > 0 {
		filter["title__icontains"] = keyword
	}

	AddLog(c.Ctx, "查看文章列表", "", "PAGE", "SUCCESS")

	categories, _ := models.GetCategories(nil)
	articles, _ := models.GetArticles(filter, 0, articlePageLimit)

	c.Data = map[interface{}]interface{}{
		"xsrfdata":   template.HTML(c.XSRFFormHTML()),
		"c":          category,
		"s":          isScroll,
		"ac":         allowComment,
		"r":          isRecommend,
		"keyword":    keyword,
		"categories": categories,
		"articles":   articles,
	}
}

// AddArticle 添加文章页面
func (c *ArticleController) AddArticle() {
	c.Layout = "layouts/master.html"
	c.TplName = "article/add.html"
	c.LayoutSections = map[string]string{
		"Style":  "article/add_style.html",
		"Script": "article/add_script.html",
	}

	categories, _ := models.GetCategories(nil)
	tags, _ := models.GetTags(nil)

	c.Data = map[interface{}]interface{}{
		"xsrfdata":   template.HTML(c.XSRFFormHTML()),
		"categories": categories,
		"tags":       tags,
	}
}

// UploadImage 上传图片
func (c *ArticleController) UploadImage() {
	c.EnableRender = false

	fileType := c.GetString("type")
	if !utils.InSlice(fileType, allowFileType) {
		AddLog(c.Ctx, "上传图片", "图片类型错误", "{code: 400000, msg: \"图片上传失败\"}", "FAIL")
		c.Data["json"] = &JSONResponse{Code: 400000, Msg: "图片上传失败"}
		c.ServeJSON()
		return
	}

	fileKey := "file"
	file, header, err := c.GetFile(fileKey)

	if err != nil {
		AddLog(c.Ctx, "上传图片", err.Error(), "{code: 400001, msg: \"图片上传失败\"}", "FAIL")
		c.Data["json"] = &JSONResponse{Code: 400001, Msg: "图片上传失败"}
		c.ServeJSON()
		return
	}

	defer file.Close()

	// 判断文件后缀
	ext := strings.ToLower(path.Ext(header.Filename))
	if !utils.InSlice(ext, allowImageExt) {
		AddLog(c.Ctx, "上传图片 "+header.Filename, "只允许上传png，gif，jpg，jpeg格式的图片", "{code: 400002, msg: \"只允许上传png，gif，jpg，jpeg格式的图片\"}", "FAIL")
		c.Data["json"] = &JSONResponse{Code: 400002, Msg: "只允许上传png，gif，jpg，jpeg格式的图片"}
		c.ServeJSON()
		return
	}

	// 判断文件大小
	size := header.Size
	if int(size) > allowImageSize {
		AddLog(c.Ctx, "上传图片 "+header.Filename, "只允许上传2M以下的图片", "{code: 400003, msg: \"只允许上传2M以下的图片\"}", "FAIL")
		c.Data["json"] = &JSONResponse{Code: 400003, Msg: "只允许上传2M以下的图片"}
		c.ServeJSON()
		return
	}

	filename := filenamePrefixMap[fileType] + time.Now().Format("20060102150405") + utils2.RandomStr(10) + ext
	if err := c.SaveToFile(fileKey, "static/upload/"+filename); err != nil {
		AddLog(c.Ctx, "上传图片 "+header.Filename, err.Error(), "{code: 400004, msg: \"图片上传失败\"}", "FAIL")
		c.Data["json"] = &JSONResponse{Code: 400004, Msg: "图片上传失败"}
		c.ServeJSON()
		return
	}

	AddLog(c.Ctx, "上传图片 "+header.Filename+" 保存名称为 "+filename, "", "{code: 200, msg: \"OK\"}", "SUCCESS")
	c.Data["json"] = &JSONResponse{Code: 200, Msg: "OK", Data: filename}
	c.ServeJSON()
}

// Post 添加文章
func (c *ArticleController) Post() {
	c.EnableRender = false

	article := &models.Article{}
	if err := c.ParseForm(article); err != nil {
		AddLog(c.Ctx, "添加文章", err.Error(), "{code: 400000, msg: \"添加文章失败\"}", "FAIL")
		c.Data["json"] = &JSONResponse{Code: 400000, Msg: "添加文章失败"}
		c.ServeJSON()
		return
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
		AddLog(c.Ctx, "添加文章", err.Error(), "{code: 400001, msg: \""+err.Error()+"\"}", "FAIL")
		c.Data["json"] = &JSONResponse{Code: 400001, Msg: err.Error()}
		c.ServeJSON()
		return
	}

	// 验证栏目是否存在
	if !models.IsCategoryExists(map[string]interface{}{"id": article.Category.Id}) {
		AddLog(c.Ctx, "添加文章", "栏目不存在", "{code: 400002, msg: \"栏目不存在\"}", "FAIL")
		c.Data["json"] = &JSONResponse{Code: 400002, Msg: "栏目不存在"}
		c.ServeJSON()
		return
	}

	// 验证封面是否存在
	if _, err := os.Stat("static/upload/" + article.Cover); err != nil {
		if os.IsNotExist(err) {
			AddLog(c.Ctx, "添加文章", "封面不存在", "{code: 400004, msg: \"封面不存在\"}", "FAIL")
			c.Data["json"] = &JSONResponse{Code: 400004, Msg: "封面不存在"}
			c.ServeJSON()
			return
		}
	}

	// 判断内容是否为空
	if len(strings.TrimLeft(article.Content, "")) == 0 {
		AddLog(c.Ctx, "添加文章", "内容不能为空", "{code: 400005, msg: \"内容不能为空\"}", "FAIL")
		c.Data["json"] = &JSONResponse{Code: 400005, Msg: "内容不能为空"}
		c.ServeJSON()
		return
	}

	// 判断是否是草稿
	if isDraft, _ := c.GetInt("is_draft"); isDraft != 1 {
		article.Status = 1
	}

	// 管理员信息
	article.Manager = &models.Manager{
		Id: ManagerInfo["uid"].(int),
	}

	if _, err := models.AddArticle(article); err != nil {
		AddLog(c.Ctx, "添加文章", err.Error(), "{code: 400006, msg: \""+err.Error()+"\"}", "FAIL")
		c.Data["json"] = &JSONResponse{Code: 400006, Msg: "添加文章失败"}
		c.ServeJSON()
		return
	}

	AddLog(c.Ctx, "添加文章", "", "{code: 200, msg: \"OK\"}", "SUCCESS")
	c.Data["json"] = &JSONResponse{Code: 200, Msg: "OK"}
	c.ServeJSON()
}

// Draft 草稿箱
func (c *ArticleController) Draft() {
	c.Layout = "layouts/master.html"
	c.TplName = "article/index.html"
}

// validData
func validData(data *models.Article) error {
	v := validation.Validation{}
	b, err := v.Valid(data)
	if err != nil {
		return err
	}

	if !b {
		return errors.New(fmt.Sprintf("%s %s", v.Errors[0].Field, v.Errors[0].Message))
	}

	return nil
}
