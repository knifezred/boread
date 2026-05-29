<script setup lang="tsx">
import { ref } from "vue";
import type { Ref } from "vue";
import { NButton, NPopconfirm, NTag } from "naive-ui";
import { useBoolean } from "@sa/hooks";
import { enableStatusRecord } from "@/constants/business";
import { fetchGetCategoryList, fetchDeleteCategory } from "@/service/api";
import { useAppStore } from "@/store/modules/app";
import {
  defaultTransform,
  useNaivePaginatedTable,
  useTableOperate,
} from "@/hooks/common/table";
import { $t } from "@/locales";
import CategoryOperateModal, {
  type OperateType,
} from "./modules/category-operate-modal.vue";

const appStore = useAppStore();

const { bool: visible, setTrue: openModal } = useBoolean();

const searchParams = ref<Api.SystemManage.CategorySearchParams>({
  current: 1,
  size: 10,
  categoryName: null,
  categoryCode: null,
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
  api: () => fetchGetCategoryList(searchParams.value),
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
      title: $t("page.manage.bookCategory.categoryName"),
      align: "left",
      width: 240,
    },
    {
      key: "description",
      title: $t("page.manage.bookCategory.description"),
      align: "left",
      ellipsis: { tooltip: true },
      minWidth: 160,
    },
    {
      key: "categoryCode",
      title: $t("page.manage.bookCategory.categoryCode"),
      align: "center",
      width: 140,
    },
    {
      key: "parentId",
      title: $t("page.manage.bookCategory.parentId"),
      align: "center",
      width: 90,
    },
    {
      key: "sortOrder",
      title: $t("page.manage.bookCategory.sortOrder"),
      align: "center",
      width: 80,
    },
    {
      key: "status",
      title: $t("page.manage.bookCategory.categoryStatus"),
      align: "center",
      width: 80,
      render: (row: Api.SystemManage.BookCategory) => {
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
      render: (row: Api.SystemManage.BookCategory) => (
        <div class="flex-center justify-end gap-8px">
          <NButton
            type="primary"
            ghost
            size="small"
            onClick={() => handleAddChild(row)}
          >
            {$t("page.manage.bookCategory.addChildCategory")}
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
const editingData: Ref<Api.SystemManage.BookCategory | null> = ref(null);

function handleAdd() {
  operateType.value = "add";
  editingData.value = null;
  openModal();
}

function handleEdit(item: Api.SystemManage.BookCategory) {
  operateType.value = "edit";
  editingData.value = { ...item };
  openModal();
}

function handleAddChild(item: Api.SystemManage.BookCategory) {
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
    <NCard
      :title="$t('page.manage.bookCategory.title')"
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
        :scroll-x="960"
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
