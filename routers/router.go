package routers

import (
	"gin-blog-example/middleware"
	"gin-blog-example/routers/api"
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

	// 登录
	r.POST("/auth", api.GetAuth)

	// 注册路由
	apiV1 := r.Group("/v1")
	// 针对部分接口进行jwt鉴权，记住这里只能写到具体接口之前，否则不生效
	apiV1.Use(middleware.JWT())
	{
		// 标签
		apiV1.GET("/tags", v1.GetTags)
		apiV1.POST("/tags", v1.AddTag)
		apiV1.PUT("/tags/:id", v1.UpdateTag)
		apiV1.DELETE("/tags/:id", v1.DeleteTag)

		// 文章
		apiV1.GET("/articles", v1.GetArticles)
		apiV1.GET("/articles/:id", v1.GetArticle)
		apiV1.POST("/articles", v1.AddArticle)
		apiV1.PUT("/articles/:id", v1.UpdateArticle)
		apiV1.DELETE("/articles/:id", v1.DeleteArticle)
	}

	return r
}
