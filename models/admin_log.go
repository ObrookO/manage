package models

// 记录日志
func AddAdminLog(data AdminLog) (int64, error) {
	return o.Insert(&data)
}

// GetAllAdminLogs 获取多条日志
func GetAllAdminLogs(filter map[string]interface{}) ([]*AdminLog, error) {
	var logs []*AdminLog

	_, err := concatFilter("admin_log", filter).OrderBy("-id").RelatedSel("manager").All(&logs)
	return logs, err
}
