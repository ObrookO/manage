package models

// 获取所有的点赞记录
func GetAllFavorRecord() ([]FavorRecord, error) {
	var records []FavorRecord

	_, err := o.QueryTable("favor_record").All(&records)
	return records, err
}
