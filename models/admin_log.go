package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type AdminLog struct {
	Id        int
	ManagerId int       // 管理员id
	Content   string    // 操作
	Ip        string    // 客户端IP
	Url       string    // 请求地址
	Method    string    // 请求方式
	Query     string    // 地址栏参数
	Headers   string    // 请求头
	Body      string    // 请求体
	Response  string    // 响应内容
	Result    string    // 操作结果
	Reason    string    // 失败原因
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

// GetAdminLogs 获取多条日志
func GetAdminLogs(filter map[string]interface{}, offset, limit int) ([]*AdminLog, error) {
	var logs []*AdminLog

	needle := orm.NewOrm().QueryTable("admin_log")
	for key, value := range filter {
		needle = needle.Filter(key, value)
	}

	_, err := needle.Offset(offset).Limit(limit).OrderBy("-id").All(&logs)
	return logs, err
}
