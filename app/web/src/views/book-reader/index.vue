<script setup lang="ts">
import { ref, computed, onMounted, nextTick, onBeforeUnmount } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { NButton, NSpin } from 'naive-ui'
import { fetchGetBook, fetchGetChapterContent } from "@/service/api"
import { useBoolean } from '@sa/hooks'
import { $t } from "@/locales"
import { formatWordCount } from '@/utils/book'
import CatalogModal from './modules/catalog-modal.vue'

defineOptions({ name: 'BookReader' })

const route = useRoute()
const router = useRouter()
const bookId = (route.params.id || route.query.id) as string

const darkMode = ref(false)
const fontSize = ref(18)
const lineHeight = ref(1.8)
const contentLoading = ref(true)
const bookLoading = ref(true)
const shiftingMore = ref(false)

const bookInfo = ref<Api.BookManage.Book | null>(null)

interface ChapterItem {
  chapterNo: number
  title: string
  content: string
  wordCount: string
  updateTime: string
}

const chapters = ref<ChapterItem[]>([])

const sortedChapters = computed(() => {
  return [...chapters.value].sort((a, b) => a.chapterNo - b.chapterNo)
})

const hasPrev = computed(() => {
  if (!chapters.value.length) return false
  const minNo = Math.min(...chapters.value.map(c => c.chapterNo))
  return minNo > 1
})

const hasNext = computed(() => {
  if (!chapters.value.length || !bookInfo.value) return false
  const maxNo = Math.max(...chapters.value.map(c => c.chapterNo))
  return maxNo < bookInfo.value.totalChapters
})

const { bool: catalogVisible, setTrue: openModel } = useBoolean()

const pageBg = computed(() => darkMode.value ? '#222' : '#fff')
const readerBg = computed(() => darkMode.value ? '#1a1a1a' : '#f2eede')
const textColor = computed(() => darkMode.value ? '#e0e0e0' : '#333')
const readerBorder = computed(() => darkMode.value ? 'rgba(255,255,255,0.06)' : 'rgba(0,0,0,0.04)')

let sentinelObserver: IntersectionObserver | null = null

async function loadOneChapter(no: number): Promise<ChapterItem | null> {
  try {
    const { data } = await fetchGetChapterContent(bookId, String(no))
    if (!data) return null
    return {
      chapterNo: no,
      title: data.title,
      content: data.content || '',
      wordCount: `${data.wordCount}${$t("page.book.detail.words")}`,
      updateTime: data.updateTime || '',
    }
  } catch {
    return null
  }
}

async function initWindow(centerNo: number) {
  chapters.value = []
  contentLoading.value = true
  const needed = [centerNo - 1, centerNo, centerNo + 1].filter(n => n >= 1)
  const results = await Promise.all(needed.map(n => loadOneChapter(n)))
  chapters.value = results.filter(Boolean) as ChapterItem[]
  chapters.value.sort((a, b) => a.chapterNo - b.chapterNo)
  router.replace({ query: { ...route.query, chapterNo: String(centerNo) } })
  contentLoading.value = false
  await nextTick()
  observeSentinel()
}

async function shiftForward() {
  if (shiftingMore.value) return
  shiftingMore.value = true
  const maxNo = Math.max(...chapters.value.map(c => c.chapterNo))
  const nextNo = maxNo + 1
  const next = await loadOneChapter(nextNo)
  shiftingMore.value = false
  if (!next) return

  const minNo = Math.min(...chapters.value.map(c => c.chapterNo))
  chapters.value = chapters.value.filter(c => c.chapterNo !== minNo)
  chapters.value.push(next)
  chapters.value.sort((a, b) => a.chapterNo - b.chapterNo)

  const centerNo = sortedChapters.value[1]?.chapterNo || nextNo
  router.replace({ query: { ...route.query, chapterNo: String(centerNo) } })

  await nextTick()
  observeSentinel()
}

async function shiftBackward() {
  if (shiftingMore.value) return
  shiftingMore.value = true
  const minNo = Math.min(...chapters.value.map(c => c.chapterNo))
  const prevNo = minNo - 1
  if (prevNo < 1) {
    shiftingMore.value = false
    return
  }
  const prev = await loadOneChapter(prevNo)
  shiftingMore.value = false
  if (!prev) return

  const maxNo = Math.max(...chapters.value.map(c => c.chapterNo))
  chapters.value = chapters.value.filter(c => c.chapterNo !== maxNo)
  chapters.value.unshift(prev)
  chapters.value.sort((a, b) => a.chapterNo - b.chapterNo)

  const centerNo = sortedChapters.value[1]?.chapterNo || prevNo
  router.replace({ query: { ...route.query, chapterNo: String(centerNo) } })

  await nextTick()
  observeSentinel()
}

function scrollToChapter(no: number) {
  const el = document.getElementById(`chapter-${no}`)
  if (el) {
    el.scrollIntoView({ behavior: 'smooth', block: 'start' })
  }
}

function prevChapter() {
  if (!hasPrev.value) return
  shiftBackward().then(() => {
    const sorted = [...chapters.value].sort((a, b) => a.chapterNo - b.chapterNo)
    scrollToChapter(sorted[1]?.chapterNo)
  })
}

function nextChapter() {
  if (!hasNext.value) return
  shiftForward().then(() => {
    const sorted = [...chapters.value].sort((a, b) => a.chapterNo - b.chapterNo)
    scrollToChapter(sorted[1]?.chapterNo)
  })
}

function goToChapter(no: number) {
  initWindow(no).then(() => {
    scrollToChapter(no)
  })
}

function observeSentinel() {
  sentinelObserver?.disconnect()
  const sentinel = document.getElementById('chapter-sentinel')
  if (!sentinel) return
  sentinelObserver = new IntersectionObserver((entries) => {
    if (entries[0].isIntersecting && !shiftingMore.value && hasNext.value) {
      shiftForward()
    }
  }, { rootMargin: '300px' })
  sentinelObserver.observe(sentinel)
}

function handleOpenCatalog() {
  openModel()
}

onMounted(async () => {
  if (!bookId) return
  bookLoading.value = true
  try {
    const { data: bookData } = await fetchGetBook(bookId)
    if (bookData) {
      bookInfo.value = bookData
    }
  } catch (err) {
    console.error('加载书籍信息失败:', err)
  } finally {
    bookLoading.value = false
  }

  const startChapter = Number(route.query.chapterNo) || 1
  await initWindow(startChapter)
})

onBeforeUnmount(() => {
  sentinelObserver?.disconnect()
})
</script>

<template>
  <div class="min-h-screen font-sans relative transition-colors duration-300"
    :style="{ backgroundColor: pageBg, color: textColor }">
    <main class="flex justify-center bg-[#F4F0E7]">
      <div class="flex gap-4 items-start w-full max-lg:max-w-full max-w-800px">
        <!-- 阅读区 -->
        <div
          class="w-full lg:px-8 px-4 lg:py-6 pt-14 pb-16 rd-1 min-h-100vh flex flex-col transition-colors duration-300 max-lg:min-h-[calc(100vh-60px)] noise-bg"
          :style="{ backgroundColor: readerBg }">

          <NSpin :show="contentLoading || bookLoading">
            <template v-if="sortedChapters.length">
              <!-- 书籍详情（仅首个窗口包含第1章时显示） -->
              <div v-if="sortedChapters[0].chapterNo === 1 && bookInfo" class="mb-10 pb-10 border-b text-center"
                :style="{ borderColor: readerBorder }">
                <div
                  v-if="bookInfo.cover"
                  class="w-130px h-180px rd-1 shadow-lg overflow-hidden mx-auto mb-6 mt-16">
                  <img :src="bookInfo.cover" :alt="bookInfo.title" class="w-full h-full object-cover" />
                </div>
                <h1 class="lg:text-3xl text-2xl font-bold mb-2 text-inherit">{{ bookInfo.title }}</h1>
                <p class="text-sm mb-8" :class="darkMode ? 'text-gray-400' : 'text-gray-500'">{{ bookInfo.author }} {{ $t("page.book.reader.writtenBy") }}</p>
                <div class="flex justify-center gap-6 lg:gap-16 mb-8">
                  <div class="flex flex-col gap-1">
                    <span class="text-base font-medium">{{ bookInfo.categoryName || $t("page.book.reader.unknown") }}</span>
                    <span class="text-xs" :class="darkMode ? 'text-gray-400' : 'text-gray-500'">{{ $t("page.book.reader.type") }}</span>
                  </div>
                  <div class="flex flex-col gap-1">
                    <span class="text-base font-medium">{{ bookInfo.serialStatus === '1' ? $t("page.book.reader.ongoing") : $t("page.book.reader.finished") }}</span>
                    <span class="text-xs" :class="darkMode ? 'text-gray-400' : 'text-gray-500'">{{ $t("page.book.reader.status") }}</span>
                  </div>
                  <div class="flex flex-col gap-1">
                    <span class="text-base font-medium">{{ formatWordCount(bookInfo.totalWords || 0) }}</span>
                    <span class="text-xs" :class="darkMode ? 'text-gray-400' : 'text-gray-500'">{{ $t("page.book.reader.words") }}</span>
                  </div>
                </div>
                <div v-if="bookInfo.intro" class="mb-6 max-w-2xl lg:mx-18 text-align-left">
                  <p v-for="(paragraph, index) in bookInfo.intro.split('\n').filter(p => p.trim())" :key="index"
                    class="text-sm leading-relaxed mb-2 last:mb-0" :class="darkMode ? 'text-gray-400' : 'text-gray-500'">
                    {{ paragraph.trim() }}
                  </p>
                </div>
              </div>

              <!-- 连续章节内容 -->
              <div class="flex-1">
                <template v-for="ch in sortedChapters" :key="ch.chapterNo">
                  <div :id="`chapter-${ch.chapterNo}`">
                    <div class="text-center mb-8 lg:mb-10 pt-2">
                      <h1 class="lg:text-2xl text-xl font-bold mb-4 text-inherit">{{ ch.title }}</h1>
                      <!-- <div class="flex justify-center gap-4 text-xs" :class="darkMode ? 'text-gray-500' : 'text-gray-400'">
                        <span>{{ $t("page.book.reader.words") }}: {{ ch.wordCount || '0' }}</span>
                      </div> -->
                    </div>
                    <div class="mb-12 text-inherit leading-relaxed"
                      :style="{ fontSize: `${fontSize}px`, lineHeight: lineHeight }">
                      <template v-for="(paragraph, index) in ch.content.split('\n').filter(p => p.trim())" :key="index">
                        <p class="mb-5 text-indent-2em tracking-normal">
                          {{ paragraph.trim() }}
                        </p>
                      </template>
                    </div>
                  </div>
                </template>
              </div>

              <!-- 滚动加载哨兵 -->
              <div id="chapter-sentinel" class="h-1"></div>

              <!-- 加载更多指示 -->
              <div v-if="shiftingMore" class="py-6 text-center text-sm"
                :class="darkMode ? 'text-gray-500' : 'text-gray-400'">
                <NSpin :show="true" />
              </div>

              <!-- 没有更多章节 -->
              <div v-if="!hasNext && sortedChapters.length > 0" class="py-12 text-center text-sm"
                :class="darkMode ? 'text-gray-500' : 'text-gray-400'">
                {{ $t("page.book.reader.noMoreChapters") }}
              </div>
            </template>

            <div v-else-if="!contentLoading && !bookLoading" class="py-20 text-center text-sm"
              :class="darkMode ? 'text-gray-500' : 'text-gray-400'">
              {{ $t("page.book.reader.noContent") }}
            </div>
          </NSpin>

          <!-- 桌面端底部导航 -->
          <div class="flex justify-center gap-6 pt-6 border-t shrink-0 max-lg:hidden"
            :style="{ borderColor: readerBorder }">
            <NButton :disabled="!hasPrev" text size="large" :style="{ color: textColor }" @click="prevChapter">
              <template #icon>
                <SvgIcon icon="solar:arrow-left-linear" size="18" />
              </template>
              {{ $t("page.book.reader.prevChapter") }}
            </NButton>
            <NButton text size="large" :style="{ color: textColor }" @click="handleOpenCatalog">
              <template #icon>
                <SvgIcon icon="solar:list-linear" size="18" />
              </template>
              {{ $t("page.book.reader.catalog") }}
            </NButton>
            <NButton :disabled="!hasNext" text size="large" :style="{ color: textColor }" @click="nextChapter">
              {{ $t("page.book.reader.nextChapter") }}
              <template #icon>
                <SvgIcon icon="solar:arrow-right-linear" size="18" />
              </template>
            </NButton>
          </div>
        </div>

        <!-- 桌面端右侧操作栏 -->
        <aside class="sticky top-72px w-56px shrink-0 flex flex-col items-center gap-5 pt-6 max-lg:hidden">
          <div
            class="w-40px h-40px rd-1 flex flex-col items-center justify-center gap-0.5 cursor-pointer transition-all duration-200"
            :class="darkMode ? 'text-gray-400 hover:bg-gray-700' : 'text-gray-500 hover:bg-gray-100'"
            @click="handleOpenCatalog"
            :title="$t('page.book.reader.catalog')">
            <SvgIcon icon="solar:list-linear" size="20" />
            <span class="text-[10px] leading-none">{{ $t("page.book.reader.catalog") }}</span>
          </div>
          <div
            class="w-40px h-40px rd-1 flex flex-col items-center justify-center gap-0.5 cursor-pointer transition-all duration-200"
            :class="darkMode ? 'text-gray-400 hover:bg-gray-700' : 'text-gray-500 hover:bg-gray-100'"
            @click="router.push({ name: 'book-detail', params: { id: bookId } })"
            :title="$t('page.book.reader.detail')">
            <SvgIcon icon="solar:info-circle-linear" size="20" />
            <span class="text-[10px] leading-none">{{ $t("page.book.reader.detail") }}</span>
          </div>
          <div
            class="w-40px h-40px rd-1 flex flex-col items-center justify-center gap-0.5 cursor-pointer transition-all duration-200"
            :class="darkMode ? 'text-gray-400 hover:bg-gray-700' : 'text-gray-500 hover:bg-gray-100'"
            :title="$t('page.book.reader.addShelf')">
            <SvgIcon icon="solar:bookmark-linear" size="20" />
            <span class="text-[10px] leading-none">{{ $t("page.book.reader.addShelf") }}</span>
          </div>
          <div
            class="w-40px h-40px rd-1 flex flex-col items-center justify-center gap-0.5 cursor-pointer transition-all duration-200"
            :class="darkMode ? 'text-gray-400 hover:bg-gray-700' : 'text-gray-500 hover:bg-gray-100'"
            :title="$t(darkMode ? 'page.book.reader.dayMode' : 'page.book.reader.nightMode')"
            @click="darkMode = !darkMode">
            <SvgIcon :icon="darkMode ? 'solar:sun-bold' : 'solar:moon-linear'" size="20" />
            <span class="text-[10px] leading-none">{{ darkMode ? $t("page.book.reader.dayMode") : $t("page.book.reader.nightMode") }}</span>
          </div>
          <div
            class="w-40px h-40px rd-1 flex flex-col items-center justify-center gap-0.5 cursor-pointer transition-all duration-200"
            :class="darkMode ? 'text-gray-400 hover:bg-gray-700' : 'text-gray-500 hover:bg-gray-100'"
            :title="$t('page.book.reader.phone')">
            <SvgIcon icon="solar:smartphone-linear" size="20" />
            <span class="text-[10px] leading-none">{{ $t("page.book.reader.phone") }}</span>
          </div>
        </aside>
      </div>

      <!-- 手机端浮动底部导航栏 -->
      <div
        class="fixed bottom-0 left-0 right-0 z-50 flex items-center justify-around px-4 py-3 lg:hidden"
        :style="{ backgroundColor: readerBg, borderTop: `1px solid ${readerBorder}` }">
        <NButton :disabled="!hasPrev" text size="small" :style="{ color: textColor }" @click="prevChapter">
          <template #icon>
            <SvgIcon icon="solar:arrow-left-linear" size="18" />
          </template>
          {{ $t("page.book.reader.prevChapter") }}
        </NButton>
        <NButton text size="small" :style="{ color: textColor }" @click="handleOpenCatalog">
          <template #icon>
            <SvgIcon icon="solar:list-linear" size="18" />
          </template>
          {{ $t("page.book.reader.catalog") }}
        </NButton>
        <NButton :disabled="!hasNext" text size="small" :style="{ color: textColor }" @click="nextChapter">
          {{ $t("page.book.reader.nextChapter") }}
          <template #icon>
            <SvgIcon icon="solar:arrow-right-linear" size="18" />
          </template>
        </NButton>
      </div>

      <!-- 手机端顶部操作栏 -->
      <div class="fixed top-0 left-0 right-0 z-50 flex items-center justify-between px-4 py-2 lg:hidden"
        :style="{ backgroundColor: readerBg, borderBottom: `1px solid ${readerBorder}` }">
        <div
          class="w-36px h-36px rd-1 flex items-center justify-center cursor-pointer"
          :class="darkMode ? 'text-gray-400 hover:bg-gray-700' : 'text-gray-500 hover:bg-gray-100'"
          @click="darkMode = !darkMode"
          :title="$t(darkMode ? 'page.book.reader.dayMode' : 'page.book.reader.nightMode')">
          <SvgIcon :icon="darkMode ? 'solar:sun-bold' : 'solar:moon-linear'" size="18" />
        </div>
        <div class="text-sm font-medium truncate max-w-[60%] text-center text-inherit">
          {{ sortedChapters.length > 0 ? sortedChapters[1]?.title || sortedChapters[0].title : '' }}
        </div>
        <div
          class="w-36px h-36px rd-1 flex items-center justify-center cursor-pointer"
          :class="darkMode ? 'text-gray-400 hover:bg-gray-700' : 'text-gray-500 hover:bg-gray-100'"
          @click="router.push({ name: 'book-detail', params: { id: bookId } })"
          :title="$t('page.book.reader.detail')">
          <SvgIcon icon="solar:info-circle-linear" size="18" />
        </div>
      </div>

      <CatalogModal
        v-model:visible="catalogVisible"
        :book-id="bookId"
        :chapter-no="sortedChapters.length > 0 ? sortedChapters[1]?.chapterNo || sortedChapters[0].chapterNo : 1"
        :dark-mode="darkMode"
        @select="goToChapter" />
    </main>
  </div>
</template>

<style scoped>
.text-indent-2em {
  text-indent: 2em;
}

.noise-bg {
  background-image: url(data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACQAAAAkCAQAAABLCVATAAAAAXNSR0IArs4c6QAAAWlJREFUeNqlloFtwzAMBFU0PxF34k5eoUtktBZVnBysixEEUWBD+idpUnrTGXNkhPvt18awXZjKyE5BE5ix8sw6IEQHfKA1kZoMF5ZNnndqy1k2vae+wTAjMBPIp+sY3QJP1JADaXtvFjv4LR1TFKA5GD4suFSQcGEhjPWRn2+zKpRLT0hBwSo3lRerdpScpbMQCgZS2cH4tHQwerJVPIQjUVBH9wFTPOMgxnRwObhWLLkKlpaJA8TnpDxBwEv1r8Uo+ImegDVX4DBXKKWt3mQnZRRMlxZ7vfxDra6j0vD8vKUtKvJ79Pt1X9W6XxZNTvphhYxcGEjneWncGVH3pM2kAs6Qlq4XDIus4x2qDKieYEsz0nTAYd96MelYZEEgElZxnJtEa4mefZpr7hHGsLLmS2uDVgPGEUadgBxwrn3zwRwGhkU2NVqy6fUEbRs1CruoCM5zlPaIIL6/biLs0edft/d7IfjhT9gfL6wnSxDYPyIAAAAASUVORK5CYII=);
  background-attachment: scroll;
}

.noise-bg > * {
  position: relative;
  z-index: 1;
}
</style>
