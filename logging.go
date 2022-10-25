package calcutron

import (
	"io"
	"log"
	"os"
)

var (
	DebugLog *log.Logger
	WarnLog  *log.Logger
	InfoLog  *log.Logger
	ErrorLog *log.Logger
)

// Called first in any package
func init() {
	file, err := os.OpenFile("logs.txt",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	flags := log.Ldate | log.Ltime | log.Lshortfile

	DebugLog = log.New(file, "DEBUG: ", flags)
	InfoLog = log.New(io.Discard, "INFO: ", flags)
	WarnLog = log.New(file, "WARNING: ", flags)
	ErrorLog = log.New(file, "ERROR: ", flags)
}
