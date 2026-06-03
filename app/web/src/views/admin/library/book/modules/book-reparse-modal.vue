<script setup lang="ts">
import { ref } from "vue";
import { NButton, NCard, NModal, NRadioGroup, NRadio, NSpace, NSpin, NAlert, NScrollbar } from "naive-ui";
import { fetchReParseChapters, fetchGetChapterRuleList, fetchGetBoundChapterRule, fetchBindChapterRule } from "@/service/api";
import { $t } from "@/locales";

defineOptions({ name: "BookReparseModal" });

interface Props { bookId: number; bookTitle: string; }
const props = defineProps<Props>();

interface Emits {
  (e: "close"): void;
  (e: "reparsed"): void;
}
const emit = defineEmits<Emits>();

const visible = defineModel<boolean>("visible", { default: false });

const rulesLoading = ref(false);
const selectedRuleId = ref<number>(0);
const ruleList = ref<Api.BookManage.BookChapterRule[]>([]);
const actionLoading = ref(false);

async function loadRules() {
  rulesLoading.value = true;
  try {
    // 获取可用规则（系统规则 + 当前用户的规则）
    const { data: ruleData } = await fetchGetChapterRuleList({
      current: 1, size: 999, ruleName: null, ruleType: null, userId: null, status: "1",
    });
    ruleList.value = ruleData?.records || [];

    // 获取当前绑定规则
    const { data: bound } = await fetchGetBoundChapterRule(props.bookId);
    if (bound) {
      selectedRuleId.value = bound.ruleId;
    } else {
      selectedRuleId.value = 0;
    }
  } catch {
    // ignore
  } finally {
    rulesLoading.value = false;
  }
}

async function handleReparse() {
  actionLoading.value = true;
  try {
    // 如果有选中规则且与当前绑定不同，先绑定
    if (selectedRuleId.value > 0) {
      const { data: currentBound } = await fetchGetBoundChapterRule(props.bookId);
      if (!currentBound || currentBound.ruleId !== selectedRuleId.value) {
        await fetchBindChapterRule({ bookId: props.bookId, ruleId: selectedRuleId.value });
      }
    }

    const { error } = await fetchReParseChapters(props.bookId, selectedRuleId.value > 0 ? selectedRuleId.value : undefined);
    if (!error) {
      window.$message?.success($t("common.updateSuccess"));
      visible.value = false;
      emit("reparsed");
    }
  } catch {
    // ignore
  } finally {
    actionLoading.value = false;
  }
}

function closeModal() {
  visible.value = false;
}
</script>

<template>
  <NModal v-model:show="visible" preset="card" :title="$t('page.admin.library.bookChapterRule.form.reParse.title')" class="w-550px" @update:show="(val) => { if (val) loadRules(); }">
    <NScrollbar class="max-h-400px pr-12px">
      <div class="mb-16px">
        <NAlert type="info" closable>
          {{ $t("page.admin.library.bookChapterRule.form.reParse.hint") }}
        </NAlert>
      </div>
      <div v-if="rulesLoading" class="flex-center py-32px">
        <NSpin />
      </div>
      <NRadioGroup v-else v-model:value="selectedRuleId" class="w-full">
        <NSpace vertical class="w-full">
          <!-- 使用系统默认规则（不绑定任何规则） -->
          <NCard
            class="rule-card" :class="[{ 'rule-card--active': selectedRuleId === 0 }]"
            size="small"
            hoverable
            @click="selectedRuleId = 0"
          >
            <NRadio :value="0" class="w-full">
              <span class="font-medium">{{ $t("page.admin.library.bookChapterRule.form.reParse.autoRule") }}</span>
            </NRadio>
          </NCard>

          <!-- 可选规则列表 -->
          <NCard
            v-for="rule in ruleList"
            :key="rule.id"
            class="rule-card" :class="[{ 'rule-card--active': selectedRuleId === rule.id }]"
            size="small"
            hoverable
            @click="selectedRuleId = rule.id"
          >
            <NRadio :value="rule.id" class="w-full">
              <div class="flex items-center gap-8px">
                <span class="font-medium">{{ rule.ruleName }}</span>
                <span class="text-#999 text-12px">
                  {{ rule.ruleType === "1" ? $t("page.admin.library.bookChapterRule.ruleTypeSystem") : $t("page.admin.library.bookChapterRule.ruleTypeCustom") }}
                </span>
              </div>
              <div class="text-12px text-#999 mt-4px">{{ rule.titlePattern }}</div>
              <div v-if="rule.description" class="text-12px text-#999 mt-2px">{{ rule.description }}</div>
            </NRadio>
          </NCard>

          <div v-if="!rulesLoading && ruleList.length === 0" class="text-center text-#999 py-16px">
            暂无可用规则
          </div>
        </NSpace>
      </NRadioGroup>
    </NScrollbar>
    <template #footer>
      <NSpace justify="end">
        <NButton @click="closeModal">{{ $t("common.cancel") }}</NButton>
        <NButton type="primary" :loading="actionLoading" @click="handleReparse">
          {{ $t("page.admin.library.bookChapterRule.form.reParse.confirm") }}
        </NButton>
      </NSpace>
    </template>
  </NModal>
</template>

<style scoped>
.rule-card {
  cursor: pointer;
  transition: all 0.2s;
  border: 1px solid transparent;
}
.rule-card:hover {
  border-color: #2080f0;
}
.rule-card--active {
  border-color: #2080f0;
  background-color: rgba(32, 128, 240, 0.05);
}
</style>
