<script setup lang="tsx">
import { ref } from "vue"
import type { Ref } from "vue"
import { NButton, NPopconfirm, NTag } from "naive-ui"
import { useBoolean } from "@sa/hooks"
import {
  fetchGetChapterRuleList,
  fetchDeleteChapterRule,
} from "@/service/api"
import { useAppStore } from "@/store/modules/app"
import { defaultTransform, useNaivePaginatedTable, useTableOperate } from "@/hooks/common/table"
import { $t } from "@/locales"
import ChapterRuleOperateModal from "./modules/chapter-rule-operate-modal.vue"

defineOptions({ name: "BookChapterRule" });

const appStore = useAppStore();
const { bool: visible, setTrue: openModal } = useBoolean();

const searchParams = ref<Api.BookManage.ChapterRuleSearchParams>({
  current: 1, size: 10, ruleName: null, scopeType: null, bookId: null, status: null,
});

const { columns, columnChecks, data, loading, pagination, getData, getDataByPage } = useNaivePaginatedTable({
  api: () => fetchGetChapterRuleList(searchParams.value),
  onPaginationParamsChange: (params) => { searchParams.value.current = params.page; searchParams.value.size = params.pageSize; },
  transform: (response) => defaultTransform(response),
  columns: () => [
    { type: "selection", align: "center", width: 48 },
    { key: "index", title: $t("common.index"), align: "center", width: 64, render: (_, index) => index + 1 },
    { key: "ruleName", title: $t("page.admin.library.bookChapterRule.ruleName"), align: "center", minWidth: 140 },
    {
      key: "scopeType", title: $t("page.admin.library.bookChapterRule.scopeType"), align: "center", width: 100,
      render: (row: { scopeType: string }) => {
        const map: Record<string, string> = { "1": $t("page.admin.library.bookChapterRule.scopeGlobal"), "2": $t("page.admin.library.bookChapterRule.scopeBook") };
        return <NTag type={row.scopeType === "1" ? "info" : "primary"}>{map[row.scopeType] || row.scopeType}</NTag>;
      },
    },
    { key: "pattern", title: $t("page.admin.library.bookChapterRule.pattern"), align: "left", ellipsis: { tooltip: true }, minWidth: 200 },
    { key: "priority", title: $t("page.admin.library.bookChapterRule.priority"), align: "center", width: 80 },
    {
      key: "status", title: $t("page.admin.library.bookChapterRule.status"), align: "center", width: 80,
      render: (row: { status: string }) => <NTag type={row.status === "1" ? "success" : "warning"}>{row.status === "1" ? $t("common.enable") : $t("common.disable")}</NTag>,
    },
    {
      key: "operate", title: $t("common.operate"), align: "center", width: 180,
      render: (row: Api.BookManage.BookChapterRule) => (
        <div class="flex-center gap-8px">
          <NButton type="primary" ghost size="small" onClick={() => handleEdit(row)}>{$t("common.edit")}</NButton>
          <NPopconfirm onPositiveClick={() => handleDelete(row.id)}>
            {{ default: () => $t("common.confirmDelete"), trigger: () => <NButton type="error" ghost size="small">{$t("common.delete")}</NButton> }}
          </NPopconfirm>
        </div>
      ),
    },
  ],
});

const { checkedRowKeys, onBatchDeleted, onDeleted } = useTableOperate(data, "id", getData);
const operateType = ref<NaiveUI.TableOperateType>("add");
const editingData: Ref<Api.BookManage.BookChapterRule | null> = ref(null);

function handleAdd() { operateType.value = "add"; editingData.value = null; openModal(); }
function handleEdit(item: Api.BookManage.BookChapterRule) { operateType.value = "edit"; editingData.value = { ...item }; openModal(); }
async function handleDelete(id: number) { const { error } = await fetchDeleteChapterRule(id); if (!error) onDeleted(); }
async function handleBatchDelete() { onBatchDeleted(); }
</script>

<template>
  <div class="flex-col-stretch gap-16px overflow-hidden lt-sm:overflow-auto">
    <NCard :title="$t('page.admin.library.bookChapterRule.title')" :bordered="false" size="small" class="card-wrapper sm:flex-1-hidden">
      <template #header-extra>
        <TableHeaderOperation v-model:columns="columnChecks" :disabled-delete="checkedRowKeys.length === 0" :loading="loading" @add="handleAdd" @delete="handleBatchDelete" @refresh="getData" />
      </template>
      <NDataTable v-model:checked-row-keys="checkedRowKeys" :columns="columns" :data="data" size="small" :flex-height="!appStore.isMobile" :scroll-x="900" :loading="loading" :row-key="(row) => row.id" remote :pagination="pagination" class="sm:h-full" />
      <ChapterRuleOperateModal v-model:visible="visible" :operate-type="operateType" :row-data="editingData" @submitted="getDataByPage" />
    </NCard>
  </div>
</template>

<style scoped></style>
