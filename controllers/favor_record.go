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
	articleFilter := map[string]interface{}{"status": 1}
	if ManagerInfo.IsAdmin != 1 {
		filter["article__manager__id"] = ManagerInfo.Id
		articleFilter["manager_id"] = ManagerInfo.Id
	}

	// 根据文章id查询
	articleId, _ := c.GetInt("ar", -1)
	if articleId > -1 {
		filter["article_id"] = articleId
	}

	addLog(c.Ctx, "查看点赞记录", "", "PAGE")

	records, _ := models.GetAllFavorRecords(filter)
	articles, _ := models.GetAllArticles(articleFilter)

	c.Data["ar"] = articleId
	c.Data["articles"] = articles
	c.Data["records"] = records
}
