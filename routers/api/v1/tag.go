package v1

import (
	"gin-blog-example/models"
	"gin-blog-example/pkg/e"
	"gin-blog-example/pkg/util"
	"gin-blog-example/settings"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
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

func AddTag(c *gin.Context) {

}

func UpdateTag(c *gin.Context) {

}

func DeleteTag(c *gin.Context) {

}
