package code

import (
	"errors"
)

// 错误码分段:
// 1000-1999: 通用参数/文件校验
// 2000-2999: 认证授权
// 3000-3999: 业务资源/冲突 (按子域细分)
// 4000-4999: 分页/搜索参数
// 5000-5999: 系统/服务器内部错误

const (
	// ========== 通用 ==========
	ParamInvalid = 1001 // 请求参数无效
	FileInvalid  = 1002 // 文件错误（太大/格式不支持/为空/重复）

	// ========== 认证授权 ==========
	AuthFailed       = 2001 // 认证失败 / 用户名或密码错误
	TokenInvalid     = 2002 // Token 无效或过期
	UserDisabled     = 2003 // 账号已禁用
	UserLocked       = 2004 // 账号已锁定
	PermissionDenied = 2005 // 权限不足
	UgreenAuthFailed = 2006 // 绿联认证失败
	UgreenNewUser    = 2007 // 绿联新用户，需手动确认注册
	UgreenPwdLogin   = 2008 // 绿联用户禁止密码登录

	// ========== 业务 - 通用资源 ==========
	ResourceConflict  = 3001 // 资源不存在/已存在/业务冲突
	ResourceProtected = 3002 // 系统内置/有子节点不可操作
	HasBoundRef       = 3003 // 仍有绑定关系不可删除

	// 章节
	ChapterNotExist         = 3041 // 章节不存在
	ChapterFileUpdateFailed = 3042 // 章节文件更新失败
	ChapterMergeNotAdjacent = 3043 // 章节非连续无法合并
	ChapterContentTooLarge  = 3044 // 章节内容过大

	// ========== 业务 - 书架 (3100-3199) ==========
	BookshelfNotFound  = 3101 // 书架记录不存在
	AlreadyInBookshelf = 3102 // 该书已在书架中

	// ========== 业务 - 阅读进度 (3200-3299) ==========
	ProgressNotFound = 3201 // 阅读进度不存在

	// ========== 业务 - 笔记 (3300-3399) ==========
	NoteNotFound = 3301 // 笔记不存在
	NoteNotOwner = 3302 // 无权修改他人笔记

	// ========== 业务 - 书评 (3400-3499) ==========
	ReviewNotFound = 3401 // 书评不存在

	// ========== 业务 - 评论 (3500-3599) ==========
	CommentNotFound       = 3501 // 评论不存在
	CommentNotOwner       = 3502 // 无权删除他人评论
	ParentCommentNotFound = 3503 // 父评论不存在

	// ========== 业务 - 角色 (3600-3699) ==========
	CharacterNotFound = 3601 // 角色不存在

	// ========== 业务 - 角色关系 (3700-3799) ==========
	CharacterRelNotFound = 3701 // 角色关系不存在

	// ========== 搜索/分页参数 ==========
	SearchParamInvalid = 4001 // 搜索参数错误

	// ========== 系统 ==========
	ServerError     = 5001 // 服务器内部错误
	DBNotConfigured = 5002 // 数据库未配置
)

// ============================================================================
//  Sentinel Error — 各业务模块统一在此定义，供 service 层引用和 MapServiceError 映射
// ============================================================================

var (
	// ---- auth ----
	ErrUserNotFound    = errors.New("user not found")
	ErrInvalidPassword = errors.New("invalid password")
	ErrUserDisabled    = errors.New("user is disabled")
	ErrUserLocked      = errors.New("user is locked")

	// ---- book ----
	ErrBookNotFound    = errors.New("书籍不存在")
	ErrBookTagInvalid  = errors.New("标签无效")
	ErrCategoryInvalid = errors.New("分类不存在")

	// ---- book_tag ----
	ErrTagNameExists = errors.New("标签名已存在")

	// ---- book_category ----
	ErrCategoryCodeExists     = errors.New("分类编码已存在")
	ErrCategoryHasChildren    = errors.New("存在子分类, 不能删除")
	ErrCategoryParentNotFound = errors.New("父分类不存在")

	// ---- book_character ----
	ErrCharacterNotFound    = errors.New("角色不存在")
	ErrCharacterRelNotFound = errors.New("角色关系不存在")
	ErrBookExists           = errors.New("书籍不存在")

	// ---- book_file ----
	ErrFileTooLarge        = errors.New("文件大小超出限制 (最大 200MB)")
	ErrFileTypeUnsupported = errors.New("不支持的文件格式，仅支持 txt/epub/mobi/pdf")
	ErrFileEmpty           = errors.New("文件为空")
	ErrUploadNotFound      = errors.New("上传记录不存在")
	ErrChapterNotFound     = errors.New("章节不存在")
	ErrFilterRuleNotFound  = errors.New("过滤规则不存在")
	ErrBookFileNotFound    = errors.New("文件记录不存在")
	ErrFileMD5Exists       = errors.New("文件 MD5 已存在（重复文件）")

	// ---- book_chapter_rule ----
	ErrRuleNotFound = errors.New("规则不存在")

	// ---- book_chapter ----
	ErrChapterContentTooLarge  = errors.New("章节内容过大")
	ErrChapterFileUpdateFailed = errors.New("章节文件更新失败")
	ErrChapterMergeNotAdjacent = errors.New("章节非连续无法合并")
	ErrChapterNoConflict       = errors.New("章节编号冲突")

	// ---- book_social ----
	ErrNoteNotFound          = errors.New("笔记不存在")
	ErrReviewNotFound        = errors.New("书评不存在")
	ErrCommentNotFound       = errors.New("评论不存在")
	ErrNoteNotOwner          = errors.New("无权修改他人笔记")
	ErrCommentNotOwner       = errors.New("无权删除他人评论")
	ErrParentCommentNotExist = errors.New("父评论不存在")
	ErrBookNotExists         = errors.New("书籍不存在")
	ErrChapterNotExists      = errors.New("章节不存在")

	// ---- book_reader ----
	ErrBookNotExist     = errors.New("书籍不存在")
	ErrProgressNotFound = errors.New("阅读进度不存在")

	// ---- book_shelf ----
	ErrBookshelfNotFound  = errors.New("书架记录不存在")
	ErrAlreadyInBookshelf = errors.New("该书已在书架中")

	// ---- dept ----
	ErrDeptCodeExists  = errors.New("部门编码已存在")
	ErrDeptHasChildren = errors.New("存在子部门, 不能删除")
	ErrDeptHasUsers    = errors.New("部门下有用户, 不能删除")
	ErrParentNotFound  = errors.New("父部门不存在")

	// ---- dict ----
	ErrDictCodeExists = errors.New("字典编码已存在")
	ErrDictSystem     = errors.New("系统内置字典不可操作")

	// ---- menu ----
	ErrMenuRouteExists = errors.New("路由名已存在")
	ErrMenuSystem      = errors.New("系统内置菜单不可删除/不可改 route_name")
	ErrMenuHasChildren = errors.New("存在子菜单, 不能删除")

	// ---- role ----
	ErrRoleCodeExists = errors.New("角色编码已存在")
	ErrRoleSystem     = errors.New("系统内置角色不可操作")
	ErrRoleHasUsers   = errors.New("角色下还有用户, 不能删除")

	// ---- user ----
	ErrUserExists = errors.New("用户名已存在")

	// ---- setting ----
	ErrSettingKeyExists   = errors.New("配置键已存在")
	ErrSettingNotEditable = errors.New("该配置项不可编辑")
	ErrSettingSystem      = errors.New("系统内置配置不可删除")

	// ---- ugreen ----
	ErrUgreenAuthFailed    = errors.New("ugreen auth failed")
	ErrUgreenNewUser       = errors.New("new ugreen user, registration required")
	ErrUgreenPasswordLogin = errors.New("ugreen user cannot login with password")
)

// Text 根据错误码返回中文消息
func Text(code int) string {
	if msg, ok := errMsgMap[code]; ok {
		return msg
	}
	return "未知错误"
}

// errMsgMap 错误码与中文消息映射
var errMsgMap = map[int]string{
	ParamInvalid:            "请求参数无效",
	FileInvalid:             "文件错误",
	AuthFailed:              "认证失败",
	TokenInvalid:            "Token 无效或过期",
	UserDisabled:            "账号已禁用",
	UserLocked:              "账号已锁定",
	PermissionDenied:        "权限不足",
	UgreenAuthFailed:        "绿联认证失败",
	UgreenNewUser:           "绿联新用户，需确认后注册",
	UgreenPwdLogin:          "绿联用户禁止密码登录",
	ResourceConflict:        "资源冲突或不存在",
	ResourceProtected:       "系统内置资源不可操作",
	HasBoundRef:             "仍有绑定关系不可删除",
	ChapterNotExist:         "章节不存在",
	ChapterFileUpdateFailed: "章节文件更新失败",
	ChapterMergeNotAdjacent: "章节非连续无法合并",
	ChapterContentTooLarge:  "章节内容过大",
	BookshelfNotFound:       "书架记录不存在",
	AlreadyInBookshelf:      "该书已在书架中",
	ProgressNotFound:        "阅读进度不存在",
	NoteNotFound:            "笔记不存在",
	NoteNotOwner:            "无权修改他人笔记",
	ReviewNotFound:          "书评不存在",
	CommentNotFound:         "评论不存在",
	CommentNotOwner:         "无权删除他人评论",
	ParentCommentNotFound:   "父评论不存在",
	CharacterNotFound:       "角色不存在",
	CharacterRelNotFound:    "角色关系不存在",
	SearchParamInvalid:      "搜索参数错误",
	ServerError:             "服务器内部错误",
	DBNotConfigured:         "数据库未配置",
}

// MapServiceError 将 sentinel error 映射为业务错误码
func MapServiceError(err error) int {
	switch {
	// auth
	case errors.Is(err, ErrUserNotFound),
		errors.Is(err, ErrInvalidPassword):
		return AuthFailed
	case errors.Is(err, ErrUgreenAuthFailed):
		return UgreenAuthFailed
	case errors.Is(err, ErrUgreenNewUser):
		return UgreenNewUser
	case errors.Is(err, ErrUgreenPasswordLogin):
		return UgreenPwdLogin
	case errors.Is(err, ErrUserDisabled):
		return UserDisabled
	case errors.Is(err, ErrUserLocked):
		return UserLocked

	// book
	case errors.Is(err, ErrBookNotFound):
		return ResourceConflict
	case errors.Is(err, ErrBookTagInvalid),
		errors.Is(err, ErrCategoryInvalid):
		return ResourceProtected

	// book_category
	case errors.Is(err, ErrCategoryCodeExists),
		errors.Is(err, ErrCategoryParentNotFound):
		return ResourceConflict
	case errors.Is(err, ErrCategoryHasChildren):
		return ResourceProtected

	// book_character
	case errors.Is(err, ErrCharacterNotFound):
		return CharacterNotFound
	case errors.Is(err, ErrCharacterRelNotFound):
		return CharacterRelNotFound
	case errors.Is(err, ErrBookExists):
		return ResourceConflict

	// book_file
	case errors.Is(err, ErrFileTooLarge),
		errors.Is(err, ErrFileTypeUnsupported),
		errors.Is(err, ErrFileEmpty),
		errors.Is(err, ErrFileMD5Exists):
		return FileInvalid
	case errors.Is(err, ErrUploadNotFound),
		errors.Is(err, ErrChapterNotFound),
		errors.Is(err, ErrBookFileNotFound),
		errors.Is(err, ErrRuleNotFound),
		errors.Is(err, ErrFilterRuleNotFound):
		return ResourceConflict

	// book_chapter
	case errors.Is(err, ErrChapterContentTooLarge):
		return ChapterContentTooLarge
	case errors.Is(err, ErrChapterFileUpdateFailed):
		return ChapterFileUpdateFailed
	case errors.Is(err, ErrChapterMergeNotAdjacent):
		return ChapterMergeNotAdjacent

	// book_reader
	case errors.Is(err, ErrProgressNotFound):
		return ProgressNotFound
	case errors.Is(err, ErrBookNotExist):
		return ResourceConflict

	// book_shelf
	case errors.Is(err, ErrBookshelfNotFound):
		return BookshelfNotFound
	case errors.Is(err, ErrAlreadyInBookshelf):
		return AlreadyInBookshelf

	// book_social
	case errors.Is(err, ErrNoteNotFound):
		return NoteNotFound
	case errors.Is(err, ErrNoteNotOwner):
		return NoteNotOwner
	case errors.Is(err, ErrReviewNotFound):
		return ReviewNotFound
	case errors.Is(err, ErrCommentNotFound):
		return CommentNotFound
	case errors.Is(err, ErrCommentNotOwner):
		return CommentNotOwner
	case errors.Is(err, ErrParentCommentNotExist):
		return ParentCommentNotFound
	case errors.Is(err, ErrBookNotExists):
		return ResourceConflict
	case errors.Is(err, ErrChapterNotExists):
		return ChapterNotExist

	// book_tag
	case errors.Is(err, ErrTagNameExists):
		return ResourceConflict

	// dept
	case errors.Is(err, ErrDeptCodeExists),
		errors.Is(err, ErrParentNotFound):
		return ResourceConflict
	case errors.Is(err, ErrDeptHasChildren),
		errors.Is(err, ErrDeptHasUsers):
		return ResourceProtected

	// dict
	case errors.Is(err, ErrDictCodeExists):
		return ResourceConflict
	case errors.Is(err, ErrDictSystem):
		return ResourceProtected

	// menu
	case errors.Is(err, ErrMenuRouteExists):
		return ResourceConflict
	case errors.Is(err, ErrMenuSystem):
		return ResourceProtected
	case errors.Is(err, ErrMenuHasChildren):
		return HasBoundRef

	// role
	case errors.Is(err, ErrRoleCodeExists):
		return ResourceConflict
	case errors.Is(err, ErrRoleSystem):
		return ResourceProtected
	case errors.Is(err, ErrRoleHasUsers):
		return HasBoundRef

	// user
	case errors.Is(err, ErrUserExists):
		return ResourceConflict

	// setting
	case errors.Is(err, ErrSettingKeyExists):
		return ResourceConflict
	case errors.Is(err, ErrSettingNotEditable),
		errors.Is(err, ErrSettingSystem):
		return ResourceProtected

	default:
		return ServerError
	}
}
