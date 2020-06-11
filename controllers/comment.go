package controllers

import (
	"manage/models"
)

type CommentController struct {
	BaseController
}

// Get 评论列表
func (c *CommentController) Get() {
	c.Layout = "layouts/master.html"
	c.TplName = "comment/index.html"
	c.LayoutSections = map[string]string{
		"Style":  "comment/index_style.html",
		"Script": "comment/index_script.html",
	}

	filter := map[string]interface{}{}

	articleId, _ := c.GetInt("ar", -1)
	if articleId > -1 {
		filter["article_id"] = articleId
	}
	accountId, _ := c.GetInt("ac", -1)
	if accountId > -1 {
		filter["account_id"] = accountId
	}

	articles, _ := models.GetAllArticles(map[string]interface{}{"status": 1})
	accounts, _ := models.GetAllAccounts(nil)
	comments, _ := models.GetAllComments(filter)

	c.Data["ar"] = articleId
	c.Data["ac"] = accountId
	c.Data["articles"] = articles
	c.Data["accounts"] = accounts
	c.Data["comments"] = comments
}

func (c *CommentController) Keyword() {
	c.Layout = "layouts/master.html"
	c.TplName = "comment/keyword.html"
}
