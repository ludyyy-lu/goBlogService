package routers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/ludyyy-lu/goBlogService/docs"
	"github.com/ludyyy-lu/goBlogService/global"
	"github.com/ludyyy-lu/goBlogService/internal/middleware"
	"github.com/ludyyy-lu/goBlogService/internal/routers/api"
	v1 "github.com/ludyyy-lu/goBlogService/internal/routers/api/v1"
	"github.com/ludyyy-lu/goBlogService/pkg/limiter"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

var methodLimiters = limiter.NewMethodLimiter().AddBuckets(
	limiter.LimiterBucketRule{
		Key:          "/auth",
		FillInterval: time.Second,
		Capacity:     10,
		Quantum:      10,
	},
)

func NewRouter() *gin.Engine {
	r := gin.New()
	if global.ServerSetting.RunMode == "debug" {
		r.Use(gin.Logger())
		r.Use(gin.Recovery())
	} else {
		r.Use(middleware.AccessLog())
		r.Use(gin.Recovery())
	}

	r.Use(middleware.RateLimiter(methodLimiters))
	r.Use(middleware.ContextTimeout(60 * time.Second))
	r.Use(middleware.Translations())

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/auth", api.GetAuth)
	article := v1.NewArticle()
	tag := v1.NewTag()

	upload := api.NewUpload()
	r.POST("/upload/file", upload.UploadFile)
	//文件服务只有提供静态资源的访问，才能在外部请求本项目HttpServer时同时提供静态资源的访问
	r.StaticFS("/static", http.Dir(global.AppSetting.UploadSavePath))
	apiv1 := r.Group("/api/v1")
	apiv1.Use(middleware.JWT())
	{
		apiv1.POST("/tags", tag.Create)
		apiv1.DELETE("/tags/:id", tag.Delete)
		apiv1.PUT("/tags/:id", tag.Update)
		apiv1.PATCH("/tags/:id/state", tag.Update)
		apiv1.GET("/tags", tag.List)

		apiv1.POST("/articles", article.Create)
		apiv1.DELETE("/articles/:id", article.Delete)
		apiv1.PUT("/articles/:id", article.Update)
		apiv1.PATCH("/articles/:id/state", article.Update)
		apiv1.GET("/articles/:id", article.Get)
		apiv1.GET("/articles", article.List)
	}
	return r
}
