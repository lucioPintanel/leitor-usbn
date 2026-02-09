package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"text/template"

	"leitor-usbn/database"
)

var tmpl *template.Template

func main() {
	dbPath := flag.String("db", "../books.db", "caminho para o arquivo sqlite")
	port := flag.Int("port", 8080, "porta HTTP")
	flag.Parse()

	abs, _ := filepath.Abs(*dbPath)
	fmt.Printf("Iniciando web UI — DB: %s\n", abs)

	db, err := database.NewDatabase(*dbPath)
	if err != nil {
		log.Fatalf("erro ao abrir DB: %v", err)
	}
	defer db.Close()

	err = db.InitSchema()
	if err != nil {
		log.Fatalf("erro ao inicializar schema: %v", err)
	}

	// carregar templates (caminho relativo ao workspace)
	tmpl = template.Must(template.ParseFiles("src/web/templates/books.html"))

	// handlers
	http.HandleFunc("/api/books", func(w http.ResponseWriter, r *http.Request) {
		books, err := db.GetBooksWithDetails()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		writeJSON(w, books)
	})

	http.HandleFunc("/ui", func(w http.ResponseWriter, r *http.Request) {
		books, err := db.GetBooksWithDetails()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data := map[string]interface{}{"Books": books}
		tmpl.Execute(w, data)
	})

	// Redirect root to UI
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/ui", http.StatusSeeOther)
	})

	// serve openapi and swagger UI static
	http.HandleFunc("/openapi.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "src/web/static/openapi.json")
	})
	http.Handle("/docs/", http.StripPrefix("/docs/", http.FileServer(http.Dir("src/web/static"))))

	addr := fmt.Sprintf(":%d", *port)
	fmt.Printf("Servidor rodando em http://localhost%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(v)
}
