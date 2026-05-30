<script setup lang="tsx">
import { ref } from "vue"
import type { Ref } from "vue"
import { NButton, NPopconfirm, NTag, NSpace } from "naive-ui"
import { useBoolean } from "@sa/hooks"
import {
  fetchGetFilterRuleList,
  fetchDeleteFilterRule,
} from "@/service/api"
import { useAppStore } from "@/store/modules/app"
import { defaultTransform, useNaivePaginatedTable, useTableOperate } from "@/hooks/common/table"
import { $t } from "@/locales"
import FilterRuleOperateModal from "./modules/filter-rule-operate-modal.vue"

defineOptions({ name: "BookFilterRule" });

const appStore = useAppStore();
const { bool: visible, setTrue: openModal } = useBoolean();

const searchParams = ref<Api.SystemManage.FilterRuleSearchParams>({
  current: 1, size: 10, ruleName: null, applyStage: null, category: null, status: null,
});

const actionMap: Record<string, string> = {
  "1": $t("page.admin.library.bookFilterRule.actionReplace"),
  "2": $t("page.admin.library.bookFilterRule.actionBlock"),
  "3": $t("page.admin.library.bookFilterRule.actionMark"),
};
const actionColor: Record<string, NaiveUI.ThemeColor> = { "1": "success", "2": "error", "3": "warning" };
const matchTypeMap: Record<string, string> = { "1": $t("page.admin.library.bookFilterRule.matchKeyword"), "2": $t("page.admin.library.bookFilterRule.matchRegex") };

const { columns, columnChecks, data, loading, pagination, getData, getDataByPage } = useNaivePaginatedTable({
  api: () => fetchGetFilterRuleList(searchParams.value),
  onPaginationParamsChange: (params) => { searchParams.value.current = params.page; searchParams.value.size = params.pageSize; },
  transform: (response) => defaultTransform(response),
  columns: () => [
    { type: "selection", align: "center", width: 48 },
    { key: "index", title: $t("common.index"), align: "center", width: 64, render: (_, index) => index + 1 },
    { key: "ruleName", title: $t("page.admin.library.bookFilterRule.ruleName"), align: "center", minWidth: 140 },
    {
      key: "matchType", title: $t("page.admin.library.bookFilterRule.matchType"), align: "center", width: 80,
      render: (row: { matchType: string }) => <NTag type="info">{matchTypeMap[row.matchType] || row.matchType}</NTag>,
    },
    { key: "pattern", title: $t("page.admin.library.bookFilterRule.pattern"), align: "left", ellipsis: { tooltip: true }, minWidth: 180 },
    {
      key: "action", title: $t("page.admin.library.bookFilterRule.action"), align: "center", width: 80,
      render: (row: { action: string }) => <NTag type={actionColor[row.action]}>{actionMap[row.action] || row.action}</NTag>,
    },
    {
      key: "status", title: $t("page.admin.library.bookFilterRule.status"), align: "center", width: 80,
      render: (row: { status: string }) => <NTag type={row.status === "1" ? "success" : "warning"}>{row.status === "1" ? $t("common.enable") : $t("common.disable")}</NTag>,
    },
    {
      key: "operate", title: $t("common.operate"), align: "center", width: 180,
      render: (row: Api.SystemManage.BookContentFilterRule) => (
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
const editingData: Ref<Api.SystemManage.BookContentFilterRule | null> = ref(null);

function handleAdd() { operateType.value = "add"; editingData.value = null; openModal(); }
function handleEdit(item: Api.SystemManage.BookContentFilterRule) { operateType.value = "edit"; editingData.value = { ...item }; openModal(); }
async function handleDelete(id: number) { const { error } = await fetchDeleteFilterRule(id); if (!error) onDeleted(); }
async function handleBatchDelete() { onBatchDeleted(); }
</script>

<template>
  <div class="flex-col-stretch gap-16px overflow-hidden lt-sm:overflow-auto">
    <NCard :title="$t('page.admin.library.bookFilterRule.title')" :bordered="false" size="small" class="card-wrapper sm:flex-1-hidden">
      <template #header-extra>
        <TableHeaderOperation v-model:columns="columnChecks" :disabled-delete="checkedRowKeys.length === 0" :loading="loading" @add="handleAdd" @delete="handleBatchDelete" @refresh="getData" />
      </template>
      <NDataTable v-model:checked-row-keys="checkedRowKeys" :columns="columns" :data="data" size="small" :flex-height="!appStore.isMobile" :scroll-x="900" :loading="loading" :row-key="(row) => row.id" remote :pagination="pagination" class="sm:h-full" />
      <FilterRuleOperateModal v-model:visible="visible" :operate-type="operateType" :row-data="editingData" @submitted="getDataByPage" />
    </NCard>
  </div>
</template>

<style scoped></style>
