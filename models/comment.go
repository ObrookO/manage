package models

// 新增评论
func AddComment(value *Comment) (int64, error) {
	return o.Insert(value)
}

// GetAllComments 获取所有评论
func GetAllComments(filter map[string]interface{}) ([]*Comment, error) {
	var comments []*Comment

	_, err := concatFilter("comment", filter).OrderBy("-id").RelatedSel().All(&comments)
	return comments, err
}

// 获取的评论数量
func GetCommentsCount(filter map[string]interface{}) (int64, error) {
	return concatFilter("comment", filter).Count()
}
