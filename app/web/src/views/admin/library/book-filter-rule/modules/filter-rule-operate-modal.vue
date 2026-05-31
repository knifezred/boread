<script setup lang="ts">
import { computed, ref, watch } from "vue"
import { NButton, NForm, NFormItem, NInput, NModal, NRadioGroup, NRadio, NScrollbar, NSpace, NSelect } from "naive-ui"
import { useFormRules, useNaiveForm } from "@/hooks/common/form"
import { fetchCreateFilterRule, fetchUpdateFilterRule } from "@/service/api"
import { $t } from "@/locales"

defineOptions({ name: "FilterRuleOperateModal" });

interface Props {
  operateType: NaiveUI.TableOperateType;
  rowData?: Api.BookManage.BookContentFilterRule | null;
}
const props = defineProps<Props>();
interface Emits { (e: "submitted"): void }
const emit = defineEmits<Emits>();
const visible = defineModel<boolean>("visible", { default: false });
const { validate, restoreValidation } = useNaiveForm();
const { defaultRequiredRule } = useFormRules();
const submitting = ref(false);

const title = computed(() => props.operateType === "add" ? $t("page.admin.library.bookFilterRule.addRule") : $t("page.admin.library.bookFilterRule.editRule"));

const model = ref<Api.BookManage.FilterRuleRequest>({
  ruleName: "", matchType: "1", pattern: "", action: "1", replacement: "***", applyStage: "1", category: null, severity: "1", description: null, status: "1",
});

const rules: Record<string, App.Global.FormRule> = {
  ruleName: defaultRequiredRule,
  pattern: defaultRequiredRule,
};

const matchTypeOptions = [
  { value: "1", label: $t("page.admin.library.bookFilterRule.matchKeyword") },
  { value: "2", label: $t("page.admin.library.bookFilterRule.matchRegex") },
];
const actionOptions = [
  { value: "1", label: $t("page.admin.library.bookFilterRule.actionReplace") },
  { value: "2", label: $t("page.admin.library.bookFilterRule.actionBlock") },
  { value: "3", label: $t("page.admin.library.bookFilterRule.actionMark") },
];
const stageOptions = [
  { value: "1", label: $t("page.admin.library.bookFilterRule.stageInput") },
  { value: "2", label: $t("page.admin.library.bookFilterRule.stageOutput") },
];
const severityOptions = [
  { value: "1", label: $t("page.admin.library.bookFilterRule.severityLow") },
  { value: "2", label: $t("page.admin.library.bookFilterRule.severityMedium") },
  { value: "3", label: $t("page.admin.library.bookFilterRule.severityHigh") },
];

function handleInitModel() {
  if (props.operateType === "edit" && props.rowData) {
    model.value = {
      ruleName: props.rowData.ruleName,
      matchType: props.rowData.matchType,
      pattern: props.rowData.pattern,
      action: props.rowData.action,
      replacement: props.rowData.replacement,
      applyStage: props.rowData.applyStage,
      category: props.rowData.category,
      severity: props.rowData.severity,
      description: props.rowData.description,
      status: props.rowData.status,
    };
  } else {
    model.value = { ruleName: "", matchType: "1", pattern: "", action: "1", replacement: "***", applyStage: "1", category: null, severity: "1", description: null, status: "1" };
  }
}

async function handleSubmit() {
  await validate();
  submitting.value = true;
  try {
    if (props.operateType === "edit" && props.rowData) {
      await fetchUpdateFilterRule(props.rowData.id, model.value);
    } else {
      await fetchCreateFilterRule(model.value);
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
    <NScrollbar class="h-480px pr-20px">
      <NForm ref="formRef" :model="model" :rules="rules" label-placement="left" :label-width="120">
        <NFormItem :label="$t('page.admin.library.bookFilterRule.ruleName')" path="ruleName">
          <NInput v-model:value="model.ruleName" :placeholder="$t('page.admin.library.bookFilterRule.form.ruleName')" />
        </NFormItem>
        <NFormItem :label="$t('page.admin.library.bookFilterRule.matchType')" path="matchType">
          <NSelect v-model:value="model.matchType" :options="matchTypeOptions" />
        </NFormItem>
        <NFormItem :label="$t('page.admin.library.bookFilterRule.pattern')" path="pattern">
          <NInput v-model:value="model.pattern" type="textarea" :rows="3" :placeholder="$t('page.admin.library.bookFilterRule.form.pattern')" />
        </NFormItem>
        <NFormItem :label="$t('page.admin.library.bookFilterRule.action')" path="action">
          <NSelect v-model:value="model.action" :options="actionOptions" />
        </NFormItem>
        <NFormItem v-if="model.action === '1'" :label="$t('page.admin.library.bookFilterRule.replacement')" path="replacement">
          <NInput v-model:value="model.replacement" :placeholder="$t('page.admin.library.bookFilterRule.form.replacement')" />
        </NFormItem>
        <NFormItem :label="$t('page.admin.library.bookFilterRule.applyStage')" path="applyStage">
          <NSelect v-model:value="model.applyStage" :options="stageOptions" />
        </NFormItem>
        <NFormItem :label="$t('page.admin.library.bookFilterRule.category')" path="category">
          <NInput v-model:value="model.category" placeholder="politics/porn/violence/ad..." />
        </NFormItem>
        <NFormItem :label="$t('page.admin.library.bookFilterRule.severity')" path="severity">
          <NSelect v-model:value="model.severity" :options="severityOptions" />
        </NFormItem>
        <NFormItem :label="$t('page.admin.library.bookFilterRule.description')" path="description">
          <NInput v-model:value="model.description" :placeholder="$t('page.admin.library.bookFilterRule.form.description')" />
        </NFormItem>
        <NFormItem :label="$t('page.admin.library.bookFilterRule.status')" path="status">
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
