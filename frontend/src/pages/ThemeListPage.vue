
<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { useThemeStore } from '@/stores/theme';
import ThemeFormModal from '@/components/theme/ThemeFormModal.vue';
import ConfirmModal from '@/components/common/ConfirmModal.vue';
import type { Theme } from '@/types/theme';

const store = useThemeStore();
const showFormModal = ref(false);
const editingTheme = ref<Theme | null>(null);
const showDeleteModal = ref(false);
const themeToDelete = ref<Theme | null>(null);
const deleteError = ref<string | null>(null);

const searchName = ref('');

onMounted(() => store.loadThemes());

function openCreateModal() {
  editingTheme.value = null;
  showFormModal.value = true;
}
function openEditModal(theme: Theme) {
  editingTheme.value = { ...theme };
  showFormModal.value = true;
}
function openDeleteModal(theme: Theme) {
  themeToDelete.value = theme;
  showDeleteModal.value = true;
}
async function confirmDelete() {
  if (!themeToDelete.value) return;
  deleteError.value = null;
  try {
    await store.deleteTheme(themeToDelete.value.id);
    showDeleteModal.value = false;
    themeToDelete.value = null;
  } catch (err: any) {
    deleteError.value = err?.response?.data?.error || 'Ошибка удаления';
  }
}
async function toggleActive(theme: Theme) {
  await store.toggleActive(theme.id, !theme.is_active);
}
function applySearch() {
  store.setFilters({ name: searchName.value });
}
</script>

<template>
  <div class="page">
    <div class="page-header">
      <h2>Темы</h2>
      <button class="btn-add" @click="openCreateModal">+ Добавить тему</button>
    </div>

    <div class="filters">
      <input v-model="searchName" placeholder="Название" @keyup.enter="applySearch" />
      <select v-model="store.filterIsActive" @change="applySearch()">
        <option value="all">Все</option>
        <option value="true">Активные</option>
        <option value="false">Неактивные</option>
      </select>
      <button class="btn-search" @click="applySearch">Найти</button>
    </div>

    <div v-if="store.isLoading" class="loading">Загрузка...</div>
    <div v-else-if="store.error" class="error-message">{{ store.error }}</div>

    <table v-else-if="store.hasThemes" class="table">
      <thead>
        <tr>
          <th>Название</th>
          <th>Активна</th>
          <th>Баллы (макс/прох)</th>
          <th>Создатель</th>
          <th>Дата</th>
          <th>Действия</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="theme in store.themes" :key="theme.id">
          <td>{{ theme.name }}</td>
          <td>
            <input type="checkbox" :checked="theme.is_active" @change="toggleActive(theme)" :disabled="store.isLoading" />
          </td>
          <td>{{ theme.max_point }} / {{ theme.check_point }}</td>
          <td>{{ theme.created_by.name }}</td>
          <td>{{ new Date(theme.created_at).toLocaleDateString() }}</td>
          <td class="actions">
            <button class="btn-edit" @click="openEditModal(theme)">Ред.</button>
            <button class="btn-delete" @click="openDeleteModal(theme)">Удалить</button>
          </td>
        </tr>
      </tbody>
    </table>

    <div v-else class="empty">Нет тем</div>

    <div v-if="store.totalPages > 1" class="pagination">
      <button :disabled="store.currentPage === 1" @click="store.setPage(store.currentPage - 1)">←</button>
      <span>Стр. {{ store.currentPage }} из {{ store.totalPages }}</span>
      <button :disabled="store.currentPage >= store.totalPages" @click="store.setPage(store.currentPage + 1)">→</button>
    </div>

    <ThemeFormModal :visible="showFormModal" :editing-theme="editingTheme" @update:visible="showFormModal = $event" @saved="showFormModal = false" />
    <ConfirmModal
      :visible="showDeleteModal"
      title="Удалить тему"
      :message="`Вы уверены, что хотите удалить тему «${themeToDelete?.name}»?`"
      confirm-text="Удалить"
      :is-loading="store.isLoading"
      :error-message="deleteError"
      @update:visible="showDeleteModal = $event"
      @confirm="confirmDelete"
      @cancel="showDeleteModal = false; themeToDelete = null; deleteError = null"
    />
  </div>
</template>

<style scoped>
.page { max-width: 1000px; margin: 0 auto; padding: 24px; }
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 24px; }
.btn-add { background: #10b981; color: white; border: none; padding: 10px 20px; border-radius: 8px; font-weight: 600; cursor: pointer; }
.filters { display: flex; gap: 8px; margin-bottom: 16px; }
.table { width: 100%; border-collapse: collapse; background: white; border-radius: 8px; overflow: hidden; }
.table th, .table td { padding: 12px; border-bottom: 1px solid #e2e8f0; }
.table th { background: #f8fafc; font-weight: 600; }
.actions { display: flex; gap: 4px; }
.btn-edit, .btn-delete { padding: 4px 8px; border: none; border-radius: 4px; font-size: 0.8rem; cursor: pointer; }
.btn-edit { background: #3b82f6; color: white; }
.btn-delete { background: #ef4444; color: white; }
.pagination { display: flex; justify-content: center; gap: 12px; margin-top: 24px; }
</style>