// src/router/index.ts
import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: () => import('@/pages/HomePage.vue'),
      meta: { title: 'Главная' },  // можно добавить заголовок и сюда
    },
    {
      path: '/departments',
      name: 'departments',
      component: () => import('@/pages/DepartmentListPage.vue'),
      meta: { requiresAuth: true, title: 'Отделы' },
    },
    {
      path: '/roles',
      name: 'roles',
      component: () => import('@/pages/RoleListPage.vue'),
      meta: { requiresAuth: true, title: 'Роли' },
    },
    {
      path: '/users',
      name: 'users',
      component: () => import('@/pages/UserListPage.vue'),
      meta: { requiresAuth: true, title: 'Пользователи' },
    },
    { 
      path: '/themes', 
      name: 'themes', 
      component: () => import('@/pages/ThemeListPage.vue'), 
      meta: { requiresAuth: true, title: 'Темы' },
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('@/pages/LoginPage.vue'),
      meta: { title: 'Вход' },
      beforeEnter: () => {
        const auth = useAuthStore()
        if (auth.isLoggedIn) {
          return { name: 'home' }
        }
      },
    },
  ],
})

// Глобальный хук для установки document.title
router.afterEach((to) => {
  const baseTitle = 'T-Platform' // или 'T-Platform – система тестирования'
  if (to.meta?.title) {
    document.title = `${to.meta.title} | ${baseTitle}`
  } else {
    document.title = baseTitle
  }
})

export default router