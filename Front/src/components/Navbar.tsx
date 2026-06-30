import { Link, NavLink, useNavigate } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'

export function Navbar() {
  const { token, isAdmin, logout } = useAuth()
  const navigate = useNavigate()

  const handleLogout = () => {
    logout()
    navigate('/')
  }

  const linkClass = ({ isActive }: { isActive: boolean }) =>
    `text-sm font-medium transition-colors ${
      isActive ? 'text-primary-600' : 'text-gray-600 hover:text-gray-900'
    }`

  return (
    <header className="sticky top-0 z-50 bg-white border-b border-gray-200 shadow-sm">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex items-center justify-between h-16">
          <Link to="/" className="flex items-center gap-2">
            <span className="text-2xl">📚</span>
            <span className="text-xl font-bold text-primary-700 font-serif">FreeLib</span>
          </Link>

          <nav className="hidden sm:flex items-center gap-6">
            <NavLink to="/" end className={linkClass}>
              Главная
            </NavLink>
            <NavLink to="/library" className={linkClass}>
              Каталог
            </NavLink>
            <NavLink to="/new" className={linkClass}>
              Новинки
            </NavLink>
            {token && (
              <NavLink to="/favorites" className={linkClass}>
                Избранное
              </NavLink>
            )}
            {isAdmin && (
              <NavLink to="/admin" className={linkClass}>
                Администратор
              </NavLink>
            )}
          </nav>

          <div className="flex items-center gap-3">
            {token ? (
              <button onClick={handleLogout} className="btn-secondary text-sm">
                Выйти
              </button>
            ) : (
              <>
                <Link to="/login" className="btn-secondary text-sm">
                  Войти
                </Link>
                <Link to="/register" className="btn-primary text-sm">
                  Регистрация
                </Link>
              </>
            )}
          </div>
        </div>
      </div>
    </header>
  )
}
