export interface User {
  id?: number;
  username: string;
  email: string;
  password?: string;
  createdAt?: Date;
}

// Структуры, соответствующие запросам на стороне сервера
export interface RegisterRequest {
  username: string;
  email: string;
  password: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

// Типы для UI (сохраняем confirmPassword для валидации на клиенте)
export interface LoginData extends LoginRequest {}

export interface RegisterData extends RegisterRequest {
  confirmPassword: string;
}
