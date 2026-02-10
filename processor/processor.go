package processor

import (
	"context"
	"fmt"
	"log"
	"leitor-usbn/api"
	"leitor-usbn/database"
	"leitor-usbn/reader"
	"sync"
	"time"
)

// ProcessorConfig contém configurações para o processador
type ProcessorConfig struct {
	MaxWorkers           int           // Número de workers para processar ISBNs em paralelo
	DelayBetweenRequests time.Duration // Delay entre requisições à API
	MaxRetries           int           // Máximo de tentativas por ISBN
	Verbose              bool          // Modo verbose
}

// ProcessResult contém o resultado do processamento de um ISBN
type ProcessResult struct {
	ISBN      string
	Success   bool
	Error     string
	Book      *database.Book
	Timestamp time.Time
}

// Processor orquestra a leitura, consulta e armazenamento de livros
type Processor struct {
	db        database.DatabasePort
	apiClient *api.BookAPIClient
	reader    reader.ISBNReader
	config    ProcessorConfig
	results   []*ProcessResult
	mu        sync.Mutex
}

// NewProcessor cria uma nova instância do processador
func NewProcessor(
	db database.DatabasePort,
	apiClient *api.BookAPIClient,
	isbnReader reader.ISBNReader,
	config ProcessorConfig,
) *Processor {
	if config.MaxWorkers <= 0 {
		config.MaxWorkers = 1
	}
	if config.DelayBetweenRequests == 0 {
		config.DelayBetweenRequests = 500 * time.Millisecond
	}
	if config.MaxRetries <= 0 {
		config.MaxRetries = 3
	}

	return &Processor{
		db:        db,
		apiClient: apiClient,
		reader:    isbnReader,
		config:    config,
		results:   make([]*ProcessResult, 0),
	}
}

// Process inicia o processamento de ISBNs
func (p *Processor) Process(ctx context.Context) error {
	if !p.reader.IsRunning() {
		return fmt.Errorf("leitor não está ativo")
	}

	if p.config.Verbose {
		log.Printf("Iniciando processamento com %d workers", p.config.MaxWorkers)
	}

	// Criar workers
	isbnChan := p.reader.Read()
	var wg sync.WaitGroup

	for i := 0; i < p.config.MaxWorkers; i++ {
		wg.Add(1)
		go p.worker(ctx, &wg, isbnChan, i+1)
	}

	wg.Wait()

	if p.config.Verbose {
		log.Println("Todos os workers finalizados")
	}

	return nil
}

// worker processa ISBNs do canal
func (p *Processor) worker(ctx context.Context, wg *sync.WaitGroup, isbnChan <-chan string, workerID int) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			if p.config.Verbose {
				log.Printf("Worker %d: contexto cancelado", workerID)
			}
			return

		case isbn, ok := <-isbnChan:
			if !ok {
				if p.config.Verbose {
					log.Printf("Worker %d: canal fechado", workerID)
				}
				return
			}

			result := p.processISBN(ctx, isbn)
			p.addResult(result)

			if p.config.Verbose {
				status := "✓"
				if !result.Success {
					status = "✗"
				}
				log.Printf("Worker %d: %s ISBN %s - %s", workerID, status, isbn, result.Error)
			}

			// Delay entre requisições
			select {
			case <-time.After(p.config.DelayBetweenRequests):
			case <-ctx.Done():
				return
			}
		}
	}
}

// processISBN processa um ISBN individual
func (p *Processor) processISBN(ctx context.Context, isbn string) *ProcessResult {
	result := &ProcessResult{
		ISBN:      isbn,
		Timestamp: time.Now(),
	}

	// Consultar API com retry
	var apiBook *api.OpenLibraryResponse
	var err error

	for attempt := 1; attempt <= p.config.MaxRetries; attempt++ {
		apiBook, err = p.apiClient.GetBookByISBN(isbn)
		if err == nil {
			break
		}

		if attempt < p.config.MaxRetries {
			select {
			case <-time.After(time.Duration(attempt*attempt) * time.Second):
			case <-ctx.Done():
				result.Error = fmt.Sprintf("Contexto cancelado na tentativa %d", attempt)
				return result
			}
		}
	}

	if err != nil {
		result.Error = fmt.Sprintf("Erro ao consultar API (tentou %d vezes): %v", p.config.MaxRetries, err)
		return result
	}

	// Converter dados
	bookData := api.ConvertToBookData(apiBook)

	// Obter ou criar autor
	var author *database.Author
	if bookData.Author != "" {
		author, err = p.db.GetOrCreateAuthor(bookData.Author)
		if err != nil {
			result.Error = fmt.Sprintf("Erro ao criar autor: %v", err)
			return result
		}
	}

	// Obter ou criar editora
	var publisher *database.Publisher
	if bookData.Publisher != "" {
		publisher, err = p.db.GetOrCreatePublisher(bookData.Publisher)
		if err != nil {
			result.Error = fmt.Sprintf("Erro ao criar editora: %v", err)
			return result
		}
	}

	// Criar estrutura de livro
	dbBook := &database.Book{
		ISBN:        bookData.ISBN,
		Title:       bookData.Title,
		PublishDate: bookData.PublishDate,
		Pages:       bookData.Pages,
		Description: bookData.Description,
		CoverURL:    bookData.CoverURL,
	}

	if author != nil {
		dbBook.AuthorID = &author.ID
	}
	if publisher != nil {
		dbBook.PublisherID = &publisher.ID
	}

	// Salvar no banco
	savedBook, err := p.db.SaveBook(dbBook)
	if err != nil {
		result.Error = fmt.Sprintf("Erro ao salvar no banco: %v", err)
		return result
	}

	result.Success = true
	result.Book = savedBook
	return result
}

// addResult adiciona um resultado de forma thread-safe
func (p *Processor) addResult(result *ProcessResult) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.results = append(p.results, result)
}

// GetResults retorna todos os resultados
func (p *Processor) GetResults() []*ProcessResult {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Retornar cópia para evitar race conditions
	results := make([]*ProcessResult, len(p.results))
	copy(results, p.results)
	return results
}

// GetStats retorna estatísticas do processamento
func (p *Processor) GetStats() map[string]interface{} {
	results := p.GetResults()

	successCount := 0
	errorCount := 0
	for _, r := range results {
		if r.Success {
			successCount++
		} else {
			errorCount++
		}
	}

	return map[string]interface{}{
		"total":   len(results),
		"success": successCount,
		"errors":  errorCount,
	}
}

// PrintSummary imprime um resumo dos resultados
func (p *Processor) PrintSummary() {
	results := p.GetResults()
	stats := p.GetStats()

	fmt.Println("\n========== RESUMO DO PROCESSAMENTO ==========")
	fmt.Printf("Total de ISBNs processados: %d\n", stats["total"])
	fmt.Printf("Sucesso: %d\n", stats["success"])
	fmt.Printf("Erros: %d\n", stats["errors"])

	if stats["errors"].(int) > 0 {
		fmt.Println("\n--- ISBNs com Erro ---")
		for _, r := range results {
			if !r.Success {
				fmt.Printf("  %s: %s\n", r.ISBN, r.Error)
			}
		}
	}

	if stats["success"].(int) > 0 {
		fmt.Println("\n--- ISBNs com Sucesso ---")
		for _, r := range results {
			if r.Success && r.Book != nil {
				fmt.Printf("  %s: %s\n", r.ISBN, r.Book.Title)
			}
		}
	}

	fmt.Println("\n✓ Processamento concluído!")
}
