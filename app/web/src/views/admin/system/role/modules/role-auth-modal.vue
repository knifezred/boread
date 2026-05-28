<script setup lang="ts">
import { h, computed, ref, watch } from 'vue';
import { NDataTable, NCheckbox, NSwitch, NSpace, NButton } from 'naive-ui';
import { fetchGetMenuTree, fetchGetRoleMenuIds, fetchGetRoleButtonIds, fetchGrantRoleMenus, fetchGrantRoleButtons } from '@/service/api';
import { $t } from '@/locales';

defineOptions({ name: 'RoleAuthModal' });

interface Props { roleId: number; }
const props = defineProps<Props>();
const visible = defineModel<boolean>('visible', { default: false });
const loading = ref(false);
const expandAll = ref(false);
const expandedRowKeys = ref<number[]>([]);
const allRowKeys = ref<number[]>([]);

const title = computed(() => $t('page.manage.role.menuAuth'));

interface TableTreeNode {
  id: number;
  menuName: string;
  buttons: Api.SystemManage.SysMenuButton[];
  children?: TableTreeNode[];
}

const tableTree = ref<TableTreeNode[]>([]);
const checks = ref<number[]>([]);

function collectAllKeys(nodes: TableTreeNode[]) {
  let keys: number[] = [];
  for (const n of nodes) {
    keys.push(n.id);
    if (n.children) keys = keys.concat(collectAllKeys(n.children));
  }
  return keys;
}

function buildTableTree(menuNodes: Api.SystemManage.MenuTreeNode[]): TableTreeNode[] {
  return menuNodes.map(n => {
    const node: TableTreeNode = {
      id: n.id,
      menuName: n.menuName,
      buttons: n.buttons ?? [],
    };
    if (n.children && n.children.length > 0) {
      node.children = buildTableTree(n.children);
    }
    return node;
  });
}

function getDescendantIds(node: TableTreeNode): number[] {
  let ids: number[] = [];
  if (node.children) {
    for (const child of node.children) {
      ids.push(child.id);
      ids = ids.concat(getDescendantIds(child));
    }
  }
  return ids;
}

function toggleMenu(id: number, checked: boolean) {
  if (checked) {
    if (!checks.value.includes(id)) checks.value.push(id);
    // 勾选父菜单自动勾选所有子菜单
    const node = findNodeById(tableTree.value, id);
    if (node) {
      const descendants = getDescendantIds(node);
      for (const did of descendants) {
        if (!checks.value.includes(did)) checks.value.push(did);
      }
    }
  } else {
    checks.value = checks.value.filter(k => k !== id);
    // 取消父菜单自动取消所有子菜单
    const node = findNodeById(tableTree.value, id);
    if (node) {
      const descendants = getDescendantIds(node);
      checks.value = checks.value.filter(k => !descendants.includes(k));
    }
  }
}

function toggleButton(buttonId: number, checked: boolean) {
  const key = -buttonId;
  if (checked) {
    if (!checks.value.includes(key)) checks.value.push(key);
  } else {
    checks.value = checks.value.filter(k => k !== key);
  }
}

function findNodeById(nodes: TableTreeNode[], id: number): TableTreeNode | null {
  for (const n of nodes) {
    if (n.id === id) return n;
    if (n.children) {
      const child = findNodeById(n.children, id);
      if (child) return child;
    }
  }
  return null;
}

const columns = computed(() => [
  {
    title: $t('page.manage.role.menuName'),
    key: 'menuName',
    width: 220,
    treeNode: true,
    render: (row: TableTreeNode) =>
      h(NCheckbox, {
        label: row.menuName,
        value: checks.value.includes(row.id),
        onUpdateValue: (val: boolean) => toggleMenu(row.id, val),
      }),
  },
  {
    title: $t('page.manage.role.buttonName'),
    key: 'buttons',
    render: (row: TableTreeNode) =>
      row.buttons?.length
        ? h(NSpace, { size: 'small', wrap: true }, () =>
            row.buttons.map(b =>
              h(NCheckbox, {
                label: b.buttonDesc ?? b.buttonCode,
                value: checks.value.includes(-b.id),
                onUpdateValue: (val: boolean) => toggleButton(b.id, val),
              }),
            ),
          )
        : null,
  },
]);

async function init() {
  loading.value = true;
  try {
    const [treeRes, menuIdsRes, buttonIdsRes] = await Promise.all([
      fetchGetMenuTree(),
      fetchGetRoleMenuIds(props.roleId),
      fetchGetRoleButtonIds(props.roleId),
    ]);
    if (!treeRes.error) {
      tableTree.value = buildTableTree(treeRes.data);
      allRowKeys.value = collectAllKeys(tableTree.value);
    }
    const checked: number[] = [];
    if (!menuIdsRes.error) checked.push(...menuIdsRes.data.map(Number));
    if (!buttonIdsRes.error) checked.push(...buttonIdsRes.data.map(id => -Number(id)));
    checks.value = checked;
    // 默认全展开
    expandAll.value = true;
    expandedRowKeys.value = allRowKeys.value;
  } finally {
    loading.value = false;
  }
}

function closeModal() { visible.value = false; }

function handleExpandChange(val: boolean) {
  expandedRowKeys.value = val ? allRowKeys.value : [];
}

async function handleSubmit() {
  const menuIds = checks.value.filter(k => k > 0).map(Number);
  const buttonIds = checks.value.filter(k => k < 0).map(k => -Number(k));
  const [menuRes, buttonRes] = await Promise.all([
    fetchGrantRoleMenus(props.roleId, menuIds),
    fetchGrantRoleButtons(props.roleId, buttonIds),
  ]);
  if (menuRes.error || buttonRes.error) {
    window.$message?.error?.(menuRes.error?.message || buttonRes.error?.message || $t('common.operateFail'));
    return;
  }
  window.$message?.success?.($t('common.modifySuccess'));
  closeModal();
}

watch(visible, val => { if (val) init(); });
</script>

<template>
  <NModal v-model:show="visible" :title="title" preset="card" class="w-800px" :loading="loading">
    <div class="flex items-center mb-12px">
      <span class="mr-12px">{{ $t('page.manage.role.expandOrCollapse') }}</span>
      <NSwitch v-model:value="expandAll" @update:value="handleExpandChange" />
    </div>
    <NDataTable
      :columns="columns"
      :data="tableTree"
      :expanded-row-keys="expandedRowKeys"
      :row-key="(row) => row.id"
      :loading="loading"
      :pagination="false"
      bordered
      single-line
      :scroll="{ y: 440, x: 750 }"
    />
    <template #footer>
      <NSpace justify="end">
        <NButton size="small" @click="closeModal">{{ $t('common.cancel') }}</NButton>
        <NButton type="primary" size="small" @click="handleSubmit">{{ $t('common.confirm') }}</NButton>
      </NSpace>
    </template>
  </NModal>
</template>

<style scoped>
:deep(.n-modal-body-wrapper) {
  overflow: hidden !important;
}
</style>