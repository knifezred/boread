<script setup lang="tsx">
import { h, ref, computed, watch } from "vue"
import {
    NButton, NDataTable, NDrawer, NDrawerContent, NInput, NSelect, NSpace,
    NTag, NModal, NPopconfirm, NScrollbar,
} from "naive-ui"
import { useBoolean } from "@sa/hooks"
import {
    fetchGetChapterList, fetchGetChapterContentByID, fetchUpdateChapterTitle,
    fetchBatchUpdateChapterTitle, fetchUpdateChapterStatus, fetchDeleteChapter,
    fetchMergeChapters, fetchFormatChapterNumbers, fetchSaveChapterContent,
} from "@/service/api"
import { defaultTransform, useNaivePaginatedTable } from "@/hooks/common/table"
import { $t } from "@/locales"
import BookReparseModal from "./book-reparse-modal.vue"

defineOptions({ name: "BookChapterDrawer" })

interface Props { bookId: number; bookTitle: string }
const props = defineProps<Props>()
const visible = defineModel<boolean>("visible", { default: false })

// ==================== 搜索 ====================
const searchTitle = ref("")
const searchStatus = ref("")

const paginationRef = ref({ page: 1, pageSize: 20 })

const searchParams = computed(() => ({
  current: paginationRef.value.page || 1,
  size: paginationRef.value.pageSize || 20,
  bookId: props.bookId || null,
  fileId: null,
  chapterNo: null,
  title: searchTitle.value || null,
  status: searchStatus.value || null,
}))

// ==================== 表格 ====================
const { columns, data, pagination, loading, getDataByPage } = useNaivePaginatedTable({
  api: () => fetchGetChapterList(searchParams.value),
  onPaginationParamsChange: (params) => {
    paginationRef.value = { page: params.page || 1, pageSize: params.pageSize || 10 }
  },
  transform: (response) => defaultTransform(response),
  columns: () => [
    { type: "selection", align: "center", width: 48 },
    {
      key: "chapterNo", title: $t("page.admin.library.book.chapterNo"), align: "center", width: 90,
      render: (row: Api.BookManage.BookChapter) => {
        return h("span", { class: "font-mono" }, `第${String(row.chapterNo).padStart(3, "0")}章`)
      },
    },
    {
      key: "volumeNo", title: "卷号", align: "center", width: 80,
      render: (row: Api.BookManage.BookChapter) => {
        if (!row.volumeNo) return null
        return h(NTag, { type: "info", size: "small" }, () => row.volumeTitle || `卷${row.volumeNo}`)
      },
    },
    {
      key: "title", title: $t("page.admin.library.book.chapterTitle"), align: "left",
      ellipsis: { tooltip: true }, minWidth: 200,
    },
    {
      key: "wordCount", title: $t("page.admin.library.book.wordCount"), align: "center", width: 80,
    },
    {
      key: "status", title: $t("page.admin.library.book.chapterStatus"), align: "center", width: 90,
      render: (row: Api.BookManage.BookChapter) => {
        const typeMap: Record<string, NaiveUI.ThemeColor> = { "1": "success", "2": "warning", "3": "error" }
        const labelMap: Record<string, string> = { "1": "启用", "2": "草稿", "3": "下架" }
        return h(NTag, { type: typeMap[row.status] || "default", size: "small" }, () => labelMap[row.status] || row.status)
      },
    },
    {
      key: "updateTime", title: $t("common.updateTime"), align: "center", width: 170,
    },
    {
      key: "operate", title: $t("common.operate"), align: "center", width: 260, fixed: "right",
      render: (row: Api.BookManage.BookChapter) => (
        <div class="flex-center gap-4px">
          <NButton size="tiny" quaternary onClick={() => openContentEdit(row)}>编辑内容</NButton>
          <NButton size="tiny" quaternary onClick={() => openTitleEdit(row)}>编辑标题</NButton>
          <NButton size="tiny" quaternary onClick={() => handleToggleStatus(row)}>
            {row.status === "1" ? "禁用" : "启用"}
          </NButton>
          <NPopconfirm onPositiveClick={() => handleDelete(row)}>
            {{
              default: () => $t("common.confirmDelete"),
              trigger: () => <NButton size="tiny" quaternary type="error">{$t("common.delete")}</NButton>,
            }}
          </NPopconfirm>
        </div>
      ),
    },
  ],
})

// ==================== 选中行 ====================
const checkedRowKeys = ref<number[]>([])

// ==================== 内容编辑 ====================
const { bool: contentEditVisible, setTrue: openContentEditModal, setFalse: closeContentEditModal } = useBoolean()
const editingContentChapter = ref<Api.BookManage.BookChapter | null>(null)
const editingContent = ref("")
const contentSaving = ref(false)

function openContentEdit(row: Api.BookManage.BookChapter) {
  editingContentChapter.value = row
  editingContent.value = ""
  openContentEditModal()
  loadChapterContent(row.id)
}

async function loadChapterContent(chapterId: number) {
  const { data: contentData, error } = await fetchGetChapterContentByID(chapterId)
  if (!error && contentData) {
    editingContent.value = contentData.content
  }
}

async function handleSaveContent() {
  const chapter = editingContentChapter.value
  if (!chapter) return
  contentSaving.value = true
  const { error } = await fetchSaveChapterContent(chapter.id, {
    bookId: chapter.bookId,
    content: editingContent.value,
  })
  contentSaving.value = false
  if (!error) {
    window.$message?.success("内容保存成功")
    closeContentEditModal()
    getDataByPage(paginationRef.value.page || 1)
  }
}

// ==================== 标题编辑 ====================
const { bool: titleEditVisible, setTrue: openTitleEditModal, setFalse: closeTitleEditModal } = useBoolean()
const editingTitleChapter = ref<Api.BookManage.BookChapter | null>(null)
const editingTitle = ref("")

function openTitleEdit(row: Api.BookManage.BookChapter) {
  editingTitleChapter.value = row
  editingTitle.value = row.title
  openTitleEditModal()
}

async function handleSaveTitle() {
  const chapter = editingTitleChapter.value
  if (!chapter) return
  const { error } = await fetchUpdateChapterTitle(chapter.id, { title: editingTitle.value })
  if (!error) {
    window.$message?.success($t("common.updateSuccess"))
    closeTitleEditModal()
    getDataByPage(paginationRef.value.page || 1)
  }
}

// ==================== 禁用/启用 ====================
async function handleToggleStatus(row: Api.BookManage.BookChapter) {
  const newStatus = row.status === "1" ? "3" : "1"
  const { error } = await fetchUpdateChapterStatus({ ids: [row.id], status: newStatus })
  if (!error) {
    window.$message?.success($t("common.updateSuccess"))
    getDataByPage(paginationRef.value.page || 1)
  }
}

// ==================== 删除 ====================
async function handleDelete(row: Api.BookManage.BookChapter) {
  const { error } = await fetchDeleteChapter(row.id)
  if (!error) {
    window.$message?.success($t("common.deleteSuccess"))
    getDataByPage(paginationRef.value.page || 1)
  }
}

// ==================== 批量禁用/启用 ====================
async function handleBatchStatus(status: string) {
  if (checkedRowKeys.value.length === 0) {
    window.$message?.warning("请先选择章节")
    return
  }
  const { error } = await fetchUpdateChapterStatus({ ids: checkedRowKeys.value, status })
  if (!error) {
    window.$message?.success($t("common.updateSuccess"))
    checkedRowKeys.value = []
    getDataByPage(paginationRef.value.page || 1)
  }
}

// ==================== 合并章节 ====================
const { bool: mergeVisible, setTrue: openMergeModal, setFalse: closeMergeModal } = useBoolean()
const allChapters = ref<Api.BookManage.BookChapter[]>([])
const mergeTargetId = ref<number | null>(null)
const merging = ref(false)

async function openMerge() {
  if (checkedRowKeys.value.length < 2) {
    window.$message?.warning("请至少选择2个章节进行合并")
    return
  }
  mergeTargetId.value = null
  const { data: chapterData } = await fetchGetChapterList({
    current: 1, size: 9999, bookId: props.bookId || null,
    fileId: null, chapterNo: null, title: null, status: null,
  })
  allChapters.value = chapterData?.records || []
  openMergeModal()
}

async function handleMerge() {
  if (!mergeTargetId.value) {
    window.$message?.warning("请选择目标章节")
    return
  }
  const sourceIds = checkedRowKeys.value.filter((id) => id !== mergeTargetId.value)
  if (sourceIds.length === 0) {
    window.$message?.warning("请至少选择2个章节进行合并")
    return
  }
  merging.value = true
  const { error } = await fetchMergeChapters({
    bookId: props.bookId,
    targetId: mergeTargetId.value,
    sourceIds,
  })
  merging.value = false
  if (!error) {
    window.$message?.success("合并成功")
    checkedRowKeys.value = []
    closeMergeModal()
    getDataByPage(paginationRef.value.page || 1)
  }
}

// ==================== 批量改标题 ====================
const { bool: batchTitleVisible, setTrue: openBatchTitleModal, setFalse: closeBatchTitleModal } = useBoolean()
const batchTitleTemplate = ref("")

function openBatchTitle() {
  if (checkedRowKeys.value.length === 0) {
    window.$message?.warning("请先选择章节")
    return
  }
  batchTitleTemplate.value = ""
  openBatchTitleModal()
}

async function handleBatchTitle() {
  if (!batchTitleTemplate.value.trim()) {
    window.$message?.warning("请输入标题模板")
    return
  }
  const { error } = await fetchBatchUpdateChapterTitle({
    ids: checkedRowKeys.value,
    title: batchTitleTemplate.value,
  })
  if (!error) {
    window.$message?.success($t("common.updateSuccess"))
    checkedRowKeys.value = []
    closeBatchTitleModal()
    getDataByPage(paginationRef.value.page || 1)
  }
}

// ==================== 格式化编号 ====================
async function handleFormatNumbers() {
  if (checkedRowKeys.value.length === 0) {
    window.$message?.warning("请先选择章节")
    return
  }
  const { error } = await fetchFormatChapterNumbers({ ids: checkedRowKeys.value })
  if (!error) {
    window.$message?.success("格式化成功")
    checkedRowKeys.value = []
    getDataByPage(paginationRef.value.page || 1)
  }
}

// ==================== 重新识别 ====================
const reparseVisible = ref(false)

function handleReparsed() {
  reparseVisible.value = false
  getDataByPage(1)
}

// ==================== 生命周期 ====================
watch(visible, (val) => {
  if (val) {
    searchTitle.value = ""
    searchStatus.value = ""
    checkedRowKeys.value = []
    getDataByPage(1)
  }
})

function closeDrawer() {
  visible.value = false
}
</script>

<template>
  <NDrawer v-model:show="visible" display-directive="show" :width="960" native-scrollbar>
    <NDrawerContent
      :title="`${$t('page.admin.library.book.chapters')} - ${props.bookTitle}`"
      :native-scrollbar="false"
      closable
      @after-leave="closeDrawer"
    >
      <!-- 搜索栏 -->
      <div class="flex items-center gap-12px mb-16px">
        <NInput
          v-model:value="searchTitle"
          placeholder="搜索章节标题"
          clearable
          style="width: 240px"
          size="small"
        />
        <NSelect
          v-model:value="searchStatus"
          placeholder="全部状态"
          :options="[
            { label: '全部', value: '' },
            { label: '启用', value: '1' },
            { label: '草稿', value: '2' },
            { label: '下架', value: '3' },
          ]"
          clearable
          style="width: 130px"
          size="small"
        />
        <NButton size="small" type="primary" ghost @click="getDataByPage(1)">
          {{ $t("common.search") }}
        </NButton>
        <NButton size="small" @click="searchTitle = ''; searchStatus = ''; getDataByPage(1)">
          {{ $t("common.reset") }}
        </NButton>
      </div>

      <!-- 章节表格 -->
      <NDataTable
        v-model:checked-row-keys="checkedRowKeys"
        :columns="columns"
        :data="data"
        size="small"
        :loading="loading"
        :pagination="pagination"
        :row-key="(row: Record<string, any>) => row.id"
        remote
        :flex-height="true"
        min-height="400"
        :scroll-x="1200"
      />

      <!-- 底部操作栏 -->
      <template #footer>
        <div class="flex items-center gap-8px">
          <NButton size="tiny" :disabled="checkedRowKeys.length === 0" @click="handleFormatNumbers">
            格式化编号
          </NButton>
          <NButton size="tiny" :disabled="checkedRowKeys.length < 2" @click="openMerge">
            合并章节
          </NButton>
          <NButton size="tiny" :disabled="checkedRowKeys.length === 0" @click="openBatchTitle">
            批量改标题
          </NButton>
          <NButton size="tiny" :disabled="checkedRowKeys.length === 0" @click="handleBatchStatus('1')">
            启用
          </NButton>
          <NButton size="tiny" :disabled="checkedRowKeys.length === 0" @click="handleBatchStatus('3')">
            禁用
          </NButton>
          <NButton size="tiny" ghost @click="reparseVisible = true">
            {{ $t("page.admin.library.bookChapterRule.form.reParse.title") }}
          </NButton>
          <div class="flex-1" />
          <NButton size="tiny" @click="closeDrawer">
            {{ $t("common.close") }}
          </NButton>
        </div>
      </template>

      <!-- 内容编辑弹窗 -->
      <NModal
        v-model:show="contentEditVisible" preset="card" title="内容编辑" style="width: 800px"
        :bordered="false" :segmented="false"
      >
        <div class="mb-12px text-14px font-medium">
          {{ editingContentChapter?.title }}
        </div>
        <NInput
          v-model:value="editingContent"
          type="textarea"
          :autosize="{ minRows: 15, maxRows: 30 }"
          placeholder="请输入章节内容"
        />
        <template #footer>
          <NSpace justify="end">
            <NButton @click="closeContentEditModal">{{ $t("common.cancel") }}</NButton>
            <NButton type="primary" :loading="contentSaving" @click="handleSaveContent">
              保存内容
            </NButton>
          </NSpace>
        </template>
      </NModal>

      <!-- 标题编辑弹窗 -->
      <NModal
        v-model:show="titleEditVisible" preset="dialog" title="编辑标题"
        positive-text="确认" negative-text="取消"
        @positive-click="handleSaveTitle" @negative-click="closeTitleEditModal"
      >
        <NInput v-model:value="editingTitle" :placeholder="$t('page.admin.library.book.chapterTitle')" />
      </NModal>

      <!-- 合并弹窗 -->
      <NModal v-model:show="mergeVisible" preset="card" title="合并章节" style="width: 560px" :bordered="false">
        <div class="mb-12px text-14px">
          请选择目标章节，选中的源章节将合并到目标章节中
        </div>
        <NScrollbar style="max-height: 360px">
          <div
            v-for="ch in allChapters" :key="ch.id"
            class="flex items-center gap-8px py-6px px-8px rounded-4px cursor-pointer hover:bg-[#f5f5f5]"
            :class="{ 'bg-[#e6f7ff]': mergeTargetId === ch.id }"
            @click="mergeTargetId = ch.id"
          >
            <div
              class="w-16px h-16px rounded-full border-2 flex-shrink-0"
              :class="mergeTargetId === ch.id ? 'border-[#2080f0] bg-[#2080f0]' : 'border-#ccc'"
            />
            <span class="font-mono text-12px w-80px text-right">
              第{{ String(ch.chapterNo).padStart(3, "0") }}章
            </span>
            <span class="flex-1 truncate">{{ ch.title }}</span>
            <NTag v-if="checkedRowKeys.includes(ch.id)" size="tiny" type="warning">
              源章节
            </NTag>
          </div>
        </NScrollbar>
        <template #footer>
          <NSpace justify="end">
            <NButton @click="closeMergeModal">{{ $t("common.cancel") }}</NButton>
            <NButton type="primary" :loading="merging" :disabled="!mergeTargetId" @click="handleMerge">
              确认合并
            </NButton>
          </NSpace>
        </template>
      </NModal>

      <!-- 批量改标题弹窗 -->
      <NModal v-model:show="batchTitleVisible" preset="card" title="批量改标题" style="width: 480px" :bordered="false">
        <div class="mb-12px text-13px text-#999">
          输入标题模板，{n} 表示章节序号，例如：第{n}章
        </div>
        <NInput v-model:value="batchTitleTemplate" placeholder="请输入标题模板" />
        <template #footer>
          <NSpace justify="end">
            <NButton @click="closeBatchTitleModal">{{ $t("common.cancel") }}</NButton>
            <NButton type="primary" @click="handleBatchTitle">
              {{ $t("common.confirm") }}
            </NButton>
          </NSpace>
        </template>
      </NModal>

      <!-- 重新识别章节 -->
      <BookReparseModal
        v-model:visible="reparseVisible"
        :book-id="props.bookId"
        :book-title="props.bookTitle"
        @reparsed="handleReparsed"
      />
    </NDrawerContent>
  </NDrawer>
</template>
