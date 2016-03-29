package webservice

import (
	"html/template"
	"net/http"
	"io/ioutil"
	"path/filepath"

	"github.com/stephenhenderson/ebooklib/lib/ebooks"
	. "github.com/stephenhenderson/ebooklib/lib/logging"
	"strings"
	"fmt"
	"strconv"
)

// NewEbookWebService initialises a new webservice with the given library
// and html template directory, returns error if there is any error loading
// the templates
func NewEbookWebService(library ebooks.Library, templateDir string) (*EbookWebService, error) {
	templates, err := loadTemplates(templateDir)
	if err != nil {
		return nil, err
	}
	return &EbookWebService{library: library, templates: templates}, nil
}

func loadTemplates(templateDir string) (map[string]*template.Template, error) {
	files, err := ioutil.ReadDir(templateDir)
	if err != nil {
		return nil, err
	}

	templateMap := make(map[string]*template.Template)
	for _, file := range(files) {
		fileName := file.Name()
		if strings.HasSuffix(fileName, ".html") {
			templatePath := filepath.Join(templateDir, fileName)
			template, err := template.New(fileName).ParseFiles(templatePath)
			if err != nil {
				return nil, err
			}
			templateMap[fileName] = template
		}
	}
	return templateMap, nil
}

type EbookWebService struct {
	library ebooks.Library
	templates map[string]*template.Template
}

// StartService starts the webserver listening on the given host
func (webservice *EbookWebService) StartService(host string) {
	Logger.Printf("Starting webservice on %s", host)
	http.HandleFunc("/", webservice.listAllHandler)
	http.HandleFunc("/add_book.html", webservice.addBookFormHandler)
	http.HandleFunc("/addBook", webservice.addBookHandler)
	http.ListenAndServe(host, nil)
}

func (webservice *EbookWebService) addBookHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Fprintf(w, "Error parsing form %v", err)
		return
	}

	title := r.PostFormValue("title")
	authors := strings.Split(r.PostFormValue("authors"), ",")
	yearStr := r.PostFormValue("year")
	year := 0
	if yearStr != "" {
		year, err = strconv.Atoi(yearStr)
		if err != nil {
			fmt.Fprintf(w, "Invalid year '%s', err=%v", yearStr, err)
		}
	}

	bookDetails := &ebooks.BookDetails{
		Title: title,
		Authors: authors,
		Year: year,
	}
	webservice.library.Add(bookDetails)
	http.Redirect(w, r, "/", http.StatusFound)
}

func (webservice *EbookWebService) addBookFormHandler(w http.ResponseWriter, r *http.Request) {
	template := webservice.templates["add_book.html"]
	template.Execute(w, nil)
}

func (webservice *EbookWebService) listAllHandler(w http.ResponseWriter, r *http.Request) {
	books := webservice.library.GetAll()
	template := webservice.templates["index.html"]
	template.Execute(w, books)
}
