package controllers

type CommentController struct {
	BaseController
}

func (c *CommentController) Get() {
	c.Layout = "layouts/master.html"
	c.TplName = "comment/index.html"
}

func (c *CommentController) Keyword() {
	c.Layout = "layouts/master.html"
	c.TplName = "comment/keyword.html"
}
