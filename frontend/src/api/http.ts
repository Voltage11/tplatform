// src/api/http.ts
import axios, { AxiosError } from 'axios';
import { useAuthStore } from '@/stores/auth';

const http = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL as string,
  headers: { 'Content-Type': 'application/json' },
});

let isRefreshing = false;
let failedQueue: Array<{
  resolve: (token: string) => void;
  reject: (error: unknown) => void;
}> = [];

function processQueue(error: unknown, token: string | null) {
  failedQueue.forEach(({ resolve, reject }) => {
    if (error) {
      reject(error);
    } else {
      resolve(token as string);
    }
  });
  failedQueue = [];
}

// Добавляем access-токен в каждый запрос
http.interceptors.request.use((config) => {
  const token = localStorage.getItem('access_token');
  if (token && config.headers) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Обработка 401 с автоматическим обновлением токена
http.interceptors.response.use(
  (response) => response,
  async (error: AxiosError) => {
    // Безопасно приводим к any, чтобы не зависеть от InternalAxiosRequestConfig
    const originalRequest = (error.config as any) || {};

    // Пропускаем, если это не 401, или уже пытались обновить, или это запрос на refresh
    if (
      error.response?.status !== 401 ||
      originalRequest._retry ||
      (originalRequest.url && originalRequest.url.includes('/auth/refresh'))
    ) {
      return Promise.reject(error);
    }

    originalRequest._retry = true;

    if (!isRefreshing) {
      isRefreshing = true;
      const refreshToken = localStorage.getItem('refresh_token');

      if (!refreshToken) {
        isRefreshing = false;
        const authStore = useAuthStore();
        await authStore.logout();
        return Promise.reject(error);
      }

      try {
        const { data } = await axios.post(
          `${import.meta.env.VITE_API_BASE_URL}/auth/refresh`,
          { refresh_token: refreshToken }
        );
        localStorage.setItem('access_token', data.access_token);
        localStorage.setItem('refresh_token', data.refresh_token);

        // Обновляем заголовок у исходного запроса
        if (originalRequest.headers) {
          originalRequest.headers.Authorization = `Bearer ${data.access_token}`;
        } else {
          originalRequest.headers = { Authorization: `Bearer ${data.access_token}` };
        }

        processQueue(null, data.access_token);
        return http(originalRequest);
      } catch (refreshError) {
        processQueue(refreshError, null);
        const authStore = useAuthStore();
        await authStore.logout();
        return Promise.reject(refreshError);
      } finally {
        isRefreshing = false;
      }
    } else {
      // Ставим запрос в очередь, если уже идёт процесс обновления токена
      return new Promise((resolve, reject) => {
        failedQueue.push({
          resolve: (token: string) => {
            if (originalRequest.headers) {
              originalRequest.headers.Authorization = `Bearer ${token}`;
            } else {
              originalRequest.headers = { Authorization: `Bearer ${token}` };
            }
            resolve(http(originalRequest));
          },
          reject: (err: unknown) => {
            reject(err);
          },
        });
      });
    }
  }
);

export default http;