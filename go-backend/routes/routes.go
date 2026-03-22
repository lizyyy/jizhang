package routes

import (
	"accounting-app/controllers"
	"accounting-app/data"
	"accounting-app/middleware"
	"accounting-app/utils"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRoutes 设置所有路由
func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	// 初始化数据层
	categoryData := data.NewCategoryData(db)
	recordData := data.NewRecordData(db)

	// 初始化控制器
	categoryController := controllers.NewCategoryController(categoryData)
	recordController := controllers.NewRecordController(recordData)

	// API 路由组
	api := r.Group("/api")
	{
		// 分类路由
		categories := api.Group("/categories")
		{
			categories.GET("", categoryController.GetAllCategories)
			categories.GET("/:id", categoryController.GetCategoryByID)
			categories.POST("", categoryController.CreateCategory)
			categories.PUT("/:id", categoryController.UpdateCategory)
			categories.DELETE("/:id", categoryController.DeleteCategory)

			// 子分类路由
			categories.POST("/subcategories", categoryController.CreateSubcategory)
			categories.PUT("/subcategories/:id", categoryController.UpdateSubcategory)
			categories.DELETE("/subcategories/:id", categoryController.DeleteSubcategory)
		}

		// 记账记录路由
		records := api.Group("/records")
		{
			records.GET("", recordController.GetAllRecords)
			records.GET("/statistics", recordController.GetStatistics)
			records.GET("/:id", recordController.GetRecordByID)
			records.POST("", recordController.CreateRecord)
			records.PUT("/:id", recordController.UpdateRecord)
			records.DELETE("/:id", recordController.DeleteRecord)
		}
	}

	// 静态文件服务 - 指向原项目的 public 目录
	publicPath := filepath.Join("..", "public")
	r.Static("/css", filepath.Join(publicPath, "css"))
	r.Static("/js", filepath.Join(publicPath, "js"))

	// 首页路由
	r.GET("/", func(c *gin.Context) {
		indexPath := filepath.Join(publicPath, "index.html")
		c.File(indexPath)
	})

	// 404 处理
	r.NoRoute(func(c *gin.Context) {
		utils.JSONError(c, 404, "接口不存在", "NOT_FOUND")
	})
}

// SetupRouter 创建并配置 Gin 引擎
func SetupRouter(db *gorm.DB) *gin.Engine {
	// 设置 Gin 模式
	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "" {
		ginMode = gin.DebugMode
	}
	gin.SetMode(ginMode)

	// 创建 Gin 引擎
	r := gin.New()

	// 使用自定义日志中间件
	r.Use(middleware.LoggerMiddleware())

	// 使用 Recovery 中间件
	r.Use(gin.Recovery())

	// CORS 配置（与原 Node.js 项目一致）
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	// 设置路由
	SetupRoutes(r, db)

	return r
}
