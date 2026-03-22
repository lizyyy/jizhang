package main

import (
	"accounting-app/data"
	"accounting-app/routes"
	"fmt"
	"log"
	"path/filepath"
)

func main() {
	// 数据库文件路径
	dbPath := filepath.Join(".", "data", "accounting.db")
	fmt.Printf("数据库路径: %s\n", dbPath)

	// 初始化数据库
	if err := data.InitDatabase(dbPath); err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}
	fmt.Println("数据库连接成功")

	// 执行数据库迁移
	if err := data.Migrate(data.GetDB()); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}
	fmt.Println("数据库迁移完成")

	// 初始化默认数据
	categoryData := data.NewCategoryData()
	if err := categoryData.InitDefaultData(); err != nil {
		log.Fatalf("默认数据初始化失败: %v", err)
	}
	fmt.Println("默认数据初始化完成")

	// 设置路由
	router := routes.SetupRouter()

	// 使用端口号7001（根据分支名称）
	port := "7001"
	fmt.Printf("服务器运行在 http://localhost:%s\n", port)

	// 启动服务
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
