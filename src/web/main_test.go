package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"leitor-usbn/api"
	"leitor-usbn/database"
)

// MockDB para testes
type MockDB struct {
	books              []*database.Book
	authors            map[int]*database.Author
	publishers         map[int]*database.Publisher
	saveBookFn         func(*database.Book) (*database.Book, error)
	getOrCreateAuthorFn func(string) (*database.Author, error)
	getOrCreatePublisherFn func(string) (*database.Publisher, error)
	getBooksFn         func() ([]*database.BookDetail, error)
	nextAuthorID       int
	nextPublisherID    int
}

func newMockDB() *MockDB {
	return &MockDB{
		books:       []*database.Book{},
		authors:     make(map[int]*database.Author),
		publishers:  make(map[int]*database.Publisher),
		nextAuthorID: 1,
		nextPublisherID: 1,
	}
}

func (m *MockDB) SaveBook(book *database.Book) (*database.Book, error) {
	if m.saveBookFn != nil {
		return m.saveBookFn(book)
	}
	book.ID = len(m.books) + 1
	m.books = append(m.books, book)
	return book, nil
}

func (m *MockDB) GetOrCreateAuthor(name string) (*database.Author, error) {
	if m.getOrCreateAuthorFn != nil {
		return m.getOrCreateAuthorFn(name)
	}
	for _, a := range m.authors {
		if a.Name == name {
			return a, nil
		}
	}
	author := &database.Author{
		ID:        m.nextAuthorID,
		Name:      name,
		CreatedAt: time.Now(),
	}
	m.authors[author.ID] = author
	m.nextAuthorID++
	return author, nil
}

func (m *MockDB) GetOrCreatePublisher(name string) (*database.Publisher, error) {
	if m.getOrCreatePublisherFn != nil {
		return m.getOrCreatePublisherFn(name)
	}
	for _, p := range m.publishers {
		if p.Name == name {
			return p, nil
		}
	}
	publisher := &database.Publisher{
		ID:        m.nextPublisherID,
		Name:      name,
		CreatedAt: time.Now(),
	}
	m.publishers[publisher.ID] = publisher
	m.nextPublisherID++
	return publisher, nil
}

func (m *MockDB) GetBooksWithDetails() ([]*database.BookDetail, error) {
	if m.getBooksFn != nil {
		return m.getBooksFn()
	}
	return []*database.BookDetail{}, nil
}

func (m *MockDB) Close() error {
	return nil
}

func (m *MockDB) InitSchema() error {
	return nil
}

func (m *MockDB) CountBooks() (int, error) {
	return len(m.books), nil
}

func (m *MockDB) GetAllBooks() ([]*database.Book, error) {
	return m.books, nil
}

func (m *MockDB) GetBookByISBN(isbn string) (*database.Book, error) {
	for _, b := range m.books {
		if b.ISBN == isbn {
			return b, nil
		}
	}
	return nil, nil
}

// TestPOSTAPIBooks testa o handler POST /api/books
func TestPOSTAPIBooks(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
		expectID       bool
		expectTitle    string
	}{
		{
			name: "Criar livro com dados completos",
			requestBody: map[string]interface{}{
				"isbn":        "0132350882",
				"title":       "Clean Code",
				"author":      "Robert C. Martin",
				"publisher":   "Prentice Hall",
				"publish_date": "2008",
				"pages":       464,
				"description": "A Handbook of Agile Software Craftsmanship",
			},
			expectedStatus: http.StatusCreated,
			expectID:       true,
			expectTitle:    "Clean Code",
		},
		{
			name: "Criar livro sem autor/editora",
			requestBody: map[string]interface{}{
				"isbn":  "020161622X",
				"title": "The Pragmatic Programmer",
			},
			expectedStatus: http.StatusCreated,
			expectID:       true,
			expectTitle:    "The Pragmatic Programmer",
		},
		{
			name: "Erro: ISBN obrigatório",
			requestBody: map[string]interface{}{
				"title": "Livro Sem ISBN",
			},
			expectedStatus: http.StatusBadRequest,
			expectID:       false,
		},
		{
			name:           "Erro: JSON inválido",
			requestBody:    nil,
			expectedStatus: http.StatusBadRequest,
			expectID:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Preparar mock
			mockDB := newMockDB()
			db = mockDB

			// Preparar request
			var body []byte
			if tt.requestBody == nil {
				body = []byte("invalid json")
			} else {
				body, _ = json.Marshal(tt.requestBody)
			}

			req := httptest.NewRequest(http.MethodPost, "/api/books", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			// Executar
			w := httptest.NewRecorder()
			handleAPICreateBook(w, req)

			// Verificar status
			if w.Code != tt.expectedStatus {
				t.Errorf("status = %d, want %d", w.Code, tt.expectedStatus)
			}

			// Verificar resposta se sucesso
			if tt.expectedStatus == http.StatusCreated {
				var result database.Book
				json.Unmarshal(w.Body.Bytes(), &result)

				if tt.expectID && result.ID == 0 {
					t.Errorf("livro salvo sem ID")
				}
				if result.Title != tt.expectTitle {
					t.Errorf("title = %q, want %q", result.Title, tt.expectTitle)
				}
			}

			// Verificar resposta se erro
			if tt.expectedStatus >= 400 {
				var errResp map[string]string
				json.Unmarshal(w.Body.Bytes(), &errResp)
				if _, hasError := errResp["error"]; !hasError {
					t.Errorf("resposta de erro sem campo 'error'")
				}
			}
		})
	}
}

// TestGETAPIBooks testa o handler GET /api/books
func TestGETAPIBooks(t *testing.T) {
	mockDB := newMockDB()
	db = mockDB

	req := httptest.NewRequest(http.MethodGet, "/api/books", nil)
	w := httptest.NewRecorder()

	handleAPIGetBooks(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}

	if w.Header().Get("Content-Type") != "application/json; charset=utf-8" {
		t.Errorf("Content-Type = %q, want application/json", w.Header().Get("Content-Type"))
	}
}

// TestGETAPISoBooksInvalidMethod testa método inválido
func TestGETAPISoBooksInvalidMethod(t *testing.T) {
	mockDB := newMockDB()
	db = mockDB

	req := httptest.NewRequest(http.MethodPut, "/api/books", nil)
	w := httptest.NewRecorder()

	handleAPIGetBooks(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("status = %d, want %d", w.Code, http.StatusMethodNotAllowed)
	}
}

type apiISBNResponse struct {
	ISBN        string `json:"isbn"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	Publisher   string `json:"publisher"`
	PublishDate string `json:"publish_date"`
	Pages       int    `json:"pages"`
	Description string `json:"description"`
	CoverURL    string `json:"cover_url"`
}

// TestAPISearchISBN testa o proxy /api/isbn
func TestAPISearchISBN(t *testing.T) {
	tests := []struct {
		name           string
		query          string
		serverStatus   int
		serverBody     string
		expectedStatus int
		expectTitle    string
		expectCover    bool
	}{
		{
			name:           "Sucesso",
			query:          "0132350882",
			serverStatus:   http.StatusOK,
			serverBody:     `{"ISBN:0132350882":{"title":"Clean Code","authors":[{"name":"Robert C. Martin"}],"publishers":[{"name":"Prentice Hall"}],"publish_date":"2008","number_of_pages":464,"description":{"value":"A Handbook"},"covers":[123]}}`,
			expectedStatus: http.StatusOK,
			expectTitle:    "Clean Code",
			expectCover:    true,
		},
		{
			name:           "Nao encontrado",
			query:          "0132350882",
			serverStatus:   http.StatusOK,
			serverBody:     `{}`,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Erro upstream",
			query:          "0132350882",
			serverStatus:   http.StatusInternalServerError,
			serverBody:     `{}`,
			expectedStatus: http.StatusBadGateway,
		},
		{
			name:           "ISBN invalido",
			query:          "123",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "ISBN ausente",
			query:          "",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.query == "" || tt.serverStatus == 0 {
				apiClient = api.NewBookAPIClient("http://127.0.0.1", 1)
			} else {
				srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(tt.serverStatus)
					_, _ = w.Write([]byte(tt.serverBody))
				}))
				defer srv.Close()

				apiClient = api.NewBookAPIClient(srv.URL, 2)
			}

			url := "/api/isbn"
			if tt.query != "" {
				url = url + "?value=" + tt.query
			}

			req := httptest.NewRequest(http.MethodGet, url, nil)
			w := httptest.NewRecorder()
			handleAPISearchISBN(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("status = %d, want %d", w.Code, tt.expectedStatus)
			}

			if tt.expectedStatus == http.StatusOK {
				var resp apiISBNResponse
				_ = json.Unmarshal(w.Body.Bytes(), &resp)
				if resp.Title != tt.expectTitle {
					t.Errorf("title = %q, want %q", resp.Title, tt.expectTitle)
				}
				if tt.expectCover && resp.CoverURL == "" {
					t.Errorf("cover_url vazio, esperado preenchido")
				}
			}
		})
	}
}

// BenchmarkPostAPIBooks testa performance do handler POST
func BenchmarkPostAPIBooks(b *testing.B) {
	mockDB := newMockDB()
	db = mockDB

	body := map[string]interface{}{
		"isbn":        "978-0-13-235088-2",
		"title":       "Clean Code",
		"author":      "Robert C. Martin",
		"publisher":   "Prentice Hall",
		"publish_date": "2008",
		"pages":       464,
	}
	bodyBytes, _ := json.Marshal(body)

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodPost, "/api/books", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		handleAPICreateBook(w, req)
	}
}
