package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	_ "boread/docs"
	v1 "boread/internal/handler/v1"
	"boread/internal/middleware"
	"boread/internal/router"
	"boread/internal/seed"
	"boread/pkg/appsignal"
	"boread/pkg/config"
	jwtPkg "boread/pkg/jwt"
	"boread/pkg/logger"
)

// @title           Boread API
// @version         1.0
// @description     小说阅读平台后端 API
// @host            localhost:8080
// @BasePath        /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	seedFlag := flag.Bool("seed", false, "运行种子初始化后退出 (idempotent)")
	portFlag := flag.Int("port", 0, "覆盖配置文件的 server.port (0=使用配置值)")
	flag.Parse()

	// 初始化 JWT 默认值（仅作为配置不可用时的兜底）
	jwtPkg.Init("boread-secret", 7200)

	// 尝试加载配置文件
	cfgExists := config.Exists("configs/config.yaml")
	if cfgExists {
		if err := config.Load("configs/config.yaml"); err != nil {
			log.Printf("WARN: 配置文件读取失败, 将使用默认值: %v", err)
			config.Cfg = nil
			cfgExists = false
		}
	}

	// 配置文件中的 JWT 设置覆盖默认值
	if cfgExists && config.Cfg != nil {
		if config.Cfg.JWT.Secret != "" {
			expire := config.Cfg.JWT.Expire
			if expire <= 0 {
				expire = 7200
			}
			jwtPkg.Init(config.Cfg.JWT.Secret, expire)
		}
	}

	// ========== 有数据库配置 → 全量启动 ==========
	if cfgExists && config.Cfg != nil && config.Cfg.Database.Host != "" {
		startFullServer(*seedFlag, *portFlag)
		return
	}

	// ========== 无数据库配置 → Setup 模式（支持自动重启） ==========
	runSetupMode(*portFlag)
}

// startFullServer 以完整模式启动：连接数据库 → 装配路由 → 启动 HTTP
func startFullServer(seedFlag bool, portFlag int) {
	cfg := config.Cfg
	logger.Init(cfg.Log.Level, cfg.Log.File)
	defer logger.Log.Sync()

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Warn),
	})
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get sql.DB: %v", err)
	}
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)

	// 种子模式: 仅初始化数据后退出
	if seedFlag {
		if err := seed.Run(context.Background(), db); err != nil {
			log.Fatalf("Seed failed: %v", err)
		}
		fmt.Println("Seed completed.")
		fmt.Println("Default admin -> username: admin  password: admin123")
		return
	}

	r := router.SetupRouter(db, cfg.CORS.AllowedOrigins...)

	port := cfg.Server.Port
	if portFlag > 0 {
		port = portFlag
	}
	addr := fmt.Sprintf(":%d", port)
	logger.Log.Info(fmt.Sprintf("Server starting on %s", addr))

	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// runSetupMode 以 Setup 模式启动：仅注册健康检查 + DB 配置路由
// 配置保存后通过 appsignal 触发自动重启进入完整模式
func runSetupMode(portFlag int) {
	logger.Init("info", "logs/boread.log")
	defer logger.Log.Sync()

	// Setup 模式也读取 CORS 配置（如果有的话）
	var corsOrigins []string
	if config.Cfg != nil {
		corsOrigins = config.Cfg.CORS.AllowedOrigins
	}

	r := gin.New()
	r.Use(middleware.Cors(corsOrigins...))
	r.Use(middleware.RequestLogger())
	r.Use(gin.Recovery())

	// 健康检查
	healthHandler := v1.NewHealthHandler()
	r.GET("/ping", healthHandler.Ping)
	r.NoRoute(healthHandler.NoRoute)
	r.NoMethod(healthHandler.NoMethod)

	// Setup 路由（无 DB 依赖）
	setupHandler := v1.NewSetupHandler()
	r.GET("/api/setup/status", setupHandler.Status)
	r.POST("/api/setup/database", setupHandler.SaveConfig)

	port := 8080
	if portFlag > 0 {
		port = portFlag
	}
	addr := fmt.Sprintf(":%d", port)
	logger.Log.Info(fmt.Sprintf("Server starting in SETUP mode on %s", addr))
	logger.Log.Info("Database not configured. Use /api/setup/database to configure.")

	// 使用可控 http.Server 以支持配置保存后自动重启
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Setup server error: %v", err)
		}
	}()

	// 阻塞等待重启信号（由 SaveConfig 在保存成功后发送）
	appsignal.WaitRestart()

	logger.Log.Info("Configuration saved, restarting server in full mode...")

	// 优雅关闭 setup 服务器
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Log.Warn(fmt.Sprintf("Setup server shutdown: %v", err))
	}

	// 重新加载配置（SaveConfig 已写入 configs/config.yaml）
	if err := config.Load("configs/config.yaml"); err != nil {
		log.Fatalf("Failed to reload config after save: %v", err)
	}

	// 重启后用配置中的 JWT 设置
	if config.Cfg.JWT.Secret != "" {
		expire := config.Cfg.JWT.Expire
		if expire <= 0 {
			expire = 7200
		}
		jwtPkg.Init(config.Cfg.JWT.Secret, expire)
	}

	// 启动完整模式（不执行 seed）
	startFullServer(false, portFlag)
}
