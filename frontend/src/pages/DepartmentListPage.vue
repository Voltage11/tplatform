
<script setup lang="ts">
import { onMounted, ref, watch } from 'vue';
import { useDepartmentStore } from '@/stores/department';
import DepartmentFormModal from '@/components/department/DepartmentFormModal.vue';
import ConfirmModal from '@/components/common/ConfirmModal.vue';
const deleteError = ref<string | null>(null);

const store = useDepartmentStore();

// Состояние модальных окон
const showFormModal = ref(false);
const editingDepartment = ref<{ id: string; name: string } | null>(null);
const showDeleteModal = ref(false);
const departmentToDelete = ref<{ id: string; name: string } | null>(null);

// Загрузка при монтировании
onMounted(() => {
  store.loadDepartments();
});

// Обработчики
function openCreateModal() {
  editingDepartment.value = null;
  showFormModal.value = true;
}

function openEditModal(department: { id: string; name: string }) {
  editingDepartment.value = { ...department };
  showFormModal.value = true;
}

function openDeleteModal(department: { id: string; name: string }) {
  departmentToDelete.value = department;
  showDeleteModal.value = true;
}

async function confirmDelete() {
  if (!departmentToDelete.value) return;
  deleteError.value = null;
  try {
    await store.deleteDepartment(departmentToDelete.value.id);
    showDeleteModal.value = false;
    departmentToDelete.value = null;
  } catch (err: any) {
    deleteError.value = err?.response?.data?.error || 'Не удалось удалить отдел';
  }
}

function onFormSaved() {
  // Действия после сохранения (можно просто закрыть)
}

// Поиск с простым debounce
let debounceTimer: ReturnType<typeof setTimeout> | null = null;
function onSearchInput(event: Event) {
  const value = (event.target as HTMLInputElement).value;
  if (debounceTimer) clearTimeout(debounceTimer);
  debounceTimer = setTimeout(() => {
    store.setSearchName(value);
  }, 300);
}
</script>

<template>
  <div class="department-list-page">
    <div class="page-header">
      <h2>Отделы</h2>
      <button class="btn-add" @click="openCreateModal">+ Добавить отдел</button>
    </div>

    <!-- Поиск -->
    <div class="search-bar">
      <input
        type="text"
        placeholder="Поиск по названию..."
        :value="store.searchName"
        @input="onSearchInput"
      />
    </div>

    <!-- Индикация загрузки -->
    <div v-if="store.isLoading" class="loading">Загрузка...</div>

    <!-- Ошибка -->
    <div v-else-if="store.error" class="error-message">
      {{ store.error }}
    </div>

    <!-- Таблица -->
    <table v-else-if="store.hasDepartments" class="department-table">
      <thead>
        <tr>
          <th>Название</th>
          <th>Дата создания</th>
          <th>Дата обновления</th>
          <th>Действия</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="dept in store.departments" :key="dept.id">
          <td>{{ dept.name }}</td>
          <td>{{ new Date(dept.created_at).toLocaleString() }}</td>
          <td>{{ new Date(dept.updated_at).toLocaleString() }}</td>
          <td class="actions">
            <button class="btn-edit" @click="openEditModal({ id: dept.id, name: dept.name })">
              Редактировать
            </button>
            <button class="btn-delete" @click="openDeleteModal({ id: dept.id, name: dept.name })">
              Удалить
            </button>
          </td>
        </tr>
      </tbody>
    </table>

    <!-- Пустой список -->
    <div v-else class="empty">Нет отделов</div>

    <!-- Пагинация -->
    <div v-if="store.totalPages > 1" class="pagination">
      <button :disabled="store.currentPage === 1" @click="store.setPage(store.currentPage - 1)">
        ←
      </button>
      <span>Страница {{ store.currentPage }} из {{ store.totalPages }}</span>
      <button
        :disabled="store.currentPage >= store.totalPages"
        @click="store.setPage(store.currentPage + 1)"
      >
        →
      </button>
    </div>

    <!-- Модальное окно создания/редактирования -->
    <DepartmentFormModal
      :visible="showFormModal"
      :editing-department="editingDepartment"
      @update:visible="showFormModal = $event"
      @saved="onFormSaved"
    />

    <!-- Модальное окно подтверждения удаления -->
    <ConfirmModal
      :visible="showDeleteModal"
      title="Подтверждение удаления"
      :message="`Вы уверены, что хотите удалить отдел «${departmentToDelete?.name}»?`"
      confirm-text="Удалить"
      :is-loading="store.isLoading"
      :error-message="deleteError"
      @update:visible="showDeleteModal = $event"
      @confirm="confirmDelete"
      @cancel="showDeleteModal = false; departmentToDelete = null; deleteError = null"
    />
  </div>
</template>

<style scoped>
.department-list-page {
  max-width: 800px;
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
.btn-add:hover {
  background-color: #059669;
}
.search-bar {
  margin-bottom: 16px;
}
.search-bar input {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid #cbd5e1;
  border-radius: 8px;
}
.loading,
.empty {
  text-align: center;
  padding: 40px;
  color: #64748b;
}
.error-message {
  background: #fee2e2;
  color: #b91c1c;
  padding: 12px;
  border-radius: 8px;
  margin-bottom: 16px;
}
.department-table {
  width: 100%;
  border-collapse: collapse;
  background: white;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0,0,0,0.05);
}
.department-table th,
.department-table td {
  padding: 12px;
  text-align: left;
  border-bottom: 1px solid #e2e8f0;
}
.department-table th {
  background-color: #f8fafc;
  font-weight: 600;
}
.actions {
  display: flex;
  gap: 8px;
}
.btn-edit,
.btn-delete {
  padding: 6px 12px;
  border: none;
  border-radius: 6px;
  font-size: 0.85rem;
  cursor: pointer;
}
.btn-edit {
  background: #3b82f6;
  color: white;
}
.btn-edit:hover {
  background: #2563eb;
}
.btn-delete {
  background: #ef4444;
  color: white;
}
.btn-delete:hover {
  background: #dc2626;
}
.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 16px;
  margin-top: 24px;
}
.pagination button {
  padding: 6px 12px;
  background: #e2e8f0;
  border: none;
  border-radius: 6px;
  cursor: pointer;
}
.pagination button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
</style>