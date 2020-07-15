package models

// IsArticleExists 判断文章是否存在
func IsArticleExists(filter map[string]interface{}) bool {
	return concatFilter("article", filter).Exist()
}

// AddArticle 添加文章
func AddArticle(article Article) (int64, error) {
	o.Insert(&article)
	m2m := o.QueryM2M(&article, "Tags")
	return m2m.Add(article.Tags)
}

// GetAllArticles 获取文章
func GetAllArticles(filter map[string]interface{}) ([]*Article, error) {
	var articles []*Article

	_, err := concatFilter("article", filter).
		RelatedSel().
		OrderBy("-id").
		All(&articles)

	// 查询标签
	for _, a := range articles {
		o.LoadRelated(a, "Tags")
		o.LoadRelated(a, "Comments")
		o.LoadRelated(a, "Favors")
	}

	return articles, err
}

// GetArticleNumOfCategory 获取文章总数
func GetArticleNumOfCategory(filter map[string]interface{}) int64 {
	num, _ := concatFilter("article", filter).Count()
	return num
}

// GetOneArticle 获取一篇文章
func GetOneArticle(filter map[string]interface{}) (Article, error) {
	var article Article

	err := concatFilter("article", filter).RelatedSel().One(&article)
	o.LoadRelated(&article, "Tags")
	return article, err
}

// GetBeforeAndAfter 获取上一篇文章的id以及下一篇文章的id
func GetBeforeAndAfter(aid int) (int, int) {
	var article Article
	var before, after int

	needle := o.QueryTable("article")
	// 获取上一篇文章的id
	if err := needle.Filter("id__lt", aid).OrderBy("-id").One(&article, "id"); err == nil {
		before = article.Id
	} else {
		before = 0
	}

	// 获取下一篇文章的id
	if err := needle.Filter("id__gt", aid).OrderBy("id").One(&article, "id"); err == nil {
		after = article.Id
	} else {
		after = 0
	}

	return before, after
}

// UpdateArticle 更新文章
func UpdateArticle(article Article, field ...string) (int64, error) {
	o.Update(&article, field...)
	m2m := o.QueryM2M(&article, "Tags")

	DeleteArticleTag(map[string]interface{}{"article": article.Id})
	return m2m.Add(article.Tags)
}

// UpdateArticleWithFilter 更新文章
func UpdateArticleWithFilter(filter, values map[string]interface{}) (int64, error) {
	return concatFilter("article", filter).Update(values)
}

// DeleteArticle 删除文章
func DeleteArticle(filter map[string]interface{}) (int64, error) {
	return concatFilter("article", filter).Delete()
}
