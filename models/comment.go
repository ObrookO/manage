package models

// 新增评论
func AddComment(value *Comment) (int64, error) {
	return o.Insert(value)
}

// 获取所有评论
func GetComments(where map[string]interface{}) ([]Comment, int64, error) {
	var comments []Comment
	needle := o.QueryTable("comment")

	for key, value := range where {
		needle = needle.Filter(key, value)
	}

	rows, err := needle.OrderBy("-id").Limit(20).All(&comments)

	return comments, rows, err
}

// 获取的评论数量
func GetCommentsCount(where map[string]interface{}) (int64, error) {
	needle := o.QueryTable("comment")

	for key, value := range where {
		needle = needle.Filter(key, value)
	}

	return needle.Count()
}
