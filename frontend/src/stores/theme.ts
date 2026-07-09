import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import type { Theme } from '@/types/theme';
import * as themeApi from '@/api/theme';

export const useThemeStore = defineStore('theme', () => {
  const themes = ref<Theme[]>([]);
  const isLoading = ref(false);
  const error = ref<string | null>(null);

  const currentPage = ref(1);
  const limit = ref(10);
  const totalCount = ref(0);
  const totalPages = ref(1);

  const filterName = ref('');
  const filterIsActive = ref<string>('all');
  const filterCreatedById = ref<string | null>(null);
  const filterCreatedAtFrom = ref<string | null>(null);
  const filterCreatedAtTo = ref<string | null>(null);

  const hasThemes = computed(() => themes.value.length > 0);

  async function loadThemes() {
    isLoading.value = true;
    error.value = null;
    try {
      const params: any = {
        page: currentPage.value,
        limit: limit.value,
      };
      if (filterName.value) params.name = filterName.value; 
      if (filterIsActive.value !== 'all') params.is_active = filterIsActive.value;
      if (filterCreatedById.value) params.created_by_id = filterCreatedById.value;
      if (filterCreatedAtFrom.value) params.created_at_from = filterCreatedAtFrom.value;
      if (filterCreatedAtTo.value) params.created_at_to = filterCreatedAtTo.value;

      const response = await themeApi.fetchThemes(params);
      themes.value = response.data;
      totalCount.value = response.pagination.total_count;
      totalPages.value = response.pagination.total_pages;
      if (currentPage.value > totalPages.value && totalPages.value > 0) {
        currentPage.value = 1;
        await loadThemes();
      }
    } catch (err: any) {
      error.value = err?.response?.data?.error || 'Ошибка загрузки тем';
    } finally {
      isLoading.value = false;
    }
  }

  async function createTheme(themeData: ThemeCreate) {
    isLoading.value = true;
    try {
      await themeApi.createTheme(themeData);
      await loadThemes();
    } catch (err: any) {
      throw err;
    } finally {
      isLoading.value = false;
    }
  }

  async function updateTheme(id: string, themeData: ThemeUpdate) {
    isLoading.value = true;
    try {
      await themeApi.updateTheme(id, themeData);
      await loadThemes();
    } catch (err: any) {
      throw err;
    } finally {
      isLoading.value = false;
    }
  }

  async function deleteTheme(id: string) {
    isLoading.value = true;
    try {
      await themeApi.deleteTheme(id);
      await loadThemes();
    } catch (err: any) {
      throw err;
    } finally {
      isLoading.value = false;
    }
  }

  async function toggleActive(id: string, isActive: boolean) {
    isLoading.value = true;
    try {
      await themeApi.setThemeActive(id, isActive);
      await loadThemes();
    } catch (err: any) {
      throw err;
    } finally {
      isLoading.value = false;
    }
  }

  function setPage(page: number) {
    if (page >= 1 && page <= totalPages.value) {
      currentPage.value = page;
      loadThemes();
    }
  }

  function setFilters(filters: {
    name?: string;
    is_active?: string;
    created_by_id?: string;
    created_at_from?: string;
    created_at_to?: string;
  }) {
    filterName.value = filters.name || '';
    filterIsActive.value = filters.is_active || 'all';
    filterCreatedById.value = filters.created_by_id || null;
    filterCreatedAtFrom.value = filters.created_at_from || null;
    filterCreatedAtTo.value = filters.created_at_to || null;
    currentPage.value = 1;
    loadThemes();
  }

  return {
    themes,
    isLoading,
    error,
    currentPage,
    limit,
    totalCount,
    totalPages,
    filterName,
    filterIsActive,
    filterCreatedById,
    filterCreatedAtFrom,
    filterCreatedAtTo,
    hasThemes,
    loadThemes,
    createTheme,
    updateTheme,
    deleteTheme,
    toggleActive,
    setPage,
    setFilters,
  };
});