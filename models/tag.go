package models

// IsTagExists 判断标签是否存在
func IsTagExists(filter map[string]interface{}) bool {
	needle := o.QueryTable("tag")
	for key, value := range filter {
		needle = needle.Filter(key, value)
	}

	return needle.Exist()
}

// AddTag 添加标签
func AddTag(data Tag) (int64, error) {
	return o.Insert(&data)
}

// GetTags 获取标签
func GetTags(filter map[string]interface{}) ([]*Tag, error) {
	var tags []*Tag

	needle := o.QueryTable("tag")
	for key, value := range filter {
		needle = needle.Filter(key, value)
	}

	_, err := needle.All(&tags)
	return tags, err
}

// GetOneTag 获取标签
func GetOneTag(filter map[string]interface{}) (Tag, error) {
	var tag Tag

	needle := o.QueryTable("tag")
	for key, value := range filter {
		needle = needle.Filter(key, value)
	}

	err := needle.One(&tag)
	return tag, err
}

// UpdateTag 更新标签
func UpdateTag(filter, values map[string]interface{}) (int64, error) {
	needle := o.QueryTable("tag")
	for key, value := range filter {
		needle = needle.Filter(key, value)
	}

	return needle.Update(values)
}

// DeleteTag 删除标签
func DeleteTag(filter map[string]interface{}) (int64, error) {
	needle := o.QueryTable("tag")
	for key, value := range filter {
		needle = needle.Filter(key, value)
	}

	return needle.Delete()
}
