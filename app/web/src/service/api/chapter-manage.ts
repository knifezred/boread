import { request } from "../request";

// =====================================================================
// 书籍章节管理 API
// =====================================================================

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

/** 读取章节内容（通过章节ID） */
export function fetchGetChapterContentByID(chapterId: string | number) {
  return request<Api.BookManage.ChapterContentResponse>({
    url: `/book/chapter/${chapterId}/content`,
    method: "get",
  });
}

/** 更新章节标题 */
export function fetchUpdateChapterTitle(id: string | number, data: Api.BookManage.ChapterTitleUpdateRequest) {
  return request<null>({
    url: `/book/chapter/${id}/title`,
    method: "put",
    data,
  });
}

/** 批量更新章节标题 */
export function fetchBatchUpdateChapterTitle(data: Api.BookManage.ChapterTitleBatchRequest) {
  return request<null>({
    url: "/book/chapter/batch-title",
    method: "put",
    data,
  });
}

/** 批量更新章节状态 */
export function fetchUpdateChapterStatus(data: Api.BookManage.ChapterStatusBatchRequest) {
  return request<null>({
    url: "/book/chapter/status",
    method: "put",
    data,
  });
}

/** 删除章节（软删除） */
export function fetchDeleteChapter(id: string | number) {
  return request<null>({
    url: `/book/chapter/${id}`,
    method: "delete",
  });
}

/** 合并章节 */
export function fetchMergeChapters(data: Api.BookManage.ChapterMergeRequest) {
  return request<null>({
    url: "/book/chapter/merge",
    method: "post",
    data,
  });
}

/** 格式化章节编号 */
export function fetchFormatChapterNumbers(data: Api.BookManage.ChapterFormatRequest) {
  return request<null>({
    url: "/book/chapter/format-numbers",
    method: "post",
    data,
  });
}

/** 保存章节内容 */
export function fetchSaveChapterContent(id: string | number, data: Api.BookManage.ChapterContentSaveRequest) {
  return request<null>({
    url: `/book/chapter/${id}/content`,
    method: "put",
    data,
  });
}


