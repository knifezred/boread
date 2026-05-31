<script setup lang="ts">
import { ref, computed, onMounted } from "vue"
import { NButton, NInput, NScrollbar, NSpin, NPagination } from "naive-ui"
import { fetchGetBookList, fetchGetHotCategoryList, fetchGetTagList } from "@/service/api"
import { useRouter } from "vue-router"
import { useDictItems } from "@/hooks/business/dict"
import { formatWordCount } from "@/utils/book"
import { $t } from "@/locales"

defineOptions({ name: "HomePage" })

const router = useRouter()

const { options: serialStatusOptions, labelMap: serialStatusLabelMap } = useDictItems("book_serial_status")
const { options: wordCountOptions } = useDictItems("book_total_words")

/** 从分类树中提取的筛选选项 */
const categoryOptions = ref<{ label: string; value: string }[]>([])

/** 加载热门分类列表，不需要子父关系 */
async function loadCategoryOptions() {
  const { data } = await fetchGetHotCategoryList()
  if (!data) return
  categoryOptions.value = [
    { label: $t("page.book.home.all"), value: "" },
    ...data.map(n => ({ label: n.categoryName, value: String(n.id) }))
  ]
}

/** 标签列表选项 */
const tagOptions = ref<{ label: string; value: string }[]>([])

async function loadTagOptions() {
  const { data } = await fetchGetTagList({
    size: 15,
    current: 1
  })
  if (!data) return
  tagOptions.value = [
    { label: $t("page.book.home.all"), value: "" },
    ...data.records.map(n => ({ label: n.tagName, value: String(n.id) }))
  ]
}

/** 搜索关键词 */
const searchKeyword = ref("")
/** 书籍列表 */
const books = ref<Api.BookManage.Book[]>([])
/** 加载状态 */
const loading = ref(false)
/** 总数 */
const total = ref(0)
/** 当前页 */
const current = ref(1)
/** 每页条数 */
const size = ref(20)

/**
 * 筛选配置（对齐后端 BookSearch 字段）
 */
const filterConfig = computed(() => ({
  categories: categoryOptions.value.length > 1 ? categoryOptions.value : [
    { label: $t("page.book.home.all"), value: "" }
  ],
  serialStatus: [
    { label: $t("page.book.home.all"), value: "" },
    ...serialStatusOptions.value.map(opt => ({ label: opt.label as string, value: opt.value as string }))
  ],
  wordCount: [
    { label: $t("page.book.home.all"), value: "" },
    ...wordCountOptions.value.map(opt => ({ label: opt.label as string, value: opt.value as string }))
  ],
  tags: tagOptions.value.length > 1 ? tagOptions.value : [
    { label: $t("page.book.home.all"), value: "" }
  ],
  updateTime: [
    { label: $t("page.book.home.all"), value: "" },
    { label: $t("page.book.home.oneWeek"), value: "7d" },
    { label: $t("page.book.home.oneMonth"), value: "30d" },
    { label: $t("page.book.home.threeMonths"), value: "90d" },
    { label: $t("page.book.home.oneYear"), value: "1y" }
  ],
  sortOptions: [
    { label: $t("page.book.home.sortPopular"), value: "popular" },
    { label: $t("page.book.home.sortCollect"), value: "collect" },
    { label: $t("page.book.home.sortWord"), value: "word" },
    { label: $t("page.book.home.sortVote"), value: "vote" },
    { label: $t("page.book.home.sortMonthly"), value: "monthly" }
  ]
}))

/** 筛选条件 */
const filters = ref<Api.BookManage.BookFilterParams>({
  categoryId: 0,
  serialStatus: "",
  wordCount: "",
  tagId: "",
  updateTime: "",
  title: "",
  sortBy: "",
  sortOrder: ""
})

onMounted(async () => {
  await Promise.all([loadCategoryOptions(), loadTagOptions()])
  loadBooks()
})

/** 解析字数范围，支持 "10-20" 区间或 "500+" 开头区间 */
function parseWordCountRange(wordCount: string) {
  if (!wordCount) return { minWords: null, maxWords: null }
  if (wordCount.endsWith("+")) {
    return { minWords: Number(wordCount.slice(0, -1)), maxWords: null }
  }
  const [min, max] = wordCount.split("-")
  return { minWords: Number(min), maxWords: Number(max) }
}

/** 解析更新时间范围 */
function parseUpdateTimeRange(updateTime: string) {
  if (!updateTime) return { updateTimeFrom: null, updateTimeTo: null }
  const [from, to] = updateTime.split("-")
  return { updateTimeFrom: Number(from), updateTimeTo: Number(to) }
}

/**
 * 加载书籍列表
 * @param page 页码
 */
async function loadBooks(page = 1) {
  loading.value = true
  current.value = page
  const { minWords, maxWords } = parseWordCountRange(filters.value.wordCount)
  const { updateTimeFrom, updateTimeTo } = parseUpdateTimeRange(filters.value.updateTime)
  const params: Record<string, unknown> = {
    current: page,
    size: size.value,
    categoryId: filters.value.categoryId || null,
    serialStatus: filters.value.serialStatus || null,
    tagId: filters.value.tagId || null,
    minWords,
    maxWords,
    updateTimeFrom,
    updateTimeTo
  }
  if (searchKeyword.value.trim()) {
    params.title = searchKeyword.value.trim()
  }
  const { data } = await fetchGetBookList(params)
  if (data) {
    books.value = data.records || []
    total.value = data.total || 0
  }
  loading.value = false
}

/** 搜索 */
function handleSearch() {
  loadBooks(1)
}

/** 筛选条件变更 */
function handleFilterChange() {
  loadBooks(1)
}

/** 排序切换 */
function handleSortChange(_value: string) {
  loadBooks(1)
}

/**
 * 点击书籍卡片，跳转到详情页
 * @param book 书籍
 */
function showChapters(book: Api.BookManage.Book) {
  router.push(`/book-detail/${book.id}`)
}
</script>

<template>
  <div class="flex gap-4 overflow-hidden lt-sm:flex-col h-full bg-layout">
    <!-- ============ 左侧筛选栏 ============ -->
    <aside class="w-240px shrink-0 lt-sm:w-full lt-sm:max-h-60 lt-sm:overflow-y-auto">
      <div class="card-wrapper p-4 bg-container h-full lt-sm:h-auto">
        <NScrollbar class="h-full">
          <BookFilter
            :config="{
              categories: filterConfig.categories,
              serialStatus: filterConfig.serialStatus,
              wordCount: filterConfig.wordCount,
              tags: filterConfig.tags,
              updateTime: filterConfig.updateTime
            }"
            :model-value="filters"
            @change="handleFilterChange"
            @update:model-value="filters = $event" />
        </NScrollbar>
      </div>
    </aside>

    <!-- ============ 右侧主内容 ============ -->
    <div class="flex-1 flex flex-col overflow-hidden min-w-0">
      <!-- 顶部搜索栏 -->
      <div class="card-wrapper p-4 bg-container mb-4">
        <div class="flex items-center gap-3 flex-wrap">
          <div class="w-280px lt-sm:w-full">
            <NInput
              v-model:value="searchKeyword"
              :placeholder="$t('page.book.home.searchPlaceholder')"
              clearable
              round
              @keyup.enter="handleSearch" />
          </div>
          <NButton type="primary" round @click="handleSearch">{{ $t("common.search") }}</NButton>
          <NButton round @click="router.push('/book-reader')">{{ $t("page.book.home.importBooks") }}</NButton>
        </div>
      </div>

      <!-- 排序栏 -->
      <div class="card-wrapper px-4 py-3 bg-container mb-4 flex items-center justify-between flex-wrap gap-2">
        <div class="flex items-center gap-6">
          <span
            v-for="item in filterConfig.sortOptions"
            :key="item.value"
            class="text-14px cursor-pointer transition-colors duration-200 select-none"
            :class="filters.sortBy === item.value
                ? 'text-primary font-600'
                : 'text-gray-400 hover:text-primary'
              "
            @click="handleSortChange(item.value)">
            {{ item.label }}
          </span>
        </div>
        <div class="text-13px text-gray-400">
          {{ $t("page.book.home.relatedWorks", { total }) }}
        </div>
      </div>

      <!-- 内容区域 -->
      <div class="flex-1 overflow-y-auto">
        <NSpin :show="loading">
          <!-- 空状态 -->
          <div v-if="books.length === 0 && !loading" class="flex flex-col items-center justify-center py-30">
            <div class="mb-5 opacity-60">
              <svg width="120" height="100" viewBox="0 0 120 100" fill="none" xmlns="http://www.w3.org/2000/svg">
                <path d="M60 10C45 10 35 20 35 35C35 50 45 60 60 60C75 60 85 50 85 35C85 20 75 10 60 10Z"
                  fill="#E6F4FF" />
                <path
                  d="M55 30C55 27.2386 57.2386 25 60 25C62.7614 25 65 27.2386 65 30C65 32.7614 62.7614 35 60 35C57.2386 35 55 32.7614 55 30Z"
                  fill="#B3D8FF" />
                <path
                  d="M40 65C40 59.4772 44.4772 55 50 55H70C75.5228 55 80 59.4772 80 65V85C80 87.7614 77.7614 90 75 90H45C42.2386 90 40 87.7614 40 85V65Z"
                  fill="#8EC2FF" />
                <path
                  d="M45 75H75C76.6569 75 78 76.3431 78 78C78 79.6569 76.6569 81 75 81H45C43.3431 81 42 79.6569 42 78C42 76.3431 43.3431 75 45 75Z"
                  fill="#E6F4FF" />
              </svg>
            </div>
            <p class="text-16px text-gray-400">{{ $t("page.book.home.noContent") }}</p>
          </div>

          <!-- 书籍网格 -->
          <div v-else class="grid grid-cols-3 gap-6 lt-sm:grid-cols-1">
            <div
              v-for="book in books"
              :key="book.id"
              class="flex gap-4 p-4 bg-container rd-2 cursor-pointer transition-all duration-200 shadow-sm hover:shadow-md hover:-translate-y-1"
              @click="showChapters(book)">
              <!-- 封面 -->
              <div class="w-100px shrink-0">
                <BookCard
                  :book="book" />
              </div>

              <!-- 书籍信息 -->
              <div class="flex-1 flex flex-col gap-2 min-w-12">
                <h3 class="text-lg m-0 truncate">
                  {{ book.title }}
                </h3>
                <div class="flex items-center gap-2 text-gray-400 text-xs">
                  <span>{{ book.author }}</span>
                  <span class="mx-1">|</span>
                  <span>{{ book.categoryName || $t("page.book.home.uncategorized") }}</span>
                  <span class="text-gray-400 mx-1">|</span>
                  <span>{{ serialStatusLabelMap[book.serialStatus] }}</span>
                </div>
                <p class="flex-1 text-sm text-gray-500 m-0 leading-6 line-clamp-3 whitespace-pre-line">
                  {{ book.intro || $t("page.book.home.noIntro") }}
                </p>
                <div class="flex items-center gap-2 justify-left text-xs text-gray-400">
                  <span>{{ formatWordCount(book.totalWords) }}</span>
                  <span class="text-gray-400 mx-1">|</span>
                  <span>{{ $t("page.book.home.latestChapter") }}</span>
                </div>
              </div>
            </div>
          </div>
        </NSpin>

        <!-- 分页 -->
        <div v-if="books.length > 0" class="flex justify-center mt-12 pb-4">
          <NPagination
            :page="current"
            :page-size="size"
            :item-count="total"
            :page-sizes="[12, 24, 48, 60]"
            show-size-picker
            @update:page="loadBooks"
            @update:page-size="(s) => { size = s; loadBooks(1) }" />
        </div>
      </div>
    </div>
  </div>
</template>
