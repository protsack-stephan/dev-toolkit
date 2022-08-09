// Package s3 this package in intention to hide aws s3 storage implementation
// under the interface that will give you the ability to user other cloud providers
// in the future
package s3

import (
	"errors"
	"fmt"
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

func contents(res *s3.ListObjectsOutput) []string {
	result := make([]string, 0)

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

	return result
}

func prefixes(res *s3.ListObjectsOutput) []string {
	result := make([]string, 0)

	for _, prefix := range res.CommonPrefixes {
		result = append(result, pathTool.Base(*prefix.Prefix))
	}

	return result
}

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

// List reads the path content or prefixes.
func (s *Storage) List(path string, options ...map[string]interface{}) ([]string, error) {
	var result []string

	input := s3.ListObjectsInput{
		Bucket: aws.String(s.bucket),
		Prefix: aws.String(path),
	}

	if len(options) > 0 {
		if options[0]["delimiter"] != nil {
			input.SetDelimiter(options[0]["delimiter"].(string))
		}
	}

	err := s.s3.ListObjectsPages(
		&input,
		// handle bulks of 1000 keys
		func(res *s3.ListObjectsOutput, _ bool) bool {
			if input.Delimiter != nil {
				result = append(result, prefixes(res)...)
			} else {
				result = append(result, contents(res)...)
			}

			return true
		},
	)

	return result, err
}

// ListWithContext reads the path content or prefixes
func (s *Storage) ListWithContext(ctx aws.Context, path string, options ...map[string]interface{}) ([]string, error) {
	var result []string

	input := s3.ListObjectsInput{
		Bucket: aws.String(s.bucket),
		Prefix: aws.String(path),
	}

	if len(options) > 0 {
		if options[0]["delimiter"] != nil {
			input.SetDelimiter(options[0]["delimiter"].(string))
		}
	}

	err := s.s3.ListObjectsPagesWithContext(
		ctx,
		&input,
		// handle bulks of 1000 keys
		func(res *s3.ListObjectsOutput, _ bool) bool {
			if input.Delimiter != nil {
				result = append(result, prefixes(res)...)
			} else {
				result = append(result, contents(res)...)
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

// WalkWithContext recursively look for files in directory
func (s *Storage) WalkWithContext(ctx aws.Context, path string, callback func(path string)) error {
	res, err := s.s3.ListObjectsWithContext(ctx, &s3.ListObjectsInput{
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

// Copy copies an object from the a path in a bucket to another path in the same or different bucket.
// 'src' and 'dst' are absolute paths of the file.
func (s *Storage) Copy(src string, dst string, options ...map[string]interface{}) error {
	bucket := s.bucket

	for _, opt := range options {
		if v, ok := opt["bucket"]; ok {
			if bkt, ok := v.(string); ok {
				bucket = bkt
			}
		}
	}

	_, err := s.s3.CopyObject(&s3.CopyObjectInput{
		Bucket:     aws.String(bucket),
		CopySource: aws.String(fmt.Sprintf("%s/%s", s.bucket, src)),
		Key:        aws.String(dst),
	})

	return err
}

// CopyWithContext copies an object from the a path in a bucket to another path in the same or different bucket.
// 'src' and 'dst' are absolute paths of the file.
func (s *Storage) CopyWithContext(ctx aws.Context, src string, dst string, options ...map[string]interface{}) error {
	bucket := s.bucket

	for _, opt := range options {
		if v, ok := opt["dstBucket"]; ok {
			if bkt, ok := v.(string); ok {
				bucket = bkt
			}
		}
	}

	_, err := s.s3.CopyObjectWithContext(ctx, &s3.CopyObjectInput{
		Bucket:     aws.String(bucket),
		CopySource: aws.String(fmt.Sprintf("%s/%s", s.bucket, src)),
		Key:        aws.String(dst),
	})

	return err
}

// Create for create interface
func (s *Storage) Create(_ string) (io.ReadWriteCloser, error) {
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

// GetWithContext gets file from s3 bucket
func (s *Storage) GetWithContext(ctx aws.Context, path string) (io.ReadCloser, error) {
	out, err := s.s3.GetObjectWithContext(ctx, &s3.GetObjectInput{
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

// PutWithContext puts file into s3 bucket
func (s *Storage) PutWithContext(ctx aws.Context, path string, body io.Reader) error {
	_, err := s.uploader.UploadWithContext(ctx, &s3manager.UploadInput{
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

// DeleteWithContext removes object from s3
func (s *Storage) DeleteWithContext(ctx aws.Context, path string) error {
	_, err := s.s3.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
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

	return &FileInfo{
		size:                      *out.ContentLength,
		acceptRanges:              *out.AcceptRanges,
		activeStatus:              *out.ArchiveStatus,
		cacheControl:              *out.CacheControl,
		contentDisposition:        *out.ContentDisposition,
		contentEncoding:           *out.ContentEncoding,
		contentLanguage:           *out.ContentLanguage,
		contentType:               *out.ContentType,
		eTag:                      *out.ETag,
		expires:                   *out.Expires,
		lastModified:              *out.LastModified,
		bucketKeyEnabled:          *out.BucketKeyEnabled,
		deleteMarker:              *out.DeleteMarker,
		expiration:                *out.Expiration,
		metadata:                  out.Metadata,
		missingMeta:               *out.MissingMeta,
		objectLockLegalHoldStatus: *out.ObjectLockLegalHoldStatus,
		objectLockMode:            *out.ObjectLockMode,
		objectLockRetainUntilDate: *out.ObjectLockRetainUntilDate,
		partsCount:                *out.PartsCount,
		replicationStatus:         *out.ReplicationStatus,
		requestCharged:            *out.RequestCharged,
		restore:                   *out.Restore,
		sseCustomerAlgorithm:      *out.SSECustomerAlgorithm,
		sseCustomerKeyMD5:         *out.SSECustomerKeyMD5,
		sseKMSKeyId:               *out.SSEKMSKeyId,
		serverSideEncryption:      *out.ServerSideEncryption,
		storageClass:              *out.StorageClass,
		versionId:                 *out.VersionId,
		websiteRedirectLocation:   *out.WebsiteRedirectLocation,
	}, nil
}
