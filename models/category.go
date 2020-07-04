package models

// IsCategoryExists 判断栏目是否存在
func IsCategoryExists(filter map[string]interface{}) bool {
	return concatFilter("category", filter).Exist()
}

// AddCategory 添加栏目
func AddCategory(data Category) (int64, error) {
	return o.Insert(&data)
}

// GetAllCategories 获取所有的栏目
func GetAllCategories(filter map[string]interface{}) ([]*Category, error) {
	var categories []*Category

	_, err := concatFilter("category", filter).RelatedSel().All(&categories)
	return categories, err
}

// GetCategory 获取某个栏目
func GetCategory(filter map[string]interface{}) (Category, error) {
	var category Category

	err := concatFilter("category", filter).One(&category)
	return category, err
}

// DeleteCategory 删除栏目
func DeleteCategory(filter map[string]interface{}) (int64, error) {
	return concatFilter("category", filter).Delete()
}

// UpdateCategoryWithFilter 更新栏目
func UpdateCategoryWithFilter(filter, values map[string]interface{}) (int64, error) {
	return concatFilter("category", filter).Update(values)
}
