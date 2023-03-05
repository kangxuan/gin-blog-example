package article_service

import (
	"encoding/json"
	"errors"
	"gin-blog-example/models"
	"gin-blog-example/pkg/gredis"
	"gin-blog-example/pkg/logging"
	"gin-blog-example/services/cache_service"
	"gin-blog-example/settings"
)

type Article struct {
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

// Add 添加文章
func (a *Article) Add() error {
	article := models.Article{
		TagId:         a.TagID,
		Title:         a.Title,
		Desc:          a.Desc,
		Content:       a.Content,
		CoverImageUrl: a.CoverImageUrl,
		CreatedBy:     a.CreatedBy,
		State:         a.State,
	}

	if !models.AddArticle(article) {
		return errors.New("添加文章失败")
	}

	return nil
}

// Update 更新文章
func (a *Article) Update() error {
	if !models.UpdateArticle(a.ID, models.Article{
		TagId:         a.TagID,
		Title:         a.Title,
		Desc:          a.Desc,
		Content:       a.Content,
		CoverImageUrl: a.CoverImageUrl,
		ModifiedBy:    a.ModifiedBy,
		State:         a.State,
	}) {
		return errors.New("更新文章失败")
	}

	return nil
}

// Get 获取文章
func (a *Article) Get() (*models.ArticleResult, error) {
	var cacheArticle *models.ArticleResult

	cache := cache_service.Article{ID: a.ID}
	key := cache.GetArticlesKey()

	if gredis.Exists(key) {
		// 缓存存在了
		data, err := gredis.Get(key)
		if err != nil {
			logging.Info(err)
		} else {
			// JsonDecode
			_ = json.Unmarshal(data, &cacheArticle)
			return cacheArticle, nil
		}
	}

	article := models.GetArticleById(a.ID)

	_ = gredis.Set(key, article, settings.RedisSetting.ExpireTime)
	return article, nil
}

// GetAll 获取文章列表
func (a *Article) GetAll() ([]*models.ArticleResult, error) {
	var (
		cacheArticles []*models.ArticleResult
	)

	cache := cache_service.Article{
		TagID: a.TagID,
		State: a.State,

		PageNum:  a.PageNum,
		PageSize: a.PageSize,
	}
	key := cache.GetArticlesKey()

	if gredis.Exists(key) {
		data, err := gredis.Get(key)
		if err != nil {
			logging.Info(err)
		} else {
			_ = json.Unmarshal(data, &cacheArticles)
			return cacheArticles, nil
		}
	}

	articles := models.GetArticles(a.PageNum, a.PageSize, a.GetMaps())
	_ = gredis.Set(key, articles, settings.RedisSetting.ExpireTime)
	return articles, nil
}

// Delete 删除文章
func (a *Article) Delete() error {
	return models.DeleteArticle(a.ID)
}

// Count 文章数量
func (a *Article) Count() (int64, error) {
	return models.GetArticlesTotal(a.GetMaps())
}

// ExistedById 根据ID获取文章是否存在
func (a *Article) ExistedById() (bool, error) {
	return models.ExistedArticleById(a.ID)
}

// GetMaps 组合查询条件maps
func (a *Article) GetMaps() map[string]interface{} {
	maps := make(map[string]interface{})

	if a.State != -1 {
		maps["state"] = a.State
	}
	if a.TagID != -1 {
		maps["tag_id"] = a.TagID
	}
	if a.Title != "" {
		maps["title"] = a.Title
	}

	return maps
}
