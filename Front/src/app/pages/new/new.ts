import { Component, OnInit, ChangeDetectorRef, OnDestroy } from '@angular/core';
import { firstValueFrom } from 'rxjs';
import { CommonModule } from '@angular/common';
import { BooksService, Book } from '../../services/books.service';
import { BookReader } from '../../components/book-reader/book-reader';
import { AuthService } from '../../services/auth';
import { FavoritesService } from '../../services/favorites.service';

@Component({
  selector: 'app-new-books',
  standalone: true,
  imports: [CommonModule, BookReader],
  templateUrl: './new.html',
  styleUrls: ['./new.scss'],
})
export class NewBooks implements OnInit, OnDestroy {
  books: Book[] = [];
  loading = false;
  error: string | null = null;

  // состояние ридера (такое же, как в Library)
  readerUrl: string | null = null;
  readerTitle = '';
  readerText: string | null = null;

  constructor(
    private booksService: BooksService,
    private cd: ChangeDetectorRef,
    private auth: AuthService,
    private favoritesService: FavoritesService
  ) {}

  private _onBookAdded = () => {
    try {
      console.log('[NewBooks] book-added event received, reloading new books');
      void this.loadNew();
    } catch (e) {
      console.warn('[NewBooks] failed to handle book-added event', e);
    }
  };

  ngOnInit(): void {
    console.log('[NewBooks] ngOnInit');
    // Поиска на фронтенде нет: просто загружаем новые книги
    void this.loadNew();
    // подписываемся на событие добавления книги, чтобы автоматически обновлять список
    try {
      window.addEventListener('book-added', this._onBookAdded as EventListener);
    } catch (e) {
      console.warn('[NewBooks] failed to attach book-added listener', e);
    }
  }

  ngOnDestroy(): void {}

  openReader(book: Book): void {
    this.loading = true;
    this.booksService.getBookText(book).subscribe(
      (txt) => {
        this.loading = false;
        if (txt && txt.length) {
          this.readerUrl = null;
          this.readerTitle = book.title;
          this.readerText = txt;
        } else if (book.readUrl) {
          this.readerUrl = book.readUrl || null;
          this.readerTitle = book.title;
          this.readerText = null;
        } else {
          this.readerUrl = null;
          this.readerTitle = book.title;
          this.readerText = 'Текст данной книги недоступен для чтения онлайн.';
        }
        try {
          this.cd.detectChanges();
        } catch (e) {}
      },
      (err) => {
        this.loading = false;
        if (book.readUrl) {
          this.readerUrl = book.readUrl;
          this.readerTitle = book.title;
          this.readerText = null;
        } else {
          this.readerUrl = null;
          this.readerTitle = book.title;
          this.readerText = 'Текст данной книги недоступен для чтения онлайн.';
        }
        try {
          this.cd.detectChanges();
        } catch (e) {}
      }
    );
  }

  closeReader(): void {
    this.readerUrl = null;
    this.readerTitle = '';
    this.readerText = null;
  }

  addToFavorites(book: Book) {
    if (!this.auth.isLoggedIn()) {
      alert('Добавлять в избранное могут только авторизованные пользователи. Пожалуйста, войдите.');
      window.dispatchEvent(new CustomEvent('open-auth', { detail: { mode: 'login' } }));
      return;
    }
    const user = (this.auth as any).currentUserSubject?.value || null;
    const userId = user?.id;
    if (!userId) {
      alert('Не удалось определить пользователя.');
      return;
    }
    const currentlyFav = this.isFavorite(book);
    try {
      const key = `favorites:${userId}`;
      const raw = localStorage.getItem(key);
      let arr: string[] = raw ? JSON.parse(raw) : [];

      if (!currentlyFav) {
        // сначала добавить локально для мгновенной обратной связи в UI
        if (!arr.includes(String(book.id))) {
          arr.push(String(book.id));
          localStorage.setItem(key, JSON.stringify(arr));
        }
        try {
          this.cd.detectChanges();
        } catch (e) {}

        // затем синхронизировать с сервером
        this.favoritesService
          .addFavorite(userId, book.id)
          .then(() => {
            try {
              this.cd.detectChanges();
            } catch (e) {}
          })
          .catch((err) => {
            console.error('addFavorite failed; keeping local state', err);
          });
        alert('Книга добавлена в избранное');
      } else {
        // сначала удалить локально
        arr = arr.filter((id) => id !== String(book.id));
        localStorage.setItem(key, JSON.stringify(arr));
        try {
          this.cd.detectChanges();
        } catch (e) {}

        this.favoritesService
          .deleteFavorite(userId, book.id)
          .then(() => {
            try {
              this.cd.detectChanges();
            } catch (e) {}
          })
          .catch((err) => {
            console.error('deleteFavorite failed; keeping local state', err);
          });
        alert('Книга удалена из избранного');
      }
    } catch (e) {
      console.error('favorites error', e);
      alert('Не удалось обновить избранное');
    }
  }

  isFavorite(book: Book): boolean {
    try {
      const user = (this.auth as any).currentUserSubject?.value || null;
      const userId = user?.id;
      if (!userId) return false;
      const key = `favorites:${userId}`;
      const raw = localStorage.getItem(key);
      const arr: string[] = raw ? JSON.parse(raw) : [];
      return arr.includes(String(book.id));
    } catch (e) {
      return false;
    }
  }

  async loadNew() {
    this.loading = true;
    this.error = null;
    try {
      // Предпочитаем сначала получить все книги с бэкенда, затем фильтровать на клиенте.
      // Это позволяет не полагаться на поиск Gutendex и гарантирует показ книг, сохранённых на сервере.
      console.log('[NewBooks] loading all books via fetchAllBooks()');
      const all = await firstValueFrom(this.booksService.fetchAllBooks());
      let items: Book[] = all && Array.isArray(all) ? all : [];
      if (!items || items.length === 0) {
        // запасной вариант: использовать searchGutendex (в случае ошибки он вернёт пустой набор)
        console.warn('[NewBooks] fetchAllBooks returned empty, falling back to searchGutendex');
        const resp: any = await firstValueFrom(
          this.booksService.searchGutendex('', 1, undefined, 100)
        );
        items = (resp && resp.items) || [];
        if (resp && resp.error) {
          console.warn('[NewBooks] books service returned empty/fallback result:', resp.error);
          this.error = 'Сервер /books недоступен — данные недоступны.';
        }
      }

      const thirtyDaysAgo = new Date();
      thirtyDaysAgo.setDate(thirtyDaysAgo.getDate() - 2);

      const hasCreatedAt = items.some((b) => !!(b as any).createdAt);
      const hasIsNewFlag = items.some((b) => (b as any).isNew);
      console.log(
        '[NewBooks] items:',
        items.length,
        'hasCreatedAt=',
        hasCreatedAt,
        'hasIsNew=',
        hasIsNewFlag
      );
      try {
        console.log(
          '[NewBooks] sample items ->',
          (items || []).slice(0, 30).map((b) => ({
            id: b.id,
            title: b.title,
            createdAt: (b as any).createdAt,
            isNew: (b as any).isNew,
          }))
        );
      } catch (e) {
        console.warn('[NewBooks] failed to log sample items', e);
      }

      // Предпочитать `createdAt` (сведётся к недавним в рамках порога), если доступно — таково прежнее поведение
      if (hasCreatedAt) {
        this.books = items.filter((b: Book) => {
          const created = (b as any).createdAt ? new Date((b as any).createdAt) : null;
          return !!created && created >= thirtyDaysAgo;
        });
      } else if (hasIsNewFlag) {
        // Запасной вариант: если сервер использует флаги `isNew` вместо `createdAt`, применять их
        this.books = items.filter((b) => !!(b as any).isNew);
      } else {
        // запасной вариант: сортировать по числовому id (сначала новые), когда возможно; иначе сохранить порядок сервера
        try {
          const byIdDesc = items.slice().sort((a, b) => {
            const ai = Number((a as any).id) || 0;
            const bi = Number((b as any).id) || 0;
            return bi - ai;
          });
          this.books = byIdDesc.slice(0, 20);
        } catch (e) {
          this.books = items.slice(0, 20);
        }
      }
    } catch (e: any) {
      console.error('[NewBooks] load error', e);
      this.error = 'Не удалось загрузить новые книги: ' + (e && e.message ? e.message : String(e));
    } finally {
      this.loading = false;
      Promise.resolve().then(() => {
        try {
          this.cd.detectChanges();
        } catch (err) {}
      });
    }
  }
}
