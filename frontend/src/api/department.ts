
import http from './http';
import type { Department, DepartmentCreate, DepartmentUpdate, PaginatedResponse } from '@/types/department';

const BASE = '/departments';

/** Получить список отделов с пагинацией и поиском по имени */
export async function fetchDepartments(page: number, limit: number, name?: string) {
  const params: Record<string, string | number> = { page, limit };
  if (name) params.name = name;
  const { data } = await http.get<PaginatedResponse<Department>>(BASE, { params });
  return data;
}

/** Получить один отдел по ID */
export async function fetchDepartment(id: string) {
  const { data } = await http.get<Department>(`${BASE}/${id}`);
  return data;
}

/** Создать новый отдел */
export async function createDepartment(department: DepartmentCreate) {
  const { data } = await http.post<Department>(BASE, department);
  return data;
}

/** Обновить отдел */
export async function updateDepartment(id: string, department: DepartmentUpdate) {
  const { data } = await http.put<Department>(`${BASE}/${id}`, department);
  return data;
}

/** Мягкое удаление отдела */
export async function deleteDepartment(id: string) {
  await http.delete(`${BASE}/${id}`);
}