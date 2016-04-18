package config

import (
	"testing"
	"github.com/stephenhenderson/ebooklib/lib/testutils"
	"io/ioutil"
	"path/filepath"
)

func TestMain(m *testing.M) {
	defer testutils.DeleteTempDirsCreatedDuringTesting()
	m.Run()
}

func TestReturnsErrorIfConfigFileNotFound(t *testing.T) {
	_, err := LoadConfigFromFile("aFileWhichDoesNotExist.json")
	if err == nil {
		t.Fatal("Expected error when config file does not exist")
	}
}

func TestReturnsErrorIfConfigIsNotValidJson(t *testing.T) {
	emptyConfigFile := tempConfigFile(t, []byte("{\"badJson"))
	_, err := LoadConfigFromFile(emptyConfigFile)
	if err == nil {
		t.Fatal("Expected error when config file contains invalid json")
	}
}

func TestReturnsErrorIfConfigIsEmpty(t *testing.T) {
	emptyConfigFile := tempConfigFile(t, []byte("{}"))
	_, err := LoadConfigFromFile(emptyConfigFile)
	if err == nil {
		t.Fatal("Expected error when config file is empty")
	}
}

func tempConfigFile(t *testing.T, data []byte) string {
	tempDir := testutils.CreateTempDir(t)
	configPath := filepath.Join(tempDir, "config.json")
	err := ioutil.WriteFile(configPath, data, 0700)
	if err != nil {
		t.Fatalf("Error writing config file: %v", err)
	}
	return configPath
}