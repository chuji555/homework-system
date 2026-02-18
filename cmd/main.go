package main

import (
	"fmt"
	"log"

	"github.com/chuji555/homework-system/dao"
	"github.com/chuji555/homework-system/router"
	"github.com/spf13/viper"
)

func init() {
	// 初始化配置文件
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Sprintf("读取配置文件失败：%v", err))
	}
	// 初始化数据库
	dao.InitDB()
}
func main() {
	// 初始化路由
	r := router.InitRouter()
	// 启动服务
	port := viper.GetString("server.port")
	log.Printf("服务器启动成功，监听端口：%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("服务器启动失败：%v", err)
	}
}
