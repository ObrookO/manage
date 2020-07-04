package models

// AddEmailLog 添加发送邮件日志
func AddEmailLog(log EmailLog) (int64, error) {
	return o.Insert(&log)
}

// GetAllEmailLogs 获取所有邮件日志
func GetAllEmailLogs(filter map[string]interface{}, field ...string) ([]*EmailLog, error) {
	var logs []*EmailLog

	_, err := concatFilter("email_log", filter).All(&logs, field...)
	return logs, err
}
