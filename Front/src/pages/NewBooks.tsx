import { useEffect, useState } from 'react'
import { booksApi } from '../api/books'
import { BookCard } from '../components/BookCard'
import type { Book } from '../types'

export function NewBooks() {
  const [books, setBooks] = useState<Book[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    booksApi.getNew({ limit: 24 })
      .then(setBooks)
      .catch(() => setError('Не удалось загрузить новинки'))
      .finally(() => setLoading(false))
  }, [])

  return (
    <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <h1 className="text-2xl font-bold text-gray-900 font-serif mb-6">Новые поступления</h1>

      {loading && (
        <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 xl:grid-cols-6 gap-4">
          {Array.from({ length: 12 }).map((_, i) => (
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
          <p className="text-4xl mb-2">📚</p>
          <p>Новых книг ещё нет</p>
        </div>
      )}

      {!loading && !error && books.length > 0 && (
        <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 xl:grid-cols-6 gap-4">
          {books.map((book) => (
            <BookCard key={book.id} book={book} />
          ))}
        </div>
      )}
    </main>
  )
}
