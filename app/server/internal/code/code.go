package code

import (
	"errors"

	"boread/internal/service"
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
	ReviewNotOwner = 3402 // 无权修改他人书评

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
	ServerError = 5001 // 服务器内部错误
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
	ReviewNotOwner:          "无权修改他人书评",
	CommentNotFound:         "评论不存在",
	CommentNotOwner:         "无权删除他人评论",
	ParentCommentNotFound:   "父评论不存在",
	CharacterNotFound:       "角色不存在",
	CharacterRelNotFound:    "角色关系不存在",
	SearchParamInvalid:      "搜索参数错误",
	ServerError:             "服务器内部错误",
}

// MapServiceError 将 service 层 sentinel error 映射为业务错误码
func MapServiceError(err error) int {
	switch {
	// auth
	case errors.Is(err, service.ErrUserNotFound),
		errors.Is(err, service.ErrInvalidPassword):
		return AuthFailed
	case errors.Is(err, service.ErrUserDisabled):
		return UserDisabled
	case errors.Is(err, service.ErrUserLocked):
		return UserLocked

	// book
	case errors.Is(err, service.ErrBookNotFound):
		return ResourceConflict
	case errors.Is(err, service.ErrBookTagInvalid),
		errors.Is(err, service.ErrCategoryInvalid):
		return ResourceProtected

	// book_category
	case errors.Is(err, service.ErrCategoryCodeExists),
		errors.Is(err, service.ErrCategoryParentNotFound):
		return ResourceConflict
	case errors.Is(err, service.ErrCategoryHasChildren):
		return ResourceProtected

	// book_character
	case errors.Is(err, service.ErrCharacterNotFound):
		return CharacterNotFound
	case errors.Is(err, service.ErrCharacterRelNotFound):
		return CharacterRelNotFound
	case errors.Is(err, service.ErrBookExists):
		return ResourceConflict

	// book_file
	case errors.Is(err, service.ErrFileTooLarge),
		errors.Is(err, service.ErrFileTypeUnsupported),
		errors.Is(err, service.ErrFileEmpty),
		errors.Is(err, service.ErrFileMD5Exists):
		return FileInvalid
	case errors.Is(err, service.ErrUploadNotFound),
		errors.Is(err, service.ErrChapterNotFound),
		errors.Is(err, service.ErrBookFileNotFound),
		errors.Is(err, service.ErrRuleNotFound),
		errors.Is(err, service.ErrFilterRuleNotFound):
		return ResourceConflict

	// book_chapter
	case errors.Is(err, service.ErrChapterContentTooLarge):
		return ChapterContentTooLarge
	case errors.Is(err, service.ErrChapterFileUpdateFailed):
		return ChapterFileUpdateFailed
	case errors.Is(err, service.ErrChapterMergeNotAdjacent):
		return ChapterMergeNotAdjacent

	// book_reader
	case errors.Is(err, service.ErrProgressNotFound):
		return ProgressNotFound
	case errors.Is(err, service.ErrBookNotExist):
		return ResourceConflict

	// book_shelf
	case errors.Is(err, service.ErrBookshelfNotFound):
		return BookshelfNotFound
	case errors.Is(err, service.ErrAlreadyInBookshelf):
		return AlreadyInBookshelf

	// book_social
	case errors.Is(err, service.ErrNoteNotFound):
		return NoteNotFound
	case errors.Is(err, service.ErrNoteNotOwner):
		return NoteNotOwner
	case errors.Is(err, service.ErrReviewNotFound):
		return ReviewNotFound
	case errors.Is(err, service.ErrCommentNotFound):
		return CommentNotFound
	case errors.Is(err, service.ErrCommentNotOwner):
		return CommentNotOwner
	case errors.Is(err, service.ErrParentCommentNotExist):
		return ParentCommentNotFound
	case errors.Is(err, service.ErrBookNotExists):
		return ResourceConflict
	case errors.Is(err, service.ErrChapterNotExists):
		return ChapterNotExist

	// book_tag
	case errors.Is(err, service.ErrTagNameExists):
		return ResourceConflict

	// dept
	case errors.Is(err, service.ErrDeptCodeExists),
		errors.Is(err, service.ErrParentNotFound):
		return ResourceConflict
	case errors.Is(err, service.ErrDeptHasChildren),
		errors.Is(err, service.ErrDeptHasUsers):
		return ResourceProtected

	// dict
	case errors.Is(err, service.ErrDictCodeExists):
		return ResourceConflict
	case errors.Is(err, service.ErrDictSystem):
		return ResourceProtected

	// menu
	case errors.Is(err, service.ErrMenuRouteExists):
		return ResourceConflict
	case errors.Is(err, service.ErrMenuSystem):
		return ResourceProtected
	case errors.Is(err, service.ErrMenuHasChildren):
		return HasBoundRef

	// role
	case errors.Is(err, service.ErrRoleCodeExists):
		return ResourceConflict
	case errors.Is(err, service.ErrRoleSystem):
		return ResourceProtected
	case errors.Is(err, service.ErrRoleHasUsers):
		return HasBoundRef

	// user
	case errors.Is(err, service.ErrUserExists):
		return ResourceConflict

	default:
		return ServerError
	}
}
