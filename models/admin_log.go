package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type AdminLog struct {
	Id        int
	ManagerId int
	Ip        string
	Url       string
	Method    string
	Query     string
	Headers   string
	Body      string
	Response  string
	Content   string
	Reason    string
	CreatedAt time.Time `orm:"auto_now_add;type(timestamp)"`
	UpdatedAt time.Time `orm:"auto_now;type(timestamp)"`
}

func init() {
	orm.RegisterModel(new(AdminLog))
}

// 记录日志
func AddAdminLog(data AdminLog) (int64, error) {
	return orm.NewOrm().Insert(&data)
}
