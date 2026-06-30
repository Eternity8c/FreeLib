import { useState } from 'react'
import type { Book } from '../types'
import { booksApi } from '../api/books'
import { useAuth } from '../context/AuthContext'

const GENRE_COLORS: Record<string, string> = {
  'Роман': 'bg-rose-100 text-rose-700',
  'Поэзия': 'bg-purple-100 text-purple-700',
  'Драма': 'bg-orange-100 text-orange-700',
  'Классика': 'bg-amber-100 text-amber-700',
  'Фантастика': 'bg-blue-100 text-blue-700',
  'Фэнтези': 'bg-indigo-100 text-indigo-700',
  'Детектив': 'bg-gray-100 text-gray-700',
  'Приключения': 'bg-green-100 text-green-700',
  'Антиутопия': 'bg-red-100 text-red-700',
  'Fiction': 'bg-sky-100 text-sky-700',
}

function genreColor(genre: string): string {
  return GENRE_COLORS[genre] ?? 'bg-violet-100 text-violet-700'
}

interface Props {
  book: Book
  onFavoriteAdded?: () => void
}

export function BookCard({ book, onFavoriteAdded }: Props) {
  const { token } = useAuth()
  const [favoriting, setFavoriting] = useState(false)
  const [downloading, setDownloading] = useState(false)
  const [msg, setMsg] = useState<string | null>(null)

  const handleFavorite = async () => {
    if (!token) return
    setFavoriting(true)
    setMsg(null)
    try {
      await booksApi.addFavorite(book.id)
      setMsg('Добавлено!')
      onFavoriteAdded?.()
    } catch {
      setMsg('Уже в избранном')
    } finally {
      setFavoriting(false)
      setTimeout(() => setMsg(null), 2000)
    }
  }

  const handleDownload = async () => {
    setDownloading(true)
    try {
      await booksApi.downloadFile(book.id, book.title)
    } catch {
      alert('Не удалось скачать файл')
    } finally {
      setDownloading(false)
    }
  }

  return (
    <div className="card flex flex-col group hover:shadow-md transition-shadow">
      <div className="bg-gradient-to-br from-primary-50 to-primary-100 h-40 flex items-center justify-center">
        <span className="text-5xl select-none">📖</span>
      </div>

      <div className="p-4 flex flex-col gap-2 flex-1">
        <span className={`text-xs font-medium px-2 py-0.5 rounded-full w-fit ${genreColor(book.genre)}`}>
          {book.genre}
        </span>

        <h3 className="font-semibold text-gray-900 leading-snug line-clamp-2 font-serif">
          {book.title}
        </h3>
        <p className="text-sm text-gray-500">{book.author}</p>

        {msg && (
          <p className="text-xs text-primary-600 font-medium">{msg}</p>
        )}

        <div className="mt-auto pt-3 flex gap-2 flex-wrap">
          <button
            onClick={handleDownload}
            disabled={downloading}
            className="btn-primary text-xs flex-1"
          >
            {downloading ? 'Загрузка…' : 'Скачать EPUB'}
          </button>
          {token && (
            <button
              onClick={handleFavorite}
              disabled={favoriting}
              className="btn-secondary text-xs px-3"
              title="В избранное"
            >
              ♥
            </button>
          )}
        </div>
      </div>
    </div>
  )
}
