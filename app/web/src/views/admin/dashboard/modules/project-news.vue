<script setup lang="ts">
import { computed } from 'vue'
import { $t } from '@/locales'

defineOptions({ name: 'ProjectNews' })

interface ActivityItem {
  id: number
  type: string
  content: string
  time: string
}

const activities = computed<ActivityItem[]>(() => [
  { id: 1, type: $t('page.home.activity.newBook'), content: $t('page.home.recentActivity.desc1'), time: '2026-06-10 14:22' },
  { id: 2, type: $t('page.home.activity.newCategory'), content: $t('page.home.recentActivity.desc2'), time: '2026-06-09 10:24' },
  { id: 3, type: $t('page.home.activity.newTag'), content: $t('page.home.recentActivity.desc3'), time: '2026-06-08 22:43' },
  { id: 4, type: $t('page.home.activity.newCharacter'), content: $t('page.home.recentActivity.desc4'), time: '2026-06-07 20:33' },
  { id: 5, type: $t('page.home.activity.editBook'), content: $t('page.home.recentActivity.desc5'), time: '2026-06-06 22:45' },
  { id: 6, type: $t('page.home.activity.newCategory'), content: $t('page.home.recentActivity.desc6'), time: '2026-06-05 15:12' },
  { id: 7, type: $t('page.home.activity.newTag'), content: $t('page.home.recentActivity.desc7'), time: '2026-06-04 09:33' },
  { id: 8, type: $t('page.home.delete'), content: '删除小说《仙王的日常生活》', time: '2026-06-03 11:45' },
])

function dotClass(type: string): string {
  const map: Record<string, string> = {
    [$t('page.home.activity.newBook')]: 'dot-book',
    [$t('page.home.activity.editBook')]: 'dot-edit',
    [$t('page.home.activity.newCategory')]: 'dot-cat',
    [$t('page.home.activity.newTag')]: 'dot-tag',
    [$t('page.home.activity.newCharacter')]: 'dot-char',
  }
  return map[type] || 'dot-edit'
}
</script>

<template>
  <NCard :bordered="false" size="small" class="!rd-10px">
    <template #header>
      <div class="flex items-center justify-between w-full">
        <span>{{ $t('page.home.recentActivity.title') }}</span>
        <a href="javascript:;" class="text-[#b8863d] text-xs no-underline hover:underline">{{ $t('page.home.recentActivity.moreNews') }}</a>
      </div>
    </template>
    <NList>
      <NListItem v-for="item in activities" :key="item.id">
        <template #prefix>
          <span class="block w-2 h-2 rd-1/2" :class="dotClass(item.type)" />
        </template>
        <NThing :title="item.content" :description="item.time" />
      </NListItem>
    </NList>
  </NCard>
</template>

<style scoped>
.dot-book { background: #d4a76a; }
.dot-edit { background: #a0763a; }
.dot-cat { background: #b8863d; }
.dot-tag { background: #8b5e2b; }
.dot-char { background: #e9bb7e; }
</style>