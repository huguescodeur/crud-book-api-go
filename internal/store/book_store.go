package store

import (
	"crud-book/internal/model"
	"database/sql"
)

type Store interface {
	GetAll() ([]*model.Book, error)
	GetByID(id int) (*model.Book, error)
	Create(book *model.Book) (*model.Book, error)
	Update(id int, book *model.Book) (*model.Book, error)
	Delete(id int) error
}

type store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &store{db: db}
}

func (s *store) GetAll() ([]*model.Book, error) {
	q := `
		SELECT id, title, author
		FROM books
	`

	rows, err := s.db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books = make([]*model.Book, 0)
	for rows.Next() {
		var b model.Book
		if err := rows.Scan(&b.ID, &b.Title, &b.Author); err != nil {
			return nil, err
		}

		books = append(books, &b)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}

func (s *store) GetByID(id int) (*model.Book, error) {
	q := `
		SELECT id, title, author
		FROM books
		WHERE id = ?
	`

	var b model.Book
	if err := s.db.QueryRow(q, id).Scan(&b.ID, &b.Title, &b.Author); err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}

		return nil, err
	}

	return &b, nil
}

func (s *store) Create(book *model.Book) (*model.Book, error) {
	q := `
		INSERT INTO books (title, author)
		VALUES (?, ?)
	`
	res, err := s.db.Exec(q, book.Title, book.Author)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	book.ID = int(id)

	return book, nil
}

func (s *store) Update(id int, book *model.Book) (*model.Book, error) {
	q := `
		UPDATE books
		SET title = ?, author = ?
		WHERE id = ?
	`

	res, err := s.db.Exec(q, book.Title, book.Author, id)
	if err != nil {
		return nil, err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}

	if n == 0 {
		return nil, sql.ErrNoRows
	}

	book.ID = id

	return book, nil
}

func (s *store) Delete(id int) error {
	q := `
		DELETE FROM books
		WHERE id = ?
	`
	res, err := s.db.Exec(q, id)
	if err != nil {
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if n == 0 {
		return sql.ErrNoRows
	}

	return nil
}
