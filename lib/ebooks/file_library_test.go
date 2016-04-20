package ebooks

import (
	"testing"
	"encoding/json"
	"io/ioutil"
	"reflect"

	"github.com/stephenhenderson/ebooklib/lib/testutils/assert"
	"github.com/stephenhenderson/ebooklib/lib/testutils"
)

func TestMain(m *testing.M) {
	defer testutils.DeleteTempDirsCreatedDuringTesting()
	m.Run()
}

func TestNewBooksAreAssignedAUniqueId(t *testing.T) {
	library := newLibraryInTempFolder(t)
	id1, err1 := library.Add(aBook("Book1", "mr writer", 2016), emptyFileMap())
	id2, err2 := library.Add(aBook("Book2", "mrs writer", 2015), emptyFileMap())

	assert.NoError(t, err1)
	assert.NoError(t, err2)
	if id1 == id2 {
		t.Fatalf("Two books cannot have the same id. id=%v", id1)
	}
}

func TestABookCanBeRetrievedAByIdAfterAdding(t *testing.T) {

	library := newLibraryInTempFolder(t)
	ebook, _ := library.Add(aBook("Book1", "mr writer", 2016), emptyFileMap())

	libraryBook, err := library.GetBookByID(ebook.ID)
	assert.NoError(t, err, "Expected to find a book but did not")

	libraryDetails := libraryBook.BookDetails
	expectedDetails := ebook.BookDetails
	if libraryDetails != expectedDetails {
		t.Fatalf("Retrieved book %v is not same as added book: %v", libraryDetails, expectedDetails)
	}
}

func TestAnEmptyLibraryContainsNoBooks(t *testing.T) {
	library := newLibraryInTempFolder(t)
	books := library.GetAll()
	if len(books) != 0 {
		t.Fatalf("Empty library should have 0 books but has %v", len(books))
	}
}

func TestALibraryContainsAllBooksAddedToIt(t *testing.T) {
	library := newLibraryInTempFolder(t)
	library.Add(aBook("Book1", "mr writer", 2016), emptyFileMap())
	library.Add(aBook("Book2", "mrs writer", 2015), emptyFileMap())

	books := library.GetAll()
	if len(books) != 2 {
		t.Fatalf("Library should contain all books added to it. Expected 2 but found %v", len(books))
	}
}

func TestSaveIndexToDiskSavesAnIndexFileInTheBaseDir_EmptyLib(t *testing.T) {
	library := newLibraryInTempFolder(t)
	err := library.SaveIndexToDisk()
	assert.NoError(t, err, "Error saving index to disk")

	indexJson, err := ioutil.ReadFile(library.fileForIndex())
	assert.NoError(t, err)
	if string(indexJson) != "{}" {
		t.Fatalf("Expected empty index '{}' but found '%v'", string(indexJson))
	}
}

func TestSaveIndexToDiskSavesAnIndexFileInTheBaseDir_NonEmptyLib(t *testing.T) {
	library := newLibraryInTempFolder(t)
	library.Add(aBook("Book1", ",mr writer", 2016), emptyFileMap())
	library.Add(aBook("Book2", "mrs writer", 2015), emptyFileMap())

	err := library.SaveIndexToDisk()
	assert.NoError(t, err, "Error saving index to disk")

	indexFileName := library.fileForIndex()
	_, err = ioutil.ReadFile(indexFileName)
	assert.NoError(t, err)

	expectedMap := library.indexToBookDetailsJsonMap()
	actualMap, err := library.bookDetailsJsonMapFromFile(indexFileName)
	assert.NoError(t, err)

	if len(actualMap) != len(expectedMap) {
		t.Fatalf("Index written to disk '%v' does not match expected '%v'", actualMap, expectedMap)
	}
	for id, bookDetails := range actualMap {
		expectedDetails, found := expectedMap[id]
		if !found {
			t.Fatalf("Found unexpected book id=%v, book=%v", id, bookDetails)
		}
		if !reflect.DeepEqual(bookDetails, expectedDetails) {
			t.Fatalf("Saved book does not match original, saved='%v', actual='%v'", bookDetails, expectedDetails)
		}
	}
}

func TestSavingABookWithAFile(t *testing.T) {
	library := newLibraryInTempFolder(t)
	bookFiles := make(map[string][]byte)
	bookFiles["file1.json"] = aJsonFile()
	book, err := library.Add(aBook("book1", "mr writer", 2016), bookFiles)

	assert.NoError(t, err)

	_, found := book.Files["file1.json"]
	if !found {
		t.Fatalf("File associated with book not saved")
	}
}

func TestAFileCanBeDeletedFromABook(t *testing.T) {
	fileName := "file1.json"
	library := newLibraryInTempFolder(t)
	bookFiles := make(map[string][]byte)
	bookData := aJsonFile()
	bookFiles[fileName] = bookData
	book, _ := library.Add(aBook("book1", "mr writer", 2016), bookFiles)

	err := library.DeleteFileFromBook(fileName, book.ID)
	assert.NoError(t, err)

	// check the file is no longer on disk
	fileLocation := library.fullPathToBookFile(fileName, book.ID)
	_, err = ioutil.ReadFile(fileLocation)
	if err == nil {
		t.Fatalf("file %s was not deleted from the file system", fileName)
	}

	// check the file is no longer referenced by the library
	book, err = library.GetBookByID(book.ID)
	assert.NoError(t, err)
	_, found := book.Files[fileName]
	if found {
		t.Fatalf("file %s still referenced in library after deletion", fileName)
	}
}

func TestReturnsAnErrorTryingToDeleteAFileWhichDoesNotExist(t *testing.T) {
	library := newLibraryInTempFolder(t)
	book, _ := library.Add(aBook("book1", "mr writer", 2016), make(map[string][]byte))

	err := library.DeleteFileFromBook("a_file_which_is_not_there", book.ID)
	if err == nil {
		t.Fatal("No error was returned trying to delete a nonexistent file")
	}
}

func TestReturnsAnErrorTryingToDeleteAFileFromABookWhichDoesNotExist(t *testing.T) {
	library := newLibraryInTempFolder(t)
	err := library.DeleteFileFromBook("a_file_which_is_not_there", 123)
	if err == nil {
		t.Fatal("No error was returned trying to delete a nonexistent file")
	}
}

func aJsonFile() []byte {
	someData := make(map[string]string)
	someData["key1"] = "value1"
	someData["key2"] = "value2"
	data, _ := json.Marshal(someData)
	return data
}

func emptyFileMap() map[string][]byte {
	return make(map[string][]byte)
}

func aBook(name string, author string, year int) *BookDetails {
	return &BookDetails{
		Title:   name,
		Authors: []string{author},
		Year:    year,
	}
}

func newLibraryInTempFolder(t *testing.T) *FileLibrary {
	folder := testutils.CreateTempDir(t)
	library, err := NewFileLibrary(folder)
	if err != nil {
		t.Fatalf("Error creating library", err)
	}

	return library
}
