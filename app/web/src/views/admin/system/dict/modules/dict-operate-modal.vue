<script setup lang="tsx">
import { computed, ref, watch } from 'vue';
import { NButton, NForm, NFormItem, NInput, NModal, NScrollbar, NSpace } from 'naive-ui';
import { enableStatusOptions } from '@/constants/business';
import { useFormRules, useNaiveForm } from '@/hooks/common/form';
import { $t } from '@/locales';
import { fetchCreateDict, fetchUpdateDict } from '@/service/api';

defineOptions({
  name: 'DictOperateModal'
});

interface Props {
  /** the type of operation */
  operateType: NaiveUI.TableOperateType;
  /** the edit row data */
  rowData?: Api.SystemManage.Dict | null;
}

const props = defineProps<Props>();

interface Emits {
  (e: 'submitted'): void;
}

const emit = defineEmits<Emits>();

export type OperateType = NaiveUI.TableOperateType;

const visible = defineModel<boolean>('visible', {
  default: false
});

const { formRef, validate, restoreValidation } = useNaiveForm();
const { defaultRequiredRule } = useFormRules();

const title = computed(() => {
  const titles: Record<NaiveUI.TableOperateType, string> = {
    add: $t('page.admin.system.dict.addDict'),
    edit: $t('page.admin.system.dict.editDict')
  };
  return titles[props.operateType];
});

type Model = Pick<Api.SystemManage.Dict, 'id' | 'dictName' | 'dictCode' | 'dictDesc' | 'status'>;

const model = ref(createDefaultModel());

function createDefaultModel(): Model {
  return {
    id: 0,
    dictName: '',
    dictCode: '',
    dictDesc: '',
    status: '1'
  };
}

type RuleKey = Extract<keyof Model, 'dictName' | 'dictCode' | 'status'>;

const rules: Record<RuleKey, App.Global.FormRule> = {
  dictName: defaultRequiredRule,
  dictCode: defaultRequiredRule,
  status: defaultRequiredRule
};

function closeModal() {
  visible.value = false;
}

async function handleSubmit() {
  await validate();

  const { id, ...restData } = model.value;
  let error = null;

  if (props.operateType === 'edit') {
    ({ error } = await fetchUpdateDict(id, restData));
  } else {
    ({ error } = await fetchCreateDict(restData));
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
        <NFormItem :label="$t('page.admin.system.dict.dictName')" path="dictName">
          <NInput v-model:value="model.dictName" :placeholder="$t('page.admin.system.dict.form.dictName')" />
        </NFormItem>
        <NFormItem :label="$t('page.admin.system.dict.dictCode')" path="dictCode">
          <NInput v-model:value="model.dictCode" :placeholder="$t('page.admin.system.dict.form.dictCode')" />
        </NFormItem>
        <NFormItem :label="$t('page.admin.system.dict.dictDesc')" path="dictDesc">
          <NInput v-model:value="model.dictDesc" :placeholder="$t('page.admin.system.dict.form.dictDesc')" />
        </NFormItem>
        <NFormItem :label="$t('page.admin.system.dict.dictStatus')" path="status">
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