<script setup lang="tsx">
import { computed, ref, onMounted } from 'vue';
import { NButton, NCard, NDataTable, NDrawer, NDrawerContent, NPopconfirm, NSpace, NTag } from 'naive-ui';
import { useBoolean } from '@sa/hooks';
import { enableStatusRecord } from '@/constants/business';
import { fetchGetDictItems, fetchDeleteDictItem } from '@/service/api';
import { $t } from '@/locales';
import DictItemOperateModal from './dict-item-operate-modal.vue';

defineOptions({
  name: 'DictItemManageDrawer'
});

interface Props {
  dictId: number;
  dictName: string;
}

const props = defineProps<Props>();

interface Emits {
  (e: 'close'): void;
}

const emit = defineEmits<Emits>();

const visible = defineModel<boolean>('visible', {
  default: false
});

const { bool: modalVisible, setTrue: openModal, setFalse: closeModal } = useBoolean();

const items = ref<Api.SystemManage.DictItem[]>([]);
const loading = ref(false);

async function loadItems() {
  loading.value = true;
  const { error, data } = await fetchGetDictItems(props.dictId);
  if (!error) {
    items.value = data;
  }
  loading.value = false;
}

function closeDrawer() {
  visible.value = false;
  emit('close');
}

const operateType = ref<NaiveUI.TableOperateType>('add');
const editingItem = ref<Api.SystemManage.DictItem | null>(null);

function handleAdd() {
  operateType.value = 'add';
  editingItem.value = null;
  openModal();
}

function handleEdit(item: Api.SystemManage.DictItem) {
  operateType.value = 'edit';
  editingItem.value = { ...item };
  openModal();
}

async function handleDelete(id: number) {
  const { error } = await fetchDeleteDictItem(id);
  if (!error) {
    window.$message?.success($t('common.deleteSuccess'));
    await loadItems();
  }
}

function onSubmitted() {
  closeModal();
  loadItems();
}

const columns = computed(() => [
  {
    key: 'index',
    title: $t('common.index'),
    align: 'center' as const,
    width: 64,
    render: (_: any, index: number) => index + 1
  },
  {
    key: 'itemLabel',
    title: $t('page.admin.system.dictItem.itemLabel'),
    align: 'center' as const,
    minWidth: 120
  },
  {
    key: 'itemValue',
    title: $t('page.admin.system.dictItem.itemValue'),
    align: 'center' as const,
    minWidth: 120
  },
  {
    key: 'itemDesc',
    title: $t('page.admin.system.dictItem.itemDesc'),
    minWidth: 160,
    render: (row: Api.SystemManage.DictItem) => row.itemDesc || '-'
  },
  {
    key: 'sortOrder',
    title: $t('page.admin.system.dictItem.sortOrder'),
    align: 'center' as const,
    width: 80
  },
  {
    key: 'status',
    title: $t('page.admin.system.dictItem.itemStatus'),
    align: 'center' as const,
    width: 100,
    render: (row: Api.SystemManage.DictItem) => {
      if (row.status === null) return null;
      const tagMap: Record<Api.Common.EnableStatus, NaiveUI.ThemeColor> = {
        1: 'success',
        2: 'warning'
      };
      return <NTag type={tagMap[row.status]}>{$t(enableStatusRecord[row.status])}</NTag>;
    }
  },
  {
    key: 'operate',
    title: $t('common.operate'),
    align: 'center' as const,
    width: 150,
    render: (row: Api.SystemManage.DictItem) => (
      <div class="flex-center gap-8px">
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
]);

onMounted(() => {
  loadItems();
});
</script>

<template>
  <NDrawer v-model:show="visible" display-directive="show" width="45%">
    <NDrawerContent :title="$t('page.admin.system.dictItem.title')" :native-scrollbar="false" closable @after-leave="closeDrawer">
      <NCard :bordered="false" size="small" :title="dictName" class="card-wrapper">
        <template #header-extra>
          <NSpace :size="12">
            <NButton size="small" @click="loadItems">
              <template #icon><icon-ic-round-refresh class="text-icon" /></template>
              {{ $t('common.refresh') }}
            </NButton>
            <NButton type="primary" size="small" @click="handleAdd">
              <template #icon><icon-ic-round-add class="text-icon" /></template>
              {{ $t('page.admin.system.dictItem.addDictItem') }}
            </NButton>
          </NSpace>
        </template>
        <NDataTable
          :columns="columns"
          :data="items"
          size="small"
          :loading="loading"
          :bordered="false"
          :row-key="row => row.id"
        />
      </NCard>
      <DictItemOperateModal
        v-model:visible="modalVisible"
        :operate-type="operateType"
        :row-data="editingItem"
        :dict-id="dictId"
        @submitted="onSubmitted"
      />
    </NDrawerContent>
  </NDrawer>
</template>