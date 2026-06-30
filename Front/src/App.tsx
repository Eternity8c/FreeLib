import { BrowserRouter, Routes, Route } from 'react-router-dom'
import { AuthProvider } from './context/AuthContext'
import { Navbar } from './components/Navbar'
import { ProtectedRoute } from './components/ProtectedRoute'
import { Home } from './pages/Home'
import { Library } from './pages/Library'
import { NewBooks } from './pages/NewBooks'
import { Favorites } from './pages/Favorites'
import { Login } from './pages/Login'
import { Register } from './pages/Register'
import { Admin } from './pages/Admin'

export default function App() {
  return (
    <BrowserRouter>
      <AuthProvider>
        <div className="min-h-screen flex flex-col">
          <Navbar />
          <div className="flex-1">
            <Routes>
              <Route path="/" element={<Home />} />
              <Route path="/library" element={<Library />} />
              <Route path="/new" element={<NewBooks />} />
              <Route path="/login" element={<Login />} />
              <Route path="/register" element={<Register />} />
              <Route
                path="/favorites"
                element={
                  <ProtectedRoute>
                    <Favorites />
                  </ProtectedRoute>
                }
              />
              <Route
                path="/admin"
                element={
                  <ProtectedRoute adminOnly>
                    <Admin />
                  </ProtectedRoute>
                }
              />
            </Routes>
          </div>
          <footer className="border-t border-gray-200 bg-white py-6 text-center text-sm text-gray-400">
            © {new Date().getFullYear()} FreeLib — Свободная электронная библиотека
          </footer>
        </div>
      </AuthProvider>
    </BrowserRouter>
  )
}
