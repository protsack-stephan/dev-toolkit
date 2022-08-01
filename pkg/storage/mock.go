package storage

import (
	"bytes"
	"context"
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

// ListWithContext get the contents of the path
func (Mock) ListWithContext(ctx context.Context, path string, options ...map[string]interface{}) ([]string, error) {
	return []string{}, nil
}

// Walk recursively look for files in directory
func (Mock) Walk(path string, callback func(path string)) error {
	return nil
}

// WalkWithContext recursively look for files in directory
func (Mock) WalkWithContext(ctx context.Context, path string, callback func(path string)) error {
	return nil
}

// Copy copies an object from the a path in a bucket to another path in the same or different bucket.
func (Mock) Copy(src string, dst string, options ...map[string]interface{}) error {
	return nil
}

// CopyWithContext copies an object from the a path in a bucket to another path in the same or different bucket.
func (Mock) CopyWithContext(ctx context.Context, src string, dst string, options ...map[string]interface{}) error {
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

// GetWithContext get object from storage
func (Mock) GetWithContext(ctx context.Context, path string) (io.ReadCloser, error) {
	return ioutil.NopCloser(bytes.NewReader([]byte{})), nil
}

// Put object into storage
func (Mock) Put(path string, body io.Reader) error {
	return nil
}

// PutWithContext object into storage
func (Mock) PutWithContext(ctx context.Context, path string, body io.Reader) error {
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

// DeleteWithContext remove object from storage
func (Mock) DeleteWithContext(ctx context.Context, path string) error {
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

// Size get file info accept-ranges
func (FileInfoMock) AcceptRanges() string {
	return "accept-ranges"
}

// Size get file info Cache-Control
func (FileInfoMock) CacheControl() string {
	return "Cache-Control"
}

// Size get file info Content-Disposition
func (FileInfoMock) ContentDisposition() string {
	return "Content-Disposition"
}

// Size get file info Content-Encoding
func (FileInfoMock) ContentEncoding() string {
	return "Content-Encoding"
}

// Size get file info Content-Language
func (FileInfoMock) ContentLanguage() string {
	return "Content-Language"
}

// Size get file info Content-Type
func (FileInfoMock) ContentType() string {
	return "Content-Type"
}

// Size get file info ETag
func (FileInfoMock) ETag() string {
	return "ETag"
}

// Size get file info Expires
func (FileInfoMock) Expires() string {
	return "Expires"
}

// Size get file info Last-Modified
func (FileInfoMock) LastModified() time.Time {
	return time.Now()
}
