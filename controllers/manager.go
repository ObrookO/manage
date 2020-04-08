package controllers

type ManagerController struct {
	BaseController
}

func (c *ManagerController) Get() {
	c.Layout = "layouts/master.html"
	c.TplName = "manager/index.html"
}
