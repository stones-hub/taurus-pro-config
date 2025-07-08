// Copyright (c) 2025 Taurus Team. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Author: yelei
// Email: 61647649@qq.com
// Date: 2025-06-13

package config

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cast"
	"gopkg.in/yaml.v3"
)

// Config 表示配置结构
type Config struct {
	data        map[string]interface{}
	PrintEnable bool
}

type ConfigOption func(*Config)

func WithPrintEnable(enable bool) ConfigOption {
	return func(c *Config) {
		c.PrintEnable = enable
	}
}

// New 创建一个新的配置实例
func New(options ...ConfigOption) *Config {
	c := &Config{
		data:        make(map[string]interface{}),
		PrintEnable: false,
	}
	for _, option := range options {
		option(c)
	}
	return c
}

// Initialize 初始化配置
func (c *Config) Initialize(configPath string, env string) error {
	// 如果提供了env文件路径，加载环境变量
	if env != "" {
		if err := c.loadEnv(env); err != nil {
			log.Printf("Error loading .env file: %v\n", err)
		}
	}

	// 加载应用配置文件
	log.Printf("Loading application configuration file: %s", configPath)
	if err := c.loadConfig(configPath); err != nil {
		return err
	}

	// 打印应用配置
	if c.PrintEnable {
		log.Printf("Configuration: %s", c.ToJSONString())
	}

	return nil
}

// loadEnv 加载环境变量文件
func (c *Config) loadEnv(envPath string) error {
	data, err := os.ReadFile(envPath)
	if err != nil {
		return fmt.Errorf("failed to read env file: %w", err)
	}

	lines := string(data)
	re := regexp.MustCompile(`(\w+)=(.+)`)
	matches := re.FindAllStringSubmatch(lines, -1)
	for _, match := range matches {
		if len(match) == 3 {
			os.Setenv(match[1], match[2])
		}
	}
	return nil
}

// loadConfig 从目录或单个文件加载配置
func (c *Config) loadConfig(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("failed to access config path: %w", err)
	}

	if info.IsDir() {
		// 递归遍历目录下的所有文件
		err := filepath.Walk(path, func(filePath string, fileInfo os.FileInfo, err error) error {
			if err != nil {
				log.Printf("Error accessing file %s: %v\n", filePath, err)
				return nil
			}

			// 跳过目录
			if fileInfo.IsDir() {
				return nil
			}

			// 加载配置文件
			if err := c.loadConfigFile(filePath); err != nil {
				log.Printf("Error loading config file %s: %v\n", filePath, err)
			}
			return nil
		})

		if err != nil {
			return fmt.Errorf("failed to walk through config directory: %w", err)
		}
	} else {
		// 加载单个配置文件
		return c.loadConfigFile(path)
	}

	log.Println("Configuration loaded successfully")
	return nil
}

// loadConfigFile 加载单个配置文件
func (c *Config) loadConfigFile(filePath string) error {
	ext := filepath.Ext(filePath)
	data, err := os.ReadFile(filePath)

	if err != nil {
		log.Printf("Failed to open config file: %v\n", err)
		return err
	}

	// 替换环境变量占位符
	content := c.replacePlaceholders(string(data))

	var config map[string]interface{}
	switch ext {
	case ".json":
		err = json.Unmarshal([]byte(content), &config)
		if err != nil {
			log.Printf("Failed to parse JSON config file: %s; error: %v\n", filePath, err)
			return err
		}
	case ".yaml", ".yml":
		err = yaml.Unmarshal([]byte(content), &config)
		if err != nil {
			log.Printf("Failed to parse YAML config file: %s; error: %v\n", filePath, err)
			return err
		}
	case ".toml":
		err = toml.Unmarshal([]byte(content), &config)
		if err != nil {
			log.Printf("Failed to parse TOML config file: %s; error: %v\n", filePath, err)
			return err
		}
	case ".xml":
		// XML需要特殊处理，因为它的结构可能更复杂
		var xmlData interface{}
		err = xml.Unmarshal([]byte(content), &xmlData)
		if err != nil {
			log.Printf("Failed to parse XML config file: %s; error: %v\n", filePath, err)
			return err
		}
		// 将XML数据转换为map
		config = make(map[string]interface{})
		if m, ok := xmlData.(map[string]interface{}); ok {
			config = m
		} else {
			config["root"] = xmlData
		}
	default:
		log.Printf("Unsupported config file format: %s\n", filePath)
		return fmt.Errorf("unsupported config file format: %s", filePath)
	}

	c.MergeMap(config)
	return nil
}

// replacePlaceholders 替换配置内容中的环境变量占位符
func (c *Config) replacePlaceholders(content string) string {
	re := regexp.MustCompile(`\$\{(\w+):([^}]+)\}`)
	return re.ReplaceAllStringFunc(content, func(match string) string {
		parts := re.FindStringSubmatch(match)
		if len(parts) == 3 {
			envVar := parts[1]
			defaultValue := parts[2]
			if value, exists := os.LookupEnv(envVar); exists {
				return value
			}
			return defaultValue
		}
		return match
	})
}

// MergeMap 合并配置映射
func (c *Config) MergeMap(data map[string]interface{}) {
	for k, v := range data {
		if existing, ok := c.data[k]; ok {
			c.data[k] = mergeValues(existing, v)
		} else {
			c.data[k] = v
		}
	}
}

// mergeValues 合并两个值，支持嵌套的 map 结构
func mergeValues(existing, new interface{}) interface{} {
	existingMap, existingOk := existing.(map[string]interface{})
	newMap, newOk := new.(map[string]interface{})

	// 如果两个值都是 map，则递归合并
	if existingOk && newOk {
		result := make(map[string]interface{})
		// 复制现有的值
		for k, v := range existingMap {
			result[k] = v
		}
		// 合并新的值
		for k, v := range newMap {
			if existing, ok := result[k]; ok {
				result[k] = mergeValues(existing, v)
			} else {
				result[k] = v
			}
		}
		return result
	}

	// 如果不是 map，则使用新值覆盖旧值
	return new
}

// Get 获取指定键的值
func (c *Config) Get(key string) interface{} {
	keys := splitKey(key)
	current := c.data

	for i, k := range keys {
		if i == len(keys)-1 {
			return current[k]
		}
		if v, ok := current[k]; ok {
			if m, ok := v.(map[string]interface{}); ok {
				current = m
			} else {
				return nil
			}
		} else {
			return nil
		}
	}
	return nil
}

// splitKey 将点分隔的键拆分为切片
func splitKey(key string) []string {
	re := regexp.MustCompile(`\.`)
	return re.Split(key, -1)
}

// GetString 获取字符串值
func (c *Config) GetString(key string) string {
	return cast.ToString(c.Get(key))
}

// GetInt 获取整数值
func (c *Config) GetInt(key string) int {
	return cast.ToInt(c.Get(key))
}

// GetBool 获取布尔值
func (c *Config) GetBool(key string) bool {
	return cast.ToBool(c.Get(key))
}

// GetFloat64 获取浮点数值
func (c *Config) GetFloat64(key string) float64 {
	return cast.ToFloat64(c.Get(key))
}

// GetStringMap 获取字符串map
func (c *Config) GetStringMap(key string) map[string]interface{} {
	return cast.ToStringMap(c.Get(key))
}

// GetStringSlice 获取字符串切片
func (c *Config) GetStringSlice(key string) []string {
	return cast.ToStringSlice(c.Get(key))
}

// ToJSONString 将配置转换为JSON字符串
func (c *Config) ToJSONString() string {
	bytes, err := json.Marshal(c.data)
	if err != nil {
		return fmt.Sprintf("Error marshaling config to JSON: %v", err)
	}
	return string(bytes)
}
