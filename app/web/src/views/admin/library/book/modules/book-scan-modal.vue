<script setup lang="ts">
import { ref } from "vue"
import { NButton, NInput, NModal, NProgress, NSpace, NText, NAlert, NTag, NScrollbar } from "naive-ui"
import { fetchScanPath } from "@/service/api"
import { $t } from "@/locales"

defineOptions({ name: "BookScanModal" });

const visible = defineModel<boolean>("visible", { default: false });
const emit = defineEmits<{ (e: "scanned"): void }>();

const scanPath = ref("");
const scanning = ref(false);
const scanResult = ref<Api.BookManage.ScanPathResponse | null>(null);
const scanError = ref("");

async function handleScan() {
  if (!scanPath.value.trim()) {
    scanError.value = $t("page.admin.library.book.scanPathEmpty");
    return;
  }
  scanning.value = true;
  scanError.value = "";
  scanResult.value = null;
  const { error, data } = await fetchScanPath(scanPath.value.trim());
  scanning.value = false;
  if (error) {
    scanError.value = error.message || $t("common.operateFail");
  } else {
    scanResult.value = data;
    emit("scanned");
  }
}

function closeModal() {
  visible.value = false;
  scanPath.value = "";
  scanResult.value = null;
  scanError.value = "";
}
</script>

<template>
  <NModal v-model:show="visible" :title="$t('page.admin.library.book.scanLocalDir')" preset="card" class="w-560px" @update:show="(val) => { if (!val) closeModal(); }">
    <NInput v-model:value="scanPath" :placeholder="$t('page.admin.library.book.scanPathPlaceholder')" :disabled="scanning" class="mb-16px" />
    <NButton type="primary" :loading="scanning" :disabled="!scanPath.trim()" @click="handleScan" class="mb-16px">
      {{ $t("page.admin.library.book.startScan") }}
    </NButton>

    <div v-if="scanning" class="mb-16px">
      <NProgress :percentage="100" :height="20" :indicator-placement="'inside'" processing />
      <NText depth="3" class="text-12px">{{ $t("page.admin.library.book.scanningHint") }}</NText>
    </div>

    <NAlert v-if="scanError" type="error" class="mb-16px" closable>
      {{ scanError }}
    </NAlert>

    <template v-if="scanResult">
      <NAlert type="success" class="mb-16px">
        {{ $t("page.admin.library.book.scanResultText", { total: scanResult.total, imported: scanResult.imported, failed: scanResult.failed }) }}
      </NAlert>
      <NScrollbar class="max-h-240px">
        <div v-for="r in scanResult.results" :key="r.uploadId" class="flex items-center gap-8px py-4px border-b border-gray-200">
          <NTag :type="r.parseStatus === '3' ? 'success' : 'error'" size="small">
            {{ r.parseStatus === "3" ? $t("common.success") : $t("common.fail") }}
          </NTag>
          <NText class="text-13px truncate flex-1">{{ r.originalName }}</NText>
          <NText v-if="r.parseMessage" depth="3" class="text-12px">{{ r.parseMessage }}</NText>
        </div>
      </NScrollbar>
    </template>

    <template #footer>
      <NSpace justify="end">
        <NButton @click="closeModal">{{ $t("common.close") }}</NButton>
      </NSpace>
    </template>
  </NModal>
</template>
