package api

import "fmt"

// BookData contém os dados padronizados do livro retornados pela API
type BookData struct {
	ISBN        string
	Title       string
	Author      string
	Publisher   string
	PublishDate string
	Pages       int
	Description string
	CoverURL    string
}

// ConvertToBookData converte a resposta da OpenLibrary para um formato padronizado
func ConvertToBookData(apiResponse *OpenLibraryResponse) *BookData {
	author := ""
	if len(apiResponse.Authors) > 0 {
		author = apiResponse.Authors[0].Name
	}

	publisher := ""
	if len(apiResponse.Publishers) > 0 {
		publisher = apiResponse.Publishers[0].Name
	}

	description := ""
	if apiResponse.Description.Value != "" {
		description = apiResponse.Description.Value
	}

	coverURL := ""
	if len(apiResponse.Covers) > 0 {
		coverID := apiResponse.Covers[0]
		coverURL = "https://covers.openlibrary.org/b/id/" + fmt.Sprintf("%d-M.jpg", coverID)
	}

	return &BookData{
		ISBN:        apiResponse.ISBN,
		Title:       apiResponse.Title,
		Author:      author,
		Publisher:   publisher,
		PublishDate: apiResponse.PublishDate,
		Pages:       apiResponse.NumberOfPages,
		Description: description,
		CoverURL:    coverURL,
	}
}
