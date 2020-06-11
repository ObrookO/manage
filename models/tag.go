package models

// IsTagExists 判断标签是否存在
func IsTagExists(filter map[string]interface{}) bool {
	return concatFilter("tag", filter).Exist()
}

// AddTag 添加标签
func AddTag(data Tag) (int64, error) {
	return o.Insert(&data)
}

// GetTags 获取标签
func GetTags(filter map[string]interface{}) ([]*Tag, error) {
	var tags []*Tag

	_, err := concatFilter("tag", filter).All(&tags)
	return tags, err
}

// GetOneTag 获取标签
func GetOneTag(filter map[string]interface{}) (Tag, error) {
	var tag Tag

	err := concatFilter("tag", filter).One(&tag)
	return tag, err
}

// UpdateTagWithFilter 更新标签
func UpdateTagWithFilter(filter, values map[string]interface{}) (int64, error) {
	return concatFilter("tag", filter).Update(values)
}

// DeleteTag 删除标签
func DeleteTag(filter map[string]interface{}) (int64, error) {
	return concatFilter("tag", filter).Delete()
}
