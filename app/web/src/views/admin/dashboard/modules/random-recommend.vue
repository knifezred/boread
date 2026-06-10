<script setup lang="ts">
import { ref } from 'vue';
import { $t } from '@/locales';
defineOptions({ name: 'RandomRecommend' });

interface RecItem {
  id: number;
  title: string;
  author: string;
  category: string;
  reason: string;
  color: string;
}

const pool: RecItem[] = [
  { id: 1, title: '剑来', author: '烽火戏诸侯', category: '玄幻奇幻', reason: '同类型阅读较多', color: '#d4a76a' },
  { id: 2, title: '诡秘之主', author: '爱潜水的乌贼', category: '玄幻奇幻', reason: '高分推荐', color: '#b8863d' },
  { id: 3, title: '凡人修仙传', author: '忘语', category: '仙侠', reason: '经典必读', color: '#8b5e2b' },
  { id: 4, title: '大奉打更人', author: '卖报小郎君', category: '玄幻奇幻', reason: '热门连载', color: '#a0763a' },
  { id: 5, title: '三体', author: '刘慈欣', category: '科幻', reason: '久未阅读', color: '#6b4f28' },
  { id: 6, title: '轮回乐园', author: '那一只蚊子', category: '无限流', reason: '同标签偏好', color: '#c49b5e' },
  { id: 7, title: '深空彼岸', author: '辰东', category: '科幻', reason: '高分推荐', color: '#5d8a3c' },
  { id: 8, title: '夜的命名术', author: '会说话的肘子', category: '都市', reason: '热门连载', color: '#3c7a8a' },
];

const list = ref<RecItem[]>([]);

function shuffle() {
  const copy = [...pool];
  const result: RecItem[] = [];
  for (let i = 0; i < 6; i++) {
    const idx = Math.floor(Math.random() * copy.length);
    result.push(copy.splice(idx, 1)[0]);
  }
  list.value = result;
}

shuffle();
</script>

<template>
  <NCard :bordered="false" class="card-wrapper" size="small">
    <template #header>
      <div class="rec-header">
        <span>{{ $t('page.home.randomRecommend') }}</span>
        <NButton size="tiny" tertiary @click="shuffle">{{ $t('page.home.refresh') }}</NButton>
      </div>
    </template>
    <NList>
      <NListItem v-for="item in list" :key="item.id">
        <template #prefix>
          <div class="cover-thumb" :style="{ background: item.color }">
            <span>{{ item.title.charAt(0) }}</span>
          </div>
        </template>
        <NThing :title="item.title">
          <template #description>
            <span class="book-meta">{{ item.author }} · {{ item.category }}</span>
            <NTag size="tiny" type="info" :bordered="false" round class="reason-tag">{{ item.reason }}</NTag>
          </template>
        </NThing>
      </NListItem>
    </NList>
  </NCard>
</template>

<style scoped>
.rec-header { display: flex; align-items: center; justify-content: space-between; width: 100%; }
.cover-thumb { width: 36px; height: 48px; border-radius: 4px; display: flex; align-items: center; justify-content: center; color: #fff; font-size: 14px; font-weight: 700; flex-shrink: 0; }
.book-title { font-size: 13px; font-weight: 600; color: var(--n-text-color); }
.book-meta { font-size: 12px; color: var(--n-text-color-disabled); margin-right: 4px; }
.reason-tag { }
</style>