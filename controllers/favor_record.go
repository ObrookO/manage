package controllers

type FavorRecordController struct {
	BaseController
}

func (c *FavorRecordController) Get() {
	c.Layout = "layouts/master.html"
	c.TplName = "favor_record/index.html"
}
