<script setup lang="ts">
import { ref, computed } from "vue"
import { useDictItems } from "@/hooks/business/dict"


/**
 * BookCard 属性接口
 */
export interface BookCardProps {
}
interface Props {
  /** 书籍数据 */
  book: Api.BookManage.Book
  /** 是否显示右上角状态标签 */
  showStatusTag?: boolean
  /** 是否显示封面下方的书名和作者 */
  showTitleAuthor?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  showStatusTag: false,
  showTitleAuthor: false
})

const emit = defineEmits<{
  click: [book: Api.BookManage.Book]
}>()

const coverLoadFailed = ref(false)

/** 根据书名 hash 生成固定渐变背景 */
const coverGradient = computed(() => {
  const gradients = [
    "linear-gradient(135deg, #667eea 0%, #764ba2 100%)",
    "linear-gradient(135deg, #f093fb 0%, #f5576c 100%)",
    "linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)",
    "linear-gradient(135deg, #43e97b 0%, #38f9d7 100%)",
    "linear-gradient(135deg, #fa709a 0%, #fee140 100%)",
    "linear-gradient(135deg, #30cfd0 0%, #330867 100%)",
    "linear-gradient(135deg, #a8edea 0%, #fed6e3 100%)",
    "linear-gradient(135deg, #ffecd2 0%, #fcb69f 100%)"
  ]
  let hash = 0
  for (let i = 0; i < props.book.title.length; i++) {
    hash = props.book.title.charCodeAt(i) + ((hash << 5) - hash)
  }
  return gradients[Math.abs(hash) % gradients.length]
})

/** 状态标签文本 */
const { labelMap: serialStatusLabelMap } = useDictItems("book_serial_status")
const statusText = computed(() => {
  return serialStatusLabelMap.value[props.book.serialStatus as string] || "连载"
})


/** 状态标签类型：连载→success，完本→warning，停更→error */
const statusType = computed(() => {
  const map: Record<string, NaiveUI.ThemeColor> = { "1": "warning", "2": "success", "3": "error" }
  return map[props.book.serialStatus] ?? "warning"
})

function handleClick() {
  emit("click", props.book)
}

function handleCoverError() {
  coverLoadFailed.value = true
}
</script>

<template>
  <div
    class="book-card cursor-pointer transition-all duration-250 hover:-translate-y-1 w-full"
    @click="handleClick">
    <!-- 封面容器 -->
    <div class="relative w-full aspect-2/3 rd-2 overflow-hidden shadow-sm mb-2">
      <!-- 真实封面 -->
      <img
        v-if="book.cover"
        :src="book.cover"
        :alt="book.title"
        class="w-full h-full object-cover bg-gray-100 transition-transform duration-300 group-hover:scale-105"
        loading="lazy"
        @error="handleCoverError">

      <!-- 占位封面（无封面时显示渐变背景 + 书名作者） -->
      <div
        v-else
        class="w-full h-full flex-center p-4 box-border"
        :style="{ background: coverGradient }">
        <div class="text-center text-white text-shadow-sm w-full">
          <div class="font-700 leading-normal mb-3 line-clamp-3 break-words">
            {{ book.title }}
          </div>
          <div class="font-500 opacity-90 truncate">
            {{ book.author }}
          </div>
        </div>
      </div>

      <!-- 右上角状态标签 -->
      <div
        v-if="showStatusTag"
        class="absolute top-2 right-2 px-1.5 py-0.5 rd-1 text-12px font-500 text-white"
        :class="statusType === 'success' ? 'bg-[rgba(16,185,129,0.9)]' : 'bg-[rgba(245,158,11,0.9)]'">
        {{ statusText }}
      </div>
    </div>

    <!-- 封面下方书籍信息 -->
    <div v-if="showTitleAuthor" class="px-0.5">
      <h3 class="text-14px font-500 m-0 leading-1.4 truncate">
        {{ book.title }}
      </h3>
      <p class="text-12px text-gray-400 m-0 truncate dark:text-gray-500">
        {{ book.author }}
      </p>
    </div>
  </div>
</template>
