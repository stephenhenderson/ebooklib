package main

import (
	"errors"
	"flag"
	"os"
	"fmt"

	"github.com/stephenhenderson/ebooklib/lib/config"
	"github.com/stephenhenderson/ebooklib/lib/ebooks"
	. "github.com/stephenhenderson/ebooklib/lib/logging"
	"github.com/stephenhenderson/ebooklib/lib/webservice"
)

func main() {
	appConfig := tryToLoadAppConfig()
	library := tryToInitializeLibrary(appConfig.LibraryPath)
	webservice := tryToInitializeWebService(library, appConfig.TemplatePath)
	webservice.StartService(appConfig.NetworkAddr)
}

func tryToLoadAppConfig() *config.AppConfig {
	appConfig, err := parseFlags()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config file: %v\n", err)
		flag.Usage()
		os.Exit(1)
	}
	return appConfig
}

func parseFlags() (*config.AppConfig, error) {
	configPath := flag.String(
		"config",
		"",
		"Path to config file containing")

	flag.Parse()
	if *configPath == "" {
		return nil, errors.New("Missing config path")
	}

	return config.LoadConfigFromFile(*configPath)
}

func tryToInitializeLibrary(libraryPath string) *ebooks.FileLibrary {
	library, err := ebooks.NewFileLibrary(libraryPath)
	if err != nil {
		Logger.Fatalf("Encountered fatal error %v", err)
	}
	return library
}

func tryToInitializeWebService(library *ebooks.FileLibrary, templatePath string) *webservice.EbookWebService {
	webservice, err := webservice.NewEbookWebService(library, templatePath)
	if err != nil {
		Logger.Fatalf("Error loading html templates from %s\nerr:\n%v",
			templatePath, err)
	}
	return webservice
}
