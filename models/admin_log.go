package models

// 记录日志
func AddAdminLog(data AdminLog) (int64, error) {
	return o.Insert(&data)
}

// GetAdminLogs 获取多条日志
func GetAdminLogs(filter map[string]interface{}, offset, limit int) ([]*AdminLog, error) {
	var logs []*AdminLog

	_, err := concatFilter("admin_log", filter).Offset(offset).Limit(limit).OrderBy("-id").All(&logs)
	return logs, err
}
