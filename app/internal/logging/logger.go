package logging

import (
	"log"
	"os"
)

// New creates a standard logger configured for the application.
func New(component string) *log.Logger {
	logger := log.New(os.Stdout, "", log.LstdFlags|log.LUTC)
	logger.SetPrefix("[" + component + "] ")
	return logger
}
