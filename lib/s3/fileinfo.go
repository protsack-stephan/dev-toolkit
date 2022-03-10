package s3

import "io/fs"

// FileInfo struct to get file information
type FileInfo struct {
	size int64
	mode fs.FileMode
}

// Size get file size
func (fi FileInfo) Size() int64 {
	return fi.size
}

// Mode get unix mode.
func (fi FileInfo) Mode() fs.FileMode {
	return fi.mode
}
