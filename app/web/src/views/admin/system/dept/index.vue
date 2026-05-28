<script setup lang="tsx">
import { ref } from 'vue';
import type { Ref } from 'vue';
import { NButton, NPopconfirm, NTag } from 'naive-ui';
import { useBoolean } from '@sa/hooks';
import { enableStatusRecord } from '@/constants/business';
import { fetchGetDeptList, fetchDeleteDept } from '@/service/api';
import { useAppStore } from '@/store/modules/app';
import { defaultTransform, useNaivePaginatedTable, useTableOperate } from '@/hooks/common/table';
import { $t } from '@/locales';
import DeptOperateModal, { type OperateType } from './modules/dept-operate-modal.vue';

const appStore = useAppStore();

const { bool: visible, setTrue: openModal } = useBoolean();

const wrapperRef = ref<HTMLElement | null>(null);

const searchParams = ref<Api.SystemManage.DeptSearchParams>({
  current: 1,
  size: 10,
  deptName: null,
  status: null,
});

const { columns, columnChecks, data, loading, pagination, getData, getDataByPage } = useNaivePaginatedTable({
  api: () => fetchGetDeptList(searchParams.value),
  onPaginationParamsChange: params => {
    searchParams.value.current = params.page;
    searchParams.value.size = params.pageSize;
  },
  transform: response => defaultTransform(response),
  columns: () => [
    {
      type: 'selection',
      align: 'center',
      width: 48
    },
    {
      key: 'deptName',
      title: $t('page.manage.dept.deptName'),
      align: 'left',
    },
    {
      key: 'deptCode',
      title: $t('page.manage.dept.deptCode'),
      align: 'center',
      width: 120
    },
    {
      key: 'leader',
      title: $t('page.manage.dept.leader'),
      align: 'center',
      width: 100
    },
    {
      key: 'parentId',
      title: $t('page.manage.dept.parentId'),
      align: 'center',
      width: 90
    },
    {
      key: 'sortOrder',
      title: $t('page.manage.dept.sortOrder'),
      align: 'center',
      width: 80
    },
    {
      key: 'status',
      title: $t('page.manage.dept.status'),
      align: 'center',
      width: 80,
      render: (row: Api.SystemManage.Dept) => {
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
      width: 230,
      render: (row: Api.SystemManage.Dept) => (
        <div class="flex-center justify-end gap-8px">
          <NButton type="primary" ghost size="small" onClick={() => handleAddChildDept(row)}>
            {$t('page.manage.dept.addChildDept')}
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
  // request
  console.log(checkedRowKeys.value);

  onBatchDeleted();
}

async function handleDelete(id: number) {
  const { error } = await fetchDeleteDept(id);
  if (!error) {
    onDeleted();
  }
}

/** the edit dept data or the parent dept data when adding a child dept */
const editingData: Ref<Api.SystemManage.Dept | null> = ref(null);

function handleEdit(item: Api.SystemManage.Dept) {
  operateType.value = 'edit';
  editingData.value = { ...item };

  openModal();
}

function handleAddChildDept(item: Api.SystemManage.Dept) {
  operateType.value = 'addChild';

  editingData.value = { ...item };

  openModal();
}
</script>

<template>
  <div ref="wrapperRef" class="flex-col-stretch gap-16px overflow-hidden lt-sm:overflow-auto">
    <NCard :title="$t('page.manage.dept.title')" :bordered="false" size="small" class="card-wrapper sm:flex-1-hidden">
      <template #header-extra>
        <TableHeaderOperation v-model:columns="columnChecks" :disabled-delete="checkedRowKeys.length === 0"
          :loading="loading" @add="handleAdd" @delete="handleBatchDelete" @refresh="getData" />
      </template>
      <NDataTable v-model:checked-row-keys="checkedRowKeys" :columns="columns" :data="data" size="small"
        :flex-height="!appStore.isMobile" :scroll-x="1088" :loading="loading" :row-key="row => row.id" remote
        :pagination="pagination" class="sm:h-full" />
      <DeptOperateModal v-model:visible="visible" :operate-type="operateType" :row-data="editingData"
        @submitted="getDataByPage" />
    </NCard>
  </div>
</template>