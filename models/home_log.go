package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type HomeLog struct {
	Id        int
	UserId    int
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
	orm.RegisterModel(new(HomeLog))
}

// GetHomeLogs 获取多条日志
func GetHomeLogs(filter map[string]interface{}, offset, limit int) ([]*HomeLog, error) {
	var logs []*HomeLog

	needle := orm.NewOrm().QueryTable("home_log")
	for key, value := range filter {
		needle = needle.Filter(key, value)
	}

	_, err := needle.Offset(offset).Limit(limit).OrderBy("-id").All(&logs)
	return logs, err
}
