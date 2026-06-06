<script setup lang="ts">
import { inject, onMounted, ref } from "vue"
import { useAuthStore } from "@/store/modules/auth"
import { useRouterPush } from "@/hooks/common/router"
import { fetchUgreenProfile } from "@/service/api/ugreen"
import { $t } from "@/locales"

defineOptions({
  name: "UgreenLogin",
});

const authStore = useAuthStore();
const { toggleLoginModule } = useRouterPush();
const disableUgreen = inject<() => void>("disableUgreen");

const loading = ref(true);
const confirming = ref(false);
const profile = ref<Api.Ugreen.UgreenProfile | null>(null);
const errorMsg = ref("");

function handleSwitchToPwd() {
  disableUgreen?.();
  toggleLoginModule("pwd-login");
}

async function handleConfirm() {
  confirming.value = true;
  await authStore.loginByUgreen();
}

onMounted(async () => {
  try {
    const { data, error } = await fetchUgreenProfile();
    if (error || !data) {
      window.$message?.error("获取用户信息失败" + error.message);
      return;
    }
    profile.value = data;
  } catch {
    window.$message?.error("获取用户信息失败");
  } finally {
    loading.value = false;
  }
});
</script>

<template>
  <div class="flex-col-center gap-20px py-24px">
    <!-- 加载中 -->
    <template v-if="loading">
      <NFlex vertical align="center" justify="center" :size="16">
        <NSpin size="large" />
        <span class="text-14px text-#999">{{
          $t("common.ugreenAuthorizing")
        }}</span>
      </NFlex>
    </template>

    <!-- 加载失败 -->
    <template v-else-if="errorMsg">
      <NResult status="error" :title="errorMsg" size="small" />
      <div class="flex-y-center justify-between gap-12px">
        <NButton
          class="flex-1"
          round
          block
          type="success"
          @click="handleSwitchToPwd"
        >
          切换到账号密码登录
        </NButton>
      </div>
    </template>

    <!-- 用户信息 + 确认登录 -->
    <template v-else-if="profile">
      <!-- 用户信息卡片 -->
      <NCard :bordered="true" size="small" class="w-full">
        <NFlex vertical :size="8">
          <div class="flex items-center gap-10px">
            <span class="text-20px">
              <SvgIcon icon="solar:user-linear"></SvgIcon>
            </span>
            <span class="text-16px font-500">{{
              profile.userName || profile.userId
            }}</span>
          </div>
          <div class="text-12px text-#999 flex items-center gap-4px">
            <span>ID:</span>
            <span class="font-mono">{{ profile.userId }}</span>
            <span class="ml-auto">{{ profile.userType }}</span>
          </div>
          <div v-if="profile.isNew" class="text-12px text-#faad14 mt-4px">
            将为你自动创建新账号
          </div>
          <div v-else class="text-12px text-#52c41a mt-4px">已有绑定账号</div>
        </NFlex>
      </NCard>

      <!-- 确认登录按钮 -->
      <div class="flex-y-center justify-between gap-12px">
        <NButton
          type="primary"
          size="large"
          round
          block
          :loading="confirming"
          @click="handleConfirm"
        >
          {{ $t("common.ugreenConfirmLogin") }}
        </NButton>
      </div>
      <!-- 切换到密码登录 -->
      <div class="flex-y-center justify-between gap-12px">
        <NButton
          class="flex-1"
          round
          block
          type="success"
          @click="handleSwitchToPwd"
        >
          切换到账号密码登录
        </NButton>
      </div>
    </template>
  </div>
</template>

<style scoped></style>
