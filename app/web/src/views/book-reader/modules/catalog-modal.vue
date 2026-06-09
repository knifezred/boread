<script setup lang="ts">
import { ref, computed, watch } from "vue"
import { NModal, NScrollbar } from "naive-ui"
import { fetchChapterList } from "@/service/api"
import { $t } from "@/locales"

defineOptions({ name: "CatalogModal" });

interface Props {
  bookId: string;
  chapterNo: number;
  darkMode: boolean;
}

const props = defineProps<Props>();

interface Emits {
  (e: "select", chapterNo: number): void;
}

const emit = defineEmits<Emits>();

const visible = defineModel<boolean>("visible", { default: false });

const chapters = ref<Api.BookManage.BookChapter[]>([]);

interface VolumeGroup {
  volumeNo: number | null;
  volumeTitle: string;
  chapters: Api.BookManage.BookChapter[];
}

/** 按分卷分组后的目录树 */
const volumeGroups = computed<VolumeGroup[]>(() => {
  const groups: VolumeGroup[] = [];
  let currentGroup: VolumeGroup | null = null;

  for (const ch of chapters.value) {
    if (!currentGroup || currentGroup.volumeNo !== ch.volumeNo) {
      currentGroup = {
        volumeNo: ch.volumeNo,
        volumeTitle: ch.volumeTitle || (ch.volumeNo
          ? $t("page.admin.library.book.volumeTitle", { no: ch.volumeNo })
          : $t("page.admin.library.book.mainText")),
        chapters: [],
      };
      groups.push(currentGroup);
    }
    currentGroup.chapters.push(ch);
  }
  return groups;
});

/** 展开状态，默认全部展开 */
const expandedVolumes = ref<Set<number | string>>(new Set());

watch(
  () => volumeGroups.value,
  (groups) => {
    const keys = groups.map((g) => g.volumeNo ?? "__main__");
    expandedVolumes.value = new Set(keys);
  },
  { immediate: true },
);

function toggleVolume(volumeNo: number | null) {
  const key = volumeNo ?? "__main__";
  const next = new Set(expandedVolumes.value);
  if (next.has(key)) {
    next.delete(key);
  } else {
    next.add(key);
  }
  expandedVolumes.value = next;
}

function isVolumeExpanded(volumeNo: number | null): boolean {
  return expandedVolumes.value.has(volumeNo ?? "__main__");
}

async function loadAllChapters() {
  const bookIdNum = Number(props.bookId);
  if (!bookIdNum) return;
  const { data } = await fetchChapterList(bookIdNum);
  if (data) {
    chapters.value = data;
  }
}

function handleSelect(no: number) {
  emit("select", no);
  visible.value = false;
}

watch(visible, (val) => {
  if (val) {
    loadAllChapters();
  }
});
</script>

<template>
  <NModal
    v-model:show="visible"
    preset="card"
    :title="$t('page.book.catalog.title')"
    class="lg:w-800px w-full h-100vh"
    segmented
  >
    <NScrollbar class="h-90vh">
      <div class="flex flex-col gap-2">
        <template v-for="group in volumeGroups" :key="group.volumeNo ?? '__main__'">
          <!-- 分卷标题 -->
          <div
            v-if="volumeGroups.length > 1"
            class="flex items-center gap-2 px-3 py-2 mt-1 cursor-pointer select-none text-xs font-medium uppercase tracking-wider"
            :class="darkMode ? 'text-gray-400 hover:text-gray-200' : 'text-gray-500 hover:text-gray-700'"
            @click="toggleVolume(group.volumeNo)"
          >
            <span class="text-[10px] transition-transform duration-200" :class="isVolumeExpanded(group.volumeNo) ? 'rotate-90' : ''">
              ▸
            </span>
            <span>{{ group.volumeTitle }}</span>
            <span class="text-[10px] opacity-50">({{ group.chapters.length }}章)</span>
          </div>

          <!-- 章节列表 -->
          <div v-show="volumeGroups.length <= 1 || isVolumeExpanded(group.volumeNo)" class="flex flex-col gap-0.5">
            <div
              v-for="ch in group.chapters"
              :key="ch.id"
              class="flex items-center gap-3 px-4 py-3 rd-1 cursor-pointer transition-colors duration-200 text-sm"
              :class="
                chapterNo === ch.chapterNo
                  ? 'text-primary font-medium bg-primary/5'
                  : darkMode
                    ? 'text-gray-300 hover:bg-gray-700'
                    : 'text-gray-700 hover:bg-gray-50'
              "
              @click="handleSelect(ch.chapterNo)"
            >
              <span
                class="text-xs shrink-0 w-8 text-right"
                :class="
                  chapterNo === ch.chapterNo
                    ? 'text-primary'
                    : darkMode
                      ? 'text-gray-500'
                      : 'text-gray-400'
                "
              >
                {{ ch.chapterNo }}
              </span>
              <span class="flex-1 overflow-hidden text-ellipsis whitespace-nowrap">
                {{ ch.title }}
              </span>
            </div>
          </div>
        </template>

        <!-- 空状态 -->
        <div
          v-if="chapters.length === 0"
          class="flex items-center justify-center py-12 text-sm"
          :class="darkMode ? 'text-gray-500' : 'text-gray-400'"
        >
          暂无章节
        </div>
      </div>
    </NScrollbar>
  </NModal>
</template>
