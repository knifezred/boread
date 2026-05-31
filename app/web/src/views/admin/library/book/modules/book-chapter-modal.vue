<script setup lang="ts">
import { h, ref, computed, watch } from "vue"
import { NButton, NDataTable, NModal, NSpace, NTag } from "naive-ui"
import { fetchGetChapterList } from "@/service/api"
import { defaultTransform, useNaivePaginatedTable } from "@/hooks/common/table"
import { $t } from "@/locales"

defineOptions({ name: "BookChapterModal" });

interface Props { bookId: number; bookTitle: string; }
const props = defineProps<Props>();
const visible = defineModel<boolean>("visible", { default: false });

const searchParams = ref({ current: 1, size: 10, bookId: props.bookId || null, fileId: null, chapterNo: null });

const { columns, data, pagination, loading, getDataByPage } = useNaivePaginatedTable({
  api: () => fetchGetChapterList(searchParams.value),
  onPaginationParamsChange: (params) => { searchParams.value.current = params.page || 1; searchParams.value.size = params.pageSize || 10; },
  transform: (response) => defaultTransform(response),
  columns: () => [
    { key: "chapterNo", title: $t("page.admin.library.book.chapterNo"), align: "center", width: 80 },
    { key: "title", title: $t("page.admin.library.book.chapterTitle"), align: "left", ellipsis: { tooltip: true }, minWidth: 200 },
    { key: "wordCount", title: $t("page.admin.library.book.wordCount"), align: "center", width: 100 },
    {
      key: "status", title: $t("page.admin.library.book.chapterStatus"), align: "center", width: 100,
      render: (row: { status: string }) => {
        const map: Record<string, NaiveUI.ThemeColor> = { "1": "success", "2": "warning", "3": "error" };
        const label: Record<string, string> = { "1": $t("common.enable"), "2": "草稿", "3": "下架" };
        return h(NTag, { type: map[row.status] || "default" }, () => label[row.status] || row.status);
      },
    },
  ],
});

const title = computed(() => `${$t("page.admin.library.book.chapters")} - ${props.bookTitle}`);

// 当 bookId 变化时，重置搜索参数并重新加载
watch(() => props.bookId, (val) => {
  searchParams.value = { current: 1, size: 10, bookId: val || null, fileId: null, chapterNo: null };
});

// 弹窗打开时刷新数据
watch(visible, (val) => {
  if (val) {
    searchParams.value.bookId = props.bookId || null;
    getDataByPage(1);
  }
});

function closeModal() { visible.value = false; }
</script>

<template>
  <NModal v-model:show="visible" :title="title" preset="card" class="w-700px" @update:show="(val) => { if (val) getDataByPage(1); }">
    <NDataTable :columns="columns" :data="data" size="small" :loading="loading" :pagination="pagination" :row-key="(row: { id: number }) => row.id" remote class="max-h-480px" />
    <template #footer>
      <NSpace justify="end">
        <NButton @click="closeModal">{{ $t("common.close") }}</NButton>
      </NSpace>
    </template>
  </NModal>
</template>
