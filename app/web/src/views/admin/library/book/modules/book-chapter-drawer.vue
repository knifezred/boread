<script setup lang="tsx">
import { ref, computed, watch } from "vue"
import {
  NButton, NInput, NSelect, NSpace,
  NTag, NModal, NScrollbar, NCheckbox,
  NDrawer, NDrawerContent, NSpin, NVirtualList,
} from "naive-ui"
import { useBoolean } from "@sa/hooks"
import {
  fetchChapterList,
  fetchBatchUpdateChapterTitle, fetchUpdateChapterStatus,
  fetchMergeChapters, fetchFormatChapterNumbers,
} from "@/service/api"
import { $t } from "@/locales"
import BookReparseModal from "./book-reparse-modal.vue"

defineOptions({ name: "BookChapterDrawer" })

interface Props { bookId: number; bookTitle: string }
const props = defineProps<Props>()
const visible = defineModel<boolean>("visible", { default: false })

// ==================== 数据 ====================
const allChapters = ref<Api.BookManage.BookChapter[]>([])
const checkedRowKeys = ref<number[]>([])
const loading = ref(false)

interface VolumeGroup {
  volumeNo: number | null
  volumeTitle: string
  chapters: Api.BookManage.BookChapter[]
}

/** 展开状态，默认全部展开 */
const expandedVolumes = ref<Set<number | string>>(new Set())

function toggleVolume(volumeNo: number | null) {
  const key = volumeNo ?? "__main__"
  const next = new Set(expandedVolumes.value)
  if (next.has(key)) next.delete(key); else next.add(key)
  expandedVolumes.value = next
}

function isVolumeExpanded(volumeNo: number | null): boolean {
  return expandedVolumes.value.has(volumeNo ?? "__main__")
}

// ==================== 搜索 ====================
const searchTitle = ref("")
const searchStatus = ref("")

const filteredChapters = computed(() => {
  let list = allChapters.value
  if (searchTitle.value) {
    const kw = searchTitle.value.toLowerCase()
    list = list.filter(ch => ch.title.toLowerCase().includes(kw))
  }
  if (searchStatus.value) {
    list = list.filter(ch => ch.status === searchStatus.value)
  }
  return list
})

/** 按分卷分组后的目录树 */
const volumeGroups = computed<VolumeGroup[]>(() => {
  const groups: VolumeGroup[] = []
  let currentGroup: VolumeGroup | null = null
  for (const ch of filteredChapters.value) {
    if (!currentGroup || currentGroup.volumeNo !== ch.volumeNo) {
      currentGroup = {
        volumeNo: ch.volumeNo,
        volumeTitle: ch.volumeTitle || (ch.volumeNo ? `第${ch.volumeNo}卷` : "正文"),
        chapters: [],
      }
      groups.push(currentGroup)
    }
    currentGroup.chapters.push(ch)
  }
  return groups
})

/** 扁平化后的虚拟列表行 */
interface FlatItem {
  key: string
  _type: 'volume' | 'chapter'
  _volumeNo: number | null
  _volumeTitle: string
  _chapterCount: number
  _expanded: boolean
  data: Api.BookManage.BookChapter | null
}

const flattenedItems = computed<FlatItem[]>(() => {
  const items: FlatItem[] = []
  for (const group of volumeGroups.value) {
    const expanded = volumeGroups.value.length <= 1 || isVolumeExpanded(group.volumeNo)
    items.push({
      key: `vol-${group.volumeNo ?? '__main__'}`,
      _type: 'volume',
      _volumeNo: group.volumeNo,
      _volumeTitle: group.volumeTitle,
      _chapterCount: group.chapters.length,
      _expanded: expanded,
      data: null,
    })
    if (expanded) {
      for (const ch of group.chapters) {
        items.push({
          key: `ch-${ch.id}`,
          _type: 'chapter',
          _volumeNo: group.volumeNo,
          _volumeTitle: group.volumeTitle,
          _chapterCount: 0,
          _expanded: true,
          data: ch,
        })
      }
    }
  }
  return items
})

async function loadAllChapters() {
  loading.value = true
  const { data } = await fetchChapterList(props.bookId)
  if (data) {
    allChapters.value = data
    // 重置展开状态：全部展开
    const groups: VolumeGroup[] = []
    let currentGroup: VolumeGroup | null = null
    for (const ch of data) {
      if (!currentGroup || currentGroup.volumeNo !== ch.volumeNo) {
        currentGroup = {
          volumeNo: ch.volumeNo,
          volumeTitle: ch.volumeTitle || (ch.volumeNo ? `第${ch.volumeNo}卷` : "正文"),
          chapters: [],
        }
        groups.push(currentGroup)
      }
      currentGroup.chapters.push(ch)
    }
    expandedVolumes.value = new Set(groups.map(g => g.volumeNo ?? "__main__"))
  }
  loading.value = false
}

// ==================== 选中行 ====================
function isChecked(chapterId: number): boolean {
  return checkedRowKeys.value.includes(chapterId)
}

function toggleCheck(chapterId: number) {
  const idx = checkedRowKeys.value.indexOf(chapterId)
  if (idx >= 0) {
    checkedRowKeys.value.splice(idx, 1)
  } else {
    checkedRowKeys.value.push(chapterId)
  }
}

function isVolumeAllChecked(group: VolumeGroup): boolean {
  return group.chapters.length > 0 && group.chapters.every(ch => checkedRowKeys.value.includes(ch.id))
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
    for (const ch of allChapters.value) {
      if (checkedRowKeys.value.includes(ch.id)) (ch as any).status = status
    }
    checkedRowKeys.value = []
  }
}

// ==================== 合并章节 ====================
const { bool: mergeVisible, setTrue: openMergeModal, setFalse: closeMergeModal } = useBoolean()
const allChaptersForMerge = ref<Api.BookManage.BookChapter[]>([])
const mergeTargetId = ref<number | null>(null)
const merging = ref(false)

async function openMerge() {
  if (checkedRowKeys.value.length < 2) {
    window.$message?.warning("请至少选择2个章节进行合并")
    return
  }
  mergeTargetId.value = null
  // 直接用现有数据，避免额外请求
  allChaptersForMerge.value = allChapters.value
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
    loadAllChapters()
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
    loadAllChapters()
  }
}

// ==================== 重新识别 ====================
const reparseVisible = ref(false)

function handleReparsed() {
  reparseVisible.value = false
  loadAllChapters()
}

// ==================== 生命周期 ====================
watch(visible, (val) => {
  if (val) {
    searchTitle.value = ""
    searchStatus.value = ""
    checkedRowKeys.value = []
    loadAllChapters()
  }
})

function closeDrawer() {
  visible.value = false
}
</script>

<template>
  <NDrawer v-model:show="visible" display-directive="show" width="60%" native-scrollbar>
    <NDrawerContent
      :title="`${$t('page.admin.library.book.chapters')} - ${props.bookTitle}`"
      closable
      @after-leave="closeDrawer"
    >
      <div class="h-full flex flex-col overflow-hidden">
        <!-- 搜索栏 -->
        <div class="flex items-center gap-12px mb-16px shrink-0">
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
          <NButton size="small" type="primary" ghost @click="loadAllChapters">
            {{ $t("common.search") }}
          </NButton>
          <NButton size="small" @click="searchTitle = ''; searchStatus = ''; loadAllChapters()">
            {{ $t("common.reset") }}
          </NButton>
          <span class="text-xs text-gray-400 ml-auto">
            共 {{ allChapters.length }} 章
          </span>
        </div>

        <!-- 章节列表（虚拟滚动 + 分卷树） -->
        <div class="flex-1" style="position: relative; min-height: 0;">
          <!-- 加载中 -->
          <div v-if="loading" class="flex items-center justify-center h-full">
            <NSpin />
          </div>

          <!-- 空状态 -->
          <div
            v-else-if="allChapters.length === 0"
            class="flex items-center justify-center h-full text-sm text-gray-400"
          >
            暂无章节
          </div>

          <!-- 虚拟列表 -->
          <NVirtualList
            v-else
            :items="flattenedItems"
            :item-size="38"
            key-field="key"
            class="h-full"
          >
            <template #default="{ item }">
              <template v-if="item._type === 'volume'">
                <!-- @ts-ignore -->
                <div
                  class="flex items-center gap-2 px-3 py-2 cursor-pointer select-none text-xs font-medium uppercase tracking-wider bg-[#fafafa] dark:bg-[#1e1e1e] border-b border-[#f0f0f0] dark:border-[#333]"
                >
                  <NCheckbox
                    :checked="isVolumeAllChecked({ volumeNo: item._volumeNo, volumeTitle: item._volumeTitle, chapters: flattenedItems.filter(f => f._type === 'chapter' && f._volumeNo === item._volumeNo).map(f => (f as any).data) })"
                    :indeterminate="(() => { const chs = flattenedItems.filter(f => f._type === 'chapter' && f._volumeNo === item._volumeNo).map(f => (f as any).data as Api.BookManage.BookChapter); return chs.some(c => checkedRowKeys.includes(c.id)) && !chs.every(c => checkedRowKeys.includes(c.id)) })()"
                    @update:checked="() => { const chs = flattenedItems.filter(f => f._type === 'chapter' && f._volumeNo === item._volumeNo).map(f => (f as any).data as Api.BookManage.BookChapter); const allChecked = chs.every(c => checkedRowKeys.includes(c.id)); for (const c of chs) { const i = checkedRowKeys.indexOf(c.id); if (allChecked && i >= 0) checkedRowKeys.splice(i, 1); else if (!allChecked && i < 0) checkedRowKeys.push(c.id) } }"
                  />
                  <span
                    class="text-[10px] transition-transform duration-200"
                    :class="isVolumeExpanded(item._volumeNo) ? 'rotate-90' : ''"
                    @click="toggleVolume(item._volumeNo)"
                  >
                    ▸
                  </span>
                  <span
                    class="text-gray-500 dark:text-gray-400"
                    :class="isVolumeExpanded(item._volumeNo) ? 'text-primary' : ''"
                    @click="toggleVolume(item._volumeNo)"
                  >
                    {{ item._volumeTitle }}
                  </span>
                  <span class="text-[10px] opacity-50">({{ item._chapterCount }}章)</span>
                </div>
              </template>
              <template v-else>
                <!-- @ts-ignore -->
                <div
                  class="flex items-center gap-3 px-4 py-3 cursor-pointer transition-colors duration-200 text-sm border-b border-[#f5f5f5] dark:border-[#2a2a2a] hover:bg-[#fafafa] dark:hover:bg-[#252525]"
                >
                  <!-- @ts-ignore -->
                  <NCheckbox
                    :checked="isChecked((item.data as any).id)"
                    @update:checked="toggleCheck((item.data as any).id)"
                  />
                  <!-- @ts-ignore -->
                  <span class="text-xs shrink-0 w-8 text-right text-gray-400">
                    {{ (item.data as any).chapterNo }}
                  </span>
                  <!-- @ts-ignore -->
                  <span class="flex-1 overflow-hidden text-ellipsis whitespace-nowrap">
                    {{ (item.data as any).title }}
                  </span>
                </div>
              </template>
            </template>
          </NVirtualList>
        </div>
      </div>

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

      <!-- 合并弹窗 -->
      <NModal v-model:show="mergeVisible" preset="card" title="合并章节" style="width: 560px" :bordered="false">
        <div class="mb-12px text-14px">
          请选择目标章节，选中的源章节将合并到目标章节中
        </div>
        <NScrollbar style="max-height: 360px">
          <div
            v-for="ch in allChaptersForMerge" :key="ch.id"
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
