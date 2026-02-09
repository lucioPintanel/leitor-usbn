package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// OpenLibraryResponse representa a resposta da API OpenLibrary
type OpenLibraryResponse struct {
	ISBN             string       `json:"isbn"`
	Title            string       `json:"title"`
	Authors          []AuthorInfo `json:"authors"`
	NumberOfPages    int          `json:"number_of_pages"`
	Publishers       []struct {
		Name string `json:"name"`
	} `json:"publishers"`
	PublishDate      string      `json:"publish_date"`
	FirstPublishDate string      `json:"first_publish_date"`
	Description      Description `json:"description"`
	Covers           []int        `json:"covers"`
}

type AuthorInfo struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

type Description struct {
	Value string `json:"value"`
}

// BookAPIClient para consumir a API de livros
type BookAPIClient struct {
	baseURL string
	timeout time.Duration
	client  *http.Client
}

// NewBookAPIClient cria uma nova instância do cliente API
func NewBookAPIClient(baseURL string, timeoutSeconds int) *BookAPIClient {
	return &BookAPIClient{
		baseURL: baseURL,
		timeout: time.Duration(timeoutSeconds) * time.Second,
		client: &http.Client{
			Timeout: time.Duration(timeoutSeconds) * time.Second,
		},
	}
}

// GetBookByISBN consulta a API OpenLibrary por ISBN
func (c *BookAPIClient) GetBookByISBN(isbn string) (*OpenLibraryResponse, error) {
	url := fmt.Sprintf("%s?bibkeys=ISBN:%s&format=json&jscmd=data", c.baseURL, isbn)

	resp, err := c.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("erro ao fazer requisição para ISBN %s: %w", isbn, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code inválido para ISBN %s: %d", isbn, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler resposta para ISBN %s: %w", isbn, err)
	}

	// Parse a resposta JSON
	var result map[string]OpenLibraryResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, fmt.Errorf("erro ao fazer parse JSON para ISBN %s: %w", isbn, err)
	}

	// Procurar a chave ISBN
	key := fmt.Sprintf("ISBN:%s", isbn)
	if book, exists := result[key]; exists {
		book.ISBN = isbn
		return &book, nil
	}

	return nil, fmt.Errorf("ISBN %s não encontrado na API", isbn)
}

// GetBookByISBNWithRetry tenta obter o livro com retry automático
func (c *BookAPIClient) GetBookByISBNWithRetry(isbn string, maxRetries int) (*OpenLibraryResponse, error) {
	var lastErr error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		book, err := c.GetBookByISBN(isbn)
		if err == nil {
			return book, nil
		}

		lastErr = err

		if attempt < maxRetries {
			// Aguardar antes de tentar novamente (backoff exponencial)
			wait := time.Duration(attempt*attempt) * time.Second
			time.Sleep(wait)
		}
	}

	return nil, fmt.Errorf("falha após %d tentativas para ISBN: %w", maxRetries, lastErr)
}
