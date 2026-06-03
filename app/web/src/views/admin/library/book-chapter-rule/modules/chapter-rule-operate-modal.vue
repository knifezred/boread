<script setup lang="ts">
import { computed, ref, watch } from "vue"
import { NButton, NCollapse, NCollapseItem, NForm, NFormItem, NInput, NInputNumber, NModal, NRadioGroup, NRadio, NScrollbar, NSpace, NAlert, NTag, NText } from "naive-ui"
import { useFormRules, useNaiveForm } from "@/hooks/common/form"
import { fetchCreateChapterRule, fetchUpdateChapterRule } from "@/service/api"
import { $t } from "@/locales"

defineOptions({ name: "ChapterRuleOperateModal" });

interface Props {
  operateType: NaiveUI.TableOperateType;
  rowData?: Api.BookManage.BookChapterRule | null;
}
const props = defineProps<Props>();
interface Emits { (e: "submitted"): void }
const emit = defineEmits<Emits>();
const visible = defineModel<boolean>("visible", { default: false });
const { validate, restoreValidation } = useNaiveForm();
const { defaultRequiredRule } = useFormRules();
const submitting = ref(false);

const title = computed(() => props.operateType === "add" ? $t("page.admin.library.bookChapterRule.addRule") : $t("page.admin.library.bookChapterRule.editRule"));

const model = ref<Api.BookManage.ChapterRuleRequest>({
  ruleName: "", ruleType: "2", titlePattern: "", groupPattern: null, minChapterLen: 100, maxChapterLen: 100000, sortOrder: 0, description: null, status: "1",
});

const rules: Record<string, App.Global.FormRule> = {
  ruleName: defaultRequiredRule,
  titlePattern: defaultRequiredRule,
};

// ==================== 正则匹配测试 ====================

const sampleText = ref("");
const testResult = ref<{
  titleMatches: { text: string; groups: string[] }[];
  groupMatches: { text: string; groups: string[] }[];
  error: string | null;
} | null>(null);

function runRegexTest() {
  const titlePattern = model.value.titlePattern;
  const groupPattern = model.value.groupPattern;

  if (!titlePattern && !groupPattern) {
    window.$message?.warning("请先填写正则表达式");
    return;
  }
  if (!sampleText.value.trim()) {
    window.$message?.warning($t("page.admin.library.bookChapterRule.form.regexTest.sampleText"));
    return;
  }

  const result: { titleMatches: { text: string; groups: string[] }[]; groupMatches: { text: string; groups: string[] }[]; error: string | null } = {
    titleMatches: [],
    groupMatches: [],
    error: null,
  };

  try {
    if (titlePattern) {
      const regex = new RegExp(titlePattern, "gm");
      const matches = sampleText.value.matchAll(regex);
      for (const m of matches) {
        const groups: string[] = [];
        for (let i = 1; i < m.length; i++) {
          if (m[i] !== undefined) {
            groups.push(m[i]);
          }
        }
        result.titleMatches.push({ text: m[0], groups });
      }
    }
  } catch (e: any) {
    result.error = `${$t("page.admin.library.bookChapterRule.form.regexTest.regexError")}（titlePattern）: ${e.message}`;
    testResult.value = result;
    return;
  }

  try {
    if (groupPattern) {
      const regex = new RegExp(groupPattern, "gm");
      const matches = sampleText.value.matchAll(regex);
      for (const m of matches) {
        const groups: string[] = [];
        for (let i = 1; i < m.length; i++) {
          if (m[i] !== undefined) {
            groups.push(m[i]);
          }
        }
        result.groupMatches.push({ text: m[0], groups });
      }
    }
  } catch (e: any) {
    result.error = `${$t("page.admin.library.bookChapterRule.form.regexTest.regexError")}（groupPattern）: ${e.message}`;
    testResult.value = result;
    return;
  }

  testResult.value = result;
}

// ==================== 表单逻辑 ====================

function handleInitModel() {
  if (props.operateType === "edit" && props.rowData) {
    model.value = {
      ruleName: props.rowData.ruleName,
      ruleType: props.rowData.ruleType,
      titlePattern: props.rowData.titlePattern,
      groupPattern: props.rowData.groupPattern,
      minChapterLen: props.rowData.minChapterLen,
      maxChapterLen: props.rowData.maxChapterLen,
      sortOrder: props.rowData.sortOrder,
      description: props.rowData.description,
      status: props.rowData.status,
    };
  } else {
    model.value = { ruleName: "", ruleType: "2", titlePattern: "", groupPattern: null, minChapterLen: 100, maxChapterLen: 100000, sortOrder: 0, description: null, status: "1" };
  }
  // 重置测试结果
  sampleText.value = "";
  testResult.value = null;
}

async function handleSubmit() {
  await validate();
  submitting.value = true;
  try {
    if (props.operateType === "edit" && props.rowData) {
      await fetchUpdateChapterRule(props.rowData.id, model.value);
    } else {
      await fetchCreateChapterRule(model.value);
    }
    window.$message?.success($t("common.updateSuccess"));
    visible.value = false;
    emit("submitted");
  } catch (err: any) {
    window.$message?.error(err.message || $t("common.operateFail"));
  } finally {
    submitting.value = false;
  }
}

function closeModal() { visible.value = false; }

watch(visible, (val) => { if (val) { handleInitModel(); restoreValidation(); } });
</script>

<template>
  <NModal v-model:show="visible" :title="title" preset="card" class="w-650px" :loading="submitting">
    <NScrollbar class="max-h-520px pr-20px">
      <NForm ref="formRef" :model="model" :rules="rules" label-placement="left" :label-width="120">
        <NFormItem :label="$t('page.admin.library.bookChapterRule.ruleName')" path="ruleName">
          <NInput v-model:value="model.ruleName" :placeholder="$t('page.admin.library.bookChapterRule.form.ruleName')" />
        </NFormItem>
        <NFormItem :label="$t('page.admin.library.bookChapterRule.ruleType')" path="ruleType">
          <NRadioGroup v-model:value="model.ruleType">
            <NRadio value="1">{{ $t("page.admin.library.bookChapterRule.ruleTypeSystem") }}</NRadio>
            <NRadio value="2">{{ $t("page.admin.library.bookChapterRule.ruleTypeCustom") }}</NRadio>
          </NRadioGroup>
        </NFormItem>
        <NFormItem :label="$t('page.admin.library.bookChapterRule.titlePattern')" path="titlePattern">
          <NInput v-model:value="model.titlePattern" type="textarea" :rows="3" :placeholder="$t('page.admin.library.bookChapterRule.form.titlePattern')" />
        </NFormItem>
        <NFormItem :label="$t('page.admin.library.bookChapterRule.groupPattern')" path="groupPattern">
          <NInput v-model:value="model.groupPattern" type="textarea" :rows="2" :placeholder="$t('page.admin.library.bookChapterRule.form.groupPattern')" />
        </NFormItem>
        <NAlert type="info" closable class="mb-12px">{{ $t("page.admin.library.bookChapterRule.rulePreview") }}</NAlert>

        <!-- 正则匹配测试 -->
        <NCollapse class="mb-12px">
          <NCollapseItem :title="$t('page.admin.library.bookChapterRule.form.regexTest.title')" name="regex-test">
            <div class="flex flex-col gap-12px">
              <NInput
                v-model:value="sampleText"
                type="textarea"
                :rows="5"
                :placeholder="$t('page.admin.library.bookChapterRule.form.regexTest.sampleText')"
              />
              <NButton size="small" secondary @click="runRegexTest">
                {{ $t("page.admin.library.bookChapterRule.form.regexTest.runTest") }}
              </NButton>

              <!-- 错误提示 -->
              <NAlert v-if="testResult?.error" type="error" closable class="mt-8px">
                {{ testResult.error }}
              </NAlert>

              <!-- 标题匹配结果 -->
              <template v-if="testResult && !testResult.error">
                <div class="font-semibold text-14px mt-8px">
                  {{ $t("page.admin.library.bookChapterRule.form.regexTest.titleMatches") }}
                  <NTag size="small" :bordered="false" class="ml-8px">
                    {{ $t("page.admin.library.bookChapterRule.form.regexTest.matchCount", { count: testResult.titleMatches.length }) }}
                  </NTag>
                </div>
                <div v-if="testResult.titleMatches.length === 0" class="text-#999 text-12px">
                  {{ $t("page.admin.library.bookChapterRule.form.regexTest.noMatch") }}
                </div>
                <div v-else class="flex flex-col gap-4px">
                  <div
                    v-for="(match, idx) in testResult.titleMatches"
                    :key="'t-' + idx"
                    class="bg-#f5f7fa dark:bg-#333 rounded-4px px-10px py-6px text-13px leading-1.6"
                  >
                    <div>
                      <NTag size="tiny" type="success" :bordered="false" class="mr-6px">{{ idx + 1 }}</NTag>
                      <NText code>{{ match.text }}</NText>
                    </div>
                    <div v-if="match.groups.length > 0" class="mt-4px text-#999 text-12px">
                      {{ $t("page.admin.library.bookChapterRule.form.regexTest.capturedGroup", { group: match.groups.join(" → ") }) }}
                    </div>
                  </div>
                </div>

                <!-- 分组匹配结果 -->
                <template v-if="model.groupPattern">
                  <div class="font-semibold text-14px mt-16px">
                    {{ $t("page.admin.library.bookChapterRule.form.regexTest.groupMatches") }}
                    <NTag size="small" :bordered="false" class="ml-8px">
                      {{ $t("page.admin.library.bookChapterRule.form.regexTest.matchCount", { count: testResult.groupMatches.length }) }}
                    </NTag>
                  </div>
                  <div v-if="testResult.groupMatches.length === 0" class="text-#999 text-12px">
                    {{ $t("page.admin.library.bookChapterRule.form.regexTest.noMatch") }}
                  </div>
                  <div v-else class="flex flex-col gap-4px">
                    <div
                      v-for="(match, idx) in testResult.groupMatches"
                      :key="'g-' + idx"
                      class="bg-#f5f7fa dark:bg-#333 rounded-4px px-10px py-6px text-13px leading-1.6"
                    >
                      <div>
                        <NTag size="tiny" type="info" :bordered="false" class="mr-6px">{{ idx + 1 }}</NTag>
                        <NText code>{{ match.text }}</NText>
                      </div>
                      <div v-if="match.groups.length > 0" class="mt-4px text-#999 text-12px">
                        {{ $t("page.admin.library.bookChapterRule.form.regexTest.capturedGroup", { group: match.groups.join(" → ") }) }}
                      </div>
                    </div>
                  </div>
                </template>
              </template>
            </div>
          </NCollapseItem>
        </NCollapse>

        <NFormItem :label="$t('page.admin.library.bookChapterRule.minChapterLen')" path="minChapterLen">
          <NInputNumber v-model:value="model.minChapterLen" :min="0" class="w-full" />
        </NFormItem>
        <NFormItem :label="$t('page.admin.library.bookChapterRule.maxChapterLen')" path="maxChapterLen">
          <NInputNumber v-model:value="model.maxChapterLen" :min="0" class="w-full" />
        </NFormItem>
        <NFormItem :label="$t('page.admin.library.bookChapterRule.sortOrder')" path="sortOrder">
          <NInputNumber v-model:value="model.sortOrder" class="w-full" />
        </NFormItem>
        <NFormItem :label="$t('page.admin.library.bookChapterRule.description')" path="description">
          <NInput v-model:value="model.description" :placeholder="$t('page.admin.library.bookChapterRule.form.description')" />
        </NFormItem>
        <NFormItem :label="$t('page.admin.library.bookChapterRule.status')" path="status">
          <NRadioGroup v-model:value="model.status">
            <NRadio value="1">{{ $t("common.enable") }}</NRadio>
            <NRadio value="2">{{ $t("common.disable") }}</NRadio>
          </NRadioGroup>
        </NFormItem>
      </NForm>
    </NScrollbar>
    <template #footer>
      <NSpace justify="end">
        <NButton @click="closeModal">{{ $t("common.cancel") }}</NButton>
        <NButton type="primary" :loading="submitting" @click="handleSubmit">{{ $t("common.confirm") }}</NButton>
      </NSpace>
    </template>
  </NModal>
</template>
