package testutils

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

// List of temp folders created during testing to be cleaned up during teardown
var tempDirs []string = []string{}

func DeleteTempDirsCreatedDuringTesting() {
	for _, dir := range tempDirs {
		err := os.RemoveAll(dir)
		if err != nil {
			fmt.Printf("Error deleting temp dir %v, err=%v\n", dir, err)
		}
	}
	tempDirs = []string{}
}

func CreateTempDir(t *testing.T) string {
	folder, err := ioutil.TempDir("", "ebook_tests")
	if err != nil {
		t.Fatalf("Error creating temp dir", err)
	}
	tempDirs = append(tempDirs, folder)
	return folder
}

func ContainsFile(targetFile string, files []os.FileInfo) bool {
	for _, file := range files {
		if file.Name() == targetFile {
			return true
		}
	}
	return false
}