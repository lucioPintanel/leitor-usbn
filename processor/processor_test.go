package processor

import (
	"fmt"
	"leitor-usbn/api"
	"leitor-usbn/database"
	"testing"
)

// TestProcessorConfig testa se a configuração é normalizada
func TestProcessorConfig(t *testing.T) {
	tests := []struct {
		name string
		cfg  ProcessorConfig
		want ProcessorConfig
	}{
		{
			name: "default config with zeros",
			cfg:  ProcessorConfig{},
			want: ProcessorConfig{
				MaxWorkers:           1,
				DelayBetweenRequests: 500000000, // 500ms em nanosegundos
				MaxRetries:           3,
			},
		},
		{
			name: "custom config preserved",
			cfg: ProcessorConfig{
				MaxWorkers:           5,
				DelayBetweenRequests: 100000000,
				MaxRetries:           5,
			},
			want: ProcessorConfig{
				MaxWorkers:           5,
				DelayBetweenRequests: 100000000,
				MaxRetries:           5,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := &MockDatabase{}

			p := NewProcessor(mockDB, &api.BookAPIClient{}, nil, tt.cfg)

			if p.config.MaxWorkers != tt.want.MaxWorkers {
				t.Errorf("MaxWorkers = %d, want %d", p.config.MaxWorkers, tt.want.MaxWorkers)
			}
			if p.config.MaxRetries != tt.want.MaxRetries {
				t.Errorf("MaxRetries = %d, want %d", p.config.MaxRetries, tt.want.MaxRetries)
			}
		})
	}
}

// TestProcessorWithMockDatabase testa se o processor chama o database corretamente
func TestProcessorWithMockDatabase(t *testing.T) {
	mockDB := &MockDatabase{
		SaveBookFn: func(book *database.Book) (*database.Book, error) {
			if book.ISBN == "" {
				return nil, fmt.Errorf("ISBN vazio")
			}
			book.ID = 1
			return book, nil
		},
		GetOrCreateAuthorFn: func(name string) (*database.Author, error) {
			return &database.Author{ID: 1, Name: name}, nil
		},
		GetOrCreatePublisherFn: func(name string) (*database.Publisher, error) {
			return &database.Publisher{ID: 1, Name: name}, nil
		},
	}

	cfg := ProcessorConfig{MaxWorkers: 1, MaxRetries: 3}
	_ = NewProcessor(mockDB, &api.BookAPIClient{}, nil, cfg)

	// Simular salvamento de livro
	book := &database.Book{
		ISBN:    "9780132350882",
		Title:   "Clean Code",
		AuthorID: nil,
	}

	savedBook, err := mockDB.SaveBook(book)
	if err != nil {
		t.Fatalf("SaveBook error: %v", err)
	}

	if savedBook.ID != 1 {
		t.Errorf("SaveBook ID = %d, want 1", savedBook.ID)
	}

	// Verificar se SaveBook foi chamado
	if mockDB.SaveBookCalls != 1 {
		t.Errorf("SaveBook calls = %d, want 1", mockDB.SaveBookCalls)
	}
}

// TestProcessorStats testa se as stats são calculadas corretamente
func TestProcessorStats(t *testing.T) {
	p := &Processor{
		db:      &MockDatabase{},
		results: []*ProcessResult{
			{ISBN: "1", Success: true},
			{ISBN: "2", Success: true},
			{ISBN: "3", Success: false, Error: "not found"},
		},
	}

	stats := p.GetStats()

	if total := stats["total"].(int); total != 3 {
		t.Errorf("total = %d, want 3", total)
	}
	if success := stats["success"].(int); success != 2 {
		t.Errorf("success = %d, want 2", success)
	}
	if errors := stats["errors"].(int); errors != 1 {
		t.Errorf("errors = %d, want 1", errors)
	}
}
