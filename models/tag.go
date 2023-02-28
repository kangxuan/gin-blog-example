package models

type Tag struct {
	Model
	Name       string `json:"name"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	DeleteOn   int    `json:"delete_on"`
	State      int    `json:"state"`
}

// GetTags 获取标签列表
func GetTags(pageNum int, pageSize int, maps interface{}) (tags []Tag) {
	db.Debug().Where(maps).Offset(pageNum).Limit(pageSize).Find(&tags)
	return
}

// GetTagTotal 获取标签数量
func GetTagTotal(maps interface{}) (total int64) {
	db.Model(&Tag{}).Where(maps).Count(&total)
	return
}
