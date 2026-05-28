package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"

	v1 "boread/internal/handler/v1"
	"boread/internal/middleware"
	"boread/internal/repository"
	"boread/internal/service"
)

// SetupRouter 装配所有路由 + 依赖注入
// 设计要点:
//   - 所有写操作受 button-level 鉴权保护
//   - 只读接口仅要求登录, 不做 button 校验
//   - /api/auth/* 是公开或登录态接口, 不参与 manage 模块的按钮鉴权
func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.New()
	r.Use(middleware.Cors())
	r.Use(middleware.RequestLogger())
	r.Use(gin.Recovery())

	// 健康检查
	healthHandler := v1.NewHealthHandler()
	r.GET("/ping", healthHandler.Ping)
	r.NoRoute(healthHandler.NoRoute)
	r.NoMethod(healthHandler.NoMethod)

	// Swagger UI
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// === Repository 层 ===
	userRepo := repository.NewSysUserRepository(db)
	deptRepo := repository.NewSysDeptRepository(db)
	roleRepo := repository.NewSysRoleRepository(db)
	menuRepo := repository.NewSysMenuRepository(db)
	dictRepo := repository.NewSysDictRepository(db)
	logRepo := repository.NewSysLogRepository(db)
	bookCategoryRepo := repository.NewBookCategoryRepository(db)
	bookTagRepo := repository.NewBookTagRepository(db)

	// === Service 层 ===
	authSvc := service.NewAuthService(userRepo, db)
	deptSvc := service.NewDeptService(deptRepo)
	roleSvc := service.NewRoleService(roleRepo)
	userSvc := service.NewUserService(userRepo)
	menuSvc := service.NewMenuService(menuRepo)
	dictSvc := service.NewDictService(dictRepo)
	logSvc := service.NewLogService(logRepo)
	bookCategorySvc := service.NewBookCategoryService(bookCategoryRepo)
	bookTagSvc := service.NewBookTagService(bookTagRepo)

	// === Handler 层 ===
	authHandler := v1.NewAuthHandler(authSvc)
	deptHandler := v1.NewDeptHandler(deptSvc)
	roleHandler := v1.NewRoleHandler(roleSvc)
	userHandler := v1.NewUserHandler(userSvc)
	menuHandler := v1.NewMenuHandler(menuSvc)
	dictHandler := v1.NewDictHandler(dictSvc)
	logHandler := v1.NewLogHandler(logSvc)
	bookCategoryHandler := v1.NewBookCategoryHandler(bookCategorySvc)
	bookTagHandler := v1.NewBookTagHandler(bookTagSvc)

	api := r.Group("/api")
	{
		// 公开
		api.POST("/auth/login", authHandler.Login)

		// 登录态
		authed := api.Group("")
		authed.Use(middleware.Auth())
		{
			authed.GET("/auth/userInfo", authHandler.GetUserInfo)
			authed.GET("/auth/menu", authHandler.GetUserMenu)
			authed.GET("/auth/buttons", authHandler.GetButtons)
		}

		// 受保护管理接口
		manage := api.Group("/manage")
		manage.Use(middleware.Auth())
		{
			// --- Dept ---
			manage.GET("/dept/tree", deptHandler.Tree)
			manage.GET("/dept/:id", deptHandler.GetByID)
			manage.POST("/dept/page", deptHandler.Page)
			manage.POST("/dept", middleware.RequireButton(authSvc, "dept:create"), deptHandler.Create)
			manage.PUT("/dept/:id", middleware.RequireButton(authSvc, "dept:update"), deptHandler.Update)
			manage.DELETE("/dept/:id", middleware.RequireButton(authSvc, "dept:delete"), deptHandler.Delete)

			// --- Role ---
			manage.GET("/role/all", roleHandler.AllBrief)
			manage.GET("/role/:id", roleHandler.GetByID)
			manage.POST("/role/page", roleHandler.Page)
			manage.GET("/role/:id/menus", roleHandler.GetMenuIDs)
			manage.GET("/role/:id/buttons", roleHandler.GetButtonIDs)
			manage.POST("/role", middleware.RequireButton(authSvc, "role:create"), roleHandler.Create)
			manage.PUT("/role/:id", middleware.RequireButton(authSvc, "role:update"), roleHandler.Update)
			manage.DELETE("/role/:id", middleware.RequireButton(authSvc, "role:delete"), roleHandler.Delete)
			manage.PUT("/role/:id/menus", middleware.RequireButton(authSvc, "role:grant"), roleHandler.GrantMenus)
			manage.PUT("/role/:id/buttons", middleware.RequireButton(authSvc, "role:grant"), roleHandler.GrantButtons)

			// --- User ---
			manage.GET("/user/:id", userHandler.GetByID)
			manage.POST("/user/page", userHandler.Page)
			manage.POST("/user", middleware.RequireButton(authSvc, "user:create"), userHandler.Create)
			manage.PUT("/user/:id", middleware.RequireButton(authSvc, "user:update"), userHandler.Update)
			manage.DELETE("/user/:id", middleware.RequireButton(authSvc, "user:delete"), userHandler.Delete)
			manage.PUT("/user/:id/reset-password", middleware.RequireButton(authSvc, "user:reset_pwd"), userHandler.ResetPassword)

			// --- Menu ---
			manage.GET("/menu/tree", menuHandler.Tree)
			manage.GET("/menu/:id", menuHandler.GetByID)
			manage.POST("/menu/page", menuHandler.Page) // 菜单分页列表
			manage.POST("/menu", middleware.RequireButton(authSvc, "menu:create"), menuHandler.Create)
			manage.PUT("/menu/:id", middleware.RequireButton(authSvc, "menu:update"), menuHandler.Update)
			manage.DELETE("/menu/:id", middleware.RequireButton(authSvc, "menu:delete"), menuHandler.Delete)
			manage.GET("/menu/buttons/:menuId", menuHandler.ListButtonsByMenu)
			manage.POST("/menu/button", middleware.RequireButton(authSvc, "menu:create"), menuHandler.CreateButton)
			manage.DELETE("/menu/button/:id", middleware.RequireButton(authSvc, "menu:delete"), menuHandler.DeleteButton)

			// --- Dict ---
			manage.GET("/dict/:id", dictHandler.GetByID)
			manage.POST("/dict/page", dictHandler.Page)
			manage.POST("/dict", middleware.RequireButton(authSvc, "dict:create"), dictHandler.Create)
			manage.PUT("/dict/:id", middleware.RequireButton(authSvc, "dict:update"), dictHandler.Update)
			manage.DELETE("/dict/:id", middleware.RequireButton(authSvc, "dict:delete"), dictHandler.Delete)
			manage.GET("/dict/items/:dictId", dictHandler.ItemsByDictID)
			manage.GET("/dict/code/:code", dictHandler.ItemsByCode) // 前端高频接口
			manage.POST("/dict/item", middleware.RequireButton(authSvc, "dict:create"), dictHandler.CreateItem)
			manage.PUT("/dict/item/:id", middleware.RequireButton(authSvc, "dict:update"), dictHandler.UpdateItem)
			manage.DELETE("/dict/item/:id", middleware.RequireButton(authSvc, "dict:delete"), dictHandler.DeleteItem)

			// --- Log ---
			manage.POST("/log/login/page", logHandler.PageLogin)
			manage.POST("/log/operation/page", logHandler.PageOperation)

			// --- Book Category ---
			manage.GET("/book-category/tree", bookCategoryHandler.Tree)
			manage.GET("/book-category/:id", bookCategoryHandler.GetByID)
			manage.POST("/book-category/page", bookCategoryHandler.Page)
			manage.POST("/book-category", middleware.RequireButton(authSvc, "book-category:create"), bookCategoryHandler.Create)
			manage.PUT("/book-category/:id", middleware.RequireButton(authSvc, "book-category:update"), bookCategoryHandler.Update)
			manage.DELETE("/book-category/:id", middleware.RequireButton(authSvc, "book-category:delete"), bookCategoryHandler.Delete)

			// --- Book Tag ---
			manage.GET("/book-tag/:id", bookTagHandler.GetByID)
			manage.POST("/book-tag/page", bookTagHandler.Page)
			manage.POST("/book-tag", middleware.RequireButton(authSvc, "book-tag:create"), bookTagHandler.Create)
			manage.PUT("/book-tag/:id", middleware.RequireButton(authSvc, "book-tag:update"), bookTagHandler.Update)
			manage.DELETE("/book-tag/:id", middleware.RequireButton(authSvc, "book-tag:delete"), bookTagHandler.Delete)
		}
	}

	return r
}
