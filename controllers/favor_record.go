package controllers

import (
	"manage/models"
)

type FavorRecordController struct {
	BaseController
}

// Get 点赞记录
func (c *FavorRecordController) Get() {
	c.Layout = "layouts/master.html"
	c.TplName = "favor_record/index.html"

	AddLog(c.Ctx, "查看点赞记录", "", "PAGE", "SUCCESS")
	records, _ := models.GetAllFavorRecords(nil)
	c.Data["records"] = records
}
