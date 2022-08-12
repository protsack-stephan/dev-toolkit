// Package s3 this package in intention to hide aws s3 storage implementation
// under the interface that will give you the ability to user other cloud providers
// in the future
package s3

import (
	"errors"
	"fmt"
	"io"
	"math"
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

const maxUploadSizeBytes = 4294967296

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

	hr, err := s.s3.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(src),
	})

	if err != nil {
		return err
	}

	if *hr.ContentLength <= maxUploadSizeBytes {
		_, err = s.s3.CopyObject(&s3.CopyObjectInput{
			Bucket:     aws.String(bucket),
			CopySource: aws.String(fmt.Sprintf("%s/%s", s.bucket, src)),
			Key:        aws.String(dst),
		})

		return err
	}

	cmr, err := s.s3.CreateMultipartUpload(&s3.CreateMultipartUploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(dst),
	})

	if err != nil {
		return err
	}

	cmu := &s3.CompletedMultipartUpload{}
	maxPart := int(math.Ceil(float64(*hr.ContentLength) / float64(maxUploadSizeBytes)))

	for prt := 0; prt < maxPart; prt++ {
		from := prt * maxUploadSizeBytes
		to := (prt * maxUploadSizeBytes) + maxUploadSizeBytes

		if prt != 0 {
			from += 1
		}

		if to > int(*hr.ContentLength) {
			to = int(*hr.ContentLength) - 1
		}

		upr, err := s.s3.UploadPartCopy(&s3.UploadPartCopyInput{
			Bucket:          aws.String(bucket),
			CopySource:      aws.String(fmt.Sprintf("%s/%s", bucket, src)),
			CopySourceRange: aws.String(fmt.Sprintf("bytes=%d-%d", from, to)),
			Key:             aws.String(dst),
			PartNumber:      aws.Int64(int64(prt) + 1),
			UploadId:        aws.String(*cmr.UploadId),
		})

		if err != nil {
			return err
		}

		cmu.Parts = append(cmu.Parts, &s3.CompletedPart{
			ETag:       upr.CopyPartResult.ETag,
			PartNumber: aws.Int64(int64(prt) + 1),
		})
	}

	_, err = s.s3.CompleteMultipartUpload(&s3.CompleteMultipartUploadInput{
		Bucket:          aws.String(bucket),
		Key:             aws.String(dst),
		UploadId:        aws.String(*cmr.UploadId),
		MultipartUpload: cmu,
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

	hr, err := s.s3.HeadObjectWithContext(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(src),
	})

	if err != nil {
		return err
	}

	if *hr.ContentLength <= maxUploadSizeBytes {
		_, err = s.s3.CopyObjectWithContext(ctx, &s3.CopyObjectInput{
			Bucket:     aws.String(bucket),
			CopySource: aws.String(fmt.Sprintf("%s/%s", s.bucket, src)),
			Key:        aws.String(dst),
		})

		return err
	}

	cmr, err := s.s3.CreateMultipartUploadWithContext(ctx, &s3.CreateMultipartUploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(dst),
	})

	if err != nil {
		return err
	}

	cmu := &s3.CompletedMultipartUpload{}
	maxPart := int(math.Ceil(float64(*hr.ContentLength) / float64(maxUploadSizeBytes)))

	for prt := 0; prt < maxPart; prt++ {
		from := prt * maxUploadSizeBytes
		to := (prt * maxUploadSizeBytes) + maxUploadSizeBytes

		if prt != 0 {
			from += 1
		}

		if to > int(*hr.ContentLength) {
			to = int(*hr.ContentLength) - 1
		}

		upr, err := s.s3.UploadPartCopyWithContext(ctx, &s3.UploadPartCopyInput{
			Bucket:          aws.String(bucket),
			CopySource:      aws.String(fmt.Sprintf("%s/%s", bucket, src)),
			CopySourceRange: aws.String(fmt.Sprintf("bytes=%d-%d", from, to)),
			Key:             aws.String(dst),
			PartNumber:      aws.Int64(int64(prt) + 1),
			UploadId:        aws.String(*cmr.UploadId),
		})

		if err != nil {
			return err
		}

		cmu.Parts = append(cmu.Parts, &s3.CompletedPart{
			ETag:       upr.CopyPartResult.ETag,
			PartNumber: aws.Int64(int64(prt) + 1),
		})
	}

	_, err = s.s3.CompleteMultipartUploadWithContext(ctx, &s3.CompleteMultipartUploadInput{
		Bucket:          aws.String(bucket),
		Key:             aws.String(dst),
		UploadId:        aws.String(*cmr.UploadId),
		MultipartUpload: cmu,
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

	file := new(FileInfo)

	if out.AcceptRanges != nil {
		file.acceptRanges = *out.AcceptRanges
	}

	if out.ArchiveStatus != nil {
		file.activeStatus = *out.ArchiveStatus
	}

	if out.CacheControl != nil {
		file.cacheControl = *out.CacheControl
	}

	if out.ContentDisposition != nil {
		file.contentDisposition = *out.ContentDisposition
	}

	if out.ContentEncoding != nil {
		file.contentEncoding = *out.ContentEncoding
	}

	if out.ContentLanguage != nil {
		file.contentLanguage = *out.ContentLanguage
	}

	if out.ContentLength != nil {
		file.size = *out.ContentLength
	}

	if out.ContentType != nil {
		file.contentType = *out.ContentType
	}

	if out.ETag != nil {
		file.eTag = *out.ETag
	}

	if out.Expires != nil {
		file.expires = *out.Expires
	}

	if out.Expiration != nil {
		file.expiration = *out.Expiration
	}

	if out.LastModified != nil {
		file.lastModified = *out.LastModified
	}

	if out.BucketKeyEnabled != nil {
		file.bucketKeyEnabled = *out.BucketKeyEnabled
	}

	if out.DeleteMarker != nil {
		file.deleteMarker = *out.DeleteMarker
	}

	if out.Metadata != nil {
		file.metadata = out.Metadata
	}

	if out.MissingMeta != nil {
		file.missingMeta = *out.MissingMeta
	}

	if out.ObjectLockLegalHoldStatus != nil {
		file.objectLockLegalHoldStatus = *out.ObjectLockLegalHoldStatus
	}

	if out.ObjectLockMode != nil {
		file.objectLockMode = *out.ObjectLockMode
	}

	if out.ObjectLockRetainUntilDate != nil {
		file.objectLockRetainUntilDate = *out.ObjectLockRetainUntilDate
	}

	if out.PartsCount != nil {
		file.partsCount = *out.PartsCount
	}

	if out.ReplicationStatus != nil {
		file.replicationStatus = *out.ReplicationStatus
	}

	if out.RequestCharged != nil {
		file.requestCharged = *out.RequestCharged
	}

	if out.Restore != nil {
		file.restore = *out.Restore
	}

	if out.SSECustomerAlgorithm != nil {
		file.sseCustomerAlgorithm = *out.SSECustomerAlgorithm
	}

	if out.SSECustomerKeyMD5 != nil {
		file.sseCustomerKeyMD5 = *out.SSECustomerKeyMD5
	}

	if out.SSEKMSKeyId != nil {
		file.sseKMSKeyId = *out.SSEKMSKeyId
	}

	if out.ServerSideEncryption != nil {
		file.serverSideEncryption = *out.ServerSideEncryption
	}

	if out.StorageClass != nil {
		file.storageClass = *out.StorageClass
	}

	if out.VersionId != nil {
		file.versionId = *out.VersionId
	}

	if out.WebsiteRedirectLocation != nil {
		file.websiteRedirectLocation = *out.WebsiteRedirectLocation
	}

	return file, nil
}
