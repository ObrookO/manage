package models

const (
	NewManager           = iota // 添加管理员
	ResetManagerPassword        // 重置密码
	RegisterAccount             // 前台用户注册
	ResetAccountPassword        // 前台用户重置密码
)

// AddEmailLog 添加发送邮件日志
func AddEmailLog(log EmailLog) (int64, error) {
	return o.Insert(&log)
}

// GetAllEmailLogs 获取所有邮件日志
func GetAllEmailLogs(filter map[string]interface{}, field ...string) ([]*EmailLog, error) {
	var logs []*EmailLog

	_, err := concatFilter("email_log", filter).OrderBy("-id").All(&logs, field...)
	return logs, err
}
