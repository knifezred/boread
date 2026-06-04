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
	bookCharacterRepo := repository.NewBookCharacterRepository(db)
	bookCharacterRelRepo := repository.NewBookCharacterRelRepository(db)
	readerReadEventRepo := repository.NewReaderReadEventRepository(db)
	readerReadStatsRepo := repository.NewReaderReadStatsRepository(db)

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
	bookReaderSvc := service.NewBookReaderService(db, readerBookshelfRepo, readerReadProgressRepo, readerReadEventRepo, bookRepo, bookChapterRepo)
	bookReadStatsSvc := service.NewBookReadStatsService(readerReadStatsRepo, bookRepo)
	bookSocialSvc := service.NewBookSocialService(db, readerNoteRepo, bookReviewRepo, bookCommentRepo, readerLikeRepo, bookRepo, bookChapterRepo, userRepo)
	bookCharacterSvc := service.NewBookCharacterService(db, bookCharacterRepo, bookRepo)
	bookCharacterRelSvc := service.NewBookCharacterRelService(db, bookCharacterRelRepo, bookCharacterRepo)

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
	bookReaderHandler := v1.NewBookReaderHandler(bookReaderSvc)
	bookReadStatsHandler := v1.NewBookReadStatsHandler(bookReadStatsSvc)
	bookshelfSvc := service.NewReaderBookshelfService(db, readerBookshelfRepo, readerReadProgressRepo, bookRepo, bookChapterRepo)
	bookshelfHandler := v1.NewBookshelfHandler(bookshelfSvc)
	noteHandler := v1.NewNoteHandler(bookSocialSvc)
	reviewHandler := v1.NewReviewHandler(bookSocialSvc)
	commentHandler := v1.NewCommentHandler(bookSocialSvc)
	likeHandler := v1.NewLikeHandler(bookSocialSvc)
	characterHandler := v1.NewCharacterHandler(bookCharacterSvc)
	characterRelHandler := v1.NewCharacterRelHandler(bookCharacterRelSvc)

	api := r.Group("/api")
	{
		// 公开
		api.POST("/auth/login", authHandler.Login)
		api.GET("/book/category/hot", bookCategoryHandler.HotList)

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
			bookCategory := api.Group("/book/category")
			bookCategory.Use(middleware.Auth())
			{
				bookCategory.GET("/tree", bookCategoryHandler.Tree)
				bookCategory.GET("/:id", bookCategoryHandler.GetByID)
				bookCategory.POST("/page", bookCategoryHandler.Page)
				bookCategory.POST("", middleware.RequireButton(authSvc, "book-category:create"), bookCategoryHandler.Create)
				bookCategory.PUT("/:id", middleware.RequireButton(authSvc, "book-category:update"), bookCategoryHandler.Update)
				bookCategory.DELETE("/:id", middleware.RequireButton(authSvc, "book-category:delete"), bookCategoryHandler.Delete)
			}

		}

		// ============== Book Modules (moved from manage) ==============

		// --- Book Tag ---
		bookTagGroup := api.Group("/book/tag")
		bookTagGroup.Use(middleware.Auth())
		{
			bookTagGroup.GET("/:id", bookTagHandler.GetByID)
			bookTagGroup.POST("/page", bookTagHandler.Page)
			bookTagGroup.POST("", middleware.RequireButton(authSvc, "book-tag:create"), bookTagHandler.Create)
			bookTagGroup.PUT("/:id", middleware.RequireButton(authSvc, "book-tag:update"), bookTagHandler.Update)
			bookTagGroup.DELETE("/:id", middleware.RequireButton(authSvc, "book-tag:delete"), bookTagHandler.Delete)
		}

		// --- Book (CRUD + File + Chapter + Rules) ---
		bookGroup := api.Group("/book")
		bookGroup.Use(middleware.Auth())
		{
			// Book CRUD
			bookGroup.GET("/:id", bookHandler.GetByID)
			bookGroup.POST("/page", bookHandler.Page)
			bookGroup.POST("", middleware.RequireButton(authSvc, "book:create"), bookHandler.Create)
			bookGroup.PUT("/:id", middleware.RequireButton(authSvc, "book:update"), bookHandler.Update)
			bookGroup.PUT("/:id/status", middleware.RequireButton(authSvc, "book:update"), bookHandler.UpdateStatus)
			bookGroup.DELETE("/:id", middleware.RequireButton(authSvc, "book:delete"), bookHandler.Delete)

			// Book File
			bookGroup.POST("/upload", middleware.RequireButton(authSvc, "book:create"), bookFileHandler.Upload)
			bookGroup.POST("/confirm-import", middleware.RequireButton(authSvc, "book:create"), bookFileHandler.ConfirmImport)
			bookGroup.POST("/scan", middleware.RequireButton(authSvc, "book:create"), bookFileHandler.ScanAll)
			bookGroup.POST("/scan-path", middleware.RequireButton(authSvc, "book:create"), bookFileHandler.ScanPath)
			bookGroup.POST("/scan/:id", middleware.RequireButton(authSvc, "book:create"), bookFileHandler.ScanByID)
			bookGroup.GET("/:id/chapter/:chapterNo", bookFileHandler.GetChapterContent)
			bookGroup.POST("/upload/page", bookFileHandler.PageUpload)
			bookGroup.POST("/file/page", bookFileHandler.PageFile)
			bookGroup.POST("/chapter/page", bookFileHandler.PageChapter)
			bookGroup.POST("/chapter/list", bookFileHandler.ListChapter)
			bookGroup.POST("/re-parse", middleware.RequireButton(authSvc, "book:update"), bookFileHandler.ReParseChapters)

			// Chapter Rule
			bookGroup.GET("/chapter-rule/:id", bookFileHandler.GetChapterRuleByID)
			bookGroup.POST("/chapter-rule/page", bookFileHandler.PageChapterRule)
			bookGroup.POST("/chapter-rule", middleware.RequireButton(authSvc, "book:create"), bookFileHandler.CreateChapterRule)
			bookGroup.PUT("/chapter-rule/:id", middleware.RequireButton(authSvc, "book:update"), bookFileHandler.UpdateChapterRule)
			bookGroup.DELETE("/chapter-rule/:id", middleware.RequireButton(authSvc, "book:delete"), bookFileHandler.DeleteChapterRule)

			// Chapter Rule Bind
			bookGroup.POST("/chapter-rule/bind", bookFileHandler.BindChapterRule)
			bookGroup.DELETE("/chapter-rule/bind/:bookId", bookFileHandler.UnbindChapterRule)
			bookGroup.GET("/chapter-rule/bind/:bookId", bookFileHandler.GetBoundChapterRule)

			// Filter Rule
			bookGroup.GET("/filter-rule/:id", bookFileHandler.GetFilterRuleByID)
			bookGroup.POST("/filter-rule/page", bookFileHandler.PageFilterRule)
			bookGroup.POST("/filter-rule", middleware.RequireButton(authSvc, "book:create"), bookFileHandler.CreateFilterRule)
			bookGroup.PUT("/filter-rule/:id", middleware.RequireButton(authSvc, "book:update"), bookFileHandler.UpdateFilterRule)
			bookGroup.DELETE("/filter-rule/:id", middleware.RequireButton(authSvc, "book:delete"), bookFileHandler.DeleteFilterRule)
		}

		// --- Book Character ---
		bookCharGroup := api.Group("/book/character")
		bookCharGroup.Use(middleware.Auth())
		{
			bookCharGroup.POST("", characterHandler.CreateCharacter)
			bookCharGroup.PUT("/:id", characterHandler.UpdateCharacter)
			bookCharGroup.DELETE("/:id", characterHandler.DeleteCharacter)
			bookCharGroup.GET("/:id", characterHandler.GetCharacter)
			bookCharGroup.POST("/page", characterHandler.PageCharacter)
			bookCharGroup.GET("/book/:bookId", characterHandler.ListByCharacterBook)

			// Character Rel
			bookCharGroup.POST("/rel", characterRelHandler.CreateRelation)
			bookCharGroup.DELETE("/rel/:id", characterRelHandler.DeleteRelation)
			bookCharGroup.GET("/rel/:id", characterRelHandler.GetRelation)
			bookCharGroup.GET("/rel/character/:characterId", characterRelHandler.ListRelationsByCharacter)
			bookCharGroup.GET("/rel/book/:bookId", characterRelHandler.ListRelationsByBook)
		}

		// --- Book Social: Note ---
		noteGroup := api.Group("/book/note")
		noteGroup.Use(middleware.Auth())
		{
			noteGroup.POST("", noteHandler.CreateNote)
			noteGroup.PUT("/:id", noteHandler.UpdateNote)
			noteGroup.DELETE("/:id", noteHandler.DeleteNote)
			noteGroup.GET("/:id", noteHandler.GetNote)
			noteGroup.POST("/page", noteHandler.PageNote)
			noteGroup.GET("/book/:bookId", noteHandler.ListNotesByBook)
		}

		// --- Book Social: Review ---
		reviewGroup := api.Group("/book/review")
		reviewGroup.Use(middleware.Auth())
		{
			reviewGroup.POST("", reviewHandler.CreateReview)
			reviewGroup.PUT("/:id", reviewHandler.UpdateReview)
			reviewGroup.DELETE("/:id", reviewHandler.DeleteReview)
			reviewGroup.GET("/:id", reviewHandler.GetReview)
			reviewGroup.POST("/page", reviewHandler.PageReview)
		}

		// --- Book Social: Comment ---
		commentGroup := api.Group("/book/comment")
		commentGroup.Use(middleware.Auth())
		{
			commentGroup.POST("", commentHandler.CreateComment)
			commentGroup.DELETE("/:id", commentHandler.DeleteComment)
			commentGroup.GET("/:id", commentHandler.GetComment)
			commentGroup.POST("/page", commentHandler.PageComment)
		}

		// --- Book Social: Like ---
		likeGroup := api.Group("/book/like")
		likeGroup.Use(middleware.Auth())
		{
			likeGroup.POST("/toggle", likeHandler.ToggleLike)
			likeGroup.POST("/status", likeHandler.GetLikeStatus)
			likeGroup.GET("/count/:targetType/:targetId", likeHandler.CountLikes)
		}

		// --- Book Shelf ---
		bookShelfGroup := api.Group("/book/shelf")
		bookShelfGroup.Use(middleware.Auth())
		{
			bookShelfGroup.POST("", bookshelfHandler.AddToBookshelf)
			bookShelfGroup.DELETE("/:bookId", bookshelfHandler.RemoveFromBookshelf)
			bookShelfGroup.PUT("/:bookId", bookshelfHandler.UpdateBookshelf)
			bookShelfGroup.POST("/page", bookshelfHandler.GetBookshelfPage)
			bookShelfGroup.GET("/groups", bookshelfHandler.ListGroups)
		}

		// --- Book Reader (阅读进度 + 阅读事件) ---
		bookReaderGroup := api.Group("/book/reader")
		bookReaderGroup.Use(middleware.Auth())
		{
			bookReaderGroup.PUT("/progress/:bookId", bookReaderHandler.ReportProgress)
			bookReaderGroup.GET("/progress/:bookId", bookReaderHandler.GetProgress)

			bookReaderGroup.POST("/read-event", bookReaderHandler.ReportEvent)
		}

		// --- Book Read Stats (阅读统计) ---
		bookReadStatsGroup := api.Group("/book/read-stats")
		bookReadStatsGroup.Use(middleware.Auth())
		{
			bookReadStatsGroup.POST("/daily", bookReadStatsHandler.GetDailyStats)
			bookReadStatsGroup.POST("/books", bookReadStatsHandler.GetBookStats)
			bookReadStatsGroup.GET("/total", bookReadStatsHandler.GetTotalStats)
		}
	}

	return r
}
