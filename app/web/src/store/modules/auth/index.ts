import { computed, reactive, ref } from 'vue'
import { useRoute } from 'vue-router'
import { defineStore } from 'pinia'
import { useLoading } from '@sa/hooks'
import { fetchGetUserInfo, fetchLogin } from '@/service/api'
import { fetchUgreenLogin } from '@/service/api/ugreen'
import { useRouterPush } from '@/hooks/common/router'
import { localStg } from '@/utils/storage'
import { SetupStoreId } from '@/enum'
import { $t } from '@/locales'
import { useRouteStore } from '../route'
import { useTabStore } from '../tab'
import { clearAuthStorage, getToken } from './shared'

export const useAuthStore = defineStore(SetupStoreId.Auth, () => {
  const route = useRoute();
  const authStore = useAuthStore();
  const routeStore = useRouteStore();
  const tabStore = useTabStore();
  const { toLogin, redirectFromLogin } = useRouterPush(false);
  const { loading: loginLoading, startLoading, endLoading } = useLoading();

  const token = ref('');

  const userInfo: Api.Auth.UserInfo = reactive({
    userId: '',
    userName: '',
    nickName: '',
    roles: [],
    buttons: []
  });

  /** is super role in static route */
  const isStaticSuper = computed(() => {
    const { VITE_AUTH_ROUTE_MODE, VITE_STATIC_SUPER_ROLE } = import.meta.env;

    return VITE_AUTH_ROUTE_MODE === 'static' && userInfo.roles.includes(VITE_STATIC_SUPER_ROLE);
  });

  /** Is login */
  const isLogin = computed(() => Boolean(token.value));

  /** Reset auth store */
  async function resetStore() {
    recordUserId();

    clearAuthStorage();

    authStore.$reset();

    if (!route.meta.constant) {
      await toLogin();
    }

    tabStore.cacheTabs();
    routeStore.resetStore();
  }

  /** Record the user ID of the previous login session Used to compare with the current user ID on next login */
  function recordUserId() {
    if (!userInfo.userId) {
      return;
    }

    // Store current user ID locally for next login comparison
    localStg.set('lastLoginUserId', userInfo.userId);
  }

  /**
   * Check if current login user is different from previous login user If different, clear all tabs
   *
   * @returns {boolean} Whether to clear all tabs
   */
  function checkTabClear(): boolean {
    if (!userInfo.userId) {
      return false;
    }

    const lastLoginUserId = localStg.get('lastLoginUserId');

    // Clear all tabs if current user is different from previous user
    if (!lastLoginUserId || lastLoginUserId !== userInfo.userId) {
      localStg.remove('globalTabs');
      tabStore.clearTabs();

      localStg.remove('lastLoginUserId');
      return true;
    }

    localStg.remove('lastLoginUserId');
    return false;
  }

  /**
   * Login
   *
   * @param userName User name
   * @param password Password
   * @param [redirect=true] Whether to redirect after login. Default is `true`
   */
  async function login(userName: string, password: string, redirect = true) {
    startLoading();

    const { data: loginToken, error } = await fetchLogin(userName, password);

    if (!error) {
      const pass = await loginByToken(loginToken);

      if (pass) {
        // Check if the tab needs to be cleared
        const isClear = checkTabClear();
        let needRedirect = redirect;

        if (isClear) {
          // If the tab needs to be cleared,it means we don't need to redirect.
          needRedirect = false;
        }
        await redirectFromLogin(needRedirect);

        window.$notification?.success({
          title: $t('page.login.common.loginSuccess'),
          content: $t('page.login.common.welcomeBack', { userName: userInfo.nickName }),
          duration: 4500
        });
      }
    } else {
      resetStore();
    }

    endLoading();
  }

  async function loginByToken(loginToken: Api.Auth.LoginToken) {
    // 1. stored in the localStorage, the later requests need it in headers
    localStg.set('token', loginToken.token);
    localStg.set('refreshToken', loginToken.refreshToken);

    // 2. get user info
    const pass = await getUserInfo();

    if (pass) {
      token.value = loginToken.token;

      return true;
    }

    return false;
  }

  async function getUserInfo() {
    const { data: info, error } = await fetchGetUserInfo();

    if (!error) {
      // update store
      Object.assign(userInfo, info);

      return true;
    }

    return false;
  }

  async function initUserInfo() {
    const maybeToken = getToken();

    if (maybeToken) {
      token.value = maybeToken;
      const pass = await getUserInfo();

      if (pass) return;
      resetStore();
    }

    // 无有效 token，尝试绿联NAS自动登录
    await loginByUgreen(false);
  }

  /** 绿联NAS登录 (在 UGOS 环境下自动获取 token 完成认证) */
  async function loginByUgreen(redirect = true) {
    // setupUGOSCore 已在 main.ts 中提前运行并将 token 存入 localStorage,
    // 这里直接从缓存读取，避免动态 import cloudWindow 在有些环境下失败
    const ugToken = localStg.get('ugreenToken') as string | null;
    if (!ugToken) return;

    try {
      const { data: loginToken, error } = await fetchUgreenLogin();

      if (!error && loginToken) {
        const pass = await loginByToken(loginToken);

        if (pass) {
          const isClear = checkTabClear();
          let needRedirect = redirect;

          if (isClear) {
            needRedirect = false;
          }
          await redirectFromLogin(needRedirect);
        }
      }
    } catch (err) {
      const errMsg = err instanceof Error ? err.message : String(err);
      window.$message?.error(`绿联账户认证失败: ${errMsg}`);
      // 非 UGOS 环境或认证失败，静默降级到普通登录页
    }
  }

  return {
    token,
    userInfo,
    isStaticSuper,
    isLogin,
    loginLoading,
    resetStore,
    login,
    loginByUgreen,
    initUserInfo
  };
});
