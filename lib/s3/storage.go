// Package s3 this package in intention to hide aws s3 storage implementation
// under the interface that will give you the ability to user other cloud providers
// in the future
package s3

import (
	"errors"
	"io"
	pathTool "path"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	s3manager "github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/protsack-stephan/dev-toolkit/pkg/storage"
)

const maxUploadParts = 20000
const partSize = 1024 * 1024 * 5 * 2

// NewStorage create new storage instance
func NewStorage(ses *session.Session, bucket string) *Storage {
	return &Storage{
		s3:     s3.New(ses),
		bucket: bucket,
		uploader: s3manager.NewUploader(ses, func(upl *s3manager.Uploader) {
			upl.MaxUploadParts = maxUploadParts
			upl.PartSize = partSize
		}),
	}
}

// Storage interface adaptation for s3
type Storage struct {
	bucket   string
	uploader *s3manager.Uploader
	s3       *s3.S3
}

// List reads the path content
// Note: it is not better performant than Walk because it still
// iterates over all objects in a path
func (s *Storage) List(path string) ([]string, error) {
	var result []string
	res, err := s.s3.ListObjects(&s3.ListObjectsInput{
		Bucket: aws.String(s.bucket),
		Prefix: aws.String(path),
	})

	if err != nil {
		return result, err
	}

	for _, object := range res.Contents {
		// Check whether the object is nested in a path
		p := strings.Split(*object.Key, "/")

		if len(p) == 1 {
			// It's a file
			result = append(result, *object.Key)
		} else if result[len(p)-1] != p[0] {
			// It's a folder
			// And it is not added yet
			// res.Contents is sorted so if p[0] is not unique it would appear last in the result
			result = append(result, p[0])
		}
	}

	return result, nil
}

// ListDelimiterPages reads the path content with a delimiter by pages
func (s *Storage) ListDelimiterPages(path string, delimiter string) ([]string, error) {
	var result []string
	err := s.s3.ListObjectsPages(
		&s3.ListObjectsInput{
			Bucket:    aws.String(s.bucket),
			Prefix:    aws.String(path),
			Delimiter: aws.String(delimiter),
		},
		// handle bulks of 1000 keys
		func(page *s3.ListObjectsOutput, lastPage bool) bool {
			for _, prefix := range page.CommonPrefixes {
				result = append(result, pathTool.Base(*prefix.Prefix))
			}

			return true
		},
	)

	return result, err
}

// Walk recursively look for files in directory
func (s *Storage) Walk(path string, callback func(path string)) error {
	res, err := s.s3.ListObjects(&s3.ListObjectsInput{
		Bucket: aws.String(s.bucket),
		Prefix: aws.String(path),
	})

	if err != nil {
		return err
	}

	for _, object := range res.Contents {
		callback(*object.Key)
	}

	return nil
}

// Create for create interface
func (s *Storage) Create(path string) (io.ReadWriteCloser, error) {
	return nil, errors.New("method unimplemented")
}

// Get file from s3 bucket
func (s *Storage) Get(path string) (io.ReadCloser, error) {
	out, err := s.s3.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
	})

	return out.Body, err
}

// Put file into s3 bucket
func (s *Storage) Put(path string, body io.Reader) error {
	_, err := s.uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
		Body:   body,
	})

	return err
}

// Link generate expiration link for s3 access
func (s *Storage) Link(path string, expire time.Duration) (string, error) {
	req, _ := s.s3.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
	})

	return req.Presign(expire)
}

// Delete remove object from s3
func (s *Storage) Delete(path string) error {
	_, err := s.s3.DeleteObject(&s3.DeleteObjectInput{
		Key:    aws.String(path),
		Bucket: aws.String(s.bucket),
	})

	return err
}

// Stat get object info
func (s *Storage) Stat(path string) (storage.FileInfo, error) {
	out, err := s.s3.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
	})

	if err != nil {
		return nil, err
	}

	return &FileInfo{*out.ContentLength}, nil
}
