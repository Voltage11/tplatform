
import http from './http';
import type { Role, RoleCreate, RoleUpdate, PaginatedResponse } from '@/types/roles';

const BASE = '/roles';

/** Получить список ролей с пагинацией и поиском по имени */
export async function fetchRoles(page: number, limit: number, name?: string) {
  const params: Record<string, string | number> = { page, limit };
  if (name) params.name = name;
  const { data } = await http.get<PaginatedResponse<Role>>(BASE, { params });
  return data;
}

/** Получить роль по ID */
export async function fetchRole(id: string) {
  const { data } = await http.get<Role>(`${BASE}/${id}`);
  return data;
}

/** Создать новую роль */
export async function createRole(role: RoleCreate) {
  const { data } = await http.post<Role>(BASE, role);
  return data;
}

/** Обновить роль */
export async function updateRole(id: string, role: RoleUpdate) {
  const { data } = await http.put<Role>(`${BASE}/${id}`, role);
  return data;
}

/** Мягкое удаление роли */
export async function deleteRole(id: string) {
  await http.delete(`${BASE}/${id}`);
}