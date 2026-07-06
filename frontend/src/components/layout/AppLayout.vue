<script setup lang="ts">
import AppHeader from './AppHeader.vue';
</script>

<template>
  <div class="app-layout">
    <AppHeader />
    <main class="main-content">
      <slot />
    </main>
  </div>
</template>

<style scoped>
.app-layout {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}
.main-content {
  flex: 1;
  padding: 32px;
  background-color: #f1f5f9;
}
</style>