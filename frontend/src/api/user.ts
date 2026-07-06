
import http from './http';
import type { User, UserCreate, UserUpdate, PaginatedResponse } from '@/types/user';

const BASE = '/users';

export async function fetchUsers(params: {
  page: number;
  limit: number;
  first_name?: string;
  last_name?: string;
  email?: string;
  department_id?: string;
  role_id?: string;
  is_active?: string;
}) {
  const { data } = await http.get<PaginatedResponse<User>>(BASE, { params });
  return data;
}

export async function fetchUser(id: string) {
  const { data } = await http.get<User>(`${BASE}/${id}`);
  return data;
}

export async function createUser(user: UserCreate) {
  const { data } = await http.post<User>(BASE, user);
  return data;
}

export async function updateUser(id: string, user: UserUpdate) {
  const { data } = await http.put<User>(`${BASE}/${id}`, user);
  return data;
}

export async function softDeleteUser(id: string) {
  await http.delete(`${BASE}/${id}`);
}

export async function hardDeleteUser(id: string) {
  await http.delete(`${BASE}/${id}/permanent`);
}

export async function setUserActive(id: string, isActive: boolean) {
  await http.patch(`${BASE}/${id}/active`, { is_active: isActive });
}