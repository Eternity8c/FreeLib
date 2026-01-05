import { Routes } from '@angular/router';
import { Library } from './pages/library/library';
import { Home } from './pages/home/home';
import { AdminBooks } from './pages/admin/admin';
import { NewBooks } from './pages/new/new';
import { Favorites } from './pages/favorites/favorites';

export const routes: Routes = [
  { path: 'library', component: Library },
  { path: 'admin', component: AdminBooks },
  { path: 'new', component: NewBooks },
  { path: 'favorites', component: Favorites },
  { path: '', component: Home },
];
