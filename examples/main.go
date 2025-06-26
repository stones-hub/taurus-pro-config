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
	log.Printf("Database Host: %s", cfg.GetString("database.host"))
	log.Printf("Database Port: %d", cfg.GetInt("database.port"))
	log.Printf("Debug Mode: %v", cfg.GetBool("debug"))
	log.Printf("API Keys: %v", cfg.GetStringSlice("api.keys"))
}
