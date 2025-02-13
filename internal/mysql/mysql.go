package mysql

import (
	"ffly-baisc/internal/config"
	"fmt"
	"log"
	"strconv"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitMySql() (*gorm.DB, error) {
	MySqlConfig := config.GlobalConfig.MySql
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		MySqlConfig.User,
		MySqlConfig.Password,
		MySqlConfig.Host,
		strconv.Itoa(MySqlConfig.Port),
		MySqlConfig.Database)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	// // 自动迁移数据库表
	// if err := AutoMigrate(db); err != nil {
	// 	log.Fatalf("Failed to migrate database: %v\n", err)
	// }

	// 获取底层的sql.DB对象
	sqlDB, _ := db.DB()

	// 设置最大空闲连接数
	sqlDB.SetMaxIdleConns(MySqlConfig.MaxIdleConns)

	// 设置最大打开连接数
	sqlDB.SetMaxOpenConns(MySqlConfig.MaxOpenConns)

	// 设置连接的最大存活时间
	sqlDB.SetConnMaxLifetime(time.Duration(MySqlConfig.ConnectionMaxLifetime) * time.Second) // n分钟

	if err != nil {
		log.Fatalf("Failed to init mysql: %v\n", err) // 这里如果出现错误，会终止整个程序的运行
	}

	// 将db赋值给全局变量
	DB = db

	return db, err
}

// func AutoMigrate(db *gorm.DB) error {
// 	return db.AutoMigrate(
// 		&model.User{},
// 		&model.Role{},
// 		&model.Permission{},
// 		&model.RolePermission{},
// 		&model.System{},
// 	)
// }
