<script setup lang="ts">
import { computed, ref, watch } from "vue"
import { NButton, NForm, NFormItem, NInput, NInputNumber, NModal, NRadioGroup, NRadio, NScrollbar, NSpace, NSelect, NAlert } from "naive-ui"
import { useFormRules, useNaiveForm } from "@/hooks/common/form"
import { fetchCreateChapterRule, fetchUpdateChapterRule } from "@/service/api"
import { $t } from "@/locales"

defineOptions({ name: "ChapterRuleOperateModal" });

interface Props {
  operateType: NaiveUI.TableOperateType;
  rowData?: Api.SystemManage.BookChapterRule | null;
}
const props = defineProps<Props>();
interface Emits { (e: "submitted"): void }
const emit = defineEmits<Emits>();
const visible = defineModel<boolean>("visible", { default: false });
const { formRef, validate, restoreValidation } = useNaiveForm();
const { defaultRequiredRule } = useFormRules();
const submitting = ref(false);

const title = computed(() => props.operateType === "add" ? $t("page.admin.library.bookChapterRule.addRule") : $t("page.admin.library.bookChapterRule.editRule"));

const model = ref<Api.SystemManage.ChapterRuleRequest>({
  ruleName: "", scopeType: "1", bookId: null, pattern: "", titleGroup: 0, minChapterLen: 100, maxChapterLen: 100000, priority: 0, description: null, status: "1",
});

const rules: Record<string, App.Global.FormRule> = {
  ruleName: defaultRequiredRule,
  pattern: defaultRequiredRule,
};

function handleInitModel() {
  if (props.operateType === "edit" && props.rowData) {
    model.value = {
      ruleName: props.rowData.ruleName,
      scopeType: props.rowData.scopeType,
      bookId: props.rowData.bookId,
      pattern: props.rowData.pattern,
      titleGroup: props.rowData.titleGroup,
      minChapterLen: props.rowData.minChapterLen,
      maxChapterLen: props.rowData.maxChapterLen,
      priority: props.rowData.priority,
      description: props.rowData.description,
      status: props.rowData.status,
    };
  } else {
    model.value = { ruleName: "", scopeType: "1", bookId: null, pattern: "", titleGroup: 0, minChapterLen: 100, maxChapterLen: 100000, priority: 0, description: null, status: "1" };
  }
}

async function handleSubmit() {
  await validate();
  submitting.value = true;
  try {
    if (props.operateType === "edit" && props.rowData) {
      await fetchUpdateChapterRule(props.rowData.id, model.value);
    } else {
      await fetchCreateChapterRule(model.value);
    }
    window.$message?.success($t("common.updateSuccess"));
    visible.value = false;
    emit("submitted");
  } catch (err: any) {
    window.$message?.error(err.message || $t("common.operateFail"));
  } finally {
    submitting.value = false;
  }
}

function closeModal() { visible.value = false; }

watch(visible, (val) => { if (val) { handleInitModel(); restoreValidation(); } });
</script>

<template>
  <NModal v-model:show="visible" :title="title" preset="card" class="w-600px" :loading="submitting">
    <NScrollbar class="h-420px pr-20px">
      <NForm ref="formRef" :model="model" :rules="rules" label-placement="left" :label-width="120">
        <NFormItem :label="$t('page.admin.library.bookChapterRule.ruleName')" path="ruleName">
          <NInput v-model:value="model.ruleName" :placeholder="$t('page.admin.library.bookChapterRule.form.ruleName')" />
        </NFormItem>
        <NFormItem :label="$t('page.admin.library.bookChapterRule.scopeType')" path="scopeType">
          <NRadioGroup v-model:value="model.scopeType">
            <NRadio value="1">{{ $t("page.admin.library.bookChapterRule.scopeGlobal") }}</NRadio>
            <NRadio value="2">{{ $t("page.admin.library.bookChapterRule.scopeBook") }}</NRadio>
          </NRadioGroup>
        </NFormItem>
        <NFormItem :label="$t('page.admin.library.bookChapterRule.pattern')" path="pattern">
          <NInput v-model:value="model.pattern" type="textarea" :rows="3" :placeholder="$t('page.admin.library.bookChapterRule.form.pattern')" />
        </NFormItem>
        <NAlert type="info" closable class="mb-12px">{{ $t("page.admin.library.bookChapterRule.rulePreview") }}</NAlert>
        <NFormItem :label="$t('page.admin.library.bookChapterRule.titleGroup')" path="titleGroup">
          <NInputNumber v-model:value="model.titleGroup" :min="0" class="w-full" />
        </NFormItem>
        <NFormItem :label="$t('page.admin.library.bookChapterRule.minChapterLen')" path="minChapterLen">
          <NInputNumber v-model:value="model.minChapterLen" :min="0" class="w-full" />
        </NFormItem>
        <NFormItem :label="$t('page.admin.library.bookChapterRule.maxChapterLen')" path="maxChapterLen">
          <NInputNumber v-model:value="model.maxChapterLen" :min="0" class="w-full" />
        </NFormItem>
        <NFormItem :label="$t('page.admin.library.bookChapterRule.priority')" path="priority">
          <NInputNumber v-model:value="model.priority" class="w-full" />
        </NFormItem>
        <NFormItem :label="$t('page.admin.library.bookChapterRule.description')" path="description">
          <NInput v-model:value="model.description" :placeholder="$t('page.admin.library.bookChapterRule.form.description')" />
        </NFormItem>
        <NFormItem :label="$t('page.admin.library.bookChapterRule.status')" path="status">
          <NRadioGroup v-model:value="model.status">
            <NRadio value="1">{{ $t("common.enable") }}</NRadio>
            <NRadio value="2">{{ $t("common.disable") }}</NRadio>
          </NRadioGroup>
        </NFormItem>
      </NForm>
    </NScrollbar>
    <template #footer>
      <NSpace justify="end">
        <NButton @click="closeModal">{{ $t("common.cancel") }}</NButton>
        <NButton type="primary" :loading="submitting" @click="handleSubmit">{{ $t("common.confirm") }}</NButton>
      </NSpace>
    </template>
  </NModal>
</template>
