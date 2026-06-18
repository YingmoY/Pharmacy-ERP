<template>
  <!-- 移动端扫码布局，单列、简洁 -->
  <div class="mobile-layout">
    <div class="mobile-header">
      <span>{{ route.meta.title || '智慧药店ERP' }}</span>
    </div>
    <div class="mobile-content">
      <router-view />
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRoute } from 'vue-router'
const route = useRoute()
</script>

<style scoped>
.mobile-layout {
  min-height: 100vh;
  background: #f0f2f5;
  display: flex;
  flex-direction: column;
}

.mobile-header {
  background: #001529;
  color: #fff;
  height: 50px;
  display: flex;
  align-items: center;
  padding: 0 16px;
  font-size: 16px;
  font-weight: 600;
  position: sticky;
  top: 0;
  z-index: 100;
}

.mobile-content {
  flex: 1;
  padding: 12px;
}
</style>
