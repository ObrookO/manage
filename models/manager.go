package models

// IsManagerExists 判断管理员是否存在
func IsManagerExists(filter map[string]interface{}) bool {
	return concatFilter("manager", filter).Exist()
}

// GetAllManagers 获取所有管理员
func GetAllManagers(filter map[string]interface{}, field ...string) ([]*Manager, error) {
	var managers []*Manager

	_, err := concatFilter("manager", filter).All(&managers, field...)
	return managers, err
}

// GetOneManager 获取管理员
func GetOneManager(filter map[string]interface{}, field ...string) (Manager, error) {
	var manager Manager

	err := concatFilter("manager", filter).One(&manager, field...)
	return manager, err
}

// AddManager 添加管理员
func AddManager(manager Manager) (int64, error) {
	return o.Insert(&manager)
}

// UpdateManager 更新管理员
func UpdateManager(filter, value map[string]interface{}) (int64, error) {
	return concatFilter("manager", filter).Update(value)
}

// DeleteManager 删除用户
func DeleteManager(filter map[string]interface{}) (int64, error) {
	return concatFilter("manager", filter).Delete()
}
