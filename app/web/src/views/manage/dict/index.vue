<script setup lang="tsx">
import { ref, type Ref, shallowRef } from 'vue';
import { NButton, NPopconfirm, NTag } from 'naive-ui';
import { useBoolean } from '@sa/hooks';
import { enableStatusRecord } from '@/constants/business';
import { fetchGetDictList, fetchDeleteDict } from '@/service/api';
import { useAppStore } from '@/store/modules/app';
import { defaultTransform, useNaivePaginatedTable, useTableOperate } from '@/hooks/common/table';
import { $t } from '@/locales';
import DictOperateModal, { type OperateType } from './modules/dict-operate-modal.vue';
import DictSearch from './modules/dict-search.vue';
import DictItemManageDrawer from './modules/dict-item-manage-drawer.vue';

const appStore = useAppStore();

const { bool: visible, setTrue: openModal } = useBoolean();

const { bool: itemDrawerVisible, setTrue: openItemDrawer, setFalse: closeItemDrawer } = useBoolean();

const searchParams = ref<Api.SystemManage.DictSearchParams>({
  current: 1,
  size: 10,
  dictName: null,
  dictCode: null,
  status: null,
});

const { columns, columnChecks, data, getData, getDataByPage, loading, mobilePagination } = useNaivePaginatedTable({
  api: () => fetchGetDictList(searchParams.value),
  transform: response => defaultTransform(response),
  onPaginationParamsChange: params => {
    searchParams.value.current = params.page;
    searchParams.value.size = params.pageSize;
  },
  columns: () => [
    {
      type: 'selection',
      align: 'center',
      width: 48
    },
    {
      key: 'index',
      title: $t('common.index'),
      align: 'center',
      width: 64,
      render: (_, index) => index + 1
    },
    {
      key: 'dictName',
      title: $t('page.manage.dict.dictName'),
      align: 'center',
      minWidth: 120
    },
    {
      key: 'dictCode',
      title: $t('page.manage.dict.dictCode'),
      align: 'center',
      minWidth: 120
    },
    {
      key: 'dictDesc',
      title: $t('page.manage.dict.dictDesc'),
      minWidth: 160
    },
    {
      key: 'status',
      title: $t('page.manage.dict.dictStatus'),
      align: 'center',
      width: 100,
      render: row => {
        if (row.status === null) {
          return null;
        }

        const tagMap: Record<Api.Common.EnableStatus, NaiveUI.ThemeColor> = {
          1: 'success',
          2: 'warning'
        };

        const label = $t(enableStatusRecord[row.status]);

        return <NTag type={tagMap[row.status]}>{label}</NTag>;
      }
    },
    {
      key: 'operate',
      title: $t('common.operate'),
      align: 'center',
      width: 240,
      render: row => (
        <div class="flex-center gap-8px">
          <NButton type="primary" ghost size="small" onClick={() => handleItems(row)}>
            {$t('page.manage.dictItem.title')}
          </NButton>
          <NButton type="primary" ghost size="small" onClick={() => handleEdit(row)}>
            {$t('common.edit')}
          </NButton>
          <NPopconfirm onPositiveClick={() => handleDelete(row.id)}>
            {{
              default: () => $t('common.confirmDelete'),
              trigger: () => (
                <NButton type="error" ghost size="small">
                  {$t('common.delete')}
                </NButton>
              )
            }}
          </NPopconfirm>
        </div>
      )
    }
  ]
});

const { checkedRowKeys, onBatchDeleted, onDeleted } = useTableOperate(data, 'id', getData);

const operateType = ref<OperateType>('add');

function handleAdd() {
  operateType.value = 'add';
  openModal();
}

async function handleBatchDelete() {
  console.log(checkedRowKeys.value);
  onBatchDeleted();
}

async function handleDelete(id: number) {
  const { error } = await fetchDeleteDict(id);
  if (!error) {
    onDeleted();
  }
}

const editingData: Ref<Api.SystemManage.Dict | null> = ref(null);

function handleEdit(item: Api.SystemManage.Dict) {
  operateType.value = 'edit';
  editingData.value = { ...item };
  openModal();
}

const currentDictItemData = shallowRef<{ dictId: number; dictName: string } | null>(null);

function handleItems(item: Api.SystemManage.Dict) {
  currentDictItemData.value = { dictId: item.id, dictName: item.dictName };
  openItemDrawer();
}
</script>

<template>
  <div class="min-h-500px flex-col-stretch gap-16px overflow-hidden lt-sm:overflow-auto">
    <DictSearch v-model:model="searchParams" @search="getDataByPage" />
    <NCard :title="$t('page.manage.dict.title')" :bordered="false" size="small" class="card-wrapper sm:flex-1-hidden">
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
        :scroll-x="962"
        :loading="loading"
        remote
        :row-key="row => row.id"
        :pagination="mobilePagination"
        class="sm:h-full"
      />
      <DictOperateModal
        v-model:visible="visible"
        :operate-type="operateType"
        :row-data="editingData"
        @submitted="getDataByPage"
      />
      <DictItemManageDrawer
        v-if="currentDictItemData"
        v-model:visible="itemDrawerVisible"
        :dict-id="currentDictItemData.dictId"
        :dict-name="currentDictItemData.dictName"
        @close="closeItemDrawer"
      />
    </NCard>
  </div>
</template>