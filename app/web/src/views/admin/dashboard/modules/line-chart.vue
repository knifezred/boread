<script setup lang="ts">
import { watch } from 'vue';
import { useAppStore } from '@/store/modules/app';
import { useEcharts } from '@/hooks/common/echarts';
import { $t } from '@/locales';
defineOptions({ name: 'LineChart' });

const appStore = useAppStore();

const { domRef, updateOptions } = useEcharts(() => ({
  tooltip: {
    trigger: 'axis',
    axisPointer: { type: 'cross', label: { backgroundColor: '#6a7985' } }
  },
  legend: { data: ['阅读时长', '阅读次数'], top: 0 },
  grid: { left: '3%', right: '4%', bottom: '3%', top: '18%' },
  xAxis: {
    type: 'category',
    boundaryGap: false,
    data: [] as string[]
  },
  yAxis: { type: 'value' },
  series: [
    {
      name: '阅读时长',
      type: 'line',
      smooth: true,
      color: '#b8863d',
      areaStyle: {
        color: {
          type: 'linear', x: 0, y: 0, x2: 0, y2: 1,
          colorStops: [
            { offset: 0.2, color: 'rgba(184, 134, 61, 0.45)' },
            { offset: 1, color: 'rgba(184, 134, 61, 0.02)' }
          ]
        }
      },
      emphasis: { focus: 'series' },
      data: [] as number[]
    },
    {
      name: '阅读次数',
      type: 'line',
      smooth: true,
      color: '#d4a76a',
      areaStyle: {
        color: {
          type: 'linear', x: 0, y: 0, x2: 0, y2: 1,
          colorStops: [
            { offset: 0.2, color: 'rgba(212, 167, 106, 0.3)' },
            { offset: 1, color: 'rgba(212, 167, 106, 0.02)' }
          ]
        }
      },
      emphasis: { focus: 'series' },
      data: []
    }
  ]
}));

function getDateStr(daysAgo: number): string {
  const d = new Date(); d.setDate(d.getDate() - daysAgo);
  return `${d.getMonth() + 1}/${d.getDate()}`;
}

async function mockData() {
  await new Promise(r => setTimeout(r, 1000));
  const dates: string[] = [];
  const dur: number[] = [];
  const cnt: number[] = [];
  for (let i = 13; i >= 0; i--) {
    dates.push(getDateStr(i));
    dur.push(Math.floor(Math.random() * 60) + 20);
    cnt.push(Math.floor(Math.random() * 15) + 5);
  }
  updateOptions(opts => {
    opts.xAxis.data = dates;
    opts.series[0].data = dur;
    opts.series[1].data = cnt;
    return opts;
  });
}

async function init() { mockData(); }
watch(() => appStore.locale, () => {});
init();
</script>

<template>
  <NCard :bordered="false" class="card-wrapper" size="small">
    <template #header>{{ $t('page.home.importTrend') }}</template>
    <div ref="domRef" class="h-300px overflow-hidden"></div>
  </NCard>
</template>

<style scoped></style>