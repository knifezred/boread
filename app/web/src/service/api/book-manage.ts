import { request } from "../request";

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