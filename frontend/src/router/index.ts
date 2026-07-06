// src/router/index.ts
import { createRouter, createWebHistory } from 'vue-router';
import { useAuthStore } from '@/stores/auth';

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: () => import('@/pages/HomePage.vue'),
    },
    {
      path: '/departments',
      name: 'departments',
      component: () => import('@/pages/DepartmentListPage.vue'),
      meta: { requiresAuth: true }, // если нужна авторизация
    },
    {
      path: '/roles',
      name: 'roles',
      component: () => import('@/pages/RoleListPage.vue'),
      meta: { requiresAuth: true }, // если нужна авторизация
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('@/pages/LoginPage.vue'),
      beforeEnter: () => {
        // Если пользователь уже авторизован, перенаправляем на главную
        const auth = useAuthStore();
        if (auth.isLoggedIn) {
          return { name: 'home' };
        }
      },
    },
  ],
});

export default router;