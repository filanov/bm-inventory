package s3wrapper

import (
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/go-openapi/swag"
	"github.com/pkg/errors"
)

type Config struct {
	S3EndpointURL      string `envconfig:"S3_ENDPOINT_URL" default:"http://cloudserver-front:8000"`
	Region             string `envconfig:"S3_REGION" default:"us-east-1"`
	S3Bucket           string `envconfig:"S3_BUCKET" default:"test"`
	AwsAccessKeyID     string `envconfig:"AWS_ACCESS_KEY_ID" default:"accessKey1"`
	AwsSecretAccessKey string `envconfig:"AWS_SECRET_ACCESS_KEY" default:"verySecretKey1"`
}

func CreateBucket(cfg *Config) error {
	client, err := NewS3Client(cfg)
	if err != nil {
		return err
	}
	if _, err = client.CreateBucket(&s3.CreateBucketInput{
		Bucket: swag.String(cfg.S3Bucket),
	}); err != nil {
		return errors.Wrapf(err, "failed to create s3 bucket %s", cfg.S3Bucket)
	}
	return nil
}

func NewS3Session(cfg *Config) (*session.Session, error) {
	HTTPTransport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 0,
		MaxIdleConnsPerHost:   4096,
		MaxIdleConns:          0,
		IdleConnTimeout:       time.Minute,
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true}, // true to enable use s3 with ip address (scality)
	}
	creds := credentials.NewStaticCredentials(cfg.AwsAccessKeyID, cfg.AwsSecretAccessKey, "")

	awsConfig := &aws.Config{
		Region:               aws.String(cfg.Region),
		Endpoint:             aws.String(cfg.S3EndpointURL),
		Credentials:          creds,
		DisableSSL:           aws.Bool(true),
		S3ForcePathStyle:     aws.Bool(true),
		S3Disable100Continue: aws.Bool(true),
		HTTPClient:           &http.Client{Transport: HTTPTransport},
	}
	awsSession, err := session.NewSession(awsConfig)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create s3 session")
	}

	return awsSession, nil
}

// NewS3Client creates new s3 client using default config along with defined env variables
func NewS3Client(cfg *Config) (*s3.S3, error) {
	awsSession, err := NewS3Session(cfg)
	if err != nil {
		return nil, err
	}

	client := s3.New(awsSession)
	if client == nil {
		return nil, errors.Errorf("failed to create s3 client")
	}
	return client, nil
}

func PutObject(cfg *Config, file io.Reader, filename string, bucket string) (string, error) {
	s3Session, err := NewS3Session(cfg)
	if err != nil {
		return "", err
	}
	// Create an uploader with the session and custom options
	uploader := s3manager.NewUploader(s3Session, func(u *s3manager.Uploader) {
		u.PartSize = 5 * 1024 * 1024 // The minimum/default allowed part size is 5MB
		u.Concurrency = 5            // default is 5
	})

	// Upload the file to S3.
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
		Body:   file,
	})

	//in case it fails to upload
	if err != nil {
		return "", errors.Errorf("failed to upload %s, error %s", filename, err.Error())
	}
	return result.Location, nil
}
