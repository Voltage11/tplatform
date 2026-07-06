import http from './http';
import type { PermissionEntity, PermissionUpdateItem } from '@/types/permission';

export async function fetchPermissions(roleId: string): Promise<PermissionEntity[]> {
  const { data } = await http.get<PermissionEntity[]>(`/roles/${roleId}/permissions`);
  return data;
}

export async function updatePermissions(roleId: string, permissions: PermissionUpdateItem[]) {
  await http.put(`/roles/${roleId}/permissions`, { permissions });
}