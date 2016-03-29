package main

import (
	. "github.com/stephenhenderson/ebooklib/lib/logging"
	"github.com/stephenhenderson/ebooklib/lib/webservice"
	"github.com/stephenhenderson/ebooklib/lib/ebooks"
)

func main() {
	baseDir := "/tmp/mylibrary"
	templateDir := "/Users/shenderson/workspace/go_paths/ebooklib/src/github.com/stephenhenderson/ebooklib/templates"
	library, err := ebooks.NewFileLibrary(baseDir)

	if err != nil {
		Logger.Fatalf("Encountered fatal error %v", err)
	}

	webservice, err := webservice.NewEbookWebService(library, templateDir)
	if err != nil {
		Logger.Fatalf("Error loading html templates from %s\nerr:\n%v",
			templateDir, err)
	}
	webservice.StartService("localhost:8080")
}
