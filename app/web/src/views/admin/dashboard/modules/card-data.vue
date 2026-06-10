<script setup lang="ts">
import { computed } from 'vue'
import { $t } from '@/locales'
import { useThemeStore } from '@/store/modules/theme'

defineOptions({ name: 'CardData' })

const themeStore = useThemeStore()
const dark = computed(() => themeStore.darkMode)

interface CardItem {
  key: string
  label: string
  value: number
  unit: string
  icon: string
  accent: string
}

const cards = computed<CardItem[]>(() => [
  { key: 'books', label: $t('page.home.totalBooks'), value: 12586, unit: '', icon: '📚', accent: '#d4a76a' },
  { key: 'words', label: $t('page.home.totalWordCount'), value: 25680, unit: 'w', icon: '📝', accent: '#b8863d' },
  { key: 'authors', label: $t('page.home.authorCount'), value: 156, unit: '', icon: '✍️', accent: '#8b5e2b' },
  { key: 'tags', label: $t('page.home.tagCount'), value: 86, unit: '', icon: '🏷️', accent: '#a0763a' },
  { key: 'rate', label: $t('page.home.completionRate'), value: 68, unit: '%', icon: '🎯', accent: '#6b4f28' },
  { key: 'weekly', label: $t('page.home.weeklyNew'), value: 47, unit: '', icon: '✨', accent: '#b8863d' },
  { key: 'reading', label: $t('page.home.totalReadingHours'), value: 2580, unit: $t('page.home.readingUnit'), icon: '⏱', accent: '#d4a76a' },
  { key: 'readWords', label: $t('page.home.totalReadingWordCount'), value: 12580, unit: 'w', icon: '📖', accent: '#8b5e2b' },
])

function fmt(n: number): string {
  if (n >= 10000) return (n / 10000).toFixed(1) + 'w'
  if (n >= 1000) return (n / 1000).toFixed(1) + 'k'
  return String(n)
}
</script>

<template>
  <NCard :bordered="false" size="small" class="!rd-10px">
    <div class="grid grid-cols-2 sm:grid-cols-4 gap-3">
      <div
        v-for="c in cards"
        :key="c.key"
        class="flex items-center gap-3 rd-8px p-3 transition-all duration-200 hover:-translate-y-1 hover:shadow-md"
        :class="[dark ? 'bg-[#292929] hover:bg-[#323232] border border-[#383838]' : 'bg-white hover:bg-[#faf8f5] border border-[#efefef]']"
        :style="{ '--accent': c.accent }"
      >
        <div
          class="flex items-center justify-center w-10 h-10 rd-8px shrink-0 text-xl"
          :style="{ background: `color-mix(in srgb, ${c.accent} 14%, transparent)` }"
        >
          {{ c.icon }}
        </div>
        <div class="flex flex-col gap-0.5 min-w-0">
          <span
            class="text-base font-bold leading-tight"
            :class="dark ? 'text-[#e0dcd6]' : 'text-[#3a3028]'"
          >
            {{ fmt(c.value) }}<small v-if="c.unit" class="text-xs font-normal ml-0.5" :class="dark ? 'text-[#777]' : 'text-[#aaa]'">{{ c.unit }}</small>
          </span>
          <span class="text-xs truncate" :class="dark ? 'text-[#888]' : 'text-[#999]'">
            {{ c.label }}
          </span>
        </div>
      </div>
    </div>
  </NCard>
</template>