<script setup lang="tsx">
import { ref, computed } from "vue"
import type { Ref } from "vue"
import { NButton, NPopconfirm, NTag, NSpace, NSelect, NInput, NTree } from "naive-ui"
import type { TreeOption } from "naive-ui"
import { useBoolean } from "@sa/hooks"
import {
  bookStatusRecord,
} from "@/constants/business"
import { useDictItems } from "@/hooks/business/dict"
import {
  fetchGetBookList,
  fetchDeleteBook,
  fetchGetCategoryTree,
  fetchUpdateBookStatus
} from "@/service/api"
import { useAppStore } from "@/store/modules/app"
import {
  defaultTransform,
  useNaivePaginatedTable,
  useTableOperate,
} from "@/hooks/common/table"
import { $t } from "@/locales"
import BookOperateModal from "./modules/book-operate-modal.vue"
import BookUploadModal from "./modules/book-upload-modal.vue"
import BookChapterModal from "./modules/book-chapter-modal.vue"
import BookScanModal from "./modules/book-scan-modal.vue"

const appStore = useAppStore()

const { bool: visible, setTrue: openModal } = useBoolean()
const { bool: uploadVisible, setTrue: openUploadModal } = useBoolean()
const { bool: scanVisible, setTrue: openScanModal } = useBoolean()
const { bool: chapterVisible, setTrue: openChapterModal } = useBoolean()

const chapterBookId = ref(0)
const chapterBookTitle = ref("")

/** 分类树原始数据 */
const categoryTree = ref<Api.SystemManage.BookCategory[]>([])
/** 当前选中的分类 ID，0 表示"全部" */
const selectedCategory = ref<number>(0)
/** NTree 的选中 key */
const treeSelectedKey = ref<string | number>("0")
/** NTree 展开的 key 列表，保持展开状态 */
const expandedKeys = ref<Array<string | number>>(["0"])

/** 将分类树转为 NTree 的 TreeOption */
const treeOptions = computed<TreeOption[]>(() => {
  function toOptions(nodes: Api.SystemManage.BookCategory[]): TreeOption[] {
    return nodes.map(n => ({
      key: n.id,
      label: n.categoryName,
      children: n.children?.length ? toOptions(n.children) : undefined,
    }))
  }
  const allOption: TreeOption = { key: "0", label: $t("page.admin.library.book.totalCategories"), children: toOptions(categoryTree.value) }
  return [allOption]
})

loadCategoryTree()

const { options: serialStatusOptions, labelMap: serialStatusLabelMap } = useDictItems("book_serial_status")
const { options: visibilityOptions, labelMap: visibilityLabelMap } = useDictItems("book_visibility")

/** 加载分类树 */
async function loadCategoryTree() {
  const { data, error } = await fetchGetCategoryTree()
  if (!error && data) {
    categoryTree.value = data
    expandedKeys.value = collectAllTreeKeys(data, ["0"])
  }
}

/** 递归收集所有树节点 key */
function collectAllTreeKeys(nodes: Api.SystemManage.BookCategory[], keys: Array<string | number> = []): Array<string | number> {
  for (const n of nodes) {
    keys.push(n.id)
    if (n.children?.length) {
      collectAllTreeKeys(n.children, keys)
    }
  }
  return keys
}

const searchParams = ref<Api.SystemManage.BookSearchParams>({
  current: 1,
  size: 10,
  title: null,
  author: null,
  categoryId: null,
  status: null,
  visibility: null,
  serialStatus: null,
  tagId: null,
})

const {
  columns,
  columnChecks,
  data,
  loading,
  pagination,
  getData,
  getDataByPage,
} = useNaivePaginatedTable({
  api: () => fetchGetBookList(searchParams.value),
  onPaginationParamsChange: (params) => {
    searchParams.value.current = params.page
    searchParams.value.size = params.pageSize
  },
  transform: (response) => defaultTransform(response),
  columns: () => [
    { type: "selection", align: "center", width: 48 },
    {
      key: "index",
      title: $t("common.index"),
      align: "center",
      width: 64,
      render: (_, index) => index + 1,
    },
    {
      key: "categoryName",
      title: $t("page.admin.library.book.categoryId"),
      align: "center",
      width: 100,
    },
    {
      key: "title",
      title: $t("page.admin.library.book.bookName"),
      align: "left",
      ellipsis: { tooltip: true },
      minWidth: 160,
    },
    {
      key: "author",
      title: $t("page.admin.library.book.author"),
      align: "center",
      width: 140,
      ellipsis: { tooltip: true },
    },
    {
      key: "serialStatus",
      title: $t("page.admin.library.book.serialStatus"),
      align: "center",
      width: 90,
      render: (row: Api.SystemManage.Book) => {
        const tagMap: Record<string, NaiveUI.ThemeColor> = { "1": "info", "2": "success", "3": "warning" }
        return <NTag type={tagMap[row.serialStatus]}>{serialStatusLabelMap.value[row.serialStatus] ?? row.serialStatus}</NTag>
      },
    },
    {
      key: "visibility",
      title: $t("page.admin.library.book.visibility"),
      align: "center",
      width: 80,
      render: (row: Api.SystemManage.Book) => {
        const tagMap: Record<string, NaiveUI.ThemeColor> = { "1": "success", "2": "warning", "3": "info" }
        return <NTag type={tagMap[row.visibility]}>{visibilityLabelMap.value[row.visibility] ?? row.visibility}</NTag>
      },
    },
    {
      key: "totalChapters",
      title: $t("page.admin.library.book.totalChapters"),
      align: "center",
      width: 80,
    },
    {
      key: "totalWords",
      title: $t("page.admin.library.book.totalWords"),
      align: "center",
      width: 80,
    },
    {
      key: "status",
      title: $t("page.admin.library.book.listingStatus"),
      align: "center",
      width: 90,
      render: (row: Api.SystemManage.Book) => {
        const tagMap: Record<string, NaiveUI.ThemeColor> = { "1": "success", "2": "warning", "3": "info", "4": "error" }
        return <NTag type={tagMap[row.status]}>{$t(bookStatusRecord[row.status])}</NTag>
      },
    },
    {
      key: "avgRating",
      title: $t("page.admin.library.book.avgRating"),
      align: "center",
      width: 80,
    },
    {
      key: "operate",
      title: $t("common.operate"),
      align: "center",
      width: 380,
      render: (row: Api.SystemManage.Book) => (
        <div class="flex-center justify-end gap-8px">
          <NPopconfirm onPositiveClick={() => handleToggleListing(row)}>
            {{
              default: () => row.status === "1" ? $t("common.confirmDelete") : $t("page.admin.library.book.statusListed"),
              trigger: () => <NButton size="small" ghost>{row.status === "1" ? $t("page.admin.library.book.statusUnlisted") : $t("page.admin.library.book.statusListed")}</NButton>,
            }}
          </NPopconfirm>
          <NButton size="small" ghost onClick={() => showChapters(row)}>{$t("page.admin.library.book.chapters")}</NButton>
          <NButton type="primary" ghost size="small" onClick={() => handleEdit(row)}>{$t("common.edit")}</NButton>
          <NPopconfirm onPositiveClick={() => handleDelete(row.id)}>
            {{ default: () => $t("common.confirmDelete"), trigger: () => <NButton type="error" ghost size="small">{$t("common.delete")}</NButton> }}
          </NPopconfirm>
        </div>
      ),
    },
  ],
})

const { checkedRowKeys, onBatchDeleted, onDeleted } = useTableOperate(data, "id", getData)

const operateType = ref<NaiveUI.TableOperateType>("add")
const editingData: Ref<Api.SystemManage.Book | null> = ref(null)

function handleTreeSelect(keys: Array<string | number>) {
  const key = keys[0]
  if (key === undefined || key === "0") {
    selectedCategory.value = 0
    searchParams.value.categoryId = null
  } else {
    selectedCategory.value = Number(key)
    searchParams.value.categoryId = Number(key)
  }
  getDataByPage(1)
}

function handleAdd() {
  operateType.value = "add"
  editingData.value = null
  openModal()
}
function handleEdit(item: Api.SystemManage.Book) {
  operateType.value = "edit"
  editingData.value = { ...item }
  openModal()
}

async function handleDelete(id: number) {
  const { error } = await fetchDeleteBook(id)
  if (!error) onDeleted()
}

async function handleScan() {
  openScanModal()
}

function showChapters(row: Api.SystemManage.Book) {
  chapterBookId.value = row.id
  chapterBookTitle.value = row.title
  openChapterModal()
}

async function handleToggleListing(row: Api.SystemManage.Book) {
  const newStatus: Api.SystemManage.BookListingStatus = row.status === "1" ? "2" : "1"
  const { error } = await fetchUpdateBookStatus(row.id, { status: newStatus })
  if (!error) {
    window.$message?.success($t("common.updateSuccess"))
    getData()
  }
}

async function handleBatchDelete() { onBatchDeleted() }

function handleSearch() { getDataByPage(1) }
function handleReset() {
  searchParams.value = { current: 1, size: 10, title: null, author: null, categoryId: null, status: null, visibility: null, serialStatus: null, tagId: null }
  selectedCategory.value = 0
  treeSelectedKey.value = "0"
  getDataByPage(1)
}
</script>

<template>
  <div class="flex gap-16px overflow-hidden lt-sm:flex-col h-full">
    <NCard :bordered="false" size="small" class="tree-card w-240px shrink-0">
      <NTree
        :data="treeOptions"
        :selected-keys="[selectedCategory === 0 ? '0' : selectedCategory]"
        :expanded-keys="expandedKeys"
        class="max-h-[calc(100vh-248px)] overflow-y-auto"
        selectable
        @update:selected-keys="handleTreeSelect"
        @update:expanded-keys="expandedKeys = $event" />
    </NCard>
    <div class="flex-col-stretch gap-16px flex-1 overflow-hidden">
      <NCard :bordered="false" size="small">
        <NSpace wrap :size="[12, 12]">
          <NInput v-model:value="searchParams.title" :placeholder="$t('page.admin.library.book.bookName')" clearable
            style="width: 160px" @keyup.enter="handleSearch" />
          <NInput v-model:value="searchParams.author" :placeholder="$t('page.admin.library.book.author')" clearable
            style="width: 160px" @keyup.enter="handleSearch" />
          <NSelect v-model:value="searchParams.serialStatus"
            :placeholder="$t('page.admin.library.book.form.serialStatus')" :options="serialStatusOptions" clearable
            style="width: 130px" />
          <NSelect v-model:value="searchParams.visibility" :placeholder="$t('page.admin.library.book.form.visibility')"
            :options="visibilityOptions" clearable style="width: 130px" />
          <NButton type="primary" @click="handleSearch">{{ $t("common.search") }}</NButton>
          <NButton @click="handleReset">{{ $t("common.reset") }}</NButton>
          <NButton @click="openUploadModal">{{ $t("page.admin.library.book.upload") }}</NButton>
          <NButton @click="handleScan">{{ $t("page.admin.library.book.scan") }}</NButton>
        </NSpace>
      </NCard>
      <NCard :title="$t('page.admin.library.book.title')" :bordered="false" size="small"
        class="card-wrapper sm:flex-1-hidden">
        <template #header-extra>
          <TableHeaderOperation v-model:columns="columnChecks" :disabled-delete="checkedRowKeys.length === 0"
            :loading="loading" @add="handleAdd" @delete="handleBatchDelete" @refresh="getData" />
        </template>
        <NDataTable v-model:checked-row-keys="checkedRowKeys" :columns="columns" :data="data" size="small"
          :flex-height="!appStore.isMobile" :scroll-x="1200" :loading="loading" :row-key="(row) => row.id" remote
          :pagination="pagination" class="sm:h-full" />
        <BookOperateModal v-model:visible="visible" :operate-type="operateType" :row-data="editingData"
          @submitted="getDataByPage" />
      </NCard>
    </div>
    <BookUploadModal v-model:visible="uploadVisible" @imported="getData" />
    <BookScanModal v-model:visible="scanVisible" @scanned="getData" />
    <BookChapterModal v-model:visible="chapterVisible" :book-id="chapterBookId" :book-title="chapterBookTitle" />
  </div>
</template>

<style scoped></style>
