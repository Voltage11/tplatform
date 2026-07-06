import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import type { Department } from '@/types/department';
import * as departmentApi from '@/api/department';

export const useDepartmentStore = defineStore('department', () => {
  const departments = ref<Department[]>([]);
  const isLoading = ref(false);
  const error = ref<string | null>(null);

  const currentPage = ref(1);
  const limit = ref(10);
  const totalCount = ref(0);
  const totalPages = ref(1);

  const searchName = ref('');

  const hasDepartments = computed(() => departments.value.length > 0);

  async function loadDepartments() {
    isLoading.value = true;
    error.value = null;
    try {
      const response = await departmentApi.fetchDepartments(
        currentPage.value,
        limit.value,
        searchName.value || undefined
      );
      departments.value = response.data;
      totalCount.value = response.pagination.total_count;
      totalPages.value = response.pagination.total_pages;
      if (currentPage.value > totalPages.value && totalPages.value > 0) {
        currentPage.value = 1;
        await loadDepartments();
        return;
      }
    } catch (err: any) {
      error.value = err?.response?.data?.error || 'Ошибка загрузки отделов';
    } finally {
      isLoading.value = false;
    }
  }

  async function createDepartment(name: string) {
    isLoading.value = true;
    // НЕ трогаем error – ошибка обрабатывается в форме
    try {
      await departmentApi.createDepartment({ name });
      await loadDepartments();
    } catch (err: any) {
      throw err; // пробрасываем, чтобы форма перехватила
    } finally {
      isLoading.value = false;
    }
  }

  async function updateDepartment(id: string, name: string) {
    isLoading.value = true;
    try {
      await departmentApi.updateDepartment(id, { name });
      await loadDepartments();
    } catch (err: any) {
      throw err;
    } finally {
      isLoading.value = false;
    }
  }

  async function deleteDepartment(id: string) {
    isLoading.value = true;
    try {
      await departmentApi.deleteDepartment(id);
      await loadDepartments();
    } catch (err: any) {
      throw err; // теперь ошибку можно поймать в ConfirmModal
    } finally {
      isLoading.value = false;
    }
  }

  function setPage(page: number) {
    if (page >= 1 && page <= totalPages.value) {
      currentPage.value = page;
      loadDepartments();
    }
  }

  function setSearchName(name: string) {
    searchName.value = name;
    currentPage.value = 1;
    loadDepartments();
  }

  return {
    departments,
    isLoading,
    error,
    currentPage,
    limit,
    totalCount,
    totalPages,
    searchName,
    hasDepartments,
    loadDepartments,
    createDepartment,
    updateDepartment,
    deleteDepartment,
    setPage,
    setSearchName,
  };
});