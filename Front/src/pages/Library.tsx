import { useEffect, useState, useCallback } from 'react'
import { booksApi } from '../api/books'
import { BookCard } from '../components/BookCard'
import type { Book } from '../types'

const GENRES = [
  'Роман', 'Поэзия', 'Драма', 'Классика', 'Фантастика', 'Фэнтези',
  'Детектив', 'Приключения', 'Антиутопия', 'Биография', 'Исторический',
  'Fiction', 'Non-Fiction',
]

const PAGE_SIZE = 12

export function Library() {
  const [books, setBooks] = useState<Book[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [genre, setGenre] = useState('')
  const [offset, setOffset] = useState(0)
  const [hasMore, setHasMore] = useState(true)

  const load = useCallback(async (genreVal: string, offsetVal: number) => {
    setLoading(true)
    setError(null)
    try {
      const data = await booksApi.getAll({ limit: PAGE_SIZE, offset: offsetVal, genre: genreVal || undefined })
      if (offsetVal === 0) {
        setBooks(data)
      } else {
        setBooks((prev) => [...prev, ...data])
      }
      setHasMore(data.length === PAGE_SIZE)
    } catch {
      setError('Не удалось загрузить книги')
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    setOffset(0)
    setBooks([])
    load(genre, 0)
  }, [genre, load])

  const loadMore = () => {
    const next = offset + PAGE_SIZE
    setOffset(next)
    load(genre, next)
  }

  return (
    <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold text-gray-900 font-serif">Каталог книг</h1>
        {books.length > 0 && (
          <span className="text-sm text-gray-500">Найдено: {books.length}+</span>
        )}
      </div>

      <div className="mb-6 flex flex-wrap gap-2">
        <button
          onClick={() => setGenre('')}
          className={`px-3 py-1.5 rounded-full text-sm font-medium transition-colors ${
            genre === '' ? 'bg-primary-600 text-white' : 'bg-white border border-gray-200 text-gray-600 hover:border-primary-400'
          }`}
        >
          Все
        </button>
        {GENRES.map((g) => (
          <button
            key={g}
            onClick={() => setGenre(g)}
            className={`px-3 py-1.5 rounded-full text-sm font-medium transition-colors ${
              genre === g ? 'bg-primary-600 text-white' : 'bg-white border border-gray-200 text-gray-600 hover:border-primary-400'
            }`}
          >
            {g}
          </button>
        ))}
      </div>

      {error && (
        <div className="text-center py-16 text-gray-500">
          <p className="text-4xl mb-2">⚠️</p>
          <p>{error}</p>
        </div>
      )}

      {!error && books.length === 0 && !loading && (
        <div className="text-center py-16 text-gray-500">
          <p className="text-4xl mb-2">📭</p>
          <p>Книг по этому жанру не найдено</p>
        </div>
      )}

      <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 xl:grid-cols-6 gap-4">
        {books.map((book) => (
          <BookCard key={book.id} book={book} />
        ))}
        {loading &&
          Array.from({ length: PAGE_SIZE }).map((_, i) => (
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

      {!loading && hasMore && books.length > 0 && (
        <div className="mt-8 flex justify-center">
          <button onClick={loadMore} className="btn-secondary">
            Загрузить ещё
          </button>
        </div>
      )}
    </main>
  )
}
