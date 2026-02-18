package dao

import (
	"fmt"

	"github.com/chuji555/homework-system/models"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// 初始化数据库连接
func InitDB() {
	var err error
	// 读取配置
	dsn := viper.GetString("mysql.dsn")
	maxOpenConns := viper.GetInt("mysql.max_open_conns")
	maxIdleConns := viper.GetInt("mysql.max_idle_conns")

	// 连接MySQL
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		// 显示SQL日志（新手调试用）
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(fmt.Sprintf("数据库连接失败：%v", err))
	}

	// 设置连接池
	sqlDB, _ := DB.DB()
	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetMaxIdleConns(maxIdleConns)

	// 自动迁移建表
	err = DB.AutoMigrate(
		&models.User{},
		&models.Homework{},
		&models.Submission{},
	)
	if err != nil {
		panic(fmt.Sprintf("建表失败：%v", err))
	}
	fmt.Println("数据库初始化成功！")
}
