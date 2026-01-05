import { Component, OnInit, ChangeDetectorRef } from '@angular/core';
import { CommonModule } from '@angular/common';
import { BooksService, Book } from '../../services/books.service';
import { AuthService } from '../../services/auth';
import { FavoritesService } from '../../services/favorites.service';
import { BookReader } from '../../components/book-reader/book-reader';

@Component({
  selector: 'app-favorites',
  standalone: true,
  imports: [CommonModule, BookReader],
  templateUrl: './favorites.html',
  styleUrls: ['./favorites.scss'],
})
export class Favorites implements OnInit {
  books: Book[] = [];
  loading = false;
  message: string | null = null;

  readerUrl: string | null = null;
  readerTitle = '';
  readerText: string | null = null;

  constructor(
    private booksService: BooksService,
    private auth: AuthService,
    private cd: ChangeDetectorRef,
    private favoritesService: FavoritesService
  ) {}

  ngOnInit(): void {
    this.loadFavorites();
  }

  async loadFavorites() {
    this.loading = true;
    this.message = null;
    try {
      const user = (this.auth as any).currentUserSubject?.value || null;
      const userId = user?.id;
      if (!userId) {
        this.message = 'Требуется вход в систему для просмотра избранного.';
        this.books = [];
        return;
      }
      // Предпочитать избранное с сервера, когда оно доступно
      try {
        const favs = await this.favoritesService.getFavorites(userId).catch((e) => {
          console.warn('favorites.getFavorites failed, falling back to localStorage', e);
          return null as unknown as Book[];
        });
        if (favs && Array.isArray(favs) && favs.length > 0) {
          this.books = favs;
          return;
        }
      } catch (e) {
        console.warn('favorites.getFavorites error, fallback to localStorage', e);
      }

      // Запасной вариант: использовать localStorage, если сервер не вернул избранное
      const key = `favorites:${userId}`;
      const raw = localStorage.getItem(key);
      const ids: string[] = raw ? JSON.parse(raw) : [];
      if (!ids || ids.length === 0) {
        this.books = [];
        return;
      }
      // получить все книги и отфильтровать по id (json-server не поддерживает batch-эндпоинт)
      const all = await this.booksService.fetchAllBooks().toPromise();
      if (all) {
        this.books = all.filter((b) => ids.includes(String(b.id)));
      } else {
        // запасной вариант: вызвать searchGutendex и отфильтровать
        const resp: any = await this.booksService.searchGutendex('', 1, undefined, 100).toPromise();
        this.books = (resp && resp.items ? resp.items : []).filter((b: Book) =>
          ids.includes(String(b.id))
        );
      }
    } catch (e) {
      console.error('favorites load error', e);
      this.message = 'Не удалось загрузить избранное.';
    } finally {
      this.loading = false;
      Promise.resolve().then(() => {
        try {
          this.cd.detectChanges();
        } catch (e) {}
      });
    }
  }

  removeFromFavorites(book: Book) {
    const user = (this.auth as any).currentUserSubject?.value || null;
    const userId = user?.id;
    if (!userId) {
      alert('Требуется вход для управления избранным.');
      window.dispatchEvent(new CustomEvent('open-auth', { detail: { mode: 'login' } }));
      return;
    }
    const key = `favorites:${userId}`;
    const raw = localStorage.getItem(key);
    let arr: string[] = raw ? JSON.parse(raw) : [];

    // оптимистичный UI: сначала удалить локально, вызвать сервер, при ошибке откатить изменения
    const prevArr = arr.slice();
    const prevBooks = this.books.slice();

    arr = arr.filter((id) => id !== String(book.id));
    try {
      localStorage.setItem(key, JSON.stringify(arr));
    } catch (e) {
      console.error('localStorage set error', e);
    }
    this.books = this.books.filter((b) => String(b.id) !== String(book.id));
    try {
      this.cd.detectChanges();
    } catch (e) {}

    // вызвать сервер для удаления; если не удалось — откатить локальные изменения
    this.favoritesService
      .deleteFavorite(userId, book.id)
      .then(() => {
        // успешно: дополнительных действий не требуется (UI уже обновлён)
      })
      .catch((err) => {
        console.error('deleteFavorite failed, reverting local change', err);
        // откатить локальное состояние
        try {
          localStorage.setItem(key, JSON.stringify(prevArr));
        } catch (e) {
          console.error('localStorage revert error', e);
        }
        this.books = prevBooks;
        try {
          this.cd.detectChanges();
        } catch (e) {}
        alert('Не удалось удалить книгу из избранного на сервере. Изменения отменены.');
      });
  }

  openReader(book: Book) {
    this.booksService.getBookText(book).subscribe(
      (txt) => {
        if (txt && txt.length) {
          this.readerText = txt;
          this.readerUrl = null;
          this.readerTitle = book.title;
        } else if (book.readUrl) {
          this.readerUrl = book.readUrl;
          this.readerText = null;
        } else {
          this.readerText = 'Текст недоступен.';
          this.readerUrl = null;
        }
        try {
          this.cd.detectChanges();
        } catch (e) {}
      },
      (err) => {
        this.readerText = 'Текст недоступен.';
        this.readerUrl = null;
      }
    );
  }

  closeReader() {
    this.readerText = null;
    this.readerUrl = null;
  }
}
