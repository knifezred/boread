<script setup lang="tsx">
import { computed, ref, watch } from 'vue';
import { NButton, NForm, NFormItem, NInput, NModal, NScrollbar, NSpace } from 'naive-ui';
import { useFormRules, useNaiveForm } from '@/hooks/common/form';
import { $t } from '@/locales';
import { fetchCreateTag, fetchUpdateTag } from '@/service/api';

defineOptions({ name: 'TagOperateModal' });

interface Props { operateType: NaiveUI.TableOperateType; rowData?: Api.SystemManage.BookTag | null; }
const props = defineProps<Props>();
interface Emits { (e: 'submitted'): void; }
const emit = defineEmits<Emits>();
const visible = defineModel<boolean>('visible', { default: false });
const { formRef, validate, restoreValidation } = useNaiveForm();
const { defaultRequiredRule } = useFormRules();

const title = computed(() => {
  const titles: Record<NaiveUI.TableOperateType, string> = {
    add: $t('page.manage.bookTag.addTag'), edit: $t('page.manage.bookTag.editTag')
  };
  return titles[props.operateType];
});

type Model = { id: number; tagName: string; description: string };
const model = ref<Model>({ id: 0, tagName: '', description: '' });
const rules: Record<'tagName', App.Global.FormRule> = { tagName: defaultRequiredRule };

function closeModal() { visible.value = false; }

async function handleSubmit() {
  await validate();
  const { id, ...restData } = model.value;
  let error = null;
  if (props.operateType === 'edit') ({ error } = await fetchUpdateTag(id, restData));
  else ({ error } = await fetchCreateTag(restData));
  if (error) window.$message?.error(error.message);
  else window.$message?.success($t('common.updateSuccess'));
  closeModal();
  emit('submitted');
}

watch(visible, () => {
  if (visible.value) {
    model.value = { id: props.rowData?.id ?? 0, tagName: props.rowData?.tagName ?? '', description: props.rowData?.description ?? '' };
    restoreValidation();
  }
});
</script>

<template>
  <NModal v-model:show="visible" :title="title" preset="card" class="w-500px">
    <NScrollbar class="h-200px pr-20px">
      <NForm ref="formRef" :model="model" :rules="rules" label-placement="left" :label-width="80">
        <NFormItem :label="$t('page.manage.bookTag.tagName')" path="tagName">
          <NInput v-model:value="model.tagName" :placeholder="$t('page.manage.bookTag.form.tagName')" />
        </NFormItem>
        <NFormItem :label="$t('page.manage.bookTag.description')" path="description">
          <NInput v-model:value="model.description" type="textarea" :placeholder="$t('page.manage.bookTag.form.description')" />
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
