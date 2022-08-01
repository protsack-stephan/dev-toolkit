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
}
