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