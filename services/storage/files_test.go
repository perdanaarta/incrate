package storage

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func TestFile(t *testing.T) {
	provider := NewFileStorageProvider("test")

	fileName := "testfile.txt"
	content := bytes.NewReader([]byte("test content"))

	err := provider.Store(fileName, content)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Check if file was created
	if _, err := os.Stat(provider.ActualPath(fileName)); os.IsNotExist(err) {
		t.Fatalf("expected file to exist, got %v", err)
	}
}

func TestFileStorageProvider(t *testing.T) {
	// Create a temporary directory for testing
	testDir := t.TempDir()
	provider := NewFileStorageProvider(testDir)

	// Test Store
	t.Run("Store", func(t *testing.T) {
		fileName := "testfile.txt"
		content := bytes.NewReader([]byte("test content"))

		err := provider.Store(fileName, content)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Check if file was created
		if _, err := os.Stat(provider.ActualPath(fileName)); os.IsNotExist(err) {
			t.Fatalf("expected file to exist, got %v", err)
		}
	})

	// Test Get
	t.Run("Get", func(t *testing.T) {
		fileName := "testfile.txt"
		expectedContent := "test content"

		content, err := provider.Get(fileName)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		var buf bytes.Buffer
		if _, err := io.Copy(&buf, content); err != nil {
			t.Fatalf("failed to read content: %v", err)
		}

		if got := buf.String(); got != expectedContent {
			t.Fatalf("expected content %q, got %q", expectedContent, got)
		}
	})

	// Test List
	t.Run("List", func(t *testing.T) {
		fileName := "testfile.txt"

		files, err := provider.List("")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(files) == 0 || files[0] != fileName {
			t.Fatalf("expected file %q to be listed, got %v", fileName, files)
		}
	})

	// Test Move
	t.Run("Move", func(t *testing.T) {
		srcFileName := "testfile.txt"
		dstFileName := "movedfile.txt"

		err := provider.Move(srcFileName, dstFileName)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if _, err := os.Stat(provider.ActualPath(srcFileName)); !os.IsNotExist(err) {
			t.Fatalf("expected file to be moved, but it still exists")
		}

		if _, err := os.Stat(provider.ActualPath(dstFileName)); os.IsNotExist(err) {
			t.Fatalf("expected moved file to exist, got %v", err)
		}
	})

	// Test Delete
	t.Run("Delete", func(t *testing.T) {
		fileName := "movedfile.txt"

		err := provider.Delete(fileName)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if _, err := os.Stat(provider.ActualPath(fileName)); !os.IsNotExist(err) {
			t.Fatalf("expected file to be deleted, but it still exists")
		}
	})
}
