package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/example/message_processor/api"
	"github.com/example/message_processor/middleware"
	"github.com/example/message_processor/models"
	"github.com/example/message_processor/storage"
)

func main() {
	// 解析命令行参数
	configFile := flag.String("config", "config.json", "Path to configuration file")
	flag.Parse()

	// 加载配置
	config, err := loadConfig(*configFile)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// 初始化数据库
	dbConfig := storage.DBConfig{
		Host:     config.Database.Host,
		Port:     config.Database.Port,
		User:     config.Database.User,
		Password: config.Database.Password,
		DBName:   config.Database.DBName,
		SSLMode:  "disable",
	}

	db := storage.NewPostgresDB(dbConfig)
	if err := db.Connect(context.Background()); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Disconnect(context.Background())

	// 初始化消息处理器
	messageProcessor := &api.DefaultMessageProcessor{}

	// 初始化API处理器
	handler := api.NewHandler(messageProcessor)

	// 初始化认证中间件
	authMiddleware := middleware.NewAuthMiddleware(config.App.JWTSecret, "API_")

	// 设置路由
	mux := setupRouter(handler, authMiddleware)

	// 创建服务器
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port),
		Handler:      mux,
		ReadTimeout:  config.Server.ReadTimeout,
		WriteTimeout: config.Server.WriteTimeout,
	}

	// 启动服务器（异步）
	go func() {
		log.Printf("Server starting on %s:%d", config.Server.Host, config.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// 等待中断信号以优雅地关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// 给服务器5秒时间完成当前请求
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}

// setupRouter 设置HTTP路由
func setupRouter(handler *api.Handler, authMiddleware *middleware.AuthMiddleware) *http.ServeMux {
	mux := http.NewServeMux()

	// 应用全局中间件
	mux.HandleFunc("/health", handler.HealthCheck)

	// 公开API（不需要认证）
	public := http.NewServeMux()
	public.HandleFunc("/api/v1/message", handler.ProcessMessageHandler)

	// 需要认证的API
	protected := http.NewServeMux()
	protected.HandleFunc("/api/v1/resource", handler.GetResourceHandler)

	// 应用认证中间件
	mux.Handle("/api/v1/message", authMiddleware.APIKeyAuth(public))
	mux.Handle("/api/v1/resource", authMiddleware.JWTAuth(protected))

	return mux
}

// loadConfig 加载配置
// 注意：这里简化了配置加载，实际项目中应该从文件或环境变量加载
func loadConfig(configFile string) (*AppConfig, error) {
	// 这里使用硬编码配置，实际项目中应该从文件加载
	return &AppConfig{
		Server: ServerConfig{
			Host:         "0.0.0.0",
			Port:         8080,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "postgres",
			Password: "postgres",
			DBName:   "message_processor",
		},
		App: AppInfoConfig{
			Name:        "message_processor",
			Version:     "1.0.0",
			Environment: "development",
			JWTSecret:   "your-secret-key",
		},
	}, nil
}

// AppConfig 应用配置
type AppConfig struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
	App      AppInfoConfig  `json:"app"`
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
	Password string `json:"password"`
	DBName   string `json:"dbname"`
}

// AppInfoConfig 应用信息配置
type AppInfoConfig struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Environment string `json:"environment"`
	JWTSecret   string `json:"jwt_secret"`
}