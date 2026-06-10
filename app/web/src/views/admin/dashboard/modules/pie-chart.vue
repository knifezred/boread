<script setup lang="ts">
import { watch } from 'vue';
import { useAppStore } from '@/store/modules/app';
import { useEcharts } from '@/hooks/common/echarts';
import { $t } from '@/locales';
defineOptions({ name: 'PieChart' });

const appStore = useAppStore();

const { domRef, updateOptions } = useEcharts(() => ({
  tooltip: { trigger: 'item' },
  legend: {
    bottom: '1%',
    left: 'center',
    itemStyle: { borderWidth: 0 }
  },
  series: [
    {
      color: ['#d4a76a', '#b8863d', '#8b5e2b', '#a0763a', '#6b4f28'],
      type: 'pie',
      radius: ['45%', '75%'],
      avoidLabelOverlap: false,
      itemStyle: {
        borderRadius: 6,
        borderColor: '#fff',
        borderWidth: 2
      },
      label: {
        show: true,
        formatter: '{b}\n{d}%',
        fontSize: 11,
        color: 'var(--n-text-color)',
        lineHeight: 16
      },
      emphasis: {
        label: { show: true, fontSize: 13, fontWeight: 'bold' },
        itemStyle: {
          shadowBlur: 10,
          shadowColor: 'rgba(180, 140, 80, 0.3)'
        }
      },
      labelLine: {
        show: true,
        lineStyle: { color: 'rgba(180, 140, 80, 0.3)' }
      },
      data: [] as { name: string; value: number }[]
    }
  ]
}));

async function mockData() {
  await new Promise(r => setTimeout(r, 800));
  updateOptions(opts => {
    opts.series[0].data = [
      { name: '玄幻奇幻', value: 35 },
      { name: '言情', value: 25 },
      { name: '科幻', value: 15 },
      { name: '悬疑', value: 15 },
      { name: '历史', value: 10 }
    ];
    return opts;
  });
}

async function init() { mockData(); }

watch(() => appStore.locale, () => { init(); });
init();
</script>

<template>
  <NCard :bordered="false" class="card-wrapper" size="small">
    <template #header>{{ $t('page.home.categoryDist') }}</template>
    <div ref="domRef" class="h-300px overflow-hidden"></div>
  </NCard>
</template>

<style scoped></style>