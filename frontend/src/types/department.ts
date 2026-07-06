export interface Department {
  id: string;
  name: string;
  created_at: string;
  updated_at: string;
}

export interface DepartmentCreate {
  name: string;
}

export interface DepartmentUpdate {
  name: string;
}

// Ответ с пагинацией
export interface PaginatedResponse<T> {
  data: T[];
  pagination: {
    current_page: number;
    limit: number;
    total_count: number;
    total_pages: number;
  };
}