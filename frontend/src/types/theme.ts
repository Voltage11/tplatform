export interface Theme {
  id: string;
  name: string;
  description?: string | null;
  is_active: boolean;
  created_by_id: string;
  created_at: string;
  date_begin?: string | null;
  date_end?: string | null;
  max_point: number;
  check_point: number;
  img_path?: string | null;
  created_by: { id: string; name: string };
}

export interface ThemeCreate {
  name: string;
  description?: string | null;
  is_active: boolean;
  date_begin?: string | null;
  date_end?: string | null;
  check_point: number;
  img_path?: string | null;
}

export interface ThemeUpdate {
  name: string;
  description?: string | null;
  is_active: boolean;
  date_begin?: string | null;
  date_end?: string | null;
  check_point: number;
  img_path?: string | null;
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