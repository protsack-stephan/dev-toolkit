package storage

import "io/fs"

// FileInfo list of file properties
type FileInfo interface {
	Size() int64
	Mode() fs.FileMode
}
