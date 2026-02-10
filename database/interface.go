package database

import (
	"context"
)

// DatabasePort define a interface para operações de banco de dados
// Facilita mocking em testes e desacopla processor de implementação específica
type DatabasePort interface {
	// InitSchema cria as tabelas se não existirem
	InitSchema() error

	// GetOrCreateAuthor obtém um autor existente ou cria um novo
	GetOrCreateAuthor(name string) (*Author, error)

	// GetOrCreatePublisher obtém uma editora existente ou cria uma nova
	GetOrCreatePublisher(name string) (*Publisher, error)

	// SaveBook salva um livro no banco de dados (cria ou atualiza)
	SaveBook(book *Book) (*Book, error)

	// GetBookByISBN obtém um livro pelo ISBN
	GetBookByISBN(isbn string) (*Book, error)

	// GetAllBooks retorna todos os livros
	GetAllBooks() ([]*Book, error)

	// GetBooksWithDetails retorna livros com nomes de autor e editora
	GetBooksWithDetails() ([]*BookDetail, error)

	// CountBooks retorna o número total de livros
	CountBooks() (int, error)

	// Close fecha a conexão com o banco de dados
	Close() error
}

// Guarantee que Database implementa DatabasePort
var _ DatabasePort = (*Database)(nil)

// GetBookByISBNWithContext obtém um livro pelo ISBN respeitando context
func GetBookByISBNWithContext(ctx context.Context, db DatabasePort, isbn string) (*Book, error) {
	// Para agora, apenas chama sem context
	// Future: implementar timeout/cancelamento respeitando context
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return db.GetBookByISBN(isbn)
	}
}
