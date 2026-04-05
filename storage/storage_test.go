package storage

import (
	"os"
	"path/filepath"
	"testing"

	"go.uber.org/zap/zaptest"
)

func TestSaveMail(t *testing.T) {
	// Arrange
	logger := zaptest.NewLogger(t)
	rootDir, subDir, mailID := "./backup", "tests", "123"
	raw := []byte("Content-Type: text/html\r\n\r\n<h1>Test</h1>")

	// Act
	labelDir := filepath.Join(rootDir, subDir)
	path := filepath.Join(labelDir, mailID)

	err := SaveMail(labelDir, mailID, raw, logger)

	// Assert
	if err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(path); err != nil {
		t.Errorf("file not created")
	}

	os.RemoveAll(rootDir)
}
