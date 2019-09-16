package routers

import (
	"blog/middleware/cors"
	"blog/pkg/setting"
	"blog/routers/api"
	v1 "blog/routers/api/v1"
	"github.com/gin-gonic/gin"
)

func InitRouter()*gin.Engine{
	r:=gin.New()
	//r.LoadHTMLGlob("views/**/*")
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	gin.SetMode(setting.RunModel)
	r.GET("/auth",api.GetAuth)

	//r.GET("/test", func(c *gin.Context) {
	//	//c.JSON(200,gin.H{
	//	//	"php":"echp 'hello world'",
	//	//})
	//	c.HTML(http.StatusOK,"index.tmpl",gin.H{
	//		"title":"main website",
	//	})
	//})
	apiV1 := r.Group("/api/v1")
	apiV1.Use(cors.Cors())
	//apiV1.Use(jwt.JWT())
	{
		apiV1.GET("/tags", v1.GetTags)
		apiV1.POST("/tags", v1.AddTag)
		apiV1.PUT("/tags/:id", v1.EditTag)
		apiV1.DELETE("/tags/:id", v1.DeleteTag)

		//获取文章列表
		apiV1.GET("/articles", v1.GetArticles)
		//获取指定文章
		apiV1.GET("/articles/:id", v1.GetArticle)
		//新建文章
		apiV1.POST("/articles", v1.AddArticle)
		//更新指定文章
		apiV1.PUT("/articles/:id", v1.EditArticle)
		//删除指定文章
		apiV1.DELETE("/articles/:id", v1.DeleteArticle)
	}


	return r
}
