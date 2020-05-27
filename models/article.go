package models

// IsArticleExists 判断文章是否存在
func IsArticleExists(filter map[string]interface{}) bool {
	needle := o.QueryTable("article")
	for key, value := range filter {
		needle = needle.Filter(key, value)
	}

	return needle.Exist()
}

// AddArticle 添加文章
func AddArticle(article *Article) (int64, error) {
	o.Insert(article)
	m2m := o.QueryM2M(article, "Tags")
	return m2m.Add(article.Tags)
}

// GetArticles 获取文章
func GetArticles(filter map[string]interface{}, offset int, limit int) ([]*Article, error) {
	var articles []*Article

	needle := o.QueryTable("article")
	for key, value := range filter {
		needle = needle.Filter(key, value)
	}

	_, err := needle.Offset(offset).Limit(limit).OrderBy("-id").RelatedSel().All(&articles)
	// 查询标签
	for _, a := range articles {
		o.LoadRelated(a, "Tags")
	}

	return articles, err
}

// GetTotal 获取文章总数
func GetTotal(where map[string]interface{}) (int64, error) {
	needle := o.QueryTable("article")

	for key, value := range where {
		needle = needle.Filter(key, value)
	}

	return needle.Count()
}

// 文章归档
func Archive() []ArticleArchive {
	var archive []ArticleArchive

	o.Raw("select date_format(created_at, '%Y年%m月') as date, date_format(created_at, '%Y/%m') as value, count(*) as sum from article group by value").QueryRows(&archive)
	return archive

}

// GetOneArticle 获取一篇文章
func GetOneArticle(filter map[string]interface{}) (Article, error) {
	var article Article

	needle := o.QueryTable("article")
	for key, value := range filter {
		needle = needle.Filter(key, value)
	}

	err := needle.One(&article)
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
func UpdateArticle(filter, values map[string]interface{}) (int64, error) {
	needle := o.QueryTable("article")
	for k, v := range filter {
		needle = needle.Filter(k, v)
	}

	return needle.Update(values)
}
