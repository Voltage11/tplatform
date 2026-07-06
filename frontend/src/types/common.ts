export interface PaginatedResponse<T> {
  data: T[];
  pagination: {
    current_page: number;
    limit: number;
    total_count: number;
    total_pages: number;
  };
}