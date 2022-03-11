package fs

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/karrick/godirwalk"
	"github.com/protsack-stephan/dev-toolkit/pkg/storage"
)

// ErrEmptyPath method call with empty file path
var ErrEmptyPath = errors.New("empty file path")

// NewStorage create new storage instance
func NewStorage(vol string) *Storage {
	loc := vol

	if vol[len(vol)-1:] != "/" {
		loc = fmt.Sprintf("%s/", vol)
	}

	return &Storage{
		loc,
	}
}

// Storage file system manipulations manager
type Storage struct {
	vol string
}

// List reads the path content
func (s Storage) List(path string, options ...map[string]interface{}) ([]string, error) {
	dir, err := s.fullPath(path)
	if err != nil {
		return []string{}, err
	}

	d, err := os.Open(dir)
	if err != nil {
		return []string{}, err
	}
	defer d.Close()

	return d.Readdirnames(-1)
}

// ListWithContext reads the path content
func (s Storage) ListWithContext(_ context.Context, path string, options ...map[string]interface{}) ([]string, error) {
	return s.List(path, options...)
}

// Walk recursively look for files in directory
func (s Storage) Walk(path string, callback func(path string)) error {
	loc, err := s.fullPath(path)

	if err != nil {
		return err
	}

	slice := len(s.vol)

	if s.vol[:2] == "./" {
		slice -= 2
	}

	return godirwalk.Walk(loc, &godirwalk.Options{
		Unsorted: true,
		ErrorCallback: func(osPathname string, err error) godirwalk.ErrorAction {
			return godirwalk.SkipNode
		},
		Callback: func(path string, de *godirwalk.Dirent) error {
			if !de.IsDir() {
				callback(path[slice:])
			}

			return nil
		},
	})
}

// WalkWithContext recursively look for files in directory
func (s Storage) WalkWithContext(_ context.Context, path string, callback func(path string)) error {
	return s.Walk(path, callback)
}

// Copy copies a file.
func (s Storage) Copy(src string, dst string, options ...map[string]interface{}) error {
	mode := 0644

	for _, opt := range options {
		if v, ok := opt["mode"]; ok {
			if m, ok := v.(int); ok {
				mode = m
			}
		}
	}

	input, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(dst, input, fs.FileMode(mode))
}

// CopyWithContext copies a file.
func (s Storage) CopyWithContext(_ context.Context, src string, dst string, options ...map[string]interface{}) error {
	return s.Copy(src, dst, options...)
}

// Create create new file or open existing one and truncate it
func (s Storage) Create(path string) (io.ReadWriteCloser, error) {
	loc, err := s.fullPath(path)

	if err != nil {
		return nil, err
	}

	dir, _ := filepath.Split(loc)
	_, err = os.Stat(dir)

	if err != nil && os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0766)
	}

	if err != nil {
		return nil, err
	}

	return os.OpenFile(loc, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0766)
}

// Get get object from storage
func (s Storage) Get(path string) (io.ReadCloser, error) {
	loc, err := s.fullPath(path)

	if err != nil {
		return nil, err
	}

	return os.Open(loc)
}

// GetWithContext get object from storage
func (s Storage) GetWithContext(_ context.Context, path string) (io.ReadCloser, error) {
	return s.Get(path)
}

// Put object into storage
func (s Storage) Put(path string, body io.Reader) error {
	loc, err := s.fullPath(path)

	if err != nil {
		return err
	}

	dir, _ := filepath.Split(loc)
	_, err = os.Stat(dir)

	if err != nil && os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0766)
	}

	if err != nil {
		return err
	}

	buff, err := ioutil.ReadAll(body)

	if err != nil {
		return err
	}

	return ioutil.WriteFile(loc, buff, 0766)
}

// PutWithContext object into storage
func (s Storage) PutWithContext(_ context.Context, path string, body io.Reader) error {
	return s.Put(path, body)
}

// Link generate expiration link for storage
func (s Storage) Link(path string, expire time.Duration) (string, error) {
	return s.fullPath(path)
}

// Delete remove object from storage
func (s *Storage) Delete(path string) error {
	loc, err := s.fullPath(path)

	if err != nil {
		return err
	}

	return os.Remove(loc)
}

// DeleteWithContext remove object from storage
func (s *Storage) DeleteWithContext(_ context.Context, path string) error {
	return s.Delete(path)
}

// Stat ge file information
func (s Storage) Stat(path string) (storage.FileInfo, error) {
	loc, err := s.fullPath(path)

	if err != nil {
		return nil, err
	}

	info, err := os.Stat(loc)

	if err != nil {
		return nil, err
	}

	return info, err
}

func (s Storage) fullPath(path string) (string, error) {
	var err error

	if len(path) <= 0 {
		err = ErrEmptyPath
	}

	return fmt.Sprintf("%s%s", s.vol, strings.TrimPrefix(path, "/")), err
}
