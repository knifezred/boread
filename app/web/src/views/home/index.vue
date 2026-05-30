<script setup lang="ts">
import { ref, computed, onMounted } from "vue"
import { NButton, NGrid, NGi, NModal, NSpace, NInput, NScrollbar, NEmpty, NSpin, NPagination } from "naive-ui"
import { fetchGetBookList, fetchGetChapterList, fetchGetChapterContent } from "@/service/api"
import BookCard from "@/components/book-card.vue"
import { useRouter } from "vue-router"
defineOptions({ name: "HomePage" });

const router = useRouter()
const searchKeyword = ref("")
const activeMenu = ref("library") // library-书库 folder-文件夹 shelf-我的书架 share-共享书架 random-随机推荐
const books = ref<Api.SystemManage.Book[]>([])
const loading = ref(false)
const total = ref(0)
const current = ref(1)
const size = ref(20)

// 导航菜单配置
const menuItems = [
  { key: 'library', label: '书库', icon: 'solar:notebook-square-linear' },
  { key: 'folder', label: '文件夹', icon: 'solar:folder-linear' },
  { key: 'shelf', label: '我的书架', icon: 'solar:book-linear' },
  { key: 'share', label: '共享书架', icon: 'solar:cash-out-linear' },
  { key: 'random', label: '随机推荐', icon: 'solar:atom-linear' },
]

// 响应式列数
const responsiveCols = computed(() => {
  const w = window.innerWidth
  if (w < 640) return 2
  if (w < 960) return 3
  if (w < 1200) return 4
  if (w < 1600) return 5
  return 6
})

// 章节弹窗
const chapterModalVisible = ref(false)
const chapterBook = ref<Api.SystemManage.Book | null>(null)
const chapters = ref<any[]>([])
const chapterLoading = ref(false)
const chapterTotal = ref(0)
const chapterPage = ref(1)

// 阅读弹窗
const readerVisible = ref(false)
const readerTitle = ref("")
const readerContent = ref("")
const readerLoading = ref(false)

onMounted(() => {
  loadBooks()
})

async function loadBooks(page = 1) {
  loading.value = true
  current.value = page
  const params: any = { current: page, size: size.value }
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

function handleSearch() {
  loadBooks(1)
}

function switchMenu(key: string) {
  activeMenu.value = key
  // 不同菜单可以加载不同数据，暂时统一加载书库
  loadBooks(1)
}

function showChapters(book: Api.SystemManage.Book) {
  router.push(`/book-detail/${book.id}`)
}

async function loadChapters() {
  if (!chapterBook.value) return
  chapterLoading.value = true
  const { data } = await fetchGetChapterList({ bookId: chapterBook.value.id, current: chapterPage.value, size: 50 })
  if (data) {
    chapters.value = data.records || []
    chapterTotal.value = data.total || 0
  }
  chapterLoading.value = false
}

async function readChapter(chapter: any) {
  if (!chapterBook.value) return
  readerLoading.value = true
  readerVisible.value = true
  readerTitle.value = `${chapterBook.value.title} - ${chapter.title}`
  const { data } = await fetchGetChapterContent(chapterBook.value.id, chapter.chapterNo)
  if (data) {
    readerContent.value = data.content || "（内容为空）"
  } else {
    readerContent.value = "（加载失败）"
  }
  readerLoading.value = false
}


</script>

<template>
  <div class="gallery-layout">
    <!-- 左侧边栏导航 -->
    <aside class="left-sidebar">
      <div class="sidebar-menu">
        <div v-for="item in menuItems" :key="item.key" class="menu-item"
          :class="{ active: activeMenu === item.key }" @click="switchMenu(item.key)">
          <SvgIcon class="menu-icon" :icon="item.icon" />
          <span class="menu-label">{{ item.label }}</span>
        </div>
      </div>
    </aside>

    <!-- 右侧主内容 -->
    <div class="main-container">
      <!-- 顶部导航栏 -->
      <header class="top-header">
        <div class="header-content">
          <div class="search-area">
            <NInput v-model:value="searchKeyword" placeholder="搜索书名、作者" clearable round
              @keyup.enter="handleSearch" class="search-input" />
          </div>
          <div class="header-actions">
            <NButton type="primary" ghost size="medium">
              <!-- <template #icon><NIcon :component="BookOutline" /></template> -->
              导入书籍
            </NButton>
          </div>
        </div>
      </header>

      <!-- 内容区域 -->
      <main class="content-area">
        <NSpin :show="loading">
          <div v-if="books.length === 0 && !loading" class="empty-container">
            <div class="empty-icon">
              <svg width="120" height="100" viewBox="0 0 120 100" fill="none" xmlns="http://www.w3.org/2000/svg">
                <path d="M60 10C45 10 35 20 35 35C35 50 45 60 60 60C75 60 85 50 85 35C85 20 75 10 60 10Z" fill="#E6F4FF"/>
                <path d="M55 30C55 27.2386 57.2386 25 60 25C62.7614 25 65 27.2386 65 30C65 32.7614 62.7614 35 60 35C57.2386 35 55 32.7614 55 30Z" fill="#B3D8FF"/>
                <path d="M40 65C40 59.4772 44.4772 55 50 55H70C75.5228 55 80 59.4772 80 65V85C80 87.7614 77.7614 90 75 90H45C42.2386 90 40 87.7614 40 85V65Z" fill="#8EC2FF"/>
                <path d="M45 75H75C76.6569 75 78 76.3431 78 78C78 79.6569 76.6569 81 75 81H45C43.3431 81 42 79.6569 42 78C42 76.3431 43.3431 75 45 75Z" fill="#E6F4FF"/>
              </svg>
            </div>
            <p class="empty-text">暂无内容</p>
          </div>
          <div v-else class="book-grid-container">
            <NGrid :x-gap="24" :y-gap="32" :cols="responsiveCols" class="book-grid">
              <NGi v-for="book in books" :key="book.id">
                <BookCard
                  :book="{
                    ...book,
                    status: book.serialStatus === '2' ? 'finished' : 'serial'
                  }"
                  :show-status-tag="true"
                  @click="showChapters(book)"
                />
              </NGi>
            </NGrid>
          </div>
        </NSpin>

        <!-- 分页 -->
        <div v-if="books.length > 0" class="pagination-bar">
          <NPagination :page="current" :page-size="size" :item-count="total" :page-sizes="[12, 24, 48, 60]"
            show-size-picker @update:page="loadBooks" @update:page-size="(s) => { size = s; loadBooks(1) }" />
        </div>
      </main>
    </div>

    <!-- 章节列表弹窗 -->
    <NModal v-model:show="chapterModalVisible" :title="chapterBook?.title || ''" preset="card" class="w-600px"
      @update:show="(val) => { if (val && chapterBook) loadChapters() }">
      <NScrollbar class="max-h-400px">
        <div v-if="chapterLoading" class="flex-center py-24px">
          <NSpin />
        </div>
        <div v-else>
          <div v-for="(ch, idx) in chapters" :key="ch.id || idx" class="chapter-item" @click="readChapter(ch)">
            <span class="chapter-no">第 {{ ch.chapterNo }} 章</span>
            <span class="chapter-title">{{ ch.title }}</span>
            <span class="chapter-words">
              {{ ch.wordCount ? (ch.wordCount > 1000 ? (ch.wordCount / 1000).toFixed(1) + 'k' : ch.wordCount) + '字' : '' }}
            </span>
          </div>
          <div v-if="chapters.length === 0" class="flex-center py-24px">
            <NEmpty description="暂无章节" />
          </div>
        </div>
      </NScrollbar>
      <template #footer>
        <NSpace justify="end">
          <NButton @click="chapterModalVisible = false">关闭</NButton>
        </NSpace>
      </template>
    </NModal>

    <!-- 阅读器弹窗 -->
    <NModal v-model:show="readerVisible" :title="readerTitle" preset="card" class="w-800px" :style="{ maxHeight: '80vh' }"
      :loading="readerLoading">
      <NScrollbar class="reader-content">
        <div class="reader-text">{{ readerContent }}</div>
      </NScrollbar>
      <template #footer>
        <NSpace justify="end">
          <NButton @click="readerVisible = false">关闭</NButton>
        </NSpace>
      </template>
    </NModal>
  </div>
</template>

<style scoped>
.gallery-layout {
  display: flex;
  height: 100vh;
  background-color: #fff;
  color: #333;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
  overflow: hidden;
}

/* 左侧边栏 */
.left-sidebar {
  width: 200px;
  flex-shrink: 0;
  background-color: #f8f9fa;
  padding: 20px 0;
  border-right: 1px solid #f0f0f0;
}

.sidebar-menu {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.menu-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 20px;
  cursor: pointer;
  transition: all 0.2s;
  color: #666;
  font-size: 15px;
}

.menu-item:hover {
  background-color: #f0f0f0;
  color: #333;
}

.menu-item.active {
  background-color: #e6f4ff;
  color: #2f96f3;
  font-weight: 500;
}

.menu-icon {
  font-size: 20px;
}

/* 右侧主容器 */
.main-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

/* 顶部导航 */
.top-header {
  height: 64px;
  border-bottom: 1px solid #f0f0f0;
  padding: 0 24px;
  display: flex;
  align-items: center;
}

.header-content {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.search-area {
  width: 40%;
  max-width: 400px;
}

.search-input {
  width: 100%;
}

.header-actions {
  display: flex;
  gap: 12px;
}

/* 内容区域 */
.content-area {
  flex: 1;
  padding: 24px;
  overflow-y: auto;
}

/* 空状态 */
.empty-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 120px 0;
}

.empty-icon {
  margin-bottom: 20px;
  opacity: 0.6;
}

.empty-text {
  font-size: 16px;
  color: #999;
}

/* 书籍网格 */
.book-grid-container {
  width: 100%;
}

.book-grid {
  width: 100%;
}

/* 分页 */
.pagination-bar {
  display: flex;
  justify-content: center;
  margin-top: 48px;
}

/* 章节列表 */
.chapter-item {
  display: flex;
  align-items: center;
  padding: 10px 16px;
  cursor: pointer;
  border-radius: 6px;
  transition: background-color 0.15s;
  gap: 12px;
}

.chapter-item:hover {
  background-color: #f5f5f5;
}

.chapter-no {
  font-size: 13px;
  color: #999;
  min-width: 72px;
}

.chapter-title {
  flex: 1;
  font-size: 14px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.chapter-words {
  font-size: 12px;
  color: #bbb;
}

/* 阅读器 */
.reader-content {
  max-height: 60vh;
  padding: 0 16px;
}

.reader-text {
  font-size: 16px;
  line-height: 1.8;
  color: #333;
  white-space: pre-wrap;
  word-wrap: break-word;
}

/* 响应式适配 */
@media (max-width: 960px) {
  .left-sidebar {
    width: 60px;
  }
  .menu-label {
    display: none;
  }
  .header-content {
    flex-direction: column;
    gap: 12px;
    padding: 12px 0;
  }
  .search-area {
    width: 100%;
    max-width: 100%;
  }
  .top-header {
    height: auto;
  }
  .content-area {
    padding: 16px;
  }
}

@media (max-width: 640px) {
  .left-sidebar {
    display: none;
  }
  .reader-content {
    max-height: 50vh;
  }
}
</style>
