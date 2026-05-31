<script setup lang="ts">
import { ref, watch } from 'vue'
import { NModal, NScrollbar } from 'naive-ui'
import { fetchChapterList } from '@/service/api'

defineOptions({ name: 'CatalogModal' })

interface Props {
  bookId: string
  chapterNo: number
  darkMode: boolean
}

const props = defineProps<Props>()

interface Emits {
  (e: 'select', chapterNo: number): void
}

const emit = defineEmits<Emits>()

const visible = defineModel<boolean>('visible', { default: false })

const chapters = ref<Api.BookManage.BookChapter[]>([])

async function loadAllChapters() {
  const bookIdNum = Number(props.bookId)
  if (!bookIdNum) return
  const { data } = await fetchChapterList(bookIdNum)
  if (data) {
    chapters.value = data
  }
}

function handleSelect(no: number) {
  emit('select', no)
  visible.value = false
}

watch(visible, (val) => {
  if (val) {
    loadAllChapters()
  }
})
</script>

<template>
  <NModal
    v-model:show="visible"
    preset="card"
    :title="$t('page.book.catalog.title')"
    class="w-800px h-100vh"
    segmented>
    <NScrollbar class="h-90vh">
      <div class="flex flex-col gap-0.5">
        <div
          v-for="ch in chapters"
          :key="ch.id"
          class="flex items-center gap-3 px-4 py-3 rd-1 cursor-pointer transition-colors duration-200 text-sm"
          :class="chapterNo === ch.chapterNo
            ? 'text-primary font-medium bg-primary/5'
            : (darkMode ? 'text-gray-300 hover:bg-gray-700' : 'text-gray-700 hover:bg-gray-50')"
          @click="handleSelect(ch.chapterNo)">
          <span class="text-xs shrink-0 w-8 text-right"
            :class="chapterNo === ch.chapterNo ? 'text-primary' : (darkMode ? 'text-gray-500' : 'text-gray-400')">{{
              ch.chapterNo }}</span>
          <span class="flex-1 overflow-hidden text-ellipsis whitespace-nowrap">{{ ch.title }}</span>
        </div>
      </div>
    </NScrollbar>
  </NModal>
</template>
