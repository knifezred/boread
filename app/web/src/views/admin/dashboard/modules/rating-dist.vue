<script setup lang="ts">
import { $t } from '@/locales';
defineOptions({ name: 'RatingDist' });

interface FavItem {
  title: string;
  count: number;
}
const favData: FavItem[] = [
  { title: '剑来', count: 1256 },
  { title: '诡秘之主', count: 982 },
  { title: '大奉打更人', count: 756 },
  { title: '凡人修仙传', count: 620 },
  { title: '三体', count: 512 },
];

function getMax(): number {
  return Math.max(...favData.map(d => d.count));
}

function barWidth(count: number): string {
  return Math.max((count / getMax()) * 100, 15) + '%';
}

interface RatingItem {
  label: string;
  value: number;
  color: string;
}
const ratingData: RatingItem[] = [
  { label: '5星', value: 45, color: '#d4a76a' },
  { label: '4星', value: 30, color: '#b8863d' },
  { label: '3星', value: 15, color: '#8b5e2b' },
  { label: '2星', value: 7, color: '#a0763a' },
  { label: '1星', value: 3, color: '#6b4f28' },
];

function fmt(n: number): string {
  if (n >= 10000) return (n / 10000).toFixed(1) + 'w';
  if (n >= 1000) return (n / 1000).toFixed(1) + 'k';
  return String(n);
}
</script>

<template>
  <NCard :bordered="false" class="card-wrapper" size="small">
    <template #header>{{ $t('page.home.ratingDist') }}</template>
    <div class="rating-grid">
      <div class="chart-section">
        <div class="chart-title">收藏 TOP</div>
        <div class="bar-list">
          <div v-for="item in favData" :key="item.title" class="bar-row">
            <span class="bar-label">{{ item.title }}</span>
            <div class="bar-track">
              <div class="bar-fill" :style="{ width: barWidth(item.count) }"></div>
            </div>
            <span class="bar-value">{{ fmt(item.count) }}</span>
          </div>
        </div>
      </div>
      <NDivider vertical class="divider" />
      <div class="chart-section">
        <div class="chart-title">{{ $t('page.home.rating') }}</div>
        <div class="stars-list">
          <div v-for="r in ratingData" :key="r.label" class="star-row">
            <span class="star-label">{{ r.label }}</span>
            <div class="star-track">
              <div class="star-fill" :style="{ width: r.value + '%', background: r.color }"></div>
            </div>
            <span class="star-value">{{ r.value }}%</span>
          </div>
        </div>
      </div>
    </div>
  </NCard>
</template>

<style scoped>
.rating-grid {
  display: grid;
  grid-template-columns: 1fr auto 1fr;
  gap: 0;
  align-items: stretch;
}
.chart-section { padding: 4px 0; }
.chart-title { font-size: 12px; color: var(--n-text-color-disabled); text-align: center; margin-bottom: 8px; }
.divider { height: auto; }

.bar-list, .stars-list { display: flex; flex-direction: column; gap: 6px; }
.bar-row, .star-row { display: flex; align-items: center; gap: 6px; font-size: 12px; }
.bar-label, .star-label { width: 56px; text-align: right; color: var(--n-text-color); white-space: nowrap; overflow: hidden; text-overflow: ellipsis; flex-shrink: 0; }
.bar-track, .star-track { flex: 1; height: 14px; background: var(--n-color, #f0f0f0); border-radius: 7px; overflow: hidden; }
.bar-fill, .star-fill { height: 100%; border-radius: 7px; transition: width 0.3s; }
.bar-fill { background: linear-gradient(90deg, #d4a76a, #8b5e2b); }
.bar-value, .star-value { width: 40px; color: var(--n-text-color-disabled); font-size: 11px; }
</style>