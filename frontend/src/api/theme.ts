import http from './http';
import type { Theme, ThemeCreate, ThemeUpdate, PaginatedResponse } from '@/types/theme';

const BASE = '/themes';

export async function fetchThemes(params: {
  page: number;
  limit: number;
  name?: string;
  is_active?: string;
  created_by_id?: string;
  created_at_from?: string; 
  created_at_to?: string;   
}) {
  const { data } = await http.get<PaginatedResponse<Theme>>(BASE, { params });
  return data;
}

export async function fetchTheme(id: string) {
  const { data } = await http.get<Theme>(`${BASE}/${id}`);
  return data;
}

export async function createTheme(theme: ThemeCreate) {
  const { data } = await http.post<Theme>(BASE, theme);
  return data;
}

export async function updateTheme(id: string, theme: ThemeUpdate) {
  const { data } = await http.put<Theme>(`${BASE}/${id}`, theme);
  return data;
}

export async function deleteTheme(id: string) {
  await http.delete(`${BASE}/${id}`);
}

export async function setThemeActive(id: string, isActive: boolean) {
  await http.patch(`${BASE}/${id}/active`, { is_active: isActive });
}