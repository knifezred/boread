package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	_ "boread/docs"
	"boread/internal/router"
	"boread/internal/seed"
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
	flag.Parse()

	if err := config.Load("configs/config.yaml"); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	cfg := config.Cfg
	logger.Init(cfg.Log.Level, cfg.Log.File)
	defer logger.Log.Sync()

	jwtPkg.Init(cfg.JWT.Secret, cfg.JWT.Expire)

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
	if *seedFlag {
		if err := seed.Run(context.Background(), db); err != nil {
			log.Fatalf("Seed failed: %v", err)
		}
		fmt.Println("Seed completed.")
		fmt.Println("Default admin -> username: admin  password: admin123")
		return
	}

	r := router.SetupRouter(db)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	logger.Log.Info(fmt.Sprintf("Server starting on %s", addr))

	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
