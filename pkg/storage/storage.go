package storage

import (
	"io"
	"time"
)

// Storage list of storage interfaces
type Storage interface {
	Lister
	Walker
	Creator
	Getter
	Putter
	Linker
	Deleter
	Stater
}

// Lister get the contents of the path
type Lister interface {
	List(path string, options ...map[string]interface{}) ([]string, error)
}

// Walker recursively look for files in directory
type Walker interface {
	Walk(path string, callback func(path string)) error
}

// Creator create newfile or open current and truncate
type Creator interface {
	Create(path string) (io.ReadWriteCloser, error)
}

// Getter get object from storage
type Getter interface {
	Get(path string) (io.ReadCloser, error)
}

// Putter move object to storage
type Putter interface {
	Put(path string, body io.Reader) error
}

// Linker get dowload link with expiration
type Linker interface {
	Link(path string, expire time.Duration) (string, error)
}

// Deleter delete object from storage
type Deleter interface {
	Delete(path string) error
}

// Stater get information about the file
type Stater interface {
	Stat(path string) (FileInfo, error)
}
