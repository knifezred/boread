<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { NButton, NCard, NTag, NSpace, NSpin, NPagination } from 'naive-ui'
import BookCard from "@/components/book-card.vue"
import { fetchGetBook, fetchGetChapterList } from "@/service/api"

defineOptions({ name: 'BookDetail' })

const route = useRoute()
const router = useRouter()
const bookId = route.params.id as string

const bookInfo = ref<Api.SystemManage.Book>({
  id: 0,
  title: '',
  author: '',
  cover: '',
  createBy: '',
  createTime: '',
  updateBy: '',
  updateTime: '',
  status: '1',
  intro: null,
  categoryId: null,
  language: '',
  serialStatus: '1',
  visibility: '1',
  primaryFileId: null,
  totalChapters: 0,
  totalWords: 0,
  aggregateStatus: '1',
  avgRating: 0,
  ratingCount: 0,
  deptId: null
})

const loading = ref(false)

const relatedBooks = ref<Api.SystemManage.Book[]>([])
const chapters = ref<Api.SystemManage.BookChapter[]>([])
const chapterTotal = ref(0)
const chapterPage = ref(1)
const chapterSize = ref(100)
const chapterLoading = ref(false)
const sortAsc = ref(true)

const sortedChapters = computed(() => {
  const list = [...chapters.value]
  list.sort((a, b) => sortAsc.value ? a.chapterNo - b.chapterNo : b.chapterNo - a.chapterNo)
  return list
})

const latestChapter = ref('')
const activeSection = ref('section-info')
const observer = ref<IntersectionObserver | null>(null)

function toggleSort() {
  sortAsc.value = !sortAsc.value
}

function scrollToSection(id: string) {
  const el = document.getElementById(id)
  if (el) {
    el.scrollIntoView({ behavior: 'smooth', block: 'start' })
  }
}

function goToReader() {
  router.push({ name: 'book-reader', query: { id: bookInfo.value.id } })
}

async function loadChapters(page = 1) {
  chapterLoading.value = true
  chapterPage.value = page
  try {
    const { data } = await fetchGetChapterList({
      bookId: Number(bookId),
      current: page,
      size: chapterSize.value
    })
    if (data) {
      chapters.value = data.records || []
      chapterTotal.value = data.total || 0
    }
  } catch (err) {
    console.error('加载章节列表失败:', err)
  } finally {
    chapterLoading.value = false
  }
}

onMounted(async () => {
  loading.value = true
  try {
    const { data: bookData } = await fetchGetBook(bookId)
    if (bookData) {
      bookInfo.value = bookData
      await loadChapters(1)

      if (bookData.totalChapters > 0) {
        try {
          const lastPage = Math.ceil(bookData.totalChapters / chapterSize.value)
          const { data: lastPageData } = await fetchGetChapterList({
            bookId: Number(bookId),
            current: lastPage,
            size: chapterSize.value
          })
          if (lastPageData?.records?.length) {
            latestChapter.value = lastPageData.records[lastPageData.records.length - 1].title
          }
        } catch {
          // 查询最新章节失败不影响主流程
        }
      }
    }
  } catch (err) {
    console.error('加载书籍详情失败:', err)
  } finally {
    loading.value = false
  }

  // 等待 DOM 渲染完成后初始化 IntersectionObserver
  await nextTick()
  const sections = ['section-info', 'section-catalog']
  const els = sections.map(id => document.getElementById(id)).filter(Boolean) as HTMLElement[]
  if (els.length) {
    observer.value = new IntersectionObserver(
      (entries) => {
        for (const entry of entries) {
          if (entry.isIntersecting) {
            activeSection.value = entry.target.id
          }
        }
      },
      { rootMargin: '-80px 0px -50% 0px' }
    )
    els.forEach(el => observer.value!.observe(el))
  }
})

onBeforeUnmount(() => {
  observer.value?.disconnect()
})
</script>

<template>
  <div class="min-h-screen bg-gray-100 px-10 py-5 font-sans">
    <div class="flex items-center gap-2 text-sm text-gray-400 mb-5">
      <span>首页</span>
      <SvgIcon icon="solar:alt-arrow-right-linear" size="14" />
      <span>仙侠频道</span>
      <SvgIcon icon="solar:alt-arrow-right-linear" size="14" />
      <span>修真文明</span>
      <SvgIcon icon="solar:alt-arrow-right-linear" size="14" />
      <span class="text-gray-700">{{ bookInfo.title }}</span>
    </div>

    <div class="flex gap-6 max-w-1600px mx-auto items-start">
      <aside class="sticky top-5 w-120px shrink-0 flex flex-col gap-2">
        <div
          class="px-4 py-3 text-[15px] cursor-pointer rd-1 relative transition-all duration-200 hover:bg-white"
          :class="activeSection === 'section-info' ? 'bg-white text-primary font-medium' : 'text-gray-500 hover:text-gray-700'"
          @click="scrollToSection('section-info')">
          <span class="absolute left-0 top-0 bottom-0 w-[3px] rd-0 rd-r-2px" :class="activeSection === 'section-info' ? 'bg-primary' : 'bg-transparent'"></span>
          作品信息
        </div>
        <div
          class="px-4 py-3 text-[15px] cursor-pointer rd-1 relative transition-all duration-200 hover:bg-white"
          :class="activeSection === 'section-catalog' ? 'bg-white text-primary font-medium' : 'text-gray-500 hover:text-gray-700'"
          @click="scrollToSection('section-catalog')">
          <span class="absolute left-0 top-0 bottom-0 w-[3px] rd-0 rd-r-2px" :class="activeSection === 'section-catalog' ? 'bg-primary' : 'bg-transparent'"></span>
          目录
        </div>
      </aside>

      <main class="flex-1 max-w-900px flex flex-col gap-4">
        <section id="section-info">
          <NCard class="rd-12px shadow-sm" :bordered="false" size="huge">
            <div class="flex gap-8">
              <div class="shrink-0 w-160px">
                <NSpin v-if="loading" show />
                <BookCard
                  v-else-if="bookInfo"
                  :book="bookInfo"
                  :show-status-tag="true"
                  class="only-cover" />
              </div>

              <div class="flex-1" v-if="bookInfo">
                <h1 class="text-4xl font-bold mb-4 text-gray-900">{{ bookInfo.title }}</h1>
                <div class="flex gap-6 text-sm text-gray-500 mb-3">
                  <span>作者: {{ bookInfo.author }}</span>
                  <span>更新时间: {{ bookInfo.updateTime }}</span>
                </div>
                <div class="text-sm text-gray-500 mb-3">
                  <span>最新章节: </span>
                  <span class="text-primary no-underline">{{ latestChapter || '加载中...' }}</span>
                </div>
                <div class="flex gap-2 mb-4" v-if="bookInfo.tags?.length">
                  <NTag v-for="tag in bookInfo.tags" :key="tag.id" size="small" round type="info" ghost>
                    {{ tag.tagName }}
                  </NTag>
                </div>
                <p class="text-sm text-gray-600 leading-relaxed mb-5">{{ bookInfo.intro ? bookInfo.intro.slice(0, 60) + '...' : '' }}</p>
                <div class="flex gap-8 mb-6">
                  <div class="text-center">
                    <span class="block text-xl font-semibold text-gray-900">{{ bookInfo.totalWords }}</span>
                    <span class="text-xs text-gray-400">字</span>
                  </div>
                  <div class="text-center">
                    <span class="block text-xl font-semibold text-gray-900">{{ bookInfo.totalChapters }}</span>
                    <span class="text-xs text-gray-400">章节</span>
                  </div>
                  <div class="text-center">
                    <span class="block text-xl font-semibold text-gray-900">{{ bookInfo.avgRating }}</span>
                    <span class="text-xs text-gray-400">评分</span>
                  </div>
                </div>
                <div class="flex items-center">
                  <NSpace size="medium">
                    <NButton size="large" ghost type="primary" @click="goToReader">
                      免费试读
                    </NButton>
                    <NButton size="large" ghost type="primary">
                      加入书架
                    </NButton>
                  </NSpace>
                </div>
              </div>
            </div>
          </NCard>
        </section>

        <section id="section-catalog">
          <NCard class="rd-12px shadow-sm" :bordered="false" size="huge">
            <template #header>
              <div class="flex items-center gap-3 w-full">
                <span class="text-xl font-semibold text-gray-900">目录</span>
                <span class="text-sm text-gray-400 font-normal">共 {{ chapterTotal }} 章</span>
                <div
                  class="flex items-center gap-1 text-xs text-gray-400 cursor-pointer px-2 py-1 rd-1 transition-all duration-200 ml-auto hover:bg-gray-100 hover:text-gray-700"
                  @click="toggleSort">
                  <SvgIcon :icon="sortAsc ? 'solar:sort-from-top-linear' : 'solar:sort-from-bottom-linear'" size="16" />
                  <span>{{ sortAsc ? '正序' : '倒序' }}</span>
                </div>
              </div>
            </template>

            <div class="flex items-center gap-3 px-5 py-4 bg-amber-50 rd-2 mb-5" v-if="latestChapter">
              <span class="text-amber-600 font-medium">最新</span>
              <span class="font-medium text-gray-900 flex-1 overflow-hidden text-ellipsis whitespace-nowrap">{{ latestChapter }}</span>
            </div>

            <NSpin :show="chapterLoading">
              <div class="grid grid-cols-3 gap-3 gap-x-6">
                <div
                  v-for="ch in sortedChapters"
                  :key="ch.id"
                  class="flex items-center gap-1.5 px-3 py-2 rd-1.5 cursor-pointer transition-colors duration-200 text-sm text-gray-700 hover:bg-gray-100"
                  @click="goToReader()">
                  <span class="flex-1 overflow-hidden text-ellipsis whitespace-nowrap">{{ ch.title }}</span>
                </div>
              </div>
            </NSpin>

            <div v-if="!chapters.length && !chapterLoading" class="py-10 text-center text-sm text-gray-400">
              暂无章节信息
            </div>

            <div class="flex justify-center mt-8" v-if="chapterTotal > chapterSize">
              <NPagination
                :page="chapterPage"
                :page-size="chapterSize"
                :item-count="chapterTotal"
                @update:page="loadChapters" />
            </div>
          </NCard>
        </section>
      </main>

      <aside class="max-xl:hidden w-300px shrink-0 flex flex-col gap-4">
        <NCard class="rd-12px shadow-sm" :bordered="false" size="small">
          <template #header>
            <div class="flex justify-between items-center w-full">
              <span>作品阅历</span>
              <a href="#" class="text-xs text-gray-400 no-underline hover:text-primary">更多</a>
            </div>
          </template>
          <div class="flex flex-col gap-1 py-3 border-b border-gray-100 last:border-b-0">
            <span class="text-xs text-gray-400">2026-05-27</span>
            <span class="text-sm text-gray-700">累积获得五千个收藏</span>
          </div>
        </NCard>

        <NCard class="rd-12px shadow-sm" :bordered="false" size="small">
          <template #header>
            <div class="flex justify-between items-center w-full">
              <span>相似小说推荐</span>
              <a href="#" class="text-xs text-gray-400 no-underline hover:text-primary">更多</a>
            </div>
          </template>
          <div class="flex flex-col gap-4">
            <div v-for="book in relatedBooks" :key="book.id" class="flex gap-3">
              <div class="shrink-0 w-50px">
                <BookCard :book="book" class="small-book-card" />
              </div>
              <div class="flex-1 overflow-hidden">
                <div class="text-sm font-medium text-gray-900 mb-1 truncate">{{ book.title }}</div>
                <div class="text-xs text-gray-400 leading-relaxed line-clamp-2">{{ book.intro?.slice(0, 40) || '' }}</div>
              </div>
            </div>
          </div>
        </NCard>
      </aside>
    </div>
  </div>
</template>

<style scoped>
.only-cover .book-info {
  display: none;
}

.small-book-card .book-info {
  display: none;
}
</style>
