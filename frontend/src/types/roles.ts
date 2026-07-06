export interface Role {
  id: string;
  name: string;
  description: string;
  created_at: string;
}

export interface RoleCreate {
  name: string;
  description: string;
}

export interface RoleUpdate {
  name: string;
  description: string;
}

export interface PaginatedResponse<T> {
  data: T[];
  pagination: {
    current_page: number;
    limit: number;
    total_count: number;
    total_pages: number;
  };
}