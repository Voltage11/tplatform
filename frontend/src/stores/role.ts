import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import type { Role } from '@/types/roles';
import * as roleApi from '@/api/role';

export const useRoleStore = defineStore('role', () => {
  const roles = ref<Role[]>([]);
  const isLoading = ref(false);
  const error = ref<string | null>(null);

  const currentPage = ref(1);
  const limit = ref(10);
  const totalCount = ref(0);
  const totalPages = ref(1);

  const searchName = ref('');

  const hasRoles = computed(() => roles.value.length > 0);

  async function loadRoles() {
    isLoading.value = true;
    error.value = null;
    try {
      const response = await roleApi.fetchRoles(
        currentPage.value,
        limit.value,
        searchName.value || undefined,
      );
      roles.value = response.data;
      totalCount.value = response.pagination.total_count;
      totalPages.value = response.pagination.total_pages;
      if (currentPage.value > totalPages.value && totalPages.value > 0) {
        currentPage.value = 1;
        await loadRoles();
        return;
      }
    } catch (err: any) {
      error.value = err?.response?.data?.error || 'Ошибка загрузки ролей';
    } finally {
      isLoading.value = false;
    }
  }

  async function createRole(name: string, description: string) {
    isLoading.value = true;
    try {
      await roleApi.createRole({ name, description });
      await loadRoles();
    } catch (err: any) {
      throw err;
    } finally {
      isLoading.value = false;
    }
  }

  async function updateRole(id: string, name: string, description: string) {
    isLoading.value = true;
    try {
      await roleApi.updateRole(id, { name, description });
      await loadRoles();
    } catch (err: any) {
      throw err;
    } finally {
      isLoading.value = false;
    }
  }

  async function deleteRole(id: string) {
    isLoading.value = true;
    try {
      await roleApi.deleteRole(id);
      await loadRoles();
    } catch (err: any) {
      throw err;
    } finally {
      isLoading.value = false;
    }
  }

  function setPage(page: number) {
    if (page >= 1 && page <= totalPages.value) {
      currentPage.value = page;
      loadRoles();
    }
  }

  function setSearchName(name: string) {
    searchName.value = name;
    currentPage.value = 1;
    loadRoles();
  }

  return {
    roles,
    isLoading,
    error,
    currentPage,
    limit,
    totalCount,
    totalPages,
    searchName,
    hasRoles,
    loadRoles,
    createRole,
    updateRole,
    deleteRole,
    setPage,
    setSearchName,
  };
});