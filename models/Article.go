package models

import (
	"gorm.io/gorm"
)

type Article struct {
	Model

	TagId int `json:"tag_id"`

	Title         string `json:"title"`
	Desc          string `json:"desc"`
	Content       string `json:"content"`
	CoverImageUrl string `json:"cover_image_url"`
	CreatedBy     string `json:"created_by"`
	ModifiedBy    string `json:"modified_by"`
	DeletedOn     int    `json:"deleted_on"`
	State         int    `json:"state"`
}

type ArticleResult struct {
	ID            int
	TagID         int
	Title         string
	Desc          string
	Content       string
	CoverImageUrl string
	State         int
	CreatedBy     string
	ModifiedBy    string
	TagName       string

	PageNum  int
	PageSize int
}

// GetArticlesTotal 获取文章总数
func GetArticlesTotal(maps map[string]interface{}) (int64, error) {
	var total int64
	if err := db.Model(&Article{}).Where(maps).Count(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

// GetArticles 获取文章列表
func GetArticles(pageNum int, pageSize int, maps map[string]interface{}) (articleResult []*ArticleResult) {
	db.Debug().Model(&Article{}).
		Select("blog_article.*, blog_tag.name as tag_name").
		Joins("left join blog_tag on blog_article.tag_id=blog_tag.id").
		Where(maps).Offset(pageNum).Limit(pageSize).Find(&articleResult)
	return
}

// GetArticleById 根据ID获取单个文章
func GetArticleById(id int) *ArticleResult {
	result := &ArticleResult{}
	db.Debug().Model(&Article{}).
		Select("blog_article.*, blog_tag.name as tag_name").
		Joins("left join blog_tag on blog_article.tag_id=blog_tag.id").
		Where("blog_article.id=?", id).First(&result)
	return result
}

// ExistedArticleById 判断文章ID是否存在
func ExistedArticleById(id int) (bool, error) {
	var article Article
	err := db.Select("id").Where("id=?", id).First(&article).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	if article.ID > 0 {
		return true, nil
	}

	return false, nil
}

// AddArticle 添加文章
func AddArticle(data Article) bool {
	if db.Model(&Article{}).Create(&data).Error != nil {
		return false
	}
	return true
}

// UpdateArticle 更新文章
func UpdateArticle(id int, data Article) bool {
	if db.Model(&Article{}).Where("id=?", id).Updates(data).Error != nil {
		return false
	}
	return true
}

// DeleteArticle 删除文章
func DeleteArticle(id int) error {
	if err := db.Where("id=?", id).Delete(&Article{}).Error; err != nil {
		return err
	}
	return nil
}
