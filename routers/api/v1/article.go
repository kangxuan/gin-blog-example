package v1

import (
	"gin-blog-example/models"
	"gin-blog-example/pkg/e"
	"gin-blog-example/pkg/util"
	"gin-blog-example/settings"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"net/http"
)

// GetArticles 获取文章列表
func GetArticles(c *gin.Context) {
	data := make(map[string]interface{})
	maps := make(map[string]interface{})
	var valid validation.Validation

	title := c.Query("title")
	if title != "" {
		maps["title"] = title
	}

	if arg := c.Query("state"); arg != "" {
		state := com.StrTo(arg).MustInt()
		valid.Range(state, 0, 1, "state").Message("状态只允许0或1")

		maps["state"] = state
	}

	if arg := c.Query("tag_id"); arg != "" {
		tagId := com.StrTo(arg).MustInt()
		valid.Min(tagId, 1, "tag_id").Message("标签ID必须大于0")

		maps["tag_id"] = tagId
	}

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		code = e.SUCCESS
		data["list"] = models.GetArticles(util.GetPage(c), settings.PageSize, maps)
		data["total"] = models.GetArticlesTotal(maps)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}

// GetArticle 获取单个文章
func GetArticle(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()
	code := e.INVALID_PARAMS

	var valid validation.Validation
	valid.Min(id, 1, "id").Message("文章ID必须大于0")

	var article interface{}
	if !valid.HasErrors() {
		code = e.SUCCESS
		article = models.GetArticleById(id)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": article,
	})
}

// AddArticle 添加文章
func AddArticle(c *gin.Context) {
	var article models.Article
	code := e.INVALID_PARAMS
	_ = c.BindJSON(&article)

	valid := validation.Validation{}
	valid.Min(article.TagId, 1, "tag_id").Message("标签ID必须大于0")
	valid.Required(article.Title, "title").Message("标题不能为空")
	valid.MaxSize(article.Title, 100, "title").Message("标题不能超过100个字符")
	valid.Required(article.Desc, "desc").Message("简述不能为空")
	valid.MaxSize(article.Desc, 255, "desc").Message("简述不能超过255个字符")
	valid.Required(article.CreatedBy, "created_by").Message("创建人不能为空")
	valid.MaxSize(article.CreatedBy, 100, "created_by").Message("创建人不能超过100个字符")
	valid.Range(article.State, 0, 1, "state").Message("状态只能是0和1")

	if !valid.HasErrors() {
		if !models.ExistedTagById(article.TagId) {
			code = e.ERROR_NOT_EXIST_TAG
		} else {
			if !models.AddArticle(article) {
				code = e.ERROR_ADD_ARTCLIE_FAIL
			}
			code = e.SUCCESS
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]interface{}),
	})
}

// UpdateArticle 更新文章
func UpdateArticle(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()
	code := e.INVALID_PARAMS
	var article models.Article
	var valid validation.Validation
	_ = c.BindJSON(&article)

	valid.Min(id, 1, "id").Message("文章ID必须大于0")
	valid.Min(article.TagId, 1, "tag_id").Message("标签ID必须大于0")
	valid.Required(article.Title, "title").Message("标题不能为空")
	valid.MaxSize(article.Title, 100, "title").Message("标题不能超过100个字符")
	valid.Required(article.Desc, "desc").Message("简述不能为空")
	valid.MaxSize(article.Desc, 255, "desc").Message("简述不能超过255个字符")
	valid.Required(article.ModifiedBy, "modified_by").Message("更新人不能为空")
	valid.MaxSize(article.ModifiedBy, 100, "modified_by").Message("更新人不能超过100个字符")
	valid.Range(article.State, 0, 1, "state").Message("状态只能是0和1")

	if !valid.HasErrors() {
		if !models.ExistedArticleById(id) {
			code = e.ERROR_NOT_EXIST_ARTICLE
		} else {
			if !models.UpdateArticle(id, article) {
				code = e.ERROR_ADD_ARTCLIE_FAIL
			} else {
				code = e.SUCCESS
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]interface{}),
	})
}

// DeleteArticle 删除文章
func DeleteArticle(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()
	code := e.INVALID_PARAMS
	var valid validation.Validation

	valid.Min(id, 1, "id").Message("文章ID必须大于0")
	if !valid.HasErrors() {
		if !models.ExistedArticleById(id) {
			code = e.ERROR_NOT_EXIST_ARTICLE
		} else {
			if !models.DeleteArticle(id) {
				code = e.ERROR_DELETE_ARTCLIE_FAIL
			} else {
				code = e.SUCCESS
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]interface{}),
	})
}
