package models

// GetHomeLogs 获取多条日志
func GetHomeLogs(filter map[string]interface{}, offset, limit int) ([]*HomeLog, error) {
	var logs []*HomeLog

	_, err := concatFilter("home_log", filter).Offset(offset).Limit(limit).OrderBy("-id").All(&logs)
	return logs, err
}
