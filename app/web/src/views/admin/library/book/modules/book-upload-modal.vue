<script setup lang="ts">
import { ref } from "vue"
import { NButton, NForm, NFormItem, NInput, NModal, NProgress, NSpace, NText, NUpload, NUploadDragger, NAlert } from "naive-ui"
import { useFormRules, useNaiveForm } from "@/hooks/common/form"
import { fetchUploadBookFile, fetchConfirmImport } from "@/service/api"
import { $t } from "@/locales"

defineOptions({ name: "BookUploadModal" });

const visible = defineModel<boolean>("visible", { default: false });
const emit = defineEmits<{ (e: "imported"): void }>();
const { validate } = useNaiveForm();
const { defaultRequiredRule } = useFormRules();

const step = ref<"upload" | "confirm">("upload");
const uploading = ref(false);
const uploadProgress = ref(0);
const importing = ref(false);
const uploadResult = ref<Api.SystemManage.FileUploadResponse | null>(null);
const importResult = ref<Api.SystemManage.ConfirmImportResponse | null>(null);
const confirmModel = ref({ title: "", author: "" });
const importError = ref("");

async function handleUpload(file: File) {
  uploading.value = true;
  uploadProgress.value = 0;
  uploadResult.value = null;
  importResult.value = null;
  importError.value = "";
  const { error, data } = await fetchUploadBookFile(file, (p) => { uploadProgress.value = p; });
  if (error) {
    importError.value = error.message || $t("common.operateFailed");
  } else {
    uploadResult.value = data;
    confirmModel.value = { title: data.suggestedTitle, author: data.suggestedAuthor };
    step.value = "confirm";
  }
  uploading.value = false;
}

function handleFileSelect(options: { file: { file: File } }) {
  handleUpload(options.file.file);
  return true;
}

async function handleConfirm() {
  await validate();
  importing.value = true;
  importError.value = "";
  const { error, data } = await fetchConfirmImport({
    uploadId: uploadResult.value!.uploadId,
    title: confirmModel.value.title,
    author: confirmModel.value.author,
  });
  if (error) {
    importError.value = error.message || $t("common.operateFailed");
    importing.value = false;
    return;
  }
  importResult.value = data;
  importing.value = false;
  window.$message?.success($t("page.admin.library.book.importSuccess", { title: data.bookTitle, count: data.chapterCount }));
  emit("imported");
}

function closeModal() {
  visible.value = false;
  resetState();
}

function resetState() {
  step.value = "upload";
  uploading.value = false;
  uploadProgress.value = 0;
  importing.value = false;
  uploadResult.value = null;
  importResult.value = null;
  importError.value = "";
  confirmModel.value = { title: "", author: "" };
}
</script>

<template>
  <NModal v-model:show="visible" :title="$t('page.admin.library.book.uploadFile')" preset="card" class="w-520px" @update:show="(val) => { if (!val) resetState(); }">
    <!-- Step 1: Upload -->
    <template v-if="step === 'upload'">
      <NUpload :multiple="false" :show-file-list="false" :disabled="uploading" @change="handleFileSelect" accept=".txt,.epub,.mobi,.pdf">
        <NUploadDragger>
          <div class="flex flex-col items-center gap-8px py-24px">
            <SvgIcon class="text-12 color-primary" icon="solar:upload-linear"></SvgIcon>
            <NText>{{ $t("page.admin.library.book.selectFile") }}</NText>
            <NText depth="3" class="text-12px">{{ $t("page.admin.library.book.fileFormat") }}</NText>
          </div>
        </NUploadDragger>
      </NUpload>
      <div v-if="uploading" class="mt-16px">
        <NProgress :percentage="uploadProgress" :height="20" :indicator-placement="'inside'" processing />
      </div>
      <NAlert v-if="importError && !uploadResult" :type="'error'" class="mt-16px" closable>
        {{ importError }}
      </NAlert>
    </template>

    <!-- Step 2: Confirm -->
    <template v-if="step === 'confirm'">
      <NAlert type="success" class="mb-16px">{{ $t("page.admin.library.book.uploadSuccess") }}</NAlert>
      <NAlert v-if="uploadResult?.matchedBookId" type="warning" class="mb-16px" closable>
        已匹配到已有小说「{{ uploadResult.matchedBookTitle }}」，确认入库后将追加为新的文件版本
      </NAlert>
      <NForm ref="formRef" :model="confirmModel" :rules="{ title: defaultRequiredRule }" label-placement="left" :label-width="80">
        <NFormItem :label="$t('page.admin.library.book.uploadTitle')" path="title">
          <NInput v-model:value="confirmModel.title" placeholder="请输入书名" />
        </NFormItem>
        <NFormItem :label="$t('page.admin.library.book.uploadAuthor')" path="author">
          <NInput v-model:value="confirmModel.author" placeholder="请输入作者" />
        </NFormItem>
      </NForm>
      <div v-if="importing" class="mb-16px">
        <NProgress :percentage="100" :height="20" :indicator-placement="'inside'" processing />
      </div>
      <NAlert v-if="importError" :type="'error'" class="mb-16px" closable>
        {{ importError }}
      </NAlert>
      <NAlert v-if="importResult" type="success" class="mb-16px">
        {{ $t("page.admin.library.book.importSuccess", { title: importResult.bookTitle, count: importResult.chapterCount }) }}
      </NAlert>
    </template>

    <template #footer>
      <NSpace justify="end">
        <template v-if="step === 'upload'">
          <NButton @click="closeModal">{{ $t("common.close") }}</NButton>
        </template>
        <template v-if="step === 'confirm'">
          <NButton @click="closeModal">{{ $t("common.close") }}</NButton>
          <NButton v-if="!importResult" type="primary" :loading="importing" @click="handleConfirm">{{ $t("page.admin.library.book.confirmImport") }}</NButton>
        </template>
      </NSpace>
    </template>
  </NModal>
</template>
