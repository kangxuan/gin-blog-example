package tag_service

import (
	"encoding/json"
	"gin-blog-example/models"
	"gin-blog-example/pkg/export"
	"gin-blog-example/pkg/gredis"
	"gin-blog-example/pkg/logging"
	"gin-blog-example/services/cache_service"
	"gin-blog-example/settings"
	"github.com/tealeg/xlsx"
	"github.com/xuri/excelize/v2"
	"io"
	"strconv"
	time2 "time"
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

// Export 导出表格
func (t *Tag) Export() (string, error) {
	tags, err := t.GetAll()
	if err != nil {
		return "", err
	}

	// 新建一个文件
	file := xlsx.NewFile()
	// 添加一个新Sheet，并命名为"标签信息"
	sheet, err := file.AddSheet("标签信息")
	if err != nil {
		return "", err
	}

	// 将标题写入一行
	titles := []string{"ID", "名称", "创建人", "创建时间", "修改人", "修改时间"}
	// 加一行
	row := sheet.AddRow()

	var cell *xlsx.Cell
	for _, title := range titles {
		// 加一个单元格
		cell = row.AddCell()
		cell.Value = title
	}

	// 将查询出来的数据写入表格
	for _, v := range tags {
		values := []string{
			strconv.Itoa(v.ID),
			v.Name,
			v.CreatedBy,
			strconv.Itoa(v.CreatedOn),
			v.ModifiedBy,
			strconv.Itoa(v.ModifiedOn),
		}

		// 加一行
		row = sheet.AddRow()
		for _, value := range values {
			// 加一个单元格
			cell = row.AddCell()
			cell.Value = value
		}
	}

	// 表格名称
	time := strconv.Itoa(int(time2.Now().Unix()))
	filename := "tags-" + time + ".xlsx"

	fullPath := export.GetExcelFullPath() + filename
	// 保存表格
	err = file.Save(fullPath)
	if err != nil {
		return "", err
	}

	return filename, nil
}

// Import 导入表格
func (t *Tag) Import(r io.Reader) error {
	// 读取文件
	xlsxReader, err := excelize.OpenReader(r)
	if err != nil {
		return err
	}

	// 读取哪个表格
	rows, err := xlsxReader.GetRows("Sheet1")
	if err != nil {
		return err
	}
	// 将数据写入数据库
	for index, row := range rows {
		if index > 0 {
			var data []string
			for _, cell := range row {
				data = append(data, cell)
			}
			err := models.AddTag(data[1], 1, data[2])
			if err != nil {
				return err
			}
		}
	}
	return nil
}
