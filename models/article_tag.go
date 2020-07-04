package models

// GetArticleNumOfTag 获取文章数量
func GetArticleNumOfTag(filter map[string]interface{}) int64 {
	num, _ := concatFilter("article_tag", filter).Count()
	return num
}

// DeleteArticleTag 删除文章标签
func DeleteArticleTag(filter map[string]interface{}) (int64, error) {
	return concatFilter("article_tag", filter).Delete()
}
