package controllers

type ArticleController struct {
	BaseController
}

// Get 文章列表
func (c *ArticleController) Get() {
	c.Layout = "layouts/master.html"
	c.TplName = "article/index.html"

	//filter := map[string]interface{} {}
	//articles, _ := models.GetArticles(filter)
}

// Draft 草稿箱
func (c *ArticleController) Draft() {
	c.Layout = "layouts/master.html"
	c.TplName = "article/index.html"
}
