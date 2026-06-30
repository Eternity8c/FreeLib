import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { booksApi } from '../api/books'
import { BookCard } from '../components/BookCard'
import type { Book } from '../types'

export function Home() {
  const [newBooks, setNewBooks] = useState<Book[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    booksApi.getNew({ limit: 6 })
      .then(setNewBooks)
      .catch(() => setError('Не удалось загрузить книги'))
      .finally(() => setLoading(false))
  }, [])

  return (
    <main>
      <section className="bg-gradient-to-br from-primary-700 to-primary-900 text-white py-20 px-4">
        <div className="max-w-3xl mx-auto text-center">
          <h1 className="text-4xl sm:text-5xl font-bold font-serif mb-4">
            Читайте свободно
          </h1>
          <p className="text-primary-200 text-lg mb-8">
            Бесплатная электронная библиотека с тысячами книг в формате EPUB
          </p>
          <div className="flex gap-4 justify-center flex-wrap">
            <Link to="/library" className="bg-white text-primary-700 font-semibold px-6 py-3 rounded-lg hover:bg-primary-50 transition-colors">
              Открыть каталог
            </Link>
            <Link to="/new" className="border border-primary-300 text-white font-semibold px-6 py-3 rounded-lg hover:bg-primary-600 transition-colors">
              Новинки
            </Link>
          </div>
        </div>
      </section>

      <section className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        <div className="flex items-center justify-between mb-6">
          <h2 className="text-2xl font-bold text-gray-900 font-serif">Новые поступления</h2>
          <Link to="/new" className="text-primary-600 text-sm font-medium hover:underline">
            Смотреть все →
          </Link>
        </div>

        {loading && (
          <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-6 gap-4">
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
          <div className="text-center py-12 text-gray-500">
            <p className="text-4xl mb-2">⚠️</p>
            <p>{error}</p>
          </div>
        )}

        {!loading && !error && newBooks.length === 0 && (
          <div className="text-center py-12 text-gray-500">
            <p className="text-4xl mb-2">📚</p>
            <p>Книги ещё не добавлены</p>
          </div>
        )}

        {!loading && !error && newBooks.length > 0 && (
          <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-6 gap-4">
            {newBooks.map((book) => (
              <BookCard key={book.id} book={book} />
            ))}
          </div>
        )}
      </section>

      <section className="bg-primary-50 py-12">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="grid sm:grid-cols-3 gap-8 text-center">
            <div>
              <div className="text-4xl mb-2">📖</div>
              <h3 className="text-lg font-semibold mb-1">Форматы EPUB</h3>
              <p className="text-gray-500 text-sm">Скачивайте книги в популярном формате для любого устройства</p>
            </div>
            <div>
              <div className="text-4xl mb-2">🆓</div>
              <h3 className="text-lg font-semibold mb-1">Бесплатно</h3>
              <p className="text-gray-500 text-sm">Все книги доступны без оплаты и регистрации</p>
            </div>
            <div>
              <div className="text-4xl mb-2">❤️</div>
              <h3 className="text-lg font-semibold mb-1">Избранное</h3>
              <p className="text-gray-500 text-sm">Сохраняйте понравившиеся книги в личный список</p>
            </div>
          </div>
        </div>
      </section>
    </main>
  )
}
