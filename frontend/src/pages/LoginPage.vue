<script setup lang="ts">
import { ref, reactive } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '@/stores/auth';

const router = useRouter();
const auth = useAuthStore();

// Реактивные данные формы
const form = reactive({
  email: '',
  password: '',
});

const errorMessage = ref<string | null>(null);
const isSubmitting = ref(false);

// Если пользователь уже авторизован, редиректим на главную
if (auth.isLoggedIn) {
  router.replace('/');
}

async function handleSubmit() {
  errorMessage.value = null;

  // Базовая валидация на клиенте
  if (!form.email || !form.password) {
    errorMessage.value = 'Пожалуйста, заполните все поля';
    return;
  }

  isSubmitting.value = true;
  try {
    await auth.login(form.email, form.password);
    // После успешного входа перенаправляем на главную (или на страницу, откуда пришли)
    router.push('/');
  } catch (err: any) {
    // Извлекаем сообщение из ответа сервера (наше API возвращает { error: "..." })
    const serverError =
      err?.response?.data?.error || 'Неверный email или пароль. Попробуйте ещё раз.';
    errorMessage.value = serverError;
  } finally {
    isSubmitting.value = false;
  }
}
</script>

<template>
  <div class="login-page">
    <h2>Вход в систему</h2>
    <form @submit.prevent="handleSubmit" class="login-form" novalidate>
      <div class="form-group">
        <label for="email">Email</label>
        <input
          id="email"
          v-model="form.email"
          type="email"
          autocomplete="email"
          placeholder="admin@example.com"
          :disabled="isSubmitting"
        />
      </div>

      <div class="form-group">
        <label for="password">Пароль</label>
        <input
          id="password"
          v-model="form.password"
          type="password"
          autocomplete="current-password"
          placeholder="••••••"
          :disabled="isSubmitting"
        />
      </div>

      <div v-if="errorMessage" class="error-message">{{ errorMessage }}</div>

      <button type="submit" :disabled="isSubmitting" class="btn-submit">
        {{ isSubmitting ? 'Вход...' : 'Войти' }}
      </button>
    </form>
  </div>
</template>

<style scoped>
.login-page {
  max-width: 400px;
  margin: 60px auto;
  padding: 32px;
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

h2 {
  text-align: center;
  margin-bottom: 24px;
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

label {
  font-weight: 500;
}

input {
  padding: 10px 12px;
  border: 1px solid #cbd5e1;
  border-radius: 8px;
  font-size: 1rem;
  transition: border-color 0.2s;
}
input:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 2px rgba(59, 130, 246, 0.2);
}

.error-message {
  background-color: #fee2e2;
  color: #b91c1c;
  padding: 10px 14px;
  border-radius: 8px;
  font-size: 0.9rem;
}

.btn-submit {
  padding: 12px;
  background-color: #3b82f6;
  color: #fff;
  border: none;
  border-radius: 8px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: background-color 0.2s;
  margin-top: 8px;
}
.btn-submit:hover {
  background-color: #2563eb;
}
.btn-submit:disabled {
  background-color: #94a3b8;
  cursor: not-allowed;
}
</style>