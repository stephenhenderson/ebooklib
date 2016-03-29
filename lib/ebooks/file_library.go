package ebooks

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

var BookNotFound = errors.New("Book not found")

const (
	IndexFileName = "index.json"
	DetailsFileName = "details.json"
)

func NewFileLibrary(baseDir string) (*FileLibrary, error) {
	err := createDirIfNotExists(baseDir)
	if err != nil {
		return nil, err
	}

	index := make(map[int]*Ebook)
	lib := &FileLibrary{baseDir: baseDir, index: index}

	existingIndexFile := lib.fileForIndex()
	if _, err := os.Stat(existingIndexFile); os.IsNotExist(err) {
		// no existing index - return empty library
		return lib, nil
	}

	// load existing library
	err = lib.loadIndexFromFile(existingIndexFile)
	if err != nil {
		return nil, err
	}
	return lib, nil
}


// A library where ebook details are persisted to the local file system
type FileLibrary struct {

	// Counter tracking the largest book id currently in the library
	maxId   int

	// All books currently in the library indexed by id
	index   map[int]*Ebook

	// Base directory where the library contents are stored
	baseDir string
}

func (lib *FileLibrary) Add(bookDetails *BookDetails) (*Ebook, error) {
	lib.maxId += 1
	ebook := &Ebook{lib.maxId, make(map[string]string), bookDetails}
	err := lib.createNewBookFiles(ebook)
	if err != nil {
		return nil, err
	}
	lib.index[ebook.Id] = ebook
	lib.SaveIndexToDisk()
	return ebook, err
}

func (lib *FileLibrary) GetBookById(id int) (*Ebook, error) {
	book, found := lib.index[id]
	if !found {
		return nil, BookNotFound
	}
	return book, nil
}

func (lib *FileLibrary) GetAll() ([]*Ebook) {
	numBooks := len(lib.index)
	books := make([]*Ebook, 0, numBooks)
	for _, book := range lib.index {
		books = append(books, book)
	}
	return books
}

func (lib *FileLibrary) SaveIndexToDisk() error {
	indexFileName := lib.fileForIndex()
	bookDetailsMap := lib.indexToBookDetailsJsonMap()

	jsonIndex, err := json.MarshalIndent(bookDetailsMap,"", " ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(indexFileName, jsonIndex, 0700)
	return err
}

func (lib *FileLibrary) indexToBookDetailsJsonMap() map[string]*BookDetails {
	bookDetailsMap := make(map[string]*BookDetails)
	for id, book := range(lib.index) {
		bookDetailsMap[strconv.Itoa(id)] = book.BookDetails
	}
	return bookDetailsMap
}

func (lib *FileLibrary) loadIndexFromFile(file string) error {
	bookDetailsJsonMap, err := lib.bookDetailsJsonMapFromFile(file)
	if err != nil {
		return err
	}
	index := make(map[int]*Ebook)
	for idStr, bookDetails := range bookDetailsJsonMap {
		id64, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return err
		}
		id := int(id64)
		index[id] = &Ebook{id, make(map[string]string), bookDetails}
	}
	lib.index = index
	return nil
}

func (lib *FileLibrary) bookDetailsJsonMapFromFile(file string) (map[string]*BookDetails, error) {
	jsonBytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	detailsMap := make(map[string]*BookDetails)
	err = json.Unmarshal(jsonBytes, &detailsMap)
	if err != nil {
		return nil, err
	}
	return detailsMap, nil
}

func (lib *FileLibrary) loadBookFromDisk(id int) (*Ebook, error) {
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

func (lib *FileLibrary) fileForIndex() string {
	return filepath.Join(lib.baseDir, IndexFileName)
}

func (lib *FileLibrary) fileForBook(id int) string {
	return filepath.Join(lib.folderForBook(id), DetailsFileName)
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

func mkDirs(dirs ...string) (err error) {
	for _, dir := range dirs {
		err = os.MkdirAll(dir, 0700)
		if err != nil {
			return err
		}
	}
	return
}

func createDirIfNotExists(dir string) (err error) {
	if _, err := os.Stat("dir"); err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(dir, 0700)
		}
	}
	return err
}