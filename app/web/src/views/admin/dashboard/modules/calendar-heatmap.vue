<script setup lang="ts">
import { ref, computed } from 'vue'
import { $t } from '@/locales'
import { useThemeStore } from '@/store/modules/theme'

defineOptions({ name: 'CalendarHeatmap' })

const themeStore = useThemeStore()
const dark = computed(() => themeStore.darkMode)

type Range = 7 | 14 | 30
const activeRange = ref<Range>(30)

const weekLabels = ['', '一', '', '三', '', '五', '']

interface DayCell {
  date: string
  dayOfWeek: number
  value: number
  intensity: number
}

const allDays = computed<DayCell[]>(() => {
  const result: DayCell[] = []
  const now = new Date()
  for (let i = 59; i >= 0; i--) {
    const d = new Date(now)
    d.setDate(d.getDate() - i)
    const value = Math.floor(Math.random() * 180) + 5
    const intensity = value > 150 ? 4 : value > 100 ? 3 : value > 60 ? 2 : value > 20 ? 1 : 0
    result.push({
      date: `${d.getMonth() + 1}/${d.getDate()}`,
      dayOfWeek: d.getDay(),
      value,
      intensity
    })
  }
  return result
})

const days = computed(() => allDays.value.slice(-activeRange.value))

const weekColumns = computed(() => {
  const groups: DayCell[][] = []
  let current: DayCell[] = []
  const startDow = days.value[0]?.dayOfWeek ?? 0

  for (let i = 0; i < startDow; i++) {
    current.push({ date: '', dayOfWeek: i, value: 0, intensity: -1 })
  }

  for (const day of days.value) {
    current.push(day)
    if (day.dayOfWeek === 6) {
      groups.push(current)
      current = []
    }
  }
  if (current.length) groups.push(current)
  return groups
})

const ranges: Range[] = [7, 14, 30]
const rangeLabels: Record<Range, string> = { 7: '7d', 14: '14d', 30: '30d' }

function intClass(i: number): string {
  if (i === -1) return 'empty'
  const lv = ['l0', 'l1', 'l2', 'l3', 'l4']
  return lv[i] || 'l0'
}
</script>

<template>
  <NCard :bordered="false" size="small" class="!rd-10px">
    <div class="flex items-center justify-between mb-3 flex-wrap gap-2">
      <div class="flex items-baseline gap-2">
        <h3 class="text-sm font-semibold m0" :class="dark ? 'text-[#e0dcd6]' : 'text-[#3a3028]'">
          {{ $t('page.home.readingCalendar') }}
        </h3>
        <span class="text-xs" :class="dark ? 'text-[#777]' : 'text-[#aaa]'">{{ activeRange }}天</span>
      </div>
      <div class="flex items-center gap-3">
        <div class="flex gap-1">
          <button
            v-for="r in ranges"
            :key="r"
            class="px-2 py-0.5 rd-3px text-xs cursor-pointer border-none transition-all duration-150"
            :class="activeRange === r
              ? (dark ? 'bg-[#b8863d] text-[#1a1208]' : 'bg-[#d4a76a] text-white')
              : (dark ? 'bg-[#333] text-[#999] hover:bg-[#3a3a3a]' : 'bg-[#f0f0f0] text-[#666] hover:bg-[#e8e8e8]')"
            @click="activeRange = r"
          >
            {{ rangeLabels[r] }}
          </button>
        </div>
        <div class="flex items-center gap-1 text-xs" :class="dark ? 'text-[#777]' : 'text-[#aaa]'">
          <span>少</span>
          <span v-for="i in 5" :key="i" class="w-3 h-3 rd-1px" :class="intClass(i - 1)" />
          <span>多</span>
        </div>
      </div>
    </div>

    <div class="flex gap-1 overflow-x-auto pb-1">
      <div class="flex flex-col gap-1 pt-0 pr-1">
        <span v-for="w in weekLabels" :key="w" class="h-3 text-[10px] leading-3" :class="dark ? 'text-[#777]' : 'text-[#aaa]'">{{ w }}</span>
      </div>
      <div class="flex gap-1">
        <div v-for="(col, ci) in weekColumns" :key="ci" class="flex flex-col gap-1">
          <div
            v-for="(cell, ri) in col"
            :key="`${ci}-${ri}`"
            class="w-3 h-3 rd-1px transition-all duration-100"
            :class="intClass(cell.intensity)"
            :title="cell.date ? `${cell.date} · ${cell.value}分钟` : ''"
          />
        </div>
      </div>
    </div>

    <div class="flex items-center gap-2 mt-2 text-xs" :class="dark ? 'text-[#555]' : 'text-[#aaa]'">
      <span>🔥 {{ $t('page.home.currentStreak') }}: <strong class="text-[#d4a76a]">12</strong></span>
      <span class="mx-1">·</span>
      <span>🏆 {{ $t('page.home.longestStreak') }}: <strong class="text-[#d4a76a]">45</strong></span>
    </div>
  </NCard>
</template>

<style scoped>
.l0, .level-0 { background: var(--n-color, #ebedf0); }
.l1, .level-1 { background: #e8d5b7; }
.l2, .level-2 { background: #d4a76a; }
.l3, .level-3 { background: #b8863d; }
.l4, .level-4 { background: #8b5e2b; }
.empty { background: transparent; }

.dark .l0 { background: #2a2a2a; }
.dark .l1 { background: #3d3020; }
.dark .l2 { background: #6b4f28; }
.dark .l3 { background: #8b6b3a; }
.dark .l4 { background: #b88a4e; }
</style>