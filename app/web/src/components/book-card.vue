<script setup lang="ts">
import { ref, computed } from 'vue'

const props = defineProps<{
  book: Api.SystemManage.Book
  showStatusTag?: boolean
  showTitleAuthor?: boolean
}>()

const emit = defineEmits<{
  click: [book: Api.SystemManage.Book]
}>()

const coverLoadFailed = ref(false)

// 判断是否显示真实封面

// 根据书名生成渐变背景
const coverGradient = computed(() => {
  const gradients = [
    'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
    'linear-gradient(135deg, #f093fb 0%, #f5576c 100%)',
    'linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)',
    'linear-gradient(135deg, #43e97b 0%, #38f9d7 100%)',
    'linear-gradient(135deg, #fa709a 0%, #fee140 100%)',
    'linear-gradient(135deg, #30cfd0 0%, #330867 100%)',
    'linear-gradient(135deg, #a8edea 0%, #fed6e3 100%)',
    'linear-gradient(135deg, #ffecd2 0%, #fcb69f 100%)',
  ]
  // 根据书名hash选渐变，保证同一本书渐变固定
  let hash = 0
  for (let i = 0; i < props.book.title.length; i++) {
    hash = props.book.title.charCodeAt(i) + ((hash << 5) - hash)
  }
  return gradients[Math.abs(hash) % gradients.length]
})

// 状态标签文本
const statusText = computed(() => {
  return props.book.serialStatus === '1' ? '完本' : '连载'
})
const statusType = computed(() => {
  return props.book.serialStatus === '1' ? 'success' : 'warning'
})

function handleClick() {
  emit('click', props.book)
}

function handleCoverError() {
  coverLoadFailed.value = true
}
</script>

<template>
  <div class="book-card" @click="handleClick">
    <!-- 封面容器 -->
    <div class="cover-wrapper">
      <!-- 真实封面 -->
      <img
        v-if="book.cover"
        :src="book.cover"
        :alt="book.title"
        class="cover-img"
        @error="handleCoverError"
        loading="lazy"
      >

      <!-- 占位封面 -->
      <div v-else class="cover-placeholder" :style="{ background: coverGradient }">
        <div class="placeholder-content">
          <div class="book-title">{{ book.title }}</div>
          <div class="book-author">{{ book.author }}</div>
        </div>
      </div>

      <!-- 状态标签 -->
      <div v-if="showStatusTag && book.status" class="status-tag" :class="statusType">
        {{ statusText }}
      </div>
    </div>

    <!-- 下方书籍信息（可选，默认显示） -->
    <div v-if="showTitleAuthor" class="book-info">
      <h3 class="title" :title="book.title">{{ book.title }}</h3>
      <p class="author" :title="book.author">{{ book.author }}</p>
    </div>
  </div>
</template>

<style scoped>
.book-card {
  cursor: pointer;
  transition: all 0.25s ease;
  width: 100%;
}

.book-card:hover {
  transform: translateY(-4px);
}

.cover-wrapper {
  position: relative;
  width: 100%;
  aspect-ratio: 2/3;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
  margin-bottom: 10px;
}

.cover-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  background-color: #f5f5f5;
  transition: transform 0.3s ease;
}

.book-card:hover .cover-img {
  transform: scale(1.05);
}

/* 占位封面样式 */
.cover-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 16px;
  box-sizing: border-box;
}

.placeholder-content {
  text-align: center;
  color: #fff;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
  width: 100%;
}

.book-title {
  font-size: clamp(16px, 3vw, 22px);
  font-weight: 700;
  line-height: 1.3;
  margin-bottom: 12px;
  /* 最多显示3行 */
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
  word-break: break-all;
}

.book-author {
  font-size: clamp(12px, 2vw, 14px);
  font-weight: 500;
  opacity: 0.9;
  /* 最多显示1行 */
  display: -webkit-box;
  -webkit-line-clamp: 1;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

/* 状态标签 */
.status-tag {
  position: absolute;
  top: 8px;
  right: 8px;
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
  color: #fff;
}

.status-tag.success {
  background-color: rgba(16, 185, 129, 0.9);
}

.status-tag.warning {
  background-color: rgba(245, 158, 11, 0.9);
}

/* 下方书籍信息 */
.book-info {
  padding: 0 2px;
}

.title {
  font-size: 14px;
  font-weight: 500;
  margin: 0 0 6px 0;
  line-height: 1.4;
  color: inherit;
  /* 最多1行 */
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.author {
  font-size: 12px;
  color: #888;
  margin: 0;
  /* 最多1行 */
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

/* 暗色模式适配 */
:deep(.dark) .author {
  color: #aaa;
}
</style>
