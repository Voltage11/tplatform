import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import type { User, TokenPair } from '@/types/user';
import * as authApi from '@/api/auth';

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const isLoading = ref(false)

  // геттеры
  const isLoggedIn = computed(() => !!user.value)
  const isAdmin = computed(() => user.value?.is_admin ?? false)

  // Установим токены в localStorage
  function saveTokens(pair: TokenPair) {
    localStorage.setItem('access_token', pair.access_token)
    localStorage.setItem('refresh_token', pair.refresh_token)
  }

  // Очистить токены
  function clearTokens() {
    localStorage.removeItem('access_token')
    localStorage.removeItem('refresh_token')
  }

  //Загрузить профиль пользователя
  async function fetchProfile() {
    try {
      isLoading.value = true
      const { user: profile } = await authApi.getProfile()
      user.value = profile
    } catch (error) {
      user.value = null
      clearTokens()
    } finally {
      isLoading.value = false
    }
  }

  // Login
  async function login(email: string, password: string) {
    isLoading.value = true
    try {
      const tokens = await authApi.login(email, password)
      saveTokens(tokens)
      await fetchProfile() // сразу получаем профиль
    } catch (error) {
      clearTokens()
      throw error // пробросим ошибку для отображения в форме
    } finally {
      isLoading.value = false
    }
  }

  // Logout
  async function logout() {
    const refreshToken = localStorage.getItem('refresh_token')
    if (refreshToken) {
      try {
        await authApi.logout(refreshToken)
      } catch (e) {}
    }
    clearTokens()
    user.value = null
  }

  // Инициализация: если есть access_token, пробуем загрузить профиль
  // initAuth() будем вызывать в main.ts при старте приложения, чтобы восстановить сессию
  async function initAuth() {
    const token = localStorage.getItem('access_token')
    if (token) {
      await fetchProfile()
    }
  }

  return {
    user,
    isLoading,
    isLoggedIn,
    isAdmin,
    login,
    logout,
    fetchProfile,
    initAuth,
  }
})
