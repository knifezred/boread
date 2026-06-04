import { request } from "../request"

// =====================================================================
// 后端: boread (Go + Gin)
// 路径约定: /api 作为 baseURL, 此处只写 /manage/* 子路径 (系统管理模块)
// 书籍业务模块已迁移至 /book/* 直接路径
// =====================================================================

// -------- Book Category --------

/** 分类树 */
export function fetchGetCategoryTree() {
  return request<Api.BookManage.BookCategory[]>({
    url: "/book/category/tree",
    method: "get",
  });
}

/** 分类分页 */
export function fetchGetCategoryList(params?: Api.BookManage.CategorySearchParams) {
  return request<Api.Common.PaginatingQueryRecord<Api.BookManage.BookCategory>>({
    url: "/book/category/page",
    method: "post",
    data: params,
  });
}

/** 分类详情 */
export function fetchGetCategory(id: string | number) {
  return request<Api.BookManage.BookCategory>({
    url: `/book/category/${id}`,
    method: "get",
  });
}

/** 新增分类 */
export function fetchCreateCategory(data: Record<string, any>) {
  return request<Api.BookManage.BookCategory>({
    url: "/book/category",
    method: "post",
    data,
  });
}

/** 编辑分类 */
export function fetchUpdateCategory(id: string | number, data: Record<string, any>) {
  return request<Api.BookManage.BookCategory>({
    url: `/book/category/${id}`,
    method: "put",
    data,
  });
}

/** 删除分类 */
export function fetchDeleteCategory(id: string | number) {
  return request<null>({
    url: `/book/category/${id}`,
    method: "delete",
  });
}

/** 热门分类列表（公开接口，无需登录） */
export function fetchGetHotCategoryList() {
  return request<Api.BookManage.HotCategoryItem[]>({
    url: "/book/category/hot",
    method: "get",
  });
}

// -------- Book Tag --------

/** 标签分页 */
export function fetchGetTagList(params?: Api.BookManage.TagSearchParams) {
  return request<Api.Common.PaginatingQueryRecord<Api.BookManage.BookTag>>({
    url: "/book/tag/page",
    method: "post",
    data: params,
  });
}

/** 标签详情 */
export function fetchGetTag(id: string | number) {
  return request<Api.BookManage.BookTag>({
    url: `/book/tag/${id}`,
    method: "get",
  });
}

/** 新增标签 */
export function fetchCreateTag(data: Record<string, any>) {
  return request<Api.BookManage.BookTag>({
    url: "/book/tag",
    method: "post",
    data,
  });
}

/** 编辑标签 */
export function fetchUpdateTag(id: string | number, data: Record<string, any>) {
  return request<Api.BookManage.BookTag>({
    url: `/book/tag/${id}`,
    method: "put",
    data,
  });
}

/** 删除标签 */
export function fetchDeleteTag(id: string | number) {
  return request<null>({
    url: `/book/tag/${id}`,
    method: "delete",
  });
}

// -------- Book (Novel) --------

/** 书籍分页 */
export function fetchGetBookList(params?: Api.BookManage.BookSearchParams) {
  return request<Api.Common.PaginatingQueryRecord<Api.BookManage.Book>>({
    url: "/book/page",
    method: "post",
    data: params,
  });
}

/** 书籍详情 */
export function fetchGetBook(id: string | number) {
  return request<Api.BookManage.Book>({
    url: `/book/${id}`,
    method: "get",
  });
}

/** 新增书籍 */
export function fetchCreateBook(data: Api.BookManage.BookRequest) {
  return request<Api.BookManage.Book>({
    url: "/book",
    method: "post",
    data,
  });
}

/** 编辑书籍 */
export function fetchUpdateBook(id: string | number, data: Api.BookManage.BookRequest) {
  return request<Api.BookManage.Book>({
    url: `/book/${id}`,
    method: "put",
    data,
  });
}

/** 删除书籍 */
export function fetchDeleteBook(id: string | number) {
  return request<null>({
    url: `/book/${id}`,
    method: "delete",
  });
}

/** 更新书籍上架状态 */
export function fetchUpdateBookStatus(id: string | number, data: Api.BookManage.BookUpdateStatusRequest) {
  return request<null>({
    url: `/book/${id}/status`,
    method: "put",
    data,
  });
}

// -------- Book File Upload & Scan --------

/** 上传小说文件 */
export function fetchUploadBookFile(file: File, onUploadProgress?: (percent: number) => void) {
  const formData = new FormData();
  formData.append("file", file);
  return request<Api.BookManage.FileUploadResponse>({
    url: "/book/upload",
    method: "post",
    data: formData,
    headers: { "Content-Type": "multipart/form-data" },
    onUploadProgress: (e) => {
      if (e.total && onUploadProgress) {
        onUploadProgress(Math.round((e.loaded / e.total) * 100));
      }
    },
  });
}

/** 确认入库 */
export function fetchConfirmImport(data: Api.BookManage.ConfirmImportRequest) {
  return request<Api.BookManage.ConfirmImportResponse>({
    url: "/book/confirm-import",
    method: "post",
    data,
  });
}

/** 扫描本地目录 */
export function fetchScanPath(path: string) {
  return request<Api.BookManage.ScanPathResponse>({
    url: "/book/scan-path",
    method: "post",
    data: { path },
  });
}

/** 批量扫描入库 */
export function fetchScanAll() {
  return request<Api.BookManage.ScanAllResponse>({
    url: "/book/scan",
    method: "post",
  });
}

/** 扫描单个上传任务 */
export function fetchScanByID(id: string | number) {
  return request<Api.BookManage.ScanResult>({
    url: `/book/scan/${id}`,
    method: "post",
  });
}

/** 上传记录分页 */
export function fetchGetUploadList(params?: Api.BookManage.UploadSearchParams) {
  return request<Api.Common.PaginatingQueryRecord<Api.BookManage.BookUpload>>({
    url: "/book/upload/page",
    method: "post",
    data: params,
  });
}

/** 文件记录分页 */
export function fetchGetFileList(params?: Api.BookManage.FileSearchParams) {
  return request<Api.Common.PaginatingQueryRecord<Api.BookManage.BookFile>>({
    url: "/book/file/page",
    method: "post",
    data: params,
  });
}

/** 章节分页 */
export function fetchGetChapterList(params?: Api.BookManage.ChapterSearchParams) {
  return request<Api.Common.PaginatingQueryRecord<Api.BookManage.BookChapter>>({
    url: "/book/chapter/page",
    method: "post",
    data: params,
  });
}

/** 章节列表（不分页） */
export function fetchChapterList(bookId: number) {
  return request<Api.BookManage.BookChapter[]>({
    url: "/book/chapter/list",
    method: "post",
    data: { bookId },
  });
}

/** 读取章节内容 */
export function fetchGetChapterContent(bookId: string | number, chapterNo: string | number) {
  return request<Api.BookManage.ChapterContentResponse>({
    url: `/book/${bookId}/chapter/${chapterNo}`,
    method: "get",
  });
}

/** 公开书籍详情 */
export function fetchGetPublicBook(id: string | number) {
  return request<Api.BookManage.Book>({
    url: `/book/${id}`,
    method: "get",
  });
}

/** 公开章节列表（不分页，一次拉取全部） */
export function fetchGetPublicChapterList(bookId: string | number) {
  return request<Api.Common.PaginatingQueryRecord<Api.BookManage.BookChapter>>({
    url: "/book/chapter/page",
    method: "post",
    data: { bookId: Number(bookId), current: 1, size: 99999 },
  });
}

/** 重新识别章节 */
export function fetchReParseChapters(bookId: number, ruleId?: number) {
  return request<Api.BookManage.ReParseResponse>({
    url: "/book/re-parse",
    method: "post",
    data: ruleId ? { bookId, ruleId } : { bookId },
  });
}

// -------- Book Chapter Rule --------

/** 章节识别规则分页 */
export function fetchGetChapterRuleList(params?: Api.BookManage.ChapterRuleSearchParams) {
  return request<Api.Common.PaginatingQueryRecord<Api.BookManage.BookChapterRule>>({
    url: "/book/chapter-rule/page",
    method: "post",
    data: params,
  });
}

/** 章节识别规则详情 */
export function fetchGetChapterRule(id: string | number) {
  return request<Api.BookManage.BookChapterRule>({
    url: `/book/chapter-rule/${id}`,
    method: "get",
  });
}

/** 新增章节识别规则 */
export function fetchCreateChapterRule(data: Api.BookManage.ChapterRuleRequest) {
  return request<Api.BookManage.BookChapterRule>({
    url: "/book/chapter-rule",
    method: "post",
    data,
  });
}

/** 编辑章节识别规则 */
export function fetchUpdateChapterRule(id: string | number, data: Api.BookManage.ChapterRuleRequest) {
  return request<Api.BookManage.BookChapterRule>({
    url: `/book/chapter-rule/${id}`,
    method: "put",
    data,
  });
}

/** 删除章节识别规则 */
export function fetchDeleteChapterRule(id: string | number) {
  return request<null>({
    url: `/book/chapter-rule/${id}`,
    method: "delete",
  });
}

// -------- Book Chapter Rule Bind --------

/** 为书籍绑定章节识别规则 */
export function fetchBindChapterRule(data: Api.BookManage.ChapterRuleBindRequest) {
  return request<Api.BookManage.ChapterRuleBindResponse>({
    url: "/book/chapter-rule/bind",
    method: "post",
    data,
  });
}

/** 解绑书籍的章节识别规则 */
export function fetchUnbindChapterRule(bookId: number) {
  return request<null>({
    url: `/book/chapter-rule/bind/${bookId}`,
    method: "delete",
  });
}

/** 获取书籍绑定的章节识别规则 */
export function fetchGetBoundChapterRule(bookId: number) {
  return request<Api.BookManage.ChapterRuleBindResponse | null>({
    url: `/book/chapter-rule/bind/${bookId}`,
    method: "get",
  });
}

// -------- Book Content Filter Rule --------

/** 内容净化规则分页 */
export function fetchGetFilterRuleList(params?: Api.BookManage.FilterRuleSearchParams) {
  return request<Api.Common.PaginatingQueryRecord<Api.BookManage.BookContentFilterRule>>({
    url: "/book/filter-rule/page",
    method: "post",
    data: params,
  });
}

/** 内容净化规则详情 */
export function fetchGetFilterRule(id: string | number) {
  return request<Api.BookManage.BookContentFilterRule>({
    url: `/book/filter-rule/${id}`,
    method: "get",
  });
}

/** 新增内容净化规则 */
export function fetchCreateFilterRule(data: Api.BookManage.FilterRuleRequest) {
  return request<Api.BookManage.BookContentFilterRule>({
    url: "/book/filter-rule",
    method: "post",
    data,
  });
}

/** 编辑内容净化规则 */
export function fetchUpdateFilterRule(id: string | number, data: Api.BookManage.FilterRuleRequest) {
  return request<Api.BookManage.BookContentFilterRule>({
    url: `/book/filter-rule/${id}`,
    method: "put",
    data,
  });
}

/** 删除内容净化规则 */
export function fetchDeleteFilterRule(id: string | number) {
  return request<null>({
    url: `/book/filter-rule/${id}`,
    method: "delete",
  });
}

// -------- Reader Bookshelf --------

/** 书架分页列表 */
export function fetchGetBookshelfPage(params?: Api.BookManage.BookshelfSearchParams) {
  return request<Api.Common.PaginatingQueryRecord<Api.BookManage.BookshelfItem>>({
    url: "/book/shelf/page",
    method: "post",
    data: params,
  });
}

/** 添加到书架 */
export function fetchAddToBookshelf(data: Api.BookManage.BookshelfRequest) {
  return request<Api.BookManage.BookshelfItem>({
    url: "/book/shelf",
    method: "post",
    data,
  });
}

/** 从书架移除 */
export function fetchRemoveFromBookshelf(bookId: string | number) {
  return request<null>({
    url: `/book/shelf/${bookId}`,
    method: "delete",
  });
}

/** 更新书架 */
export function fetchUpdateBookshelf(bookId: string | number, data: Api.BookManage.BookshelfUpdateRequest) {
  return request<Api.BookManage.BookshelfItem>({
    url: `/book/shelf/${bookId}`,
    method: "put",
    data,
  });
}

/** 获取书架分组列表 */
export function fetchListBookshelfGroups() {
  return request<Api.BookManage.BookshelfGroupItem[]>({
    url: "/book/shelf/groups",
    method: "get",
  });
}

// -------- Reader Read Progress --------

/** 上报阅读进度 */
export function fetchReportReadProgress(bookId: string | number, data: Api.BookManage.ReadProgressRequest) {
  return request<Api.BookManage.ReadProgressResponse>({
    url: `/book/reader/progress/${bookId}`,
    method: "put",
    data,
  });
}

/** 获取阅读进度 */
export function fetchGetReadProgress(bookId: string | number) {
  return request<Api.BookManage.ReadProgressResponse>({
    url: `/book/reader/progress/${bookId}`,
    method: "get",
  });
}

// -------- Reader Read Event --------

/** 上报阅读事件(心跳) */
export function fetchReportReadEvent(data: Api.BookManage.ReadEventRequest) {
  return request<null>({
    url: "/book/reader/read-event",
    method: "post",
    data,
  });
}
