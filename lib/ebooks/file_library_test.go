package ebooks

import (
	"fmt"
	"testing"

	assert "github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
)

// List of temp folders created during testing to be cleaned up during teardown
var tempFolders []string = []string{}

func TestMain(m *testing.M) {
	defer deleteTempFoldersCreatedDuringTesting()
	m.Run()
}

func TestNewBooksAreAssignedAUniqueId(t *testing.T) {
	library := newLibraryInTempFolder(t)
	id1, err1 := library.Add(aBook("Book1", "mr writer", 2016))
	id2, err2 := library.Add(aBook("Book2", "mrs writer", 2015))

	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.NotEqual(t, id1, id2, "Two books cannot have same id")
}

func TestABookCanBeRetrievedAByIdAfterAdding(t *testing.T) {

	library := newLibraryInTempFolder(t)
	ebook, _ := library.Add(aBook("Book1", "mr writer", 2016))

	libraryBook, err := library.GetBookById(ebook.Id)
	assert.NoError(t, err, "Expected to find a book but did not")
	assert.Equal(t, "Book1", libraryBook.Title)
	assert.Equal(t, "mr writer", libraryBook.Authors[0])
	assert.Equal(t, 2016, libraryBook.Year)
}

func TestAnEmptyLibraryContainsNoBooks(t *testing.T) {
	library := newLibraryInTempFolder(t)
	books := library.GetAll()
	assert.Empty(t, books)
}

func TestALibraryContainsAllBooksAddedToIt(t *testing.T) {
	library := newLibraryInTempFolder(t)
	library.Add(aBook("Book1", "mr writer", 2016))
	library.Add(aBook("Book2", "mrs writer", 2015))

	books := library.GetAll()
	assert.Equal(t, 2, len(books))
}

func TestSaveIndexToDiskSavesAnIndexFileInTheBaseDir_EmptyLib(t *testing.T) {
	library := newLibraryInTempFolder(t)
	err := library.SaveIndexToDisk()
	assert.NoError(t, err, "Error saving index to disk")

	indexJson, err := ioutil.ReadFile(library.fileForIndex())
	assert.NoError(t, err)
	assert.Equal(t, "{}", string(indexJson))
}

func TestSaveIndexToDiskSavesAnIndexFileInTheBaseDir_NonEmptyLib(t *testing.T) {
	library := newLibraryInTempFolder(t)
	library.Add(aBook("Book1", ",mr writer", 2016))
	library.Add(aBook("Book2", "mrs writer", 2015))

	err := library.SaveIndexToDisk()
	assert.NoError(t, err, "Error saving index to disk")

	indexFileName := library.fileForIndex()
	_, err = ioutil.ReadFile(indexFileName)
	assert.NoError(t, err)

	expectedMap := library.indexToBookDetailsJsonMap()
	actualMap, err := library.bookDetailsJsonMapFromFile(indexFileName)
	assert.NoError(t, err)
	assert.Equal(t, expectedMap, actualMap)
}

func aBook(name string, author string, year int) *BookDetails {
	return &BookDetails{
		Title:    name,
		Authors: []string{author},
		Year:    year,
	}
}

func deleteTempFoldersCreatedDuringTesting() {
	for _, folder := range tempFolders {
		err := os.RemoveAll(folder)
		if err != nil {
			fmt.Printf("Error deleting temp dir %v, err=%v\n", folder, err)
		}
	}
}

func createTempFolder(t *testing.T) string {
	folder, err := ioutil.TempDir("", "ebook_tests")
	if err != nil {
		t.Fatalf("Error creating temp dir", err)
	}
	tempFolders = append(tempFolders, folder)
	return folder
}

func newLibraryInTempFolder(t *testing.T) *FileLibrary {
	folder := createTempFolder(t)
	library, err := NewFileLibrary(folder)
	if err != nil {
		t.Fatalf("Error creating library", err)
	}

	return library
}

func containsFile(targetFile string, files []os.FileInfo) bool {
	for _, file := range files {
		if file.Name() == targetFile {
			return true
		}
	}
	return false
}
