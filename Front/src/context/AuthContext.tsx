import { createContext, useContext, useState, useCallback, type ReactNode } from 'react'
import { usersApi } from '../api/users'

interface AuthContextValue {
  token: string | null
  isAdmin: boolean
  login: (email: string, password: string) => Promise<void>
  register: (email: string, password: string) => Promise<void>
  logout: () => void
}

const AuthContext = createContext<AuthContextValue | null>(null)

function parseAdmin(token: string): boolean {
  try {
    const payload = JSON.parse(atob(token.split('.')[1]))
    return payload.is_admin === true || payload.isAdmin === true
  } catch {
    return false
  }
}

export function AuthProvider({ children }: { children: ReactNode }) {
  const [token, setToken] = useState<string | null>(() => localStorage.getItem('token'))
  const [isAdmin, setIsAdmin] = useState<boolean>(() => {
    const t = localStorage.getItem('token')
    return t ? parseAdmin(t) : false
  })

  const login = useCallback(async (email: string, password: string) => {
    const res = await usersApi.login(email, password)
    localStorage.setItem('token', res.jwt_token)
    setToken(res.jwt_token)
    setIsAdmin(parseAdmin(res.jwt_token))
  }, [])

  const register = useCallback(async (email: string, password: string) => {
    await usersApi.register(email, password)
  }, [])

  const logout = useCallback(() => {
    localStorage.removeItem('token')
    setToken(null)
    setIsAdmin(false)
  }, [])

  return (
    <AuthContext.Provider value={{ token, isAdmin, login, register, logout }}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const ctx = useContext(AuthContext)
  if (!ctx) throw new Error('useAuth must be used within AuthProvider')
  return ctx
}
