package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type Manager struct {
	Id        int
	Username  string
	Nickname  string
	Email     string
	Password  string
	Avatar    string
	CreatedAt time.Time `orm:"auto_now_add;type(timestamp)"`
	UpdatedAt time.Time `orm:"auto_now;type(timestamp)"`
}

func init() {
	orm.RegisterModelWithPrefix("admin_", new(Manager))
}

// IsManagerExists 判断管理员是否存在
func IsManagerExists(filter map[string]interface{}) bool {
	needle := orm.NewOrm().QueryTable("admin_manager")
	for key, value := range filter {
		needle = needle.Filter(key, value)
	}

	return needle.Exist()
}

// GetOneManager 获取某个管理员信息
func GetOneManager(filter map[string]interface{}) (Manager, error) {
	var manager Manager

	needle := orm.NewOrm().QueryTable("admin_manager")
	for key, value := range filter {
		needle = needle.Filter(key, value)
	}

	err := needle.One(&manager)
	return manager, err
}

// AddManager 添加管理员
func AddManager(manager Manager) (int64, error) {
	return orm.NewOrm().Insert(&manager)
}
