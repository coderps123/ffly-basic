package main

import (
	"ffly-baisc/internal/config"
	"ffly-baisc/internal/db"
	"ffly-baisc/internal/router"
	"fmt"
	"log"
	"os"
)

func main() {
	// 获取环境变量
	appEnv := os.Getenv("APP_ENV") // APP_ENV 需要在环境变量中设置， 如：windows set APP_ENV=development，如：linux export APP_ENV=production
	if appEnv == "" {
		appEnv = "development"
		fmt.Println("APP_ENV is not set, use default value: development")
	}

	// 构建配置文件路径
	configPath := fmt.Sprintf("config/config_%s.yaml", appEnv)

	// 初始化配置
	if err := config.Init(configPath); err != nil {
		log.Fatalf("Failed to init config: %v\n", err)
	}

	// 初始化数据库
	db.InitDB()

	// 初始化路由服务
	router.Init()
}
