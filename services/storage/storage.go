package storage

import (
	"io"
)

/*
Interface to store into storage and stuff
*/
type StorageProvider interface {

	// List all item inside the path
	List(path string) ([]string, error)

	// Check if file or directory exist
	IsExist(path string) bool

	// Get the file
	Get(filename string) (io.Reader, error)

	// Move the file or dir
	Move(srcPath string, dstPath string) error

	// Store a file (Automatically create parent directory)
	Store(filename string, content io.Reader) error

	// Delete directory and files
	Delete(path string) error
}
