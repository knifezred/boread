<script setup lang="ts">
import { ref, computed, onMounted } from "vue"
import { useRoute, useRouter } from "vue-router"
import { NButton, NCard, NTag, NSpace, NSpin, NPagination } from "naive-ui"
import BookCard from "@/components/book-card.vue"
import {
  fetchGetBook,
  fetchGetChapterList,
  fetchReParseChapters,
  fetchAddToBookshelf,
  fetchGetReadProgress,
} from "@/service/api"
import { $t } from "@/locales"
import { formatWordCount, formatTime } from "@/utils/book"

defineOptions({ name: "BookDetail" });

const route = useRoute();
const router = useRouter();
const bookId = route.params.id as string;

const bookInfo = ref<Api.BookManage.Book>({
  id: 0,
  title: "",
  author: "",
  cover: "",
  createBy: "",
  createTime: "",
  updateBy: "",
  updateTime: "",
  status: "1",
  intro: null,
  categoryId: null,
  language: "",
  serialStatus: "1",
  visibility: "1",
  primaryFileId: null,
  totalChapters: 0,
  totalWords: 0,
  aggregateStatus: "1",
  avgRating: 0,
  ratingCount: 0,
  deptId: null,
});

const loading = ref(false);

/** 是否有阅读进度（"继续阅读"vs"立即阅读"） */
const hasReadProgress = ref(false);
const lastReadChapterNo = ref(1);

const relatedBooks = ref<Api.BookManage.Book[]>([]);
const chapters = ref<Api.BookManage.BookChapter[]>([]);
const chapterTotal = ref(0);
const chapterPage = ref(1);
const chapterSize = ref(42);
const chapterLoading = ref(false);
const sortAsc = ref(true);

const sortedChapters = computed(() => {
  const list = [...chapters.value];
  list.sort((a, b) =>
    sortAsc.value ? a.chapterNo - b.chapterNo : b.chapterNo - a.chapterNo,
  );
  return list;
});

const displayLimit = ref(10);
const showAll = ref(true);

const displayChapters = computed(() => {
  if (showAll.value) return sortedChapters.value;
  return sortedChapters.value.slice(0, displayLimit.value);
});

function toggleShowAll() {
  showAll.value = !showAll.value;
}

function toggleSort() {
  sortAsc.value = !sortAsc.value;
}

function goToReader(chapterNo?: number) {
  const target = chapterNo ?? lastReadChapterNo.value;
  router.push({
    name: "book-reader",
    query: { id: bookInfo.value.id, chapterNo: String(target) },
  });
}

async function loadChapters(page = 1) {
  chapterLoading.value = true;
  chapterPage.value = page;
  try {
    const { data } = await fetchGetChapterList({
      bookId: Number(bookId),
      current: page,
      size: chapterSize.value,
    });
    if (data) {
      chapters.value = data.records || [];
      chapterTotal.value = data.total || 0;
    }
  } catch (err) {
    console.error("加载章节列表失败:", err);
  } finally {
    chapterLoading.value = false;
  }
}

async function handleReParse() {
  const dialogRef = window.$dialog?.warning({
    title: $t("page.book.detail.reParseTitle"),
    content: $t("page.book.detail.reParseConfirm"),
    positiveText: $t("common.confirm"),
    negativeText: $t("common.cancel"),
    onPositiveClick: async () => {
      try {
        const { data } = await fetchReParseChapters(Number(bookId));
        if (data) {
          window.$message?.success(
            $t("page.book.detail.reParseSuccess", {
              old: data.oldCount,
              new: data.newCount,
            }),
          );
          bookInfo.value.totalChapters = data.newCount;
          bookInfo.value.totalWords = data.totalWords;
          await loadChapters(1);
        }
        dialogRef?.destroy();
      } catch (err: any) {
        window.$message?.error(
          err.message || $t("page.book.detail.reParseFailed"),
        );
        dialogRef?.destroy();
      }
    },
  });
}

/** 添加到书架 */
async function handleAddToShelf() {
  try {
    const { data } = await fetchAddToBookshelf({ bookId: bookInfo.value.id });
    if (data) {
      window.$message?.success($t("page.book.detail.addShelfSuccess"));
    }
  } catch (err: any) {
    window.$message?.error(
      err.message || $t("page.book.detail.addShelfFailed"),
    );
  }
}

onMounted(async () => {
  loading.value = true;
  try {
    const { data: bookData } = await fetchGetBook(bookId);
    if (bookData) {
      bookInfo.value = bookData;
      await loadChapters(1);
    }
    // 查询阅读进度
    const { data: progressData } = await fetchGetReadProgress(bookId);
    if (progressData?.chapterNo) {
      hasReadProgress.value = true;
      lastReadChapterNo.value = progressData.chapterNo;
    }
  } catch (err) {
    console.error("加载书籍详情失败:", err);
  } finally {
    loading.value = false;
  }
});
</script>

<template>
  <div class="min-h-screen bg-gray-100 lg:px-10 px-4 lg:py-5 py-4 font-sans">
    <!-- 面包屑 -->
    <div class="flex items-center gap-2 text-sm mb-5">
      <span>{{ $t("page.book.detail.breadcrumbHome") }}</span>
      <SvgIcon icon="solar:alt-arrow-right-linear" size="14" />
      <span>{{
        bookInfo.categoryName || $t("page.book.home.uncategorized")
      }}</span>
      <SvgIcon icon="solar:alt-arrow-right-linear" size="14" />
      <span>{{ bookInfo.title }}</span>
    </div>

    <div class="flex gap-6 max-w-1600px mx-auto items-start">
      <main class="flex-1 min-w-0 flex flex-col gap-4">
        <section id="section-info">
          <NCard class="rd-12px shadow-sm" :bordered="false" size="huge">
            <div class="flex lg:flex-row flex-col gap-6 lg:gap-8">
              <div class="shrink-0 lg:w-160px w-120px mx-auto lg:mx-0">
                <NSpin v-if="loading" show />
                <BookCard
                  v-else-if="bookInfo"
                  :book="bookInfo"
                  :show-status-tag="true"
                  class="only-cover"
                />
              </div>

              <div v-if="bookInfo" class="flex-1 text-center lg:text-left">
                <h1 class="lg:text-4xl text-2xl font-bold mb-4">
                  {{ bookInfo.title }}
                </h1>
                <div
                  class="flex lg:flex-row flex-col lg:gap-6 gap-1 text-sm mb-3"
                >
                  <span>
                    {{ $t("page.book.detail.author") }}:
                    {{ bookInfo.author }}
                  </span>
                  <span>
                    {{ $t("page.book.detail.updateTime") }}:
                    {{ formatTime(bookInfo.updateTime) }}
                  </span>
                </div>
                <div class="text-sm mb-3">
                  <span>{{ $t("page.book.detail.latestChapter") }}: </span>
                  <span class="text-primary no-underline">{{
                    bookInfo.latestChapterTitle || $t("page.book.reader.loading")
                  }}</span>
                </div>
                <div
                  v-if="bookInfo.tags?.length"
                  class="flex gap-2 mb-4 justify-center lg:justify-start"
                >
                  <NTag
                    v-for="tag in bookInfo.tags"
                    :key="tag.id"
                    size="small"
                    round
                    type="info"
                    ghost
                  >
                    {{ tag.tagName }}
                  </NTag>
                </div>
                <div
                  class="flex lg:gap-8 gap-6 justify-center lg:justify-start my-6"
                >
                  <div class="inline-flex items-end gap-1.5">
                    <span class="text-xl">{{
                      formatWordCount(bookInfo.totalWords)
                    }}</span>
                    <span class="text-xs">{{
                      $t("page.book.detail.words")
                    }}</span>
                  </div>
                  <div class="inline-flex items-end gap-1.5">
                    <span class="text-xl">{{ bookInfo.totalChapters }}</span>
                    <span class="text-xs">{{
                      $t("page.book.detail.chapters")
                    }}</span>
                  </div>
                  <!--
 <div class="inline-flex items-end gap-1.5">
                    <span class="text-xl">{{ bookInfo.avgRating }}</span>
                    <span class="text-xs">{{ $t("page.book.detail.rating") }}</span>
                  </div>
-->
                </div>

                <div class="flex justify-center lg:justify-start">
                  <NSpace size="medium">
                    <NButton
                      size="large"
                      ghost
                      type="primary"
                      @click="goToReader()"
                    >
                      {{ hasReadProgress ? $t("page.book.detail.continueRead") : $t("page.book.detail.readNow") }}
                    </NButton>
                    <NButton size="large" ghost type="primary" @click="handleAddToShelf">
                      {{ $t("page.book.detail.addToShelf") }}
                    </NButton>
                  </NSpace>
                </div>
              </div>
            </div>
          </NCard>
        </section>

        <NCard
          v-if="bookInfo.intro"
          class="rd-12px shadow-sm"
          :bordered="false"
          size="huge"
        >
          <NH2>{{ $t("page.book.detail.introTitle") }}</NH2>
          <p class="leading-relaxed mb-5 whitespace-break-spaces">
            {{ bookInfo.intro }}
          </p>
        </NCard>

        <section id="section-catalog">
          <NCard class="rd-12px shadow-sm" :bordered="false" size="huge">
            <template #header>
              <div class="flex items-center gap-3 w-full">
                <span class="text-xl font-semibold">{{
                  $t("page.book.detail.catalog")
                }}</span>
                <span class="text-sm font-normal">{{
                  $t("page.book.detail.totalChapters", {
                    total: chapterTotal,
                  })
                }}</span>
                <div class="ml-auto flex items-center gap-2">
                  <NButton size="tiny" quaternary @click="handleReParse">
                    {{ $t("page.book.detail.reParse") }}
                  </NButton>
                  <div
                    class="flex items-center gap-1 text-xs cursor-pointer px-2 py-1 rd-1 transition-all duration-200 hover:bg-gray-100"
                    @click="toggleSort"
                  >
                    <SvgIcon
                      :icon="
                        sortAsc
                          ? 'solar:sort-from-top-linear'
                          : 'solar:sort-from-bottom-linear'
                      "
                      size="16"
                    />
                    <span>{{
                      sortAsc
                        ? $t("page.book.detail.ascSort")
                        : $t("page.book.detail.descSort")
                    }}</span>
                  </div>
                </div>
              </div>
            </template>

            <NSpin :show="chapterLoading">
              <div class="grid lg:grid-cols-3 grid-cols-1 gap-3 gap-x-6">
                <div
                  v-for="ch in displayChapters"
                  :key="ch.id"
                  class="flex items-center gap-1.5 px-3 py-2 rd-1.5 cursor-pointer transition-colors duration-200 text-sm hover:bg-gray-100"
                  @click="goToReader(ch.chapterNo)"
                >
                  <span
                    class="flex-1 overflow-hidden text-ellipsis whitespace-nowrap"
                  >
                    {{ ch.title }}
                  </span>
                </div>
              </div>
            </NSpin>

            <!-- 移动端展开/收起 -->
            <div
              v-if="sortedChapters.length > displayLimit"
              class="flex justify-center mt-4 lg:hidden"
            >
              <NButton text size="small" @click="toggleShowAll">
                {{
                  showAll
                    ? $t("page.book.detail.collapse")
                    : $t("page.book.detail.expand", {
                      count: sortedChapters.length - displayLimit,
                    })
                }}
                <template #icon>
                  <SvgIcon
                    :icon="
                      showAll
                        ? 'solar:alt-arrow-up-linear'
                        : 'solar:alt-arrow-down-linear'
                    "
                    size="16"
                  />
                </template>
              </NButton>
            </div>

            <div
              v-if="!chapters.length && !chapterLoading"
              class="py-10 text-center text-sm"
            >
              {{ $t("page.book.detail.noChapters") }}
            </div>

            <div
              v-if="chapterTotal > chapterSize"
              class="flex justify-center mt-8"
            >
              <NPagination
                :page="chapterPage"
                :page-size="chapterSize"
                :item-count="chapterTotal"
                @update:page="loadChapters"
              />
            </div>
          </NCard>
        </section>

        <!-- 推荐位（移动端：目录下方） -->
        <div class="flex flex-col gap-4 lg:hidden">
          <NCard class="rd-12px shadow-sm" :bordered="false" size="small">
            <template #header>
              <div class="flex justify-between items-center w-full">
                <span>{{ $t("page.book.detail.authorOtherWorks") }}</span>
                <a href="#" class="text-xs no-underline hover:text-primary">{{ relatedBooks.length }}
                  {{ $t("page.book.detail.books") }}</a>
              </div>
            </template>
            <div class="grid grid-cols-2 gap-4">
              <div
                v-for="book in relatedBooks"
                :key="book.id"
                class="flex gap-3"
              >
                <div class="shrink-0 w-50px">
                  <BookCard :book="book" class="small-book-card" />
                </div>
                <div class="flex-1 overflow-hidden">
                  <div class="text-sm font-medium mb-1 truncate">
                    {{ book.title }}
                  </div>
                  <div class="text-xs leading-relaxed line-clamp-2">
                    {{ book.intro?.slice(0, 40) || "" }}
                  </div>
                </div>
              </div>
            </div>
          </NCard>
          <NCard class="rd-12px shadow-sm" :bordered="false" size="small">
            <template #header>
              <div class="flex justify-between items-center w-full">
                <span>{{ $t("page.book.detail.similarRecommend") }}</span>
                <a href="#" class="text-xs no-underline hover:text-primary">{{
                  $t("page.book.detail.more")
                }}</a>
              </div>
            </template>
            <div class="grid grid-cols-2 gap-4">
              <div
                v-for="book in relatedBooks"
                :key="book.id"
                class="flex gap-3"
              >
                <div class="shrink-0 w-50px">
                  <BookCard :book="book" class="small-book-card" />
                </div>
                <div class="flex-1 overflow-hidden">
                  <div class="text-sm font-medium mb-1 truncate">
                    {{ book.title }}
                  </div>
                  <div class="text-xs leading-relaxed line-clamp-2">
                    {{ book.intro?.slice(0, 40) || "" }}
                  </div>
                </div>
              </div>
            </div>
          </NCard>
        </div>
      </main>

      <!-- 桌面端右侧推荐栏 -->
      <aside class="max-lg:hidden w-300px shrink-0 flex flex-col gap-4">
        <NCard class="rd-12px shadow-sm" :bordered="false" size="small">
          <template #header>
            <div class="flex justify-between items-center w-full">
              <span>{{ $t("page.book.detail.authorOtherWorks") }}</span>
              <a href="#" class="text-xs no-underline hover:text-primary">{{ relatedBooks.length }} {{ $t("page.book.detail.books") }}</a>
            </div>
          </template>
          <div class="flex flex-col gap-4">
            <div v-for="book in relatedBooks" :key="book.id" class="flex gap-3">
              <div class="shrink-0 w-50px">
                <BookCard :book="book" class="small-book-card" />
              </div>
              <div class="flex-1 overflow-hidden">
                <div class="text-sm font-medium mb-1 truncate">
                  {{ book.title }}
                </div>
                <div class="text-xs leading-relaxed line-clamp-2">
                  {{ book.intro?.slice(0, 40) || "" }}
                </div>
              </div>
            </div>
          </div>
        </NCard>
        <NCard class="rd-12px shadow-sm" :bordered="false" size="small">
          <template #header>
            <div class="flex justify-between items-center w-full">
              <span>{{ $t("page.book.detail.similarRecommend") }}</span>
              <a href="#" class="text-xs no-underline hover:text-primary">{{
                $t("page.book.detail.more")
              }}</a>
            </div>
          </template>
          <div class="flex flex-col gap-4">
            <div v-for="book in relatedBooks" :key="book.id" class="flex gap-3">
              <div class="shrink-0 w-50px">
                <BookCard :book="book" class="small-book-card" />
              </div>
              <div class="flex-1 overflow-hidden">
                <div class="text-sm font-medium mb-1 truncate">
                  {{ book.title }}
                </div>
                <div class="text-xs leading-relaxed line-clamp-2">
                  {{ book.intro?.slice(0, 40) || "" }}
                </div>
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
