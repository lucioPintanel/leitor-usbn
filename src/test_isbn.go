package main

import (
	"fmt"
	"log"
	"leitor-usbn/api"
	"leitor-usbn/database"
)

func main() {
	fmt.Println("=== Teste de Consulta ISBN e População do DB ===\n")

	// ISBN do livro "Clean Code" de Robert C. Martin
	testISBN := "0132350882"

	// 1. Inicializar banco de dados
	fmt.Println("[1] Inicializando banco de dados...")
	db, err := database.NewDatabase("./books.db")
	if err != nil {
		log.Fatalf("Erro ao inicializar banco de dados: %v", err)
	}
	defer db.Close()

	// Criar schema
	err = db.InitSchema()
	if err != nil {
		log.Fatalf("Erro ao criar schema: %v", err)
	}
	fmt.Println("✓ Banco de dados inicializado com sucesso\n")

	// 2. Criar cliente API
	fmt.Println("[2] Inicializando cliente da API OpenLibrary...")
	client := api.NewBookAPIClient("https://openlibrary.org/api/books", 10)
	fmt.Println("✓ Cliente API criado\n")

	// 3. Consultar ISBN na API
	fmt.Printf("[3] Consultando ISBN: %s na API...\n", testISBN)
	apiBook, err := client.GetBookByISBN(testISBN)
	if err != nil {
		log.Fatalf("Erro ao consultar API: %v", err)
	}
	fmt.Println("✓ Livro encontrado na API\n")

	// 4. Converter para formato padrão
	fmt.Println("[4] Convertendo dados da API...")
	bookData := api.ConvertToBookData(apiBook)
	fmt.Println("✓ Dados convertidos\n")

	// 5. Salvar no banco de dados
	fmt.Println("[5] Salvando no banco de dados...")

	// Obter ou criar autor
	author, err := db.GetOrCreateAuthor(bookData.Author)
	if err != nil {
		log.Fatalf("Erro ao salvar autor: %v", err)
	}

	// Obter ou criar editora
	publisher, err := db.GetOrCreatePublisher(bookData.Publisher)
	if err != nil {
		log.Fatalf("Erro ao salvar editora: %v", err)
	}

	// Criar estrutura de livro para salvar
	dbBook := &database.Book{
		ISBN:        bookData.ISBN,
		Title:       bookData.Title,
		AuthorID:    nil,
		PublisherID: nil,
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
		log.Fatalf("Erro ao salvar livro: %v", err)
	}
	fmt.Println("✓ Livro salvo no banco de dados\n")

	// 6. Recuperar e exibir dados
	fmt.Println("[6] Recuperando dados do banco de dados...")
	retrievedBook, err := db.GetBookByISBN(testISBN)
	if err != nil {
		log.Fatalf("Erro ao recuperar livro: %v", err)
	}

	fmt.Println("\n========== RESULTADO FINAL ==========")
	fmt.Printf("ID no DB:     %d\n", retrievedBook.ID)
	fmt.Printf("ISBN:         %s\n", retrievedBook.ISBN)
	fmt.Printf("Título:       %s\n", retrievedBook.Title)
	if retrievedBook.AuthorID != nil {
		fmt.Printf("ID do Autor:  %d\n", *retrievedBook.AuthorID)
	}
	if retrievedBook.PublisherID != nil {
		fmt.Printf("ID Editora:   %d\n", *retrievedBook.PublisherID)
	}
	fmt.Printf("Data Pub.:    %s\n", retrievedBook.PublishDate)
	fmt.Printf("Páginas:      %d\n", retrievedBook.Pages)
	fmt.Printf("Descrição:    %s\n", truncateString(retrievedBook.Description, 100))
	fmt.Printf("URL da Capa:  %s\n", retrievedBook.CoverURL)
	fmt.Printf("Criado em:    %s\n", retrievedBook.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("Atualizado:   %s\n", retrievedBook.UpdatedAt.Format("2006-01-02 15:04:05"))

	// 7. Contar total de livros
	fmt.Println("\n[7] Verificando estatísticas...")
	count, err := db.CountBooks()
	if err != nil {
		log.Fatalf("Erro ao contar livros: %v", err)
	}
	fmt.Printf("Total de livros no DB: %d\n", count)

	// 8. Exibir autor salvo
	if author != nil {
		fmt.Printf("\nAutor salvo: ID=%d, Nome=%s\n", author.ID, author.Name)
	}

	// 9. Exibir editora salva
	if publisher != nil {
		fmt.Printf("Editora salva: ID=%d, Nome=%s\n", publisher.ID, publisher.Name)
	}

	fmt.Println("\n✓ Teste concluído com sucesso!")
}

// truncateString limita o tamanho de uma string para exibição
func truncateString(s string, maxLen int) string {
	if len(s) > maxLen {
		return s[:maxLen] + "..."
	}
	return s
}
