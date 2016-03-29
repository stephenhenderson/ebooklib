package logging

import (
	"os"
	"log"
)

var Logger = log.New(os.Stdout, "", log.Ldate | log.Ltime)
