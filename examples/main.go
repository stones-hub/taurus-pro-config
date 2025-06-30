package main

import (
	"log"

	"github.com/stones-hub/taurus-pro-config/pkg/config"
)

func main() {
	// 创建配置实例
	cfg := config.New()
	cfg.PrintEnable = true // 启用配置打印

	// 初始化配置
	if err := cfg.Initialize("config/", ".env.local"); err != nil {
		log.Fatalf("Failed to initialize config: %v", err)
	}

	// 使用配置
	log.Printf("App Name: %s", cfg.GetString("app_name"))
	log.Printf("Version: %s", cfg.GetString("version"))
	log.Printf("HTTP Address: %s", cfg.GetString("http.address"))
	log.Printf("HTTP Port: %d", cfg.GetInt("http.port"))
	log.Printf("HTTP Authorization: %s", cfg.GetString("http.authorization"))
	log.Printf("HTTP Read Timeout: %d", cfg.GetInt("http.read_timeout"))
}
