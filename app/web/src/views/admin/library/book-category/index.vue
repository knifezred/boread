<script setup lang="tsx">
import { ref } from "vue"
import type { Ref } from "vue"
import { NButton, NPopconfirm, NTag, NSpace, NSelect, NInput } from "naive-ui"
import { useBoolean } from "@sa/hooks"
import { enableStatusOptions, enableStatusRecord } from "@/constants/business"
import { fetchGetCategoryList, fetchDeleteCategory } from "@/service/api"
import { useAppStore } from "@/store/modules/app"
import {
  defaultTransform,
  useNaivePaginatedTable,
  useTableOperate,
} from "@/hooks/common/table"
import { $t } from "@/locales"
import CategoryOperateModal, {
  type OperateType,
} from "./modules/category-operate-modal.vue"

const appStore = useAppStore();

const { bool: visible, setTrue: openModal } = useBoolean();

/** isHot 筛选值: "true"/"false"/null，传给后端时转为 boolean */
const isHotValue = ref<string | null>(null)

const searchParams = ref<Api.BookManage.CategorySearchParams>({
  current: 1,
  size: 10,
  categoryName: null,
  categoryCode: null,
  isHot: null,
  status: null,
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
  api: () => {
    const params = { ...searchParams.value };
    if (isHotValue.value !== null) {
      params.isHot = isHotValue.value === "true";
    } else {
      params.isHot = null;
    }
    return fetchGetCategoryList(params);
  },
  onPaginationParamsChange: (params) => {
    searchParams.value.current = params.page;
    searchParams.value.size = params.pageSize;
  },
  transform: (response) => defaultTransform(response),
  columns: () => [
    {
      type: "selection",
      align: "center",
      width: 48,
    },
    {
      key: "categoryName",
      title: $t("page.admin.library.bookCategory.categoryName"),
      align: "left",
      width: 240,
    },
    {
      key: "categoryCode",
      title: $t("page.admin.library.bookCategory.categoryCode"),
      align: "center",
      width: 140,
    },
    {
      key: "description",
      title: $t("page.admin.library.bookCategory.description"),
      align: "left",
      ellipsis: { tooltip: true },
      minWidth: 160,
    },
    {
      key: "sortOrder",
      title: $t("page.admin.library.bookCategory.sortOrder"),
      align: "center",
      width: 80,
    },
    {
      key: "isHot",
      title: $t("page.admin.library.bookCategory.isHot"),
      align: "center",
      width: 100,
      render: (row: Api.BookManage.BookCategory) => {
        if (row.isHot) {
          return <NTag type="error" bordered={false}>{$t('common.yesOrNo.yes')}</NTag>;
        }
        return <NTag type="default" bordered={false}>{$t('common.yesOrNo.no')}</NTag>;
      },
    },
    {
      key: "status",
      title: $t("page.admin.library.bookCategory.categoryStatus"),
      align: "center",
      width: 80,
      render: (row: Api.BookManage.BookCategory) => {
        if (row.status === null) return null;
        const tagMap: Record<Api.Common.EnableStatus, NaiveUI.ThemeColor> = {
          1: "success",
          2: "warning",
        };
        return (
          <NTag type={tagMap[row.status]}>
            {$t(enableStatusRecord[row.status])}
          </NTag>
        );
      },
    },
    {
      key: "operate",
      title: $t("common.operate"),
      align: "center",
      width: 260,
      render: (row: Api.BookManage.BookCategory) => (
        <div class="flex-center justify-end gap-8px">
          <NButton
            type="primary"
            ghost
            size="small"
            onClick={() => handleAddChild(row)}
          >
            {$t("page.admin.library.bookCategory.addChildCategory")}
          </NButton>
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

const operateType = ref<OperateType>("add");
const editingData: Ref<Api.BookManage.BookCategory | null> = ref(null);

function handleSearch() {
  getDataByPage(1);
}

function handleReset() {
  isHotValue.value = null
  searchParams.value = {
    current: 1,
    size: 10,
    categoryName: null,
    categoryCode: null,
    isHot: null,
    status: null,
  };
}

function handleAdd() {
  operateType.value = "add";
  editingData.value = null;
  openModal();
}

function handleEdit(item: Api.BookManage.BookCategory) {
  operateType.value = "edit";
  editingData.value = { ...item };
  openModal();
}

function handleAddChild(item: Api.BookManage.BookCategory) {
  operateType.value = "addChild";
  editingData.value = { ...item };
  openModal();
}

async function handleDelete(id: number) {
  const { error } = await fetchDeleteCategory(id);
  if (!error) onDeleted();
}

async function handleBatchDelete() {
  onBatchDeleted();
}
</script>

<template>
  <div class="flex-col-stretch gap-16px overflow-hidden lt-sm:overflow-auto">
    <NCard :bordered="false" size="small">
      <NSpace wrap :size="[12, 12]">
        <NInput v-model:value="searchParams.categoryName" :placeholder="$t('page.admin.library.bookCategory.categoryName')" clearable
          style="width: 160px" @keyup.enter="handleSearch" />
        <NInput v-model:value="searchParams.categoryCode" :placeholder="$t('page.admin.library.bookCategory.categoryCode')" clearable
          style="width: 140px" @keyup.enter="handleSearch" />
        <NSelect v-model:value="isHotValue" :placeholder="$t('page.admin.library.bookCategory.isHot')"
          :options="[{value: 'true', label: $t('common.yesOrNo.yes')}, {value: 'false', label: $t('common.yesOrNo.no')}]" clearable style="width: 130px" />
        <NSelect v-model:value="searchParams.status" :placeholder="$t('page.admin.library.bookCategory.categoryStatus')"
          :options="enableStatusOptions" clearable style="width: 130px" />
        <NButton type="primary" @click="handleSearch">{{ $t("common.search") }}</NButton>
        <NButton @click="handleReset">{{ $t("common.reset") }}</NButton>
      </NSpace>
    </NCard>
    <NCard
      :title="$t('page.admin.library.bookCategory.title')"
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
        :scroll-x="1060"
        :loading="loading"
        :row-key="(row) => row.id"
        remote
        :pagination="pagination"
        class="sm:h-full"
      />
      <CategoryOperateModal
        v-model:visible="visible"
        :operate-type="operateType"
        :row-data="editingData"
        @submitted="getDataByPage"
      />
    </NCard>
  </div>
</template>
