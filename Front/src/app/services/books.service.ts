import { Injectable } from '@angular/core';
import { Observable, of, throwError, firstValueFrom } from 'rxjs';
import { map, catchError, timeout } from 'rxjs/operators';
import { HttpClient, HttpParams } from '@angular/common/http';
import { environment } from '../../environments/environment';

export interface Book {
  id: string | number;
  title: string;
  localizedTitle?: string;
  author?: string;
  authors?: string[];
  description?: string;
  genre?: string;
  content?: string;
  _sampleText?: string;
  coverUrl?: string;
  coverURL?: string;
  cover_url?: string;
  _coverFailed?: boolean;
  _justSaved?: boolean;
  readUrl?: string;
  createdAt?: string;
  hasText?: boolean;
  textSource?: 'direct' | 'jina' | 'iframe' | 'unknown';
  source?: string;
}

@Injectable({ providedIn: 'root' })
export class BooksService {
  // Константы сети: единый таймаут и базовый URL
  private readonly NETWORK_TIMEOUT_MS = 10000;
  // нормализовать объекты книг, полученные с бэкенда, в формат фронтенда
  private normalizeBook(raw: any): Book {
    if (!raw) return raw;
    const b: any = { ...(raw as any) };
    // бэкенд иногда использует `coverURL` (с заглавными буквами) — привести к `coverUrl`
    if (b.coverURL && !b.coverUrl) {
      b.coverUrl = b.coverURL;
    }
    // бэкенд также может использовать snake_case `cover_url`
    if (b.cover_url && !b.coverUrl) {
      b.coverUrl = b.cover_url;
    }
    // бэкенд иногда присылает строку `author` — преобразовать в массив `authors`
    if (b.author && !b.authors) {
      try {
        b.authors = String(b.author)
          .split(',')
          .map((s: string) => s.trim())
          .filter(Boolean);
      } catch (e) {
        b.authors = [String(b.author)];
      }
    }
    // перенести содержимое `content` из бэкенда в `_sampleText` на фронтенде, чтобы форма администратора видела текст
    if (b.content && !(b as any)._sampleText) {
      try {
        (b as any)._sampleText = b.content;
      } catch (e) {
        (b as any)._sampleText = String(b.content || '');
      }
    }
    // преобразовать snake_case `created_at` в `createdAt`, используемый на фронтенде
    if (b.created_at && !b.createdAt) {
      b.createdAt = b.created_at;
      // считать нулевое время Go (0001-01-01) отсутствующим, чтобы работала логика запасного варианта на фронтенде
      try {
        if (typeof b.createdAt === 'string' && String(b.createdAt).startsWith('0001')) {
          delete b.createdAt;
        }
      } catch (e) {}
    }
    // преобразовать `is_new` в `isNew`
    if ((b as any).is_new !== undefined && (b as any).isNew === undefined) {
      (b as any).isNew = (b as any).is_new;
    }
    return b as Book;
  }

  /**
   * Загрузить файл обложки на бэкенд. Ожидается multipart/form-data с полем 'file'.
   * Сервер должен вернуть либо строку (путь), либо объект { path: string }.
   */
  uploadCover(file: File) {
    try {
      const fd = new FormData();
      fd.append('file', file, file.name);
      return this.http.post<any>(`${this.baseUrl}/upload`, fd).pipe(
        map((resp: any) => {
          if (!resp) return '';
          if (typeof resp === 'string') return resp;
          if (resp.path) return resp.path;
          if (resp.url) return resp.url;
          return '';
        }),
        catchError((err) => {
          console.error('[BooksService] uploadCover failed', err);
          return of('');
        })
      );
    } catch (e) {
      return of('');
    }
  }

  private normalizeBooks(arr: any[] | undefined | null): Book[] {
    if (!arr || !Array.isArray(arr)) return [];
    return arr.map((x) => this.normalizeBook(x));
  }

  // Убраны локальные запасные книги (sampleBooks). В случае отсутствия бэкенда
  // сервис теперь возвращает пустые наборы или пробрасывает ошибки, чтобы
  // вызывающий код явно обрабатывал отсутствие данных.

  private predefinedGenres: string[] = [
    'Роман',
    'Поэзия',
    'Драма',
    'Классика',
    'Фантастика',
    'Фэнтези',
    'Детектив',
    'Приключения',
    'Антиутопия',
    'Публицистика',
    'Эссе',
    'Исторический',
    'Биография',
    'Научная литература',
    'Детская',
    'Юмор',
  ];

  constructor(private http: HttpClient) {}

  private baseUrl = `${environment.apiUrl}/api`;

  searchGutendex(
    q: string,
    page = 1,
    languages?: string,
    pageSize = 20
  ): Observable<{
    total: number;
    items: Book[];
    fallback?: boolean;
    error?: string;
    headerMissing?: boolean;
  }> {
    const query = q && q.trim() ? q.trim() : '';

    // Попытаться получить данные с json-server по пути /books с постраничной навигацией
    try {
      let params = new HttpParams().set('_page', String(page)).set('_limit', String(pageSize));
      if (query) params = params.set('q', query);
      // json-server поддерживает полнотекстовый поиск через параметр q
      return this.http
        .get<Book[]>(`${this.baseUrl}/books`, { params: params, observe: 'response' as any })
        .pipe(
          map((resp: any) => {
            const items: Book[] = this.normalizeBooks(resp.body || []);
            const rawCount = resp.headers ? resp.headers.get('X-Total-Count') : null;
            const headerMissing = rawCount === null || rawCount === undefined;
            const total = headerMissing ? items.length : Number(rawCount);
            return { total, items, headerMissing };
          }),
          catchError((err) => {
            // Логируем ошибку и возвращаем пустой результат — вызывающий код
            // должен корректно обработать отсутствие данных (показать сообщение
            // об ошибке или внутренний запас).
            console.error('[BooksService] failed to fetch /books', err);
            return of({
              total: 0,
              items: [],
              fallback: false,
              error: (err && err.message) || String(err),
            });
          })
        );
    } catch (e) {
      // Если попытка инициировать запрос упала (например ошибка в формировании параметров),
      // вернуть пустой результат и зафиксировать ошибку в консоли.
      console.error('[BooksService] searchGutendex failed before HTTP call', e);
      return of({ total: 0, items: [] });
    }
  }

  getBookText(book: Book): Observable<string> {
    if (book._sampleText) return of(book._sampleText);
    return throwError(() => new Error('no-text-available'));
  }

  findReadableUrl(book: Book): Observable<string | undefined> {
    return of(book.readUrl);
  }

  async addBook(book: Partial<Book>): Promise<Book> {
    try {
      // Собрать полезную нагрузку для сервера. НЕ отправлять клиентский id: сервер
      // использует числовой авто-генерируемый id (uint). Отправка строкового id вызовет
      // `json: cannot unmarshal string into Go struct field Book.id of type uint`.
      const payloadForServer: any = {
        ...book,
        createdAt: (book as any).createdAt || new Date().toISOString(),
      };

      // нормализовать ключи, чтобы соответствовать ожиданиям бэкенда
      // бэкенд может использовать `cover_url` (snake_case) или `coverURL` — привести `coverUrl` к нужному виду
      if (
        (payloadForServer as any).coverUrl &&
        !(payloadForServer as any).cover_url &&
        !(payloadForServer as any).coverURL
      ) {
        // предпочитаем snake_case, который часто используется в Go-бэкенде
        payloadForServer.cover_url = (payloadForServer as any).coverUrl;
        delete (payloadForServer as any).coverUrl;
      }
      if ((payloadForServer as any).authors && !(payloadForServer as any).author) {
        // объединить массив authors в строку author (бэкенд ожидает поле author)
        try {
          payloadForServer.author = (payloadForServer as any).authors.join(', ');
        } catch (e) {
          payloadForServer.author = String((payloadForServer as any).authors || '');
        }
        // сохранить поле authors для фронтенда, но бэкенд будет использовать `author`
      }

      // перенести фронтендовый `_sampleText` или `content` в поле `content` для бэкенда
      if ((payloadForServer as any)._sampleText && !(payloadForServer as any).content) {
        payloadForServer.content = (payloadForServer as any)._sampleText;
        // не отправлять внутреннее поле `_sampleText` на сервер
        delete (payloadForServer as any)._sampleText;
      }

      // POST на endpoint создания на сервере (baseUrl уже содержит /api)
      // небольшой таймаут, чтобы избежать зависания в случае отсутствия ответа от сервера
      console.log(
        '[BooksService] POST /create payload:',
        payloadForServer && {
          sizeHint: String((payloadForServer.cover_url || payloadForServer.coverURL || '').length),
          hasCover: !!(payloadForServer.cover_url || payloadForServer.coverURL),
        }
      );
      const raw = await firstValueFrom(
        this.http
          .post<any>(`${this.baseUrl}/create`, payloadForServer)
          .pipe(timeout(this.NETWORK_TIMEOUT_MS))
      );

      // Если сервер вернул полный объект книги, нормализовать и вернуть его.
      let result: any = raw;
      const looksLikeBook =
        raw &&
        (raw.title || raw.authors || raw.genre || raw.coverUrl || raw.coverURL || raw.cover_url);
      if (!looksLikeBook && raw && raw.id) {
        try {
          const fetched = await firstValueFrom(
            this.http
              .get<Book>(`${this.baseUrl}/books/${encodeURIComponent(String(raw.id))}`)
              .pipe(timeout(this.NETWORK_TIMEOUT_MS))
          );
          result = fetched;
        } catch (e) {
          // не удалось получить полный объект — возвращаем сырое подтверждение
          result = raw;
        }
      }
      return this.normalizeBook(result as any);
    } catch (err) {
      console.error('[BooksService] addBook failed', err);
      // Пробрасываем ошибку вверх — UI/вызывающий код должен обработать её
      throw err;
    }
  }
  async updateBook(book: Partial<Book> & { id: string }): Promise<Book> {
    if (!book || !book.id) throw new Error('missing-id');
    try {
      // собрать payload для сервера: удалить id из тела и нормализовать имена полей
      const payload: any = { ...(book as any) };
      delete payload.id;
      // нормализовать coverUrl -> coverURL (иногда бэкенд ожидает заглавную форму)
      // или использовать snake_case `cover_url`, если это принято на сервере
      if (payload.coverUrl && !payload.cover_url && !payload.coverURL) {
        payload.cover_url = payload.coverUrl;
        delete payload.coverUrl;
      } else if (payload.coverUrl && !payload.coverURL) {
        payload.coverURL = payload.coverUrl;
        delete payload.coverUrl;
      }
      // при необходимости преобразовать authors[] -> author (строка)
      if (payload.authors && !payload.author) {
        try {
          payload.author = (payload.authors || []).join(', ');
        } catch (e) {
          payload.author = String(payload.authors || '');
        }
      }

      // перенести фронтендовый `_sampleText` или `content` в поле `content` для бэкенда
      if ((payload as any)._sampleText && !(payload as any).content) {
        (payload as any).content = (payload as any)._sampleText;
        delete (payload as any)._sampleText;
      }

      // ВАЖНО: маршрут сервера использует единственное /book/{id}
      // отправляем PATCH; сервер может вернуть минимальный ответ (например {message, id})
      // отправляем PATCH с таймаутом, чтобы избежать бесконечного ожидания
      console.log(
        '[BooksService] PATCH /book/' + String(book.id) + ' payload:',
        payload && {
          sizeHint: String((payload.cover_url || payload.coverURL || '').length),
          hasCover: !!(payload.cover_url || payload.coverURL),
        }
      );
      const raw = await firstValueFrom(
        this.http
          .patch<any>(`${this.baseUrl}/book/${encodeURIComponent(String(book.id))}`, payload)
          .pipe(timeout(this.NETWORK_TIMEOUT_MS))
      );

      // Если сервер вернул полный объект книги, нормализовать и вернуть его.
      // Если сервер вернул только подтверждение ({message, id}), попытаться получить книгу GET-запросом.
      let toNormalize: any = raw;
      const looksLikeBook =
        raw &&
        (raw.title || raw.authors || raw.genre || raw.coverUrl || raw.coverURL || raw.cover_url);
      if (!looksLikeBook && raw && raw.id) {
        try {
          const fetched = await firstValueFrom(
            this.http
              .get<Book>(`${this.baseUrl}/books/${encodeURIComponent(String(raw.id))}`)
              .pipe(timeout(this.NETWORK_TIMEOUT_MS))
          );
          toNormalize = fetched;
        } catch (e) {
          toNormalize = raw;
        }
      }

      return this.normalizeBook(toNormalize as any);
    } catch (err) {
      console.error('[BooksService] updateBook failed', err);
      // Пробрасываем ошибку вверх — вызывающий код должен обработать откат/сообщение
      throw err;
    }
  }

  fetchAllBooks(): Observable<Book[] | null> {
    return this.http.get<Book[]>(`${this.baseUrl}/books`).pipe(
      map((items) => this.normalizeBooks(items || [])),
      catchError((err) => {
        console.error('[BooksService] fetchAllBooks failed', err);
        return of(null);
      })
    );
  }

  listGenres(): Observable<string[]> {
    return this.fetchAllBooks().pipe(
      map((items) => {
        const source = items && items.length ? items : [];
        const set = new Set<string>();
        source.forEach((b) => {
          if (b.genre && b.genre.toString().trim()) {
            try {
              const parts = String(b.genre)
                .split(',')
                .map((s) => s.trim())
                .filter(Boolean);
              parts.forEach((p) => set.add(p));
            } catch (e) {
              set.add(b.genre.toString().trim());
            }
          }
        });
        this.predefinedGenres.forEach((g) => set.add(g));
        return Array.from(set).sort();
      })
    );
  }
}
