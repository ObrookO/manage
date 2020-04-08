package controllers

type RoleController struct {
	BaseController
}

func (c *RoleController) Get() {
	c.Layout = "layouts/master.html"
	c.TplName = "role/index.html"
}
