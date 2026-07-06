/** Пользователь, как приходит с бэкенда (без пароля) */
export interface User {
  id: string;
  first_name: string;
  second_name?: string | null;
  last_name: string;
  email: string;
  department?: { id: string; name: string } | null;
  role?: { id: string; name: string } | null;
  is_active: boolean;
  is_admin: boolean;
  created_at: string;
  updated_at: string;
}

/** Ответ после логина */
export interface TokenPair {
  access_token: string;
  refresh_token: string;
  expires_in: number;
}

export interface UserCreate {
  first_name: string;
  second_name?: string;
  last_name: string;
  email: string;
  password: string;
  department_id?: string | null;
  role_id?: string | null;
  is_active?: boolean;
  is_admin?: boolean;
}

export interface UserUpdate {
  first_name?: string;
  second_name?: string | null;
  last_name?: string;
  email?: string;
  department_id?: string | null;
  role_id?: string | null;
  is_active?: boolean;
  is_admin?: boolean;
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