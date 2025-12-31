// Package logger provides structured JSON logging utilities.
package logger

import (
	"io"
	"log"
)

// New creates a new logger that writes to the given output with the specified prefix.
// The logger uses standard log flags for timestamp and file information.
func New(output io.Writer, prefix string) *log.Logger {
	return log.New(output, prefix, log.LstdFlags)
}
