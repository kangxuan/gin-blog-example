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
func GetTags(pageNum int, pageSize int, maps interface{}) (tags []Tag, err error) {
	err = db.Debug().Where(maps).Offset(pageNum).Limit(pageSize).Find(&tags).Error
	if err != nil {
		return nil, err
	}
	return
}

// GetTagTotal 获取标签数量
func GetTagTotal(maps interface{}) (total int64, err error) {
	err = db.Model(&Tag{}).Where(maps).Count(&total).Error
	if err != nil {
		return 0, err
	}
	return
}

// ExistedTagByName 判断标签名称
func ExistedTagByName(name string) (bool, error) {
	var tag Tag
	err := db.Select("id").Where("name=?", name).First(&tag).Error
	if err != nil {
		return false, err
	}
	if tag.ID > 0 {
		return true, nil
	}
	return false, nil
}

// ExistedTagById 判断标签ID
func ExistedTagById(id int) (bool, error) {
	var tag Tag
	err := db.Select("id").Where("id=?", id).Find(&tag).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

// AddTag 添加标签
func AddTag(name string, state int, createBy string) error {
	err := db.Create(&Tag{
		Name:      name,
		State:     state,
		CreatedBy: createBy,
	}).Error
	if err != nil {
		return err
	}
	return nil
}

// UpdateTag 更新标签
func UpdateTag(id int, data interface{}) error {
	err := db.Model(&Tag{}).Where("id=?", id).Updates(data).Error
	if err != nil {
		return err
	}
	return nil
}

// DeleteTag 删除标签
func DeleteTag(id int) error {
	err := db.Model(&Tag{}).Where("id=?", id).Delete(&Tag{}).Error
	if err != nil {
		return err
	}
	return nil
}
