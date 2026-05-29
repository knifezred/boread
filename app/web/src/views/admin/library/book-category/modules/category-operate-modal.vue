<script setup lang="tsx">
import { computed, ref, watch } from "vue";
import {
  NButton,
  NForm,
  NFormItem,
  NInput,
  NInputNumber,
  NModal,
  NScrollbar,
  NRadioGroup,
  NRadio,
  NSpace,
} from "naive-ui";
import { enableStatusOptions } from "@/constants/business";
import { useFormRules, useNaiveForm } from "@/hooks/common/form";
import { $t } from "@/locales";
import { fetchCreateCategory, fetchUpdateCategory } from "@/service/api";

defineOptions({ name: "CategoryOperateModal" });

export type OperateType = NaiveUI.TableOperateType | "addChild";

interface Props {
  operateType: OperateType;
  rowData?: Api.SystemManage.BookCategory | null;
}

const props = defineProps<Props>();
interface Emits {
  (e: "submitted"): void;
}
const emit = defineEmits<Emits>();
const visible = defineModel<boolean>("visible", { default: false });
const { formRef, validate, restoreValidation } = useNaiveForm();
const { defaultRequiredRule } = useFormRules();

const title = computed(() => {
  const titles: Record<OperateType, string> = {
    add: $t("page.admin.library.bookCategory.addCategory"),
    addChild: $t("page.admin.library.bookCategory.addChildCategory"),
    edit: $t("page.admin.library.bookCategory.editCategory"),
  };
  return titles[props.operateType];
});

type Model = {
  id: number;
  parentId: number;
  categoryName: string;
  categoryCode: string;
  description: string;
  sortOrder: number;
  status: Api.Common.EnableStatus;
};
const model = ref<Model>(createDefaultModel());
function createDefaultModel(): Model {
  return {
    id: 0,
    parentId: 0,
    categoryName: "",
    categoryCode: "",
    description: "",
    sortOrder: 0,
    status: "1",
  };
}
type RuleKey = Extract<keyof Model, "categoryName" | "categoryCode" | "status">;
const rules: Record<RuleKey, App.Global.FormRule> = {
  categoryName: defaultRequiredRule,
  categoryCode: defaultRequiredRule,
  status: defaultRequiredRule,
};

function closeModal() {
  visible.value = false;
}

async function handleSubmit() {
  await validate();
  const { id, ...restData } = model.value;
  let error = null;
  if (props.operateType === "edit")
    ({ error } = await fetchUpdateCategory(id, restData));
  else ({ error } = await fetchCreateCategory(restData));
  if (error) window.$message?.error(error.message);
  else window.$message?.success($t("common.updateSuccess"));
  closeModal();
  emit("submitted");
}

function handleInitModel() {
  model.value = createDefaultModel();
  if (!props.rowData) return;
  if (props.operateType === "addChild")
    Object.assign(model.value, { parentId: props.rowData.id });
  if (props.operateType === "edit") Object.assign(model.value, props.rowData);
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
    <NScrollbar class="h-400px pr-20px">
      <NForm
        ref="formRef"
        :model="model"
        :rules="rules"
        label-placement="left"
        :label-width="100"
      >
        <NFormItem
          :label="$t('page.admin.library.bookCategory.parentId')"
          path="parentId"
        >
          <NInputNumber
            v-model:value="model.parentId"
            :min="0"
            disabled
            class="w-full"
          />
        </NFormItem>
        <NFormItem
          :label="$t('page.admin.library.bookCategory.categoryName')"
          path="categoryName"
        >
          <NInput
            v-model:value="model.categoryName"
            :placeholder="$t('page.admin.library.bookCategory.form.categoryName')"
          />
        </NFormItem>
        <NFormItem
          :label="$t('page.admin.library.bookCategory.categoryCode')"
          path="categoryCode"
        >
          <NInput
            v-model:value="model.categoryCode"
            :placeholder="$t('page.admin.library.bookCategory.form.categoryCode')"
          />
        </NFormItem>
        <NFormItem
          :label="$t('page.admin.library.bookCategory.description')"
          path="description"
        >
          <NInput
            v-model:value="model.description"
            type="textarea"
            :placeholder="$t('page.admin.library.bookCategory.form.description')"
          />
        </NFormItem>
        <NFormItem
          :label="$t('page.admin.library.bookCategory.sortOrder')"
          path="sortOrder"
        >
          <NInputNumber
            v-model:value="model.sortOrder"
            :min="0"
            class="w-full"
          />
        </NFormItem>
        <NFormItem
          :label="$t('page.admin.library.bookCategory.categoryStatus')"
          path="status"
        >
          <NRadioGroup v-model:value="model.status">
            <NRadio
              v-for="item in enableStatusOptions"
              :key="item.value"
              :value="item.value"
              :label="$t(item.label)"
            />
          </NRadioGroup>
        </NFormItem>
      </NForm>
    </NScrollbar>
    <template #footer>
      <NSpace justify="end" :size="16">
        <NButton @click="closeModal">{{ $t("common.cancel") }}</NButton>
        <NButton type="primary" @click="handleSubmit">
          {{ $t("common.confirm") }}
        </NButton>
      </NSpace>
    </template>
  </NModal>
</template>
