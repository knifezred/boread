<script setup lang="tsx">
import { computed, ref, watch } from 'vue';
import { NButton, NForm, NFormItem, NInput, NInputNumber, NModal, NScrollbar, NRadioGroup, NRadio, NSpace } from 'naive-ui';
import { enableStatusOptions } from '@/constants/business';
import { useFormRules, useNaiveForm } from '@/hooks/common/form';
import { $t } from '@/locales';
import { fetchCreateDept, fetchUpdateDept } from '@/service/api';

defineOptions({
  name: 'DeptOperateModal'
});

export type OperateType = NaiveUI.TableOperateType | 'addChild';

interface Props {
  /** the type of operation */
  operateType: OperateType;
  /** the edit dept data or the parent dept data when adding a child dept */
  rowData?: Api.SystemManage.Dept | null;
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
  const titles: Record<OperateType, string> = {
    add: $t('page.admin.system.dept.addDept'),
    addChild: $t('page.admin.system.dept.addChildDept'),
    edit: $t('page.admin.system.dept.editDept')
  };
  return titles[props.operateType];
});

type Model = {
  id: number;
  parentId: number;
  deptName: string;
  deptCode: string;
  leader: string;
  sortOrder: number;
  status: Api.Common.EnableStatus;
};

const model = ref<Model>(createDefaultModel());

function createDefaultModel(): Model {
  return {
    id: 0,
    parentId: 0,
    deptName: '',
    deptCode: '',
    leader: '',
    sortOrder: 0,
    status: '1'
  };
}

type RuleKey = Extract<keyof Model, 'deptName' | 'deptCode' | 'status'>;

const rules: Record<RuleKey, App.Global.FormRule> = {
  deptName: defaultRequiredRule,
  deptCode: defaultRequiredRule,
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
    ({ error } = await fetchUpdateDept(id, restData));
  } else {
    ({ error } = await fetchCreateDept(restData));
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

  if (!props.rowData) return;

  if (props.operateType === 'addChild') {
    const { id } = props.rowData;
    Object.assign(model.value, { parentId: id });
  }

  if (props.operateType === 'edit') {
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
        <NFormItem :label="$t('page.admin.system.dept.parentId')" path="parentId">
          <NInputNumber v-model:value="model.parentId" :min="0" disabled class="w-full" />
        </NFormItem>
        <NFormItem :label="$t('page.admin.system.dept.deptName')" path="deptName">
          <NInput v-model:value="model.deptName" :placeholder="$t('page.admin.system.dept.form.deptName')" />
        </NFormItem>
        <NFormItem :label="$t('page.admin.system.dept.deptCode')" path="deptCode">
          <NInput v-model:value="model.deptCode" :placeholder="$t('page.admin.system.dept.form.deptCode')" />
        </NFormItem>
        <NFormItem :label="$t('page.admin.system.dept.leader')" path="leader">
          <NInput v-model:value="model.leader" :placeholder="$t('page.admin.system.dept.form.leader')" />
        </NFormItem>
        <NFormItem :label="$t('page.admin.system.dept.sortOrder')" path="sortOrder">
          <NInputNumber v-model:value="model.sortOrder" :min="0" class="w-full" />
        </NFormItem>
        <NFormItem :label="$t('page.admin.system.dept.status')" path="status">
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
