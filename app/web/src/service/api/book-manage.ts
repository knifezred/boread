import { request } from "../request"

// =====================================================================
// 后端: boread (Go + Gin)
// 路径约定: /api 作为 baseURL, 此处只写 /manage/* 子路径
// =====================================================================

// -------- Book Category --------

/** 分类树 */
export function fetchGetCategoryTree() {
  return request<Api.SystemManage.BookCategory[]>({
    url: "/manage/book-category/tree",
    method: "get",
  });
}

/** 分类分页 */
export function fetchGetCategoryList(params?: Api.SystemManage.CategorySearchParams) {
  return request<Api.Common.PaginatingQueryRecord<Api.SystemManage.BookCategory>>({
    url: "/manage/book-category/page",
    method: "post",
    data: params,
  });
}

/** 分类详情 */
export function fetchGetCategory(id: string | number) {
  return request<Api.SystemManage.BookCategory>({
    url: `/manage/book-category/${id}`,
    method: "get",
  });
}

/** 新增分类 */
export function fetchCreateCategory(data: Record<string, any>) {
  return request<Api.SystemManage.BookCategory>({
    url: "/manage/book-category",
    method: "post",
    data,
  });
}

/** 编辑分类 */
export function fetchUpdateCategory(id: string | number, data: Record<string, any>) {
  return request<Api.SystemManage.BookCategory>({
    url: `/manage/book-category/${id}`,
    method: "put",
    data,
  });
}

/** 删除分类 */
export function fetchDeleteCategory(id: string | number) {
  return request<null>({
    url: `/manage/book-category/${id}`,
    method: "delete",
  });
}

// -------- Book Tag --------

/** 标签分页 */
export function fetchGetTagList(params?: Api.SystemManage.TagSearchParams) {
  return request<Api.Common.PaginatingQueryRecord<Api.SystemManage.BookTag>>({
    url: "/manage/book-tag/page",
    method: "post",
    data: params,
  });
}

/** 标签详情 */
export function fetchGetTag(id: string | number) {
  return request<Api.SystemManage.BookTag>({
    url: `/manage/book-tag/${id}`,
    method: "get",
  });
}

/** 新增标签 */
export function fetchCreateTag(data: Record<string, any>) {
  return request<Api.SystemManage.BookTag>({
    url: "/manage/book-tag",
    method: "post",
    data,
  });
}

/** 编辑标签 */
export function fetchUpdateTag(id: string | number, data: Record<string, any>) {
  return request<Api.SystemManage.BookTag>({
    url: `/manage/book-tag/${id}`,
    method: "put",
    data,
  });
}

/** 删除标签 */
export function fetchDeleteTag(id: string | number) {
  return request<null>({
    url: `/manage/book-tag/${id}`,
    method: "delete",
  });
}

// -------- Book (Novel) --------

/** 书籍分页 */
export function fetchGetBookList(params?: Api.SystemManage.BookSearchParams) {
  return request<Api.Common.PaginatingQueryRecord<Api.SystemManage.Book>>({
    url: "/manage/book/page",
    method: "post",
    data: params,
  });
}

/** 书籍详情 */
export function fetchGetBook(id: string | number) {
  return request<Api.SystemManage.Book>({
    url: `/manage/book/${id}`,
    method: "get",
  });
}

/** 新增书籍 */
export function fetchCreateBook(data: Api.SystemManage.BookRequest) {
  return request<Api.SystemManage.Book>({
    url: "/manage/book",
    method: "post",
    data,
  });
}

/** 编辑书籍 */
export function fetchUpdateBook(id: string | number, data: Api.SystemManage.BookRequest) {
  return request<Api.SystemManage.Book>({
    url: `/manage/book/${id}`,
    method: "put",
    data,
  });
}

/** 删除书籍 */
export function fetchDeleteBook(id: string | number) {
  return request<null>({
    url: `/manage/book/${id}`,
    method: "delete",
  });
}

/** 更新书籍上架状态 */
export function fetchUpdateBookStatus(id: string | number, data: Api.SystemManage.BookUpdateStatusRequest) {
  return request<null>({
    url: `/manage/book/${id}/status`,
    method: "put",
    data,
  });
}

// -------- Book File Upload & Scan --------

/** 上传小说文件 */
export function fetchUploadBookFile(file: File, onUploadProgress?: (percent: number) => void) {
  const formData = new FormData();
  formData.append("file", file);
  return request<Api.SystemManage.FileUploadResponse>({
    url: "/manage/book/upload",
    method: "post",
    data: formData,
    headers: { "Content-Type": "multipart/form-data" },
    onUploadProgress: (e: ProgressEvent) => {
      if (e.total && onUploadProgress) {
        onUploadProgress(Math.round((e.loaded / e.total) * 100));
      }
    },
  });
}

/** 确认入库 */
export function fetchConfirmImport(data: Api.SystemManage.ConfirmImportRequest) {
  return request<Api.SystemManage.ConfirmImportResponse>({
    url: "/manage/book/confirm-import",
    method: "post",
    data,
  });
}

/** 扫描本地目录 */
export function fetchScanPath(path: string) {
  return request<Api.SystemManage.ScanPathResponse>({
    url: "/manage/book/scan-path",
    method: "post",
    data: { path },
  });
}

/** 批量扫描入库 */
export function fetchScanAll() {
  return request<Api.SystemManage.ScanAllResponse>({
    url: "/manage/book/scan",
    method: "post",
  });
}

/** 扫描单个上传任务 */
export function fetchScanByID(id: string | number) {
  return request<Api.SystemManage.ScanResult>({
    url: `/manage/book/scan/${id}`,
    method: "post",
  });
}

/** 上传记录分页 */
export function fetchGetUploadList(params?: Api.SystemManage.UploadSearchParams) {
  return request<Api.Common.PaginatingQueryRecord<Api.SystemManage.BookUpload>>({
    url: "/manage/book/upload/page",
    method: "post",
    data: params,
  });
}

/** 文件记录分页 */
export function fetchGetFileList(params?: Api.SystemManage.FileSearchParams) {
  return request<Api.Common.PaginatingQueryRecord<Api.SystemManage.BookFile>>({
    url: "/manage/book/file/page",
    method: "post",
    data: params,
  });
}

/** 章节分页 */
export function fetchGetChapterList(params?: Api.SystemManage.ChapterSearchParams) {
  return request<Api.Common.PaginatingQueryRecord<Api.SystemManage.BookChapter>>({
    url: "/manage/book/chapter/page",
    method: "post",
    data: params,
  });
}

/** 读取章节内容 */
export function fetchGetChapterContent(bookId: string | number, chapterNo: string | number) {
  return request<Api.SystemManage.ChapterContentResponse>({
    url: `/manage/book/${bookId}/chapter/${chapterNo}`,
    method: "get",
  });
}

// -------- Book Chapter Rule --------

/** 章节识别规则分页 */
export function fetchGetChapterRuleList(params?: Api.SystemManage.ChapterRuleSearchParams) {
  return request<Api.Common.PaginatingQueryRecord<Api.SystemManage.BookChapterRule>>({
    url: "/manage/book/chapter-rule/page",
    method: "post",
    data: params,
  });
}

/** 章节识别规则详情 */
export function fetchGetChapterRule(id: string | number) {
  return request<Api.SystemManage.BookChapterRule>({
    url: `/manage/book/chapter-rule/${id}`,
    method: "get",
  });
}

/** 新增章节识别规则 */
export function fetchCreateChapterRule(data: Api.SystemManage.ChapterRuleRequest) {
  return request<Api.SystemManage.BookChapterRule>({
    url: "/manage/book/chapter-rule",
    method: "post",
    data,
  });
}

/** 编辑章节识别规则 */
export function fetchUpdateChapterRule(id: string | number, data: Api.SystemManage.ChapterRuleRequest) {
  return request<Api.SystemManage.BookChapterRule>({
    url: `/manage/book/chapter-rule/${id}`,
    method: "put",
    data,
  });
}

/** 删除章节识别规则 */
export function fetchDeleteChapterRule(id: string | number) {
  return request<null>({
    url: `/manage/book/chapter-rule/${id}`,
    method: "delete",
  });
}

// -------- Book Content Filter Rule --------

/** 内容净化规则分页 */
export function fetchGetFilterRuleList(params?: Api.SystemManage.FilterRuleSearchParams) {
  return request<Api.Common.PaginatingQueryRecord<Api.SystemManage.BookContentFilterRule>>({
    url: "/manage/book/filter-rule/page",
    method: "post",
    data: params,
  });
}

/** 内容净化规则详情 */
export function fetchGetFilterRule(id: string | number) {
  return request<Api.SystemManage.BookContentFilterRule>({
    url: `/manage/book/filter-rule/${id}`,
    method: "get",
  });
}

/** 新增内容净化规则 */
export function fetchCreateFilterRule(data: Api.SystemManage.FilterRuleRequest) {
  return request<Api.SystemManage.BookContentFilterRule>({
    url: "/manage/book/filter-rule",
    method: "post",
    data,
  });
}

/** 编辑内容净化规则 */
export function fetchUpdateFilterRule(id: string | number, data: Api.SystemManage.FilterRuleRequest) {
  return request<Api.SystemManage.BookContentFilterRule>({
    url: `/manage/book/filter-rule/${id}`,
    method: "put",
    data,
  });
}

/** 删除内容净化规则 */
export function fetchDeleteFilterRule(id: string | number) {
  return request<null>({
    url: `/manage/book/filter-rule/${id}`,
    method: "delete",
  });
}
