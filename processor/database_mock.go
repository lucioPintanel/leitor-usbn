package processor

import (
	"leitor-usbn/database"
)

// MockDatabase é um mock de database.DatabasePort para testes
type MockDatabase struct {
	// Comportamentos configuráveis
	SaveBookFn       func(*database.Book) (*database.Book, error)
	GetOrCreateAuthorFn func(name string) (*database.Author, error)
	GetOrCreatePublisherFn func(name string) (*database.Publisher, error)
	CountBooksFn     func() (int, error)
	InitSchemaFn     func() error
	CloseFn          func() error

	// Rastreamento de chamadas (para asserts)
	SaveBookCalls      int
	GetOrCreateAuthorCalls int
	GetOrCreatePublisherCalls int
	CountBooksCalls    int
}

// SaveBook implementa DatabasePort
func (m *MockDatabase) SaveBook(book *database.Book) (*database.Book, error) {
	m.SaveBookCalls++
	if m.SaveBookFn != nil {
		return m.SaveBookFn(book)
	}
	return book, nil
}

// GetOrCreateAuthor implementa DatabasePort
func (m *MockDatabase) GetOrCreateAuthor(name string) (*database.Author, error) {
	m.GetOrCreateAuthorCalls++
	if m.GetOrCreateAuthorFn != nil {
		return m.GetOrCreateAuthorFn(name)
	}
	return &database.Author{ID: 1, Name: name}, nil
}

// GetOrCreatePublisher implementa DatabasePort
func (m *MockDatabase) GetOrCreatePublisher(name string) (*database.Publisher, error) {
	m.GetOrCreatePublisherCalls++
	if m.GetOrCreatePublisherFn != nil {
		return m.GetOrCreatePublisherFn(name)
	}
	return &database.Publisher{ID: 1, Name: name}, nil
}

// CountBooks implementa DatabasePort
func (m *MockDatabase) CountBooks() (int, error) {
	m.CountBooksCalls++
	if m.CountBooksFn != nil {
		return m.CountBooksFn()
	}
	return 0, nil
}

// InitSchema implementa DatabasePort
func (m *MockDatabase) InitSchema() error {
	if m.InitSchemaFn != nil {
		return m.InitSchemaFn()
	}
	return nil
}

// Close implementa DatabasePort
func (m *MockDatabase) Close() error {
	if m.CloseFn != nil {
		return m.CloseFn()
	}
	return nil
}

// GetBookByISBN implementa DatabasePort
func (m *MockDatabase) GetBookByISBN(isbn string) (*database.Book, error) {
	return nil, nil
}

// GetAllBooks implementa DatabasePort
func (m *MockDatabase) GetAllBooks() ([]*database.Book, error) {
	return nil, nil
}

// GetBooksWithDetails implementa DatabasePort
func (m *MockDatabase) GetBooksWithDetails() ([]*database.BookDetail, error) {
	return nil, nil
}
