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
	c.LayoutSections = map[string]string{
		"Script": "favor_record/index_script.html",
	}

	filter := map[string]interface{}{}
	if ManagerInfo.IsAdmin != 1 {
		filter["article__manager__id"] = ManagerInfo.Id
	}

	AddLog(c.Ctx, "查看点赞记录", "", "PAGE")
	records, _ := models.GetAllFavorRecords(filter)
	c.Data["records"] = records
}
