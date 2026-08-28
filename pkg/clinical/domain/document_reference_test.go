package domain

import (
	"github.com/savannahghi/scalarutils"
	"testing"
)

func TestFHIRDocumentReference_GetDocumentType(t *testing.T) {
	code := "161360"

	type fields struct {
		Type *FHIRCodeableConcept
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Happy case: get code",
			fields: fields{
				Type: &FHIRCodeableConcept{
					ID: new(string),
					Coding: []*FHIRCoding{
						{
							Code: (*scalarutils.Code)(&code),
						},
					},
					Text: "",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &FHIRDocumentReference{
				Type: tt.fields.Type,
			}
			_, err := d.GetDocumentType()
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRDocumentReference.GetDocumentType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFHIRDocumentReference_GetDocumentAttachment(t *testing.T) {
	contentType := scalarutils.Code("application/json")
	url := scalarutils.URL("http://example.com")
	title := "Patient Referral Document"

	type fields struct {
		Content []FHIRDocumentReferenceContent
	}
	tests := []struct {
		name    string
		fields  fields
		want    *FHIRAttachment
		wantErr bool
	}{
		{
			name: "Happy case: get document attachment",
			fields: fields{
				Content: []FHIRDocumentReferenceContent{
					{
						Attachment: FHIRAttachment{
							ContentType: &contentType,
							URL:         &url,
							Title:       &title,
						},
					},
				},
			},
			wantErr: false,
			want: &FHIRAttachment{
				ContentType: &contentType,
				URL:         &url,
				Title:       &title,
			},
		},
		{
			name: "Sad case: content is empty",
			fields: fields{
				Content: []FHIRDocumentReferenceContent{},
			},
			wantErr: true,
			want:    nil,
		},
		{
			name:    "Sad case: document reference is nil",
			fields:  fields{nil},
			wantErr: true,
			want:    nil,
		},
		{
			name: "Sad case: no attachment URL found",
			fields: fields{
				Content: []FHIRDocumentReferenceContent{
					{
						Attachment: FHIRAttachment{
							ContentType: &contentType,
							URL:         nil,
							Title:       &title,
						},
					},
				},
			},
			wantErr: true,
			want:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var d *FHIRDocumentReference
			if tt.fields.Content != nil {
				d = &FHIRDocumentReference{
					Content: tt.fields.Content,
				}
			}
			_, err := d.GetDocumentAttachment()
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRDocumentReference.GetDocumentAttachment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
