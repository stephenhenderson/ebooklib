package ebooks

import (
	"encoding/json"
	"path/filepath"
	"os"
	"io/ioutil"
	"strconv"
	"errors"
)

var BookNotFound = errors.New("Book not found")

type BookDetails struct {
	Name    string
	Authors []string
	Year    int
}

func (book *BookDetails) ToJson() []byte {
	bookJson, _ := json.Marshal(book)
	return bookJson
}

type Ebook struct {
	Id int
	Files map[string]string
	*BookDetails
}

type Library interface {
	Add(book *BookDetails) (*Ebook, error)
	GetBookById(id int) (book *Ebook, found bool)
}

type FileLibrary struct {
	maxId int
	baseDir string
}

func (lib *FileLibrary) Add(bookDetails *BookDetails) (ebook *Ebook, err error) {
	lib.maxId += 1
	ebook = &Ebook{lib.maxId, make(map[string]string), bookDetails}
	err = lib.createNewBookFiles(ebook)
	return ebook, err
}

func (lib *FileLibrary) GetBookById(id int) (*Ebook, error) {
	// does folder exist
	details := &BookDetails{}
	detailsFile := lib.fileForBook(id)
	detailsJson, err := ioutil.ReadFile(detailsFile)
	if err != nil {
		return nil, BookNotFound
	}

	err = json.Unmarshal(detailsJson, details)
	if err != nil {
		return nil, err
	}

	book := &Ebook{id, make(map[string]string), details}
	return book, nil
}

func (lib *FileLibrary) fileForBook(id int) string {
	return filepath.Join(lib.folderForBook(id), "details.json")
}

func (lib *FileLibrary) folderForBook(id int) string {
	return filepath.Join(lib.baseDir, strconv.Itoa(id))
}

func (lib *FileLibrary) createNewBookFiles(book *Ebook) (err error) {
	// Create directory structure
	bookFolder := lib.folderForBook(book.Id)
	filesFolder := filepath.Join(bookFolder, "files")
	mkDirs(bookFolder, filesFolder)
	if err != nil {
		return err
	}

	// Write the book details descriptor
	bookDetailsFile := lib.fileForBook(book.Id)
	jsonDescriptor := book.ToJson()
	err = ioutil.WriteFile(bookDetailsFile, jsonDescriptor, 0700)
	return err
}

func mkDirs(dirs... string) (err error) {
	for _, dir := range(dirs) {
		err = os.MkdirAll(dir, 0700)
		if err != nil {
			return err
		}
	}
	return
}

func NewFileLibrary(baseDir string) (*FileLibrary, error) {
	err := createDirIfNotExists(baseDir)
	return &FileLibrary{baseDir: baseDir}, err
}

func createDirIfNotExists(dir string) (err error) {
	if _, err := os.Stat("dir"); err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(dir, 0700)
		}
	}
	return err
}