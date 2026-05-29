<script setup lang="tsx">
import { ref } from "vue";
import { NButton, NPopconfirm, NTag } from "naive-ui";
import { enableStatusRecord, dataScopeRecord } from "@/constants/business";
import { fetchGetRoleList } from "@/service/api";
import { useAppStore } from "@/store/modules/app";
import {
  defaultTransform,
  useNaivePaginatedTable,
  useTableOperate,
} from "@/hooks/common/table";
import { $t } from "@/locales";
import RoleOperateDrawer from "./modules/role-operate-drawer.vue";
import RoleSearch from "./modules/role-search.vue";
import RoleAuthModal from "./modules/role-auth-modal.vue";
import { useBoolean } from "@sa/hooks";

const appStore = useAppStore();

const searchParams = ref<Api.SystemManage.RoleSearchParams>({
  current: 1,
  size: 10,
  roleName: null,
  roleCode: null,
  status: null,
});

const { bool: authVisible, setTrue: openAuthModal } = useBoolean();
const currentAuthRoleId = ref<number>(0);
const currentAuthRoleName = ref<string>("");

const {
  columns,
  columnChecks,
  data,
  loading,
  getData,
  getDataByPage,
  mobilePagination,
} = useNaivePaginatedTable({
  api: () => fetchGetRoleList(searchParams.value),
  transform: (response) => defaultTransform(response),
  onPaginationParamsChange: (params) => {
    searchParams.value.current = params.page;
    searchParams.value.size = params.pageSize;
  },
  columns: () => [
    {
      type: "selection",
      align: "center",
      width: 48,
    },
    {
      key: "index",
      title: $t("common.index"),
      width: 64,
      align: "center",
      render: (_, index) => index + 1,
    },
    {
      key: "roleName",
      title: $t("page.manage.role.roleName"),
      align: "center",
      minWidth: 120,
    },
    {
      key: "roleCode",
      title: $t("page.manage.role.roleCode"),
      align: "center",
      minWidth: 120,
    },
    {
      key: "roleDesc",
      title: $t("page.manage.role.roleDesc"),
      minWidth: 120,
    },
    {
      key: "dataScope",
      title: $t("page.manage.role.dataScope.title"),
      minWidth: 120,
      render: (row) => {
        if (row.dataScope === null) {
          return null;
        }

        const tagMap: Record<Api.SystemManage.DataScope, NaiveUI.ThemeColor> = {
          1: "error",
          2: "warning",
          3: "info",
          4: "primary",
          5: "default",
        };

        const label = $t(dataScopeRecord[row.dataScope]);

        return <NTag type={tagMap[row.dataScope]}>{label}</NTag>;
      },
    },
    {
      key: "status",
      title: $t("page.manage.role.roleStatus"),
      align: "center",
      width: 100,
      render: (row) => {
        if (row.status === null) {
          return null;
        }

        const tagMap: Record<Api.Common.EnableStatus, NaiveUI.ThemeColor> = {
          1: "success",
          2: "warning",
        };

        const label = $t(enableStatusRecord[row.status]);

        return <NTag type={tagMap[row.status]}>{label}</NTag>;
      },
    },
    {
      key: "operate",
      title: $t("common.operate"),
      align: "center",
      width: 220,
      render: (row) => (
        <div class="flex-center gap-8px">
          <NButton
            type="primary"
            ghost
            size="small"
            onClick={() => edit(row.id)}
          >
            {$t("common.edit")}
          </NButton>
          <NButton
            type="info"
            ghost
            size="small"
            onClick={() => {
              currentAuthRoleId.value = row.id;
              currentAuthRoleName.value = row.roleName;
              openAuthModal();
            }}
          >
            {$t("page.manage.role.menuAuth")}
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

const {
  drawerVisible,
  operateType,
  editingData,
  handleAdd,
  handleEdit,
  checkedRowKeys,
  onBatchDeleted,
  onDeleted,
  // closeDrawer
} = useTableOperate(data, "id", getData);

async function handleBatchDelete() {
  // request
  console.log(checkedRowKeys.value);

  onBatchDeleted();
}

function handleDelete(id: number) {
  // request
  console.log(id);

  onDeleted();
}

function edit(id: number) {
  handleEdit(id);
}
</script>

<template>
  <div
    class="min-h-500px flex-col-stretch gap-16px overflow-hidden lt-sm:overflow-auto"
  >
    <RoleSearch v-model:model="searchParams" @search="getDataByPage" />
    <NCard
      :title="$t('page.manage.role.title')"
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
        :scroll-x="702"
        :loading="loading"
        remote
        :row-key="(row) => row.id"
        :pagination="mobilePagination"
        class="sm:h-full"
      />
      <RoleOperateDrawer
        v-model:visible="drawerVisible"
        :operate-type="operateType"
        :row-data="editingData"
        @submitted="getDataByPage"
      />
      <RoleAuthModal
        v-model:visible="authVisible"
        :role-id="currentAuthRoleId"
        :role-name="currentAuthRoleName"
        @submitted="getDataByPage"
      />
    </NCard>
  </div>
</template>

<style scoped></style>
