package v1

import (
	"gin-blog-example/models"
	"gin-blog-example/pkg/e"
	"gin-blog-example/pkg/logging"
	"gin-blog-example/pkg/util"
	"gin-blog-example/settings"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"log"
	"net/http"
)

// GetTags 获取标签列表
func GetTags(c *gin.Context) {
	name := c.Query("name")
	maps := make(map[string]interface{})
	data := make(map[string]interface{})

	if name != "" {
		maps["name"] = name
	}

	var state = -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		maps["state"] = state
	}

	code := e.SUCCESS
	data["list"] = models.GetTags(util.GetPage(c), settings.PageSize, maps)
	data["total"] = models.GetTagTotal(maps)
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}

// AddTag godoc
// @Summary 新增文章标签
// @Produce  json
// @Param tag body models.Tag true "Add Tag"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/v1/tags [post]
func AddTag(c *gin.Context) {
	// 绑定JSON数据
	var tag models.Tag
	_ = c.BindJSON(&tag)

	code := e.INVALID_PARAMS
	// 参数验证
	valid := validation.Validation{}
	valid.Required(tag.Name, "name").Message("名称不能为空")
	valid.MaxSize(tag.Name, 100, "name").Message("名称最长为100字符")
	valid.Required(tag.CreatedBy, "create_by").Message("创建人不能为空")
	valid.MaxSize(tag.CreatedBy, 100, "create_by").Message("创建人最长为100字符")
	valid.Range(tag.State, 0, 1, "state").Message("状态只允许0或1")

	if !valid.HasErrors() {
		// 标签名存在校验
		if models.ExistedTagByName(tag.Name) {
			code = e.ERROR_EXIST_TAG
		} else {
			if !models.AddTag(tag.Name, tag.State, tag.CreatedBy) {
				code = e.ERROR_ADD_TAG_FAIL
			}
			code = e.SUCCESS
		}
	} else {
		for _, err := range valid.Errors {
			logging.Info(err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.MsgFlags[code],
	})
}

// UpdateTag 更新标签
func UpdateTag(c *gin.Context) {
	// 通过
	id := com.StrTo(c.Param("id")).MustInt()
	var tag models.Tag
	_ = c.BindJSON(&tag)

	valid := validation.Validation{}
	valid.Required(id, "id").Message("标签ID不能为空")
	valid.Required(tag.ModifiedBy, "modified_by").Message("修改人不能为空")
	valid.MaxSize(tag.ModifiedBy, 100, "modified_by").Message("修改人最长为100字符")
	valid.MaxSize(tag.Name, 100, "name").Message("标签名称最长为100字符")
	valid.Range(tag.State, 0, 1, "state").Message("标签状态只能是0和1")

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		if !models.ExistedTagById(id) {
			code = e.ERROR_NOT_EXIST_TAG
		} else {
			models.UpdateTag(id, tag)
			code = e.SUCCESS
		}
	} else {
		for _, v := range valid.Errors {
			log.Println(v.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
	})
}

// DeleteTag 删除标签
func DeleteTag(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Required(id, "id").Message("标签ID不能为空")

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		if !models.ExistedTagById(id) {
			code = e.ERROR_NOT_EXIST_TAG
		} else {
			models.DeleteTag(id)
			code = e.SUCCESS
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
	})
}
