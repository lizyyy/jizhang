package main

import (
	"accounting-app/data"
	"accounting-app/routes"
	"fmt"
	"os"
)

func main() {
	// 初始化数据库
	if err := data.InitDatabase(); err != nil {
		fmt.Printf("数据库初始化失败: %v\n", err)
		os.Exit(1)
	}

	// 获取数据库实例
	db := data.GetDB()

	// 设置路由
	r := routes.SetupRouter(db)

	// 从环境变量或分支名获取端口
	port := os.Getenv("PORT")
	if port == "" {
		// 使用分支名作为端口号（7003）
		port = "7003"
	}

	fmt.Printf("服务器运行在 http://localhost:%s\n", port)

	// 启动服务器
	if err := r.Run(":" + port); err != nil {
		fmt.Printf("服务器启动失败: %v\n", err)
		os.Exit(1)
	}
}
