package controllers

import (
	"html/template"
	"manage/models"
	"path"
	"strings"
	"time"

	utils2 "github.com/ObrookO/go-utils"

	"github.com/astaxie/beego/utils"
)

type ArticleController struct {
	BaseController
}

var (
	articlePageLimit = 20
	allowImageExt    = []string{".png", ".jpg", ".jpeg", ".gif"}
	allowImageSize   = 1 << 21 // 2M
)

// Get 文章列表
func (c *ArticleController) Get() {
	c.Layout = "layouts/master.html"
	c.TplName = "article/index.html"

	AddLog(c.Ctx, "查看文章列表", "", "PAGE", "SUCCESS")

	categories, _ := models.GetCategories(nil)

	filter := map[string]interface{}{}
	articles, _ := models.GetArticles(filter, 0, articlePageLimit)

	c.Data = map[interface{}]interface{}{
		"xsrfdata":   template.HTML(c.XSRFFormHTML()),
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

func (c *ArticleController) UploadCover() {
	c.EnableRender = false

	fileKey := "cover"
	file, header, err := c.GetFile(fileKey)
	if err != nil {
		AddLog(c.Ctx, "上传封面图", err.Error(), "{code: 400000, msg: \"封面上传失败\"}", "FAIL")
		c.Data["json"] = &JSONResponse{Code: 400000, Msg: "封面上传失败"}
		c.ServeJSON()
		return
	}

	defer file.Close()

	// 判断文件后缀
	ext := strings.ToLower(path.Ext(header.Filename))
	if !utils.InSlice(ext, allowImageExt) {
		AddLog(c.Ctx, "上传封面图 "+header.Filename, "只允许上传png，gif，jpg，jpeg格式的图片", "{code: 400001, msg: \"只允许上传png，gif，jpg，jpeg格式的图片\"}", "FAIL")
		c.Data["json"] = &JSONResponse{Code: 400001, Msg: "只允许上传png，gif，jpg，jpeg格式的图片"}
		c.ServeJSON()
		return
	}

	// 判断文件大小
	size := header.Size
	if int(size) > allowImageSize {
		AddLog(c.Ctx, "上传封面图 "+header.Filename, "只允许上传2M以下的图片", "{code: 400002, msg: \"只允许上传2M以下的图片\"}", "FAIL")
		c.Data["json"] = &JSONResponse{Code: 400002, Msg: "只允许上传2M以下的图片"}
		c.ServeJSON()
		return
	}

	filename := "C" + time.Now().Format("20060102150405") + utils2.RandomStr(10) + ext
	if err := c.SaveToFile(fileKey, "static/upload/"+filename); err != nil {
		AddLog(c.Ctx, "上传封面图 "+header.Filename, "文件上传失败", "{code: 400003, msg: \"文件上传失败\"}", "FAIL")
		c.Data["json"] = &JSONResponse{Code: 400003, Msg: "文件上传失败"}
		c.ServeJSON()
		return
	}

	AddLog(c.Ctx, "上传封面图 "+header.Filename+" 保存名称为 "+filename, "", "{code: 200, msg: \"OK\"}", "SUCCESS")
	c.Data["json"] = &JSONResponse{Code: 200, Msg: "OK", Data: filename}
	c.ServeJSON()
}

func (c *ArticleController) Post() {
	c.EnableRender = false
}

// Draft 草稿箱
func (c *ArticleController) Draft() {
	c.Layout = "layouts/master.html"
	c.TplName = "article/index.html"
}
