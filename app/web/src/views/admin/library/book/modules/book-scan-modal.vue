<script setup lang="ts">
import { ref } from "vue"
import { NButton, NInput, NModal, NProgress, NSpace, NText, NAlert, NTag } from "naive-ui"
import { fetchScanPath } from "@/service/api"
import { $t } from "@/locales"

defineOptions({ name: "BookScanModal" });

const visible = defineModel<boolean>("visible", { default: false });
const emit = defineEmits<{ (e: "scanned"): void }>();

const scanPath = ref("");
const scanning = ref(false);
const scanResult = ref<Api.SystemManage.ScanPathResponse | null>(null);
const scanError = ref("");

async function handleScan() {
  if (!scanPath.value.trim()) {
    scanError.value = "请输入扫描路径";
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
  <NModal v-model:show="visible" title="扫描本地目录" preset="card" class="w-560px" @update:show="(val) => { if (!val) closeModal(); }">
    <NInput v-model:value="scanPath" placeholder="请输入本地目录路径，如 D:\books" :disabled="scanning" class="mb-16px" />
    <NButton type="primary" :loading="scanning" :disabled="!scanPath.trim()" @click="handleScan" class="mb-16px">
      开始扫描
    </NButton>

    <div v-if="scanning" class="mb-16px">
      <NProgress :percentage="100" :height="20" :indicator-placement="'inside'" processing />
      <NText depth="3" class="text-12px">正在扫描并入库，请稍候...</NText>
    </div>

    <NAlert v-if="scanError" type="error" class="mb-16px" closable>
      {{ scanError }}
    </NAlert>

    <template v-if="scanResult">
      <NAlert type="success" class="mb-16px">
        扫描完成：共发现 {{ scanResult.total }} 个文件，成功入库 {{ scanResult.imported }} 个，失败 {{ scanResult.failed }} 个
      </NAlert>
      <NScrollbar class="max-h-240px">
        <div v-for="r in scanResult.results" :key="r.uploadId" class="flex items-center gap-8px py-4px border-b border-gray-200">
          <NTag :type="r.parseStatus === '3' ? 'success' : 'error'" size="small">
            {{ r.parseStatus === "3" ? "成功" : "失败" }}
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
