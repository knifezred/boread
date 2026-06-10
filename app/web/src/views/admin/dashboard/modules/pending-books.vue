<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { $t } from '@/locales';
defineOptions({ name: 'PendingBooks' });

const router = useRouter();

interface PendingItem {
  id: number;
  title: string;
  author: string;
  missing: string[];
  color: string;
}

const books = ref<PendingItem[]>([
  { id: 1, title: '沧元图', author: '我吃西红柿', missing: ['缺简介', '缺标签'], color: '#8b5e2b' },
  { id: 2, title: '深空彼岸', author: '辰东', missing: ['缺分类'], color: '#a0763a' },
  { id: 3, title: '夜的命名术', author: '会说话的肘子', missing: ['缺简介'], color: '#6b4f28' },
  { id: 4, title: '长夜余火', author: '爱潜水的乌贼', missing: ['缺标签'], color: '#b8863d' },
  { id: 5, title: '星门', author: '老鹰吃小鸡', missing: ['缺简介', '缺分类'], color: '#c49b5e' },
  { id: 6, title: '不科学御兽', author: '轻泉流响', missing: ['缺标签'], color: '#d4a76a' },
]);
</script>

<template>
  <NCard :bordered="false" class="card-wrapper" size="small">
    <template #header>{{ $t('page.home.pendingBooks') }}</template>
    <NList>
      <NListItem v-for="book in books" :key="book.id">
        <template #prefix>
          <div class="cover-thumb" :style="{ background: book.color }">
            <span>{{ book.title.charAt(0) }}</span>
          </div>
        </template>
        <NThing :title="book.title">
          <template #description>
            <span class="book-meta">{{ book.author }}</span>
            <template v-for="m in book.missing" :key="m">
              <NTag size="tiny" type="warning" :bordered="false" round class="missing-tag">{{ m }}</NTag>
            </template>
          </template>
        </NThing>
        <template #suffix>
          <NButton size="tiny" tertiary @click="router.push('/admin/library/book')">
            {{ $t('page.home.goEdit') }}
          </NButton>
        </template>
      </NListItem>
    </NList>
  </NCard>
</template>

<style scoped>
.cover-thumb { width: 36px; height: 48px; border-radius: 4px; display: flex; align-items: center; justify-content: center; color: #fff; font-size: 14px; font-weight: 700; flex-shrink: 0; }
.book-title { font-size: 13px; font-weight: 600; color: var(--n-text-color); margin-right: 4px; }
.missing-tag { margin: 0 2px; }
.book-meta { font-size: 12px; color: var(--n-text-color-disabled); }
</style>