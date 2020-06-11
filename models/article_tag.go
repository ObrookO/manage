package models

// DeleteArticleTag 删除文章标签
func DeleteArticleTag(filter map[string]interface{}) (int64, error) {
	return concatFilter("article_tag", filter).Delete()
}
