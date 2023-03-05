package v1

import (
	"gin-blog-example/models"
	"gin-blog-example/pkg/app"
	"gin-blog-example/pkg/e"
	"gin-blog-example/pkg/util"
	"gin-blog-example/services/tag_service"
	"gin-blog-example/settings"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"net/http"
)

// GetTags 获取标签列表
func GetTags(c *gin.Context) {
	name := c.Query("name")
	state := -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
	}

	data := make(map[string]interface{})

	appG := app.Gin{C: c}

	tagService := tag_service.Tag{
		Name:     name,
		State:    state,
		PageNum:  util.GetPage(c),
		PageSize: settings.AppSetting.PageSize,
	}

	tags, err := tagService.GetAll()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_TAG_FAIL, nil)
		return
	}

	count, err := tagService.Count()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_COUNT_TAG_FAIL, nil)
		return
	}

	data["list"] = tags
	data["total"] = count
	appG.Response(http.StatusOK, e.SUCCESS, data)
}

// AddTag godoc
// @Summary 新增文章标签
// @Produce  json
// @Param tag body models.Tag true "Add Tag"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/v1/tags [post]
func AddTag(c *gin.Context) {
	// 绑定JSON数据
	var (
		tag  models.Tag
		appG = app.Gin{C: c}
	)
	_ = c.BindJSON(&tag)

	// 参数验证
	valid := validation.Validation{}
	valid.Required(tag.Name, "name").Message("名称不能为空")
	valid.MaxSize(tag.Name, 100, "name").Message("名称最长为100字符")
	valid.Required(tag.CreatedBy, "create_by").Message("创建人不能为空")
	valid.MaxSize(tag.CreatedBy, 100, "create_by").Message("创建人最长为100字符")
	valid.Range(tag.State, 0, 1, "state").Message("状态只允许0或1")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	tagService := tag_service.Tag{
		Name:     tag.Name,
		CreateBy: tag.CreatedBy,
		State:    tag.State,
	}

	tagNameExisted, err := tagService.ExistedTagByName()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EXIST_TAG_FAIL, nil)
		return
	}
	if tagNameExisted {
		appG.Response(http.StatusOK, e.ERROR_EXIST_TAG, nil)
		return
	}

	err = tagService.Add()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_ADD_TAG_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// UpdateTag 更新标签
func UpdateTag(c *gin.Context) {
	// 通过
	id := com.StrTo(c.Param("id")).MustInt()
	var (
		tag  models.Tag
		appG = app.Gin{C: c}
	)
	_ = c.BindJSON(&tag)

	valid := validation.Validation{}
	valid.Required(id, "id").Message("标签ID不能为空")
	valid.Required(tag.ModifiedBy, "modified_by").Message("修改人不能为空")
	valid.MaxSize(tag.ModifiedBy, 100, "modified_by").Message("修改人最长为100字符")
	valid.MaxSize(tag.Name, 100, "name").Message("标签名称最长为100字符")
	valid.Range(tag.State, 0, 1, "state").Message("标签状态只能是0和1")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	tagService := tag_service.Tag{
		ID:         id,
		Name:       tag.Name,
		ModifiedBy: tag.CreatedBy,
		State:      tag.State,
	}
	tagExisted, err := tagService.ExistedTagById()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EXIST_TAG_FAIL, nil)
		return
	}
	if !tagExisted {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_TAG, nil)
		return
	}

	err = models.UpdateTag(id, tag)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_UPDATE_TAG_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// DeleteTag 删除标签
func DeleteTag(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()
	appG := app.Gin{C: c}

	valid := validation.Validation{}
	valid.Required(id, "id").Message("标签ID不能为空")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	tagService := tag_service.Tag{ID: id}
	err := tagService.Delete()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_DELETE_TAG_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
