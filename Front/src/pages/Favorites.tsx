import { useEffect, useState } from 'react'
import { booksApi } from '../api/books'
import { BookCard } from '../components/BookCard'
import type { Book } from '../types'

export function Favorites() {
  const [books, setBooks] = useState<Book[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const load = () => {
    setLoading(true)
    booksApi.getFavorites()
      .then(setBooks)
      .catch(() => setError('Не удалось загрузить избранное'))
      .finally(() => setLoading(false))
  }

  useEffect(() => { load() }, [])

  return (
    <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <h1 className="text-2xl font-bold text-gray-900 font-serif mb-6">Избранное</h1>

      {loading && (
        <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 xl:grid-cols-6 gap-4">
          {Array.from({ length: 6 }).map((_, i) => (
            <div key={i} className="card animate-pulse">
              <div className="bg-gray-200 h-40" />
              <div className="p-4 space-y-2">
                <div className="bg-gray-200 h-3 w-16 rounded" />
                <div className="bg-gray-200 h-4 rounded" />
                <div className="bg-gray-200 h-3 w-24 rounded" />
              </div>
            </div>
          ))}
        </div>
      )}

      {error && (
        <div className="text-center py-16 text-gray-500">
          <p className="text-4xl mb-2">⚠️</p>
          <p>{error}</p>
        </div>
      )}

      {!loading && !error && books.length === 0 && (
        <div className="text-center py-16 text-gray-500">
          <p className="text-5xl mb-4">♡</p>
          <p className="text-lg font-medium mb-2">Список пуст</p>
          <p className="text-sm">Добавляйте книги в избранное прямо с карточки</p>
        </div>
      )}

      {!loading && !error && books.length > 0 && (
        <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 xl:grid-cols-6 gap-4">
          {books.map((book) => (
            <BookCard key={book.id} book={book} onFavoriteAdded={load} />
          ))}
        </div>
      )}
    </main>
  )
}
