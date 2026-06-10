import { ref, watch, type Component } from 'vue'
import { $t } from '@/locales'

import HeaderBanner from '@/views/admin/dashboard/modules/header-banner.vue'
import CardData from '@/views/admin/dashboard/modules/card-data.vue'
import LibraryTrend from '@/views/admin/dashboard/modules/library-trend.vue'
import PieChart from '@/views/admin/dashboard/modules/pie-chart.vue'
import RecentImports from '@/views/admin/dashboard/modules/recent-imports.vue'
import WordCountDist from '@/views/admin/dashboard/modules/word-count-dist.vue'
import RecentlyEdited from '@/views/admin/dashboard/modules/recently-edited.vue'
import RatingDist from '@/views/admin/dashboard/modules/rating-dist.vue'
import PendingBooks from '@/views/admin/dashboard/modules/pending-books.vue'
import RandomRecommend from '@/views/admin/dashboard/modules/random-recommend.vue'
import ReadingHours from '@/views/admin/dashboard/modules/reading-hours.vue'
import WeeklyPattern from '@/views/admin/dashboard/modules/weekly-pattern.vue'
import CalendarHeatmap from '@/views/admin/dashboard/modules/calendar-heatmap.vue'
import LineChart from '@/views/admin/dashboard/modules/line-chart.vue'
import ReadingByCategory from '@/views/admin/dashboard/modules/reading-by-category.vue'
import ReadingByTag from '@/views/admin/dashboard/modules/reading-by-tag.vue'
import RadarChart from '@/views/admin/dashboard/modules/radar-chart.vue'
import TopAuthors from '@/views/admin/dashboard/modules/top-authors.vue'
import TagDistribution from '@/views/admin/dashboard/modules/tag-distribution.vue'
import ProjectNews from '@/views/admin/dashboard/modules/project-news.vue'

export interface ModuleDef {
  id: string
  name: string
  comp: Component
  span: number
}

const DEFAULT_ORDER: string[] = [
  'header-banner', 'card-data',
  'library-trend', 'pie-chart',
  'recent-imports', 'word-count-dist',
  'recently-edited', 'rating-dist',
  'pending-books', 'random-recommend',
  'reading-hours', 'weekly-pattern', 'calendar-heatmap',
  'line-chart', 'reading-by-category', 'reading-by-tag',
  'radar-chart', 'top-authors', 'tag-distribution',
  'project-news',
]

const MODULE_MAP: Record<string, ModuleDef> = {
  'header-banner': { id: 'header-banner', name: '顶部横幅', comp: HeaderBanner, span: 24 },
  'card-data': { id: 'card-data', name: '关键指标', comp: CardData, span: 24 },
  'library-trend': { id: 'library-trend', name: '书库增长趋势', comp: LibraryTrend, span: 16 },
  'pie-chart': { id: 'pie-chart', name: '分类分布', comp: PieChart, span: 8 },
  'recent-imports': { id: 'recent-imports', name: '近期入库', comp: RecentImports, span: 14 },
  'word-count-dist': { id: 'word-count-dist', name: '字数分布', comp: WordCountDist, span: 10 },
  'recently-edited': { id: 'recently-edited', name: '最新编辑', comp: RecentlyEdited, span: 14 },
  'rating-dist': { id: 'rating-dist', name: '收藏与评分', comp: RatingDist, span: 10 },
  'pending-books': { id: 'pending-books', name: '待完善书籍', comp: PendingBooks, span: 14 },
  'random-recommend': { id: 'random-recommend', name: '随机推荐', comp: RandomRecommend, span: 10 },
  'reading-hours': { id: 'reading-hours', name: '阅读时段', comp: ReadingHours, span: 8 },
  'weekly-pattern': { id: 'weekly-pattern', name: '每周模式', comp: WeeklyPattern, span: 8 },
  'calendar-heatmap': { id: 'calendar-heatmap', name: '阅读热力', comp: CalendarHeatmap, span: 8 },
  'line-chart': { id: 'line-chart', name: '阅读趋势', comp: LineChart, span: 8 },
  'reading-by-category': { id: 'reading-by-category', name: '分类阅读时长', comp: ReadingByCategory, span: 8 },
  'reading-by-tag': { id: 'reading-by-tag', name: '标签阅读时长', comp: ReadingByTag, span: 8 },
  'radar-chart': { id: 'radar-chart', name: '阅读偏好', comp: RadarChart, span: 8 },
  'top-authors': { id: 'top-authors', name: '作者 TOP', comp: TopAuthors, span: 8 },
  'tag-distribution': { id: 'tag-distribution', name: '标签分布', comp: TagDistribution, span: 8 },
  'project-news': { id: 'project-news', name: '书库动态', comp: ProjectNews, span: 24 },
}

const STORAGE_KEY = 'dashboard-layout-v1'

export const SIZE_OPTIONS = [
  { span: 24, label: '全宽' },
  { span: 12, label: '1/2' },
  { span: 8, label: '1/3' },
  { span: 6, label: '1/4' },
]

interface LayoutStore {
  visible: string[]
  hidden: string[]
  sizes: Record<string, number>
}

function loadLayout(): LayoutStore {
  try {
    const raw = localStorage.getItem(STORAGE_KEY)
    if (raw) {
      const parsed = JSON.parse(raw) as LayoutStore
      if (Array.isArray(parsed.visible) && Array.isArray(parsed.hidden)) {
        return { ...parsed, sizes: parsed.sizes || {} }
      }
    }
  } catch { /* ignore */ }
  return { visible: [...DEFAULT_ORDER], hidden: [], sizes: {} }
}

function saveLayout(store: LayoutStore) {
  localStorage.setItem(STORAGE_KEY, JSON.stringify(store))
}

export function useDashboardLayout() {
  const isEditing = ref(false)

  const store = ref<LayoutStore>(loadLayout())

  const visibleModules = ref<ModuleDef[]>(
    store.value.visible.map(id => {
      const def = MODULE_MAP[id]
      if (!def) return null
      const savedSpan = store.value.sizes[id]
      return { ...def, span: savedSpan || def.span }
    }).filter(Boolean) as ModuleDef[]
  )

  const hiddenModules = ref<ModuleDef[]>(
    store.value.hidden.map(id => {
      const def = MODULE_MAP[id]
      if (!def) return null
      const savedSpan = store.value.sizes[id]
      return { ...def, span: savedSpan || def.span }
    }).filter(Boolean) as ModuleDef[]
  )

  function persist() {
    store.value.visible = visibleModules.value.map(m => m.id)
    store.value.hidden = hiddenModules.value.map(m => m.id)
    store.value.sizes = {}
    for (const mod of visibleModules.value) {
      const defaultSpan = MODULE_MAP[mod.id]?.span
      if (mod.span !== defaultSpan) {
        store.value.sizes[mod.id] = mod.span
      }
    }
    for (const mod of hiddenModules.value) {
      const defaultSpan = MODULE_MAP[mod.id]?.span
      if (mod.span !== defaultSpan) {
        store.value.sizes[mod.id] = mod.span
      }
    }
    saveLayout(store.value)
  }

  function removeModule(id: string) {
    const idx = visibleModules.value.findIndex(m => m.id === id)
    if (idx === -1) return
    const [removed] = visibleModules.value.splice(idx, 1)
    hiddenModules.value.push(removed)
    persist()
    buildRows()
  }

  function addModule(id: string) {
    const idx = hiddenModules.value.findIndex(m => m.id === id)
    if (idx === -1) return
    const [added] = hiddenModules.value.splice(idx, 1)
    visibleModules.value.push(added)
    persist()
    buildRows()
  }

  function setModuleSize(id: string, span: number) {
    const mod = visibleModules.value.find(m => m.id === id)
    if (mod) {
      mod.span = span
      persist()
      buildRows()
    }
  }

  function resetLayout() {
    visibleModules.value = DEFAULT_ORDER.map(id => ({ ...MODULE_MAP[id] })).filter(Boolean) as ModuleDef[]
    hiddenModules.value = []
    store.value.sizes = {}
    saveLayout(store.value)
    buildRows()
  }

  function toggleEdit() {
    isEditing.value = !isEditing.value
  }

  const rows = ref<ModuleDef[][]>([])

  function buildRows() {
    const groups: ModuleDef[][] = []
    let current: ModuleDef[] = []
    let currentSpan = 0
    for (const mod of visibleModules.value) {
      if (mod.span === 24) {
        if (current.length) {
          groups.push(current)
          current = []
          currentSpan = 0
        }
        groups.push([mod])
      } else if (currentSpan + mod.span <= 24) {
        current.push(mod)
        currentSpan += mod.span
      } else {
        groups.push(current)
        current = [mod]
        currentSpan = mod.span
      }
    }
    if (current.length) groups.push(current)
    rows.value = groups
  }

  watch(visibleModules, buildRows, { immediate: true, deep: true })

  return {
    isEditing,
    visibleModules,
    hiddenModules,
    rows,
    removeModule,
    addModule,
    setModuleSize,
    resetLayout,
    toggleEdit,
    persist,
  }
}