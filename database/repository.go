package database

import (
	"database/sql"
	"fmt"
	"time"
)

// Author representa um autor no banco de dados
type Author struct {
	ID        int
	Name      string
	CreatedAt time.Time
}

// Publisher representa uma editora no banco de dados
type Publisher struct {
	ID        int
	Name      string
	CreatedAt time.Time
}

// Book representa um livro no banco de dados
type Book struct {
	ID          int
	ISBN        string
	Title       string
	AuthorID    *int
	PublisherID *int
	PublishDate string
	Pages       int
	Description string
	CoverURL    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// GetOrCreateAuthor obtém um autor existente ou cria um novo
func (db *Database) GetOrCreateAuthor(name string) (*Author, error) {
	if name == "" {
		return nil, nil
	}

	var author Author

	// Tentar obter autor existente
	err := db.conn.QueryRow(
		"SELECT id, name, created_at FROM authors WHERE name = ?",
		name,
	).Scan(&author.ID, &author.Name, &author.CreatedAt)

	if err == nil {
		return &author, nil
	}

	if err != sql.ErrNoRows {
		return nil, fmt.Errorf("erro ao buscar autor: %w", err)
	}

	// Criar novo autor
	result, err := db.conn.Exec(
		"INSERT INTO authors (name) VALUES (?)",
		name,
	)

	if err != nil {
		return nil, fmt.Errorf("erro ao criar autor: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("erro ao obter ID do autor: %w", err)
	}

	author.ID = int(id)
	author.Name = name
	author.CreatedAt = time.Now()

	return &author, nil
}

// GetOrCreatePublisher obtém uma editora existente ou cria uma nova
func (db *Database) GetOrCreatePublisher(name string) (*Publisher, error) {
	if name == "" {
		return nil, nil
	}

	var publisher Publisher

	// Tentar obter editora existente
	err := db.conn.QueryRow(
		"SELECT id, name, created_at FROM publishers WHERE name = ?",
		name,
	).Scan(&publisher.ID, &publisher.Name, &publisher.CreatedAt)

	if err == nil {
		return &publisher, nil
	}

	if err != sql.ErrNoRows {
		return nil, fmt.Errorf("erro ao buscar editora: %w", err)
	}

	// Criar nova editora
	result, err := db.conn.Exec(
		"INSERT INTO publishers (name) VALUES (?)",
		name,
	)

	if err != nil {
		return nil, fmt.Errorf("erro ao criar editora: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("erro ao obter ID da editora: %w", err)
	}

	publisher.ID = int(id)
	publisher.Name = name
	publisher.CreatedAt = time.Now()

	return &publisher, nil
}

// SaveBook salva um livro no banco de dados (cria ou atualiza)
func (db *Database) SaveBook(book *Book) (*Book, error) {
	now := time.Now()

	// Verificar se o livro já existe
	var existingID int
	err := db.conn.QueryRow("SELECT id FROM books WHERE isbn = ?", book.ISBN).Scan(&existingID)

	if err == nil {
		// Atualizar livro existente
		_, err := db.conn.Exec(`
			UPDATE books 
			SET title = ?, author_id = ?, publisher_id = ?, 
			    publish_date = ?, pages = ?, description = ?, 
			    cover_url = ?, updated_at = ?
			WHERE isbn = ?
		`,
			book.Title, book.AuthorID, book.PublisherID,
			book.PublishDate, book.Pages, book.Description,
			book.CoverURL, now, book.ISBN,
		)

		if err != nil {
			return nil, fmt.Errorf("erro ao atualizar livro: %w", err)
		}

		book.UpdatedAt = now
		return book, nil
	}

	if err != sql.ErrNoRows {
		return nil, fmt.Errorf("erro ao verificar livro existente: %w", err)
	}

	// Criar novo livro
	result, err := db.conn.Exec(`
		INSERT INTO books (isbn, title, author_id, publisher_id, publish_date, pages, description, cover_url, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		book.ISBN, book.Title, book.AuthorID, book.PublisherID,
		book.PublishDate, book.Pages, book.Description, book.CoverURL, now, now,
	)

	if err != nil {
		return nil, fmt.Errorf("erro ao criar livro: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("erro ao obter ID do livro: %w", err)
	}

	book.ID = int(id)
	book.CreatedAt = now
	book.UpdatedAt = now

	return book, nil
}

// GetBookByISBN obtém um livro pelo ISBN
func (db *Database) GetBookByISBN(isbn string) (*Book, error) {
	var book Book

	err := db.conn.QueryRow(`
		SELECT id, isbn, title, author_id, publisher_id, publish_date, pages, description, cover_url, created_at, updated_at
		FROM books
		WHERE isbn = ?
	`, isbn).Scan(&book.ID, &book.ISBN, &book.Title, &book.AuthorID, &book.PublisherID,
		&book.PublishDate, &book.Pages, &book.Description, &book.CoverURL, &book.CreatedAt, &book.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("erro ao buscar livro: %w", err)
	}

	return &book, nil
}

// GetAllBooks retorna todos os livros com informações de autor e editora
func (db *Database) GetAllBooks() ([]*Book, error) {
	rows, err := db.conn.Query(`
		SELECT b.id, b.isbn, b.title, b.author_id, b.publisher_id, b.publish_date, 
		       b.pages, b.description, b.cover_url, b.created_at, b.updated_at
		FROM books b
		ORDER BY b.created_at DESC
	`)

	if err != nil {
		return nil, fmt.Errorf("erro ao buscar livros: %w", err)
	}
	defer rows.Close()

	var books []*Book
	for rows.Next() {
		var book Book
		err := rows.Scan(&book.ID, &book.ISBN, &book.Title, &book.AuthorID, &book.PublisherID,
			&book.PublishDate, &book.Pages, &book.Description, &book.CoverURL, &book.CreatedAt, &book.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("erro ao fazer scan do livro: %w", err)
		}
		books = append(books, &book)
	}

	return books, rows.Err()
}

// CountBooks retorna o número total de livros
func (db *Database) CountBooks() (int, error) {
	var count int
	err := db.conn.QueryRow("SELECT COUNT(*) FROM books").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("erro ao contar livros: %w", err)
	}
	return count, nil
}
