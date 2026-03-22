package routes

import (
	"accounting-app/controllers"
	"accounting-app/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/cors"
)

// SetupRouter 设置路由
func SetupRouter() *gin.Engine {
	router := gin.Default()

	// 使用日志中间件
	router.Use(middleware.LoggingMiddleware())

	// CORS 配置 - 与原Node.js项目一致
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		ExposedHeaders:   []string{"Content-Length"},
		AllowCredentials: true,
	})

	// 使用CORS中间件
	router.Use(func(ctx *gin.Context) {
		c.HandlerFunc(ctx.Writer, ctx.Request)
		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(http.StatusOK)
			return
		}
		ctx.Next()
	})

	// 静态文件服务 - 提供前端public目录
	router.Static("/css", "./public/css")
	router.Static("/js", "./public/js")
	router.StaticFile("/", "./public/index.html")

	// API路由组
	api := router.Group("/api")
	{
		// 分类路由
		categories := api.Group("/categories")
		{
			categoryController := controllers.NewCategoryController()

			// 获取所有分类
			categories.GET("", categoryController.GetAllCategories)
			// 获取单个分类
			categories.GET("/:id", categoryController.GetCategoryByID)
			// 创建分类
			categories.POST("", categoryController.CreateCategory)
			// 更新分类
			categories.PUT("/:id", categoryController.UpdateCategory)
			// 删除分类
			categories.DELETE("/:id", categoryController.DeleteCategory)
			// 创建子分类
			categories.POST("/subcategories", categoryController.CreateSubcategory)
			// 更新子分类
			categories.PUT("/subcategories/:id", categoryController.UpdateSubcategory)
			// 删除子分类
			categories.DELETE("/subcategories/:id", categoryController.DeleteSubcategory)
		}

		// 记账记录路由
		records := api.Group("/records")
		{
			recordController := controllers.NewRecordController()

			// 获取所有记账记录
			records.GET("", recordController.GetAllRecords)
			// 获取统计信息
			records.GET("/statistics", recordController.GetStatistics)
			// 获取单个记账记录
			records.GET("/:id", recordController.GetRecordByID)
			// 创建记账记录
			records.POST("", recordController.CreateRecord)
			// 更新记账记录
			records.PUT("/:id", recordController.UpdateRecord)
			// 删除记账记录
			records.DELETE("/:id", recordController.DeleteRecord)
		}
	}

	// 404处理
	router.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"Errno":     404,
			"Errmsg":    "接口不存在",
			"ErrorCode": "NOT_FOUND",
		})
	})

	return router
}
