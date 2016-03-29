package main

import (
	"fmt"
	"github.com/stephenhenderson/ebooklib/lib/webservice"
	"github.com/stephenhenderson/ebooklib/lib/ebooks"
	"log"
)

func main() {
	fmt.Println("Starting webserver")
	baseDir := "/tmp/mylibrary"
	library, err := ebooks.NewFileLibrary(baseDir)
	if err != nil {
		log.Fatalf("Encountered fatal error %v", err)
	}

	webservice := webservice.NewEbookWebService(library)
	webservice.StartService("localhost:8080")
}
