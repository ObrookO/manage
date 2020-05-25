package models

// IsManagerExists 判断管理员是否存在
func IsManagerExists(filter map[string]interface{}) bool {
	needle := o.QueryTable("manager")
	for key, value := range filter {
		needle = needle.Filter(key, value)
	}

	return needle.Exist()
}

// GetOneManager 获取某个管理员信息
func GetOneManager(filter map[string]interface{}) (Manager, error) {
	var manager Manager

	needle := o.QueryTable("manager")
	for key, value := range filter {
		needle = needle.Filter(key, value)
	}

	err := needle.One(&manager)
	return manager, err
}

// AddManager 添加管理员
func AddManager(manager Manager) (int64, error) {
	return o.Insert(&manager)
}
