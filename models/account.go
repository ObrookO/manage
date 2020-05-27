package models

// IsAccountExists 判断账号是否存在
func IsAccountExists(filter map[string]interface{}) bool {
	needle := o.QueryTable("account")
	for key, value := range filter {
		needle = needle.Filter(key, value)
	}

	return needle.Exist()
}

// AddAccount 添加账号
func AddAccount(data *Account) (int64, error) {
	return o.Insert(data)
}

// UpdateAccount 更新账号信息
func UpdateAccount(filter, values map[string]interface{}) (int64, error) {
	needle := o.QueryTable("account")
	for key, value := range filter {
		needle = needle.Filter(key, value)
	}

	return needle.Update(values)
}

// GetOneAccount 获取账号信息
func GetOneAccount(filter map[string]interface{}) (Account, error) {
	var account Account

	needle := o.QueryTable("account")
	for key, value := range filter {
		needle = needle.Filter(key, value)
	}

	err := needle.One(&account)
	return account, err
}

// GetAccounts 获取账号信息
func GetAccounts(filter map[string]interface{}) ([]*Account, error) {
	var accounts []*Account

	needle := o.QueryTable("account")
	for key, value := range filter {
		needle = needle.Filter(key, value)
	}

	_, err := needle.All(&accounts)
	return accounts, err
}
