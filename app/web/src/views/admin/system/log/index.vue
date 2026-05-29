<script setup lang="tsx">
import { ref } from 'vue';
import { NButton } from 'naive-ui';
import { fetchGetLoginLogList, fetchGetOperationLogList } from '@/service/api';
import { useAppStore } from '@/store/modules/app';
import { defaultTransform, useNaivePaginatedTable } from '@/hooks/common/table';
import { $t } from '@/locales';

const appStore = useAppStore();

const activeTab = ref<'login' | 'operation'>('login');

const loginSearchParams = ref<Api.SystemManage.LoginLogSearchParams>({
  current: 1,
  size: 10,
  userName: null,
  loginIp: null,
  loginType: null,
  loginResult: null,
  startTime: null,
  endTime: null,
});

const operationSearchParams = ref<Api.SystemManage.OperationLogSearchParams>({
  current: 1,
  size: 10,
  userName: null,
  module: null,
  action: null,
  clientIp: null,
  startTime: null,
  endTime: null,
});

const {
  columns: loginColumns,
  data: loginData,
  loading: loginLoading,
  pagination: loginPagination,
  getData: getLoginData,
} = useNaivePaginatedTable({
  api: () => fetchGetLoginLogList(loginSearchParams.value),
  transform: response => defaultTransform(response),
  onPaginationParamsChange: params => {
    loginSearchParams.value.current = params.page;
    loginSearchParams.value.size = params.pageSize;
  },
  columns: () => [
    {
      key: 'index',
      title: $t('common.index'),
      align: 'center',
      width: 64,
      render: (_, index) => index + 1
    },
    {
      key: 'userName',
      title: $t('page.admin.system.log.userName'),
      align: 'center',
      minWidth: 100
    },
    {
      key: 'loginIp',
      title: $t('page.admin.system.log.loginIp'),
      align: 'center',
      width: 140
    },
    {
      key: 'loginType',
      title: $t('page.admin.system.log.loginType'),
      align: 'center',
      width: 100
    },
    {
      key: 'loginResult',
      title: $t('page.admin.system.log.loginResult'),
      align: 'center',
      width: 100
    },
    {
      key: 'loginTime',
      title: $t('common.updateTime'),
      align: 'center',
      width: 180
    }
  ]
});

const {
  columns: operationColumns,
  data: operationData,
  loading: operationLoading,
  pagination: operationPagination,
  getData: getOperationData,
} = useNaivePaginatedTable({
  api: () => fetchGetOperationLogList(operationSearchParams.value),
  transform: response => defaultTransform(response),
  onPaginationParamsChange: params => {
    operationSearchParams.value.current = params.page;
    operationSearchParams.value.size = params.pageSize;
  },
  columns: () => [
    {
      key: 'index',
      title: $t('common.index'),
      align: 'center',
      width: 64,
      render: (_, index) => index + 1
    },
    {
      key: 'userName',
      title: $t('page.admin.system.log.userName'),
      align: 'center',
      minWidth: 100
    },
    {
      key: 'module',
      title: $t('page.admin.system.log.module'),
      align: 'center',
      width: 120
    },
    {
      key: 'action',
      title: $t('page.admin.system.log.action'),
      align: 'center',
      width: 120
    },
    {
      key: 'clientIp',
      title: $t('page.admin.system.log.clientIp'),
      align: 'center',
      width: 140
    },
    {
      key: 'operateTime',
      title: $t('common.updateTime'),
      align: 'center',
      width: 180
    }
  ]
});
</script>

<template>
  <div class="min-h-500px flex-col-stretch gap-16px overflow-hidden lt-sm:overflow-auto">
    <NCard :bordered="false" size="small" class="card-wrapper">
      <NTabs v-model:value="activeTab" type="line" animated>
        <NTabPane name="login" :tab="$t('page.admin.system.log.loginLog')" />
        <NTabPane name="operation" :tab="$t('page.admin.system.log.operationLog')" />
      </NTabs>
    </NCard>

    <NCard v-if="activeTab === 'login'" :title="$t('page.admin.system.log.loginLog')" :bordered="false" size="small" class="card-wrapper sm:flex-1-hidden">
      <template #header-extra>
        <NButton size="small" @click="getLoginData">
          <template #icon><icon-ic-round-refresh class="text-icon" /></template>
          {{ $t('common.refresh') }}
        </NButton>
      </template>
      <NDataTable
        :columns="loginColumns"
        :data="loginData"
        size="small"
        :flex-height="!appStore.isMobile"
        :scroll-x="962"
        :loading="loginLoading"
        remote
        :row-key="row => row.id"
        :pagination="loginPagination"
        class="sm:h-full"
      />
    </NCard>

    <NCard v-if="activeTab === 'operation'" :title="$t('page.admin.system.log.operationLog')" :bordered="false" size="small" class="card-wrapper sm:flex-1-hidden">
      <template #header-extra>
        <NButton size="small" @click="getOperationData">
          <template #icon><icon-ic-round-refresh class="text-icon" /></template>
          {{ $t('common.refresh') }}
        </NButton>
      </template>
      <NDataTable
        :columns="operationColumns"
        :data="operationData"
        size="small"
        :flex-height="!appStore.isMobile"
        :scroll-x="962"
        :loading="operationLoading"
        remote
        :row-key="row => row.id"
        :pagination="operationPagination"
        class="sm:h-full"
      />
    </NCard>
  </div>
</template>