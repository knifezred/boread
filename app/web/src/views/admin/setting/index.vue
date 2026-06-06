<script setup lang="ts">
import { ref, reactive } from 'vue'
import { NButton, NCard, NDynamicInput, NForm, NFormItem, NInput, NInputNumber, NSpin, NTabPane, NTabs } from 'naive-ui'
import { fetchBatchSaveSettings, fetchGetSettingsByCategory } from '@/service/api'
import { $t } from '@/locales'

// ---------------------------------------------------------------------------
// Category constants
// ---------------------------------------------------------------------------
const CAT_DATABASE = 'database';
const CAT_SCAN = 'scan';
const CAT_OUTPUT = 'output';
const CAT_RECOGNITION = 'recognition';
const CAT_TAG = 'tag';
const CAT_TITLE_FORMAT = 'title_format';

// ---------------------------------------------------------------------------
// Shared loading / saving state
// ---------------------------------------------------------------------------
const activeTab = ref(CAT_DATABASE);
const loadingMap = reactive<Record<string, boolean>>({});
const savingMap = reactive<Record<string, boolean>>({});

// ---------------------------------------------------------------------------
// 1. Database config state
// ---------------------------------------------------------------------------
const db = reactive({ host: '', port: 3306, username: '', password: '', dbname: '' });

// ---------------------------------------------------------------------------
// 2. Scan directories state
// ---------------------------------------------------------------------------
const scanDirs = ref<string[]>([]);

// ---------------------------------------------------------------------------
// 3. Output directory state
// ---------------------------------------------------------------------------
const outputDir = reactive({ path: '' });

// ---------------------------------------------------------------------------
// 4. Title/Author recognition rules state
// ---------------------------------------------------------------------------
const recognition = reactive({ titlePattern: '', authorPattern: '' });

// ---------------------------------------------------------------------------
// 5. Tag extraction rules state
// ---------------------------------------------------------------------------
const tagRules = ref<{ name: string; pattern: string }[]>([]);

function addTagRule() {
  tagRules.value.push({ name: '', pattern: '' });
}

function removeTagRule(index: number) {
  tagRules.value.splice(index, 1);
}

// ---------------------------------------------------------------------------
// 6. Title format state
// ---------------------------------------------------------------------------
const titleFormat = reactive({ template: '' });

// ---------------------------------------------------------------------------
// JSON helpers
// ---------------------------------------------------------------------------
function parseJsonArray(str: string): any[] {
  try {
    return JSON.parse(str);
  } catch {
    return [];
  }
}

function getVal(map: Api.SystemManage.SettingCategoryMap, key: string, fallback = ''): string {
  return map[key]?.value ?? fallback;
}

// ---------------------------------------------------------------------------
// Apply loaded data to reactive state
// ---------------------------------------------------------------------------
function applyLoaded(cat: string, map: Api.SystemManage.SettingCategoryMap) {
  if (cat === CAT_DATABASE) {
    db.host = getVal(map, 'host', 'localhost');
    db.port = Number(getVal(map, 'port', '3306'));
    db.username = getVal(map, 'username', 'root');
    db.password = getVal(map, 'password', '');
    db.dbname = getVal(map, 'dbname', 'boread');
  } else if (cat === CAT_SCAN) {
    scanDirs.value = parseJsonArray(getVal(map, 'dirs', '[]'));
  } else if (cat === CAT_OUTPUT) {
    outputDir.path = getVal(map, 'path', '');
  } else if (cat === CAT_RECOGNITION) {
    recognition.titlePattern = getVal(map, 'title_pattern', '');
    recognition.authorPattern = getVal(map, 'author_pattern', '');
  } else if (cat === CAT_TAG) {
    tagRules.value = parseJsonArray(getVal(map, 'rules', '[]')).map((r: any) => ({
      name: r?.name ?? '',
      pattern: r?.pattern ?? '',
    }));
  } else if (cat === CAT_TITLE_FORMAT) {
    titleFormat.template = getVal(map, 'template', '');
  }
}

// ---------------------------------------------------------------------------
// Load category settings from backend
// ---------------------------------------------------------------------------
async function loadCategory(category: string) {
  if (loadingMap[category]) return;
  loadingMap[category] = true;
  try {
    const { data, error } = await fetchGetSettingsByCategory(category);
    if (!error && data) {
      applyLoaded(category, data);
    }
  } finally {
    loadingMap[category] = false;
  }
}

// ---------------------------------------------------------------------------
// Batch save settings
// ---------------------------------------------------------------------------
async function saveCategory(
  category: string,
  items: { key: string; value: string; valueType: string }[],
) {
  savingMap[category] = true;
  try {
    const { error } = await fetchBatchSaveSettings({ category, items });
    if (!error) {
      window.$message?.success($t('page.admin.system.setting.saveSuccess'));
    } else {
      window.$message?.error($t('page.admin.system.setting.saveFailed'));
    }
  } finally {
    savingMap[category] = false;
  }
}

// ---------------------------------------------------------------------------
// Save handlers for each tab
// ---------------------------------------------------------------------------
function saveDatabase() {
  return saveCategory(CAT_DATABASE, [
    { key: 'host', value: db.host, valueType: 'string' },
    { key: 'port', value: String(db.port), valueType: 'number' },
    { key: 'username', value: db.username, valueType: 'string' },
    { key: 'password', value: db.password, valueType: 'string' },
    { key: 'dbname', value: db.dbname, valueType: 'string' },
  ]);
}

function saveScanDirs() {
  return saveCategory(CAT_SCAN, [
    { key: 'dirs', value: JSON.stringify(scanDirs.value.filter(Boolean)), valueType: 'array' },
  ]);
}

function saveOutputDir() {
  return saveCategory(CAT_OUTPUT, [
    { key: 'path', value: outputDir.path, valueType: 'string' },
  ]);
}

function saveRecognition() {
  return saveCategory(CAT_RECOGNITION, [
    { key: 'title_pattern', value: recognition.titlePattern, valueType: 'string' },
    { key: 'author_pattern', value: recognition.authorPattern, valueType: 'string' },
  ]);
}

function saveTagRules() {
  const valid = tagRules.value.filter(r => r.name || r.pattern);
  return saveCategory(CAT_TAG, [
    { key: 'rules', value: JSON.stringify(valid), valueType: 'array' },
  ]);
}

function saveTitleFormat() {
  return saveCategory(CAT_TITLE_FORMAT, [
    { key: 'template', value: titleFormat.template, valueType: 'string' },
  ]);
}

// ---------------------------------------------------------------------------
// Tab change handler: lazy load on first view
// ---------------------------------------------------------------------------
const loadedTabs = reactive<Set<string>>(new Set());

function handleTabChange(name: string) {
  if (!loadedTabs.has(name)) {
    loadedTabs.add(name);
    loadCategory(name);
  }
}

// Initial load
loadedTabs.add(CAT_DATABASE);
loadCategory(CAT_DATABASE);
</script>

<template>
  <div class="min-h-500px flex-col-stretch gap-16px overflow-hidden lt-sm:overflow-auto">
    <NCard :title="$t('page.admin.system.setting.title')" :bordered="false" size="small" class="card-wrapper sm:flex-1-hidden">
      <NTabs v-model:value="activeTab" type="line" animated @update:value="handleTabChange">
        <!-- Tab 1: 数据库配置 -->
        <NTabPane name="database" :tab="$t('page.admin.system.setting.tab.database')">
          <NSpin :show="loadingMap[CAT_DATABASE]">
            <NForm label-placement="left" :label-width="120" class="max-w-600px py-16px">
              <NFormItem :label="$t('page.admin.system.setting.database.host')">
                <NInput v-model:value="db.host" :placeholder="$t('page.admin.system.setting.database.hostPlaceholder')" />
              </NFormItem>
              <NFormItem :label="$t('page.admin.system.setting.database.port')">
                <NInputNumber v-model:value="db.port" :min="1" :max="65535" :placeholder="$t('page.admin.system.setting.database.portPlaceholder')" class="w-full" />
              </NFormItem>
              <NFormItem :label="$t('page.admin.system.setting.database.username')">
                <NInput v-model:value="db.username" :placeholder="$t('page.admin.system.setting.database.usernamePlaceholder')" />
              </NFormItem>
              <NFormItem :label="$t('page.admin.system.setting.database.password')">
                <NInput v-model:value="db.password" type="password" show-password-on="click" :placeholder="$t('page.admin.system.setting.database.passwordPlaceholder')" />
              </NFormItem>
              <NFormItem :label="$t('page.admin.system.setting.database.dbname')">
                <NInput v-model:value="db.dbname" :placeholder="$t('page.admin.system.setting.database.dbnamePlaceholder')" />
              </NFormItem>
            </NForm>
            <NButton type="primary" :loading="savingMap[CAT_DATABASE]" @click="saveDatabase">
              {{ $t('page.admin.system.setting.save') }}
            </NButton>
          </NSpin>
        </NTabPane>

        <!-- Tab 2: 扫书目录 -->
        <NTabPane name="scan" :tab="$t('page.admin.system.setting.tab.scanDir')">
          <NSpin :show="loadingMap[CAT_SCAN]">
            <div class="max-w-800px py-16px">
              <NDynamicInput
                v-model:value="scanDirs"
                :placeholder="$t('page.admin.system.setting.scanDir.pathPlaceholder')"
                :min="0"
                :show-sort-button="true"
              >
                <template #create-button-default>
                  {{ $t('page.admin.system.setting.scanDir.addPath') }}
                </template>
              </NDynamicInput>
            </div>
            <NButton type="primary" :loading="savingMap[CAT_SCAN]" @click="saveScanDirs">
              {{ $t('page.admin.system.setting.save') }}
            </NButton>
          </NSpin>
        </NTabPane>

        <!-- Tab 3: 输出目录 -->
        <NTabPane name="output" :tab="$t('page.admin.system.setting.tab.outputDir')">
          <NSpin :show="loadingMap[CAT_OUTPUT]">
            <NForm label-placement="left" :label-width="120" class="max-w-600px py-16px">
              <NFormItem :label="$t('page.admin.system.setting.outputDir.path')">
                <NInput v-model:value="outputDir.path" :placeholder="$t('page.admin.system.setting.outputDir.pathPlaceholder')" />
              </NFormItem>
            </NForm>
            <NButton type="primary" :loading="savingMap[CAT_OUTPUT]" @click="saveOutputDir">
              {{ $t('page.admin.system.setting.save') }}
            </NButton>
          </NSpin>
        </NTabPane>

        <!-- Tab 4: 标题/作者识别规则 -->
        <NTabPane name="recognition" :tab="$t('page.admin.system.setting.tab.recognitionRule')">
          <NSpin :show="loadingMap[CAT_RECOGNITION]">
            <NForm label-placement="left" :label-width="120" class="max-w-600px py-16px">
              <NFormItem :label="$t('page.admin.system.setting.recognitionRule.titlePattern')">
                <NInput v-model:value="recognition.titlePattern" :placeholder="$t('page.admin.system.setting.recognitionRule.titlePatternPlaceholder')" />
              </NFormItem>
              <NFormItem :label="$t('page.admin.system.setting.recognitionRule.authorPattern')">
                <NInput v-model:value="recognition.authorPattern" :placeholder="$t('page.admin.system.setting.recognitionRule.authorPatternPlaceholder')" />
              </NFormItem>
            </NForm>
            <NButton type="primary" :loading="savingMap[CAT_RECOGNITION]" @click="saveRecognition">
              {{ $t('page.admin.system.setting.save') }}
            </NButton>
          </NSpin>
        </NTabPane>

        <!-- Tab 5: 标签提取规则 -->
        <NTabPane name="tag" :tab="$t('page.admin.system.setting.tab.tagRule')">
          <NSpin :show="loadingMap[CAT_TAG]">
            <div class="max-w-800px py-16px">
              <div v-for="(rule, index) in tagRules" :key="index" class="mb-12px flex gap-8px items-center">
                <NInput v-model:value="rule.name" :placeholder="$t('page.admin.system.setting.tagRule.namePlaceholder')" class="w-200px" />
                <NInput v-model:value="rule.pattern" :placeholder="$t('page.admin.system.setting.tagRule.patternPlaceholder')" class="flex-1" />
                <NButton quaternary type="error" @click="removeTagRule(index)">
                  {{ $t('common.delete') }}
                </NButton>
              </div>
              <NButton class="mb-16px" dashed @click="addTagRule">
                {{ $t('page.admin.system.setting.tagRule.addRule') }}
              </NButton>
            </div>
            <NButton type="primary" :loading="savingMap[CAT_TAG]" @click="saveTagRules">
              {{ $t('page.admin.system.setting.save') }}
            </NButton>
          </NSpin>
        </NTabPane>

        <!-- Tab 6: 标题格式化 -->
        <NTabPane name="title_format" :tab="$t('page.admin.system.setting.tab.titleFormat')">
          <NSpin :show="loadingMap[CAT_TITLE_FORMAT]">
            <NForm label-placement="left" :label-width="120" class="max-w-600px py-16px">
              <NFormItem :label="$t('page.admin.system.setting.titleFormat.template')">
                <NInput v-model:value="titleFormat.template" :placeholder="$t('page.admin.system.setting.titleFormat.templatePlaceholder')" />
              </NFormItem>
            </NForm>
            <NButton type="primary" :loading="savingMap[CAT_TITLE_FORMAT]" @click="saveTitleFormat">
              {{ $t('page.admin.system.setting.save') }}
            </NButton>
          </NSpin>
        </NTabPane>
      </NTabs>
    </NCard>
  </div>
</template>
