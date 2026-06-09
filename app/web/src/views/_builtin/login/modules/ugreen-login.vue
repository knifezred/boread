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
  // 新用户确认后才允许创建账号
  await authStore.loginByUgreen(true, true);
}

/** 获取绿联用户信息，失败时延迟重试一次（网关首次可能未就绪） */
async function loadProfile() {
  loading.value = true;
  errorMsg.value = "";

  const { data, error } = await fetchUgreenProfile();
  if (!error && data) {
    await handleProfileLoaded(data);
    return;
  }

  // 首次失败 → 等待 1.5s 后重试
  await new Promise((r) => setTimeout(r, 1500));
  const retry = await fetchUgreenProfile();
  if (!retry.error && retry.data) {
    await handleProfileLoaded(retry.data);
    return;
  }

  // 重试仍失败
  const err = retry.error ?? error;
  errorMsg.value = (err?.message ?? $t("page.login.ugreenLogin.profileError"));
  loading.value = false;
}

/** profile 加载成功后根据是否新用户决定自动登录或显示确认 UI */
async function handleProfileLoaded(data: Api.Ugreen.UgreenProfile) {
  if (!data.isNew) {
    // 已有绑定账户 → 自动登录，无需用户操作
    profile.value = data;
    await authStore.loginByUgreen(true, false);
  } else {
    // 新用户 → 显示确认 UI，用户手动授权后才创建账号
    profile.value = data;
  }
  loading.value = false;
}

onMounted(() => {
  loadProfile();
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
      <div class="flex-col-center gap-12px w-full">
        <NButton
          type="primary"
          round
          block
          :loading="loading"
          @click="loadProfile"
        >
          重试
        </NButton>
        <NButton
          round
          block
          type="success"
          @click="handleSwitchToPwd"
        >
          切换到账号密码登录
        </NButton>
      </div>
    </template>

    <!-- 新用户：展示绿联用户信息 + 授权确认 -->
    <template v-else-if="profile">
      <NCard :bordered="true" size="small" class="w-full">
        <NFlex vertical :size="12">
          <!-- 用户头像 + 账号名 -->
          <div class="flex items-center gap-10px">
            <div class="flex-center size-40px rounded-full bg-primary/10">
              <SvgIcon icon="solar:user-linear" class="text-22px text-primary" />
            </div>
            <div class="flex-1">
              <div class="text-16px font-500">{{ profile.userName || "绿联用户" }}</div>
              <div class="text-12px text-#999">{{ profile.userType === "admin" ? "管理员" : "普通用户" }}</div>
            </div>
          </div>

          <!-- 用户 ID -->
          <div class="rounded-6px bg-#f5f5f5 px-12px py-8px dark:bg-#ffffff10">
            <div class="text-12px text-#999">绿联用户 ID</div>
            <div class="mt-2px font-mono text-14px font-500">{{ profile.userId }}</div>
          </div>

          <!-- 提示 -->
          <div class="text-12px text-#faad14">
            当前绿联账号尚未绑定本地账户，点击授权将自动创建并登录
          </div>
        </NFlex>
      </NCard>

      <!-- 授权按钮 -->
      <NButton
        type="primary"
        size="large"
        round
        block
        :loading="confirming"
        @click="handleConfirm"
      >
        授权登录并创建账号
      </NButton>

      <!-- 切换到密码登录 -->
      <NButton
        round
        block
        type="success"
        @click="handleSwitchToPwd"
      >
        切换到账号密码登录
      </NButton>
    </template>
  </div>
</template>

<style scoped></style>
