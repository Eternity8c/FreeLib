package book_service

import (
	"FreeLib/internal/core/domain"
	core_errors "FreeLib/internal/core/errors"
	"context"
	"fmt"
)

type BookService struct {
	bookRepository BookRepository
}

type BookRepository interface {
	CreateBook(ctx context.Context, book domain.Book) (domain.Book, error)
	GetBooks(ctx context.Context, limit *int, offset *int) ([]domain.Book, error)
	GetNewBooks(ctx context.Context, limit *int, offset *int) ([]domain.Book, error)
	GetBook(ctx context.Context, id int) (domain.Book, error)
	FavoriteBook(ctx context.Context, userID int, bookID int) (int, domain.Book, error)
	GetFavoriteBooks(ctx context.Context, userID int) ([]domain.Book, error)
}

func NewBookService(bookRepository BookRepository) *BookService {
	return &BookService{
		bookRepository: bookRepository,
	}
}

func (s *BookService) CreateBook(ctx context.Context, book domain.Book) (domain.Book, error) {
	if err := book.Validate(); err != nil {
		return domain.Book{}, fmt.Errorf("validate book domain: %w", err)
	}

	domainBook, err := s.bookRepository.CreateBook(ctx, book)
	if err != nil {
		return domain.Book{}, fmt.Errorf("create book: %w", err)
	}

	return domainBook, nil
}

func (s *BookService) FavoriteBook(ctx context.Context, userID int, bookID int) (int, domain.Book, error) {
	if userID == domain.UninitializedID {
		return domain.UninitializedID, domain.Book{}, fmt.Errorf("userID must be non-negative: %w", core_errors.ErrInvalidArgumment)
	}
	if bookID <= 0 {
		return domain.UninitializedID, domain.Book{}, fmt.Errorf("bookID must be non-negative: %w", core_errors.ErrInvalidArgumment)
	}

	uID, domainBook, err := s.bookRepository.FavoriteBook(ctx, userID, bookID)
	if err != nil {
		return domain.UninitializedID, domain.Book{}, fmt.Errorf("favorite book from repository: %w", err)
	}

	return uID, domainBook, err
}

func (s *BookService) GetBook(ctx context.Context, id int) (domain.Book, error) {
	if err := validateID(id); err != nil {
		return domain.Book{}, fmt.Errorf("validate ID: %w", err)
	}

	book, err := s.bookRepository.GetBook(ctx, id)
	if err != nil {
		return domain.Book{}, fmt.Errorf("get book from repository: %w", err)
	}

	return book, nil
}

func (s *BookService) GetBooks(ctx context.Context, limit *int, offset *int) ([]domain.Book, error) {
	if err := validateLimitOffset(limit, offset); err != nil {
		return nil, fmt.Errorf("validate limit offset: %w", err)
	}

	books, err := s.bookRepository.GetBooks(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get books from repository: %w", err)
	}

	return books, nil
}

func (s *BookService) GetFavoriteBooks(ctx context.Context, userID int) ([]domain.Book, error) {
	bookDomains, err := s.bookRepository.GetFavoriteBooks(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get book domains from repository: %w", err)
	}

	return bookDomains, nil
}

func (s *BookService) GetNewBooks(ctx context.Context, limit *int, offset *int) ([]domain.Book, error) {
	if err := validateLimitOffset(limit, offset); err != nil {
		return nil, fmt.Errorf("validate limit offset: %w", err)
	}

	books, err := s.bookRepository.GetNewBooks(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get new books from repository")
	}

	return books, nil
}
