package s3

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"time"

	"github.com/uswitch/drone-cache/cache"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3Manager"
)

type s3Cache struct {
	client *s3.S3
	bucket string
}

func (c *s3Cache) List(root string) ([]os.FileInfo, error) {
	return nil, errors.New("List not implemented for S3")
}

func (c *s3Cache) Get(p string) (io.ReadCloser, error) {
	resp, err := c.client.GetObject(&s3.GetObjectInput{
		Bucket: &c.bucket,
		Key:    &p,
	})

	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

func (c *s3Cache) Put(p string, t time.Duration, src io.Reader) error {
	uploader := s3manager.NewUploaderWithClient(c.client)

	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: &c.bucket,
		Key:    &p,
		Body:   src,
	})

	return err
}

func (c *s3Cache) Remove(p string) error {
	_, err := c.client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: &c.bucket,
		Key:    &p,
	})

	return err
}

func (c *s3Cache) Close() error {
	return nil
}

func New(bucket string) (cache.Cache, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}

	return &s3Cache{
		client: s3.New(sess),
		bucket: bucket,
	}, nil
}

type configuration struct {
	Bucket string
}

func FromJSON(raw string) (cache.Cache, error) {
	var config configuration

	if err := json.Unmarshal([]byte(raw), &config); err != nil {
		return nil, err
	}

	return New(config.Bucket)
}
