package main

import (
	"context"
	"fmt"
	"log"
	"leitor-usbn/api"
	"leitor-usbn/database"
	"leitor-usbn/reader"
	"time"
)

func main() {
	fmt.Println("=== Teste Etapa 4: Leitura de ISBNs ===\n")

	// Inicializar banco de dados
	fmt.Println("[1] Inicializando banco de dados...")
	db, err := database.NewDatabase("./books_etapa4.db")
	if err != nil {
		log.Fatalf("Erro ao inicializar banco de dados: %v", err)
	}
	defer db.Close()

	err = db.InitSchema()
	if err != nil {
		log.Fatalf("Erro ao criar schema: %v", err)
	}
	fmt.Println("✓ Banco de dados inicializado\n")

	// Inicializar cliente API
	fmt.Println("[2] Inicializando cliente API...")
	apiClient := api.NewBookAPIClient("https://openlibrary.org/api/books", 10)
	fmt.Println("✓ Cliente API criado\n")

	// Criar leitor de arquivo
	fmt.Println("[3] Configurando leitor de ISBNs (arquivo)...")
	readerConfig := reader.ReaderConfig{
		FilePath: "./config/isbn_list.txt",
		Verbose:  true,
	}

	isbnReader := reader.NewFileISBNReader(readerConfig)
	fmt.Printf("✓ Leitor de arquivo criado: %s\n\n", isbnReader.GetType())

	// Iniciar leitor com contexto
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = isbnReader.Start(ctx)
	if err != nil {
		log.Fatalf("Erro ao iniciar leitor: %v", err)
	}

	fmt.Println("[4] Iniciando leitura de ISBNs...\n")

	// Processar ISBNs
	processedCount := 0
	successCount := 0
	errorCount := 0

	for isbn := range isbnReader.Read() {
		processedCount++

		fmt.Printf("\n--- Processando ISBN %d: %s ---\n", processedCount, isbn)

		// Consultar API
		apiBook, err := apiClient.GetBookByISBN(isbn)
		if err != nil {
			fmt.Printf("❌ Erro ao consultar API: %v\n", err)
			errorCount++
			continue
		}

		// Converter dados
		bookData := api.ConvertToBookData(apiBook)

		// Salvar no banco
		author, err := db.GetOrCreateAuthor(bookData.Author)
		if err != nil {
			fmt.Printf("❌ Erro ao criar autor: %v\n", err)
			errorCount++
			continue
		}

		publisher, err := db.GetOrCreatePublisher(bookData.Publisher)
		if err != nil {
			fmt.Printf("❌ Erro ao criar editora: %v\n", err)
			errorCount++
			continue
		}

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

		_, err = db.SaveBook(dbBook)
		if err != nil {
			fmt.Printf("❌ Erro ao salvar livro: %v\n", err)
			errorCount++
			continue
		}

		fmt.Printf("✓ Sucesso! Livro: %s\n", bookData.Title)
		if author != nil {
			fmt.Printf("  Autor: %s\n", author.Name)
		}
		if publisher != nil {
			fmt.Printf("  Editora: %s\n", publisher.Name)
		}
		successCount++

		// Pequeno delay entre requisições para não sobrecarregar a API
		time.Sleep(500 * time.Millisecond)
	}

	// Resultados finais
	fmt.Println("\n\n========== RESUMO FINAL ==========")
	fmt.Printf("Total de ISBNs processados: %d\n", processedCount)
	fmt.Printf("Sucesso: %d\n", successCount)
	fmt.Printf("Erros: %d\n", errorCount)

	// Estatísticas do banco
	totalBooks, err := db.CountBooks()
	if err == nil {
		fmt.Printf("\nTotal de livros no banco de dados: %d\n", totalBooks)
	}

	fmt.Println("\n✓ Teste concluído!")
}
