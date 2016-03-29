package webservice

import (
	"fmt"
	"net/http"

	"github.com/stephenhenderson/ebooklib/lib/ebooks"
	. "github.com/stephenhenderson/ebooklib/lib/logging"
)

func NewEbookWebService(library ebooks.Library) *EbookWebService{
	return &EbookWebService{library: library}
}

type EbookWebService struct {
	library ebooks.Library
}

// StartService starts the webserver listening on the given host
func (webservice *EbookWebService) StartService(host string) {
	Logger.Printf("Starting webservice on %s", host)
	http.HandleFunc("/", webservice.listAllHandler)
	http.ListenAndServe(host, nil)
}

func (webservice *EbookWebService) listAllHandler(w http.ResponseWriter, r *http.Request) {
	books := webservice.library.GetAll()
	for _, book := range(books) {
		fmt.Fprintf(w, "Book: %s\n", string(book.ToJson()))
	}
}


