package ebooks

import (
	"encoding/json"
)

// Meta information about a book
type BookDetails struct {
	Title   string
	Authors []string
	Year    int
	Tags    []string
}

func (book *BookDetails) ToJson() []byte {
	bookJson, _ := json.Marshal(book)
	return bookJson
}

type Ebook struct {
	ID    int
	Files map[string]string
	Image string
	*BookDetails
}

type Library interface {
	// Add a new book to the library
	Add(book *BookDetails, files map[string][]byte) (*Ebook, error)

	// Gets a single book with a given id if it exists
	GetBookByID(id int) (*Ebook, error)

	// Gets all books in the library
	GetAll() []*Ebook
}
