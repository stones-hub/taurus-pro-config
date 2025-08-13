# Taurus Pro Config

[![Go Version](https://img.shields.io/badge/Go-1.24.2+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/stones-hub/taurus-pro-config)](https://goreportcard.com/report/github.com/stones-hub/taurus-pro-config)

一个功能强大、灵活易用的Go语言配置管理库，支持多种配置文件格式和环境变量集成。

## 📋 目录

- [特性](#特性)
- [安装](#安装)
- [快速开始](#快速开始)
- [配置格式支持](#配置格式支持)
- [API 文档](#api-文档)
- [配置示例](#配置示例)
- [环境变量支持](#环境变量支持)
- [贡献指南](#贡献指南)
- [许可证](#许可证)

## ✨ 特性

- 🔧 **多格式支持**: 支持 YAML、JSON、TOML、XML 等多种配置文件格式
- 🌍 **环境变量集成**: 支持环境变量占位符和 .env 文件加载
- 📁 **目录递归加载**: 支持从目录递归加载所有配置文件
- 🔄 **配置合并**: 智能合并多个配置文件，支持嵌套结构
- 🎯 **类型安全**: 提供类型安全的配置值获取方法
- 🚀 **高性能**: 基于内存的配置管理，快速访问
- 🛠️ **易于使用**: 简洁的 API 设计，开箱即用
- 📊 **调试友好**: 支持配置打印和 JSON 导出

## 📦 安装

```bash
go get github.com/stones-hub/taurus-pro-config
```

## 🚀 快速开始

### 基本用法

```go
package main

import (
    "log"
    "github.com/stones-hub/taurus-pro-config/pkg/config"
)

func main() {
    // 创建配置实例
    cfg := config.New()
    
    // 初始化配置
    if err := cfg.Initialize("config/", ".env.local"); err != nil {
        log.Fatalf("Failed to initialize config: %v", err)
    }
    
    // 获取配置值
    appName := cfg.GetString("app_name")
    port := cfg.GetInt("http.port")
    
    log.Printf("App: %s, Port: %d", appName, port)
}
```

### 使用选项模式

```go
// 启用配置打印
cfg := config.New(config.WithPrintEnable(true))
```

## 📁 配置格式支持

### 支持的格式

| 格式 | 扩展名 | 说明 |
|------|--------|------|
| YAML | `.yaml`, `.yml` | 推荐使用，支持注释和复杂结构 |
| JSON | `.json` | 标准格式，易于程序处理 |
| TOML | `.toml` | 人类友好的配置格式 |
| XML | `.xml` | 传统格式，支持复杂结构 |

### 配置文件结构

```
config/
├── config.yaml          # 主配置文件
├── autoload/            # 自动加载目录
│   ├── http/
│   │   └── http.yaml    # HTTP 服务配置
│   ├── db/
│   │   └── db.yaml      # 数据库配置
│   ├── redis/
│   │   └── redis.toml   # Redis 配置
│   └── logger/
│       └── logger.yaml  # 日志配置
└── .env.local           # 环境变量文件
```

## 📚 API 文档

### 核心方法

#### 创建和初始化

```go
// 创建新的配置实例
func New(options ...ConfigOption) *Config

// 初始化配置
func (c *Config) Initialize(configPath string, env string) error
```

#### 配置值获取

```go
// 获取任意类型的值
func (c *Config) Get(key string) interface{}

// 获取字符串值
func (c *Config) GetString(key string) string

// 获取整数值
func (c *Config) GetInt(key string) int

// 获取布尔值
func (c *Config) GetBool(key string) bool

// 获取浮点数值
func (c *Config) GetFloat64(key string) float64

// 获取字符串映射
func (c *Config) GetStringMap(key string) map[string]interface{}

// 获取字符串切片
func (c *Config) GetStringSlice(key string) []string
```

#### 工具方法

```go
// 转换为 JSON 字符串
func (c *Config) ToJSONString() string

// 合并配置映射
func (c *Config) MergeMap(data map[string]interface{})
```

### 配置选项

```go
// 启用配置打印
func WithPrintEnable(enable bool) ConfigOption
```

## ⚙️ 配置示例

### 主配置文件 (config.yaml)

```yaml
# 系统版本
version: "${VERSION:v1.0.0}"
app_name: "${APP_NAME:taurus}"
print_enable: true
```

### HTTP 服务配置 (http/http.yaml)

```yaml
http:
  address: "${SERVER_ADDRESS:0.0.0.0}"
  port: ${SERVER_PORT:8080}
  read_timeout: 60
  write_timeout: 60
  idle_timeout: 60
  authorization: "Bearer ${AUTHORIZATION:123456}"
```

### 数据库配置 (db/db.yaml)

```yaml
database:
  driver: "${DB_DRIVER:mysql}"
  host: "${DB_HOST:localhost}"
  port: ${DB_PORT:3306}
  username: "${DB_USERNAME:root}"
  password: "${DB_PASSWORD:}"
  database: "${DB_NAME:taurus}"
  max_open_conns: 100
  max_idle_conns: 10
```

### Redis 配置 (redis/redis.toml)

```toml
[redis]
host = "${REDIS_HOST:localhost}"
port = 6379
password = "${REDIS_PASSWORD:}"
db = 0
pool_size = 10
```

## 🌍 环境变量支持

### 环境变量占位符

配置文件中支持使用环境变量占位符，格式为 `${ENV_VAR:default_value}`：

```yaml
# 使用环境变量，如果未设置则使用默认值
database:
  host: "${DB_HOST:localhost}"
  port: ${DB_PORT:3306}
```

### .env 文件支持

支持加载 `.env` 文件来设置环境变量：

```bash
# .env.local
VERSION=v1.0.0
APP_NAME=taurus-pro
SERVER_ADDRESS=0.0.0.0
SERVER_PORT=8080
DB_HOST=localhost
DB_PORT=3306
REDIS_HOST=localhost
```

## 🔧 高级用法

### 配置合并策略

当多个配置文件包含相同键时，库会智能合并配置：

- 对于简单值（字符串、数字、布尔值），后加载的配置会覆盖先前的值
- 对于嵌套的映射结构，会递归合并，保留所有配置项

### 目录递归加载

支持从目录递归加载所有配置文件：

```go
// 加载 config/ 目录下的所有配置文件
cfg.Initialize("config/", "")
```

### 类型转换

库使用 `github.com/spf13/cast` 进行安全的类型转换，支持：

- 字符串 ↔ 数字
- 字符串 ↔ 布尔值
- 字符串 ↔ 浮点数
- 接口类型转换

## 📝 完整示例

查看 [examples/](examples/) 目录获取更多使用示例：

- [基本配置示例](examples/main.go)
- [配置文件示例](examples/config/)
- [各种配置格式示例](examples/config/autoload/)

## 🤝 贡献指南

我们欢迎所有形式的贡献！请遵循以下步骤：

1. Fork 本仓库
2. 创建您的特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交您的更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开一个 Pull Request

### 开发环境设置

```bash
# 克隆仓库
git clone https://github.com/stones-hub/taurus-pro-config.git
cd taurus-pro-config

# 安装依赖
go mod download

# 运行测试
go test ./...

# 运行示例
go run examples/main.go
```

## 📄 许可证

本项目采用 Apache License 2.0 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 👥 作者

- **yelei** - *初始工作* - [61647649@qq.com](mailto:61647649@qq.com)

## 🙏 致谢

感谢以下开源项目的支持：

- [go-toml](https://github.com/pelletier/go-toml) - TOML 解析器
- [yaml.v3](https://github.com/go-yaml/yaml) - YAML 解析器
- [cast](https://github.com/spf13/cast) - 类型转换工具

## 📞 支持

如果您遇到问题或有建议，请：

1. 查看 [Issues](../../issues) 页面
2. 创建新的 Issue
3. 发送邮件至 [61647649@qq.com](mailto:61647649@qq.com)

---

**Taurus Pro Config** - 让配置管理变得简单而强大 🚀
