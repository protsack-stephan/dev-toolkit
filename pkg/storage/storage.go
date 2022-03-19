package storage

import (
	"context"
	"io"
	"time"
)

// Storage list of storage interfaces
type Storage interface {
	Lister
	ListerWithContext
	Walker
	WalkerWithContext
	Copier
	CopierWithContext
	Creator
	Getter
	GetterWithContext
	Putter
	PutterWithContext
	Linker
	Deleter
	DeleterWithContext
	Stater
}

// Lister get the contents of the path
type Lister interface {
	List(path string, options ...map[string]interface{}) ([]string, error)
}

// ListerWithContext get the contents of the path
type ListerWithContext interface {
	ListWithContext(ctx context.Context, path string, options ...map[string]interface{}) ([]string, error)
}

// Walker recursively look for files in directory
type Walker interface {
	Walk(path string, callback func(path string)) error
}

// WalkerWithContext recursively look for files in directory
type WalkerWithContext interface {
	WalkWithContext(ctx context.Context, path string, callback func(path string)) error
}

// Copier copies an object from the a path in a bucket to another path in the same or different bucket.
// 'src' and 'dst' are absolute paths of the file.
type Copier interface {
	Copy(src string, dst string, options ...map[string]interface{}) error
}

// CopierWithContext copies an object from the a path in a bucket to another path in the same or different bucket.
// 'src' and 'dst' are absolute paths of the file.
type CopierWithContext interface {
	CopyWithContext(ctx context.Context, src string, dst string, options ...map[string]interface{}) error
}

// Creator create newfile or open current and truncate
type Creator interface {
	Create(path string) (io.ReadWriteCloser, error)
}

// Getter get object from storage
type Getter interface {
	Get(path string) (io.ReadCloser, error)
}

// GetterWithContext get object from storage
type GetterWithContext interface {
	GetWithContext(ctx context.Context, path string) (io.ReadCloser, error)
}

// Putter move object to storage
type Putter interface {
	Put(path string, body io.Reader) error
}

// PutterWithContext moves object to storage
type PutterWithContext interface {
	PutWithContext(ctx context.Context, path string, body io.Reader) error
}

// Linker get dowload link with expiration
type Linker interface {
	Link(path string, expire time.Duration) (string, error)
}

// Deleter delete object from storage
type Deleter interface {
	Delete(path string) error
}

// DeleterWithContext delete object from storage
type DeleterWithContext interface {
	DeleteWithContext(ctx context.Context, path string) error
}

// Stater get information about the file
type Stater interface {
	Stat(path string) (FileInfo, error)
}
