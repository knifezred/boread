<script setup lang="ts">
import { $t } from '@/locales';
defineOptions({ name: 'RecentlyEdited' });

interface EditItem {
  id: number;
  title: string;
  author: string;
  type: string;
  time: string;
  color: string;
}

const items: EditItem[] = [
  { id: 1, title: '剑来', author: '烽火戏诸侯', type: '更新章节', time: '2026-06-10 14:22', color: '#d4a76a' },
  { id: 2, title: '诡秘之主', author: '爱潜水的乌贼', type: '修改信息', time: '2026-06-09 22:43', color: '#b8863d' },
  { id: 3, title: '大奉打更人', author: '卖报小郎君', type: '更新章节', time: '2026-06-08 22:45', color: '#8b5e2b' },
  { id: 4, title: '凡人修仙传', author: '忘语', type: '添加标签', time: '2026-06-07 09:33', color: '#a0763a' },
  { id: 5, title: '三体', author: '刘慈欣', type: '修改信息', time: '2026-06-06 16:55', color: '#6b4f28' },
  { id: 6, title: '轮回乐园', author: '那一只蚊子', type: '添加角色', time: '2026-06-05 11:20', color: '#c49b5e' },
];
</script>

<template>
  <NCard :bordered="false" class="card-wrapper" size="small">
    <template #header>{{ $t('page.home.recentlyEdited') }}</template>
    <NList>
      <NListItem v-for="item in items" :key="item.id">
        <template #prefix>
          <div class="cover-thumb" :style="{ background: item.color }">
            <span>{{ item.title.charAt(0) }}</span>
          </div>
        </template>
        <NThing :title="item.title">
          <template #description>
            <span class="book-meta">{{ item.author }} · {{ item.time }}</span>
          </template>
          <template #header-extra>
            <NTag size="tiny" :bordered="false" round class="edit-tag">{{ item.type }}</NTag>
          </template>
        </NThing>
      </NListItem>
    </NList>
  </NCard>
</template>

<style scoped>
.cover-thumb {
  width: 36px; height: 48px; border-radius: 4px;
  display: flex; align-items: center; justify-content: center;
  color: #fff; font-size: 14px; font-weight: 700; flex-shrink: 0;
}
.book-title { font-size: 13px; font-weight: 600; color: var(--n-text-color); }
.edit-tag { margin-left: 6px; }
.book-meta { font-size: 12px; color: var(--n-text-color-disabled); }
</style>