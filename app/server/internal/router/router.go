package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"

	_ "boread/docs"
	v1 "boread/internal/handler/v1"
	"boread/internal/middleware"
	"boread/internal/repository"
	"boread/internal/service"
)

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.New()

	r.Use(middleware.Cors())
	r.Use(middleware.RequestLogger())
	r.Use(gin.Recovery())

	healthHandler := v1.NewHealthHandler()
	r.GET("/ping", healthHandler.Ping)
	r.NoRoute(healthHandler.NoRoute)
	r.NoMethod(healthHandler.NoMethod)

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := v1.NewUserHandler(userService)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api/v1")
	{
		auth := api.Group("")
		auth.Use(middleware.Auth())
		{
			auth.GET("/user/profile", userHandler.GetProfile)
			auth.PUT("/user/profile", userHandler.UpdateProfile)
			auth.DELETE("/user/:id", userHandler.Delete)
			auth.GET("/user/:id", userHandler.GetByID)
		}

		api.POST("/user/register", userHandler.Register)
		api.POST("/user/login", userHandler.Login)
		api.GET("/users", userHandler.List)
	}

	return r
}