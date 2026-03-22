package routes

import (
	"go-accounting/config"
	"go-accounting/controllers"
	"go-accounting/data"
	"go-accounting/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRouter(publicPath string) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middlewares.LoggerMiddleware())

	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	categoryRepo := data.NewCategoryRepository(config.DB)
	recordRepo := data.NewRecordRepository(config.DB)

	categoryController := controllers.NewCategoryController(categoryRepo)
	recordController := controllers.NewRecordController(recordRepo)

	api := r.Group("/api")
	{
		categories := api.Group("/categories")
		{
			categories.GET("", categoryController.GetAll)
			categories.GET("/:id", categoryController.GetByID)
			categories.POST("", categoryController.Create)
			categories.PUT("/:id", categoryController.Update)
			categories.DELETE("/:id", categoryController.Delete)
			categories.POST("/subcategories", categoryController.CreateSubcategory)
			categories.PUT("/subcategories/:id", categoryController.UpdateSubcategory)
			categories.DELETE("/subcategories/:id", categoryController.DeleteSubcategory)
		}

		records := api.Group("/records")
		{
			records.GET("", recordController.GetAll)
			records.GET("/statistics", recordController.GetStatistics)
			records.GET("/:id", recordController.GetByID)
			records.POST("", recordController.Create)
			records.PUT("/:id", recordController.Update)
			records.DELETE("/:id", recordController.Delete)
		}
	}

	r.Static("/js", publicPath+"/js")
	r.Static("/css", publicPath+"/css")
	r.StaticFile("/", publicPath+"/index.html")

	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"Errno":     404,
			"Errmsg":    "接口不存在",
			"ErrorCode": "NOT_FOUND",
		})
	})

	return r
}
