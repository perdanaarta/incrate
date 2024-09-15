package storage

import "io"

type StorageProvider interface {
	List(path string) ([]string, error)
	Get(filepath string) (error, io.Reader)
	Delete(filepath string) error
	Move(srcPath string, dstPath string) error
}
