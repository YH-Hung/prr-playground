package logger

import (
	"bytes"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	var buf bytes.Buffer
	logger := New(&buf, "[TEST] ")

	logger.Println("test message")

	output := buf.String()
	if !strings.Contains(output, "[TEST]") {
		t.Errorf("Logger output missing prefix: %v", output)
	}
	if !strings.Contains(output, "test message") {
		t.Errorf("Logger output missing message: %v", output)
	}
}

func TestNewMultipleLoggers(t *testing.T) {
	var buf1, buf2 bytes.Buffer

	logger1 := New(&buf1, "[LOG1] ")
	logger2 := New(&buf2, "[LOG2] ")

	logger1.Println("message 1")
	logger2.Println("message 2")

	out1 := buf1.String()
	out2 := buf2.String()

	if !strings.Contains(out1, "[LOG1]") || !strings.Contains(out1, "message 1") {
		t.Errorf("Logger1 output incorrect: %v", out1)
	}
	if !strings.Contains(out2, "[LOG2]") || !strings.Contains(out2, "message 2") {
		t.Errorf("Logger2 output incorrect: %v", out2)
	}
}
