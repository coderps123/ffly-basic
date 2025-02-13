package main

import (
	"ffly-baisc/internal/config"
	"ffly-baisc/internal/mysql"
	"ffly-baisc/internal/router"
	"log"
)

func main() {
	configPath := "config/config.yaml"
	if err := config.Init(configPath); err != nil {
		log.Fatalf("Failed to init config: %v\n", err)
	}

	// 初始化mysql
	// GORM会自动管理连接池和连接的生命周期。所以不需要我们手动关闭连接。
	mysql.InitMySql()

	// 初始化路由
	router.InitRouter()
}
