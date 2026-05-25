package main

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"boread/internal/router"
	"boread/pkg/config"
	jwtPkg "boread/pkg/jwt"
	"boread/pkg/logger"
)

func main() {
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

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get sql.DB: %v", err)
	}
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)

	r := router.SetupRouter(db)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	logger.Log.Info(fmt.Sprintf("Server starting on %s", addr))

	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}