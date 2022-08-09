package storage

import "time"

// FileInfo list of file properties
type FileInfo interface {
	Size() int64
	AcceptRanges() string
	CacheControl() string
	ContentDisposition() string
	ContentEncoding() string
	ContentLanguage() string
	ContentType() string
	ETag() string
	Expires() string
	LastModified() time.Time
	ActiveStatus() string
	DeleteMarker() bool
	Expiration() string
	Metadata() map[string]*string
	MissingMeta() int64
	ObjectLockLegalHoldStatus() string
	ObjectLockMode() string
	ObjectLockRetainUntilDate() time.Time
	PartsCount() int64
	ReplicationStatus() string
	RequestCharged() string
	Restore() string
	SSECustomerAlgorithm() string
	SSECustomerKeyMD5() string
	SSEKMSKeyId() string
	ServerSideEncryption() string
	StorageClass() string
	VersionId() string
	WebsiteRedirectLocation() string
}
