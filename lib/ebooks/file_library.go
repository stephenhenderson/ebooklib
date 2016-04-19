package ebooks

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	. "github.com/stephenhenderson/ebooklib/lib/logging"
	"fmt"
	"github.com/FiloSottile/gvt/fileutils"
)

var BookNotFound = errors.New("Book not found")

const (
	IndexFileName = "index.json"
)

// NewFileLibrary opens a file library in the given directory. If existing
// library files are found they are loaded otherwise a new empty library
// is created. If the directory does not exist we attempt to create it.
func NewFileLibrary(baseDir string) (*FileLibrary, error) {
	Logger.Printf("Opening library in %s\n", baseDir)
	err := createDirIfNotExists(baseDir)
	if err != nil {
		return nil, err
	}

	index := make(map[int]*Ebook)
	lib := &FileLibrary{BaseDir: baseDir, index: index}

	existingIndexFile := lib.fileForIndex()
	if _, err := os.Stat(existingIndexFile); os.IsNotExist(err) {
		Logger.Println("No existing index found, creating emptry library")
		return lib, nil
	}

	// load existing library
	Logger.Println("Found existing index file, loading...")
	err = lib.loadIndexFromFile(existingIndexFile)
	if err != nil {
		return nil, err
	}
	Logger.Printf("Loaded library with %v books\n", len(lib.index))
	return lib, nil
}


// A library where ebook details are persisted to the local file system
type FileLibrary struct {

	// Counter tracking the largest book id currently in the library
	maxID   int

	// All books currently in the library indexed by id
	index   map[int]*Ebook

	// Base directory where the library contents are stored
	BaseDir string
}

func (lib *FileLibrary) Add(bookDetails *BookDetails, files map[string][]byte) (*Ebook, error) {
	var err error
	lib.maxID += 1
	ebook := &Ebook{lib.maxID, make(map[string]string), "", bookDetails}
	if err = lib.createNewBookFiles(ebook); err != nil {
		return nil, err
	}

	for fileName, data := range(files) {
		Logger.Printf("Adding files for book=%v, file=%v", bookDetails, fileName)
		if err := lib.AddFileToBook(ebook, fileName, data); err != nil {
			return nil, err
		}
	}

	lib.index[ebook.ID] = ebook
	err = lib.SaveIndexToDisk()
	return ebook, err
}

func (lib *FileLibrary) AddFileToBook(book *Ebook, name string, data []byte) error {
	filePath := lib.fullPathToBookFile(name, book.ID)
	if err := ioutil.WriteFile(filePath, data, 0700); err != nil {
		return err
	}

	// update map with path of file
	book.Files[name] = lib.relativePathToBookFile(name, book.ID)
	return nil
}

func (lib *FileLibrary) DeleteFileFromBook(fileName string, bookID int) error {
	book, found := lib.index[bookID]
	if !found {
		return fmt.Errorf("cannot delete file %s from book with id=%d, no book found with that id", fileName, bookID)
	}

	_, exists := book.Files[fileName]
	if !exists {
		return fmt.Errorf("no file with name %s found for book with id=%d", fileName, bookID)
	}

	err := fileutils.RemoveAll(lib.fullPathToBookFile(fileName, bookID))
	delete(book.Files, fileName)
	return err
}

func (lib *FileLibrary) GetBookByID(id int) (*Ebook, error) {
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
	maxId := 0
	for idStr, bookDetails := range bookDetailsJsonMap {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return err
		}

		book := &Ebook{id, make(map[string]string), "", bookDetails}
		err = lib.loadFilesForBook(book)
		if err != nil {
			return err
		}

		index[id] = book
		if id > maxId {
			maxId = id
		}
	}
	lib.index = index
	lib.maxID = maxId
	return nil
}

func (lib *FileLibrary) loadFilesForBook(book *Ebook) error {
	filesPath := filepath.Join(lib.folderForBook(book.ID), "files")
	files, err := ioutil.ReadDir(filesPath)
	if err != nil {
		return err
	}
	for _, file := range(files) {
		fileName := file.Name()
		book.Files[fileName] = lib.relativePathToBookFile(fileName, book.ID)
	}
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

func (lib *FileLibrary) fileForIndex() string {
	return filepath.Join(lib.BaseDir, IndexFileName)
}

func (lib *FileLibrary) folderForBook(id int) string {
	return filepath.Join(lib.BaseDir, strconv.Itoa(id))
}

func (lib *FileLibrary) relativePathToBookFile(fileName string, bookID int) string {
	return filepath.Join(strconv.Itoa(bookID), "files", fileName)
}

func (lib *FileLibrary) fullPathToBookFile(fileName string, bookID int) string {
	return filepath.Join(lib.folderForBook(bookID), "files", fileName)
}

func (lib *FileLibrary) createNewBookFiles(book *Ebook) error {
	// Create directory structure
	bookFolder := lib.folderForBook(book.ID)
	filesFolder := filepath.Join(bookFolder, "files")
	err := mkDirs(bookFolder, filesFolder)
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