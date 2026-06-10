<script setup lang="ts">
import { watch } from 'vue';
import { useAppStore } from '@/store/modules/app';
import { useEcharts } from '@/hooks/common/echarts';
import { $t } from '@/locales';
defineOptions({ name: 'ReadingByCategory' });

const appStore = useAppStore();

const { domRef, updateOptions } = useEcharts(() => ({
  tooltip: { trigger: 'axis', axisPointer: { type: 'shadow' } },
  grid: { left: '1%', right: '6%', bottom: '3%', top: '3%', containLabel: true },
  xAxis: { type: 'value', axisLabel: { show: false }, splitLine: { show: false } },
  yAxis: { type: 'category', data: ['都市', '历史', '悬疑', '科幻', '仙侠', '言情', '玄幻奇幻'], axisLine: { show: false }, axisTick: { show: false } },
  series: [{ type: 'bar', barWidth: 12, label: { show: true, position: 'right', formatter: '{c}h', fontSize: 10 }, itemStyle: { borderRadius: [0, 3, 3, 0] }, data: [{ value: 30, itemStyle: { color: '#6b4f28' } }, { value: 40, itemStyle: { color: '#8b5e2b' } }, { value: 55, itemStyle: { color: '#a0763a' } }, { value: 60, itemStyle: { color: '#b8863d' } }, { value: 85, itemStyle: { color: '#c49b5e' } }, { value: 70, itemStyle: { color: '#d4a76a' } }, { value: 128, itemStyle: { color: '#e9bb7e' } }] }]
}));

async function init() { await new Promise(r => setTimeout(r, 800)); }
watch(() => appStore.locale, () => {});
init();
</script>

<template>
  <NCard :bordered="false" class="card-wrapper" size="small">
    <template #header>{{ $t('page.home.readingByCategory') }}</template>
    <div ref="domRef" class="h-300px overflow-hidden"></div>
  </NCard>
</template>

<style scoped></style>