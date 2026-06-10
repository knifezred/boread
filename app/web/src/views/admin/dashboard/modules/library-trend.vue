<script setup lang="ts">
import { watch } from 'vue';
import { useAppStore } from '@/store/modules/app';
import { useEcharts } from '@/hooks/common/echarts';
import { $t } from '@/locales';
defineOptions({ name: 'LibraryTrend' });

const appStore = useAppStore();

const { domRef, updateOptions } = useEcharts(() => ({
  tooltip: {
    trigger: 'axis',
    axisPointer: { type: 'cross' }
  },
  legend: {
    data: [$t('page.home.monthlyNew'), $t('page.home.cumulativeTotal')],
    top: 0
  },
  grid: {
    left: '3%',
    right: '4%',
    bottom: '3%',
    top: '18%'
  },
  xAxis: {
    type: 'category',
    data: [] as string[]
  },
  yAxis: [
    {
      type: 'value',
      name: $t('page.home.monthlyNew'),
      splitLine: { lineStyle: { color: 'rgba(180,150,110,0.12)' } }
    },
    {
      type: 'value',
      name: $t('page.home.cumulativeTotal'),
      splitLine: { show: false }
    }
  ],
  series: [
    {
      name: $t('page.home.monthlyNew'),
      type: 'bar',
      yAxisIndex: 0,
      barWidth: '50%',
      data: [] as number[],
      itemStyle: {
        borderRadius: [4, 4, 0, 0],
        color: {
          type: 'linear', x: 0, y: 0, x2: 0, y2: 1,
          colorStops: [
            { offset: 0, color: '#d4a76a' },
            { offset: 1, color: '#b8863d' }
          ]
        }
      }
    },
    {
      name: $t('page.home.cumulativeTotal'),
      type: 'line',
      yAxisIndex: 1,
      smooth: true,
      data: [] as number[],
      color: '#8b5e2b',
      areaStyle: {
        color: {
          type: 'linear', x: 0, y: 0, x2: 0, y2: 1,
          colorStops: [
            { offset: 0.2, color: 'rgba(139, 94, 43, 0.3)' },
            { offset: 1, color: 'rgba(139, 94, 43, 0.02)' }
          ]
        }
      }
    }
  ]
}));

async function mockData() {
  await new Promise(r => setTimeout(r, 800));

  const months = ['8月', '9月', '10月', '11月', '12月', '1月', '2月', '3月', '4月', '5月', '6月', '7月'];
  const monthly = [15, 22, 18, 30, 25, 20, 28, 35, 18, 40, 32, 47];
  let cumulative = 12000;
  const cumArr = monthly.map(v => { cumulative += v; return cumulative; });

  updateOptions(opts => {
    opts.xAxis.data = months;
    opts.series[0].data = monthly;
    opts.series[1].data = cumArr;
    return opts;
  });
}

function updateLocale() {
  updateOptions((opts, factory) => {
    const o = factory();
    opts.legend.data = o.legend.data;
    opts.yAxis[0].name = o.yAxis[0].name;
    opts.yAxis[1].name = o.yAxis[1].name;
    opts.series[0].name = o.series[0].name;
    opts.series[1].name = o.series[1].name;
    return opts;
  });
}

async function init() { mockData(); }

watch(() => appStore.locale, () => updateLocale());
init();
</script>

<template>
  <NCard :bordered="false" class="card-wrapper" size="small">
    <template #header>{{ $t('page.home.libraryTrend') }}</template>
    <div ref="domRef" class="h-300px overflow-hidden"></div>
  </NCard>
</template>

<style scoped></style>