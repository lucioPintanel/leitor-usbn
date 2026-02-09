package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// Database encapsula a conexão e operações com SQLite
type Database struct {
	conn *sql.DB
}

// NewDatabase cria uma nova conexão com o banco de dados
func NewDatabase(filepath string) (*Database, error) {
	conn, err := sql.Open("sqlite3", filepath)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir banco de dados: %w", err)
	}

	// Testar a conexão
	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("erro ao conectar ao banco de dados: %w", err)
	}

	return &Database{conn: conn}, nil
}

// InitSchema cria as tabelas se não existirem
func (db *Database) InitSchema() error {
	schema := `
	-- Tabela de autores
	CREATE TABLE IF NOT EXISTS authors (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	-- Tabela de editoras
	CREATE TABLE IF NOT EXISTS publishers (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	-- Tabela de livros
	CREATE TABLE IF NOT EXISTS books (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		isbn TEXT NOT NULL UNIQUE,
		title TEXT NOT NULL,
		author_id INTEGER,
		publisher_id INTEGER,
		publish_date TEXT,
		pages INTEGER,
		description TEXT,
		cover_url TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (author_id) REFERENCES authors(id),
		FOREIGN KEY (publisher_id) REFERENCES publishers(id)
	);

	-- Índices para melhorar performance
	CREATE INDEX IF NOT EXISTS idx_books_isbn ON books(isbn);
	CREATE INDEX IF NOT EXISTS idx_books_author_id ON books(author_id);
	CREATE INDEX IF NOT EXISTS idx_books_publisher_id ON books(publisher_id);
	CREATE INDEX IF NOT EXISTS idx_authors_name ON authors(name);
	CREATE INDEX IF NOT EXISTS idx_publishers_name ON publishers(name);
	`

	_, err := db.conn.Exec(schema)
	if err != nil {
		return fmt.Errorf("erro ao criar schema: %w", err)
	}

	return nil
}

// Close fecha a conexão com o banco de dados
func (db *Database) Close() error {
	return db.conn.Close()
}

// GetConnection retorna a conexão bruta (se necessário)
func (db *Database) GetConnection() *sql.DB {
	return db.conn
}
