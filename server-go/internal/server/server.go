package server

import (
	"context"
	"net/http"
	"time"

	"vue2-blog-server/internal/auth"
	"vue2-blog-server/internal/handler"
	"vue2-blog-server/internal/middleware"
	"vue2-blog-server/internal/pkg"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Server struct {
	router     *gin.Engine
	logger     *zap.Logger
	jwtManager *auth.JWTManager
	handlers   *Handlers
}

type Handlers struct {
	AuthHandler    *handler.AuthHandler
	ArticleHandler *handler.ArticleHandler
	MessageHandler *handler.MessageHandler
	VisitorHandler *handler.VisitorHandler
	UploadHandler  *handler.UploadHandler
}

func NewServer(logger *zap.Logger, jwtManager *auth.JWTManager, handlers *Handlers) *Server {
	return &Server{
		logger:     logger,
		jwtManager: jwtManager,
		handlers:   handlers,
	}
}

func (s *Server) Initialize() {
	// 设置生产模式
	gin.SetMode(gin.ReleaseMode)
	
	s.router = gin.New()
	
	// 添加全局中间件
	s.router.Use(middleware.CORSMiddleware())
	s.router.Use(middleware.StructuredLoggerMiddleware(s.logger))
	s.router.Use(middleware.RecoveryMiddleware(s.logger))

	// 静态文件服务 - 添加头像目录映射
	s.router.Static("/img", "./public/img")
	s.router.Static("/public", "./public")
	s.router.Static("/avatar", "./public/avatar")

	// 健康检查
	s.router.GET("/healthz", s.healthCheck)

	// 注册路由
	s.registerRoutes()
}

func (s *Server) registerRoutes() {
	api := s.router.Group("/api/v1")
	
	// 认证相关路由
	auth := api.Group("/auth")
	{
		auth.POST("/register", s.handlers.AuthHandler.Register)
		auth.POST("/login", s.handlers.AuthHandler.Login)
		auth.POST("/logout", s.handlers.AuthHandler.Logout)
		auth.POST("/refresh", s.handlers.AuthHandler.RefreshToken)
		
		// 需要认证的路由
		authRequired := auth.Group("")
		authRequired.Use(middleware.AuthMiddleware(s.jwtManager))
		{
			authRequired.GET("/profile", s.handlers.AuthHandler.Profile)
		}
	}

	// 文章相关路由
	articles := api.Group("/articles")
	{
		articles.GET("", s.handlers.ArticleHandler.GetList)
		articles.GET("/hot", s.handlers.ArticleHandler.GetHot)
		articles.GET("/info", s.handlers.ArticleHandler.GetInfo)
		articles.POST("/search", s.handlers.ArticleHandler.Search)
		articles.GET("/:id", s.handlers.ArticleHandler.GetByID)
	}

	// 留言相关路由
	messages := api.Group("/messages")
	{
		messages.GET("", s.handlers.MessageHandler.GetList)
		
		// 需要认证的路由
		messageAuth := messages.Group("")
		messageAuth.Use(middleware.AuthMiddleware(s.jwtManager))
		{
			messageAuth.POST("", s.handlers.MessageHandler.Create)
			messageAuth.POST("/:id/replies", s.handlers.MessageHandler.Reply)
		}
	}

	// 访客相关路由
	visitors := api.Group("/visitors")
	{
		visitors.GET("", s.handlers.VisitorHandler.GetList)
	}

	// 上传相关路由
	upload := api.Group("/upload")
	upload.Use(middleware.AuthMiddleware(s.jwtManager))
	{
		upload.POST("/avatar", s.handlers.UploadHandler.UploadAvatar)
	}

	// 兼容旧版接口路由
	s.registerLegacyRoutes()
}

func (s *Server) registerLegacyRoutes() {
	// 兼容旧版接口
	legacy := s.router.Group("")
	
	// 文章接口兼容 - 修复请求方式为POST并支持form数据
	legacy.POST("/article/getInfo", s.wrapFormHandler(s.handlers.ArticleHandler.GetInfo))
	legacy.POST("/article/getHot", s.wrapFormHandler(s.handlers.ArticleHandler.GetHot))
	legacy.POST("/article/getShow", s.wrapFormHandler(s.handlers.ArticleHandler.GetList))
	legacy.POST("/article", s.wrapFormHandler(func(c *gin.Context) {
		// 处理单篇文章查询 - 通过_id参数
		id := c.PostForm("_id")
		if id != "" {
			c.Param("id")
			c.Set("id", id)
			s.handlers.ArticleHandler.GetByID(c)
		} else {
			pkg.ValidateError(c, "缺少文章ID参数")
		}
	}))
	legacy.POST("/article/search", s.wrapFormHandler(func(c *gin.Context) {
		// 处理搜索请求 - keywords -> content
		keywords := c.PostForm("keywords")
		if keywords != "" {
			c.Set("content", keywords)
			s.handlers.ArticleHandler.Search(c)
		} else {
			pkg.ValidateError(c, "缺少搜索关键词")
		}
	}))
	
	// 添加延伸阅读接口
	legacy.POST("/article/extend", s.wrapFormHandler(s.handlers.ArticleHandler.GetExtend))

	// 登录接口兼容
	login := legacy.Group("/login")
	{
		login.POST("", s.wrapFormHandler(s.handlers.AuthHandler.Login))
		login.POST("/logout", s.wrapFormHandler(s.handlers.AuthHandler.Logout))
		
		// 可选认证中间件
		login.Use(middleware.OptionalAuthMiddleware(s.jwtManager))
		login.POST("/ifLogin", s.wrapFormHandler(s.handlers.AuthHandler.CheckLoginLegacy))
	}

	// 注册接口兼容 - 添加验证码接口
	register := legacy.Group("/register")
	{
		register.POST("", s.wrapFormHandler(s.handlers.AuthHandler.Register))
		register.POST("/vcode", s.wrapFormHandler(s.handlers.AuthHandler.GenerateVCode))
		register.POST("/checkVcode", s.wrapFormHandler(s.handlers.AuthHandler.CheckVCode))
	}

	// 留言接口兼容 - 修复请求方式和参数映射
	message := legacy.Group("/message")
	{
		message.POST("/getList", s.wrapFormHandler(s.handlers.MessageHandler.GetList))
		
		messageAuth := message.Group("")
		messageAuth.Use(middleware.AuthMiddleware(s.jwtManager))
		{
			messageAuth.POST("/commit", s.wrapFormHandler(s.handlers.MessageHandler.Create))
			messageAuth.POST("/childCommit", s.wrapFormHandler(func(c *gin.Context) {
				// 修复参数映射 - parentId -> messageId
				parentId := c.PostForm("parentId")
				if parentId != "" {
					c.Param("id")
					c.Set("id", parentId)
					s.handlers.MessageHandler.Reply(c)
				} else {
					pkg.ValidateError(c, "缺少父留言ID")
				}
			}))
		}
	}

	// 访客接口兼容 - 修复请求方式
	legacy.POST("/visitor", s.wrapFormHandler(s.handlers.VisitorHandler.GetList))

	// 上传接口兼容
	uploadLegacy := legacy.Group("/upload")
	uploadLegacy.Use(middleware.AuthMiddleware(s.jwtManager))
	{
		uploadLegacy.POST("/avatar", s.handlers.UploadHandler.UploadAvatarLegacy)
	}
}

// wrapFormHandler 包装处理器以支持表单数据解析
func (s *Server) wrapFormHandler(handler gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 确保能够解析表单数据
		c.Request.ParseForm()
		handler(c)
	}
}

func (s *Server) healthCheck(c *gin.Context) {
	pkg.Success(c, gin.H{
		"status":    "ok",
		"timestamp": time.Now(),
		"service":   "vue2-blog-server",
	})
}

func (s *Server) Start(addr string) error {
	s.logger.Info("Starting server", zap.String("address", addr))
	
	srv := &http.Server{
		Addr:         addr,
		Handler:      s.router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	return srv.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("Shutting down server...")
	
	srv := &http.Server{
		Handler: s.router,
	}
	
	return srv.Shutdown(ctx)
}