package v1

import (
	"gin-blog-example/models"
	"gin-blog-example/pkg/app"
	"gin-blog-example/pkg/e"
	"gin-blog-example/pkg/qrcode"
	"gin-blog-example/pkg/util"
	"gin-blog-example/services/article_service"
	"gin-blog-example/services/tag_service"
	"gin-blog-example/settings"
	"github.com/astaxie/beego/validation"
	"github.com/boombuler/barcode/qr"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"net/http"
)

var (
	QRCODE_URL = "https://github.com/EDDYCJY/blog#gin%E7%B3%BB%E5%88%97%E7%9B%AE%E5%BD%95"
)

// GetArticles godoc
// @Summary 获取文章列表
// @Produce  json
// @Param title path string false "标题"
// @Param tag_id path int false "标签ID"
// @Param state path int false "状态 0和1"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/articles [get]
func GetArticles(c *gin.Context) {
	appG := app.Gin{C: c}

	data := make(map[string]interface{})
	var valid validation.Validation

	title := c.Query("title")

	state := -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
	}

	tagId := -1
	if arg := c.Query("tag_id"); arg != "" {
		tagId = com.StrTo(arg).MustInt()
		valid.Min(tagId, 1, "tag_id").Message("标签ID必须大于0")
	}

	code := e.INVALID_PARAMS
	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	articleService := article_service.Article{
		Title:    title,
		State:    state,
		TagID:    tagId,
		PageNum:  util.GetPage(c),
		PageSize: settings.AppSetting.PageSize,
	}
	articles, err := articleService.GetAll()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_GET_ARTICLES_FAIL, nil)
		return
	}

	total, err := articleService.Count()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_COUNT_ARTICLE_FAIL, nil)
	}

	data["list"] = articles
	data["total"] = total

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}

// GetArticle godoc
// @Summary 获取单篇文章
// @Produce  json
// @Param id path int true "文章ID"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/articles/{id} [get]
func GetArticle(c *gin.Context) {
	appG := app.Gin{C: c}
	id := com.StrTo(c.Param("id")).MustInt()

	var valid validation.Validation
	valid.Min(id, 1, "id").Message("文章ID必须大于0")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	articleService := article_service.Article{ID: id}
	existed, err := articleService.ExistedById()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_CHECK_EXIST_ARTICLE_FAIL, nil)
		return
	}
	if !existed {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_ARTICLE, nil)
		return
	}

	article, err := articleService.Get()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_GET_ARTICLE_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, article)
}

// AddArticle 添加文章
func AddArticle(c *gin.Context) {
	var article models.Article
	appG := app.Gin{C: c}
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
	valid.MaxSize(article.CoverImageUrl, 255, "cover_image_url").Message("封面图不能超过255个字符")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	articleService := article_service.Article{
		TagID:         article.TagId,
		Title:         article.Title,
		Desc:          article.Desc,
		CreatedBy:     article.CreatedBy,
		State:         article.State,
		CoverImageUrl: article.CoverImageUrl,
	}

	tagService := tag_service.Tag{ID: article.ID}
	tagExisted, err := tagService.ExistedTagById()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EXIST_TAG_FAIL, nil)
		return
	}
	if !tagExisted {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_TAG, nil)
	}

	err = articleService.Add()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_ADD_ARTCLIE_FAIL, nil)
	}
}

// UpdateArticle 更新文章
func UpdateArticle(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()
	var article models.Article
	var valid validation.Validation
	var appG = app.Gin{C: c}
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

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	articleService := article_service.Article{
		ID:         id,
		TagID:      article.TagId,
		Title:      article.Title,
		Desc:       article.Desc,
		ModifiedBy: article.ModifiedBy,
		State:      article.State,
	}
	articleExisted, err := articleService.ExistedById()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_CHECK_EXIST_ARTICLE_FAIL, nil)
		return
	}
	if !articleExisted {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_ARTICLE, nil)
		return
	}

	tagService := tag_service.Tag{ID: article.TagId}
	tagExisted, err := tagService.ExistedTagById()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EXIST_TAG_FAIL, nil)
		return
	}
	if !tagExisted {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_TAG, nil)
		return
	}

	err = articleService.Update()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_ARTICLES_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// DeleteArticle 删除文章
func DeleteArticle(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()
	appG := app.Gin{C: c}
	var valid validation.Validation

	valid.Min(id, 1, "id").Message("文章ID必须大于0")
	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	articleService := article_service.Article{ID: id}
	articleExisted, err := articleService.ExistedById()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_CHECK_EXIST_ARTICLE_FAIL, nil)
		return
	}
	if !articleExisted {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_ARTICLE, nil)
		return
	}

	err = articleService.Delete()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_DELETE_ARTCLIE_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// GenerateArticlePoster 生成二维码
func GenerateArticlePoster(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)

	qrc := qrcode.NewQrCode(QRCODE_URL, 300, 300, qr.M, qr.Auto)
	path := qrcode.GetQrCodeFullPath()
	_, _, err := qrc.Encode(path)
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
