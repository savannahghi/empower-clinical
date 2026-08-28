package clinical_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
	"github.com/savannahghi/scalarutils"
	"github.com/savannahghi/silurlshortener"
	"github.com/stretchr/testify/mock"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/common"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/empower-clinical/pkg/clinical/domain"
	"github.com/savannahghi/empower-clinical/pkg/clinical/usecases/clinical"
	usecaseMock "github.com/savannahghi/empower-clinical/pkg/clinical/usecases/mock"
)

// TODO: Investigate why this function makes a real http request

// func TestUseCasesClinicalImpl_GenerateReferralReportPDF(t *testing.T) {
// 	ctx := usecaseMock.AddTenantIdentifierContext(context.Background())

// 	status := dto.ObservationStatusFinal
// 	valueConcept := "222"
// 	UUID := gofakeit.UUID()
// 	instant := "2023-05-28T14:20:00Z"
// 	currentTime := time.Now().Format(time.RFC3339)
// 	pagedFHIRObservations := &domain.PagedFHIRObservations{
// 		Observations: []domain.FHIRObservation{
// 			{
// 				ID:     new(string),
// 				Status: (*domain.ObservationStatusEnum)(&status),
// 				Code: &domain.FHIRCodeableConcept{
// 					ID: new(string),
// 					Coding: []*domain.FHIRCoding{{
// 						Display: gofakeit.BS(),
// 						Code:    (*scalarutils.Code)(&valueConcept),
// 					}},
// 				},
// 				EffectiveDateTime: (*scalarutils.DateTime)(&currentTime),
// 				Subject: &domain.FHIRReference{
// 					ID: new(string),
// 				},
// 				Encounter: &domain.FHIRReference{
// 					ID: new(string),
// 				},
// 				ValueQuantity: &domain.FHIRQuantity{
// 					Value: 100,
// 					Unit:  "cm",
// 				},
// 				ValueCodeableConcept: &domain.FHIRCodeableConcept{
// 					Text: valueConcept,
// 				},
// 				ValueString:  new(string),
// 				ValueBoolean: new(bool),
// 				ValueInteger: new(string),
// 				ValueRange: &domain.FHIRRange{
// 					Low: domain.FHIRQuantity{
// 						Value: 100,
// 						Unit:  "cm",
// 					},
// 					High: domain.FHIRQuantity{
// 						Value: 100,
// 						Unit:  "cm",
// 					},
// 				},
// 				ValueRatio: &domain.FHIRRatio{
// 					Numerator: domain.FHIRQuantity{
// 						Value: 100,
// 						Unit:  "cm",
// 					},
// 					Denominator: domain.FHIRQuantity{
// 						Value: 100,
// 						Unit:  "cm",
// 					},
// 				},
// 				ValueSampledData: &domain.FHIRSampledData{
// 					ID: &UUID,
// 				},
// 				ValueTime: &time.Time{},
// 				ValueDateTime: &scalarutils.Date{
// 					Year:  2000,
// 					Month: 1,
// 					Day:   1,
// 				},
// 				ValuePeriod: &domain.FHIRPeriod{
// 					Start: scalarutils.DateTime(time.Wednesday.String()),
// 					End:   scalarutils.DateTime(time.Thursday.String()),
// 				},
// 				EffectiveInstant: (*scalarutils.Instant)(&instant),
// 			},
// 		},
// 		HasNextPage:     false,
// 		NextPageURL:     "",
// 		HasPreviousPage: false,
// 		PreviousPageURL: "",
// 		TotalCount:      0,
// 	}

// 	noteText := scalarutils.Markdown(gofakeit.BeerName())
// 	serviceCode := scalarutils.Code("TEST")
// 	ref := mock.Anything
// 	servicerequest := &domain.FHIRServiceRequestRelayPayload{
// 		Resource: &domain.FHIRServiceRequest{
// 			ID:         &UUID,
// 			Text:       &domain.FHIRNarrative{},
// 			Identifier: []*domain.FHIRIdentifier{},
// 			Status:     domain.ServiceRequestStatusActive,
// 			Intent:     domain.ServiceRequestIntentDirective,
// 			Subject: &domain.FHIRReference{
// 				ID:        &UUID,
// 				Reference: &ref,
// 				Display:   mock.Anything,
// 			},
// 			Encounter: &domain.FHIRReference{
// 				Display:   mock.Anything,
// 				Reference: &ref,
// 			},
// 			Note: []*domain.FHIRAnnotation{
// 				{
// 					Text: &noteText,
// 				},
// 			},
// 			Code: &domain.FHIRCodeableReference{
// 				Concept: &domain.FHIRCodeableConcept{
// 					Coding: []*domain.FHIRCoding{
// 						{
// 							Code: &serviceCode,
// 						},
// 					},
// 				},
// 			},
// 			Performer: []*domain.FHIRReference{
// 				{
// 					Display: gofakeit.Name(),
// 				},
// 			},
// 		},
// 	}

// 	birthDate := scalarutils.Date{Year: 2010, Month: 12, Day: 1}
// 	telecom := gofakeit.Phone()
// 	healthSystem := scalarutils.URI("HEALTH_ID")
// 	nationalSystem := scalarutils.URI("NATIONAL_ID")
// 	patientCode := scalarutils.Code(gofakeit.Name())
// 	gender := domain.PatientGenderEnumFemale
// 	patient := &domain.FHIRPatientRelayPayload{
// 		Resource: &domain.FHIRPatient{
// 			ID:        &UUID,
// 			BirthDate: &birthDate,
// 			Name: []*domain.FHIRHumanName{
// 				{
// 					Text: gofakeit.Name(),
// 				},
// 			},
// 			Telecom: []*domain.FHIRContactPoint{
// 				{
// 					Value: &telecom,
// 				},
// 			},
// 			Identifier: []*domain.FHIRIdentifier{
// 				{
// 					Type: domain.FHIRCodeableConcept{
// 						Coding: []*domain.FHIRCoding{
// 							{
// 								System: &healthSystem,
// 								Code:   &patientCode,
// 							},
// 							{
// 								System: &nationalSystem,
// 								Code:   &patientCode,
// 							},
// 						},
// 					},
// 				},
// 			},
// 			Gender: &gender,
// 		},
// 	}

// 	ID := uuid.NewString()
// 	orgName := mock.Anything
// 	orgPhone := gofakeit.Phone()
// 	phoneSystem := domain.ContactPointSystemEnumPhone
// 	emailSystem := domain.ContactPointSystemEnumEmail
// 	orgEmail := gofakeit.Email()
// 	organization := &domain.FHIROrganizationRelayPayload{
// 		Resource: &domain.FHIROrganization{
// 			ID:   &ID,
// 			Name: &orgName,
// 			Contact: []domain.FHIROrganizationContact{
// 				{
// 					Telecom: []domain.FHIRContactPoint{
// 						{
// 							System: &phoneSystem,
// 							Value:  &orgPhone,
// 						},
// 						{
// 							System: &emailSystem,
// 							Value:  &orgEmail,
// 						},
// 					},
// 				},
// 			},
// 		},
// 	}

// 	urlresponse := &silurlshortener.ShortenURLResponse{
// 		ShortURL: mock.Anything,
// 	}

// 	implicitRules := mock.Anything
// 	language := mock.Anything
// 	documentRef := &domain.FHIRDocumentReference{
// 		ID:            uuid.NewString(),
// 		Status:        domain.DocumentReferenceStatusEnumCurrent,
// 		ImplicitRules: &implicitRules,
// 		Language:      &language,
// 		Meta: &domain.FHIRMeta{
// 			VersionID:   mock.Anything,
// 			Source:      mock.Anything,
// 			LastUpdated: gofakeit.Date(),
// 		},
// 		Text:      &domain.FHIRNarrative{},
// 		Extension: []domain.FHIRExtension{},
// 	}

// 	conceptpayload := &domain.Concept{
// 		ConceptClass: mock.Anything,
// 		DataType:     mock.Anything,
// 		ID:           gofakeit.UUID(),
// 	}

// 	tenantIDs := &dto.TenantIdentifiers{
// 		OrganizationID: uuid.NewString(),
// 		FacilityID:     uuid.NewString(),
// 	}

// 	media := &dto.Media{
// 		ID:          uuid.NewString(),
// 		PatientID:   uuid.NewString(),
// 		PatientName: gofakeit.Name(),
// 		MediaLink:   mock.Anything,
// 		Name:        gofakeit.Name(),
// 		SignedURL:   gofakeit.URL(),
// 		ContentType: "pdf",
// 	}

// 	type args struct {
// 		ctx              context.Context
// 		serviceRequestID string
// 	}
// 	tests := []struct {
// 		name    string
// 		setup   func(mh *usecaseMock.Mocks) args
// 		wantErr bool
// 	}{
// 		{
// 			name: "Happy Case - Successfully generate a referral report pdf",
// 			setup: func(mh *usecaseMock.Mocks) args {
// 				mh.FHIR.EXPECT().GetFHIRServiceRequest(mock.Anything, mock.Anything).
// 					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRServiceRequestRelayPayload, error) {
// 						return servicerequest, nil
// 					})
// 				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
// 					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
// 						return patient, nil
// 					})
// 				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
// 					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
// 						return organization, nil
// 					}).Twice()
// 				mh.FHIR.EXPECT().SearchPatientObservations(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
// 					RunAndReturn(func(ctx context.Context, bundleID string, searchParameters map[string]interface{}, tenant dto.TenantIdentifiers, pagination serverutils.PaginationInput) (*domain.PagedFHIRObservations, error) {
// 						return pagedFHIRObservations, nil
// 					})
// 				mh.URLShortener.EXPECT().Shorten(mock.Anything, mock.Anything).
// 					RunAndReturn(func(ctx context.Context, payload *silurlshortener.ShortenURLPayload) (*silurlshortener.ShortenURLResponse, error) {
// 						return urlresponse, nil
// 					})
// 				mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
// 					RunAndReturn(func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
// 						return conceptpayload, nil
// 					})
// 				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
// 					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
// 						return tenantIDs, nil
// 					})
// 				mh.FHIR.EXPECT().CreateFHIRDocumentReference(mock.Anything, mock.Anything).
// 					RunAndReturn(func(ctx context.Context, input *domain.FHIRDocumentReferenceInput) (*domain.FHIRDocumentReference, error) {
// 						return documentRef, nil
// 					})
// 				mh.Upload.EXPECT().UploadMedia(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
// 					RunAndReturn(func(ctx context.Context, name string, file io.Reader, contentType string) (*dto.Media, error) {
// 						return media, nil
// 					})

// 				return args{ctx: ctx, serviceRequestID: uuid.NewString()}
// 			},
// 			wantErr: true,
// 		},
// 			{
// 				name: "Sad Case - Missing service request ID",
// 				setup: func(mh *usecaseMock.Mocks) args {
// 					return args{ctx: ctx}
// 				},
// 				wantErr: true,
// 			},
// 			{
// 				name: "Sad Case - Fail to get service request",
// 				setup: func(mh *usecaseMock.Mocks) args {
// 					mh.FHIR.EXPECT().GetFHIRServiceRequest(mock.Anything, mock.Anything).
// 						RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRServiceRequestRelayPayload, error) {
// 							return nil, fmt.Errorf("an error occurred")
// 						})

// 					return args{ctx: ctx, serviceRequestID: uuid.NewString()}
// 				},
// 				wantErr: true,
// 			},
// 			{
// 				name: "Sad Case - Fail to get patient",
// 				setup: func(mh *usecaseMock.Mocks) args {
// 					mh.FHIR.EXPECT().GetFHIRServiceRequest(mock.Anything, mock.Anything).
// 						RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRServiceRequestRelayPayload, error) {
// 							return servicerequest, nil
// 						})
// 					mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
// 						RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
// 							return nil, fmt.Errorf("an error occurred")
// 						})

// 					return args{ctx: ctx, serviceRequestID: uuid.NewString()}
// 				},
// 				wantErr: true,
// 			},
// 			{
// 				name: "Sad Case - unable to upload media",
// 				setup: func(mh *usecaseMock.Mocks) args {
// 					mh.FHIR.EXPECT().GetFHIRServiceRequest(mock.Anything, mock.Anything).
// 						RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRServiceRequestRelayPayload, error) {
// 							return servicerequest, nil
// 						})
// 					mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
// 						RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
// 							return patient, nil
// 						})
// 					mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
// 						RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
// 							return organization, nil
// 						})
// 					mh.FHIR.EXPECT().SearchPatientObservations(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
// 						RunAndReturn(func(ctx context.Context, bundleID string, searchParameters map[string]interface{}, tenant dto.TenantIdentifiers, pagination serverutils.PaginationInput) (*domain.PagedFHIRObservations, error) {
// 							return pagedFHIRObservations, nil
// 						})
// 					mh.Upload.EXPECT().UploadMedia(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
// 						RunAndReturn(func(ctx context.Context, name string, file io.Reader, contentType string) (*dto.Media, error) {
// 							return nil, fmt.Errorf("an error occurred")
// 						})

// 					return args{ctx: ctx, serviceRequestID: uuid.NewString()}
// 				},
// 				wantErr: true,
// 			},
// 			{
// 				name: "Sad Case - unable to create FHIR document reference",
// 				setup: func(mh *usecaseMock.Mocks) args {
// 					mh.FHIR.EXPECT().GetFHIRServiceRequest(mock.Anything, mock.Anything).
// 						RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRServiceRequestRelayPayload, error) {
// 							return servicerequest, nil
// 						})
// 					mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
// 						RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
// 							return patient, nil
// 						})
// 					mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
// 						RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
// 							return organization, nil
// 						})
// 					mh.FHIR.EXPECT().SearchPatientObservations(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
// 						RunAndReturn(func(ctx context.Context, bundleID string, searchParameters map[string]interface{}, tenant dto.TenantIdentifiers, pagination serverutils.PaginationInput) (*domain.PagedFHIRObservations, error) {
// 							return pagedFHIRObservations, nil
// 						})
// 					mh.URLShortener.EXPECT().Shorten(mock.Anything, mock.Anything).
// 						RunAndReturn(func(ctx context.Context, payload *silurlshortener.ShortenURLPayload) (*silurlshortener.ShortenURLResponse, error) {
// 							return urlresponse, nil
// 						})
// 					mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
// 						RunAndReturn(func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
// 							return conceptpayload, nil
// 						})
// 					mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
// 						RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
// 							return tenantIDs, nil
// 						})
// 					mh.FHIR.EXPECT().CreateFHIRDocumentReference(mock.Anything, mock.Anything).
// 						RunAndReturn(func(ctx context.Context, input *domain.FHIRDocumentReferenceInput) (*domain.FHIRDocumentReference, error) {
// 							return nil, fmt.Errorf("an error occurred")
// 						})
// 					mh.Upload.EXPECT().UploadMedia(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
// 						RunAndReturn(func(ctx context.Context, name string, file io.Reader, contentType string) (*dto.Media, error) {
// 							return media, nil
// 						})

// 					return args{ctx: ctx, serviceRequestID: uuid.NewString()}
// 				},
// 				wantErr: true,
// 			},
// 			{
// 				name: "Sad Case - unable to get terminology concept",
// 				setup: func(mh *usecaseMock.Mocks) args {
// 					mh.FHIR.EXPECT().GetFHIRServiceRequest(mock.Anything, mock.Anything).
// 						RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRServiceRequestRelayPayload, error) {
// 							return servicerequest, nil
// 						})
// 					mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
// 						RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
// 							return patient, nil
// 						})
// 					mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
// 						RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
// 							return organization, nil
// 						})
// 					mh.FHIR.EXPECT().SearchPatientObservations(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
// 						RunAndReturn(func(ctx context.Context, bundleID string, searchParameters map[string]interface{}, tenant dto.TenantIdentifiers, pagination serverutils.PaginationInput) (*domain.PagedFHIRObservations, error) {
// 							return pagedFHIRObservations, nil
// 						})
// 					mh.URLShortener.EXPECT().Shorten(mock.Anything, mock.Anything).
// 						RunAndReturn(func(ctx context.Context, payload *silurlshortener.ShortenURLPayload) (*silurlshortener.ShortenURLResponse, error) {
// 							return urlresponse, nil
// 						})
// 					mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
// 						RunAndReturn(func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
// 							return nil, fmt.Errorf("an error  occurred")
// 						})
// 					mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
// 						RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
// 							return tenantIDs, nil
// 						})
// 					mh.FHIR.EXPECT().CreateFHIRDocumentReference(mock.Anything, mock.Anything).
// 						RunAndReturn(func(ctx context.Context, input *domain.FHIRDocumentReferenceInput) (*domain.FHIRDocumentReference, error) {
// 							return documentRef, nil
// 						})
// 					mh.Upload.EXPECT().UploadMedia(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
// 						RunAndReturn(func(ctx context.Context, name string, file io.Reader, contentType string) (*dto.Media, error) {
// 							return media, nil
// 						})

// 					return args{ctx: ctx, serviceRequestID: uuid.NewString()}
// 				},
// 				wantErr: true,
// 			},
// 			{
// 				name: "Sad Case - Fail to get organization",
// 				setup: func(mh *usecaseMock.Mocks) args {
// 					mh.FHIR.EXPECT().GetFHIRServiceRequest(mock.Anything, mock.Anything).
// 						RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRServiceRequestRelayPayload, error) {
// 							return servicerequest, nil
// 						})
// 					mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
// 						RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
// 							return patient, nil
// 						})
// 					mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
// 						RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
// 							return nil, fmt.Errorf("an error occurred")
// 						})
// 					return args{ctx: ctx, serviceRequestID: uuid.NewString()}
// 				},
// 				wantErr: true,
// 			},
// 			{
// 				name: "Sad Case - unable to get tenant meta tags",
// 				setup: func(mh *usecaseMock.Mocks) args {
// 					mh.FHIR.EXPECT().GetFHIRServiceRequest(mock.Anything, mock.Anything).
// 						RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRServiceRequestRelayPayload, error) {
// 							return servicerequest, nil
// 						})
// 					mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
// 						RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
// 							return patient, nil
// 						})
// 					mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
// 						RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
// 							return nil, fmt.Errorf("an error occurred")
// 						})
// 					return args{ctx: ctx, serviceRequestID: uuid.NewString()}
// 				},
// 				wantErr: true,
// 			},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
// 			args := tt.setup(&mock)

// 			if _, err := clinicalUsecase.GenerateReferralReportPDF(args.ctx, args.serviceRequestID); (err != nil) != tt.wantErr {
// 				t.Errorf("UseCasesClinicalImpl.GenerateReferralReportPDF() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }

func TestUseCasesClinicalImpl_CreateDocumentReference(t *testing.T) {
	ctx := usecaseMock.AddTenantIdentifierContext(context.Background())
	ID := gofakeit.UUID()
	ref := fmt.Sprintf("ServiceRequest/%s", ID)
	encounterRef := fmt.Sprintf("Encounter/%s", ID)
	url := gofakeit.URL()
	mimeType := "application/json"
	title := fmt.Sprintf("%s's Document Reference", gofakeit.Name())
	subjectRef := fmt.Sprintf("Subject/%s", gofakeit.UUID())

	payload := &clinical.DocumentReferencePayload{
		Subject: &domain.FHIRReferenceInput{
			Reference: &subjectRef,
		},
		Attachment: &domain.FHIRAttachment{
			ContentType: (*scalarutils.Code)(&mimeType),
			URL:         (*scalarutils.URL)(&url),
			Title:       &title,
		},
		Encounter: &domain.FHIRReference{
			ID:        &ID,
			Reference: &encounterRef,
			Display:   ID,
		},
		BasedOn: []*domain.FHIRReferenceInput{
			{
				Reference: &ref,
				Display:   ID,
			},
		},
		TerminologySystem: common.ReferralLOINCTerminologySystem,
	}

	implicitRules := mock.Anything
	language := mock.Anything
	documentRef := &domain.FHIRDocumentReference{
		ID:            uuid.NewString(),
		Status:        domain.DocumentReferenceStatusEnumCurrent,
		ImplicitRules: &implicitRules,
		Language:      &language,
		Meta: &domain.FHIRMeta{
			VersionID:   mock.Anything,
			Source:      mock.Anything,
			LastUpdated: gofakeit.Date(),
		},
		Text:      &domain.FHIRNarrative{},
		Extension: []domain.FHIRExtension{},
	}

	conceptpayload := &domain.Concept{
		ConceptClass: mock.Anything,
		DataType:     mock.Anything,
		ID:           gofakeit.UUID(),
	}

	tenantIDs := &dto.TenantIdentifiers{
		OrganizationID: uuid.NewString(),
		FacilityID:     uuid.NewString(),
	}

	orgName := mock.Anything
	organization := &domain.FHIROrganizationRelayPayload{
		Resource: &domain.FHIROrganization{
			ID:   &ID,
			Name: &orgName,
		},
	}

	type args struct {
		ctx     context.Context
		payload *clinical.DocumentReferencePayload
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy case: create document reference",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
						return conceptpayload, nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return organization, nil
					})
				mh.FHIR.EXPECT().CreateFHIRDocumentReference(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input *domain.FHIRDocumentReferenceInput) (*domain.FHIRDocumentReference, error) {
						return documentRef, nil
					})

				return args{ctx: ctx, payload: payload}
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to create document reference",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
						return conceptpayload, nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return organization, nil
					})
				mh.FHIR.EXPECT().CreateFHIRDocumentReference(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input *domain.FHIRDocumentReferenceInput) (*domain.FHIRDocumentReference, error) {
						return nil, fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, payload: payload}
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get concept",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, payload: payload}
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get tenant meta tags",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
						return conceptpayload, nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, payload: payload}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			if err := clinicalUsecase.CreateDocumentReference(args.ctx, args.payload); (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.CreateDocumentReference() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUseCasesClinicalImpl_ShareReferralForm(t *testing.T) {
	ctx := usecaseMock.AddTenantIdentifierContext(context.Background())
	formInput := &dto.ShareReferralFormInput{
		ServiceRequestID: gofakeit.UUID(),
		WorkstationID:    gofakeit.UUID(),
		BranchID:         gofakeit.UUID(),
	}

	resourceID := gofakeit.UUID()
	contentType := "application/pdf"
	title := "Test Title"
	url := gofakeit.URL()
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

	serviceRequestPayload := &domain.FHIRServiceRequestRelayPayload{
		Resource: &domain.FHIRServiceRequest{
			ID: new(string),
			Extension: []*domain.FHIRExtension{
				{
					URL: "http://savannahghi.org/fhir/StructureDefinition/referred-facility",
					Extension: []domain.Extension{
						{
							URL:         "facilityName",
							ValueString: "Nairobi Hospital",
						},
						{
							URL:         "facilityEmail",
							ValueString: gofakeit.Email(),
						},
					},
				},
			},
		},
	}

	tenantIDs := &dto.TenantIdentifiers{
		OrganizationID: uuid.NewString(),
		FacilityID:     uuid.NewString(),
	}

	ID := uuid.NewString()
	phoneSystem := domain.ContactPointSystemEnumPhone
	emailSystem := domain.ContactPointSystemEnumEmail
	phone := gofakeit.Phone()
	email := gofakeit.Email()
	name := gofakeit.Name()
	patient := &domain.FHIRPatientRelayPayload{
		Resource: &domain.FHIRPatient{
			ID: &ID,
			Telecom: []*domain.FHIRContactPoint{
				{
					System: &phoneSystem,
					Value:  &phone,
				},
				{
					System: &emailSystem,
					Value:  &email,
				},
			},
			Name: []*domain.FHIRHumanName{
				{
					Given: []*string{
						&name,
					},
				},
			},
		},
	}

	urlresponse := &silurlshortener.ShortenURLResponse{
		ShortURL: gofakeit.URL(),
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	type args struct {
		ctx   context.Context
		input *dto.ShareReferralFormInput
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy case: share referral form",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchFHIRDocumentReference(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, searchParams map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRDocumentReference, error) {
						return pagedFHIRDocumentRef, nil
					})
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return patient, nil
					})
				mh.FHIR.EXPECT().GetFHIRServiceRequest(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRServiceRequestRelayPayload, error) {
						return serviceRequestPayload, nil
					})
				mh.URLShortener.EXPECT().Shorten(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, payload *silurlshortener.ShortenURLPayload) (*silurlshortener.ShortenURLResponse, error) {
						return urlresponse, nil
					})
				mh.Advantage.EXPECT().SendSMS(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, workstationID, branchID string, payload dto.SMSPayload) error {
						return nil
					})
				mh.MailSender.EXPECT().SendRawEmail(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params *ses.SendRawEmailInput, optFns ...func(*ses.Options)) (*ses.SendRawEmailOutput, error) {
						return nil, nil
					})

				return args{ctx: ctx, input: formInput}
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to share referral form",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchFHIRDocumentReference(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, searchParams map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRDocumentReference, error) {
						return pagedFHIRDocumentRef, nil
					})
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return patient, nil
					})
				mh.FHIR.EXPECT().GetFHIRServiceRequest(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRServiceRequestRelayPayload, error) {
						return serviceRequestPayload, nil
					})
				mh.URLShortener.EXPECT().Shorten(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, payload *silurlshortener.ShortenURLPayload) (*silurlshortener.ShortenURLResponse, error) {
						return urlresponse, nil
					})
				mh.Advantage.EXPECT().SendSMS(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, workstationID, branchID string, payload dto.SMSPayload) error {
						return nil
					})
				mh.MailSender.EXPECT().SendRawEmail(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params *ses.SendRawEmailInput, optFns ...func(*ses.Options)) (*ses.SendRawEmailOutput, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, input: formInput}
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get tenant identifiers",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, input: formInput}
			},
			wantErr: true,
		},
		{
			name: "Sad case: no document references found",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchFHIRDocumentReference(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, searchParams map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRDocumentReference, error) {
						return &domain.PagedFHIRDocumentReference{
							DocumentReferences: []domain.FHIRDocumentReference{},
						}, nil
					})

				return args{ctx: ctx, input: formInput}
			},
			wantErr: true,
		},
		{
			name: "Sad case: fail to search document reference",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchFHIRDocumentReference(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, searchParams map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRDocumentReference, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, input: formInput}
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get patient",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchFHIRDocumentReference(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, searchParams map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRDocumentReference, error) {
						return pagedFHIRDocumentRef, nil
					})
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, input: formInput}
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to send SMS",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchFHIRDocumentReference(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, searchParams map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRDocumentReference, error) {
						return pagedFHIRDocumentRef, nil
					})
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return patient, nil
					})
				mh.FHIR.EXPECT().GetFHIRServiceRequest(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRServiceRequestRelayPayload, error) {
						return serviceRequestPayload, nil
					})
				mh.URLShortener.EXPECT().Shorten(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, payload *silurlshortener.ShortenURLPayload) (*silurlshortener.ShortenURLResponse, error) {
						return urlresponse, nil
					})
				mh.Advantage.EXPECT().SendSMS(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, workstationID, branchID string, payload dto.SMSPayload) error {
						return fmt.Errorf("an error occurred")
					})
				mh.MailSender.EXPECT().SendRawEmail(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params *ses.SendRawEmailInput, optFns ...func(*ses.Options)) (*ses.SendRawEmailOutput, error) {
						return nil, nil
					})

				return args{ctx: ctx, input: formInput}
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get subject associated with the document reference",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchFHIRDocumentReference(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, searchParams map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRDocumentReference, error) {
						return &domain.PagedFHIRDocumentReference{
							DocumentReferences: []domain.FHIRDocumentReference{
								{
									ID:       resourceID,
									Meta:     &domain.FHIRMeta{},
									Type:     &domain.FHIRCodeableConcept{},
									Category: []domain.FHIRCodeableConcept{},
								},
							},
							HasNextPage:     false,
							NextCursor:      "",
							HasPreviousPage: false,
							PreviousCursor:  "",
							TotalCount:      0,
						}, nil
					})

				return args{ctx: ctx, input: formInput}
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to shorten URL",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchFHIRDocumentReference(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, searchParams map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRDocumentReference, error) {
						return pagedFHIRDocumentRef, nil
					})
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return patient, nil
					})
				mh.FHIR.EXPECT().GetFHIRServiceRequest(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRServiceRequestRelayPayload, error) {
						return serviceRequestPayload, nil
					})
				mh.URLShortener.EXPECT().Shorten(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, payload *silurlshortener.ShortenURLPayload) (*silurlshortener.ShortenURLResponse, error) {
						return nil, fmt.Errorf("an error occurred")
					})
				mh.MailSender.EXPECT().SendRawEmail(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params *ses.SendRawEmailInput, optFns ...func(*ses.Options)) (*ses.SendRawEmailOutput, error) {
						return nil, nil
					})

				return args{ctx: ctx, input: formInput}
			},
			wantErr: true,
		},
		{
			name: "Sad case: fail to get service request",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchFHIRDocumentReference(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, searchParams map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRDocumentReference, error) {
						return pagedFHIRDocumentRef, nil
					})
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return patient, nil
					})
				mh.FHIR.EXPECT().GetFHIRServiceRequest(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRServiceRequestRelayPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, input: formInput}
			},
			wantErr: true,
		},
		{
			name: "Sad case: fail to send email",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchFHIRDocumentReference(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, searchParams map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRDocumentReference, error) {
						return pagedFHIRDocumentRef, nil
					})
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return patient, nil
					})
				mh.FHIR.EXPECT().GetFHIRServiceRequest(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRServiceRequestRelayPayload, error) {
						return serviceRequestPayload, nil
					})
				mh.URLShortener.EXPECT().Shorten(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, payload *silurlshortener.ShortenURLPayload) (*silurlshortener.ShortenURLResponse, error) {
						return urlresponse, nil
					})
				mh.Advantage.EXPECT().SendSMS(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, workstationID, branchID string, payload dto.SMSPayload) error {
						return nil
					})
				mh.MailSender.EXPECT().SendRawEmail(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params *ses.SendRawEmailInput, optFns ...func(*ses.Options)) (*ses.SendRawEmailOutput, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, input: formInput}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			httpmock.RegisterResponder(http.MethodGet, url, func(r *http.Request) (*http.Response, error) {
				return httpmock.NewJsonResponse(http.StatusOK, []byte("Referral form"))
			})

			_, err := clinicalUsecase.ShareReferralForm(args.ctx, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.ShareReferralForm() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
