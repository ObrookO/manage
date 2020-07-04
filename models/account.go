package models

// AddAccount 添加账号
func AddAccount(data *Account) (int64, error) {
	return o.Insert(data)
}

// UpdateAccountWithFilter 更新账号信息
func UpdateAccountWithFilter(filter, values map[string]interface{}) (int64, error) {
	return concatFilter("account", filter).Update(values)
}

// GetOneAccount 获取账号信息
func GetOneAccount(filter map[string]interface{}) (Account, error) {
	var account Account

	err := concatFilter("account", filter).One(&account)
	return account, err
}

// GetAllAccounts 获取账号信息
func GetAllAccounts(filter map[string]interface{}) ([]*Account, error) {
	var accounts []*Account

	_, err := concatFilter("account", filter).All(&accounts)
	return accounts, err
}
