<script setup lang="tsx">
import { ref, onMounted } from "vue"
import { useRouter } from "vue-router"
import {
    NButton,
    NEmpty,
    NPagination,
    NSpin,
    NInput,
    NModal,
    NForm,
    NFormItem,
    NDropdown,
    useMessage,
} from "naive-ui"
import { useBoolean } from "@sa/hooks"
import {
    fetchGetBookshelfPage,
    fetchRemoveFromBookshelf,
    fetchUpdateBookshelf,
    fetchListBookshelfGroups,
} from "@/service/api"
import { $t } from "@/locales"

const router = useRouter();
const message = useMessage();

/** 加载状态 */
const loading = ref(false);
/** 书架数据 */
const bookshelfData = ref<Api.BookManage.BookshelfItem[]>([]);
/** 总数 */
const total = ref(0);
/** 当前页码 */
const current = ref(1);
/** 每页数量 */
const pageSize = ref(24);

/** 分组列表 */
const groups = ref<Api.BookManage.BookshelfGroupItem[]>([]);
/** 当前选中的分组 */
const selectedGroup = ref<string>("");
/** 搜索关键词 */
const searchKeyword = ref("");

/** 操作弹窗相关 */
const { bool: showModal, setTrue: openModal, setFalse: closeModal } = useBoolean();
const editingBookId = ref<number>(0);
const editingGroupName = ref("");
const editingIsTop = ref(false);

/** 下拉菜单选中项 */
function handleDropdownSelect(key: string, item: Api.BookManage.BookshelfItem) {
  switch (key) {
    case "toggleTop":
      handleToggleTop(item);
      break;
    case "changeGroup":
      editingBookId.value = item.bookId;
      editingGroupName.value = item.groupName;
      editingIsTop.value = item.isTop;
      openModal();
      break;
    case "remove":
      handleRemove(item);
      break;
  }
}

/** 加载分组 */
async function loadGroups() {
  const { data, error } = await fetchListBookshelfGroups();
  if (!error && data) {
    groups.value = data;
  }
}

/** 加载书架数据 */
async function loadData() {
  loading.value = true;
  try {
    const params: Api.BookManage.BookshelfSearchParams = {
      current: current.value,
      size: pageSize.value,
      groupName: selectedGroup.value || null,
      keyword: searchKeyword.value || null,
    };
    const { data, error } = await fetchGetBookshelfPage(params);
    if (!error && data) {
      bookshelfData.value = data.records;
      total.value = data.total;
    }
  } catch (err) {
    console.error("加载书架数据失败:", err);
  } finally {
    loading.value = false;
  }
}

/** 切换分组 */
function handleGroupChange(groupName: string) {
  selectedGroup.value = groupName;
  current.value = 1;
  loadData();
}

/** 搜索 */
function handleSearch() {
  current.value = 1;
  loadData();
}

/** 重置搜索 */
function handleReset() {
  selectedGroup.value = "";
  searchKeyword.value = "";
  current.value = 1;
  loadData();
}

/** 切换置顶 */
async function handleToggleTop(item: Api.BookManage.BookshelfItem) {
  const { error } = await fetchUpdateBookshelf(item.bookId, { isTop: !item.isTop });
  if (!error) {
    message.success($t("common.updateSuccess"));
    loadData();
  }
}

/** 从书架移除 */
async function handleRemove(item: Api.BookManage.BookshelfItem) {
  const { error } = await fetchRemoveFromBookshelf(item.bookId);
  if (!error) {
    message.success($t("common.deleteSuccess"));
    loadData();
  }
}

/** 更新分组 */
async function handleUpdateGroup() {
  const { error } = await fetchUpdateBookshelf(editingBookId.value, {
    groupName: editingGroupName.value,
  });
  if (!error) {
    message.success($t("common.updateSuccess"));
    closeModal();
    loadData();
    loadGroups();
  }
}

/** 点击进入书籍详情或继续阅读 */
function handleBookClick(item: Api.BookManage.BookshelfItem) {
  if (item.readPercent > 0 && item.chapterNo) {
    router.push({ name: "book-reader", query: { id: item.bookId, chapterNo: String(item.chapterNo) } });
  } else {
    router.push({ name: "book-detail", params: { id: item.bookId } });
  }
}

/** 分页变化 */
function handlePageChange(page: number) {
  current.value = page;
  loadData();
}

/** 下拉菜单选项 */
const dropdownOptions = (item: Api.BookManage.BookshelfItem) => [
  {
    label: item.isTop ? $t("page.bookshelf.cancelTop") : $t("page.bookshelf.setTop"),
    key: "toggleTop",
  },
  {
    label: $t("page.bookshelf.changeGroup"),
    key: "changeGroup",
  },
  {
    label: $t("page.bookshelf.remove"),
    key: "remove",
  },
];

/** 格式化时间 */
function formatTime(time: string | null): string {
  if (!time) return "";
  return time.slice(0, 10);
}

/** 格式化进度百分比 */
function formatPercent(percent: number | undefined): string {
  if (percent == null || Number.isNaN(percent) || percent <= 0) return $t("page.bookshelf.notRead");
  if (percent >= 100) return $t("page.bookshelf.finished");
  return `${percent.toFixed(1)}%`;
}

onMounted(() => {
  loadGroups();
  loadData();
});
</script>

<template>
  <div class="flex-col-stretch gap-16px h-full">
    <!-- 搜索栏 -->
    <div class="flex items-center gap-12px flex-wrap">
      <NInput
        v-model:value="searchKeyword"
        :placeholder="$t('page.bookshelf.searchPlaceholder')"
        class="w-240px"
        clearable
        @keyup.enter="handleSearch"
      />
      <NButton type="primary" ghost @click="handleSearch">
        {{ $t("common.search") }}
      </NButton>
      <NButton @click="handleReset">
        {{ $t("common.reset") }}
      </NButton>
    </div>

    <!-- 分组标签 -->
    <div class="flex items-center gap-8px flex-wrap">
      <NButton
        :type="selectedGroup === '' ? 'primary' : 'default'"
        size="small"
        @click="handleGroupChange('')"
      >
        {{ $t("page.bookshelf.all") }}
      </NButton>
      <NButton
        v-for="g in groups"
        :key="g.groupName"
        :type="selectedGroup === g.groupName ? 'primary' : 'default'"
        size="small"
        @click="handleGroupChange(g.groupName)"
      >
        {{ g.groupName }} ({{ g.bookCount }})
      </NButton>
    </div>

    <!-- 书籍网格 -->
    <div class="flex-1 overflow-y-auto">
      <NSpin :show="loading">
        <template v-if="bookshelfData.length === 0 && !loading">
          <NEmpty :description="$t('page.bookshelf.empty')" class="mt-80px" />
        </template>
        <div v-else class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 gap-16px">
          <div
            v-for="item in bookshelfData"
            :key="item.id"
            class="group relative cursor-pointer rd-2 overflow-hidden bg-white dark:bg-[#1e1e1e] shadow-sm hover:shadow-md transition-shadow duration-200"
          >
            <!-- 更多操作按钮 -->
            <NDropdown
              :options="dropdownOptions(item)"
              trigger="click"
              placement="bottom-end"
              @select="(key: string) => handleDropdownSelect(key, item)"
            >
              <div
                class="absolute top-2 right-2 z-10 w-24px h-24px flex-center rd-1 bg-[rgba(0,0,0,0.4)] text-white opacity-0 group-hover:opacity-100 transition-opacity cursor-pointer"
              >
                <span class="i-carbon:overflow-menu-vertical text-14px" />
              </div>
            </NDropdown>

            <!-- 置顶标记 -->
            <div
              v-if="item.isTop"
              class="absolute top-2 left-2 z-10 px-1.5 py-0.5 rd-1 text-11px font-500 text-white bg-[rgba(245,158,11,0.9)]"
            >
              {{ $t("page.bookshelf.top") }}
            </div>

            <!-- 封面 -->
            <div class="w-full aspect-2/3 overflow-hidden" @click="handleBookClick(item)">
              <img
                v-if="item.bookCover"
                :src="item.bookCover"
                :alt="item.bookTitle"
                class="w-full h-full object-cover transition-transform duration-300 group-hover:scale-105"
                loading="lazy"
              />
              <div
                v-else
                class="w-full h-full flex-center p-4 box-border"
                :style="{
                  background: `linear-gradient(135deg, hsl(${(item.bookId * 137.508) % 360}, 70%, 60%), hsl(${(item.bookId * 137.508 + 60) % 360}, 70%, 45%))`,
                }"
              >
                <div class="text-center text-white text-shadow-sm w-full">
                  <div class="font-700 leading-normal mb-2 line-clamp-3 break-words text-14px">
                    {{ item.bookTitle }}
                  </div>
                  <div class="font-500 opacity-90 truncate text-12px">
                    {{ item.bookAuthor }}
                  </div>
                </div>
              </div>
            </div>

            <!-- 书籍信息 -->
            <div class="p-8px" @click="handleBookClick(item)">
              <h3 class="text-13px font-500 m-0 truncate">
                {{ item.bookTitle }}
              </h3>
              <p class="text-12px text-gray-400 m-0 truncate dark:text-gray-500">
                {{ item.bookAuthor }}
              </p>
              <!-- 进度条 -->
              <div class="mt-6px">
                <div class="flex items-center justify-between text-11px text-gray-400 mb-2px">
                  <span>{{ formatPercent(item.readPercent) }}</span>
                  <span v-if="item.lastReadTime">{{ formatTime(item.lastReadTime) }}</span>
                </div>
                <div class="w-full h-3px bg-gray-200 rd-1 dark:bg-gray-700 overflow-hidden">
                  <div
                    class="h-full rd-1 transition-all duration-300"
                    :class="item.readPercent >= 100 ? 'bg-green-500' : 'bg-blue-500'"
                    :style="{ width: Math.min(item.readPercent ?? 0, 100) + '%' }"
                  />
                </div>
              </div>
            </div>
          </div>
        </div>
      </NSpin>
    </div>

    <!-- 分页 -->
    <div class="flex justify-end pt-8px">
      <NPagination
        :page="current"
        :page-size="pageSize"
        :page-sizes="[12, 24, 48, 96]"
        :item-count="total"
        show-size-picker
        @update:page="handlePageChange"
        @update:page-size="(size: number) => { pageSize = size; current = 1; loadData(); }"
      />
    </div>

    <!-- 修改分组弹窗 -->
    <NModal
      :show="showModal"
      :title="$t('page.bookshelf.changeGroup')"
      preset="card"
      class="w-360px"
      @close="closeModal"
      @update:show="(val: boolean) => !val && closeModal()"
    >
      <NForm>
        <NFormItem :label="$t('page.bookshelf.groupName')">
          <NInput v-model:value="editingGroupName" :placeholder="$t('page.bookshelf.groupPlaceholder')" />
        </NFormItem>
        <div class="flex justify-end gap-8px mt-16px">
          <NButton @click="closeModal">{{ $t("common.cancel") }}</NButton>
          <NButton type="primary" @click="handleUpdateGroup">{{ $t("common.confirm") }}</NButton>
        </div>
      </NForm>
    </NModal>
  </div>
</template>

<style scoped></style>
