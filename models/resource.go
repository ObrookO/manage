package models

// GetAllResource 获取所有干货收藏
func GetAllResource(filter map[string]interface{}, field ...string) ([]*Resource, error) {
	var list []*Resource

	_, err := concatFilter("resource", filter).All(&list, field...)
	return list, err
}

// GetOneResource 获取干货收藏
func GetOneResource(filter map[string]interface{}, field ...string) (Resource, error) {
	var resource Resource

	err := concatFilter("resource", filter).One(&resource, field...)
	return resource, err
}

// AddResource 添加干货收藏
func AddResource(resource Resource) (int64, error) {
	return o.Insert(&resource)
}

// DeleteResource 删除干货收藏
func DeleteResource(filter map[string]interface{}) (int64, error) {
	return concatFilter("resource", filter).Delete()
}
