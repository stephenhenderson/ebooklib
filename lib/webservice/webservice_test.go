package webservice

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"bytes"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/stephenhenderson/ebooklib/lib/ebooks"
	"github.com/stephenhenderson/ebooklib/lib/testutils"
	"encoding/json"
	"io/ioutil"
	"path"
	"fmt"
	"errors"
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
	library, err := ebooks.NewFileLibrary(testutils.CreateTempDir(t))
	if err != nil {
		t.Fatalf("Error creating new library %v")
	}
	webservice, err := NewEbookWebService(library, "../../templates/")
	if err != nil {
		t.Fatalf("Error creating new webservice %v", err)
	}

	ts := httptest.NewServer(http.HandlerFunc(webservice.addBookHandler))
	defer ts.Close()

	request := newfileUploadRequest(ts.URL, map[string]string{
		"title":   "Title",
		"authors": "mr writer,mrs writer",
		"year":    "2016",
		"tags":    "tag1,tag2",
	}, aJsonFile(t), t)


	skipRedirectsErr := errors.New("Don't follow redirects")
	client := &http.Client{}
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error { return skipRedirectsErr }
	resp, _ := client.Do(request)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusFound {
		t.Fatalf("Expected status code %d but got %s", http.StatusFound, resp.Status)
	}
}

// Creates a new file upload http request with optional extra params
func newfileUploadRequest(uri string, bookDetails map[string]string, filePath string, t *testing.T) *http.Request {
	file, err := os.Open(filePath)
	if err != nil {
		t.Errorf("Error creating local book file %v", err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	contentType := writer.FormDataContentType()
	part, err := writer.CreateFormFile("files", filepath.Base(filePath))
	if err != nil {
		t.Errorf("Error creating local book file %v", err)
	}
	_, err = io.Copy(part, file)

	for key, val := range bookDetails {
		_ = writer.WriteField(key, val)
	}

	err = writer.Close()
	if err != nil {
		t.Errorf("Error creating local book file %v", err)
	}

	req, err := http.NewRequest("POST", uri, body)
	if err != nil {
		t.Errorf("Error creating local book file %v", err)
	}
	req.Header.Add("Content-Type", contentType)
	fmt.Printf("Content type is %s\n", contentType)
	return req
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