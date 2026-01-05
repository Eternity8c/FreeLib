import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { firstValueFrom } from 'rxjs';
import { Book } from './books.service';

@Injectable({ providedIn: 'root' })
export class FavoritesService {
  private baseUrl = 'http://localhost:8080/api';
  constructor(private http: HttpClient) {}

  async addFavorite(userId: string | number, bookId: string | number): Promise<any> {
    const payload = { book_id: bookId };
    return firstValueFrom(
      this.http.post<any>(`${this.baseUrl}/users/${userId}/favorites`, payload)
    );
  }

  async deleteFavorite(userId: string | number, bookId: string | number): Promise<any> {
    return firstValueFrom(
      this.http.delete<any>(`${this.baseUrl}/users/${userId}/favorites/book/${bookId}`)
    );
  }

  async getFavorites(userId: string | number): Promise<Book[]> {
    return firstValueFrom(this.http.get<Book[]>(`${this.baseUrl}/users/${userId}/favorites`));
  }
}
