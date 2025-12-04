package models

import (
	"encoding/json"
	"os"
	"time"
)

// Config 应用配置模型
// 包含应用程序的所有配置信息
// 实现了JSON序列化和反序列化方法

// Config 配置结构体
type Config struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
	Logging  LoggingConfig  `json:"logging"`
	App      AppConfig      `json:"app"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host         string        `json:"host"`
	Port         int           `json:"port"`
	ReadTimeout  time.Duration `json:"read_timeout"`
	WriteTimeout time.Duration `json:"write_timeout"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"-"` // 密码不序列化
	DBName   string `json:"dbname"`
}

// LoggingConfig 日志配置
type LoggingConfig struct {
	Level  string `json:"level"`
	Format string `json:"format"`
	Path   string `json:"path"`
}

// AppConfig 应用配置
type AppConfig struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Env     string `json:"env"`
}

// MarshalJSON 自定义JSON序列化方法
func (c Config) MarshalJSON() ([]byte, error) {
	// 创建一个匿名结构体用于JSON序列化
	type Alias Config
	return json.Marshal(&struct {
		Alias
		// 自定义字段格式
		Server struct {
			ServerConfig
			ReadTimeout  string `json:"read_timeout"`
			WriteTimeout string `json:"write_timeout"`
		} `json:"server"`
	}{
		Alias: (Alias)(c),
		Server: struct {
			ServerConfig
			ReadTimeout  string `json:"read_timeout"`
			WriteTimeout string `json:"write_timeout"`
		}{
			ServerConfig: c.Server,
			ReadTimeout:  c.Server.ReadTimeout.String(),
			WriteTimeout: c.Server.WriteTimeout.String(),
		},
	})
}

// UnmarshalJSON 自定义JSON反序列化方法
func (c *Config) UnmarshalJSON(data []byte) error {
	// 创建一个匿名结构体用于JSON反序列化
	type Alias Config
	aux := &struct {
		*Alias
		Server struct {
			ServerConfig
			ReadTimeout  string `json:"read_timeout"`
			WriteTimeout string `json:"write_timeout"`
		} `json:"server"`
	}{
		Alias: (*Alias)(c),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// 解析时间字符串
	readTimeout, err := time.ParseDuration(aux.Server.ReadTimeout)
	if err != nil {
		return err
	}
	c.Server.ReadTimeout = readTimeout

	writeTimeout, err := time.ParseDuration(aux.Server.WriteTimeout)
	if err != nil {
		return err
	}
	c.Server.WriteTimeout = writeTimeout

	return nil
}

// LoadFromFile 从文件加载配置
func LoadConfigFromFile(filePath string) (Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	var config Config
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return Config{}, err
	}

	return config, nil
}

// SaveToFile 将配置保存到文件
func (c Config) SaveConfigToFile(filePath string) error {
	jsonData, err := c.MarshalJSON()
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, jsonData, 0644)
}