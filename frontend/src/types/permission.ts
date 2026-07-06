export interface PermissionAction {
  name: string;
  description: string;
  is_active: boolean;
}

export interface PermissionEntity {
  entity: {
    name: string;
    description: string;
  };
  actions: PermissionAction[];
}

export interface PermissionUpdateItem {
  entity_name: string;
  action_name: string;
}