import http from './http'
import type { TokenPair, User } from '@/types/user'

/** Отправить логин/пароль, получить токены */
export async function login(email: string, password: string): Promise<TokenPair> {
  const { data } = await http.post<TokenPair>('/auth/login', { email, password })

  return data
}

/** Выход (отзыв refresh-токена) */
export async function logout(refreshToken: string): Promise<void> {
  await http.post('/auth/logout', { refresh_token: refreshToken })
}

/** Получить профиль текущего пользователя */
export async function getProfile(): Promise<{ user: User }> {
  const { data } = await http.get<{ user: User }>('/profile')
  return data
}
