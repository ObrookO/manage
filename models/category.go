package models

// IsCategoryExists 判断栏目是否存在
func IsCategoryExists(filter map[string]interface{}) bool {
	needle := o.QueryTable("category")
	for key, value := range filter {
		needle = needle.Filter(key, value)
	}

	return needle.Exist()
}

// AddCategory 添加栏目
func AddCategory(data Category) (int64, error) {
	return o.Insert(&data)
}

// GetCategories 获取所有的栏目
func GetCategories(filter map[string]interface{}) ([]*Category, error) {
	var categories []*Category

	needle := o.QueryTable("category")
	for key, value := range filter {
		needle = needle.Filter(key, value)
	}

	_, err := needle.All(&categories)
	return categories, err
}

// GetCategory 获取某个栏目
func GetCategory(filter map[string]interface{}) (Category, error) {
	var category Category

	needle := o.QueryTable("category")
	for key, value := range filter {
		needle = needle.Filter(key, value)
	}

	err := needle.One(&category)
	return category, err
}

// DeleteCategory 删除栏目
func DeleteCategory(filter map[string]interface{}) (int64, error) {
	needle := o.QueryTable("category")
	for key, value := range filter {
		needle = needle.Filter(key, value)
	}

	return needle.Delete()
}

// UpdateCategory 更新栏目
func UpdateCategory(filter, values map[string]interface{}) (int64, error) {
	needle := o.QueryTable("category")
	for key, value := range filter {
		needle = needle.Filter(key, value)
	}

	return needle.Update(values)
}
