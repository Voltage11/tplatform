<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const auth = useAuthStore()

async function handleLogout() {
  await auth.logout()
  router.push('/') // после выхода отправляем на главную
}
</script>

<template>
  <header class="app-header">
    <div class="header-left">
      <router-link to="/" class="logo">T-Platform</router-link>
      <nav class="nav-links">
        <!-- Заглушки разделов (позже заменим на router-link) -->
        <router-link to="#" class="nav-link">Тесты</router-link>
        <router-link to="/users" class="nav-link">Пользователи</router-link>
        <router-link to="/departments" class="nav-link">Отделы</router-link>
        <router-link to="/roles" class="nav-link">Роли</router-link>
        <router-link to="/themes" class="nav-link">Темы</router-link>
      </nav>
    </div>

    <div class="header-right">
      <template v-if="auth.isLoading">
        <span class="loading">Загрузка...</span>
      </template>
      <template v-else-if="auth.isLoggedIn && auth.user">
        <span class="user-name">{{ auth.user.first_name }} {{ auth.user.last_name }}</span>
        <button class="btn-logout" @click="handleLogout">Выйти</button>
      </template>
      <template v-else>
        <router-link to="/login" class="btn-login">Войти</router-link>
      </template>
    </div>
  </header>
</template>

<style scoped>
.app-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 24px;
  height: 56px;
  background-color: #1e293b;
  color: #f8fafc;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 32px;
}

.logo {
  font-size: 1.4rem;
  font-weight: bold;
  color: #f8fafc;
  text-decoration: none;
}

.nav-links {
  display: flex;
  gap: 20px;
}

.nav-link {
  color: #cbd5e1;
  text-decoration: none;
  font-weight: 500;
  transition: color 0.2s;
}
.nav-link:hover {
  color: #fff;
}
.nav-link.router-link-active {
  color: #fff;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.user-name {
  font-weight: 500;
}

.btn-logout,
.btn-login {
  background-color: #ef4444;
  border: none;
  color: white;
  padding: 6px 16px;
  border-radius: 6px;
  font-size: 0.9rem;
  cursor: pointer;
  text-decoration: none;
  transition: background-color 0.2s;
  display: inline-block;
}
.btn-logout:hover {
  background-color: #dc2626;
}
.btn-login {
  background-color: #3b82f6;
}
.btn-login:hover {
  background-color: #2563eb;
}

.loading {
  font-style: italic;
}
</style>
