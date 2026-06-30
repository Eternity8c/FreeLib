import { api } from './client'
import type { Book } from '../types'

export interface GetBooksParams {
  limit?: number
  offset?: number
  genre?: string
}

export const booksApi = {
  getAll: ({ limit = 20, offset = 0, genre }: GetBooksParams = {}) => {
    let path = `/books?limit=${limit}&offset=${offset}`
    if (genre) path += `&genre=${encodeURIComponent(genre)}`
    return api.get<Book[]>(path)
  },

  getNew: ({ limit = 20, offset = 0 }: { limit?: number; offset?: number } = {}) =>
    api.get<Book[]>(`/books/new?limit=${limit}&offset=${offset}`),

  getById: (id: number) => api.get<Book>(`/book?id=${id}`),

  getFavorites: () => api.get<Book[]>('/books/favorite'),

  addFavorite: (bookId: number) => api.post('/book', { book_id: bookId }),

  create: (form: FormData) => api.postForm<Book>('/books', form),

  update: (form: FormData) => api.putForm<Book>('/book', form),

  delete: (id: number) => api.delete(`/book?id=${id}`),

  downloadFile: async (id: number, title: string) => {
    const token = localStorage.getItem('token')
    const headers: Record<string, string> = {}
    if (token) headers['Authorization'] = `Bearer ${token}`

    const res = await fetch(`/api/book/file?id=${id}`, { headers })
    if (!res.ok) throw new Error('Не удалось скачать файл')
    const blob = await res.blob()
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `${title}.epub`
    a.click()
    URL.revokeObjectURL(url)
  },
}
