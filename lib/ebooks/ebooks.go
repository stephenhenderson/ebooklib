package ebooks

import (
	"encoding/json"
	"github.com/stephenhenderson/ebooklib/lib/utils"
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

func (book *BookDetails) Equals(anotherBook *BookDetails) bool {
	if book.Title != anotherBook.Title {
		return false
	}
	if book.Year != anotherBook.Year {
		return false
	}
	if !utils.StringSliceEquals(book.Authors, anotherBook.Authors) {
		return false
	}
	if !utils.StringSliceEquals(book.Tags, anotherBook.Tags) {
		return false
	}
	return true
}

// A book in a library with a unique id, details and a collection of
// files.
type Ebook struct {
	// Unique id for this book in the library
	ID    int

	// Files associated with this book (typically the ebook file(s) but
	// could be supporting code, etc)
	Files map[string]string

	// Optional image of the book
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
