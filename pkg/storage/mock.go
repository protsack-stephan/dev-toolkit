package storage

import (
	"bytes"
	"context"
	"io"
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

// Select filters the contents of an object based on SQL statement, and returns only records that match the specified SQL expression.
func (Mock) Select(path string, query string, options ...map[string]interface{}) (string, error) {
	return "", nil
}

// SelectWithContext filters the contents of an object based on SQL statement, and returns only records that match the specified SQL expression.
func (Mock) SelectWithContext(ctx context.Context, path string, query string, options ...map[string]interface{}) (string, error) {
	return "", nil
}

// Create for create object in storage
func (Mock) Create(path string) (io.ReadWriteCloser, error) {
	return nil, nil
}

// Get get object from storage
func (Mock) Get(path string) (io.ReadCloser, error) {
	return io.NopCloser(bytes.NewReader([]byte{})), nil
}

// GetWithContext get object from storage
func (Mock) GetWithContext(ctx context.Context, path string) (io.ReadCloser, error) {
	return io.NopCloser(bytes.NewReader([]byte{})), nil
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

// ActiveStatus get ActiveStatus
func (FileInfoMock) ActiveStatus() string {
	return "activeStatus"
}

// BucketKeyEnabled get BucketKeyEnabled
func (FileInfoMock) BucketKeyEnabled() bool {
	return true
}

// DeleteMarker get DeleteMarker
func (FileInfoMock) DeleteMarker() bool {
	return true
}

// Expiration get Expiration
func (FileInfoMock) Expiration() string {
	return "expiration"
}

// Metadata get Metadata
func (FileInfoMock) Metadata() map[string]*string {
	return map[string]*string{}
}

// MissingMeta get MissingMeta
func (FileInfoMock) MissingMeta() int64 {
	return 1
}

// ObjectLockLegalHoldStatus get ObjectLockLegalHoldStatus
func (FileInfoMock) ObjectLockLegalHoldStatus() string {
	return "objectLockLegalHoldStatus"
}

// ObjectLockMode get ObjectLockMode
func (FileInfoMock) ObjectLockMode() string {
	return "objectLockMode"
}

// ObjectLockRetainUntilDate get ObjectLockRetainUntilDate
func (FileInfoMock) ObjectLockRetainUntilDate() time.Time {
	return time.Now()
}

// PartsCount get PartsCount
func (FileInfoMock) PartsCount() int64 {
	return 1
}

// ReplicationStatus get ReplicationStatus
func (FileInfoMock) ReplicationStatus() string {
	return "replicationStatus"
}

// RequestCharged get RequestCharged
func (FileInfoMock) RequestCharged() string {
	return "requestCharged"
}

// Restore get Restore
func (FileInfoMock) Restore() string {
	return "restore"
}

// SSECustomerAlgorithm get SSECustomerAlgorithm
func (FileInfoMock) SSECustomerAlgorithm() string {
	return "sseCustomerAlgorithm"
}

// SSECustomerKeyMD5 get SSECustomerKeyMD5
func (FileInfoMock) SSECustomerKeyMD5() string {
	return "sseCustomerKeyMD5"
}

// SSEKMSKeyId get SSEKMSKeyId
func (FileInfoMock) SSEKMSKeyId() string {
	return "sseKMSKeyId"
}

// ServerSideEncryption get ServerSideEncryption
func (FileInfoMock) ServerSideEncryption() string {
	return "serverSideEncryption"
}

// StorageClass get StorageClass
func (FileInfoMock) StorageClass() string {
	return "storageClass"
}

// VersionId get VersionId
func (FileInfoMock) VersionId() string {
	return "versionId"
}

// WebsiteRedirectLocation get WebsiteRedirectLocation
func (FileInfoMock) WebsiteRedirectLocation() string {
	return "websiteRedirectLocation"
}
