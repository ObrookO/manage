package controllers

type IndexController struct {
	BaseController
}

func (c *IndexController) Get() {
	c.Layout = "layouts/master.html"
	c.TplName = "index/index.html"
}
