package webservice

import (
	"testing"
	"github.com/stephenhenderson/ebooklib/lib/ebooks"
	"github.com/stephenhenderson/ebooklib/lib/testutils"
)

func TestMain(m *testing.M) {
	defer testutils.DeleteTempDirsCreatedDuringTesting()
	m.Run()
}

func TestNewEbookWebServiceReturnsErrorIfTemplatesDirDoesNotExist(t *testing.T) {
	_, err := NewEbookWebService(&ebooks.FileLibrary{}, "/non-existant-dir")
	if err == nil {
		t.Fatal("Expected error creating webservice with non-existant template dir")
	}
}

func TestNewEbookWebServiceReturnsErrorIfTemplatesDirIsEmpty(t *testing.T) {
	emptyTempDir := testutils.CreateTempDir(t)
	_, err := NewEbookWebService(&ebooks.FileLibrary{}, emptyTempDir)
	if err == nil {
		t.Fatal("Expected error creating webservice with non-existant template dir")
	}
}
