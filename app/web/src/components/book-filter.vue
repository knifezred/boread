<script setup lang="ts">
import { reactive } from "vue"
import { $t } from "@/locales"

export interface BookFilterProps {
  config: {
    /** 分类列表 */
    categories: { label: string; value: string }[]
    /** 连载状态列表 */
    serialStatus: { label: string; value: string }[]
    /** 字数区间列表 */
    wordCount: { label: string; value: string }[]
    /** 标签列表 */
    tags: { label: string; value: string }[]
    /** 更新时间范围列表 */
    updateTime: { label: string; value: string }[]
  }
  modelValue: Api.BookManage.BookFilterParams
}

const props = defineProps<BookFilterProps>()
const emit = defineEmits<{
  "update:modelValue": [values: Api.BookManage.BookFilterParams]
  change: []
}>()

/** 内部响应式筛选值 */
const values = reactive<Api.BookManage.BookFilterParams>({
  categoryId: props.modelValue.categoryId,
  serialStatus: props.modelValue.serialStatus,
  wordCount: props.modelValue.wordCount,
  tagId: props.modelValue.tagId,
  updateTime: props.modelValue.updateTime,
  title: "",
  sortBy: "",
  sortOrder: ""
})

function updateFilter<K extends keyof Api.BookManage.BookFilterParams>(key: K, value: Api.BookManage.BookFilterParams[K]) {
  values[key] = value
  emit("update:modelValue", { ...values })
  emit("change")
}
</script>

<template>
  <div class="flex flex-col gap-6">
    <!-- 分类筛选 -->
    <div class="filter-group">
      <div class="text-14px font-600 mb-3 text-gray-800 dark:text-gray-200">{{ $t("page.book.filter.category") }}</div>
      <div class="flex flex-wrap gap-2">
        <span
          v-for="item in config.categories"
          :key="item.value"
          class="px-2 py-1 text-13px rd-1 cursor-pointer transition-all duration-200 select-none"
          :class="values.categoryId === Number(item.value)
              ? 'bg-primary text-white'
              : 'bg-gray-100 text-gray-500 hover:bg-primary-100 hover:text-primary dark:bg-gray-700 dark:text-gray-400 dark:hover:bg-primary-900 dark:hover:text-primary-200'
            "
          @click="updateFilter('categoryId', Number(item.value))">
          {{ item.label }}
        </span>
      </div>
    </div>

    <!-- 连载状态筛选 -->
    <div class="filter-group">
      <div class="text-14px font-600 mb-3 text-gray-800 dark:text-gray-200">{{ $t("page.book.filter.serialStatus") }}</div>
      <div class="flex flex-wrap gap-2">
        <span
          v-for="item in config.serialStatus"
          :key="item.value"
          class="px-2 py-1 text-13px rd-1 cursor-pointer transition-all duration-200 select-none"
          :class="values.serialStatus === item.value
              ? 'bg-primary text-white'
              : 'bg-gray-100 text-gray-500 hover:bg-primary-100 hover:text-primary dark:bg-gray-700 dark:text-gray-400 dark:hover:bg-primary-900 dark:hover:text-primary-200'
            "
          @click="updateFilter('serialStatus', item.value)">
          {{ item.label }}
        </span>
      </div>
    </div>

    <!-- 字数区间筛选 -->
    <div class="filter-group">
      <div class="text-14px font-600 mb-3 text-gray-800 dark:text-gray-200">{{ $t("page.book.filter.wordCount") }}</div>
      <div class="flex flex-wrap gap-2">
        <span
          v-for="item in config.wordCount"
          :key="item.value"
          class="px-2 py-1 text-13px rd-1 cursor-pointer transition-all duration-200 select-none"
          :class="values.wordCount === item.value
              ? 'bg-primary text-white'
              : 'bg-gray-100 text-gray-500 hover:bg-primary-100 hover:text-primary dark:bg-gray-700 dark:text-gray-400 dark:hover:bg-primary-900 dark:hover:text-primary-200'
            "
          @click="updateFilter('wordCount', item.value)">
          {{ item.label }}
        </span>
      </div>
    </div>

    <!-- 标签筛选 -->
    <div class="filter-group">
      <div class="text-14px font-600 mb-3 text-gray-800 dark:text-gray-200">{{ $t("page.book.filter.tags") }}</div>
      <div class="flex flex-wrap gap-2">
        <span
          v-for="item in config.tags"
          :key="item.value"
          class="px-2 py-1 text-13px rd-1 cursor-pointer transition-all duration-200 select-none"
          :class="values.tagId === item.value
              ? 'bg-primary text-white'
              : 'bg-gray-100 text-gray-500 hover:bg-primary-100 hover:text-primary dark:bg-gray-700 dark:text-gray-400 dark:hover:bg-primary-900 dark:hover:text-primary-200'
            "
          @click="updateFilter('tagId', item.value)">
          {{ item.label }}
        </span>
      </div>
    </div>

    <!-- 更新时间筛选 -->
    <div class="filter-group">
      <div class="text-14px font-600 mb-3 text-gray-800 dark:text-gray-200">{{ $t("page.book.filter.updateTime") }}</div>
      <div class="flex flex-wrap gap-2">
        <span
          v-for="item in config.updateTime"
          :key="item.value"
          class="px-2 py-1 text-13px rd-1 cursor-pointer transition-all duration-200 select-none"
          :class="values.updateTime === item.value
              ? 'bg-primary text-white'
              : 'bg-gray-100 text-gray-500 hover:bg-primary-100 hover:text-primary dark:bg-gray-700 dark:text-gray-400 dark:hover:bg-primary-900 dark:hover:text-primary-200'
            "
          @click="updateFilter('updateTime', item.value)">
          {{ item.label }}
        </span>
      </div>
    </div>
  </div>
</template>
