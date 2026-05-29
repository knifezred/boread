<script setup lang="tsx">
import { ref } from 'vue';
import type { Ref } from 'vue';
import { NButton, NPopconfirm } from 'naive-ui';
import { useBoolean } from '@sa/hooks';
import { fetchGetTagList, fetchDeleteTag } from '@/service/api';
import { useAppStore } from '@/store/modules/app';
import { defaultTransform, useNaivePaginatedTable, useTableOperate } from '@/hooks/common/table';
import { $t } from '@/locales';
import TagOperateModal from './modules/tag-operate-modal.vue';

const appStore = useAppStore();
const { bool: visible, setTrue: openModal } = useBoolean();

const searchParams = ref<Api.SystemManage.TagSearchParams>({
  current: 1, size: 10, tagName: null,
});

const { columns, columnChecks, data, loading, pagination, getData, getDataByPage } = useNaivePaginatedTable({
  api: () => fetchGetTagList(searchParams.value),
  onPaginationParamsChange: params => { searchParams.value.current = params.page; searchParams.value.size = params.pageSize; },
  transform: response => defaultTransform(response),
  columns: () => [
    { type: 'selection', align: 'center', width: 48 },
    { key: 'index', title: $t('common.index'), align: 'center', width: 64, render: (_, index) => index + 1 },
    { key: 'tagName', title: $t('page.admin.library.bookTag.tagName'), align: 'center', minWidth: 140 },
    { key: 'description', title: $t('page.admin.library.bookTag.description'), align: 'left', ellipsis: { tooltip: true }, minWidth: 160 },
    { key: 'usageCount', title: $t('page.admin.library.bookTag.usageCount'), align: 'center', width: 120 },
    {
      key: 'operate', title: $t('common.operate'), align: 'center', width: 180,
      render: (row: Api.SystemManage.BookTag) => (
        <div class="flex-center gap-8px">
          <NButton type="primary" ghost size="small" onClick={() => handleEdit(row)}>{$t('common.edit')}</NButton>
          <NPopconfirm onPositiveClick={() => handleDelete(row.id)}>
            {{ default: () => $t('common.confirmDelete'), trigger: () => <NButton type="error" ghost size="small">{$t('common.delete')}</NButton> }}
          </NPopconfirm>
        </div>
      )
    }
  ]
});

const { checkedRowKeys, onBatchDeleted, onDeleted } = useTableOperate(data, 'id', getData);

const operateType = ref<NaiveUI.TableOperateType>('add');
const editingData: Ref<Api.SystemManage.BookTag | null> = ref(null);

function handleAdd() { operateType.value = 'add'; editingData.value = null; openModal(); }
function handleEdit(item: Api.SystemManage.BookTag) { operateType.value = 'edit'; editingData.value = { ...item }; openModal(); }
async function handleDelete(id: number) { const { error } = await fetchDeleteTag(id); if (!error) onDeleted(); }
async function handleBatchDelete() { onBatchDeleted(); }
</script>

<template>
  <div class="flex-col-stretch gap-16px overflow-hidden lt-sm:overflow-auto">
    <NCard :title="$t('page.admin.library.bookTag.title')" :bordered="false" size="small" class="card-wrapper sm:flex-1-hidden">
      <template #header-extra>
        <TableHeaderOperation v-model:columns="columnChecks" :disabled-delete="checkedRowKeys.length === 0"
          :loading="loading" @add="handleAdd" @delete="handleBatchDelete" @refresh="getData" />
      </template>
      <NDataTable v-model:checked-row-keys="checkedRowKeys" :columns="columns" :data="data" size="small"
        :flex-height="!appStore.isMobile" :scroll-x="620" :loading="loading" :row-key="row => row.id" remote
        :pagination="pagination" class="sm:h-full" />
      <TagOperateModal v-model:visible="visible" :operate-type="operateType" :row-data="editingData" @submitted="getDataByPage" />
    </NCard>
  </div>
</template>