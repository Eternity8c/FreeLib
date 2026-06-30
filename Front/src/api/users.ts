import { api } from './client'

export const usersApi = {
  register: (email: string, password: string) =>
    api.post<{ id: number; email: string }>('/register', { email, password }),

  login: (email: string, password: string) =>
    api.post<{ jwt_token: string }>('/login', { email, password }),
}
