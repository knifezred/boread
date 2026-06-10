<script setup lang="ts">
import { watch } from 'vue';
import { useAppStore } from '@/store/modules/app';
import { useEcharts } from '@/hooks/common/echarts';
import { $t } from '@/locales';
defineOptions({ name: 'WeeklyPattern' });

const appStore = useAppStore();

const { domRef, updateOptions } = useEcharts(() => ({
  tooltip: {
    trigger: 'axis',
    axisPointer: { type: 'shadow' }
  },
  grid: {
    left: '3%',
    right: '3%',
    bottom: '3%',
    top: '12%',
    containLabel: true
  },
  xAxis: {
    type: 'category',
    data: ['周一', '周二', '周三', '周四', '周五', '周六', '周日'],
    axisLabel: { fontSize: 11 }
  },
  yAxis: {
    type: 'value',
    name: 'min',
    splitLine: { lineStyle: { color: 'rgba(180,150,110,0.12)' } }
  },
  series: [
    {
      type: 'bar',
      barWidth: '55%',
      data: [
        { value: 45, itemStyle: { color: '#b8863d' } },
        { value: 38, itemStyle: { color: '#b8863d' } },
        { value: 52, itemStyle: { color: '#b8863d' } },
        { value: 42, itemStyle: { color: '#b8863d' } },
        { value: 35, itemStyle: { color: '#b8863d' } },
        { value: 78, itemStyle: { color: '#d4a76a' } },
        { value: 95, itemStyle: { color: '#d4a76a' } }
      ],
      label: {
        show: true,
        position: 'top',
        formatter: '{c}min',
        fontSize: 10,
        color: 'var(--n-text-color-disabled)'
      },
      itemStyle: {
        borderRadius: [4, 4, 0, 0]
      }
    }
  ]
}));

async function init() {
  await new Promise(r => setTimeout(r, 800));
}

watch(() => appStore.locale, () => {});
init();
</script>

<template>
  <NCard :bordered="false" class="card-wrapper" size="small">
    <template #header>{{ $t('page.home.weeklyPattern') }}</template>
    <div ref="domRef" class="h-300px overflow-hidden"></div>
  </NCard>
</template>

<style scoped></style>