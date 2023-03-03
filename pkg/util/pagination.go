package util

import (
	"gin-blog-example/settings"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
)

// GetPage 获取分页
func GetPage(c *gin.Context) int {
	result := 0
	page, _ := com.StrTo(c.Query("page")).Int()
	if page > 0 {
		result = (page - 1) * settings.AppSetting.PageSize
	}
	return result
}
