package foundation_test

import (
	"context"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/scalarutils"
	"github.com/savannahghi/silurlshortener"
	"github.com/stretchr/testify/mock"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/empower-clinical/pkg/clinical/domain"
	usecaseMock "github.com/savannahghi/empower-clinical/pkg/clinical/usecases/mock"
)

func TestUseCasesClinicalImpl_UploadMedia(t *testing.T) {
	ctx := usecaseMock.AddTenantIdentifierContext(context.Background())
	ID := uuid.NewString()
	encounter := &domain.FHIREncounterRelayPayload{
		Resource: &domain.FHIREncounter{
			ID: &ID,
			Subject: &domain.FHIRReference{
				ID:      &ID,
				Display: ID,
			},
		},
	}

	orgName := mock.Anything
	organization := &domain.FHIROrganizationRelayPayload{
		Resource: &domain.FHIROrganization{
			ID:   &ID,
			Name: &orgName,
		},
	}

	patient := domain.FHIRPatientRelayPayload{
		Resource: &domain.FHIRPatient{
			ID: &ID,
			Name: []*domain.FHIRHumanName{
				{
					Text: gofakeit.Name(),
				},
			},
		},
	}

	docref := &domain.FHIRDocumentReference{
		ID: ID,
	}

	media := &dto.Media{
		ID:          ID,
		ContentType: "pdf",
		SignedURL:   gofakeit.URL(),
		Name:        gofakeit.Name(),
	}

	tenantIDs := &dto.TenantIdentifiers{
		OrganizationID: ID,
		FacilityID:     ID,
	}

	type args struct {
		ctx              context.Context
		encounterID      string
		serviceRequestID string
		file             io.Reader
		contentType      string
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "happy case: upload media",
			setup: func(mh *usecaseMock.Mocks) args {
				ID := uuid.NewString()

				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return organization, nil
					})
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return &patient, nil
					})
				mh.Upload.EXPECT().UploadMedia(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, name string, file io.Reader, contentType string) (*dto.Media, error) {
						return media, nil
					})
				mh.URLShortener.EXPECT().Shorten(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, payload *silurlshortener.ShortenURLPayload) (*silurlshortener.ShortenURLResponse, error) {
						return &silurlshortener.ShortenURLResponse{
							ShortURL: "test",
						}, nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().CreateFHIRDocumentReference(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input *domain.FHIRDocumentReferenceInput) (*domain.FHIRDocumentReference, error) {
						return docref, nil
					})

				return args{ctx: ctx, encounterID: ID, serviceRequestID: ID, file: strings.NewReader("test"), contentType: "application/json"}
			},
			wantErr: false,
		},
		{
			name: "sad case: unable to get encounter",
			setup: func(mh *usecaseMock.Mocks) args {
				ID := uuid.NewString()

				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return organization, nil
					})
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, encounterID: ID, serviceRequestID: ID, file: strings.NewReader("test"), contentType: "application/json"}
			},
			wantErr: true,
		},
		{
			name: "sad case: unable to upload media",
			setup: func(mh *usecaseMock.Mocks) args {
				ID := uuid.NewString()

				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return organization, nil
					})
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return &patient, nil
					})
				mh.Upload.EXPECT().UploadMedia(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, name string, file io.Reader, contentType string) (*dto.Media, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, encounterID: ID, serviceRequestID: ID, file: strings.NewReader("test"), contentType: "application/json"}
			},
			wantErr: true,
		},
		{
			name: "sad case: unable to create fhir media",
			setup: func(mh *usecaseMock.Mocks) args {
				ID := uuid.NewString()

				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return organization, nil
					})
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return &patient, nil
					})
				mh.Upload.EXPECT().UploadMedia(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, name string, file io.Reader, contentType string) (*dto.Media, error) {
						return media, nil
					})
				mh.URLShortener.EXPECT().Shorten(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, payload *silurlshortener.ShortenURLPayload) (*silurlshortener.ShortenURLResponse, error) {
						return &silurlshortener.ShortenURLResponse{
							ShortURL: "test",
						}, nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().CreateFHIRDocumentReference(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input *domain.FHIRDocumentReferenceInput) (*domain.FHIRDocumentReference, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, encounterID: ID, serviceRequestID: ID, file: strings.NewReader("test"), contentType: "application/json"}
			},
			wantErr: true,
		},
		{
			name: "sad case: missing facility identifier in context",
			setup: func(mh *usecaseMock.Mocks) args {
				ID := uuid.NewString()
				return args{ctx: context.Background(), encounterID: ID, serviceRequestID: ID, file: strings.NewReader("test"), contentType: "application/json"}
			},
			wantErr: true,
		},
		{
			name: "sad case: unable to get fhir organisation",
			setup: func(mh *usecaseMock.Mocks) args {
				ID := uuid.NewString()

				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, encounterID: ID, serviceRequestID: ID, file: strings.NewReader("test"), contentType: "application/json"}
			},
			wantErr: true,
		},
		{
			name: "sad case: unable to get patient",
			setup: func(mh *usecaseMock.Mocks) args {
				ID := uuid.NewString()

				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return organization, nil
					})
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, encounterID: ID, serviceRequestID: ID, file: strings.NewReader("test"), contentType: "application/json"}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get tenant meta tags",
			setup: func(mh *usecaseMock.Mocks) args {
				ID := uuid.NewString()

				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return organization, nil
					})
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return &patient, nil
					})
				mh.Upload.EXPECT().UploadMedia(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, name string, file io.Reader, contentType string) (*dto.Media, error) {
						return media, nil
					})
				mh.URLShortener.EXPECT().Shorten(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, payload *silurlshortener.ShortenURLPayload) (*silurlshortener.ShortenURLResponse, error) {
						return &silurlshortener.ShortenURLResponse{
							ShortURL: "test",
						}, nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, encounterID: ID, serviceRequestID: ID, file: strings.NewReader("test"), contentType: "application/json"}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to shorten url",
			setup: func(mh *usecaseMock.Mocks) args {
				ID := uuid.NewString()

				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return organization, nil
					})
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return &patient, nil
					})
				mh.Upload.EXPECT().UploadMedia(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, name string, file io.Reader, contentType string) (*dto.Media, error) {
						return media, nil
					})
				mh.URLShortener.EXPECT().Shorten(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, payload *silurlshortener.ShortenURLPayload) (*silurlshortener.ShortenURLResponse, error) {
						return nil, fmt.Errorf("error")
					})

				return args{ctx: ctx, encounterID: ID, serviceRequestID: ID, file: strings.NewReader("test"), contentType: "application/json"}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			_, err := clinicalUsecase.UploadMedia(args.ctx, args.encounterID, args.serviceRequestID, args.file, args.contentType)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.UploadMedia() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCasesClinicalImpl_ListPatientMedia(t *testing.T) {
	ctx := usecaseMock.AddTenantIdentifierContext(context.Background())
	first := 10
	pagination := dto.Pagination{
		First: &first,
	}

	ID := uuid.NewString()
	encounter := &domain.FHIREncounterRelayPayload{
		Resource: &domain.FHIREncounter{
			ID: &ID,
			Subject: &domain.FHIRReference{
				ID:      &ID,
				Display: ID,
			},
		},
	}

	tenantIDs := &dto.TenantIdentifiers{
		OrganizationID: ID,
		FacilityID:     ID,
	}

	resourceID := gofakeit.UUID()
	contentType := "application/pdf"
	title := "Test Title"
	url := mock.Anything
	pagedFHIRDocumentRef := &domain.PagedFHIRDocumentReference{
		DocumentReferences: []domain.FHIRDocumentReference{
			{
				ID: resourceID,
				Context: []*domain.FHIRReference{
					{
						Reference: new(string),
						ID:        new(string),
					},
				},
				Subject: &domain.FHIRReference{
					ID: &resourceID,
				},
				Content: []domain.FHIRDocumentReferenceContent{
					{
						ID: resourceID,
						Attachment: domain.FHIRAttachment{
							URL:         (*scalarutils.URL)(&url),
							ContentType: (*scalarutils.Code)(&contentType),
							Title:       &title,
						},
					},
				},
			},
		},
		HasNextPage:     false,
		NextCursor:      "",
		HasPreviousPage: false,
		PreviousCursor:  "",
		TotalCount:      0,
	}

	urlresponse := &silurlshortener.ShortenURLResponse{
		ShortURL: mock.Anything,
	}

	type args struct {
		ctx                         context.Context
		patientID, serviceRequestID string
		pagination                  dto.Pagination
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy case: list patient media",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchFHIRDocumentReference(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, searchParams map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRDocumentReference, error) {
						return pagedFHIRDocumentRef, nil
					})
				mh.URLShortener.EXPECT().Shorten(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, payload *silurlshortener.ShortenURLPayload) (*silurlshortener.ShortenURLResponse, error) {
						return urlresponse, nil
					})

				return args{ctx: ctx, patientID: fmt.Sprintf("Patient/%s", gofakeit.UUID()), serviceRequestID: gofakeit.UUID(), pagination: pagination}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - invalid pagination",
			setup: func(mh *usecaseMock.Mocks) args {
				pagination := dto.Pagination{
					First: &first,
					Last:  &first,
				}

				return args{ctx: ctx, patientID: fmt.Sprintf("Patient/%s", gofakeit.UUID()), serviceRequestID: gofakeit.UUID(), pagination: pagination}
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get tenant identifiers",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, patientID: fmt.Sprintf("Patient/%s", gofakeit.UUID()), serviceRequestID: gofakeit.UUID(), pagination: pagination}
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to list patient media",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchFHIRDocumentReference(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, searchParams map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRDocumentReference, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, patientID: fmt.Sprintf("Patient/%s", gofakeit.UUID()), serviceRequestID: gofakeit.UUID(), pagination: pagination}
			},
			wantErr: true,
		},
		{
			name: "sad case: unable to get encounter",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, patientID: fmt.Sprintf("Patient/%s", gofakeit.UUID()), serviceRequestID: gofakeit.UUID(), pagination: pagination}
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to shorten media url",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchFHIRDocumentReference(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, searchParams map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRDocumentReference, error) {
						return pagedFHIRDocumentRef, nil
					})
				mh.URLShortener.EXPECT().Shorten(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, payload *silurlshortener.ShortenURLPayload) (*silurlshortener.ShortenURLResponse, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, patientID: fmt.Sprintf("Patient/%s", gofakeit.UUID()), serviceRequestID: gofakeit.UUID(), pagination: pagination}
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			_, err := clinicalUsecase.ListPatientMedia(args.ctx, args.patientID, args.serviceRequestID, args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.ListPatientMedia() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
