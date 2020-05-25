package models

// 记录日志
func AddAdminLog(data AdminLog) (int64, error) {
	return o.Insert(&data)
}

// GetAdminLogs 获取多条日志
func GetAdminLogs(filter map[string]interface{}, offset, limit int) ([]*AdminLog, error) {
	var logs []*AdminLog

	needle := o.QueryTable("admin_log")
	for key, value := range filter {
		needle = needle.Filter(key, value)
	}

	_, err := needle.Offset(offset).Limit(limit).OrderBy("-id").All(&logs)
	return logs, err
}
