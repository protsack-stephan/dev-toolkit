package fs

import "time"

// FileInfo struct to get file information
type FileInfo struct {
	size               int64
	acceptRanges       string
	cacheControl       string
	contentDisposition string
	contentEncoding    string
	contentLanguage    string
	contentType        string
	eTag               string
	expires            string
	lastModified       time.Time
}

// Size get file size
func (fi FileInfo) Size() int64 {
	return fi.size
}

// AcceptRanges get accept-ranges
func (fi FileInfo) AcceptRanges() string {
	return fi.acceptRanges
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

// ETag get ETag
func (fi FileInfo) ETag() string {
	return fi.eTag
}

// Expires get Expires
func (fi FileInfo) Expires() string {
	return fi.expires
}

// LastModified get Last-Modified
func (fi FileInfo) LastModified() time.Time {
	return fi.lastModified
}
