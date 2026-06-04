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
	bookRepo := repository.NewBookRepository(db)
	bookTagRelRepo := repository.NewBookTagRelRepository(db)
	bookFileRepo := repository.NewBookFileRepository(db)
	bookUploadRepo := repository.NewBookUploadRepository(db)
	bookChapterRepo := repository.NewBookChapterRepository(db)
	bookChapterRuleRepo := repository.NewBookChapterRuleRepository(db)
	bookChapterRuleRelRepo := repository.NewBookChapterRuleRelRepository(db)
	bookFilterRuleRepo := repository.NewBookContentFilterRuleRepository(db)
	readerBookshelfRepo := repository.NewReaderBookshelfRepository(db)
	readerReadProgressRepo := repository.NewReaderReadProgressRepository(db)
	readerNoteRepo := repository.NewReaderNoteRepository(db)
	bookReviewRepo := repository.NewBookReviewRepository(db)
	bookCommentRepo := repository.NewBookChapterCommentRepository(db)
	readerLikeRepo := repository.NewReaderLikeRepository(db)

	// === Service 层 ===
	authSvc := service.NewAuthService(userRepo, bookChapterRuleRepo, db)
	deptSvc := service.NewDeptService(deptRepo)
	roleSvc := service.NewRoleService(roleRepo)
	userSvc := service.NewUserService(userRepo)
	menuSvc := service.NewMenuService(menuRepo)
	dictSvc := service.NewDictService(dictRepo)
	logSvc := service.NewLogService(logRepo)
	bookCategorySvc := service.NewBookCategoryService(bookCategoryRepo)
	bookTagSvc := service.NewBookTagService(bookTagRepo)
	bookSvc := service.NewBookService(db, bookRepo, bookTagRelRepo, bookCategoryRepo, bookTagRepo, bookChapterRepo)
	bookFileSvc := service.NewBookFileService(db, bookRepo, bookFileRepo, bookUploadRepo, bookChapterRepo, bookChapterRuleRepo, bookChapterRuleRelRepo, bookFilterRuleRepo, bookCategoryRepo, bookTagRepo)
	readerSvc := service.NewReaderService(db, readerBookshelfRepo, readerReadProgressRepo, bookRepo, bookChapterRepo)
	bookSocialSvc := service.NewBookSocialService(db, readerNoteRepo, bookReviewRepo, bookCommentRepo, readerLikeRepo, bookRepo, bookChapterRepo, userRepo)

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
	bookHandler := v1.NewBookHandler(bookSvc)
	bookFileHandler := v1.NewBookFileHandler(bookFileSvc)
	readerHandler := v1.NewReaderHandler(readerSvc)
	noteHandler := v1.NewNoteHandler(bookSocialSvc)
	reviewHandler := v1.NewReviewHandler(bookSocialSvc)
	commentHandler := v1.NewCommentHandler(bookSocialSvc)
	likeHandler := v1.NewLikeHandler(bookSocialSvc)

	api := r.Group("/api")
	{
		// 公开
		api.POST("/auth/login", authHandler.Login)
		api.GET("/book-category/hot", bookCategoryHandler.HotList)

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

			// --- Book ---
			manage.GET("/book/:id", bookHandler.GetByID)
			manage.POST("/book/page", bookHandler.Page)
			manage.POST("/book", middleware.RequireButton(authSvc, "book:create"), bookHandler.Create)
			manage.PUT("/book/:id", middleware.RequireButton(authSvc, "book:update"), bookHandler.Update)
			manage.PUT("/book/:id/status", middleware.RequireButton(authSvc, "book:update"), bookHandler.UpdateStatus)
			manage.DELETE("/book/:id", middleware.RequireButton(authSvc, "book:delete"), bookHandler.Delete)

			// --- Book File ---
			manage.POST("/book/upload", middleware.RequireButton(authSvc, "book:create"), bookFileHandler.Upload)
			manage.POST("/book/confirm-import", middleware.RequireButton(authSvc, "book:create"), bookFileHandler.ConfirmImport)
			manage.POST("/book/scan", middleware.RequireButton(authSvc, "book:create"), bookFileHandler.ScanAll)
			manage.POST("/book/scan-path", middleware.RequireButton(authSvc, "book:create"), bookFileHandler.ScanPath)
			manage.POST("/book/scan/:id", middleware.RequireButton(authSvc, "book:create"), bookFileHandler.ScanByID)
			manage.GET("/book/:id/chapter/:chapterNo", bookFileHandler.GetChapterContent)
			manage.POST("/book/upload/page", bookFileHandler.PageUpload)
			manage.POST("/book/file/page", bookFileHandler.PageFile)
			manage.POST("/book/chapter/page", bookFileHandler.PageChapter)
			manage.POST("/book/chapter/list", bookFileHandler.ListChapter)
			manage.POST("/book/re-parse", middleware.RequireButton(authSvc, "book:update"), bookFileHandler.ReParseChapters)

			// --- Book Chapter Rule ---
			manage.GET("/book/chapter-rule/:id", bookFileHandler.GetChapterRuleByID)
			manage.POST("/book/chapter-rule/page", bookFileHandler.PageChapterRule)
			manage.POST("/book/chapter-rule", middleware.RequireButton(authSvc, "book:create"), bookFileHandler.CreateChapterRule)
			manage.PUT("/book/chapter-rule/:id", middleware.RequireButton(authSvc, "book:update"), bookFileHandler.UpdateChapterRule)
			manage.DELETE("/book/chapter-rule/:id", middleware.RequireButton(authSvc, "book:delete"), bookFileHandler.DeleteChapterRule)

			// --- Book Chapter Rule Rel ---
			manage.POST("/book/chapter-rule/bind", bookFileHandler.BindChapterRule)
			manage.DELETE("/book/chapter-rule/bind/:bookId", bookFileHandler.UnbindChapterRule)
			manage.GET("/book/chapter-rule/bind/:bookId", bookFileHandler.GetBoundChapterRule)

			// --- Book Filter Rule ---
			manage.GET("/book/filter-rule/:id", bookFileHandler.GetFilterRuleByID)
			manage.POST("/book/filter-rule/page", bookFileHandler.PageFilterRule)
			manage.POST("/book/filter-rule", middleware.RequireButton(authSvc, "book:create"), bookFileHandler.CreateFilterRule)
			manage.PUT("/book/filter-rule/:id", middleware.RequireButton(authSvc, "book:update"), bookFileHandler.UpdateFilterRule)
			manage.DELETE("/book/filter-rule/:id", middleware.RequireButton(authSvc, "book:delete"), bookFileHandler.DeleteFilterRule)

			// --- Reader Bookshelf ---
			manage.POST("/bookshelf", readerHandler.AddToBookshelf)
			manage.DELETE("/bookshelf/:bookId", readerHandler.RemoveFromBookshelf)
			manage.PUT("/bookshelf/:bookId", readerHandler.UpdateBookshelf)
			manage.POST("/bookshelf/page", readerHandler.GetBookshelfPage)
			manage.GET("/bookshelf/groups", readerHandler.ListGroups)

			// --- Reader Progress ---
			manage.PUT("/reader/progress/:bookId", readerHandler.ReportProgress)
			manage.GET("/reader/progress/:bookId", readerHandler.GetProgress)

			// --- Book Social: Reader Note ---
			manage.POST("/reader/note", noteHandler.CreateNote)
			manage.PUT("/reader/note/:id", noteHandler.UpdateNote)
			manage.DELETE("/reader/note/:id", noteHandler.DeleteNote)
			manage.GET("/reader/note/:id", noteHandler.GetNote)
			manage.POST("/reader/note/page", noteHandler.PageNote)
			manage.GET("/reader/note/book/:bookId", noteHandler.ListNotesByBook)

			// --- Book Social: Book Review ---
			manage.POST("/book-review", reviewHandler.CreateReview)
			manage.PUT("/book-review/:id", reviewHandler.UpdateReview)
			manage.DELETE("/book-review/:id", reviewHandler.DeleteReview)
			manage.GET("/book-review/:id", reviewHandler.GetReview)
			manage.POST("/book-review/page", reviewHandler.PageReview)

			// --- Book Social: Chapter Comment ---
			manage.POST("/book/comment", commentHandler.CreateComment)
			manage.DELETE("/book/comment/:id", commentHandler.DeleteComment)
			manage.GET("/book/comment/:id", commentHandler.GetComment)
			manage.POST("/book/comment/page", commentHandler.PageComment)

			// --- Book Social: Like ---
			manage.POST("/like/toggle", likeHandler.ToggleLike)
			manage.POST("/like/status", likeHandler.GetLikeStatus)
			manage.GET("/like/count/:targetType/:targetId", likeHandler.CountLikes)
		}
	}

	return r
}
