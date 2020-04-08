package controllers

type ArticleController struct {
	BaseController
}

func (c *ArticleController) Get() {
	c.Layout = "layouts/master.html"
	c.TplName = "article/index.html"
}

func (c *ArticleController) Draft() {
	c.Layout = "layouts/master.html"
	c.TplName = "article/index.html"
}
