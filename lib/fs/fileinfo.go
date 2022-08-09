package fs

import "time"

// FileInfo struct to get file information
type FileInfo struct {
	size                      int64
	acceptRanges              string
	activeStatus              string
	bucketKeyEnabled          bool
	cacheControl              string
	contentDisposition        string
	contentEncoding           string
	contentLanguage           string
	contentType               string
	deleteMarker              bool
	eTag                      string
	expiration                string
	expires                   string
	lastModified              time.Time
	metadata                  map[string]*string
	missingMeta               int64
	objectLockLegalHoldStatus string
	objectLockMode            string
	objectLockRetainUntilDate time.Time
	partsCount                int64
	replicationStatus         string
	requestCharged            string
	restore                   string
	sseCustomerAlgorithm      string
	sseCustomerKeyMD5         string
	sseKMSKeyId               string
	serverSideEncryption      string
	storageClass              string
	versionId                 string
	websiteRedirectLocation   string
}

// Size get file size
func (fi FileInfo) Size() int64 {
	return fi.size
}

// AcceptRanges get accept-ranges
func (fi FileInfo) AcceptRanges() string {
	return fi.acceptRanges
}

// ActiveStatus get ActiveStatus
func (fi FileInfo) ActiveStatus() string {
	return fi.activeStatus
}

// BucketKeyEnabled get BucketKeyEnabled
func (fi FileInfo) BucketKeyEnabled() bool {
	return fi.bucketKeyEnabled
}

// CacheControl get Cache-Control
func (fi FileInfo) CacheControl() string {
	return fi.cacheControl
}

// ContentDisposition get Content-Disposition
func (fi FileInfo) ContentDisposition() string {
	return fi.contentDisposition
}

// ContentEncoding get Content-Encoding
func (fi FileInfo) ContentEncoding() string {
	return fi.contentEncoding
}

// ContentLanguage get Content-Language
func (fi FileInfo) ContentLanguage() string {
	return fi.contentLanguage
}

// ContentType get Content-Type
func (fi FileInfo) ContentType() string {
	return fi.contentType
}

// DeleteMarker get DeleteMarker
func (fi FileInfo) DeleteMarker() bool {
	return fi.deleteMarker
}

// ETag get ETag
func (fi FileInfo) ETag() string {
	return fi.eTag
}

// Expires get Expires
func (fi FileInfo) Expires() string {
	return fi.expires
}

// Expiration get Expiration
func (fi FileInfo) Expiration() string {
	return fi.expiration
}

// LastModified get Last-Modified
func (fi FileInfo) LastModified() time.Time {
	return fi.lastModified
}

// Metadata get Metadata
func (fi FileInfo) Metadata() map[string]*string {
	return fi.metadata
}

// MissingMeta get MissingMeta
func (fi FileInfo) MissingMeta() int64 {
	return fi.missingMeta
}

// ObjectLockLegalHoldStatus get ObjectLockLegalHoldStatus
func (fi FileInfo) ObjectLockLegalHoldStatus() string {
	return fi.objectLockLegalHoldStatus
}

// ObjectLockMode get ObjectLockMode
func (fi FileInfo) ObjectLockMode() string {
	return fi.objectLockMode
}

// ObjectLockRetainUntilDate get ObjectLockRetainUntilDate
func (fi FileInfo) ObjectLockRetainUntilDate() time.Time {
	return fi.objectLockRetainUntilDate
}

// PartsCount get PartsCount
func (fi FileInfo) PartsCount() int64 {
	return fi.partsCount
}

// ReplicationStatus get ReplicationStatus
func (fi FileInfo) ReplicationStatus() string {
	return fi.replicationStatus
}

// RequestCharged get RequestCharged
func (fi FileInfo) RequestCharged() string {
	return fi.requestCharged
}

// Restore get Restore
func (fi FileInfo) Restore() string {
	return fi.restore
}

// SSECustomerAlgorithm get SSECustomerAlgorithm
func (fi FileInfo) SSECustomerAlgorithm() string {
	return fi.sseCustomerAlgorithm
}

// SSECustomerKeyMD5 get SSECustomerKeyMD5
func (fi FileInfo) SSECustomerKeyMD5() string {
	return fi.sseCustomerKeyMD5
}

// SSEKMSKeyId get SSEKMSKeyId
func (fi FileInfo) SSEKMSKeyId() string {
	return fi.sseKMSKeyId
}

// ServerSideEncryption get ServerSideEncryption
func (fi FileInfo) ServerSideEncryption() string {
	return fi.serverSideEncryption
}

// StorageClass get StorageClass
func (fi FileInfo) StorageClass() string {
	return fi.storageClass
}

// VersionId get VersionId
func (fi FileInfo) VersionId() string {
	return fi.versionId
}

// WebsiteRedirectLocation get WebsiteRedirectLocation
func (fi FileInfo) WebsiteRedirectLocation() string {
	return fi.websiteRedirectLocation
}
