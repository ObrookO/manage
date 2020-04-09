package models

import (
	_ "fmt"
	"time"

	"github.com/astaxie/beego/orm"
)

type Account struct {
	Id           int
	Username     string    // 用户名
	Email        string    // 邮箱
	Password     string    // 密码
	Avatar       string    // 头像
	AllowComment uint8     // 是否允许评论
	Status       uint8     // 状态
	CreatedAt    time.Time `orm:"auto_now_add;type(timestamp)"`
	UpdatedAt    time.Time `orm:"auto_now;type(timestamp)"`
}

func init() {
	orm.RegisterModelWithPrefix("home_", new(Account))
}

// IsAccountExists 判断账号是否存在
func IsAccountExists(filter map[string]interface{}) bool {
	needle := orm.NewOrm().QueryTable("home_account")
	for key, value := range filter {
		needle = needle.Filter(key, value)
	}

	return needle.Exist()
}

// AddAccount 添加账号
func AddAccount(data *Account) (int64, error) {
	return orm.NewOrm().Insert(data)
}

// UpdateAccount 更新账号信息
func UpdateAccount(filter, values map[string]interface{}) (int64, error) {
	needle := orm.NewOrm().QueryTable("home_account")
	for key, value := range filter {
		needle = needle.Filter(key, value)
	}

	return needle.Update(values)
}

// GetOneAccount 获取账号信息
func GetOneAccount(filter map[string]interface{}) (Account, error) {
	var account Account

	needle := orm.NewOrm().QueryTable("home_account")
	for key, value := range filter {
		needle = needle.Filter(key, value)
	}

	err := needle.One(&account)
	return account, err
}

// GetAccounts 获取账号信息
func GetAccounts(filter map[string]interface{}) ([]*Account, error) {
	var accounts []*Account

	needle := orm.NewOrm().QueryTable("home_account")
	for key, value := range filter {
		needle = needle.Filter(key, value)
	}

	_, err := needle.All(&accounts)
	return accounts, err
}
