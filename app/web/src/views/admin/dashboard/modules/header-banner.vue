<script setup lang="ts">
import { computed } from 'vue'
import { useAppStore } from '@/store/modules/app'
import { useAuthStore } from '@/store/modules/auth'
import { $t } from '@/locales'
defineOptions({ name: 'HeaderBanner' });

const appStore = useAppStore();
const authStore = useAuthStore();

const gap = computed(() => (appStore.isMobile ? 0 : 16));

const stats = computed(() => [
  { label: $t('page.home.bookCategories'), value: '12' },
  { label: $t('page.home.bookTags'), value: '86' },
  { label: $t('page.home.authorCount'), value: '156' },
  { label: $t('page.home.bookCharacters'), value: '1,024' }
]);

const newBooksCount = 47;
</script>

<template>
  <div class="banner">
    <div class="banner-bg">
      <div class="banner-content">
        <div class="user-section">
          <div class="avatar">
            <img src="@/assets/imgs/soybean.jpg" alt="avatar" />
          </div>
          <div class="greeting">
            <h2>{{ $t('page.home.greeting', { userName: authStore.userInfo.nickName }) }}</h2>
            <p>{{ $t('page.home.weatherDesc', { newBooks: newBooksCount }) }}</p>
          </div>
        </div>
        <div class="stats-row">
          <div v-for="s in stats" :key="s.label" class="stat-item">
            <span class="stat-value">{{ s.value }}</span>
            <span class="stat-label">{{ s.label }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.banner {
  border-radius: 8px;
  overflow: hidden;
}

.banner-bg {
  background: linear-gradient(135deg, #1a1208 0%, #2a1f0e 40%, #3d2b14 100%);
  padding: 24px 28px;
  color: #f5efe6;
}

.banner-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 16px;
}

.user-section {
  display: flex;
  align-items: center;
  gap: 16px;
}

.avatar {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  overflow: hidden;
  border: 2px solid rgba(245, 239, 230, 0.3);
  flex-shrink: 0;
}

.avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.greeting h2 {
  margin: 0;
  font-size: 17px;
  font-weight: 600;
  color: #f5efe6;
}

.greeting p {
  margin: 4px 0 0;
  font-size: 13px;
  color: rgba(245, 239, 230, 0.65);
}

.stats-row {
  display: flex;
  gap: 28px;
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2px;
}

.stat-value {
  font-size: 20px;
  font-weight: 700;
  color: #d4a76a;
  line-height: 1.2;
}

.stat-label {
  font-size: 12px;
  color: rgba(245, 239, 230, 0.6);
}

@media (max-width: 640px) {
  .banner-bg { padding: 18px 16px; }
  .banner-content { flex-direction: column; align-items: flex-start; }
  .stats-row { width: 100%; justify-content: space-around; }
  .avatar { width: 48px; height: 48px; }
  .greeting h2 { font-size: 15px; }
}
</style>