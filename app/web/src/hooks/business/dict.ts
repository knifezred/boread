import { computed, ref } from 'vue';
import type { Ref } from 'vue';
import { fetchGetDictItemsByCode } from '@/service/api';
import { SelectOption } from 'naive-ui/es/select/src/interface';

const cache = new Map<string, SelectOption[]>();

export function useDictItems(code: string) {
  const options: Ref<SelectOption[]> = ref([]);
  const loading = ref(false);

  const labelMap = computed(() => {
    const map: Record<string, string> = {};
    for (const opt of options.value) {
      map[opt.value as string] = opt.label as string;
    }
    return map;
  });

  async function load() {
    if (cache.has(code)) {
      options.value = cache.get(code)!;
      return;
    }
    loading.value = true;
    const { data, error } = await fetchGetDictItemsByCode(code);
    if (!error && data) {
      const opts = data
        .filter(item => item.status === '1')
        .map(item => ({ value: item.itemValue, label: item.itemLabel }));
      cache.set(code, opts);
      options.value = opts;
    }
    loading.value = false;
  }

  load();

  return { options, loading, labelMap };
}
