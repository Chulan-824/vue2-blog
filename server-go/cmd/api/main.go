package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"vue2-blog-server/internal/auth"
	"vue2-blog-server/internal/config"
	"vue2-blog-server/internal/handler"
	"vue2-blog-server/internal/repository"
	"vue2-blog-server/internal/server"
	"vue2-blog-server/internal/service"
	"vue2-blog-server/pkg/logger"

	"go.uber.org/zap"
)

func main() {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化日志
	zapLogger, err := logger.New(&cfg.Log)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer zapLogger.Sync()

	zapLogger.Info("Starting vue2-blog-server...")

	// 初始化数据库
	db, err := repository.NewDatabase(cfg.Database.URI, cfg.Database.DBName)
	if err != nil {
		zapLogger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	zapLogger.Info("Connected to database successfully")

	// 初始化Repository层
	userRepo := repository.NewUserRepository(db)
	articleRepo := repository.NewArticleRepository(db)
	messageRepo := repository.NewMessageRepository(db)
	visitorRepo := repository.NewVisitorRepository(db)

	// 初始化Service层
	userService := service.NewUserService(userRepo, visitorRepo)
	articleService := service.NewArticleService(articleRepo)
	messageService := service.NewMessageService(messageRepo, userRepo)
	visitorService := service.NewVisitorService(visitorRepo, userRepo)
	uploadService := service.NewUploadService(userService, cfg.Upload.Path, cfg.Upload.BaseURL)

	// 初始化JWT管理器
	jwtManager := auth.NewJWTManager(cfg.JWT.Secret, cfg.JWT.TokenExpires, cfg.JWT.RefreshExpires)

	// 初始化Handler层
	handlers := &server.Handlers{
		AuthHandler:    handler.NewAuthHandler(userService, jwtManager),
		ArticleHandler: handler.NewArticleHandler(articleService),
		MessageHandler: handler.NewMessageHandler(messageService),
		VisitorHandler: handler.NewVisitorHandler(visitorService),
		UploadHandler:  handler.NewUploadHandler(uploadService),
	}

	// 初始化服务器
	srv := server.NewServer(zapLogger, jwtManager, handlers)
	srv.Initialize()

	// 启动服务器
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	
	// 在单独的goroutine中启动服务器
	go func() {
		if err := srv.Start(addr); err != nil {
			zapLogger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	zapLogger.Info("Server started successfully", zap.String("address", addr))

	// 等待中断信号以优雅地关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	zapLogger.Info("Shutting down server...")

	// 创建一个5秒的上下文来关闭服务器
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		zapLogger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	zapLogger.Info("Server exited")
}