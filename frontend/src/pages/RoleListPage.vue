
<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { useRoleStore } from '@/stores/role';
import RoleFormModal from '@/components/role/RoleFormModal.vue';
import ConfirmModal from '@/components/common/ConfirmModal.vue';
import PermissionEditor from '@/components/role/PermissionEditor.vue';

const store = useRoleStore();

const showFormModal = ref(false);
const editingRole = ref<{ id: string; name: string; description: string } | null>(null);
const showDeleteModal = ref(false);
const roleToDelete = ref<{ id: string; name: string } | null>(null);
const deleteError = ref<string | null>(null);

const showPermissionsModal = ref(false);
const selectedRoleForPermissions = ref<{ id: string; name: string } | null>(null);

onMounted(() => {
  store.loadRoles();
});

function openCreateModal() {
  editingRole.value = null;
  showFormModal.value = true;
}

function openEditModal(role: { id: string; name: string; description: string }) {
  editingRole.value = { ...role };
  showFormModal.value = true;
}

function openDeleteModal(role: { id: string; name: string }) {
  roleToDelete.value = role;
  showDeleteModal.value = true;
}

function openPermissionsModal(role: { id: string; name: string }) {
  selectedRoleForPermissions.value = role;
  showPermissionsModal.value = true;
}

async function confirmDelete() {
  if (!roleToDelete.value) return;
  deleteError.value = null;
  try {
    await store.deleteRole(roleToDelete.value.id);
    showDeleteModal.value = false;
    roleToDelete.value = null;
  } catch (err: any) {
    deleteError.value = err?.response?.data?.error || 'Не удалось удалить роль';
  }
}

function onFormSaved() {
  // форма закрывается сама
}

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
  <div class="role-list-page">
    <div class="page-header">
      <h2>Роли</h2>
      <button class="btn-add" @click="openCreateModal">+ Добавить роль</button>
    </div>

    <div class="search-bar">
      <input
        type="text"
        placeholder="Поиск по названию..."
        :value="store.searchName"
        @input="onSearchInput"
      />
    </div>

    <div v-if="store.isLoading" class="loading">Загрузка...</div>

    <div v-else-if="store.error" class="error-message">
      {{ store.error }}
    </div>

    <table v-else-if="store.hasRoles" class="role-table">
      <thead>
        <tr>
          <th>Название</th>
          <th>Описание</th>
          <th>Создана</th>
          <th>Действия</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="role in store.roles" :key="role.id">
          <td>{{ role.name }}</td>
          <td>{{ role.description }}</td>
          <td>{{ new Date(role.created_at).toLocaleString() }}</td>
          <td class="actions">
            <button class="btn-edit" @click="openEditModal(role)">
              Редактировать
            </button>
            <button class="btn-permissions" @click="openPermissionsModal(role)">
              Права
            </button>
            <button class="btn-delete" @click="openDeleteModal(role)">
              Удалить
            </button>
          </td>
        </tr>
      </tbody>
    </table>

    <div v-else class="empty">Нет ролей</div>

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

    <RoleFormModal
      :visible="showFormModal"
      :editing-role="editingRole"
      @update:visible="showFormModal = $event"
      @saved="onFormSaved"
    />

    <ConfirmModal
      :visible="showDeleteModal"
      title="Подтверждение удаления"
      :message="`Вы уверены, что хотите удалить роль «${roleToDelete?.name}»?`"
      confirm-text="Удалить"
      :error-message="deleteError"
      :is-loading="store.isLoading"
      @update:visible="showDeleteModal = $event"
      @confirm="confirmDelete"
      @cancel="showDeleteModal = false; roleToDelete = null; deleteError = null"
    />

    <PermissionEditor
      :visible="showPermissionsModal"
      :role-id="selectedRoleForPermissions?.id ?? null"
      :role-name="selectedRoleForPermissions?.name"
      @update:visible="showPermissionsModal = $event"
      @saved="store.loadRoles()"
    />
  </div>
</template>

<style scoped>
.role-list-page {
  max-width: 900px;
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
.role-table {
  width: 100%;
  border-collapse: collapse;
  background: white;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0,0,0,0.05);
}
.role-table th,
.role-table td {
  padding: 12px;
  text-align: left;
  border-bottom: 1px solid #e2e8f0;
}
.role-table th {
  background-color: #f8fafc;
  font-weight: 600;
}
.actions {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}
.btn-edit,
.btn-permissions,
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
.btn-permissions {
  background: #8b5cf6;
  color: white;
}
.btn-permissions:hover {
  background: #7c3aed;
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