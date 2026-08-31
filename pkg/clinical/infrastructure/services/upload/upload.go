package upload

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/dto"
)

// ServiceUpload holds the upload service methods
type ServiceUpload interface {
	UploadMedia(ctx context.Context, name string, file io.Reader, contentType string) (*dto.Media, error)
}

// ServiceUploadImpl stores media in any S3-compatible object store.
//
// This was Google Cloud Storage. It is now S3-compatible so the service can run
// against MinIO, Ceph, AWS or anything else that speaks the protocol, and so that
// an operator is not required to hold Google Cloud credentials to start the
// process. Configuration is by environment:
//
//	STORAGE_ENDPOINT     base URL of the object store. Set for MinIO
//	                     (http://localhost:9000); leave unset for real AWS S3.
//	STORAGE_ACCESS_KEY   static credentials. When unset the default AWS chain is
//	STORAGE_SECRET_KEY   used instead (instance role, shared config, environment).
//	STORAGE_REGION       defaults to us-east-1, which MinIO accepts.
//	CLINICAL_BUCKET_NAME the bucket. Created on first use if absent.
type ServiceUploadImpl struct {
	client   *s3.Client
	presign  *s3.PresignClient
	bucket   string
	endpoint string
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return fallback
}

// NewServiceUpload returns a new instance of the upload service.
//
// Unlike the Google Cloud implementation this replaced, it does not panic when
// storage is unconfigured: building an S3 client performs no I/O and validates
// no credentials, so a misconfiguration surfaces on the first upload rather than
// preventing the process from starting at all.
func NewServiceUpload(ctx context.Context) *ServiceUploadImpl {
	region := env("STORAGE_REGION", "us-east-1")
	endpoint := os.Getenv("STORAGE_ENDPOINT")

	opts := []func(*awsconfig.LoadOptions) error{awsconfig.WithRegion(region)}

	if key, secret := os.Getenv("STORAGE_ACCESS_KEY"), os.Getenv("STORAGE_SECRET_KEY"); key != "" && secret != "" {
		opts = append(opts, awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(key, secret, ""),
		))
	}

	cfg, err := awsconfig.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		// Deliberately not fatal. Storage is used by two endpoints; the rest of
		// the service is unaffected and should still come up.
		cfg = aws.Config{Region: region}
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		if endpoint != "" {
			o.BaseEndpoint = aws.String(endpoint)
			// MinIO and most self-hosted stores do not do virtual-host addressing.
			o.UsePathStyle = true
		}
	})

	return &ServiceUploadImpl{
		client:   client,
		presign:  s3.NewPresignClient(client),
		bucket:   env("CLINICAL_BUCKET_NAME", "clinical-media"),
		endpoint: strings.TrimSuffix(endpoint, "/"),
	}
}

// ensureBucket creates the bucket when it does not already exist. Real S3
// deployments normally provision the bucket out of band, but a local or
// first-run deployment has nothing to upload into otherwise.
func (u *ServiceUploadImpl) ensureBucket(ctx context.Context) error {
	_, err := u.client.HeadBucket(ctx, &s3.HeadBucketInput{Bucket: aws.String(u.bucket)})
	if err == nil {
		return nil
	}

	var notFound *types.NotFound
	var noSuchBucket *types.NoSuchBucket

	if !errors.As(err, &notFound) && !errors.As(err, &noSuchBucket) {
		return err
	}

	_, err = u.client.CreateBucket(ctx, &s3.CreateBucketInput{Bucket: aws.String(u.bucket)})

	return err
}

// UploadMedia stores the reader's contents and returns a link to it.
func (u *ServiceUploadImpl) UploadMedia(
	ctx context.Context,
	name string,
	file io.Reader,
	contentType string,
) (*dto.Media, error) {
	if err := u.ensureBucket(ctx); err != nil {
		return nil, fmt.Errorf("unable to prepare bucket %q: %w", u.bucket, err)
	}

	body, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("unable to read media: %w", err)
	}

	input := &s3.PutObjectInput{
		Bucket: aws.String(u.bucket),
		Key:    aws.String(name),
		Body:   strings.NewReader(string(body)),
	}
	if contentType != "" {
		input.ContentType = aws.String(contentType)
	}

	if _, err := u.client.PutObject(ctx, input); err != nil {
		return nil, fmt.Errorf("unable to upload media %q: %w", name, err)
	}

	// The Google implementation signed for a hundred years. SigV4 caps a
	// presigned URL at seven days, so callers that persist this URL must be able
	// to re-sign. Referral reports shorten it before storing, which hides the
	// expiry behind a stable short link.
	signed, err := u.presign.PresignGetObject(ctx,
		&s3.GetObjectInput{Bucket: aws.String(u.bucket), Key: aws.String(name)},
		s3.WithPresignExpires(7*24*time.Hour),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to sign url for %q: %w", name, err)
	}

	return &dto.Media{
		MediaLink:   fmt.Sprintf("%s/%s/%s", u.endpoint, u.bucket, name),
		Name:        name,
		ContentType: contentType,
		SignedURL:   signed.URL,
	}, nil
}
