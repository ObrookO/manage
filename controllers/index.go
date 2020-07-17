package controllers

import (
	"manage/models"
	"time"
)

type IndexController struct {
	BaseController
}

func (c *IndexController) Get() {
	c.Layout = "layouts/master.html"
	c.TplName = "index/index.html"

	tFilter := map[string]interface{}{"created_at__startswith": time.Now().Format("2006-01-02"), "status": 1}
	mFilter := map[string]interface{}{"created_at__startswith": time.Now().Format("2006-01"), "status": 1}
	aFilter := map[string]interface{}{"status": 1}

	if ManagerInfo.IsAdmin != 1 {
		tFilter["manager_id"] = ManagerInfo.Id
		mFilter["manager_id"] = ManagerInfo.Id
		aFilter["manager_id"] = ManagerInfo.Id
	}

	todayArticleAmount, _ := models.GetArticleAmount(tFilter) // 今日发布文章数
	monthArticleAmount, _ := models.GetArticleAmount(mFilter) // 本月发布文章数
	articleTotalAmount, _ := models.GetArticleAmount(aFilter) // 文章总数

	t2Filter := map[string]interface{}{"created_at__startswith": time.Now().Format("2006-01-02")}
	m2Filter := map[string]interface{}{"created_at__startswith": time.Now().Format("2006-01")}

	todayAccountAmount, _ := models.GetAccountAmount(t2Filter)                              // 今日新增账号数
	monthAccountAmount, _ := models.GetAccountAmount(m2Filter)                              // 本月新增账号数
	accountAmount, _ := models.GetAccountAmount(map[string]interface{}{"status": 1})        // 可用账号数
	disableAccountAmount, _ := models.GetAccountAmount(map[string]interface{}{"status": 0}) // 不可用账号数

	c.Data["tAmount"] = todayArticleAmount
	c.Data["mAmount"] = monthArticleAmount
	c.Data["aAmount"] = articleTotalAmount
	c.Data["t2Amount"] = todayAccountAmount
	c.Data["m2Amount"] = monthAccountAmount
	c.Data["a2Amount"] = accountAmount
	c.Data["da2Amount"] = disableAccountAmount
}
