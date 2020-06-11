package models

// GetAllFavorRecords 获取所有的点赞记录
func GetAllFavorRecords(filter map[string]interface{}) ([]*FavorRecord, error) {
	var records []*FavorRecord

	_, err := concatFilter("favor_record", filter).RelatedSel().All(&records)
	return records, err
}
