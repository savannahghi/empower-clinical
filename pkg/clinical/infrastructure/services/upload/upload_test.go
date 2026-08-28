package upload_test

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/savannahghi/empower-clinical/pkg/clinical/infrastructure/services/upload"
)

func TestServiceUploadImpl_UploadMedia(t *testing.T) {
	type args struct {
		ctx         context.Context
		name        string
		file        io.Reader
		contentType string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Happy case: upload media",
			args: args{
				ctx:         context.Background(),
				name:        gofakeit.BeerName(),
				file:        strings.NewReader("example file contents"),
				contentType: "application/json",
			},
			want:    "example file contents",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeUpload := upload.NewServiceUpload(context.Background())

			_, err := fakeUpload.UploadMedia(tt.args.ctx, tt.args.name, tt.args.file, tt.args.contentType)
			if (err != nil) != tt.wantErr {
				t.Errorf("ServiceUploadImpl.UploadMedia() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
