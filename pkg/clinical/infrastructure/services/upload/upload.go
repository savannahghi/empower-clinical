package upload

import (
	"context"
	"io"
	"net/http"
	"time"

	"cloud.google.com/go/storage"
	"github.com/savannahghi/serverutils"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/dto"
	"google.golang.org/api/option"
)

// ServiceUpload holds the upload service methods
type ServiceUpload interface {
	UploadMedia(ctx context.Context, name string, file io.Reader, contentType string) (*dto.Media, error)
}

// ServiceUploadImpl represents upload service implementations
type ServiceUploadImpl struct {
	Client storage.Client
}

// NewServiceUpload returns new instance of upload service
func NewServiceUpload(ctx context.Context) *ServiceUploadImpl {
	credentials := serverutils.MustGetEnvVar("GOOGLE_APPLICATION_CREDENTIALS")

	client, err := storage.NewClient(ctx, option.WithAuthCredentialsFile(option.ServiceAccount, credentials))
	if err != nil {
		panic(err)
	}

	defer client.Close()

	return &ServiceUploadImpl{
		Client: *client,
	}
}

// UploadMedia uploads media to GCS
func (u *ServiceUploadImpl) UploadMedia(ctx context.Context, name string, file io.Reader, contentType string) (*dto.Media, error) {
	bucketName := serverutils.MustGetEnvVar("CLINICAL_BUCKET_NAME")

	object := u.Client.Bucket(bucketName).Object(name)

	wc := object.NewWriter(ctx)
	wc.ContentType = contentType
	wc.ChunkSize = 256 * 1024 // 256 KB chunk size

	if _, err := io.Copy(wc, file); err != nil {
		wc.Close()
		return nil, err
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 50*time.Second)
	defer cancel()

	if err := wc.Close(); err != nil {
		return nil, err
	}

	url, err := object.Attrs(timeoutCtx)
	if err != nil {
		return nil, err
	}

	signedURL, err := u.Client.Bucket(bucketName).SignedURL(name, &storage.SignedURLOptions{
		Method:  http.MethodGet,
		Scheme:  storage.SigningSchemeV2,
		Expires: time.Now().Add(24 * time.Hour * 365 * 100),
	})
	if err != nil {
		return nil, err
	}

	output := &dto.Media{
		MediaLink:   url.MediaLink,
		Name:        url.Name,
		ContentType: url.ContentType,
		SignedURL:   signedURL,
	}

	return output, nil
}
