package ebooks

import (
	"encoding/json"
)

// Meta information about a book
type BookDetails struct {
	Title   string
	Authors []string
	Year    int
}

func (book *BookDetails) ToJson() []byte {
	bookJson, _ := json.Marshal(book)
	return bookJson
}

type Ebook struct {
	Id    int
	Files map[string]string
	*BookDetails
}

type Library interface {
	// Add a new book to the library
	Add(book *BookDetails) (*Ebook, error)

	// Gets a single book with a given id if it exists
	GetBookById(id int) (*Ebook, error)

	// Gets all books in the library
	GetAll() ([]*Ebook)
}


