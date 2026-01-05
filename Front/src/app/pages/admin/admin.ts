import { Component, ChangeDetectorRef } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { BooksService, Book } from '../../services/books.service';
import { AuthService } from '../../services/auth';
import { Router } from '@angular/router';
import { firstValueFrom } from 'rxjs';

@Component({
  selector: 'app-admin-books',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './admin.html',
  styleUrls: ['./admin.scss'],
})
export class AdminBooks {
  title = '';
  authors = '';
  genre = '';
  coverUrl = '';
  hasText = false;
  sampleText = '';

  message: string | null = null;
  isSubmitting = false;
  // список существующих книг для редактирования
  books: Book[] = [];
  editingId: string | null = null;
  availableGenres: string[] = [];
  selectedGenresForForm: string[] = [];
  newGenreInput = '';
  isAdmin = false;
  // конструктор и внедрение зависимостей
  constructor(
    private booksService: BooksService,
    private auth: AuthService,
    private router: Router, // Router для навигации
    private cd: ChangeDetectorRef
  ) {
    this.auth.currentUser$.subscribe((u) => {
      const wasAdmin = this.isAdmin;
      const newIsAdmin = !!(u as any)?.isAdmin;
      // отложенное присвоение, чтобы избежать ExpressionChangedAfterItHasBeenCheckedError
      if (newIsAdmin !== wasAdmin) {
        Promise.resolve().then(() => {
          this.isAdmin = newIsAdmin;
          // Если пользователь стал админом и книги ещё не загружены, получить их
          if (this.isAdmin && (!this.books || this.books.length === 0)) {
            void this.loadBooks();
          }
          try {
            this.cd.detectChanges();
          } catch (e) {}
        });
      }
    });
  }
  // Примечание: загрузка файлов удалена — обложка указывается только URL в `coverUrl`

  addAvailableGenres(): void {
    const raw = (this.newGenreInput || '')
      .split(',')
      .map((s) => s.trim())
      .filter(Boolean);
    if (!raw.length) return;
    const set = new Set(this.availableGenres || []);
    raw.forEach((g) => set.add(g));
    this.availableGenres = Array.from(set).sort();
    // также отметить в форме только что добавленные жанры
    const selSet = new Set(this.selectedGenresForForm || []);
    raw.forEach((g) => selSet.add(g));
    this.selectedGenresForForm = Array.from(selSet);
    this.newGenreInput = '';
  }

  // Обработчик переключателя в шаблоне: добавить/удалить жанр в выборе
  onToggleGenreModel(g: string, checked: boolean): void {
    try {
      const set = new Set(this.selectedGenresForForm || []);
      if (checked) set.add(g);
      else set.delete(g);
      this.selectedGenresForForm = Array.from(set);
    } catch (e) {
      // запасной вариант: восстановить простой массив
      if (checked) {
        if (!this.selectedGenresForForm.includes(g)) this.selectedGenresForForm.push(g);
      } else {
        this.selectedGenresForForm = (this.selectedGenresForForm || []).filter((x) => x !== g);
      }
    }
  }

  async ngOnInit(): Promise<void> {
    // загружать существующие книги для редактирования только если при инициализации пользователь уже админ
    if (this.isAdmin) {
      await this.loadBooks();
    }
  }

  private async loadBooks(): Promise<void> {
    try {
      const all = await firstValueFrom(this.booksService.fetchAllBooks());
      if (all) {
        this.books = all;
        // также заполнить список доступных жанров
        const gs = await firstValueFrom(this.booksService.listGenres());
        this.availableGenres = gs || [];
        // гарантировать обновление представления Angular после асинхронной загрузки
        try {
          this.cd.detectChanges();
        } catch (e) {}
        return;
      }
      // запасной вариант: попробовать searchGutendex (в случае ошибки он вернёт пустой набор)
      const searchResult = await firstValueFrom(
        this.booksService.searchGutendex('', 1, undefined, 100)
      );
      if (searchResult && searchResult.items) this.books = searchResult.items;
      const gs2 = await firstValueFrom(this.booksService.listGenres());
      this.availableGenres = gs2 || [];
      try {
        this.cd.detectChanges();
      } catch (e) {}
    } catch (e) {
      console.warn('Failed to load books for admin edit', e);
    }
  }

  startEdit(b: Book) {
    this.editingId = String(b.id);
    // заполнить форму значениями книги
    this.title = b.title || '';
    // поддерживать и `authors: string[]`, и устаревшее `author: string`
    if (b.authors && b.authors.length) this.authors = (b.authors || []).join(', ');
    else if ((b as any).author) this.authors = String((b as any).author);
    else this.authors = '';
    this.genre = b.genre || '';
    // разбить CSV в массив для мультиселекта
    this.selectedGenresForForm = (b.genre || '')
      .split(',')
      .map((s) => s.trim())
      .filter(Boolean);
    try {
      const set = new Set(this.availableGenres || []);
      this.selectedGenresForForm.forEach((g) => set.add(g));
      this.availableGenres = Array.from(set).sort();
    } catch (e) {
      // игнорировать
    }
    // поле publishYear удалено из формы
    // поддержка вариантов поля обложки: cover_url, coverURL или coverUrl
    this.coverUrl = (b.coverUrl as any) || (b as any).coverURL || (b as any).cover_url || '';
    this.hasText = !!b.hasText;
    this.sampleText = b._sampleText || '';
    this.message = null;
  }

  async saveEdit(): Promise<void> {
    if (this.isSubmitting) return;
    if (!this.editingId) return;
    if (!this.title || !this.authors) {
      this.message = 'Заполните название и авторов';
      return;
    }
    // предотвращаем отправку локальных путей файлов в coverUrl
    if (this.coverUrl && /(^[a-zA-Z]:\\)|(^file:\/\/)/.test(this.coverUrl)) {
      this.message =
        'Ссылка указывает на локальный файл. Пожалуйста, загрузите изображение через форму или используйте http(s)-ссылку.';
      return;
    }
    const manual = (this.genre || '')
      .split(',')
      .map((s) => s.trim())
      .filter(Boolean);
    const mergedGenres = Array.from(new Set([...this.selectedGenresForForm, ...manual]));

    const payload: Partial<Book> & { id: string } = {
      id: this.editingId,
      title: this.title,
      authors: this.authors.split(',').map((s) => s.trim()),
      genre: mergedGenres.join(', '),
      // поле publishYear пропущено в payload
      coverUrl: this.coverUrl,
      hasText: this.hasText,
      _sampleText: this.hasText ? this.sampleText : undefined,
    };
    this.isSubmitting = true;
    try {
      console.log('[AdminBooks] saveEdit() payload:', payload);
      // обложка предоставляется через `coverUrl` (здесь нет загрузки файлов на клиенте)

      const updated = await this.booksService.updateBook(payload);
      console.log('[AdminBooks] saveEdit() server response:', updated);
      this.message = `Книга "${updated?.title || payload.title}" обновлена`;
      // обновить локальный список
      const idx = this.books.findIndex((x) => x.id === updated.id);
      if (idx >= 0) this.books[idx] = updated;
      else this.books.push(updated);
      // пометить обновлённую книгу, чтобы пользователь увидел индикатор сохранения
      try {
        const jidx = this.books.findIndex((x) => x.id === updated.id);
        if (jidx >= 0) {
          (this.books[jidx] as any)._justSaved = true;
          try {
            this.cd.detectChanges();
          } catch (e) {}
          setTimeout(() => {
            try {
              delete (this.books[jidx] as any)._justSaved;
            } catch (e) {}
          }, 3000);
        }
      } catch (e) {
        // игнорировать
      }
      try {
        const mg = (mergedGenres || []).slice();
        const set = new Set(this.availableGenres || []);
        mg.forEach((g) => set.add(g));
        this.availableGenres = Array.from(set).sort();
        this.selectedGenresForForm = mg;
        try {
          this.cd.detectChanges();
        } catch (e) {}
      } catch (e) {
        // игнорировать
      }
      // выйти из режима редактирования и очистить форму
      this.editingId = null;
      this.title = '';
      this.authors = '';
      this.genre = '';
      this.selectedGenresForForm = [];
      // поле publishYear очищено (поле удалено)
      this.coverUrl = '';
      this.hasText = false;
      this.sampleText = '';
    } catch (err) {
      console.error(err);
      this.message = 'Ошибка при обновлении книги';
    } finally {
      this.isSubmitting = false;
      setTimeout(() => (this.message = null), 4000);
    }
  }

  async add(): Promise<void> {
    if (this.isSubmitting) return;
    if (!this.title || !this.authors) {
      this.message = 'Заполните название и авторов';
      return;
    }
    // предотвращаем отправку локальных путей файлов в coverUrl
    if (this.coverUrl && /(^[a-zA-Z]:\\)|(^file:\/\/)/.test(this.coverUrl)) {
      this.message = 'Ссылка указывает на локальный файл. Пожалуйста, используйте http(s)-ссылку.';
      return;
    }

    const manual = (this.genre || '')
      .split(',')
      .map((s) => s.trim())
      .filter(Boolean);
    const mergedGenres = Array.from(new Set([...this.selectedGenresForForm, ...manual]));

    const payload: Partial<Book> = {
      title: this.title,
      authors: this.authors.split(',').map((s) => s.trim()),
      genre: mergedGenres.join(', '),
      coverUrl: this.coverUrl,
      hasText: this.hasText,
      _sampleText: this.hasText ? this.sampleText : undefined,
      createdAt: new Date().toISOString(),
      // пометить как новое, чтобы страницы «новинки», полагающиеся на isNew, показывали его
      isNew: true,
    } as Partial<Book>;

    this.isSubmitting = true;
    try {
      console.log('[AdminBooks] add() payload:', payload);
      const created = await this.booksService.addBook(payload);
      console.log('[AdminBooks] add() server response:', created);
      this.message = `Книга "${created?.title || payload.title}" добавлена`;
      // добавить в локальный список и пометить как только что сохранённую
      this.books.unshift(created as Book);
      try {
        (this.books[0] as any)._justSaved = true;
        try {
          this.cd.detectChanges();
        } catch (e) {}
        setTimeout(() => {
          try {
            delete (this.books[0] as any)._justSaved;
          } catch (e) {}
        }, 3000);
      } catch (e) {}

      // уведомить другие части приложения, что книга добавлена, чтобы они могли обновиться
      try {
        window.dispatchEvent(new CustomEvent('book-added', { detail: created }));
      } catch (e) {
        // игнорировать, если dispatch не сработает в некоторых окружениях
      }

      // объединить жанры в список доступных
      try {
        const mg = (mergedGenres || []).slice();
        const set = new Set(this.availableGenres || []);
        mg.forEach((g) => set.add(g));
        this.availableGenres = Array.from(set).sort();
        this.selectedGenresForForm = mg;
        try {
          this.cd.detectChanges();
        } catch (e) {}
      } catch (e) {}

      // очистить форму
      this.title = '';
      this.authors = '';
      this.genre = '';
      this.selectedGenresForForm = [];
      this.coverUrl = '';
      this.hasText = false;
      this.sampleText = '';
    } catch (err) {
      console.error(err);
      this.message = 'Ошибка при добавлении книги';
    } finally {
      this.isSubmitting = false;
      setTimeout(() => (this.message = null), 4000);
    }
  }

  cancelEdit() {
    this.editingId = null;
    this.title = '';
    this.authors = '';
    this.genre = '';
    this.selectedGenresForForm = [];
    // поле publishYear удалено
    this.coverUrl = '';
    this.hasText = false;
    this.sampleText = '';
    this.message = null;
  }
  // обработка input файла для загрузки обложки
  // загрузка файлов удалена: обложка принимается только как URL через поле `coverUrl`
}
