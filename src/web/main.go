package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"text/template"
	"time"

	"leitor-usbn/database"
)

var (
	tmpl *template.Template
	db   *database.Database
)

func main() {
	dbPath := flag.String("db", "../books.db", "caminho para o arquivo sqlite")
	port := flag.Int("port", 8080, "porta HTTP")
	flag.Parse()

	abs, _ := filepath.Abs(*dbPath)
	fmt.Printf("Iniciando web UI — DB: %s\n", abs)

	var err error
	db, err = database.NewDatabase(*dbPath)
	if err != nil {
		log.Fatalf("erro ao abrir DB: %v", err)
	}
	defer db.Close()

	err = db.InitSchema()
	if err != nil {
		log.Fatalf("erro ao inicializar schema: %v", err)
	}

	// carregar templates
	tmpl = template.Must(template.ParseFiles(
		"src/web/templates/books.html",
		"src/web/templates/add-book.html",
	))

	// Redirect root to UI
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/ui", http.StatusSeeOther)
	})

	// Listar livros (GET /ui)
	http.HandleFunc("/ui", handleGetBooks)

	// API: listar livros (GET /api/books) e criar livro (POST /api/books)
	http.HandleFunc("/api/books", handleAPIBooks)

	// Página de adicionar livro (GET /add)
	http.HandleFunc("/add", handleAddBookPage)

	// serve openapi and swagger UI static
	http.HandleFunc("/openapi.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "src/web/static/openapi.json")
	})
	http.Handle("/docs/", http.StripPrefix("/docs/", http.FileServer(http.Dir("src/web/static"))))

	addr := fmt.Sprintf(":%d", *port)
	fmt.Printf("Servidor rodando em http://localhost%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

// handleGetBooks retorna a página HTML com a lista de livros
func handleGetBooks(w http.ResponseWriter, r *http.Request) {
	books, err := db.GetBooksWithDetails()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := map[string]interface{}{"Books": books}
	if err := tmpl.ExecuteTemplate(w, "books.html", data); err != nil {
		fmt.Printf("erro ao renderizar template: %v\n", err)
	}
}

// handleAPIGetBooks retorna lista de livros em JSON
func handleAPIGetBooks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "metodo nao permitido", http.StatusMethodNotAllowed)
		return
	}

	books, err := db.GetBooksWithDetails()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, books)
}

// handleAPIBooks despacha GET (listar) ou POST (criar) para /api/books
func handleAPIBooks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleAPIGetBooks(w, r)
	case http.MethodPost:
		handleAPICreateBook(w, r)
	default:
		http.Error(w, "metodo nao permitido", http.StatusMethodNotAllowed)
	}
}

// handleAPICreateBook recebe JSON e persiste livro no banco
func handleAPICreateBook(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ISBN        string `json:"isbn"`
		Title       string `json:"title"`
		Author      string `json:"author"`
		Publisher   string `json:"publisher"`
		PublishDate string `json:"publish_date"`
		Pages       int    `json:"pages"`
		Description string `json:"description"`
		CoverURL    string `json:"cover_url"`
	}

	// parse JSON
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "erro ao ler JSON: " + err.Error()})
		return
	}

	// validar ISBN
	if input.ISBN == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "ISBN obrigatório"})
		return
	}

	// criar ou obter author/publisher
	var authorID, publisherID *int

	if input.Author != "" {
		author, err := db.GetOrCreateAuthor(input.Author)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "erro ao processar autor: " + err.Error()})
			return
		}
		authorID = &author.ID
	}

	if input.Publisher != "" {
		publisher, err := db.GetOrCreatePublisher(input.Publisher)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "erro ao processar editora: " + err.Error()})
			return
		}
		publisherID = &publisher.ID
	}

	// montar book
	book := &database.Book{
		ISBN:        input.ISBN,
		Title:       input.Title,
		AuthorID:    authorID,
		PublisherID: publisherID,
		PublishDate: input.PublishDate,
		Pages:       input.Pages,
		Description: input.Description,
		CoverURL:    input.CoverURL,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// salvar
	saved, err := db.SaveBook(book)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "erro ao salvar livro: " + err.Error()})
		return
	}

	// retornar livro salvo
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(saved)
}

// handleAddBookPage retorna a página HTML do formulário
func handleAddBookPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "metodo nao permitido", http.StatusMethodNotAllowed)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "add-book.html", nil); err != nil {
		fmt.Printf("erro ao renderizar template: %v\n", err)
		http.Error(w, "erro ao carregar página", http.StatusInternalServerError)
	}
}

// writeJSON escreve resposta JSON
func writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(data)
}
