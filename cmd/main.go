package main

import (
	"ffly-baisc/internal/config"
	"ffly-baisc/internal/db"
	"ffly-baisc/internal/router"
	"log"
)

func main() {
	configPath := "config/config.yaml"
	if err := config.Init(configPath); err != nil {
		log.Fatalf("Failed to init config: %v\n", err)
	}

	// 初始化数据库
	db.InitDB()

	// 初始化路由
	router.InitRouter()
}
