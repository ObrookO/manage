package models

// GetAllHomeLogs 获取所有前台日志
func GetAllHomeLogs(filter map[string]interface{}) ([]*HomeLog, error) {
	var logs []*HomeLog

	_, err := concatFilter("home_log", filter).OrderBy("-id").All(&logs)
	return logs, err
}
