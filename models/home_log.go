package models

// GetHomeLogs 获取多条日志
func GetHomeLogs(filter map[string]interface{}, offset, limit int) ([]*HomeLog, error) {
	var logs []*HomeLog

	needle := o.QueryTable("home_log")
	for key, value := range filter {
		needle = needle.Filter(key, value)
	}

	_, err := needle.Offset(offset).Limit(limit).OrderBy("-id").All(&logs)
	return logs, err
}
