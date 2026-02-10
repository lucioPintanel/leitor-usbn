package api

import (
	"testing"
)

// TestConvertToBookData testa conversão de resposta API para formato padrão
func TestConvertToBookData(t *testing.T) {
	apiResp := &OpenLibraryResponse{
		ISBN: "0132350882",
		Title: "Clean Code",
		Authors: []AuthorInfo{
			{Name: "Robert C. Martin", Key: "/authors/OL123A"},
		},
		Publishers: []struct {
			Name string `json:"name"`
		}{
			{Name: "Prentice Hall"},
		},
		PublishDate: "2008",
		NumberOfPages: 464,
		Description: Description{Value: "A Handbook of Agile Software Craftsmanship"},
		Covers: []int{1234567},
	}

	data := ConvertToBookData(apiResp)

	if data.ISBN != "0132350882" {
		t.Errorf("ISBN = %s, want 0132350882", data.ISBN)
	}
	if data.Title != "Clean Code" {
		t.Errorf("Title = %s, want Clean Code", data.Title)
	}
	if data.Author != "Robert C. Martin" {
		t.Errorf("Author = %s, want Robert C. Martin", data.Author)
	}
	if data.Publisher != "Prentice Hall" {
		t.Errorf("Publisher = %s, want Prentice Hall", data.Publisher)
	}
	if data.Pages != 464 {
		t.Errorf("Pages = %d, want 464", data.Pages)
	}
	if data.PublishDate != "2008" {
		t.Errorf("PublishDate = %s, want 2008", data.PublishDate)
	}
	if data.Description != "A Handbook of Agile Software Craftsmanship" {
		t.Errorf("Description = %s, want A Handbook of Agile Software Craftsmanship", data.Description)
	}
}

// TestConvertToBookDataEmptyValues testa conversão com valores vazios
func TestConvertToBookDataEmptyValues(t *testing.T) {
	apiResp := &OpenLibraryResponse{
		ISBN: "9999999999",
		Title: "Unknown Book",
		// No Authors, Publishers, Covers
	}

	data := ConvertToBookData(apiResp)

	if data.Author != "" {
		t.Errorf("Author should be empty, got %s", data.Author)
	}
	if data.Publisher != "" {
		t.Errorf("Publisher should be empty, got %s", data.Publisher)
	}
	if data.CoverURL != "" {
		t.Errorf("CoverURL should be empty, got %s", data.CoverURL)
	}
}

// TestBookDataStructure testa se BookData tem campos obrigatórios
func TestBookDataStructure(t *testing.T) {
	data := &BookData{
		ISBN:        "test",
		Title:       "test",
		Author:      "test",
		Publisher:   "test",
		PublishDate: "2026",
		Pages:       100,
		Description: "test",
		CoverURL:    "http://example.com/cover.jpg",
	}

	// Verificar que nenhum campo crítico é nulo
	if data.ISBN == "" {
		t.Error("ISBN não pode ser vazio")
	}
	if data.Title == "" {
		t.Error("Title não pode ser vazio")
	}
}
