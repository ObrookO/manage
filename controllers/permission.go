package controllers

type PermissionController struct {
	BaseController
}

func (c *PermissionController) Get() {
	c.Layout = "layouts/master.html"
	c.TplName = "permission/index.html"
}
