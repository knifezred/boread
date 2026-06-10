<script setup lang="ts">
import { $t } from '@/locales';
defineOptions({ name: 'RecentImports' });

interface BookItem {
  id: number;
  title: string;
  author: string;
  category: string;
  time: string;
  color: string;
}

const books: BookItem[] = [
  { id: 1, title: '剑来', author: '烽火戏诸侯', category: '玄幻奇幻', time: '2026-06-10', color: '#d4a76a' },
  { id: 2, title: '三体', author: '刘慈欣', category: '科幻', time: '2026-06-09', color: '#b8863d' },
  { id: 3, title: '大奉打更人', author: '卖报小郎君', category: '玄幻奇幻', time: '2026-06-08', color: '#8b5e2b' },
  { id: 4, title: '诡秘之主', author: '爱潜水的乌贼', category: '玄幻奇幻', time: '2026-06-07', color: '#a0763a' },
  { id: 5, title: '凡人修仙传', author: '忘语', category: '仙侠', time: '2026-06-06', color: '#6b4f28' },
  { id: 6, title: '轮回乐园', author: '那一只蚊子', category: '无限流', time: '2026-06-05', color: '#c49b5e' },
];
</script>

<template>
  <NCard :bordered="false" class="card-wrapper" size="small">
    <template #header>{{ $t('page.home.recentImports') }}</template>
    <NList>
      <NListItem v-for="book in books" :key="book.id">
        <template #prefix>
          <div class="cover-thumb" :style="{ background: book.color }">
            <span>{{ book.title.charAt(0) }}</span>
          </div>
        </template>
        <NThing :title="book.title">
          <template #description>
            <span class="book-meta">{{ book.author }} · {{ book.category }}</span>
          </template>
        </NThing>
        <template #suffix>
          <span class="import-time">{{ book.time }}</span>
        </template>
      </NListItem>
    </NList>
  </NCard>
</template>

<style scoped>
.cover-thumb {
  width: 36px;
  height: 48px;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-size: 14px;
  font-weight: 700;
  flex-shrink: 0;
}
.book-title {
  font-size: 13px;
  font-weight: 600;
  color: var(--n-text-color);
}
.book-meta {
  font-size: 12px;
  color: var(--n-text-color-disabled);
}
.import-time {
  font-size: 11px;
  color: var(--n-text-color-disabled);
  white-space: nowrap;
}
</style>