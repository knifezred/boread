<script setup lang="ts">
import { watch } from 'vue';
import { useAppStore } from '@/store/modules/app';
import { useEcharts } from '@/hooks/common/echarts';
import { $t } from '@/locales';
defineOptions({ name: 'ReadingHours' });

const appStore = useAppStore();

const { domRef, updateOptions } = useEcharts(() => ({
  tooltip: {
    trigger: 'axis',
    axisPointer: { type: 'shadow' }
  },
  grid: {
    left: '2%',
    right: '3%',
    bottom: '3%',
    top: '5%',
    containLabel: true
  },
  xAxis: {
    type: 'category',
    data: [
      '00-02', '02-04', '04-06', '06-08',
      '08-10', '10-12', '12-14', '14-16',
      '16-18', '18-20', '20-22', '22-24'
    ],
    axisLabel: { fontSize: 10, rotate: 30 }
  },
  yAxis: {
    type: 'value',
    name: 'min'
  },
  series: [
    {
      type: 'bar',
      barWidth: '70%',
      data: [5, 2, 1, 8, 15, 25, 30, 20, 22, 35, 55, 40],
      itemStyle: {
        borderRadius: [3, 3, 0, 0],
        color: {
          type: 'linear',
          x: 0, y: 0, x2: 0, y2: 1,
          colorStops: [
            { offset: 0, color: '#d4a76a' },
            { offset: 1, color: '#8b5e2b' }
          ]
        }
      },
      emphasis: {
        itemStyle: {
          color: '#b8863d'
        }
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
    <template #header>{{ $t('page.home.readingHours') }}</template>
    <div ref="domRef" class="h-280px overflow-hidden"></div>
  </NCard>
</template>

<style scoped></style>