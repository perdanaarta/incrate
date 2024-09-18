package storage

import (
	"errors"
	"io"
	"os"
	"path/filepath"
)

type FileStorageProvider struct {
	StoragePath string
}

func NewFileStorageProvider(storagePath string) *FileStorageProvider {
	return &FileStorageProvider{
		StoragePath: storagePath,
	}
}

func (s *FileStorageProvider) ActualPath(path string) string {
	return filepath.Join(s.StoragePath, path)
}

func (s *FileStorageProvider) List(path string) (filenames []string, err error) {
	path = s.ActualPath(path)

	entries, err := os.ReadDir(path)
	for _, entry := range entries {
		filenames = append(filenames, entry.Name())
	}

	return
}

func (s *FileStorageProvider) Delete(path string) (err error) {
	path = s.ActualPath(path)
	err = os.RemoveAll(path)
	return
}

func (s *FileStorageProvider) Move(srcPath string, dstPath string) (err error) {
	srcPath = s.ActualPath(srcPath)
	dstPath = s.ActualPath(dstPath)

	err = os.Rename(srcPath, dstPath)
	return
}

func (s *FileStorageProvider) Store(filename string, content io.Reader) (err error) {
	filename = s.ActualPath(filename)

	if err = os.MkdirAll(filepath.Dir(filename), 0770); err != nil {
		return
	}

	file, err := os.Create(filename)
	if err != nil {
		return
	}
	defer file.Close()

	_, err = io.Copy(file, content)
	return
}

func (s *FileStorageProvider) Get(filepath string) (content io.Reader, err error) {
	filepath = s.ActualPath(filepath)

	info, err := os.Stat(filepath)
	if err != nil {
		return
	}

	if info.IsDir() {
		return nil, errors.New("path is a directory")
	}

	content, err = os.Open(filepath)

	return
}

func (s *FileStorageProvider) IsExist(path string) (ok bool) {
	if _, err := os.Stat(path); err != nil {
		return false
	}

	return ok
}
