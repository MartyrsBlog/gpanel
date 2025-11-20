package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"gpanel/internal/api"
	"gpanel/internal/config"
	"gpanel/internal/middleware"
	"gpanel/internal/plugin"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化数据库
	if err := config.InitDB(cfg); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 创建 Gin 引擎
	r := gin.Default()

	// 全局中间件
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.ErrorHandler())

	// CORS 配置
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// 静态文件服务 (前端)
	r.Use(static.Serve("/", static.LocalFile("./web/dist", false)))

	// API 路由组
	apiGroup := r.Group("/api")
	{
		// 认证路由
		authGroup := apiGroup.Group("/auth")
		{
			authGroup.POST("/login", api.Login)
			authGroup.POST("/logout", api.Logout)
			authGroup.GET("/info", middleware.JWTAuth(), api.GetUserInfo)
		}

		// 系统监控路由
		systemGroup := apiGroup.Group("/system")
		systemGroup.Use(middleware.JWTAuth())
		{
			systemGroup.GET("/monitor", api.GetSystemMonitor)
			systemGroup.GET("/processes", api.GetProcessList)
			systemGroup.GET("/disk", api.GetDiskInfo)
		}

		// 网站管理路由
		siteGroup := apiGroup.Group("/site")
		siteGroup.Use(middleware.JWTAuth())
		{
			siteGroup.GET("/list", api.GetSiteList)
			siteGroup.POST("/create", api.CreateSite)
			siteGroup.POST("/delete", api.DeleteSite)
			siteGroup.POST("/ssl", api.ApplySSL)
		}

		// 文件管理路由
		fileGroup := apiGroup.Group("/file")
		fileGroup.Use(middleware.JWTAuth())
		{
			fileGroup.GET("/list", api.GetFileList)
			fileGroup.POST("/upload", api.UploadFile)
			fileGroup.GET("/download", api.DownloadFile)
			fileGroup.POST("/extract", api.ExtractFile)
			fileGroup.POST("/chmod", api.ChmodFile)
			fileGroup.POST("/rename", api.RenameFile)
			fileGroup.POST("/delete", api.DeleteFile)
		}

		// 数据库管理路由
		dbGroup := apiGroup.Group("/database")
		dbGroup.Use(middleware.JWTAuth())
		{
			dbGroup.GET("/list", api.GetDatabaseList)
			dbGroup.POST("/create", api.CreateDatabase)
			dbGroup.POST("/delete", api.DeleteDatabase)
			dbGroup.GET("/backup", api.BackupDatabase)
		}

		// Docker 管理路由
		dockerGroup := apiGroup.Group("/docker")
		dockerGroup.Use(middleware.JWTAuth())
		{
			dockerGroup.GET("/containers", api.GetContainers)
			dockerGroup.POST("/container/start", api.StartContainer)
			dockerGroup.POST("/container/stop", api.StopContainer)
			dockerGroup.GET("/container/logs", api.GetContainerLogs)
			dockerGroup.GET("/images", api.GetImages)
			dockerGroup.POST("/image/pull", api.PullImage)
		}
	}

	// 插件系统
	pluginLoader := plugin.NewLoader()
	if err := pluginLoader.LoadPlugins("plugins", apiGroup); err != nil {
		log.Printf("Failed to load plugins: %v", err)
	}

	// Swagger 文档路由
	r.GET("/swagger/*any", gin.WrapH(http.StripPrefix("/swagger/", http.FileServer(http.Dir("./docs/swagger")))))

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"message": "GPanel is running",
		})
	})

	port := cfg.Server.Port
	if port == "" {
		port = "8080"
	}

	fmt.Printf("🚀 GPanel server starting on port %s\n", port)
	fmt.Printf("📖 Swagger documentation: http://localhost:%s/swagger/index.html\n", port)
	fmt.Printf("🌐 Web interface: http://localhost:%s\n", port)

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
		os.Exit(1)
	}
}