package models

type Tag struct {
	Model
	Name       string `json:"name"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	DeletedOn  int    `json:"deleted_on"`
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

// ExistedTagByName 判断标签名称
func ExistedTagByName(name string) bool {
	var tag Tag
	db.Select("id").Where("name=?", name).First(&tag)
	if tag.ID > 0 {
		return true
	}
	return false
}

// AddTag 添加标签
func AddTag(name string, state int, createBy string) bool {
	err := db.Create(&Tag{
		Name:      name,
		State:     state,
		CreatedBy: createBy,
	}).Error
	if err != nil {
		return false
	}
	return true
}
