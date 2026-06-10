<script setup lang="ts">
import { watch } from 'vue';
import { useAppStore } from '@/store/modules/app';
import { useEcharts } from '@/hooks/common/echarts';
import { $t } from '@/locales';
defineOptions({ name: 'ReadingByTag' });

const appStore = useAppStore();

const { domRef, updateOptions } = useEcharts(() => ({
  tooltip: { trigger: 'axis', axisPointer: { type: 'shadow' } },
  grid: { left: '1%', right: '6%', bottom: '3%', top: '3%', containLabel: true },
  xAxis: { type: 'value', axisLabel: { show: false }, splitLine: { show: false } },
  yAxis: { type: 'category', data: ['黑暗流', '种田', '重生', '无敌流', '脑洞', '后宫', '系统流', '穿越'], axisLine: { show: false }, axisTick: { show: false } },
  series: [{ type: 'bar', barWidth: 12, label: { show: true, position: 'right', formatter: '{c}h', fontSize: 10 }, itemStyle: { borderRadius: [0, 3, 3, 0] }, data: [{ value: 22, itemStyle: { color: '#6b4f28' } }, { value: 28, itemStyle: { color: '#8b5e2b' } }, { value: 35, itemStyle: { color: '#a0763a' } }, { value: 40, itemStyle: { color: '#b8863d' } }, { value: 48, itemStyle: { color: '#c49b5e' } }, { value: 55, itemStyle: { color: '#d4a76a' } }, { value: 65, itemStyle: { color: '#e9bb7e' } }, { value: 80, itemStyle: { color: '#f0c88a' } }] }]
}));

async function init() { await new Promise(r => setTimeout(r, 800)); }
watch(() => appStore.locale, () => {});
init();
</script>

<template>
  <NCard :bordered="false" class="card-wrapper" size="small">
    <template #header>{{ $t('page.home.readingByTag') }}</template>
    <div ref="domRef" class="h-300px overflow-hidden"></div>
  </NCard>
</template>

<style scoped></style>