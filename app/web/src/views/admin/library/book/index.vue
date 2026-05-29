<script setup lang="tsx">
import { ref } from "vue"
import type { Ref } from "vue"
import { NButton, NPopconfirm, NTag, NSpace, NSelect, NInput } from "naive-ui"
import { useBoolean } from "@sa/hooks"
import {
  bookStatusRecord,
} from "@/constants/business"
import { useDictItems } from "@/hooks/business/dict"
import {
  fetchGetBookList,
  fetchDeleteBook,
  fetchGetCategoryTree,
  fetchUpdateBookStatus,
} from "@/service/api"
import { useAppStore } from "@/store/modules/app"
import {
  defaultTransform,
  useNaivePaginatedTable,
  useTableOperate,
} from "@/hooks/common/table"
import { $t } from "@/locales"
import BookOperateModal from "./modules/book-operate-modal.vue"

const appStore = useAppStore();

const { bool: visible, setTrue: openModal } = useBoolean();

const categoryOptions = ref<CommonType.Option<number>[]>([]);
loadCategoryOptions();

const { options: serialStatusOptions, labelMap: serialStatusLabelMap } = useDictItems("book_serial_status");
const { options: visibilityOptions, labelMap: visibilityLabelMap } = useDictItems("book_visibility");

async function loadCategoryOptions() {
  const { data, error } = await fetchGetCategoryTree();
  if (!error && data) {
    categoryOptions.value = flattenTree(data, 0);
  }
}

function flattenTree(
  nodes: Api.SystemManage.BookCategory[],
  depth: number,
): CommonType.Option<number>[] {
  let result: CommonType.Option<number>[] = [];
  for (const n of nodes) {
    const prefix = "\u3000".repeat(depth);
    result.push({ value: n.id, label: `${prefix}${n.categoryName}` });
    if (n.children?.length) {
      result = result.concat(flattenTree(n.children, depth + 1));
    }
  }
  return result;
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
});

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
    searchParams.value.current = params.page;
    searchParams.value.size = params.pageSize;
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
        const tagMap: Record<string, NaiveUI.ThemeColor> = {
          "1": "info",
          "2": "success",
          "3": "warning",
        };
        return (
          <NTag type={tagMap[row.serialStatus]}>
            {serialStatusLabelMap.value[row.serialStatus] ?? row.serialStatus}
          </NTag>
        );
      },
    },
    {
      key: "visibility",
      title: $t("page.admin.library.book.visibility"),
      align: "center",
      width: 80,
      render: (row: Api.SystemManage.Book) => {
        const tagMap: Record<string, NaiveUI.ThemeColor> = {
          "1": "success",
          "2": "warning",
          "3": "info",
        };
        return (
          <NTag type={tagMap[row.visibility]}>
            {visibilityLabelMap.value[row.visibility] ?? row.visibility}
          </NTag>
        );
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
        const tagMap: Record<string, NaiveUI.ThemeColor> = {
          "1": "success",
          "2": "warning",
          "3": "info",
          "4": "error",
        };
        return (
          <NTag type={tagMap[row.status]}>
            {$t(bookStatusRecord[row.status])}
          </NTag>
        );
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
      width: 280,
      render: (row: Api.SystemManage.Book) => (
        <div class="flex-center justify-end gap-8px">
          <NPopconfirm onPositiveClick={() => handleToggleListing(row)}>
            {{
              default: () =>
                row.status === "1"
                  ? $t("common.confirmDelete")
                  : $t("page.admin.library.book.statusListed"),
              trigger: () => (
                <NButton size="small" ghost>
                  {row.status === "1"
                    ? $t("page.admin.library.book.statusUnlisted")
                    : $t("page.admin.library.book.statusListed")}
                </NButton>
              ),
            }}
          </NPopconfirm>
          <NButton
            type="primary"
            ghost
            size="small"
            onClick={() => handleEdit(row)}
          >
            {$t("common.edit")}
          </NButton>
          <NPopconfirm onPositiveClick={() => handleDelete(row.id)}>
            {{
              default: () => $t("common.confirmDelete"),
              trigger: () => (
                <NButton type="error" ghost size="small">
                  {$t("common.delete")}
                </NButton>
              ),
            }}
          </NPopconfirm>
        </div>
      ),
    },
  ],
});

const { checkedRowKeys, onBatchDeleted, onDeleted } = useTableOperate(
  data,
  "id",
  getData,
);

const operateType = ref<NaiveUI.TableOperateType>("add");
const editingData: Ref<Api.SystemManage.Book | null> = ref(null);

function handleAdd() {
  operateType.value = "add";
  editingData.value = null;
  openModal();
}
function handleEdit(item: Api.SystemManage.Book) {
  operateType.value = "edit";
  editingData.value = { ...item };
  openModal();
}

async function handleDelete(id: number) {
  const { error } = await fetchDeleteBook(id);
  if (!error) onDeleted();
}

async function handleToggleListing(row: Api.SystemManage.Book) {
  const newStatus: Api.SystemManage.BookListingStatus =
    row.status === "1" ? "2" : "1";
  const { error } = await fetchUpdateBookStatus(row.id, { status: newStatus });
  if (!error) {
    window.$message?.success($t("common.updateSuccess"));
    getData();
  }
}

async function handleBatchDelete() {
  onBatchDeleted();
}

function handleSearch() {
  getDataByPage(1);
}
function handleReset() {
  searchParams.value = {
    current: 1,
    size: 10,
    title: null,
    author: null,
    categoryId: null,
    status: null,
    visibility: null,
    serialStatus: null,
    tagId: null,
  };
  getDataByPage(1);
}
</script>

<template>
  <div class="flex-col-stretch gap-16px overflow-hidden lt-sm:overflow-auto">
    <NCard :bordered="false" size="small">
      <NSpace wrap :size="[12, 12]">
        <NInput
          v-model:value="searchParams.title"
          :placeholder="$t('page.admin.library.book.bookName')"
          clearable
          style="width: 160px"
          @keyup.enter="handleSearch"
        />
        <NInput
          v-model:value="searchParams.author"
          :placeholder="$t('page.admin.library.book.author')"
          clearable
          style="width: 160px"
          @keyup.enter="handleSearch"
        />
        <NSelect
          v-model:value="searchParams.serialStatus"
          :placeholder="$t('page.admin.library.book.form.serialStatus')"
          :options="serialStatusOptions"
          clearable
          style="width: 130px"
        />
        <NSelect
          v-model:value="searchParams.visibility"
          :placeholder="$t('page.admin.library.book.form.visibility')"
          :options="visibilityOptions"
          clearable
          style="width: 130px"
        />
        <NSelect
          v-model:value="searchParams.categoryId"
          :placeholder="$t('page.admin.library.book.categoryId')"
          :options="categoryOptions"
          clearable
          filterable
          style="width: 160px"
        />
        <NButton type="primary" @click="handleSearch">
          {{ $t("common.search") }}
        </NButton>
        <NButton @click="handleReset">{{ $t("common.reset") }}</NButton>
      </NSpace>
    </NCard>
    <NCard
      :title="$t('page.admin.library.book.title')"
      :bordered="false"
      size="small"
      class="card-wrapper sm:flex-1-hidden"
    >
      <template #header-extra>
        <TableHeaderOperation
          v-model:columns="columnChecks"
          :disabled-delete="checkedRowKeys.length === 0"
          :loading="loading"
          @add="handleAdd"
          @delete="handleBatchDelete"
          @refresh="getData"
        />
      </template>
      <NDataTable
        v-model:checked-row-keys="checkedRowKeys"
        :columns="columns"
        :data="data"
        size="small"
        :flex-height="!appStore.isMobile"
        :scroll-x="1200"
        :loading="loading"
        :row-key="(row) => row.id"
        remote
        :pagination="pagination"
        class="sm:h-full"
      />
      <BookOperateModal
        v-model:visible="visible"
        :operate-type="operateType"
        :row-data="editingData"
        @submitted="getDataByPage"
      />
    </NCard>
  </div>
</template>
