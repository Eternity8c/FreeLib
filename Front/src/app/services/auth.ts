import { Injectable } from '@angular/core';
import { BehaviorSubject, firstValueFrom } from 'rxjs';
import { User, LoginData, RegisterData, RegisterRequest } from '../models/user';
import { HttpClient } from '@angular/common/http';

@Injectable({
  providedIn: 'root',
})
export class AuthService {
  private currentUserSubject = new BehaviorSubject<User | null>(null);
  public currentUser$ = this.currentUserSubject.asObservable();

  private baseUrl = 'http://localhost:8080';

  constructor(private http: HttpClient) {
    const saved = localStorage.getItem('currentUser');
    if (saved) {
      try {
        const parsed = JSON.parse(saved);
        this.currentUserSubject.next(this.normalizeUser(parsed));
      } catch (e) {}
    }
  }

  // Нормализовать объект пользователя с сервера в удобный для фронтенда формат
  private normalizeUser(u: any): any {
    if (!u) return null;
    const out: any = { ...u };
    // привести snake_case is_admin -> isAdmin и привести числовые/строковые значения к булеву
    if (out.isAdmin === undefined) {
      const raw = out.is_admin !== undefined ? out.is_admin : out.isAdmin;
      out.isAdmin = !!(raw === true || raw === 'true' || raw === 't' || raw === 1 || raw === '1');
    }
    // преобразовать created_at -> createdAt, если поле присутствует
    if (out.createdAt === undefined && out.created_at !== undefined) {
      out.createdAt = out.created_at;
    }
    return out;
  }

  async login(loginData: LoginData): Promise<boolean> {
    // Попытаться отправить POST /api/auntificate с JSON телом (серверный endpoint аутентификации)
    try {
      const body = { email: loginData.email, password: loginData.password };
      const resp = await firstValueFrom(
        this.http.post<any>(`${this.baseUrl}/api/login`, body, {
          headers: { 'Content-Type': 'application/json' },
        })
      );

      const user: User | null = resp?.user ? resp.user : resp && resp.id ? resp : null;
      if (user) {
        const norm = this.normalizeUser(user);
        this.currentUserSubject.next(norm);
        localStorage.setItem('currentUser', JSON.stringify(norm));
        return true;
      }
    } catch (err) {
      // Если /api/auntificate недоступен или возвращает ошибку, попробовать /login как запасной вариант
      try {
        const body = { email: loginData.email, password: loginData.password };
        const resp = await firstValueFrom(
          this.http.post<any>(`${this.baseUrl}/login`, body, {
            headers: { 'Content-Type': 'application/json' },
          })
        );
        const user: User | null = resp?.user ? resp.user : resp && resp.id ? resp : null;
        if (user) {
          const norm = this.normalizeUser(user);
          this.currentUserSubject.next(norm);
          localStorage.setItem('currentUser', JSON.stringify(norm));
          return true;
        }
      } catch (postErr) {
        console.error('Login error', postErr);
        return false;
      }
    }

    // Если выполнение дошло до этого места, вход не удался
    return false;
  }

  async register(registerData: RegisterData): Promise<boolean> {
    try {
      const payload: RegisterRequest & { createdAt: string } = {
        username: registerData.username,
        email: registerData.email,
        password: registerData.password,
        createdAt: new Date().toISOString(),
      };

      const created = await firstValueFrom(
        this.http.post<User>(`${this.baseUrl}/api/register`, payload, {
          headers: { 'Content-Type': 'application/json' },
        })
      );

      if (created) {
        const norm = this.normalizeUser(created);
        this.currentUserSubject.next(norm);
        localStorage.setItem('currentUser', JSON.stringify(norm));
        return true;
      }
      return false;
    } catch (err) {
      const status = (err as any)?.status;
      if (status === 409) {
        return false;
      }
      console.error('Register error', err);
      return false;
    }
  }

  logout(): void {
    this.currentUserSubject.next(null);
    localStorage.removeItem('currentUser');
  }

  isLoggedIn(): boolean {
    return this.currentUserSubject.value !== null;
  }
}
