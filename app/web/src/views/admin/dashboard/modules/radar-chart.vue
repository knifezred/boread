<script setup lang="ts">
import { watch } from 'vue';
import { useAppStore } from '@/store/modules/app';
import { useEcharts } from '@/hooks/common/echarts';
import { $t } from '@/locales';
defineOptions({ name: 'RadarChart' });

const appStore = useAppStore();

const { domRef, updateOptions } = useEcharts(() => ({
  tooltip: {
    trigger: 'item'
  },
  legend: {
    data: [$t('page.home.readingPreferences')],
    top: 0,
    left: 'center',
    textStyle: { fontSize: 12 }
  },
  radar: {
    indicator: [
      { name: '玄幻奇幻', max: 100 },
      { name: '言情', max: 100 },
      { name: '科幻', max: 100 },
      { name: '悬疑', max: 100 },
      { name: '历史', max: 100 },
      { name: '都市', max: 100 }
    ],
    shape: 'circle',
    splitNumber: 4,
    axisName: {
      color: 'var(--n-text-color)',
      fontSize: 11
    },
    splitLine: {
      lineStyle: {
        color: 'rgba(180, 150, 110, 0.15)'
      }
    },
    splitArea: {
      areaStyle: {
        color: ['rgba(180, 150, 110, 0.02)', 'rgba(180, 150, 110, 0.05)']
      }
    },
    axisLine: {
      lineStyle: {
        color: 'rgba(180, 150, 110, 0.2)'
      }
    }
  },
  series: [
    {
      type: 'radar',
      data: [
        {
          value: [90, 70, 85, 55, 40, 65],
          name: $t('page.home.readingPreferences'),
          areaStyle: {
            color: 'rgba(180, 140, 80, 0.25)'
          },
          lineStyle: {
            color: '#b8863d',
            width: 2
          },
          itemStyle: {
            color: '#b8863d'
          }
        }
      ]
    }
  ]
}));

async function mockData() {
  await new Promise(r => setTimeout(r, 800));
}

async function init() { mockData(); }

watch(() => appStore.locale, () => {
  updateOptions((opts, factory) => {
    const o = factory();
    opts.legend.data = o.legend.data;
    opts.series[0].data[0].name = o.series[0].data[0].name;
    return opts;
  });
});

init();
</script>

<template>
  <NCard :bordered="false" class="card-wrapper" size="small">
    <template #header>{{ $t('page.home.readingPreferences') }}</template>
    <div ref="domRef" class="h-280px overflow-hidden"></div>
  </NCard>
</template>

<style scoped></style>