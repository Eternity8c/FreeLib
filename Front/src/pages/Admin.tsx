import { useEffect, useState, type FormEvent } from 'react'
import { booksApi } from '../api/books'
import type { Book } from '../types'

type Mode = 'list' | 'create' | 'edit'

interface BookForm {
  id?: number
  title: string
  author: string
  genre: string
  file: File | null
}

const emptyForm = (): BookForm => ({ title: '', author: '', genre: '', file: null })

export function Admin() {
  const [books, setBooks] = useState<Book[]>([])
  const [loadingList, setLoadingList] = useState(true)
  const [mode, setMode] = useState<Mode>('list')
  const [form, setForm] = useState<BookForm>(emptyForm())
  const [submitting, setSubmitting] = useState(false)
  const [msg, setMsg] = useState<{ text: string; type: 'success' | 'error' } | null>(null)
  const [deletingId, setDeletingId] = useState<number | null>(null)

  const loadBooks = () => {
    setLoadingList(true)
    booksApi.getAll({ limit: 100 })
      .then(setBooks)
      .finally(() => setLoadingList(false))
  }

  useEffect(() => { loadBooks() }, [])

  const showMsg = (text: string, type: 'success' | 'error') => {
    setMsg({ text, type })
    setTimeout(() => setMsg(null), 3000)
  }

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault()
    if (!form.file && mode === 'create') {
      showMsg('Прикрепите EPUB файл', 'error')
      return
    }
    if (!form.file && mode === 'edit') {
      showMsg('Прикрепите EPUB файл', 'error')
      return
    }

    setSubmitting(true)
    try {
      const fd = new FormData()
      fd.append('title', form.title)
      fd.append('author', form.author)
      fd.append('genre', form.genre)
      if (form.file) fd.append('epub', form.file)

      if (mode === 'create') {
        await booksApi.create(fd)
        showMsg('Книга создана', 'success')
      } else if (mode === 'edit' && form.id != null) {
        fd.append('id', String(form.id))
        await booksApi.update(fd)
        showMsg('Книга обновлена', 'success')
      }

      setForm(emptyForm())
      setMode('list')
      loadBooks()
    } catch (err) {
      showMsg(err instanceof Error ? err.message : 'Ошибка', 'error')
    } finally {
      setSubmitting(false)
    }
  }

  const handleEdit = (book: Book) => {
    setForm({ id: book.id, title: book.title, author: book.author, genre: book.genre, file: null })
    setMode('edit')
  }

  const handleDelete = async (id: number) => {
    if (!confirm('Удалить книгу?')) return
    setDeletingId(id)
    try {
      await booksApi.delete(id)
      showMsg('Книга удалена', 'success')
      loadBooks()
    } catch {
      showMsg('Не удалось удалить', 'error')
    } finally {
      setDeletingId(null)
    }
  }

  return (
    <main className="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold text-gray-900 font-serif">Панель администратора</h1>
        {mode === 'list' ? (
          <button onClick={() => { setForm(emptyForm()); setMode('create') }} className="btn-primary">
            + Добавить книгу
          </button>
        ) : (
          <button onClick={() => { setMode('list'); setForm(emptyForm()) }} className="btn-secondary">
            ← Назад
          </button>
        )}
      </div>

      {msg && (
        <div className={`mb-4 px-4 py-3 rounded-lg text-sm font-medium ${
          msg.type === 'success' ? 'bg-green-50 text-green-700 border border-green-200' : 'bg-red-50 text-red-700 border border-red-200'
        }`}>
          {msg.text}
        </div>
      )}

      {(mode === 'create' || mode === 'edit') && (
        <div className="card p-6 mb-8">
          <h2 className="text-lg font-semibold mb-4">
            {mode === 'create' ? 'Новая книга' : `Редактировать: ${form.title}`}
          </h2>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="grid sm:grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Название</label>
                <input
                  type="text"
                  required
                  value={form.title}
                  onChange={(e) => setForm((f) => ({ ...f, title: e.target.value }))}
                  className="input"
                  placeholder="Название книги"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Автор</label>
                <input
                  type="text"
                  required
                  value={form.author}
                  onChange={(e) => setForm((f) => ({ ...f, author: e.target.value }))}
                  className="input"
                  placeholder="Имя Фамилия"
                />
              </div>
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Жанр</label>
              <input
                type="text"
                required
                value={form.genre}
                onChange={(e) => setForm((f) => ({ ...f, genre: e.target.value }))}
                className="input"
                placeholder="Например: Фантастика"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                EPUB файл {mode === 'edit' && <span className="text-gray-400 font-normal">(обязательно для обновления)</span>}
              </label>
              <input
                type="file"
                accept=".epub"
                required
                onChange={(e) => setForm((f) => ({ ...f, file: e.target.files?.[0] ?? null }))}
                className="block w-full text-sm text-gray-500 file:mr-4 file:py-2 file:px-4 file:rounded-lg file:border-0 file:text-sm file:font-medium file:bg-primary-50 file:text-primary-700 hover:file:bg-primary-100"
              />
            </div>
            <div className="flex gap-3 pt-2">
              <button type="submit" disabled={submitting} className="btn-primary">
                {submitting ? 'Сохраняем…' : mode === 'create' ? 'Создать' : 'Сохранить'}
              </button>
              <button
                type="button"
                onClick={() => { setMode('list'); setForm(emptyForm()) }}
                className="btn-secondary"
              >
                Отмена
              </button>
            </div>
          </form>
        </div>
      )}

      {mode === 'list' && (
        <div className="card overflow-hidden">
          {loadingList ? (
            <div className="p-8 text-center text-gray-400">Загрузка…</div>
          ) : books.length === 0 ? (
            <div className="p-8 text-center text-gray-400">Книг нет. Добавьте первую!</div>
          ) : (
            <table className="w-full text-sm">
              <thead className="bg-gray-50 border-b border-gray-200">
                <tr>
                  <th className="text-left px-4 py-3 font-medium text-gray-600">ID</th>
                  <th className="text-left px-4 py-3 font-medium text-gray-600">Название</th>
                  <th className="text-left px-4 py-3 font-medium text-gray-600 hidden sm:table-cell">Автор</th>
                  <th className="text-left px-4 py-3 font-medium text-gray-600 hidden md:table-cell">Жанр</th>
                  <th className="px-4 py-3" />
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-100">
                {books.map((book) => (
                  <tr key={book.id} className="hover:bg-gray-50">
                    <td className="px-4 py-3 text-gray-400 font-mono">{book.id}</td>
                    <td className="px-4 py-3 font-medium text-gray-900">{book.title}</td>
                    <td className="px-4 py-3 text-gray-500 hidden sm:table-cell">{book.author}</td>
                    <td className="px-4 py-3 hidden md:table-cell">
                      <span className="bg-primary-100 text-primary-700 text-xs font-medium px-2 py-0.5 rounded-full">
                        {book.genre}
                      </span>
                    </td>
                    <td className="px-4 py-3">
                      <div className="flex gap-2 justify-end">
                        <button
                          onClick={() => handleEdit(book)}
                          className="text-xs text-primary-600 hover:text-primary-800 font-medium"
                        >
                          Изменить
                        </button>
                        <button
                          onClick={() => handleDelete(book.id)}
                          disabled={deletingId === book.id}
                          className="text-xs text-red-500 hover:text-red-700 font-medium"
                        >
                          {deletingId === book.id ? '…' : 'Удалить'}
                        </button>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </div>
      )}
    </main>
  )
}
