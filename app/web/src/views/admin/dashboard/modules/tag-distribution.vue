<script setup lang="ts">
import { watch } from 'vue';
import { useAppStore } from '@/store/modules/app';
import { useEcharts } from '@/hooks/common/echarts';
import { $t } from '@/locales';
defineOptions({ name: 'TagDistribution' });

const appStore = useAppStore();

const { domRef, updateOptions } = useEcharts(() => ({
  tooltip: {
    trigger: 'axis',
    axisPointer: { type: 'shadow' }
  },
  grid: {
    left: '1%',
    right: '8%',
    bottom: '3%',
    top: '3%',
    containLabel: true
  },
  xAxis: {
    type: 'value',
    axisLabel: { show: false }
  },
  yAxis: {
    type: 'category',
    data: ['黑暗流', '种田', '脑洞', '无敌流', '重生', '后宫', '系统流', '穿越'],
    axisLine: { show: false },
    axisTick: { show: false },
    axisLabel: {
      fontSize: 11,
      fontWeight: 500,
      color: 'var(--n-text-color)'
    }
  },
  series: [
    {
      type: 'bar',
      barWidth: 12,
      data: [
        { value: 38, itemStyle: { color: '#a0763a' } },
        { value: 42, itemStyle: { color: '#a87f44' } },
        { value: 55, itemStyle: { color: '#b4894d' } },
        { value: 65, itemStyle: { color: '#bf9357' } },
        { value: 76, itemStyle: { color: '#ca9d61' } },
        { value: 85, itemStyle: { color: '#d4a76a' } },
        { value: 98, itemStyle: { color: '#dfb174' } },
        { value: 120, itemStyle: { color: '#e9bb7e' } }
      ],
      label: {
        show: true,
        position: 'right',
        formatter: '{c}本',
        fontSize: 11,
        color: 'var(--n-text-color-disabled)'
      },
      itemStyle: {
        borderRadius: [0, 4, 4, 0],
        shadowBlur: 4,
        shadowColor: 'rgba(180, 140, 80, 0.15)'
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
    <template #header>{{ $t('page.home.tagDist') }}</template>
    <div ref="domRef" class="h-300px overflow-hidden"></div>
  </NCard>
</template>

<style scoped></style>