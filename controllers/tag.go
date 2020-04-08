package controllers

type TagController struct {
	BaseController
}

func (c *TagController) Get() {
	c.Layout = "layouts/master.html"
	c.TplName = "tag/index.html"
}
