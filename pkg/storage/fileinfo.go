package storage

// FileInfo list of file properties
type FileInfo interface {
	Size() int64
}
