package webservice

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"path/filepath"
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

func TestAddBookParseBookDetailsFromForm(t *testing.T) {
	webservice := newWebserviceWithEmptyLibrary(t)
	ts := httptest.NewServer(http.HandlerFunc(webservice.addBookHandler))
	defer ts.Close()

	request := newAddBookRequest(ts.URL, map[string]string{
		"title":   "Title",
		"authors": "mr writer,mrs writer",
		"year":    "2016",
		"tags":    "tag1,tag2",
	}, aJsonFile(t), t)

	skipRedirectsErr := errors.New("Don't follow redirects")
	client := &http.Client{}

	// Don't follow redirects
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error { return skipRedirectsErr }
	resp, _ := client.Do(request)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusFound {
		t.Fatalf("Expected status code %d but got %s", http.StatusFound, resp.Status)
	}

	allBooks := webservice.library.GetAll()
	if len(allBooks) != 1 {
		t.Fatalf("Expected only 1 book in library but found %d", len(allBooks))
	}
	book := allBooks[0]
	expectedBookDetails := &ebooks.BookDetails{
		Authors: []string{"mr writer","mrs writer"},
		Tags: []string{"tag1","tag2"},
		Title: "Title",
		Year: 2016,
	}
	if !book.BookDetails.Equals(expectedBookDetails) {
		t.Fatalf("Retrieved book %v is not same as added book: %v", book.BookDetails, expectedBookDetails)
	}
}



func newWebserviceWithEmptyLibrary(t *testing.T) *EbookWebService {
	library, err := ebooks.NewFileLibrary(testutils.CreateTempDir(t))
	if err != nil {
		t.Fatalf("Error creating new library %v")
	}
	webservice, err := NewEbookWebService(library, "../../templates/")
	if err != nil {
		t.Fatalf("Error creating new webservice %v", err)
	}
	return webservice
}

// Creates a multi-part form request with a single file and given form parameters
func newAddBookRequest(uri string, bookDetails map[string]string, filePath string, t *testing.T) *http.Request {
	body, contentType := createAddBookMultiPartFormBody(bookDetails, filePath, t)
	req, err := http.NewRequest("POST", uri, body)
	if err != nil {
		t.Errorf("Error creating local book file %v", err)
	}
	req.Header.Add("Content-Type", contentType)
	return req
}

func createAddBookMultiPartFormBody(bookDetails map[string]string, filePath string, t *testing.T) (body *bytes.Buffer, contentType string) {
	body = &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	contentType = writer.FormDataContentType()

	addBookFileToMultiPartWriter(writer, filePath, t)
	for key, val := range bookDetails {
		_ = writer.WriteField(key, val)
	}

	// need to close the writer before sending so it adds the boundary
	// line to the output.
	err := writer.Close()
	if err != nil {
		t.Errorf("Error creating local book file %v", err)
	}
	return
}

func addBookFileToMultiPartWriter(writer *multipart.Writer, filePath string, t *testing.T) {
	file, err := os.Open(filePath)
	if err != nil {
		t.Errorf("Error creating local book file %v", err)
	}
	defer file.Close()

	part, err := writer.CreateFormFile("files", filepath.Base(filePath))
	if err != nil {
		t.Errorf("Error creating local book file %v", err)
	}
	_, err = io.Copy(part, file)
}

func aJsonFile(t *testing.T) string {
	tmpDir := testutils.CreateTempDir(t)
	tmpFileName := path.Join(tmpDir, "myfile.json")

	someData := make(map[string]string)
	someData["key1"] = "value1"
	someData["key2"] = "value2"
	data, _ := json.Marshal(someData)
	ioutil.WriteFile(tmpFileName, data, 0700)
	return tmpFileName
}
