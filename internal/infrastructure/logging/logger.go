package logging

import (
	"log"
	"os"
)

func New() *log.Logger {
	return log.New(os.Stderr, "", log.LstdFlags|log.Lmicroseconds)
}
