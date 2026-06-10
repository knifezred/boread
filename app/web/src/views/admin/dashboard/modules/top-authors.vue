<script setup lang="ts">
import { watch } from 'vue';
import { useAppStore } from '@/store/modules/app';
import { useEcharts } from '@/hooks/common/echarts';
import { $t } from '@/locales';
defineOptions({ name: 'TopAuthors' });

const appStore = useAppStore();

const { domRef, updateOptions } = useEcharts(() => ({
  tooltip: {
    trigger: 'axis',
    axisPointer: { type: 'shadow' }
  },
  grid: {
    left: '1%',
    right: '6%',
    bottom: '3%',
    top: '5%',
    containLabel: true
  },
  xAxis: {
    type: 'value',
    axisLabel: { show: false },
    splitLine: { show: false }
  },
  yAxis: {
    type: 'category',
    data: ['烽火戏诸侯', '爱潜水的乌贼', '忘语', '卖报小郎君', '刘慈欣', '老鹰吃小鸡'],
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
      barWidth: 14,
      data: [
        { value: 128, itemStyle: { color: '#d4a76a' } },
        { value: 96, itemStyle: { color: '#c49b5e' } },
        { value: 85, itemStyle: { color: '#b48e52' } },
        { value: 72, itemStyle: { color: '#a48246' } },
        { value: 60, itemStyle: { color: '#94753a' } },
        { value: 48, itemStyle: { color: '#84692e' } }
      ],
      label: {
        show: true,
        position: 'right',
        formatter: '{c}h',
        fontSize: 11,
        color: 'var(--n-text-color-disabled)'
      },
      itemStyle: {
        borderRadius: [0, 4, 4, 0],
        shadowBlur: 4,
        shadowColor: 'rgba(180,140,80,0.15)'
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
    <template #header>{{ $t('page.home.topAuthors') }}</template>
    <div ref="domRef" class="h-300px overflow-hidden"></div>
  </NCard>
</template>

<style scoped></style>