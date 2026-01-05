import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { AuthModalComponent } from './components/auth-modal/auth-modal';
import { AuthService } from './services/auth';
import { Library } from './pages/library/library';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [CommonModule, RouterModule, FormsModule, AuthModalComponent],
  templateUrl: './app.html',
  styleUrls: ['./app.scss'],
})
export class App {
  showAuthModal = false;
  currentUser: any = null;
  authMode: 'login' | 'register' = 'login';

  books = [
    {
      title: 'Мастер и Маргарита',
      author: 'Михаил Булгаков',
    },
    {
      title: '1984',
      author: 'Джордж Оруэлл',
    },
    {
      title: 'Преступление и наказание',
      author: 'Фёдор Достоевский',
    },
    {
      title: 'Гарри Поттер',
      author: 'Дж. К. Роулинг',
    },
  ];

  constructor(private authService: AuthService) {
    this.authService.currentUser$.subscribe((user) => {
      this.currentUser = user;
    });
    // глобальный слушатель для открытия модального окна аутентификации из других компонентов
    window.addEventListener('open-auth', (e: any) => {
      const mode = e && e.detail && e.detail.mode ? e.detail.mode : 'login';
      this.openAuthModal(mode);
    });
  }

  openAuthModal(mode: 'login' | 'register'): void {
    this.authMode = mode;
    this.showAuthModal = true;
  }

  closeAuthModal(): void {
    this.showAuthModal = false;
  }

  logout(): void {
    this.authService.logout();
  }

  isOnHomePage(): boolean {
    return window.location.pathname === '/';
  }
}
