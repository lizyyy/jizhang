package main

import (
	"fmt"
	"go-accounting/config"
	"go-accounting/routes"
	"log"
	"os"
	"path/filepath"
)

var publicPath string

func main() {
	execPath, err := os.Executable()
	if err != nil {
		log.Fatal("获取执行路径失败:", err)
	}

	execDir := filepath.Dir(execPath)
	projectRoot := execDir

	if _, err := os.Stat(filepath.Join(execDir, "public")); os.IsNotExist(err) {
		projectRoot = filepath.Dir(execDir)
		if _, err := os.Stat(filepath.Join(projectRoot, "public")); os.IsNotExist(err) {
			cwd, _ := os.Getwd()
			if _, err := os.Stat(filepath.Join(cwd, "public")); err == nil {
				projectRoot = cwd
			} else {
				projectRoot = filepath.Dir(cwd)
			}
		}
	}

	publicPath = filepath.Join(projectRoot, "public")

	dataPath := filepath.Join(projectRoot, "data")
	if err := os.MkdirAll(dataPath, 0755); err != nil {
		log.Fatal("创建数据目录失败:", err)
	}

	if err := config.InitDatabase(dataPath); err != nil {
		log.Fatal("初始化数据库失败:", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "7002"
	}

	r := routes.SetupRouter(publicPath)

	fmt.Printf("服务器运行在 http://localhost:%s\n", port)
	fmt.Printf("静态文件目录: %s\n", publicPath)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("启动服务器失败:", err)
	}
}
