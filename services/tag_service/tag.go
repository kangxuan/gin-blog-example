package tag_service

import (
	"encoding/json"
	"gin-blog-example/models"
	"gin-blog-example/pkg/gredis"
	"gin-blog-example/pkg/logging"
	"gin-blog-example/services/cache_service"
	"gin-blog-example/settings"
)

type Tag struct {
	ID         int
	Name       string
	CreateBy   string
	ModifiedBy string
	State      int

	PageNum  int
	PageSize int
}

// GetAll 获取标签列表
func (t *Tag) GetAll() ([]models.Tag, error) {
	var (
		tags, cacheTag []models.Tag
	)
	cache := cache_service.Tag{
		State: t.State,
		Name:  t.Name,

		PageNum:  t.PageNum,
		PageSize: t.PageSize,
	}

	key := cache.GetTagsKey()

	if gredis.Exists(key) {
		data, err := gredis.Get(key)
		if err != nil {
			logging.Info(err)
		} else {
			_ = json.Unmarshal(data, &cacheTag)
			return cacheTag, nil
		}
	}

	tags, err := models.GetTags(t.PageNum, t.PageSize, t.getMaps())
	if err != nil {
		return nil, err
	}

	_ = gredis.Set(key, tags, settings.RedisSetting.ExpireTime)
	return tags, nil
}

// Count 获取标签总数
func (t *Tag) Count() (int64, error) {
	return models.GetTagTotal(t.getMaps())
}

// Add 添加标签
func (t *Tag) Add() error {
	return models.AddTag(t.Name, t.State, t.CreateBy)
}

// Update 更新标签
func (t *Tag) Update() error {
	data := make(map[string]interface{})
	data["modified_by"] = t.ModifiedBy
	data["name"] = t.Name
	if t.State >= 0 {
		data["state"] = t.State
	}
	return models.UpdateTag(t.ID, data)
}

// Delete 删除标签
func (t *Tag) Delete() error {
	return models.DeleteTag(t.ID)
}

// ExistedTagByName 根据Name判断Tag是否存在
func (t *Tag) ExistedTagByName() (bool, error) {
	return models.ExistedTagByName(t.Name)
}

// ExistedTagById 根据Id判断Tag是否存在
func (t *Tag) ExistedTagById() (bool, error) {
	return models.ExistedTagById(t.ID)
}

func (t *Tag) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})

	if t.Name != "" {
		maps["name"] = t.Name
	}
	if t.State >= 0 {
		maps["state"] = t.State
	}

	return maps
}
