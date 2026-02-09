package database

import (
	"database/sql"
	"fmt"
	"time"
)

// BookDetail representa um livro com nomes de autor e editora
type BookDetail struct {
	ID           int
	ISBN         string
	Title        string
	AuthorID     *int
	AuthorName   string
	PublisherID  *int
	PublisherName string
	PublishDate  string
	Pages        int
	Description  string
	CoverURL     string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// GetBooksWithDetails retorna livros junto com nome do autor e editora
func (db *Database) GetBooksWithDetails() ([]*BookDetail, error) {
	query := `
	SELECT b.id, b.isbn, b.title, b.author_id, a.name as author_name, b.publisher_id, p.name as publisher_name,
	       b.publish_date, b.pages, b.description, b.cover_url, b.created_at, b.updated_at
	FROM books b
	LEFT JOIN authors a ON b.author_id = a.id
	LEFT JOIN publishers p ON b.publisher_id = p.id
	ORDER BY b.created_at DESC
	`

	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erro ao consultar livros com detalhes: %w", err)
	}
	defer rows.Close()

	var results []*BookDetail
	for rows.Next() {
		var d BookDetail
		var authorName sql.NullString
		var publisherName sql.NullString
		var authorID sql.NullInt64
		var publisherID sql.NullInt64

		err := rows.Scan(&d.ID, &d.ISBN, &d.Title, &authorID, &authorName, &publisherID, &publisherName,
			&d.PublishDate, &d.Pages, &d.Description, &d.CoverURL, &d.CreatedAt, &d.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("erro ao ler linha de resultado: %w", err)
		}

		if authorID.Valid {
			tmp := int(authorID.Int64)
			d.AuthorID = &tmp
		}
		if authorName.Valid {
			d.AuthorName = authorName.String
		}
		if publisherID.Valid {
			tmp := int(publisherID.Int64)
			d.PublisherID = &tmp
		}
		if publisherName.Valid {
			d.PublisherName = publisherName.String
		}

		results = append(results, &d)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro em rows: %w", err)
	}

	return results, nil
}
