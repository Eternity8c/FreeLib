export interface Book {
  id: number
  title: string
  author: string
  genre: string
}

export interface User {
  id: number
  email: string
}

export interface AuthState {
  token: string | null
  user: User | null
  isAdmin: boolean
}
