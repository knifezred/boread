<script setup lang="tsx">
import { computed, ref, watch } from 'vue';
import { NButton, NForm, NFormItem, NInput, NInputNumber, NModal, NScrollbar, NSpace } from 'naive-ui';
import { enableStatusOptions } from '@/constants/business';
import { useFormRules, useNaiveForm } from '@/hooks/common/form';
import { $t } from '@/locales';
import { fetchCreateDictItem, fetchUpdateDictItem } from '@/service/api/system-manage';

defineOptions({
  name: 'DictItemOperateModal'
});

interface Props {
  /** the type of operation */
  operateType: NaiveUI.TableOperateType;
  /** the edit row data */
  rowData?: Api.SystemManage.DictItem | null;
  /** the parent dict id (required when adding) */
  dictId: number;
}

const props = defineProps<Props>();

interface Emits {
  (e: 'submitted'): void;
}

const emit = defineEmits<Emits>();

const visible = defineModel<boolean>('visible', {
  default: false
});

const { formRef, validate, restoreValidation } = useNaiveForm();
const { defaultRequiredRule } = useFormRules();

const title = computed(() => {
  const titles: Record<NaiveUI.TableOperateType, string> = {
    add: $t('page.admin.system.dictItem.addDictItem'),
    edit: $t('page.admin.system.dictItem.editDictItem')
  };
  return titles[props.operateType];
});

type Model = Pick<Api.SystemManage.DictItem, 'id' | 'dictId' | 'itemLabel' | 'itemValue' | 'itemDesc' | 'sortOrder' | 'status'>;

const model = ref(createDefaultModel());

function createDefaultModel(): Model {
  return {
    id: 0,
    dictId: props.dictId,
    itemLabel: '',
    itemValue: '',
    itemDesc: '',
    sortOrder: 0,
    status: '1'
  };
}

type RuleKey = Extract<keyof Model, 'itemLabel' | 'itemValue' | 'status'>;

const rules: Record<RuleKey, App.Global.FormRule> = {
  itemLabel: defaultRequiredRule,
  itemValue: defaultRequiredRule,
  status: defaultRequiredRule
};

function closeModal() {
  visible.value = false;
}

async function handleSubmit() {
  await validate();

  const { id, dictId, ...restData } = model.value;
  let error = null;

  if (props.operateType === 'edit') {
    ({ error } = await fetchUpdateDictItem(id, { dictId, ...restData }));
  } else {
    ({ error } = await fetchCreateDictItem({ dictId, ...restData }));
  }

  if (error) {
    window.$message?.error(error.message);
  } else {
    window.$message?.success($t('common.updateSuccess'));
  }

  closeModal();
  emit('submitted');
}

function handleInitModel() {
  model.value = createDefaultModel();

  if (props.operateType === 'edit' && props.rowData) {
    Object.assign(model.value, props.rowData);
  }
}

watch(visible, () => {
  if (visible.value) {
    handleInitModel();
    restoreValidation();
  }
});
</script>

<template>
  <NModal v-model:show="visible" :title="title" preset="card" class="w-600px">
    <NScrollbar class="h-360px pr-20px">
      <NForm ref="formRef" :model="model" :rules="rules" label-placement="left" :label-width="100">
        <NFormItem :label="$t('page.admin.system.dictItem.itemLabel')" path="itemLabel">
          <NInput v-model:value="model.itemLabel" :placeholder="$t('page.admin.system.dictItem.form.itemLabel')" />
        </NFormItem>
        <NFormItem :label="$t('page.admin.system.dictItem.itemValue')" path="itemValue">
          <NInput v-model:value="model.itemValue" :placeholder="$t('page.admin.system.dictItem.form.itemValue')" />
        </NFormItem>
        <NFormItem :label="$t('page.admin.system.dictItem.itemDesc')" path="itemDesc">
          <NInput v-model:value="model.itemDesc" :placeholder="$t('page.admin.system.dictItem.form.itemDesc')" />
        </NFormItem>
        <NFormItem :label="$t('page.admin.system.dictItem.sortOrder')" path="sortOrder">
          <NInputNumber v-model:value="model.sortOrder" class="w-full" />
        </NFormItem>
        <NFormItem :label="$t('page.admin.system.dictItem.itemStatus')" path="status">
          <NRadioGroup v-model:value="model.status">
            <NRadio v-for="item in enableStatusOptions" :key="item.value" :value="item.value" :label="$t(item.label)" />
          </NRadioGroup>
        </NFormItem>
      </NForm>
    </NScrollbar>
    <template #footer>
      <NSpace justify="end" :size="16">
        <NButton @click="closeModal">{{ $t('common.cancel') }}</NButton>
        <NButton type="primary" @click="handleSubmit">{{ $t('common.confirm') }}</NButton>
      </NSpace>
    </template>
  </NModal>
</template>