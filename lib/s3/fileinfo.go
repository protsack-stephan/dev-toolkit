package s3

// FileInfo struct to get file information
type FileInfo struct {
	size int64
}

// Size get file size
func (fi FileInfo) Size() int64 {
	return fi.size
}
