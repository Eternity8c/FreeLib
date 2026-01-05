package repository

import "FreeLib/internal/models"

type BookRepository interface {
	GetAll() ([]models.Book, error)
	GetByID(id uint) (*models.Book, error)
	Create(book *models.Book) error
	Delete(id uint) error
	Update(book *models.Book) error
	AddFavorite(userID uint, bookID uint) error
	DeleteFavorite(userID uint, bookID uint) error
	GetAllFavorite(userID uint) ([]models.Book, error)
}
