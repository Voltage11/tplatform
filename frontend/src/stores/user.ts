// src/stores/user.ts
import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import type { User } from '@/types/user';
import * as userApi from '@/api/user';

export const useUserStore = defineStore('user', () => {
  const users = ref<User[]>([]);
  const isLoading = ref(false);
  const error = ref<string | null>(null);

  const currentPage = ref(1);
  const limit = ref(10);
  const totalCount = ref(0);
  const totalPages = ref(1);

  const filterFirstName = ref('');
  const filterLastName = ref('');
  const filterEmail = ref('');
  const filterDepartmentId = ref<string | null>(null);
  const filterRoleId = ref<string | null>(null);
  const filterIsActive = ref<string>('all'); // 'true', 'false', 'all'

  const hasUsers = computed(() => users.value.length > 0);

  async function loadUsers() {
    isLoading.value = true;
    error.value = null;
    try {
      const params: any = {
        page: currentPage.value,
        limit: limit.value,
      };
      if (filterFirstName.value) params.first_name = filterFirstName.value;
      if (filterLastName.value) params.last_name = filterLastName.value;
      if (filterEmail.value) params.email = filterEmail.value;
      if (filterDepartmentId.value) params.department_id = filterDepartmentId.value;
      if (filterRoleId.value) params.role_id = filterRoleId.value;
      if (filterIsActive.value !== 'all') params.is_active = filterIsActive.value;

      const response = await userApi.fetchUsers(params);
      users.value = response.data;
      totalCount.value = response.pagination.total_count;
      totalPages.value = response.pagination.total_pages;
      if (currentPage.value > totalPages.value && totalPages.value > 0) {
        currentPage.value = 1;
        await loadUsers();
      }
    } catch (err: any) {
      error.value = err?.response?.data?.error || 'Ошибка загрузки пользователей';
    } finally {
      isLoading.value = false;
    }
  }

  async function createUser(userData: userApi.UserCreate) {
    isLoading.value = true;
    try {
      await userApi.createUser(userData);
      await loadUsers();
    } catch (err: any) {
      throw err;
    } finally {
      isLoading.value = false;
    }
  }

  async function updateUser(id: string, userData: userApi.UserUpdate) {
    isLoading.value = true;
    try {
      await userApi.updateUser(id, userData);
      await loadUsers();
    } catch (err: any) {
      throw err;
    } finally {
      isLoading.value = false;
    }
  }

  async function softDeleteUser(id: string) {
    isLoading.value = true;
    try {
      await userApi.softDeleteUser(id);
      await loadUsers();
    } catch (err: any) {
      throw err;
    } finally {
      isLoading.value = false;
    }
  }

  async function hardDeleteUser(id: string) {
    isLoading.value = true;
    try {
      await userApi.hardDeleteUser(id);
      await loadUsers();
    } catch (err: any) {
      throw err;
    } finally {
      isLoading.value = false;
    }
  }

  async function toggleActive(id: string, isActive: boolean) {
    isLoading.value = true;
    try {
      await userApi.setUserActive(id, isActive);
      await loadUsers();
    } catch (err: any) {
      throw err;
    } finally {
      isLoading.value = false;
    }
  }

  function setPage(page: number) {
    if (page >= 1 && page <= totalPages.value) {
      currentPage.value = page;
      loadUsers();
    }
  }

  function setFilters(filters: {
    first_name?: string;
    last_name?: string;
    email?: string;
    department_id?: string | null;
    role_id?: string | null;
    is_active?: string;
  }) {
    filterFirstName.value = filters.first_name || '';
    filterLastName.value = filters.last_name || '';
    filterEmail.value = filters.email || '';
    filterDepartmentId.value = filters.department_id || null;
    filterRoleId.value = filters.role_id || null;
    filterIsActive.value = filters.is_active || 'all';
    currentPage.value = 1;
    loadUsers();
  }

  return {
    users,
    isLoading,
    error,
    currentPage,
    limit,
    totalCount,
    totalPages,
    filterFirstName,
    filterLastName,
    filterEmail,
    filterDepartmentId,
    filterRoleId,
    filterIsActive,
    hasUsers,
    loadUsers,
    createUser,
    updateUser,
    softDeleteUser,
    hardDeleteUser,
    toggleActive,
    setPage,
    setFilters,
  };
});