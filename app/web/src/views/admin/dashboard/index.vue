<script setup lang="ts">
import { computed, ref } from 'vue'
import { VueDraggable } from 'vue-draggable-plus'
import { useAppStore } from '@/store/modules/app'
import { useThemeStore } from '@/store/modules/theme'
import { useDashboardLayout, type ModuleDef } from './composables/useDashboardLayout'

const appStore = useAppStore()
const themeStore = useThemeStore()
const dark = computed(() => themeStore.darkMode)

const {
  isEditing,
  visibleModules,
  hiddenModules,
  rows,
  removeModule,
  addModule,
  setModuleSize,
  resetLayout,
  toggleEdit,
} = useDashboardLayout()

const gap = computed(() => (appStore.isMobile ? 0 : 12))

function onDragEnd() {
  // rows auto-rebuild via watcher
}

const sizeLabels: Record<number, string> = { 6: '1/4', 8: '1/3', 12: '1/2', 24: '全宽' }
function getSizeLabel(span: number) {
  return sizeLabels[span] || `${span}/24`
}

const resizeState = ref<{
  modId: string
  startX: number
  startSpan: number
  currentTargetSpan: number
} | null>(null)

const RESIZE_OPTIONS = [6, 8, 12, 24]

function startResize(e: MouseEvent, mod: ModuleDef) {
  e.preventDefault()
  const el = e.currentTarget as HTMLElement
  const parent = el.parentElement!
  const rect = parent.getBoundingClientRect()

  resizeState.value = {
    modId: mod.id,
    startX: e.clientX,
    startSpan: mod.span,
    currentTargetSpan: mod.span,
  }

  const mousemove = (ev: MouseEvent) => {
    if (!resizeState.value) return
    const deltaX = ev.clientX - resizeState.value.startX
    const ratio = deltaX / rect.width
    let target = resizeState.value.startSpan + Math.round(ratio * 24)
    target = Math.max(6, Math.min(24, target))
    target = RESIZE_OPTIONS.reduce((a, b) => Math.abs(b - target) < Math.abs(a - target) ? b : a)
    resizeState.value.currentTargetSpan = target
  }

  const mouseup = () => {
    document.removeEventListener('mousemove', mousemove)
    document.removeEventListener('mouseup', mouseup)
    if (resizeState.value && resizeState.value.currentTargetSpan !== resizeState.value.startSpan) {
      setModuleSize(resizeState.value.modId, resizeState.value.currentTargetSpan)
    }
    resizeState.value = null
  }

  document.addEventListener('mousemove', mousemove)
  document.addEventListener('mouseup', mouseup)
}
</script>

<template>
  <div class="dashboard">
    <!-- Edit Toolbar -->
    <div class="sticky top-0 z-50 flex items-center justify-between px-4 py-2 mb-0 rd-10px"
      :class="dark ? 'bg-[#1a1a1a]/90 backdrop-blur-sm' : 'bg-white/90 backdrop-blur-sm'"
    >
      <span class="text-sm font-medium" :class="dark ? 'text-[#e0dcd6]' : 'text-[#3a3028]'">
        看板
      </span>
      <div class="flex items-center gap-2">
        <template v-if="isEditing">
          <NButton size="tiny" quaternary @click="resetLayout">
            重置
          </NButton>
          <NButton size="tiny" @click="toggleEdit">
            完成编辑
          </NButton>
        </template>
        <NButton v-else size="tiny" quaternary @click="toggleEdit">
          ✏️ 编辑
        </NButton>
      </div>
    </div>

    <!-- Add Module Panel (edit mode) -->
    <transition name="fade">
      <div v-if="isEditing && hiddenModules.length" class="flex flex-wrap gap-2 px-4 py-3 rd-10px"
        :class="dark ? 'bg-[#222] border border-[#333]' : 'bg-[#faf8f5] border border-[#e8e0d5]'"
      >
        <span class="text-xs leading-8" :class="dark ? 'text-[#888]' : 'text-[#999]'">隐藏模块：</span>
        <NButton
          v-for="mod in hiddenModules"
          :key="mod.id"
          size="tiny"
          secondary
          @click="addModule(mod.id)"
        >
          + {{ mod.name }}
        </NButton>
      </div>
    </transition>

    <!-- Edit mode badge on modules -->
    <!-- Draggable module rows -->
    <VueDraggable
      v-model="visibleModules"
      :filter="'.no-drag'"
      :prevent-on-filter="false"
      ghost-class="ghost"
      :animation="250"
      @end="onDragEnd"
    >
      <template v-for="(mod, idx) in visibleModules" :key="mod.id">
        <!-- Full-width modules -->
        <div v-if="mod.span === 24" :key="mod.id" class="relative mb-3">
          <!-- Edit mode overlay -->
          <div v-if="isEditing" class="absolute inset-0 z-10 pointer-events-none">
            <div class="flex items-start justify-end p-1.5">
              <NButton
                size="tiny"
                type="error"
                class="pointer-events-auto no-drag"
                @click="removeModule(mod.id)"
              >
                ✕ 删除
              </NButton>
            </div>
          </div>
          <!-- Resize handle on right edge -->
          <div v-if="isEditing"
            class="absolute right-0 top-0 bottom-0 w-4 z-20 cursor-ew-resize pointer-events-auto no-drag flex items-center justify-center opacity-0 hover:opacity-100 transition-opacity rounded-r-lg"
            :class="dark ? 'hover:bg-[#333]/50' : 'hover:bg-[#e8e0d5]/50'"
            @mousedown="startResize($event, mod)"
          >
            <div class="w-0.5 h-8 rounded-full bg-current opacity-40"></div>
          </div>
          <!-- Resize indicator during drag -->
          <div v-if="resizeState?.modId === mod.id"
            class="absolute inset-0 z-30 rounded-lg pointer-events-none"
            :class="dark ? 'ring-2 ring-[#b8863d]/60' : 'ring-2 ring-[#b8863d]/60'"
          >
            <div class="absolute bottom-2 right-2 text-xs px-2 py-0.5 rounded-full bg-[#b8863d] text-white font-medium shadow">
              {{ getSizeLabel(resizeState.currentTargetSpan) }}
            </div>
          </div>
          <component :is="mod.comp" />
        </div>

        <!-- Row modules (span < 24) - handled via grouping -->
        <div v-else-if="false" />
      </template>
    </VueDraggable>

    <!-- Rows for non-full-width modules -->
    <template v-for="(row, ri) in rows" :key="`row-${ri}`">
      <!-- skip full-width modules handled above -->
      <template v-if="row.every(m => m.span < 24)">
        <NGrid :x-gap="gap" :y-gap="12" responsive="screen" item-responsive>
          <NGi
            v-for="mod in row"
            :key="mod.id"
            :span="24"
            s:span="24"
            :m:span="mod.span"
            class="relative"
          >
            <!-- Edit mode overlay -->
            <div v-if="isEditing" class="absolute inset-0 z-10 pointer-events-none">
              <div class="flex items-start justify-end p-1.5">
                <NButton
                  size="tiny"
                  type="error"
                  class="pointer-events-auto no-drag"
                  @click="removeModule(mod.id)"
                >
                  ✕ 删除
                </NButton>
              </div>
            </div>
            <!-- Resize handle on right edge -->
            <div v-if="isEditing"
              class="absolute right-0 top-0 bottom-0 w-4 z-20 cursor-ew-resize pointer-events-auto no-drag flex items-center justify-center opacity-0 hover:opacity-100 transition-opacity rounded-r-lg"
              :class="dark ? 'hover:bg-[#333]/50' : 'hover:bg-[#e8e0d5]/50'"
              @mousedown="startResize($event, mod)"
            >
              <div class="w-0.5 h-8 rounded-full bg-current opacity-40"></div>
            </div>
            <!-- Resize indicator during drag -->
            <div v-if="resizeState?.modId === mod.id"
              class="absolute inset-0 z-30 rounded-lg pointer-events-none"
              :class="dark ? 'ring-2 ring-[#b8863d]/60' : 'ring-2 ring-[#b8863d]/60'"
            >
              <div class="absolute bottom-2 right-2 text-xs px-2 py-0.5 rounded-full bg-[#b8863d] text-white font-medium shadow">
                {{ getSizeLabel(resizeState.currentTargetSpan) }}
              </div>
            </div>
            <component :is="mod.comp" />
          </NGi>
        </NGrid>
      </template>
    </template>
  </div>
</template>

<style scoped>
.dashboard {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.ghost {
  opacity: 0.5;
  outline: 2px dashed #b8863d;
  outline-offset: 2px;
  border-radius: 10px;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>