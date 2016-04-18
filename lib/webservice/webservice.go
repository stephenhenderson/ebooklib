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

const (
	addBookTemplate = "add_book.html"
	indexTemplate = "index.html"
	viewBookTemplate = "view_book.html"
)

// NewEbookWebService initialises a new webservice with the given library
// and html template directory, returns error if there is any error loading
// the templates
func NewEbookWebService(library *ebooks.FileLibrary, templateDir string) (*EbookWebService, error) {
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

	err = validateTemplates(templateMap)
	return templateMap, err
}

func validateTemplates(templateMap map[string]*template.Template) error {
	expectedTemplates := []string{addBookTemplate, viewBookTemplate, indexTemplate}
	for _, template := range(expectedTemplates) {
		_, found := templateMap[template]
		if !found {
			return fmt.Errorf("missing required template %s", template)
		}
	}
	return nil
}

type EbookWebService struct {
	library *ebooks.FileLibrary
	templates map[string]*template.Template
}

// StartService starts the webserver listening on the given host
func (webservice *EbookWebService) StartService(host string) {
	Logger.Printf("Starting webservice on %s", host)
	http.HandleFunc("/", webservice.listAllHandler)
	http.HandleFunc("/" + addBookTemplate, webservice.addBookFormHandler)
	http.HandleFunc("/" + viewBookTemplate, webservice.viewBookHandler)

	http.Handle("/download_book/", http.StripPrefix("/download_book/", http.FileServer(http.Dir(webservice.library.BaseDir))))
	http.HandleFunc("/addBook", webservice.addBookHandler)

	http.ListenAndServe(host, nil)
}

func (webservice *EbookWebService) viewBookHandler(w http.ResponseWriter, r *http.Request) {
	bookID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "No book with this id", http.StatusNotFound)
		return
	}

	book, err := webservice.library.GetBookByID(bookID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = webservice.templates["view_book.html"].Execute(w, book)
	if err != nil {
		http.Error(w, "No book with this id", http.StatusNotFound)
		return
	}
}

func (webservice *EbookWebService) addBookHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(100000)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	title := r.MultipartForm.Value["title"][0]
	authors := strings.Split(r.MultipartForm.Value["authors"][0], ",")
	yearStr := r.MultipartForm.Value["year"][0]
	year := 0
	if yearStr != "" {
		year, err = strconv.Atoi(yearStr)
		if err != nil {
			fmt.Fprintf(w, "Invalid year '%s', err=%v", yearStr, err)
		}
	}

	bookFiles := make(map[string][]byte)
	fileHeaders := r.MultipartForm.File["files"]
	Logger.Printf("File headers: %v", fileHeaders)
	for _, fileHeader := range fileHeaders {
		//for each fileheader, get a handle to the actual file

		file, err := fileHeader.Open()
		defer file.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data, err := ioutil.ReadAll(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		bookFiles[fileHeader.Filename] = data
	}

	bookDetails := &ebooks.BookDetails{
		Title: title,
		Authors: authors,
		Year: year,
	}

	webservice.library.Add(bookDetails, bookFiles)
	http.Redirect(w, r, "/", http.StatusFound)
}

func (webservice *EbookWebService) addBookFormHandler(w http.ResponseWriter, r *http.Request) {
	template := webservice.templates[addBookTemplate]
	template.Execute(w, nil)
}

func (webservice *EbookWebService) listAllHandler(w http.ResponseWriter, r *http.Request) {
	books := webservice.library.GetAll()
	template := webservice.templates[indexTemplate]
	err := template.Execute(w, books)
	if err != nil {
		fmt.Fprintf(w, "Unexpected error:%v", err)
	}
}
