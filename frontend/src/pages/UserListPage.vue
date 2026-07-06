
<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { useUserStore } from '@/stores/user';
import UserFormModal from '@/components/user/UserFormModal.vue';
import ConfirmModal from '@/components/common/ConfirmModal.vue';
import type { User } from '@/types/user';

const store = useUserStore();

const showFormModal = ref(false);
const editingUser = ref<User | null>(null);
const showDeleteModal = ref(false);
const userToDelete = ref<User | null>(null);
const deleteError = ref<string | null>(null);
const isHardDelete = ref(false);

const searchFirstName = ref('');
const searchLastName = ref('');
const searchEmail = ref('');

onMounted(() => store.loadUsers());

function openCreateModal() {
  editingUser.value = null;
  showFormModal.value = true;
}

function openEditModal(user: User) {
  editingUser.value = { ...user };
  showFormModal.value = true;
}

function openSoftDeleteModal(user: User) {
  userToDelete.value = user;
  isHardDelete.value = false;
  showDeleteModal.value = true;
}

function openHardDeleteModal(user: User) {
  userToDelete.value = user;
  isHardDelete.value = true;
  showDeleteModal.value = true;
}

async function confirmDelete() {
  if (!userToDelete.value) return;
  deleteError.value = null;
  try {
    if (isHardDelete.value) {
      await store.hardDeleteUser(userToDelete.value.id);
    } else {
      await store.softDeleteUser(userToDelete.value.id);
    }
    showDeleteModal.value = false;
    userToDelete.value = null;
  } catch (err: any) {
    deleteError.value = err?.response?.data?.error || 'Ошибка удаления';
  }
}

async function toggleActive(user: User) {
  try {
    await store.toggleActive(user.id, !user.is_active);
  } catch (e) {
    // ошибка в консоль
  }
}

function applySearch() {
  store.setFilters({
    first_name: searchFirstName.value,
    last_name: searchLastName.value,
    email: searchEmail.value,
  });
}
</script>

<template>
  <div class="user-list-page">
    <div class="page-header">
      <h2>Пользователи</h2>
      <button class="btn-add" @click="openCreateModal">+ Добавить пользователя</button>
    </div>

    <!-- Фильтры -->
    <div class="filters">
      <input v-model="searchFirstName" placeholder="Имя" @keyup.enter="applySearch" />
      <input v-model="searchLastName" placeholder="Фамилия" @keyup.enter="applySearch" />
      <input v-model="searchEmail" placeholder="Email" @keyup.enter="applySearch" />
      <button class="btn-search" @click="applySearch">Найти</button>
      <button class="btn-reset" @click="searchFirstName=''; searchLastName=''; searchEmail=''; applySearch()">Сброс</button>
    </div>

    <div v-if="store.isLoading" class="loading">Загрузка...</div>
    <div v-else-if="store.error" class="error-message">{{ store.error }}</div>

    <table v-else-if="store.hasUsers" class="user-table">
      <thead>
        <tr>
          <th>ФИО</th>
          <th>Email</th>
          <th>Отдел</th>
          <th>Роль</th>
          <th>Активен</th>
          <th>Админ</th>
          <th>Действия</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="user in store.users" :key="user.id">
          <td>{{ user.last_name }} {{ user.first_name }} {{ user.second_name }}</td>
          <td>{{ user.email }}</td>
          <td>{{ user.department?.name || '—' }}</td>
          <td>{{ user.role?.name || '—' }}</td>
          <td>
            <input type="checkbox" :checked="user.is_active" @change="toggleActive(user)" :disabled="store.isLoading" />
          </td>
          <td>{{ user.is_admin ? 'Да' : 'Нет' }}</td>
          <td class="actions">
            <button class="btn-edit" @click="openEditModal(user)">Ред.</button>
            <button class="btn-soft-delete" @click="openSoftDeleteModal(user)">Удалить</button>
            <button class="btn-hard-delete" @click="openHardDeleteModal(user)">Полное удаление</button>
          </td>
        </tr>
      </tbody>
    </table>

    <div v-else class="empty">Нет пользователей</div>

    <div v-if="store.totalPages > 1" class="pagination">
      <button :disabled="store.currentPage === 1" @click="store.setPage(store.currentPage - 1)">←</button>
      <span>Стр. {{ store.currentPage }} из {{ store.totalPages }}</span>
      <button :disabled="store.currentPage >= store.totalPages" @click="store.setPage(store.currentPage + 1)">→</button>
    </div>

    <UserFormModal
      :visible="showFormModal"
      :editing-user="editingUser"
      @update:visible="showFormModal = $event"
      @saved="showFormModal = false"
    />

    <ConfirmModal
      :visible="showDeleteModal"
      :title="isHardDelete ? 'Полное удаление' : 'Мягкое удаление'"
      :message="`Вы уверены, что хотите удалить пользователя «${userToDelete?.last_name} ${userToDelete?.first_name}»?`"
      confirm-text="Удалить"
      :is-loading="store.isLoading"
      :error-message="deleteError"
      @update:visible="showDeleteModal = $event"
      @confirm="confirmDelete"
      @cancel="showDeleteModal = false; userToDelete = null; deleteError = null"
    />
  </div>
</template>

<style scoped>
.user-list-page {
  max-width: 1100px;
  margin: 0 auto;
  padding: 24px;
}
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}
.btn-add {
  background-color: #10b981;
  color: white;
  border: none;
  padding: 10px 20px;
  border-radius: 8px;
  font-weight: 600;
  cursor: pointer;
}
.filters {
  display: flex;
  gap: 8px;
  margin-bottom: 16px;
  flex-wrap: wrap;
}
.filters input {
  padding: 6px 10px;
  border: 1px solid #cbd5e1;
  border-radius: 6px;
}
.btn-search, .btn-reset {
  padding: 6px 12px;
  border: none;
  border-radius: 6px;
  cursor: pointer;
}
.btn-search { background: #3b82f6; color: white; }
.btn-reset { background: #e2e8f0; color: #1e293b; }
.loading, .empty { text-align: center; padding: 40px; color: #64748b; }
.error-message { background: #fee2e2; color: #b91c1c; padding: 12px; border-radius: 8px; margin-bottom: 16px; }
.user-table { width: 100%; border-collapse: collapse; background: white; border-radius: 8px; overflow: hidden; box-shadow: 0 2px 8px rgba(0,0,0,0.05); }
.user-table th, .user-table td { padding: 12px; text-align: left; border-bottom: 1px solid #e2e8f0; }
.user-table th { background-color: #f8fafc; font-weight: 600; }
.actions { display: flex; gap: 4px; flex-wrap: wrap; }
.btn-edit, .btn-soft-delete, .btn-hard-delete {
  padding: 4px 8px;
  border: none;
  border-radius: 4px;
  font-size: 0.8rem;
  cursor: pointer;
}
.btn-edit { background: #3b82f6; color: white; }
.btn-soft-delete { background: #ef4444; color: white; }
.btn-hard-delete { background: #b91c1c; color: white; }
.pagination { display: flex; justify-content: center; align-items: center; gap: 16px; margin-top: 24px; }
.pagination button { padding: 6px 12px; background: #e2e8f0; border: none; border-radius: 6px; cursor: pointer; }
.pagination button:disabled { opacity: 0.5; cursor: not-allowed; }
</style>