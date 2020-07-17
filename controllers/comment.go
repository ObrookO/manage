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
		"Script": "comment/index_script.html",
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
	// 根据账号查询
	accountId, _ := c.GetInt("ac", -1)
	if accountId > -1 {
		filter["account_id"] = accountId
	}

	addLog(c.Ctx, "查看评论列表", "", "PAGE")

	articles, _ := models.GetAllArticles(articleFilter)
	accounts, _ := models.GetAllAccounts(nil)
	comments, _ := models.GetAllComments(filter)

	c.Data["ar"] = articleId
	c.Data["ac"] = accountId
	c.Data["articles"] = articles
	c.Data["accounts"] = accounts
	c.Data["comments"] = comments
}

// Delete 删除评论
func (c *CommentController) Delete() {
	id, _ := c.GetInt("id")

	comment, _ := models.GetOneComment(map[string]interface{}{"id": id})
	if comment.Id == 0 {
		addLog(c.Ctx, "删除评论", "评论不存在", "{\"code\": 400000, \"msg\":\"评论不存在\"}")
		c.Data["json"] = &JSONResponse{Code: 400000, Msg: "评论不存在"}
		c.ServeJSON()
		return
	}

	logContent := "删除评论，文章作者：" + comment.Article.Manager.Username + "，文章标题：" + comment.Article.Title + "，评论账号：" + comment.Account.
		Username + "，评论内容：" + comment.
		OriginalContent

	// 判断权限
	if ManagerInfo.IsAdmin != 1 {
		if ManagerInfo.Id != comment.Article.Manager.Id {
			addLog(c.Ctx, logContent, "非法操作", "{\"code\": 500, \"msg\":\"非法操作\"}")
			c.Data["json"] = &JSONResponse{Code: 500, Msg: "非法操作"}
			c.ServeJSON()
			return
		}
	}

	if _, err := models.DeleteComment(map[string]interface{}{"id": id}); err != nil {
		addLog(c.Ctx, logContent, err.Error(), "{\"code\": 400001, \"msg\":\"操作失败\"}")
		c.Data["json"] = &JSONResponse{Code: 400001, Msg: "操作失败"}
		c.ServeJSON()
		return
	}

	addLog(c.Ctx, logContent, "", "{\"code\": 200, \"msg\":\"OK\"}")
	c.Data["json"] = &JSONResponse{Code: 200, Msg: "OK"}
	c.ServeJSON()
}
