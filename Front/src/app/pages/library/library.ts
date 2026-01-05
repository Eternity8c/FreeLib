import { Component, OnInit, ChangeDetectorRef, OnDestroy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { BooksService, Book } from '../../services/books.service';
import { BookReader } from '../../components/book-reader/book-reader';
import { AuthService } from '../../services/auth';
import { FavoritesService } from '../../services/favorites.service';

@Component({
  selector: 'app-library',
  standalone: true,
  imports: [CommonModule, RouterModule, FormsModule, BookReader],
  templateUrl: './library.html',
  styleUrls: ['./library.scss'],
})
export class Library implements OnInit {
  books: Book[] = [];
  unfiltered: Book[] = [];
  total = 0;
  page = 1;
  pageSize = 3;
  defaultPageSize = 3;
  private expandedToTotal = false;
  totalPages = 1;
  readerUrl: string | null = null;
  readerTitle = '';
  readerText: string | null = null;

  loading = false;
  error: string | null = null;
  statusMessage: string | null = null;

  genres: string[] = [];
  selectedGenres: Set<string> = new Set();
  showOnlyWithCover = false;
  showOnlyWithText = false;

  constructor(
    private booksService: BooksService,
    private cd: ChangeDetectorRef,
    private auth: AuthService,
    private favoritesService: FavoritesService
  ) {}

  ngOnInit(): void {
    // Поиска на фронтенде нет: просто загружаем жанры и книги
    this.booksService.listGenres().subscribe((g) => {
      this.genres = g || [];
      this.load();
    });
  }

  toggleCover(event: Event) {
    const inp = event.target as HTMLInputElement;
    this.showOnlyWithCover = !!inp.checked;
    this.applyAllFilters();
  }

  toggleTextOnly(event: Event) {
    const inp = event.target as HTMLInputElement;
    this.showOnlyWithText = !!inp.checked;
    this.applyAllFilters();
  }

  private applyAllFilters() {
    try {
      let arr = (this.unfiltered || []).slice();

      if (this.selectedGenres && this.selectedGenres.size > 0) {
        arr = arr.filter((b) => {
          if (!b.genre) return false;
          const bookGenres = b.genre
            .toString()
            .split(',')
            .map((s) => s.trim())
            .filter(Boolean);
          return bookGenres.some((bg) => this.selectedGenres.has(bg));
        });
      }

      if (this.showOnlyWithCover) {
        arr = arr.filter((b) => !!(b.coverUrl && b.coverUrl.toString().trim()));
      }

      if (this.showOnlyWithText) {
        arr = arr.filter((b) => !!(b.readUrl || b._sampleText || b.hasText));
      }
      this.books = arr.slice(0, this.pageSize);
      this.totalPages = Math.max(1, Math.ceil((arr.length || 0) / this.pageSize));
      Promise.resolve().then(() => {
        try {
          this.cd.detectChanges();
        } catch (e) {}
      });
    } catch (e) {}
  }

  ngOnDestroy(): void {}

  load(): void {
    this.loading = true;
    this.error = null;
    this.statusMessage = 'Загрузка...';
    console.log('[Library] load start', { page: this.page });

    // вызываем json-server /books (запасной вариант Gutendex) — поиск на фронтенде не поддерживается
    this.booksService.searchGutendex('', this.page, undefined, this.pageSize).subscribe({
      next: (res: any) => this.handleSearchResponse(res),
      error: (err) => {
        console.error('Gutendex error', err);
        try {
          const key = `gutendex:cache:all:page:${this.page}`;
          const raw = localStorage.getItem(key);
          if (raw) {
            const cached = JSON.parse(raw);
            this.total = cached.total || 0;
            this.unfiltered = cached.items || [];
            this.applyAllFilters();
            this.statusMessage = 'Gutendex недоступен — показаны закешированные результаты.';
            this.loading = false;
            return;
          }
        } catch (e) {}

        // Нет кеша — очистим список и покажем простое сообщение об ошибке.
        this.statusMessage = 'Gutendex недоступен — проверьте соединение.';
        this.loading = false;
        this.unfiltered = [];
        this.applyAllFilters();
        try {
          this.cd.detectChanges();
        } catch (e) {}
      },
    });
  }

  goPrev(): void {
    if (this.page > 1) {
      this.page -= 1;
      this.load();
    }
  }

  goNext(): void {
    if (this.page < this.totalPages) {
      this.page += 1;
      this.load();
    }
  }

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
        } catch (e) {
          // игнорировать
        }
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
        } catch (e) {
          // игнорировать
        }
      }
    );
  }

  closeReader(): void {
    this.readerUrl = null;
    this.readerTitle = '';
    this.readerText = null;
  }

  // вызывается, когда загрузка <img> завершилась ошибкой
  onCoverError(book: Book, ev: Event): void {
    try {
      const img = ev?.target as HTMLImageElement | null;
      const src =
        (img && img.src) ||
        (book as any).coverUrl ||
        (book as any)['coverURL'] ||
        (book as any)['cover_url'];
      console.warn('[Library] cover load failed', { id: book.id, src, event: ev });
      // пометить книгу, чтобы шаблон мог показать заглушку
      (book as any)._coverFailed = true;
      try {
        this.cd.detectChanges();
      } catch (e) {}
    } catch (e) {
      // игнорировать
    }
  }

  toggleGenre(g: string, event: Event) {
    const inp = event.target as HTMLInputElement;
    if (inp.checked) this.selectedGenres.add(g);
    else this.selectedGenres.delete(g);
    // Применять фильтры локально без повторной загрузки с сервера
    this.applyAllFilters();
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

    // оптимистичный UI: сразу обновляем localStorage и интерфейс, затем синхронизируем с сервером
    const currentlyFav = this.isFavorite(book);
    try {
      const key = `favorites:${userId}`;
      const raw = localStorage.getItem(key);
      let arr: string[] = raw ? JSON.parse(raw) : [];

      if (!currentlyFav) {
        // сначала добавить локально
        if (!arr.includes(String(book.id))) {
          arr.push(String(book.id));
          localStorage.setItem(key, JSON.stringify(arr));
        }
        try {
          this.cd.detectChanges();
        } catch (e) {}

        // затем уведомить сервер (fire-and-forget)
        this.favoritesService
          .addFavorite(userId, book.id)
          .then(() => {
            // сервер подтвердил; опционально можно уведомить пользователя
            try {
              this.cd.detectChanges();
            } catch (e) {}
          })
          .catch((err) => {
            console.error('addFavorite failed; keeping local state', err);
            // намеренно сохраняем локальное оптимистичное состояние, чтобы UI оставался отзывчивым
          });
        alert('Книга добавлена в избранное');
      } else {
        // сначала удалить локально
        arr = arr.filter((id) => id !== String(book.id));
        localStorage.setItem(key, JSON.stringify(arr));
        try {
          this.cd.detectChanges();
        } catch (e) {}

        // затем уведомить сервер
        this.favoritesService
          .deleteFavorite(userId, book.id)
          .then(() => {
            try {
              this.cd.detectChanges();
            } catch (e) {}
          })
          .catch((err) => {
            console.error('deleteFavorite failed; keeping local state', err);
            // сохраняем локальное изменение для отзывчивости интерфейса
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
      // предпочитаем быстную проверку в localStorage для синхронизации; авторитетное состояние — на сервере
      const key = `favorites:${userId}`;
      const raw = localStorage.getItem(key);
      const arr: string[] = raw ? JSON.parse(raw) : [];
      return arr.includes(String(book.id));
    } catch (e) {
      return false;
    }
  }

  private handleSearchResponse(res: any) {
    console.log('[Library] /books response', res);
    this.total = res.total;
    this.unfiltered = res.items || [];

    this.applyAllFilters();
    if (res.error) {
      console.warn('[Library] books service returned empty/fallback result:', res.error);
      this.statusMessage = 'Сервер /books недоступен — данные недоступны.';
    }

    if (
      !res.fallback &&
      this.page === 1 &&
      Array.isArray(res.items) &&
      res.items.length < this.pageSize
    ) {
      console.warn(
        '[Library] first page returned fewer items than pageSize — trying fetchAllBooks()'
      );
      this.booksService.fetchAllBooks().subscribe((all) => {
        if (all && all.length > (res.items ? res.items.length : 0)) {
          console.log(
            '[Library] fetchAllBooks succeeded (first-page short), replacing items',
            all.length
          );
          if (all.length > this.pageSize) {
            this.expandedToTotal = true;
            this.pageSize = all.length;
          }
          this.unfiltered = all || [];
          this.applyAllFilters();
          this.total = all.length;
          this.totalPages = Math.max(1, Math.ceil(this.total / this.pageSize));
          this.statusMessage = `Книг найдено: ${this.total}`;
          this.loading = false;
          Promise.resolve().then(() => {
            try {
              this.cd.detectChanges();
            } catch (e) {}
          });
        }
      });
    }

    if (res.headerMissing) {
      console.warn('[Library] X-Total-Count header missing — attempting fetchAllBooks()');
      this.booksService.fetchAllBooks().subscribe((all) => {
        if (all && all.length) {
          console.log(
            '[Library] fetchAllBooks succeeded (headerMissing), replacing items',
            all.length
          );
          if (all.length > this.pageSize) {
            this.expandedToTotal = true;
            this.pageSize = all.length;
          }
          this.unfiltered = all || [];
          this.applyAllFilters();
          this.total = all.length;
          this.totalPages = Math.max(1, Math.ceil(this.total / this.pageSize));
          this.statusMessage = `Книг найдено: ${this.total}`;
          this.loading = false;
          Promise.resolve().then(() => {
            try {
              this.cd.detectChanges();
            } catch (e) {}
          });
        }
      });
    }

    if (!this.expandedToTotal && this.total > this.pageSize) {
      this.expandedToTotal = true;
      this.pageSize = this.total;
      // перезагрузить один раз с увеличенным pageSize
      setTimeout(() => this.load(), 0);
      return;
    }
    this.totalPages = Math.max(1, Math.ceil(this.total / this.pageSize));
    this.loading = false;

    if (!this.books.length) {
      Promise.resolve().then(() => {
        this.statusMessage = 'Ничего не найдено.';
        this.error = null;
      });
    } else {
      Promise.resolve().then(() => {
        this.statusMessage = `Книг найдено: ${this.total}`;
        this.error = null;
      });
    }

    const checkN = Math.min(this.books.length, 8);
    for (let i = 0; i < checkN; i++) {
      const b = this.books[i];
      this.booksService.findReadableUrl(b).subscribe((url) => {
        if (url) {
          b.hasText = true;
          b.readUrl = url;
          b.textSource = 'direct';
        } else {
          b.hasText = !!b._sampleText;
        }
      });
    }

    Promise.resolve().then(() => {
      try {
        this.cd.detectChanges();
      } catch (e) {}
    });
  }
}
