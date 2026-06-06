<script setup lang="ts">
import { ref } from "vue"
import { useRouter } from "vue-router"
import {
  NButton,
  NCard,
  NForm,
  NFormItem,
  NInput,
  NInputNumber,
  NSpace,
  useMessage,
} from "naive-ui"
import { fetchSaveDatabaseConfig } from "@/service/api"
import { $t } from "@/locales"

defineOptions({ name: "SetupPage" });

const router = useRouter();
const message = useMessage();

const saving = ref(false);
const saved = ref(false);

const formValue = ref({
  host: "127.0.0.1",
  port: 3306,
  username: "root",
  password: "",
  dbname: "boread",
});

async function handleSave() {
  saving.value = true;
  try {
    const { error } = await fetchSaveDatabaseConfig(formValue.value);
    if (error) {
      message.error(
        $t("page.setup.saveFailed") + ": " + (error.message || error),
      );
      return;
    }
    saved.value = true;
    message.success($t("page.setup.saveSuccess"));
  } finally {
    saving.value = false;
  }
}
</script>

<template>
  <div
    class="relative size-full flex-center overflow-hidden bg-#f5f5f5 dark:bg-#1e1e1e"
  >
    <NCard :bordered="false" class="relative z-4 w-auto rd-12px shadow-md">
      <div class="w-420px lt-sm:w-300px">
        <header class="text-center pb-24px">
          <SystemLogo class="size-64px lt-sm:size-48px mx-auto" />
          <h3 class="text-24px text-primary font-500 mt-12px">
            {{ $t("system.title") }}
          </h3>
        </header>
        <main>
          <template v-if="!saved">
            <h3 class="text-18px text-primary font-medium mb-16px">
              {{ $t("page.setup.title") }}
            </h3>
            <p class="text-13px text-gray-500 mb-20px">
              {{ $t("page.setup.description") }}
            </p>
            <NForm :model="formValue" label-placement="top">
              <NFormItem :label="$t('page.setup.host')" path="host">
                <NInput
                  v-model:value="formValue.host"
                  placeholder="127.0.0.1"
                />
              </NFormItem>
              <NFormItem :label="$t('page.setup.port')" path="port">
                <NInputNumber
                  v-model:value="formValue.port"
                  :min="1"
                  :max="65535"
                  class="w-full"
                />
              </NFormItem>
              <NFormItem :label="$t('page.setup.dbname')" path="dbname">
                <NInput v-model:value="formValue.dbname" placeholder="boread" />
              </NFormItem>
              <NFormItem :label="$t('page.setup.username')" path="username">
                <NInput v-model:value="formValue.username" placeholder="root" />
              </NFormItem>
              <NFormItem :label="$t('page.setup.password')" path="password">
                <NInput
                  v-model:value="formValue.password"
                  type="password"
                  show-password-on="click"
                  placeholder=""
                />
              </NFormItem>
              <NSpace vertical :size="16" class="pt-8px">
                <NButton
                  type="primary"
                  size="large"
                  round
                  block
                  :loading="saving"
                  @click="handleSave"
                >
                  {{ $t("page.setup.saveAndTest") }}
                </NButton>
              </NSpace>
            </NForm>
          </template>
          <template v-else>
            <div class="text-center py-20px">
              <div
                class="i-carbon:checkmark-filled text-48px text-green-500 mx-auto mb-16px"
              />
              <h3 class="text-20px text-primary font-500 mb-12px">
                {{ $t("page.setup.successTitle") }}
              </h3>
              <p class="text-14px text-gray-500 mb-24px">
                {{ $t("page.setup.successDescription") }}
              </p>
              <NButton
                size="large"
                round
                type="primary"
                @click="router.push({ name: 'login' })"
              >
                {{ $t("page.setup.goToLogin") }}
              </NButton>
            </div>
          </template>
        </main>
      </div>
    </NCard>
  </div>
</template>

<style scoped></style>
