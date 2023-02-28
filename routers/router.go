package routers

import (
	"gin-blog-example/pkg/e"
	v1 "gin-blog-example/routers/api/v1"
	"gin-blog-example/settings"
	"github.com/gin-gonic/gin"
)

// InitRouter 初始化路由
func InitRouter() *gin.Engine {
	r := gin.New() // 不使用gin.Default()，为了不打印Warning
	//r := gin.Default()
	r.Use(gin.Logger(), gin.Recovery())

	// 设置模式
	gin.SetMode(settings.RunMode)
	r.GET("/test", func(c *gin.Context) {
		c.JSON(e.SUCCESS, gin.H{
			"message": "Ok",
		})
	})

	// 注册路由
	apiV1 := r.Group("/v1")
	{
		apiV1.GET("/tags", v1.GetTags)
		apiV1.POST("/tags", v1.AddTag)
		apiV1.PUT("/tags/:id", v1.UpdateTag)
		apiV1.DELETE("/tags/:id", v1.DeleteTag)
	}

	return r
}
