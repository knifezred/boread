<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { NButton, NSpin } from 'naive-ui'
import { fetchGetBook, fetchGetChapterContent } from "@/service/api"
import { useBoolean } from '@sa/hooks'
import { $t } from "@/locales"
import CatalogModal from './modules/catalog-modal.vue'

defineOptions({ name: 'BookReader' })

const route = useRoute()
const router = useRouter()
const bookId = (route.params.id || route.query.id) as string
const chapterNo = ref(Number(route.query.chapterNo) || 1)

const darkMode = ref(false)
const fontSize = ref(18)
const lineHeight = ref(1.8)
const chapterLoading = ref(false)
const bookLoading = ref(true)

const bookInfo = ref<Api.BookManage.Book | null>(null)
const chapterInfo = ref({
  id: '',
  title: '',
  updateTime: '',
  wordCount: '',
  content: ''
})

const { bool: catalogVisible, setTrue: openModel } = useBoolean()

const pageBg = computed(() => darkMode.value ? '#222' : '#fff')
const readerBg = computed(() => darkMode.value ? '#1a1a1a' : '#f2eede')
const textColor = computed(() => darkMode.value ? '#e0e0e0' : '#333')
const readerBorder = computed(() => darkMode.value ? 'rgba(255,255,255,0.06)' : 'rgba(0,0,0,0.04)')

function handleOpenCatalog() {
  openModel()
}

function goToChapter(no: number) {
  chapterNo.value = no
  router.replace({ query: { ...route.query, chapterNo: String(no) } })
  loadChapter()
}

function prevChapter() {
  if (chapterNo.value > 1) {
    chapterNo.value--
    router.replace({ query: { ...route.query, chapterNo: String(chapterNo.value) } })
    loadChapter()
  }
}

function nextChapter() {
  chapterNo.value++
  router.replace({ query: { ...route.query, chapterNo: String(chapterNo.value) } })
  loadChapter()
}

async function loadChapter() {
  if (!bookId) return
  chapterLoading.value = true
  try {
    const { data } = await fetchGetChapterContent(bookId, String(chapterNo.value))
    if (data) {
      chapterInfo.value = {
        id: String(data.id),
        title: data.title,
        updateTime: data.updateTime || '',
        wordCount: `${data.wordCount}${$t("page.book.detail.words")}`,
        content: data.content || ''
      }
    }
  } catch (err) {
    console.error('加载章节内容失败:', err)
  } finally {
    chapterLoading.value = false
  }
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
  loadChapter()
})
</script>

<template>
  <div class="min-h-screen font-sans relative transition-colors duration-300"
    :style="{ backgroundColor: pageBg, color: textColor }">
    <!-- 主内容 -->
    <main class="flex justify-center">
      <div class="flex gap-4 items-start w-full max-w-800px">
        <div class="w-full px-8 py-6 rd-1 min-h-100vh flex flex-col transition-colors duration-300"
          :style="{ backgroundColor: readerBg }">
          <div class="flex-1">
            <!-- 书籍详情（仅第一章） -->
            <div v-if="chapterNo == 1 && bookInfo" class="mb-10 pb-10 border-b text-center"
              :style="{ borderColor: readerBorder }">
              <div
                v-if="bookInfo.cover"
                class="w-130px h-180px rd-1 shadow-lg overflow-hidden mx-auto mb-6 mt-16">
                <img :src="bookInfo.cover" :alt="bookInfo.title" class="w-full h-full object-cover" />
              </div>
              <h1 class="text-3xl font-bold mb-2 text-inherit">{{ bookInfo.title }}</h1>
              <p class="text-sm mb-8" :class="darkMode ? 'text-gray-400' : 'text-gray-500'">{{ bookInfo.author }} {{ $t("page.book.reader.writtenBy") }}</p>
              <div class="flex justify-center gap-16 mb-8">
                <div class="flex flex-col gap-1">
                  <span class="text-base font-medium">{{ bookInfo.categoryName || $t("page.book.reader.unknown") }}</span>
                  <span class="text-xs" :class="darkMode ? 'text-gray-400' : 'text-gray-500'">{{ $t("page.book.reader.type") }}</span>
                </div>
                <div class="flex flex-col gap-1">
                  <span class="text-base font-medium">{{ bookInfo.serialStatus === '1' ? $t("page.book.reader.ongoing") : $t("page.book.reader.finished") }}</span>
                  <span class="text-xs" :class="darkMode ? 'text-gray-400' : 'text-gray-500'">{{ $t("page.book.reader.status") }}</span>
                </div>
                <div class="flex flex-col gap-1">
                  <span class="text-base font-medium">{{ bookInfo.totalWords || 0 }}{{ $t("page.book.detail.words") }}</span>
                  <span class="text-xs" :class="darkMode ? 'text-gray-400' : 'text-gray-500'">{{ $t("page.book.reader.words") }}</span>
                </div>
              </div>
              <div v-if="bookInfo.intro" class="mb-6 max-w-2xl mx-18 text-align-left">
                <p v-for="(paragraph, index) in bookInfo.intro.split('\n').filter(p => p.trim())" :key="index"
                  class="text-sm leading-relaxed mb-2 last:mb-0" :class="darkMode ? 'text-gray-400' : 'text-gray-500'">
                  {{ paragraph.trim() }}
                </p>
              </div>
            </div>

            <!-- 章节标题 -->
            <div class="text-center mb-10 pt-2">
              <h1 class="text-2xl font-bold mb-3 text-inherit">{{ chapterInfo.title || $t("page.book.reader.loading") }}</h1>
              <div class="flex justify-center gap-4 text-xs" :class="darkMode ? 'text-gray-500' : 'text-gray-400'">
                <span>{{ $t("page.book.reader.words") }}: {{ chapterInfo.wordCount || '0' }}</span>
                <span>{{ $t("page.book.detail.updateTime") }}: {{ chapterInfo.updateTime || '-' }}</span>
              </div>
            </div>

            <!-- 正文 -->
            <div class="mb-12 text-inherit leading-relaxed"
              :style="{ fontSize: `${fontSize}px`, lineHeight: lineHeight }">
              <NSpin :show="chapterLoading || bookLoading">
                <template v-if="chapterInfo.content">
                  <template v-for="(paragraph, index) in chapterInfo.content.split('\n').filter(p => p.trim())"
                    :key="index">
                    <p class="mb-5 text-indent-2em tracking-normal">
                      {{ paragraph.trim() }}
                    </p>
                  </template>
                </template>
                <div v-else-if="!chapterLoading && !bookLoading" class="py-20 text-center text-sm"
                  :class="darkMode ? 'text-gray-500' : 'text-gray-400'">
                  {{ $t("page.book.reader.noContent") }}
                </div>
              </NSpin>
            </div>
          </div>

          <!-- 上下章 -->
          <div class="flex justify-center gap-6 pt-6 border-t shrink-0" :style="{ borderColor: readerBorder }">
            <NButton :disabled="chapterNo <= 1" text size="large" :style="{ color: textColor }" @click="prevChapter">
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
            <NButton text size="large" :style="{ color: textColor }" @click="nextChapter">
              {{ $t("page.book.reader.nextChapter") }}
              <template #icon>
                <SvgIcon icon="solar:arrow-right-linear" size="18" />
              </template>
            </NButton>
          </div>
        </div>

        <!-- 右侧操作栏 -->
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
      <!-- 目录弹窗 -->
      <CatalogModal
        v-model:visible="catalogVisible"
        :book-id="bookId"
        :chapter-no="chapterNo"
        :dark-mode="darkMode"
        @select="goToChapter" />
    </main>

  </div>
</template>

<style scoped>
.text-indent-2em {
  text-indent: 2em;
}
</style>
