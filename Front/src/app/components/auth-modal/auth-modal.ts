import { Component, EventEmitter, Output, Input, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { AuthService } from '../../services/auth';
import { LoginData, RegisterData } from '../../models/user';

@Component({
  selector: 'app-auth-modal',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './auth-modal.html',
  styleUrls: ['./auth-modal.scss'],
})
export class AuthModalComponent implements OnInit {
  @Input() mode!: 'login' | 'register';
  @Output() modalClosed = new EventEmitter<void>();

  isLoginMode = true;
  loginData: LoginData = { email: '', password: '' };
  registerData: RegisterData = {
    username: '',
    email: '',
    password: '',
    confirmPassword: '',
  };

  constructor(private authService: AuthService) {}

  ngOnInit() {
    if (this.mode) {
      this.isLoginMode = this.mode === 'login';
    }
  }

  closeModal(): void {
    this.modalClosed.emit();
  }

  onLogin(): void {
    (async () => {
      const ok = await this.authService.login(this.loginData);
      if (ok) {
        this.closeModal();
      } else {
        alert('Неверный email или пароль');
      }
    })();
  }

  onRegister(): void {
    if (this.registerData.password !== this.registerData.confirmPassword) {
      alert('Пароли не совпадают');
      return;
    }

    (async () => {
      const ok = await this.authService.register(this.registerData);
      if (ok) {
        this.closeModal();
      } else {
        alert('Ошибка регистрации — возможно пользователь с таким email уже существует');
      }
    })();
  }
}
