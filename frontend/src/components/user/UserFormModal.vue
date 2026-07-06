
<script setup lang="ts">
import { ref, watch, computed } from 'vue';
import { useUserStore } from '@/stores/user';
import { useFieldErrors } from '@/composables/useFieldErrors';
import AsyncSelect from '@/components/common/AsyncSelect.vue';
import * as departmentApi from '@/api/department';
import * as roleApi from '@/api/role';
import type { UserUpdate } from '@/types/user';

const props = defineProps<{
  visible: boolean;
  editingUser?: import('@/types/user').User | null;
}>();

const emit = defineEmits<{
  (e: 'update:visible', value: boolean): void;
  (e: 'saved'): void;
}>();

const isCreate = computed(() => !props.editingUser);

const firstName = ref('');
const secondName = ref('');
const lastName = ref('');
const email = ref('');
const password = ref('');
const isActive = ref(true);
const isAdmin = ref(false);
const departmentId = ref<string | null>(null);
const roleId = ref<string | null>(null);
const isSubmitting = ref(false);

const { fieldErrors, generalError, setServerError, clearErrors } = useFieldErrors();

// Заполнение формы при открытии
watch(() => props.visible, (newVal) => {
  if (newVal) {
    if (props.editingUser) {
      const u = props.editingUser;
      firstName.value = u.first_name;
      secondName.value = u.second_name || '';
      lastName.value = u.last_name;
      email.value = u.email;
      isActive.value = u.is_active;
      isAdmin.value = u.is_admin;
      departmentId.value = u.department?.id || null;
      roleId.value = u.role?.id || null;
    } else {
      firstName.value = '';
      secondName.value = '';
      lastName.value = '';
      email.value = '';
      password.value = '';
      isActive.value = true;
      isAdmin.value = false;
      departmentId.value = null;
      roleId.value = null;
    }
    clearErrors();
  }
});

async function handleSubmit() {
  if (!firstName.value.trim() || !lastName.value.trim() || !email.value.trim()) {
    fieldErrors.value = {
      first_name: !firstName.value.trim() ? 'Имя обязательно' : '',
      last_name: !lastName.value.trim() ? 'Фамилия обязательна' : '',
      email: !email.value.trim() ? 'Email обязателен' : '',
    };
    return;
  }
  if (isCreate.value && !password.value.trim()) {
    fieldErrors.value = { password: 'Пароль обязателен' };
    return;
  }

  isSubmitting.value = true;
  try {
    const store = useUserStore();
    if (isCreate.value) {
      await store.createUser({
        first_name: firstName.value,
        second_name: secondName.value || undefined,
        last_name: lastName.value,
        email: email.value,
        password: password.value,
        department_id: departmentId.value,
        role_id: roleId.value,
        is_active: isActive.value,
        is_admin: isAdmin.value,
      });
    } else {
      const updateData: UserUpdate = {
        first_name: firstName.value,
        second_name: secondName.value || null,
        last_name: lastName.value,
        email: email.value,
        department_id: departmentId.value,
        role_id: roleId.value,
        is_active: isActive.value,
        is_admin: isAdmin.value,
      };
      await store.updateUser(props.editingUser!.id, updateData);
    }
    emit('saved');
    emit('update:visible', false);
  } catch (err: any) {
    setServerError(err);
  } finally {
    isSubmitting.value = false;
  }
}

async function fetchDepartments(search: string) {
  const { data } = await departmentApi.fetchDepartments(1, 5, search || undefined);
  return data.map(d => ({ id: d.id, name: d.name }));
}

async function fetchRoles(search: string) {
  const { data } = await roleApi.fetchRoles(1, 5, search || undefined);
  return data.map(r => ({ id: r.id, name: r.name }));
}

const initialDepartments = ref<{ id: string; name: string }[]>([]);
const initialRoles = ref<{ id: string; name: string }[]>([]);

// Подгружаем первые 5 отделов и ролей при монтировании (можно сделать в watch visible)
watch(() => props.visible, async (newVal) => {
  if (newVal) {
    try {
      const deptRes = await departmentApi.fetchDepartments(1, 5);
      initialDepartments.value = deptRes.data.map(d => ({ id: d.id, name: d.name }));
      const roleRes = await roleApi.fetchRoles(1, 5);
      initialRoles.value = roleRes.data.map(r => ({ id: r.id, name: r.name }));
    } catch (e) {
      // оставляем пустыми
    }
  }
});

function handleCancel() {
  emit('update:visible', false);
}
</script>

<template>
  <div v-if="visible" class="modal-overlay" @click.self="handleCancel">
    <div class="modal-container">
      <h3>{{ isCreate ? 'Создать пользователя' : 'Редактировать пользователя' }}</h3>
      <form @submit.prevent="handleSubmit" class="user-form">
        <div class="form-row">
          <div class="form-group">
            <label>Имя *</label>
            <input v-model="firstName" type="text" :disabled="isSubmitting" :class="{ 'input-error': fieldErrors['first_name'] }" />
            <span v-if="fieldErrors['first_name']" class="field-error">{{ fieldErrors['first_name'] }}</span>
          </div>
          <div class="form-group">
            <label>Фамилия *</label>
            <input v-model="lastName" type="text" :disabled="isSubmitting" :class="{ 'input-error': fieldErrors['last_name'] }" />
            <span v-if="fieldErrors['last_name']" class="field-error">{{ fieldErrors['last_name'] }}</span>
          </div>
        </div>
        <div class="form-row">
          <div class="form-group">
            <label>Отчество</label>
            <input v-model="secondName" type="text" :disabled="isSubmitting" />
          </div>
          <div class="form-group">
            <label>Email *</label>
            <input v-model="email" type="email" :disabled="isSubmitting" :class="{ 'input-error': fieldErrors['email'] }" />
            <span v-if="fieldErrors['email']" class="field-error">{{ fieldErrors['email'] }}</span>
          </div>
        </div>
        <div v-if="isCreate" class="form-group">
          <label>Пароль *</label>
          <input v-model="password" type="password" :disabled="isSubmitting" :class="{ 'input-error': fieldErrors['password'] }" />
          <span v-if="fieldErrors['password']" class="field-error">{{ fieldErrors['password'] }}</span>
        </div>
        <div class="form-row">
          <div class="form-group">
            <label>Отдел</label>
            <AsyncSelect
              v-model="departmentId"
              :fetch-options="fetchDepartments"
              :initial-options="initialDepartments"
              placeholder="Выберите отдел"
              :disabled="isSubmitting"
            />
          </div>
          <div class="form-group">
            <label>Роль</label>
            <AsyncSelect
              v-model="roleId"
              :fetch-options="fetchRoles"
              :initial-options="initialRoles"
              placeholder="Выберите роль"
              :disabled="isSubmitting"
            />
          </div>
        </div>
        <div class="form-row">
          <div class="form-group checkbox-group">
            <label>
              <input v-model="isActive" type="checkbox" :disabled="isSubmitting" />
              Активен
            </label>
          </div>
          <div class="form-group checkbox-group">
            <label>
              <input v-model="isAdmin" type="checkbox" :disabled="isSubmitting" />
              Администратор
            </label>
          </div>
        </div>

        <div v-if="generalError" class="error-message">{{ generalError }}</div>

        <div class="modal-actions">
          <button type="button" class="btn-cancel" :disabled="isSubmitting" @click="handleCancel">Отмена</button>
          <button type="submit" class="btn-save" :disabled="isSubmitting">
            {{ isSubmitting ? 'Сохранение...' : (isCreate ? 'Создать' : 'Обновить') }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<style scoped>
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}
.modal-container {
  background: white;
  padding: 24px;
  border-radius: 12px;
  min-width: 500px;
  max-width: 600px;
  box-shadow: 0 8px 24px rgba(0,0,0,0.2);
}
h3 {
  margin-top: 0;
}
.user-form {
  display: flex;
  flex-direction: column;
  gap: 16px;
}
.form-row {
  display: flex;
  gap: 16px;
}
.form-group {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 6px;
}
label {
  font-weight: 500;
}
input[type="text"],
input[type="email"],
input[type="password"] {
  padding: 8px 12px;
  border: 1px solid #cbd5e1;
  border-radius: 8px;
  font-size: 1rem;
}
.input-error {
  border-color: #ef4444;
}
.field-error {
  color: #ef4444;
  font-size: 0.85rem;
}
.checkbox-group label {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: normal;
}
.error-message {
  background: #fee2e2;
  color: #b91c1c;
  padding: 8px;
  border-radius: 8px;
}
.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 16px;
}
.btn-cancel,
.btn-save {
  padding: 8px 20px;
  border: none;
  border-radius: 8px;
  font-weight: 600;
  cursor: pointer;
}
.btn-cancel { background: #e2e8f0; color: #1e293b; }
.btn-save { background: #3b82f6; color: white; }
.btn-save:hover { background: #2563eb; }
.btn-cancel:hover { background: #cbd5e1; }
button:disabled { opacity: 0.6; cursor: not-allowed; }
</style>