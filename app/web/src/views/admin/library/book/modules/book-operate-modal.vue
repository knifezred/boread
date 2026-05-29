<script setup lang="tsx">
import { computed, ref, watch } from "vue";
import {
  NButton,
  NForm,
  NFormItem,
  NInput,
  NModal,
  NScrollbar,
  NRadioGroup,
  NRadio,
  NSpace,
  NSelect,
} from "naive-ui";
import { useBoolean } from "@sa/hooks";
import { useFormRules, useNaiveForm } from "@/hooks/common/form";
import { useDictItems } from "@/hooks/business/dict";
import { $t } from "@/locales";
import {
  fetchCreateBook,
  fetchUpdateBook,
  fetchGetCategoryTree,
  fetchGetTagList,
} from "@/service/api";

defineOptions({ name: "BookOperateModal" });

interface Props {
  operateType: NaiveUI.TableOperateType;
  rowData?: Api.SystemManage.Book | null;
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
  const titles: Record<NaiveUI.TableOperateType, string> = {
    add: $t("page.admin.library.book.addBook"),
    edit: $t("page.admin.library.book.editBook"),
  };
  return titles[props.operateType];
});

type Model = {
  id: number;
  title: string;
  author: string;
  cover: string | null;
  intro: string | null;
  categoryId: number | null;
  language: string;
  serialStatus: Api.SystemManage.SerialStatus;
  visibility: Api.SystemManage.Visibility;
  tagIds: number[];
};

const model = ref<Model>({
  id: 0,
  title: "",
  author: "",
  cover: null,
  intro: null,
  categoryId: null,
  language: "zh-CN",
  serialStatus: "1",
  visibility: "1",
  tagIds: [],
});

const rules: Record<
  "title" | "author" | "serialStatus" | "visibility",
  App.Global.FormRule
> = {
  title: defaultRequiredRule,
  author: defaultRequiredRule,
  serialStatus: defaultRequiredRule,
  visibility: defaultRequiredRule,
};

const categoryOptions = ref<CommonType.Option<number>[]>([]);
const tagOptions = ref<CommonType.Option<number>[]>([]);
const {
  bool: loadingCategories,
  setTrue: startLoadingCategories,
  setFalse: doneLoadingCategories,
} = useBoolean();
const {
  bool: loadingTags,
  setTrue: startLoadingTags,
  setFalse: doneLoadingTags,
} = useBoolean();
const { labelMap: serialStatusLabelMap } = useDictItems("book_serial_status");
const { labelMap: visibilityLabelMap } = useDictItems("book_visibility");
async function loadCategoryOptions() {
  startLoadingCategories();
  const { data, error } = await fetchGetCategoryTree();
  if (!error && data) {
    categoryOptions.value = flattenTree(data, 0);
  }
  doneLoadingCategories();
}

function flattenTree(
  nodes: Api.SystemManage.BookCategory[],
  depth: number,
): CommonType.Option<number>[] {
  let result: CommonType.Option<number>[] = [];
  for (const n of nodes) {
    const prefix = "\u3000".repeat(depth);
    result.push({ value: n.id, label: `${prefix}${n.categoryName}` });
    if (n.children?.length) {
      result = result.concat(flattenTree(n.children, depth + 1));
    }
  }
  return result;
}

async function loadTagOptions() {
  startLoadingTags();
  const { data, error } = await fetchGetTagList({
    current: 1,
    size: 200,
    tagName: null,
  });
  if (!error && data) {
    tagOptions.value = data.records.map((t) => ({
      value: t.id,
      label: t.tagName,
    }));
  }
  doneLoadingTags();
}

function closeModal() {
  visible.value = false;
}

async function handleSubmit() {
  await validate();
  const { id, ...restData } = model.value;
  const payload = {
    ...restData,
    cover: restData.cover || null,
    intro: restData.intro || null,
  };
  let error = null;
  if (props.operateType === "edit")
    ({ error } = await fetchUpdateBook(id, payload));
  else ({ error } = await fetchCreateBook(payload));
  if (error) window.$message?.error(error.message);
  else window.$message?.success($t("common.updateSuccess"));
  closeModal();
  emit("submitted");
}

watch(visible, () => {
  if (visible.value) {
    model.value = {
      id: props.rowData?.id ?? 0,
      title: props.rowData?.title ?? "",
      author: props.rowData?.author ?? "",
      cover: props.rowData?.cover ?? null,
      intro: props.rowData?.intro ?? null,
      categoryId: props.rowData?.categoryId ?? null,
      language: props.rowData?.language ?? "zh-CN",
      serialStatus: props.rowData?.serialStatus ?? "1",
      visibility: props.rowData?.visibility ?? "1",
      tagIds: props.rowData?.tagIds ?? [],
    };
    restoreValidation();
    loadCategoryOptions();
    loadTagOptions();
  }
});
</script>

<template>
  <NModal v-model:show="visible" :title="title" preset="card" class="w-700px">
    <NScrollbar class="h-500px pr-20px">
      <NForm
        ref="formRef"
        :model="model"
        :rules="rules"
        label-placement="left"
        :label-width="100"
      >
        <NFormItem :label="$t('page.admin.library.book.bookName')" path="title">
          <NInput
            v-model:value="model.title"
            :placeholder="$t('page.admin.library.book.form.title')"
          />
        </NFormItem>
        <NFormItem :label="$t('page.admin.library.book.author')" path="author">
          <NInput
            v-model:value="model.author"
            :placeholder="$t('page.admin.library.book.form.author')"
          />
        </NFormItem>
        <NFormItem :label="$t('page.admin.library.book.cover')" path="cover">
          <NInput
            v-model:value="model.cover"
            :placeholder="$t('page.admin.library.book.form.cover')"
          />
        </NFormItem>
        <NFormItem :label="$t('page.admin.library.book.intro')" path="intro">
          <NInput
            v-model:value="model.intro"
            type="textarea"
            :placeholder="$t('page.admin.library.book.form.intro')"
          />
        </NFormItem>
        <NFormItem
          :label="$t('page.admin.library.book.categoryId')"
          path="categoryId"
        >
          <NSelect
            v-model:value="model.categoryId"
            :options="categoryOptions"
            :loading="loadingCategories"
            :placeholder="$t('page.admin.library.book.form.categoryId')"
            clearable
            filterable
            class="w-full"
          />
        </NFormItem>
        <NFormItem
          :label="$t('page.admin.library.book.language')"
          path="language"
        >
          <NInput
            v-model:value="model.language"
            :placeholder="$t('page.admin.library.book.form.language')"
          />
        </NFormItem>
        <NFormItem
          :label="$t('page.admin.library.book.serialStatus')"
          path="serialStatus"
        >
          <NRadioGroup v-model:value="model.serialStatus">
            <NRadio
              v-for="(label, value) in serialStatusLabelMap"
              :key="value"
              :value="value"
              :label="label"
            />
          </NRadioGroup>
        </NFormItem>
        <NFormItem
          :label="$t('page.admin.library.book.visibility')"
          path="visibility"
        >
          <NRadioGroup v-model:value="model.visibility">
            <NRadio
              v-for="(label, value) in visibilityLabelMap"
              :key="value"
              :value="value"
              :label="label"
            />
          </NRadioGroup>
        </NFormItem>
        <NFormItem :label="$t('page.admin.library.book.tags')" path="tagIds">
          <NSelect
            v-model:value="model.tagIds"
            :options="tagOptions"
            :loading="loadingTags"
            :placeholder="$t('page.admin.library.book.form.tags')"
            multiple
            filterable
            class="w-full"
          />
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
