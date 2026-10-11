package service

import (
	"crud-book/internal/model"
	"crud-book/internal/store"
	"database/sql"
	"errors"
	"strings"
)

type Service struct {
	store store.Store
}

func NewService(s store.Store) *Service {
	return &Service{store: s}
}

var (
	ErrNotFound          = errors.New("data not found")
	ErrAllFieldsRequired = errors.New("all fields required")
)

func (s *Service) GetAllBook() ([]*model.Book, error) {
	return s.store.GetAll()
}

func (s *Service) GetBookByID(id int) (*model.Book, error) {
	book, err := s.store.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return book, nil
}

func (s *Service) CreateBook(book *model.Book) (*model.Book, error) {
	if len(strings.TrimSpace(book.Title)) == 0 || len(strings.TrimSpace(book.Author)) == 0 {
		return nil, ErrAllFieldsRequired
	}

	createdBook, err := s.store.Create(book)
	if err != nil {
		return nil, err
	}

	return createdBook, nil
}

func (s *Service) UpdateBook(id int, book *model.Book) (*model.Book, error) {
	b, err := s.store.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}

		return nil, err
	}

	if len(strings.TrimSpace(book.Title)) == 0 {
		book.Title = b.Title
	}

	if len(strings.TrimSpace(book.Author)) == 0 {
		book.Author = b.Author
	}

	updated, err := s.store.Update(id, book)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return updated, nil
}

func (s *Service) DeleteBook(id int) error {
	if err := s.store.Delete(id); err != nil {
		if err == sql.ErrNoRows {
			return ErrNotFound
		}

		return err
	}

	return nil
}
