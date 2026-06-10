<script setup lang="ts">
import { watch } from 'vue';
import { useAppStore } from '@/store/modules/app';
import { useEcharts } from '@/hooks/common/echarts';
import { $t } from '@/locales';
defineOptions({ name: 'WordCountDist' });

const appStore = useAppStore();

const { domRef, updateOptions } = useEcharts(() => ({
  tooltip: { trigger: 'item' },
  legend: { bottom: '1%', left: 'center' },
  series: [{
    color: ['#d4a76a', '#b8863d', '#8b5e2b', '#a0763a'],
    type: 'pie',
    radius: ['40%', '70%'],
    itemStyle: { borderRadius: 6, borderColor: '#fff', borderWidth: 2 },
    label: { show: true, formatter: '{b}\n{d}%', fontSize: 11, color: 'var(--n-text-color)' },
    emphasis: { label: { show: true, fontSize: 13 } },
    data: [
      { name: $t('page.home.short'), value: 15 },
      { name: $t('page.home.medium'), value: 35 },
      { name: $t('page.home.long'), value: 32 },
      { name: $t('page.home.extraLong'), value: 18 }
    ]
  }]
}));

async function init() { await new Promise(r => setTimeout(r, 800)); }
watch(() => appStore.locale, () => init());
init();
</script>

<template>
  <NCard :bordered="false" class="card-wrapper" size="small">
    <template #header>{{ $t('page.home.wordCountDist') }}</template>
    <div ref="domRef" class="h-300px overflow-hidden"></div>
  </NCard>
</template>

<style scoped></style>