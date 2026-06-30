import { useState, type FormEvent } from 'react'
import { useNavigate, Link } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'

export function Register() {
  const { register } = useAuth()
  const navigate = useNavigate()
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault()
    setError(null)
    setLoading(true)
    try {
      await register(email, password)
      navigate('/login')
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Ошибка регистрации')
    } finally {
      setLoading(false)
    }
  }

  return (
    <main className="min-h-[calc(100vh-4rem)] flex items-center justify-center px-4 py-12">
      <div className="w-full max-w-md">
        <div className="text-center mb-8">
          <span className="text-4xl">📚</span>
          <h1 className="mt-3 text-2xl font-bold text-gray-900 font-serif">Создать аккаунт</h1>
          <p className="mt-1 text-gray-500 text-sm">Присоединяйтесь к FreeLib бесплатно</p>
        </div>

        <div className="card p-8">
          <form onSubmit={handleSubmit} className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Email
              </label>
              <input
                type="email"
                required
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                className="input"
                placeholder="user@example.com"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Пароль
              </label>
              <input
                type="password"
                required
                minLength={8}
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                className="input"
                placeholder="Минимум 8 символов"
              />
            </div>

            {error && (
              <div className="bg-red-50 border border-red-200 text-red-700 text-sm px-4 py-3 rounded-lg">
                {error}
              </div>
            )}

            <button type="submit" disabled={loading} className="btn-primary w-full">
              {loading ? 'Регистрируем…' : 'Зарегистрироваться'}
            </button>
          </form>
        </div>

        <p className="text-center mt-4 text-sm text-gray-600">
          Уже есть аккаунт?{' '}
          <Link to="/login" className="text-primary-600 font-medium hover:underline">
            Войти
          </Link>
        </p>
      </div>
    </main>
  )
}
