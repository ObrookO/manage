package controllers

type LogController struct {
	BaseController
}

func (c *LogController) HomeLog() {
	c.Layout = "layouts/master.html"
	c.TplName = "log/home.html"
}

func (c *LogController) AdminLog() {
	c.Layout = "layouts/master.html"
	c.TplName = "log/admin.html"
}
