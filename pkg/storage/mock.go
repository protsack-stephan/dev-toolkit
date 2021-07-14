package storage

import (
	"bytes"
	"io"
	"io/ioutil"
	"time"
)

// NewMock create new mock instance
func NewMock() *Mock {
	return new(Mock)
}

// Mock storage object for testing
type Mock struct{}

// List get the contents of the path
func (Mock) List(path string, options ...map[string]interface{}) ([]string, error) {
	return []string{}, nil
}

// Walk recursively look for files in directory
func (Mock) Walk(path string, callback func(path string)) error {
	return nil
}

// Create for create object in storage
func (Mock) Create(path string) (io.ReadWriteCloser, error) {
	return nil, nil
}

// Get get object from storage
func (Mock) Get(path string) (io.ReadCloser, error) {
	return ioutil.NopCloser(bytes.NewReader([]byte{})), nil
}

// Put object into storage
func (Mock) Put(path string, body io.Reader) error {
	return nil
}

// Link generate expiration link for storage
func (Mock) Link(path string, expire time.Duration) (string, error) {
	return "", nil
}

// Delete remove object from storage
func (Mock) Delete(path string) error {
	return nil
}

// Stat get object info
func (Mock) Stat(path string) (FileInfo, error) {
	return new(FileInfoMock), nil
}

// FileInfoMock mock for file information
type FileInfoMock struct{}

// Size get file info size
func (FileInfoMock) Size() int64 {
	return 0
}
