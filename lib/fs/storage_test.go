package fs

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/protsack-stephan/dev-toolkit/pkg/storage"
	"github.com/stretchr/testify/assert"
)

const storageTestVol = "./testdata"
const storageTestPath = "test.txt"
const copyDestPath = "dest.txt"
const storageTestWalkPath = "testwalk.txt"

var storageTestExpire = time.Second * 1
var storageTestData = []byte("hello storage")

func testStorage(storage storage.Storage) error {
	return nil
}

func TestStorage(t *testing.T) {
	store := NewStorage(storageTestVol)
	assert := assert.New(t)
	assert.Nil(testStorage(store))
	ctx := context.Background()

	t.Run("List path's content", func(t *testing.T) {
		content, err := store.List("/")
		assert.NoError(err)
		assert.Equal(content, []string{storageTestWalkPath})
	})

	t.Run("list path's content with context", func(t *testing.T) {
		content, err := store.ListWithContext(ctx, "/")
		assert.NoError(err)
		assert.Equal(content, []string{storageTestWalkPath})
	})

	t.Run("walk path", func(t *testing.T) {
		assert.NoError(store.Walk("/", func(path string) {
			assert.Equal(storageTestWalkPath, path)
		}))
	})

	t.Run("walk path with context", func(t *testing.T) {
		assert.NoError(store.WalkWithContext(ctx, "/", func(path string) {
			assert.Equal(storageTestWalkPath, path)
		}))
	})

	t.Run("create file", func(t *testing.T) {
		file, err := store.Create(storageTestPath)
		assert.NoError(err)
		file.Close()
		assert.NoError(os.Remove(fmt.Sprintf("%s/%s", storageTestVol, storageTestPath)))
	})

	t.Run("put file", func(t *testing.T) {
		assert.NoError(store.Put(storageTestPath, bytes.NewReader(storageTestData)))
	})

	t.Run("put file with context", func(t *testing.T) {
		assert.NoError(store.PutWithContext(ctx, storageTestPath, bytes.NewReader(storageTestData)))
	})

	t.Run("stat file", func(t *testing.T) {
		info, err := store.Stat(storageTestPath)
		assert.NoError(err)
		assert.NotZero(info.Size())
	})

	t.Run("get file", func(t *testing.T) {
		body, err := store.Get(storageTestPath)
		assert.NoError(err)
		defer body.Close()

		data, err := ioutil.ReadAll(body)
		assert.NoError(err)
		assert.Equal(storageTestData, data)
	})

	t.Run("get file with context", func(t *testing.T) {
		body, err := store.GetWithContext(ctx, storageTestPath)
		assert.NoError(err)
		defer body.Close()

		data, err := ioutil.ReadAll(body)
		assert.NoError(err)
		assert.Equal(storageTestData, data)
	})

	t.Run("link file", func(t *testing.T) {
		loc, err := store.Link(storageTestPath, storageTestExpire)
		assert.NoError(err)
		assert.Equal(fmt.Sprintf("%s/%s", storageTestVol, storageTestPath), loc)
	})

	t.Run("delete file", func(t *testing.T) {
		assert.NoError(store.Delete(storageTestPath))
	})

	t.Run("delete file with context", func(t *testing.T) {
		assert.NoError(store.PutWithContext(ctx, storageTestPath, bytes.NewReader(storageTestData)))
		assert.NoError(store.DeleteWithContext(ctx, storageTestPath))
	})

	t.Run("copy file with permission argument", func(t *testing.T) {
		assert.NoError(store.Put(storageTestPath, bytes.NewReader(storageTestData)))
		options := []map[string]interface{}{
			{"mode": 0711},
		}
		assert.NoError(store.Copy(fmt.Sprintf("%s/%s", storageTestVol, storageTestPath), fmt.Sprintf("%s/%s", storageTestVol, copyDestPath), options...))
		assert.NoError(store.compareFileContent(copyDestPath, storageTestData))
		assert.NoError(compareFileMode(copyDestPath, options[0]["mode"].(int)))
		assert.NoError(os.Remove(fmt.Sprintf("%s/%s", storageTestVol, storageTestPath)))
		assert.NoError(os.Remove(fmt.Sprintf("%s/%s", storageTestVol, copyDestPath)))
	})

	t.Run("copy file with context", func(t *testing.T) {
		assert.NoError(store.Put(storageTestPath, bytes.NewReader(storageTestData)))
		assert.NoError(store.CopyWithContext(ctx, fmt.Sprintf("%s/%s", storageTestVol, storageTestPath), fmt.Sprintf("%s/%s", storageTestVol, copyDestPath)))
		assert.NoError(store.compareFileContent(copyDestPath, storageTestData))
		assert.NoError(compareFileMode(copyDestPath, 0644))
		assert.NoError(os.Remove(fmt.Sprintf("%s/%s", storageTestVol, storageTestPath)))
		assert.NoError(os.Remove(fmt.Sprintf("%s/%s", storageTestVol, copyDestPath)))
	})
}

func (store *Storage) compareFileContent(filePath string, content []byte) error {
	body, err := store.Get(copyDestPath)

	if err != nil {
		return err
	}
	defer body.Close()

	data, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}

	if res := bytes.Compare(data, content); res != 0 {
		return errors.New("Contents not equal")
	}

	return nil
}

func compareFileMode(filePath string, mode int) error {
	info, err := os.Stat(fmt.Sprintf("%s/%s", storageTestVol, filePath))

	if err != nil {
		return err
	}

	if info.Mode() != fs.FileMode(mode) {
		return errors.New("Wrong file permission set.")
	}

	return nil
}
