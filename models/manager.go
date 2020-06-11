package models

// IsManagerExists 判断管理员是否存在
func IsManagerExists(filter map[string]interface{}) bool {
	return concatFilter("manager", filter).Exist()
}

// GetOneManager 获取某个管理员信息
func GetOneManager(filter map[string]interface{}) (Manager, error) {
	var manager Manager

	err := concatFilter("manager", filter).One(&manager)
	return manager, err
}

// AddManager 添加管理员
func AddManager(manager Manager) (int64, error) {
	return o.Insert(&manager)
}
