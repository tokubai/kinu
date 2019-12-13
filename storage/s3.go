package storage

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	awsSession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/tokubai/kinu/logger"
)

type S3Storage struct {
	Storage

	client *s3.S3

	region         string
	bucket         string
	bucketBasePath string
}

type S3StorageItem struct {
	StorageItem

	Object *s3.Object
}

func openS3Storage() (Storage, error) {
	s := &S3Storage{}
	err := s.Open()
	if err != nil {
		return nil, logger.ErrorDebug(err)
	}
	return s, nil
}

func (s *S3Storage) Open() error {
	s.region = os.Getenv("KINU_S3_REGION")
	if len(s.region) == 0 {
		return &ErrInvalidStorageOption{Message: "KINU_S3_REGION system env is required"}
	}

	s.bucket = os.Getenv("KINU_S3_BUCKET")
	if len(s.bucket) == 0 {
		return &ErrInvalidStorageOption{Message: "KINU_S3_BUCKET system env is required"}
	}

	s.bucketBasePath = os.Getenv("KINU_S3_BUCKET_BASE_PATH")

	s.client = s3.New(awsSession.New(), &aws.Config{Region: aws.String(s.region)})

	logger.WithFields(logrus.Fields{
		"bucket":    s.bucket,
		"base_path": s.bucketBasePath,
		"region":    s.region,
	}).Debug("open s3 storage")

	return nil
}

func (s *S3Storage) BuildKey(key string) string {
	if s.bucketBasePath == "/" {
		return key
	} else if len(s.bucketBasePath) == 0 {
		return key
	} else if strings.HasSuffix(s.bucketBasePath, "/") {
		return s.bucketBasePath + key
	} else {
		return s.bucketBasePath + "/" + key
	}
}

func (s *S3Storage) Fetch(key string) (*Object, error) {
	key = s.BuildKey(key)

	params := &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}

	logger.WithFields(logrus.Fields{
		"bucket": s.bucket,
		"key":    key,
	}).Debug("start get object from s3")

	resp, err := s.client.GetObject(params)

	if reqerr, ok := err.(awserr.RequestFailure); ok && reqerr.StatusCode() == http.StatusNotFound {
		return nil, ErrImageNotFound
	} else if err != nil {
		return nil, logger.ErrorDebug(err)
	}

	logger.WithFields(logrus.Fields{
		"bucket": s.bucket,
		"key":    key,
	}).Debug("found object from s3")

	defer resp.Body.Close()

	object := &Object{
		Metadata: make(map[string]string, 0),
	}
	for k, v := range resp.Metadata {
		object.Metadata[k] = *v
	}
	object.Body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, logger.ErrorDebug(err)
	}

	return object, nil
}

func (s *S3Storage) PutFromBlob(key string, image []byte, contentType string, metadata map[string]string) error {
	tmpfile, err := ioutil.TempFile("", "kinu-upload")
	if err != nil {
		return logger.ErrorDebug(err)
	}
	_, err = tmpfile.Write(image)
	if err != nil {
		return logger.ErrorDebug(err)
	}

	defer func() {
		tmpfile.Close()
		os.Remove(tmpfile.Name())
	}()

	return s.Put(key, tmpfile, contentType, metadata)
}

func (s *S3Storage) Put(key string, imageFile io.ReadSeeker, contentType string, metadata map[string]string) error {
	putMetadata := make(map[string]*string, 0)
	for k, v := range metadata {
		putMetadata[k] = aws.String(v)
	}

	_, err := imageFile.Seek(0, 0)
	if err != nil {
		return logger.ErrorDebug(err)
	}

	_, err = s.client.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(s.BuildKey(key)),
		ContentType: aws.String(contentType),
		Body:        imageFile,
		Metadata:    putMetadata,
	})

	logger.WithFields(logrus.Fields{
		"bucket": s.bucket,
		"key":    s.BuildKey(key),
	}).Debug("put to s3")

	if err != nil {
		return logger.ErrorDebug(err)
	}

	return nil
}

func (s *S3Storage) List(key string) ([]StorageItem, error) {
	resp, err := s.client.ListObjects(&s3.ListObjectsInput{
		Bucket: aws.String(s.bucket),
		Prefix: aws.String(s.BuildKey(key)),
	})

	if err != nil {
		return nil, logger.ErrorDebug(err)
	}

	logger.WithFields(logrus.Fields{
		"bucket": s.bucket,
		"key":    s.BuildKey(key),
	}).Debug("start list object from s3")

	items := make([]StorageItem, 0)
	for _, object := range resp.Contents {
		logger.WithFields(logrus.Fields{
			"key": &object.Key,
		}).Debug("found object")
		item := S3StorageItem{Object: object}
		items = append(items, &item)
	}

	return items, nil
}

func (s *S3Storage) Move(from string, to string) error {
	fromKey := s.bucket + "/" + from
	toKey := s.bucketBasePath + "/" + to

	_, err := s.client.CopyObject(&s3.CopyObjectInput{
		Bucket:     aws.String(s.bucket),
		CopySource: aws.String(fromKey),
		Key:        aws.String(toKey),
	})

	logger.WithFields(logrus.Fields{
		"from": fromKey,
		"to":   toKey,
	}).Debug("move s3 object start")

	if reqerr, ok := err.(awserr.RequestFailure); ok && reqerr.StatusCode() == http.StatusNotFound {
		return ErrImageNotFound
	} else if err != nil {
		return logger.ErrorDebug(err)
	}
	_, err = s.client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(from),
	})

	if reqerr, ok := err.(awserr.RequestFailure); ok && reqerr.StatusCode() == http.StatusNotFound {
		return ErrImageNotFound
	} else if err != nil {
		return logger.ErrorDebug(err)
	}

	return nil
}

func (s *S3StorageItem) IsValid() bool {
	if len(s.Extension()) == 0 {
		return false
	}

	if len(s.ImageSize()) == 0 {
		return false
	}

	return true
}

func (s *S3StorageItem) Key() string {
	return *s.Object.Key
}

func (s *S3StorageItem) Filename() string {
	path := strings.Split(s.Key(), "/")
	return path[len(path)-1]
}

func (s *S3StorageItem) Extension() string {
	path := strings.Split(*s.Object.Key, ".")
	return path[len(path)-1]
}

// KeyFormat: :image_type/:id/:id.:size.:format
func (s *S3StorageItem) ImageSize() string {
	path := strings.Split(s.Key(), ".")
	return path[len(path)-2]
}
