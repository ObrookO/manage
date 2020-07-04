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

	AddLog(c.Ctx, "查看点赞记录", "", "PAGE")
	records, _ := models.GetAllFavorRecords(map[string]interface{}{"article__manager__id": ManagerInfo.Id})
	c.Data["records"] = records
}
