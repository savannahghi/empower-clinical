package clinical_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/scalarutils"
	"github.com/stretchr/testify/mock"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/empower-clinical/pkg/clinical/domain"
	usecaseMock "github.com/savannahghi/empower-clinical/pkg/clinical/usecases/mock"
)

func TestUseCasesClinicalImpl_RecordMammographyResult(t *testing.T) {
	url := gofakeit.URL()
	ctx := usecaseMock.AddTenantIdentifierContext(context.Background())

	input := dto.DiagnosticReportInput{
		Date:        &scalarutils.Date{Year: 2025, Month: 1, Day: 1},
		EncounterID: "12345678905432345",
		Note:        "Test",
		Findings:    "BIRADS 0",
		Media: []*dto.Media{
			{
				ID:        gofakeit.UUID(),
				MediaLink: url,
				Name:      gofakeit.BeerName(),
			},
		},
		UsageContext: dto.BreastCancerScreeningTypeEnum,
	}

	ID := uuid.NewString()
	ref := "Reference/12345"
	encounter := &domain.FHIREncounterRelayPayload{
		Resource: &domain.FHIREncounter{
			ID:     &ID,
			Status: domain.EncounterStatusEnumInProgress,
			Subject: &domain.FHIRReference{
				ID:        &ID,
				Reference: &ref,
				Display:   mock.Anything,
			},
		},
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

	status := domain.ObservationStatusEnumRegistered
	valueString := gofakeit.BeerName()
	noteText := scalarutils.Markdown(mock.Anything)
	instant := scalarutils.Instant(time.Now().GoString())
	system := scalarutils.URI(gofakeit.BeerName())
	code := scalarutils.Code("code")
	observation := &domain.FHIRObservation{
		ID: &ID,
		Text: &domain.FHIRNarrative{
			Div: scalarutils.XHTML(gofakeit.BeerName()),
		},
		ValueString: &valueString,
		Status:      &status,
		Code: &domain.FHIRCodeableConcept{
			Coding: []*domain.FHIRCoding{
				{
					System:  &system,
					Display: gofakeit.BeerName(),
					Code:    &code,
				},
			},
		},
		Subject: &domain.FHIRReference{
			ID:      &ID,
			Display: gofakeit.BeerName(),
		},
		Encounter: &domain.FHIRReference{
			Display: gofakeit.BeerName(),
		},
		Note: []*domain.FHIRAnnotation{
			{
				ID:   &ID,
				Text: &noteText,
			},
		},
		Category: []*domain.FHIRCodeableConcept{
			{
				ID: &ID,
				Coding: []*domain.FHIRCoding{
					{
						ID: &ID,
					},
				},
				Text: mock.Anything,
			},
		},
		EffectiveInstant: &instant,
	}

	type args struct {
		ctx   context.Context
		input dto.DiagnosticReportInput
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy case: record mammography report",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
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
				mh.FHIR.EXPECT().CreateFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
						return observation, nil
					})
				mh.FHIR.EXPECT().CreateFHIRDiagnosticReport(mock.Anything, mock.Anything).
					RunAndReturn(func(context1 context.Context, input *domain.FHIRDiagnosticReportInput) (*domain.FHIRDiagnosticReport, error) {
						issued := mock.Anything
						diagnosticsReport := &domain.FHIRDiagnosticReport{
							ID: &ID,
							Encounter: &domain.FHIRReference{
								ID: &ID,
							},
							Subject: &domain.FHIRReference{
								ID: &ID,
							},
							Status: domain.DiagnosticReportStatusAmended,
							Issued: &issued,
						}
						return diagnosticsReport, nil
					})

				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to get meta tags",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
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

				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get concept",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to create FHIR observation",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
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
				mh.FHIR.EXPECT().CreateFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to create FHIR diagnostic report",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
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
				mh.FHIR.EXPECT().CreateFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
						return observation, nil
					})
				mh.FHIR.EXPECT().CreateFHIRDiagnosticReport(mock.Anything, mock.Anything).
					RunAndReturn(func(context1 context.Context, input *domain.FHIRDiagnosticReportInput) (*domain.FHIRDiagnosticReport, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
		{
			name: "Sad case: required field omitted",
			setup: func(mh *usecaseMock.Mocks) args {
				input := dto.DiagnosticReportInput{
					Date: &scalarutils.Date{Year: 2025, Month: 1, Day: 1},
					Note: "Test",
				}
				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			_, err := clinicalUsecase.RecordMammographyResult(args.ctx, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.RecordMammographyResult() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCasesClinicalImpl_RecordBiopsy(t *testing.T) {
	ctx := usecaseMock.AddTenantIdentifierContext(context.Background())

	input := dto.DiagnosticReportInput{
		Date:        &scalarutils.Date{Year: 2025, Month: 1, Day: 1},
		EncounterID: "12345678905432345",
		Note:        "Go for biopsy test",
		Media: []*dto.Media{
			{
				ID:        gofakeit.UUID(),
				MediaLink: gofakeit.URL(),
				Name:      gofakeit.BeerName(),
			},
		},
		Findings: gofakeit.HipsterSentence(20),
	}

	ID := uuid.NewString()
	ref := "Reference/12345"
	encounter := &domain.FHIREncounterRelayPayload{
		Resource: &domain.FHIREncounter{
			ID:     &ID,
			Status: domain.EncounterStatusEnumInProgress,
			Subject: &domain.FHIRReference{
				ID:        &ID,
				Reference: &ref,
				Display:   mock.Anything,
			},
		},
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

	status := domain.ObservationStatusEnumRegistered
	valueString := gofakeit.BeerName()
	noteText := scalarutils.Markdown(mock.Anything)
	instant := scalarutils.Instant(time.Now().GoString())
	system := scalarutils.URI(gofakeit.BeerName())
	code := scalarutils.Code("code")
	observation := &domain.FHIRObservation{
		ID: &ID,
		Text: &domain.FHIRNarrative{
			Div: scalarutils.XHTML(gofakeit.BeerName()),
		},
		ValueString: &valueString,
		Status:      &status,
		Code: &domain.FHIRCodeableConcept{
			Coding: []*domain.FHIRCoding{
				{
					System:  &system,
					Display: gofakeit.BeerName(),
					Code:    &code,
				},
			},
		},
		Subject: &domain.FHIRReference{
			ID:      &ID,
			Display: gofakeit.BeerName(),
		},
		Encounter: &domain.FHIRReference{
			Display: gofakeit.BeerName(),
		},
		Note: []*domain.FHIRAnnotation{
			{
				ID:   &ID,
				Text: &noteText,
			},
		},
		Category: []*domain.FHIRCodeableConcept{
			{
				ID: &ID,
				Coding: []*domain.FHIRCoding{
					{
						ID: &ID,
					},
				},
				Text: mock.Anything,
			},
		},
		EffectiveInstant: &instant,
	}

	type args struct {
		ctx   context.Context
		input dto.DiagnosticReportInput
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy case: successfully record biopsy test",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
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
				mh.FHIR.EXPECT().CreateFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
						return observation, nil
					})
				mh.FHIR.EXPECT().CreateFHIRDiagnosticReport(mock.Anything, mock.Anything).
					RunAndReturn(func(context1 context.Context, input *domain.FHIRDiagnosticReportInput) (*domain.FHIRDiagnosticReport, error) {
						issued := mock.Anything
						diagnosticsReport := &domain.FHIRDiagnosticReport{
							ID: &ID,
							Encounter: &domain.FHIRReference{
								ID: &ID,
							},
							Subject: &domain.FHIRReference{
								ID: &ID,
							},
							Status: domain.DiagnosticReportStatusAmended,
							Issued: &issued,
						}
						return diagnosticsReport, nil
					})

				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to successfully record biopsy test",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
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
				mh.FHIR.EXPECT().CreateFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
						return observation, nil
					})
				mh.FHIR.EXPECT().CreateFHIRDiagnosticReport(mock.Anything, mock.Anything).
					RunAndReturn(func(context1 context.Context, input *domain.FHIRDiagnosticReportInput) (*domain.FHIRDiagnosticReport, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to record observation",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
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
				mh.FHIR.EXPECT().CreateFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, input: input}
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
				mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
						return conceptpayload, nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get organization",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
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

				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			_, err := clinicalUsecase.RecordBiopsy(args.ctx, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.RecordBiopsy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCasesClinicalImpl_RecordMRI(t *testing.T) {
	ctx := usecaseMock.AddTenantIdentifierContext(context.Background())

	input := dto.DiagnosticReportInput{
		Date:        &scalarutils.Date{Year: 2025, Month: 1, Day: 1},
		EncounterID: "12345678905432345",
		Note:        "No Tumours observed",
		Media: []*dto.Media{
			{
				ID:        gofakeit.UUID(),
				MediaLink: gofakeit.URL(),
				Name:      gofakeit.BeerName(),
			},
		},
		Findings: "BIRADS 5",
	}

	ID := uuid.NewString()
	ref := "Reference/12345"
	encounter := &domain.FHIREncounterRelayPayload{
		Resource: &domain.FHIREncounter{
			ID:     &ID,
			Status: domain.EncounterStatusEnumInProgress,
			Subject: &domain.FHIRReference{
				ID:        &ID,
				Reference: &ref,
				Display:   mock.Anything,
			},
		},
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

	status := domain.ObservationStatusEnumRegistered
	valueString := gofakeit.BeerName()
	noteText := scalarutils.Markdown(mock.Anything)
	instant := scalarutils.Instant(time.Now().GoString())
	system := scalarutils.URI(gofakeit.BeerName())
	code := scalarutils.Code("code")
	observation := &domain.FHIRObservation{
		ID: &ID,
		Text: &domain.FHIRNarrative{
			Div: scalarutils.XHTML(gofakeit.BeerName()),
		},
		ValueString: &valueString,
		Status:      &status,
		Code: &domain.FHIRCodeableConcept{
			Coding: []*domain.FHIRCoding{
				{
					System:  &system,
					Display: gofakeit.BeerName(),
					Code:    &code,
				},
			},
		},
		Subject: &domain.FHIRReference{
			ID:      &ID,
			Display: gofakeit.BeerName(),
		},
		Encounter: &domain.FHIRReference{
			Display: gofakeit.BeerName(),
		},
		Note: []*domain.FHIRAnnotation{
			{
				ID:   &ID,
				Text: &noteText,
			},
		},
		Category: []*domain.FHIRCodeableConcept{
			{
				ID: &ID,
				Coding: []*domain.FHIRCoding{
					{
						ID: &ID,
					},
				},
				Text: mock.Anything,
			},
		},
		EffectiveInstant: &instant,
	}

	type args struct {
		ctx   context.Context
		input dto.DiagnosticReportInput
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy case: successfully record mri results",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
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
				mh.FHIR.EXPECT().CreateFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
						return observation, nil
					})
				mh.FHIR.EXPECT().CreateFHIRDiagnosticReport(mock.Anything, mock.Anything).
					RunAndReturn(func(context1 context.Context, input *domain.FHIRDiagnosticReportInput) (*domain.FHIRDiagnosticReport, error) {
						issued := mock.Anything
						diagnosticsReport := &domain.FHIRDiagnosticReport{
							ID: &ID,
							Encounter: &domain.FHIRReference{
								ID: &ID,
							},
							Subject: &domain.FHIRReference{
								ID: &ID,
							},
							Status: domain.DiagnosticReportStatusAmended,
							Issued: &issued,
						}
						return diagnosticsReport, nil
					})

				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to successfully record mri results",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
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
				mh.FHIR.EXPECT().CreateFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
						return observation, nil
					})
				mh.FHIR.EXPECT().CreateFHIRDiagnosticReport(mock.Anything, mock.Anything).
					RunAndReturn(func(context1 context.Context, input *domain.FHIRDiagnosticReportInput) (*domain.FHIRDiagnosticReport, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
		{
			name: "Sad case: fail validation",
			setup: func(mh *usecaseMock.Mocks) args {
				input := dto.DiagnosticReportInput{
					Date: &scalarutils.Date{Year: 2025, Month: 1, Day: 1},
					Media: []*dto.Media{
						{
							ID:        gofakeit.UUID(),
							MediaLink: gofakeit.URL(),
							Name:      gofakeit.BeerName(),
						},
					},
				}

				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to record observation",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
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
				mh.FHIR.EXPECT().CreateFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, input: input}
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
				mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
						return conceptpayload, nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get organization",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
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
				mh.FHIR.EXPECT().CreateFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)
			_, err := clinicalUsecase.RecordMRI(args.ctx, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.RecordMRI() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCasesClinicalImpl_RecordUltrasound(t *testing.T) {
	ctx := usecaseMock.AddTenantIdentifierContext(context.Background())

	input := dto.DiagnosticReportInput{
		Date:        &scalarutils.Date{Year: 2025, Month: 1, Day: 1},
		EncounterID: "12345678905432345",
		Note:        "No lumps felt",
		Media: []*dto.Media{
			{
				ID:        gofakeit.UUID(),
				MediaLink: gofakeit.URL(),
				Name:      gofakeit.BeerName(),
			},
		},
		Findings: "BIRADS 3",
	}

	ID := uuid.NewString()
	ref := "Reference/12345"
	encounter := &domain.FHIREncounterRelayPayload{
		Resource: &domain.FHIREncounter{
			ID:     &ID,
			Status: domain.EncounterStatusEnumInProgress,
			Subject: &domain.FHIRReference{
				ID:        &ID,
				Reference: &ref,
				Display:   mock.Anything,
			},
		},
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

	status := domain.ObservationStatusEnumRegistered
	valueString := gofakeit.BeerName()
	noteText := scalarutils.Markdown(mock.Anything)
	instant := scalarutils.Instant(time.Now().GoString())
	system := scalarutils.URI(gofakeit.BeerName())
	code := scalarutils.Code("code")
	observation := &domain.FHIRObservation{
		ID: &ID,
		Text: &domain.FHIRNarrative{
			Div: scalarutils.XHTML(gofakeit.BeerName()),
		},
		ValueString: &valueString,
		Status:      &status,
		Code: &domain.FHIRCodeableConcept{
			Coding: []*domain.FHIRCoding{
				{
					System:  &system,
					Display: gofakeit.BeerName(),
					Code:    &code,
				},
			},
		},
		Subject: &domain.FHIRReference{
			ID:      &ID,
			Display: gofakeit.BeerName(),
		},
		Encounter: &domain.FHIRReference{
			Display: gofakeit.BeerName(),
		},
		Note: []*domain.FHIRAnnotation{
			{
				ID:   &ID,
				Text: &noteText,
			},
		},
		Category: []*domain.FHIRCodeableConcept{
			{
				ID: &ID,
				Coding: []*domain.FHIRCoding{
					{
						ID: &ID,
					},
				},
				Text: mock.Anything,
			},
		},
		EffectiveInstant: &instant,
	}

	type args struct {
		ctx   context.Context
		input dto.DiagnosticReportInput
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy case: record ultrasound",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
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
				mh.FHIR.EXPECT().CreateFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
						return observation, nil
					})
				mh.FHIR.EXPECT().CreateFHIRDiagnosticReport(mock.Anything, mock.Anything).
					RunAndReturn(func(context1 context.Context, input *domain.FHIRDiagnosticReportInput) (*domain.FHIRDiagnosticReport, error) {
						issued := mock.Anything
						diagnosticsReport := &domain.FHIRDiagnosticReport{
							ID: &ID,
							Encounter: &domain.FHIRReference{
								ID: &ID,
							},
							Subject: &domain.FHIRReference{
								ID: &ID,
							},
							Status: domain.DiagnosticReportStatusAmended,
							Issued: &issued,
						}
						return diagnosticsReport, nil
					})

				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to record ultrasound results",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
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
				mh.FHIR.EXPECT().CreateFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
						return observation, nil
					})
				mh.FHIR.EXPECT().CreateFHIRDiagnosticReport(mock.Anything, mock.Anything).
					RunAndReturn(func(context1 context.Context, input *domain.FHIRDiagnosticReportInput) (*domain.FHIRDiagnosticReport, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
		{
			name: "Sad case: fail validation",
			setup: func(mh *usecaseMock.Mocks) args {
				input := dto.DiagnosticReportInput{
					Date: &scalarutils.Date{Year: 2025, Month: 1, Day: 1},
					Media: []*dto.Media{
						{
							ID:        gofakeit.UUID(),
							MediaLink: gofakeit.URL(),
							Name:      gofakeit.BeerName(),
						},
					},
				}

				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to record observation",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
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
				mh.FHIR.EXPECT().CreateFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, input: input}
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
				mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
						return conceptpayload, nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get organization",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
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

				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			_, err := clinicalUsecase.RecordUltrasound(args.ctx, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.RecordUltrasound() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCasesClinicalImpl_RecordCBE(t *testing.T) {
	ctx := usecaseMock.AddTenantIdentifierContext(context.Background())

	input := dto.DiagnosticReportInput{
		Date:        &scalarutils.Date{Year: 2025, Month: 1, Day: 1},
		EncounterID: "12345678905432345",
		Note:        "No lumps felt",
		Media: []*dto.Media{
			{
				ID:        gofakeit.UUID(),
				MediaLink: gofakeit.URL(),
				Name:      gofakeit.BeerName(),
			},
		},
		Findings: "BIRADS 1",
	}

	ID := uuid.NewString()
	ref := "Reference/12345"
	encounter := &domain.FHIREncounterRelayPayload{
		Resource: &domain.FHIREncounter{
			ID:     &ID,
			Status: domain.EncounterStatusEnumInProgress,
			Subject: &domain.FHIRReference{
				ID:        &ID,
				Reference: &ref,
				Display:   mock.Anything,
			},
		},
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

	status := domain.ObservationStatusEnumRegistered
	valueString := gofakeit.BeerName()
	noteText := scalarutils.Markdown(mock.Anything)
	instant := scalarutils.Instant(time.Now().GoString())
	system := scalarutils.URI(gofakeit.BeerName())
	code := scalarutils.Code("code")
	observation := &domain.FHIRObservation{
		ID: &ID,
		Text: &domain.FHIRNarrative{
			Div: scalarutils.XHTML(gofakeit.BeerName()),
		},
		ValueString: &valueString,
		Status:      &status,
		Code: &domain.FHIRCodeableConcept{
			Coding: []*domain.FHIRCoding{
				{
					System:  &system,
					Display: gofakeit.BeerName(),
					Code:    &code,
				},
			},
		},
		Subject: &domain.FHIRReference{
			ID:      &ID,
			Display: gofakeit.BeerName(),
		},
		Encounter: &domain.FHIRReference{
			Display: gofakeit.BeerName(),
		},
		Note: []*domain.FHIRAnnotation{
			{
				ID:   &ID,
				Text: &noteText,
			},
		},
		Category: []*domain.FHIRCodeableConcept{
			{
				ID: &ID,
				Coding: []*domain.FHIRCoding{
					{
						ID: &ID,
					},
				},
				Text: mock.Anything,
			},
		},
		EffectiveInstant: &instant,
	}

	type args struct {
		ctx   context.Context
		input dto.DiagnosticReportInput
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy case: record CBE test",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
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
				mh.FHIR.EXPECT().CreateFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
						return observation, nil
					})
				mh.FHIR.EXPECT().CreateFHIRDiagnosticReport(mock.Anything, mock.Anything).
					RunAndReturn(func(context1 context.Context, input *domain.FHIRDiagnosticReportInput) (*domain.FHIRDiagnosticReport, error) {
						issued := mock.Anything
						diagnosticsReport := &domain.FHIRDiagnosticReport{
							ID: &ID,
							Encounter: &domain.FHIRReference{
								ID: &ID,
							},
							Subject: &domain.FHIRReference{
								ID: &ID,
							},
							Status: domain.DiagnosticReportStatusAmended,
							Issued: &issued,
						}
						return diagnosticsReport, nil
					})

				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to record CBE test",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
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
				mh.FHIR.EXPECT().CreateFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
						return observation, nil
					})
				mh.FHIR.EXPECT().CreateFHIRDiagnosticReport(mock.Anything, mock.Anything).
					RunAndReturn(func(context1 context.Context, input *domain.FHIRDiagnosticReportInput) (*domain.FHIRDiagnosticReport, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
		{
			name: "Sad case: fail validation",
			setup: func(mh *usecaseMock.Mocks) args {
				input := dto.DiagnosticReportInput{
					Date: &scalarutils.Date{Year: 2025, Month: 1, Day: 1},
					Media: []*dto.Media{
						{
							ID:        gofakeit.UUID(),
							MediaLink: gofakeit.URL(),
							Name:      gofakeit.BeerName(),
						},
					},
				}

				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to record observation",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
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
				mh.FHIR.EXPECT().CreateFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, input: input}
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
				mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
						return conceptpayload, nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get organization",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
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

				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			_, err := clinicalUsecase.RecordCBE(args.ctx, &args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.RecordCBE() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCasesClinicalImpl_RecordPapSmear(t *testing.T) {
	ctx := usecaseMock.AddTenantIdentifierContext(context.Background())

	input := dto.DiagnosticReportInput{
		Date:        &scalarutils.Date{Year: 2025, Month: 1, Day: 1},
		EncounterID: gofakeit.UUID(),
		Note:        "Additional information",
		Media: []*dto.Media{
			{
				ID:        gofakeit.UUID(),
				MediaLink: gofakeit.URL(),
				Name:      gofakeit.BeerName(),
			},
		},
		Findings: "ASCUS or greater",
	}

	ID := uuid.NewString()
	ref := "Reference/12345"
	encounter := &domain.FHIREncounterRelayPayload{
		Resource: &domain.FHIREncounter{
			ID:     &ID,
			Status: domain.EncounterStatusEnumInProgress,
			Subject: &domain.FHIRReference{
				ID:        &ID,
				Reference: &ref,
				Display:   mock.Anything,
			},
		},
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

	status := domain.ObservationStatusEnumRegistered
	valueString := gofakeit.BeerName()
	noteText := scalarutils.Markdown(mock.Anything)
	instant := scalarutils.Instant(time.Now().GoString())
	system := scalarutils.URI(gofakeit.BeerName())
	code := scalarutils.Code("code")
	observation := &domain.FHIRObservation{
		ID: &ID,
		Text: &domain.FHIRNarrative{
			Div: scalarutils.XHTML(gofakeit.BeerName()),
		},
		ValueString: &valueString,
		Status:      &status,
		Code: &domain.FHIRCodeableConcept{
			Coding: []*domain.FHIRCoding{
				{
					System:  &system,
					Display: gofakeit.BeerName(),
					Code:    &code,
				},
			},
		},
		Subject: &domain.FHIRReference{
			ID:      &ID,
			Display: gofakeit.BeerName(),
		},
		Encounter: &domain.FHIRReference{
			Display: gofakeit.BeerName(),
		},
		Note: []*domain.FHIRAnnotation{
			{
				ID:   &ID,
				Text: &noteText,
			},
		},
		Category: []*domain.FHIRCodeableConcept{
			{
				ID: &ID,
				Coding: []*domain.FHIRCoding{
					{
						ID: &ID,
					},
				},
				Text: mock.Anything,
			},
		},
		EffectiveInstant: &instant,
	}

	type args struct {
		ctx   context.Context
		input dto.DiagnosticReportInput
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy case: record pap smear",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
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
				mh.FHIR.EXPECT().CreateFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
						return observation, nil
					})
				mh.FHIR.EXPECT().CreateFHIRDiagnosticReport(mock.Anything, mock.Anything).
					RunAndReturn(func(context1 context.Context, input *domain.FHIRDiagnosticReportInput) (*domain.FHIRDiagnosticReport, error) {
						issued := mock.Anything
						diagnosticsReport := &domain.FHIRDiagnosticReport{
							ID: &ID,
							Encounter: &domain.FHIRReference{
								ID: &ID,
							},
							Subject: &domain.FHIRReference{
								ID: &ID,
							},
							Status: domain.DiagnosticReportStatusAmended,
							Issued: &issued,
						}
						return diagnosticsReport, nil
					})

				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to record pap smear",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
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
				mh.FHIR.EXPECT().CreateFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
						return observation, nil
					})
				mh.FHIR.EXPECT().CreateFHIRDiagnosticReport(mock.Anything, mock.Anything).
					RunAndReturn(func(context1 context.Context, input *domain.FHIRDiagnosticReportInput) (*domain.FHIRDiagnosticReport, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
		{
			name: "Sad case: invalid input(no encounter id)",
			setup: func(mh *usecaseMock.Mocks) args {
				input := dto.DiagnosticReportInput{
					Date: &scalarutils.Date{Year: 2025, Month: 1, Day: 1},
					Note: "Additional information",
					Media: []*dto.Media{
						{
							ID:        gofakeit.UUID(),
							MediaLink: gofakeit.URL(),
							Name:      gofakeit.BeerName(),
						},
					},
					Findings: "ASCUS or greater",
				}

				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get tenant meta tags",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
						return conceptpayload, nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to create observation",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
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
				mh.FHIR.EXPECT().CreateFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			_, err := clinicalUsecase.RecordPapSmear(args.ctx, &args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.RecordPapSmear() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCasesClinicalImpl_CreateLabOrderResult(t *testing.T) {
	ctx := usecaseMock.AddTenantIdentifierContext(context.Background())
	input := &dto.TestOrderResult{
		ServiceRequestID: gofakeit.UUID(),
		Test: []dto.TestOrderObservation{
			{
				Test:    gofakeit.Name(),
				Value:   gofakeit.Name(),
				Finding: gofakeit.Name(),
			},
		},
	}

	ID := uuid.NewString()
	ref := "Reference/12345"
	encounter := &domain.FHIREncounterRelayPayload{
		Resource: &domain.FHIREncounter{
			ID:     &ID,
			Status: domain.EncounterStatusEnumInProgress,
			Subject: &domain.FHIRReference{
				ID:        &ID,
				Reference: &ref,
				Display:   mock.Anything,
			},
		},
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

	status := domain.ObservationStatusEnumRegistered
	valueString := gofakeit.BeerName()
	noteText := scalarutils.Markdown(mock.Anything)
	instant := scalarutils.Instant(time.Now().GoString())
	system := scalarutils.URI(gofakeit.BeerName())
	code := scalarutils.Code("code")
	observation := &domain.FHIRObservation{
		ID: &ID,
		Text: &domain.FHIRNarrative{
			Div: scalarutils.XHTML(gofakeit.BeerName()),
		},
		ValueString: &valueString,
		Status:      &status,
		Code: &domain.FHIRCodeableConcept{
			Coding: []*domain.FHIRCoding{
				{
					System:  &system,
					Display: gofakeit.BeerName(),
					Code:    &code,
				},
			},
		},
		Subject: &domain.FHIRReference{
			ID:      &ID,
			Display: gofakeit.BeerName(),
		},
		Encounter: &domain.FHIRReference{
			Display: gofakeit.BeerName(),
		},
		Note: []*domain.FHIRAnnotation{
			{
				ID:   &ID,
				Text: &noteText,
			},
		},
		Category: []*domain.FHIRCodeableConcept{
			{
				ID: &ID,
				Coding: []*domain.FHIRCoding{
					{
						ID: &ID,
					},
				},
				Text: mock.Anything,
			},
		},
		EffectiveInstant: &instant,
	}

	servicerequest := &domain.FHIRServiceRequestRelayPayload{
		Resource: &domain.FHIRServiceRequest{
			ID:         &ID,
			Text:       &domain.FHIRNarrative{},
			Identifier: []*domain.FHIRIdentifier{},
			Status:     domain.ServiceRequestStatusActive,
			Intent:     domain.ServiceRequestIntentDirective,
			Code: &domain.FHIRCodeableReference{
				Concept: &domain.FHIRCodeableConcept{
					Coding: []*domain.FHIRCoding{},
				},
			},
			Subject: &domain.FHIRReference{
				ID:        &ID,
				Reference: &ref,
				Display:   mock.Anything,
			},
			Encounter: &domain.FHIRReference{
				Display:   mock.Anything,
				Reference: &ref,
			},
		},
	}

	type args struct {
		ctx   context.Context
		input *dto.TestOrderResult
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy Case: Successfull add results",
			setup: func(mh *usecaseMock.Mocks) args {

				mh.FHIR.EXPECT().GetFHIRServiceRequest(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRServiceRequestRelayPayload, error) {
						return servicerequest, nil
					})
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return organization, nil
					})
				mh.FHIR.EXPECT().CreateFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
						return observation, nil
					})
				mh.FHIR.EXPECT().CreateFHIRDiagnosticReport(mock.Anything, mock.Anything).
					RunAndReturn(func(context1 context.Context, input *domain.FHIRDiagnosticReportInput) (*domain.FHIRDiagnosticReport, error) {
						issued := mock.Anything
						diagnosticsReport := &domain.FHIRDiagnosticReport{
							ID: &ID,
							Encounter: &domain.FHIRReference{
								ID: &ID,
							},
							Subject: &domain.FHIRReference{
								ID: &ID,
							},
							Status: domain.DiagnosticReportStatusAmended,
							Issued: &issued,
						}
						return diagnosticsReport, nil
					})

				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to get tenant identifiers",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRServiceRequest(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRServiceRequestRelayPayload, error) {
						return servicerequest, nil
					})
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
		{
			name: "Sad case: fail to add results",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRServiceRequest(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRServiceRequestRelayPayload, error) {
						return servicerequest, nil
					})
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return organization, nil
					})
				mh.FHIR.EXPECT().CreateFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
						return observation, nil
					})
				mh.FHIR.EXPECT().CreateFHIRDiagnosticReport(mock.Anything, mock.Anything).
					RunAndReturn(func(context1 context.Context, input *domain.FHIRDiagnosticReportInput) (*domain.FHIRDiagnosticReport, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
		{
			name: "Sad case: missing service request ID",
			setup: func(mh *usecaseMock.Mocks) args {
				input := &dto.TestOrderResult{
					Test: []dto.TestOrderObservation{
						{
							Test:    gofakeit.Name(),
							Value:   gofakeit.Name(),
							Finding: gofakeit.Name(),
						},
					},
				}
				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			got, err := clinicalUsecase.CreateLabOrderResult(args.ctx, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.CreateLabOrderResult() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got == nil {
					t.Errorf("expected a response but got %v", got)
					return
				}
			}
		})
	}
}

func TestUseCasesClinicalImpl_RecordPSA(t *testing.T) {
	ctx := usecaseMock.AddTenantIdentifierContext(context.Background())

	input := dto.PSAInput{
		DiagnosticInput: &dto.DiagnosticReportInput{
			Date:        &scalarutils.Date{Year: 2025, Month: 1, Day: 1},
			EncounterID: gofakeit.UUID(),
			Note:        "Additional information",
			Media: []*dto.Media{
				{
					ID:        gofakeit.UUID(),
					MediaLink: gofakeit.URL(),
					Name:      gofakeit.BeerName(),
				},
			},
			Findings: "ASCUS or greater",
		},
		PSAType: dto.ProstaticSerumAntigen,
	}

	ID := uuid.NewString()
	ref := "Reference/12345"
	encounter := &domain.FHIREncounterRelayPayload{
		Resource: &domain.FHIREncounter{
			ID:     &ID,
			Status: domain.EncounterStatusEnumInProgress,
			Subject: &domain.FHIRReference{
				ID:        &ID,
				Reference: &ref,
				Display:   mock.Anything,
			},
		},
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

	status := domain.ObservationStatusEnumRegistered
	valueString := gofakeit.BeerName()
	noteText := scalarutils.Markdown(mock.Anything)
	instant := scalarutils.Instant(time.Now().GoString())
	system := scalarutils.URI(gofakeit.BeerName())
	code := scalarutils.Code("code")
	observation := &domain.FHIRObservation{
		ID: &ID,
		Text: &domain.FHIRNarrative{
			Div: scalarutils.XHTML(gofakeit.BeerName()),
		},
		ValueString: &valueString,
		Status:      &status,
		Code: &domain.FHIRCodeableConcept{
			Coding: []*domain.FHIRCoding{
				{
					System:  &system,
					Display: gofakeit.BeerName(),
					Code:    &code,
				},
			},
		},
		Subject: &domain.FHIRReference{
			ID:      &ID,
			Display: gofakeit.BeerName(),
		},
		Encounter: &domain.FHIRReference{
			Display: gofakeit.BeerName(),
		},
		Note: []*domain.FHIRAnnotation{
			{
				ID:   &ID,
				Text: &noteText,
			},
		},
		Category: []*domain.FHIRCodeableConcept{
			{
				ID: &ID,
				Coding: []*domain.FHIRCoding{
					{
						ID: &ID,
					},
				},
				Text: mock.Anything,
			},
		},
		EffectiveInstant: &instant,
	}

	type args struct {
		ctx   context.Context
		input dto.PSAInput
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy case: record PSA",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
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
				mh.FHIR.EXPECT().CreateFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
						return observation, nil
					})
				mh.FHIR.EXPECT().CreateFHIRDiagnosticReport(mock.Anything, mock.Anything).
					RunAndReturn(func(context1 context.Context, input *domain.FHIRDiagnosticReportInput) (*domain.FHIRDiagnosticReport, error) {
						issued := mock.Anything
						diagnosticsReport := &domain.FHIRDiagnosticReport{
							ID: &ID,
							Encounter: &domain.FHIRReference{
								ID: &ID,
							},
							Subject: &domain.FHIRReference{
								ID: &ID,
							},
							Status: domain.DiagnosticReportStatusAmended,
							Issued: &issued,
						}
						return diagnosticsReport, nil
					})

				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to record PSA",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
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
				mh.FHIR.EXPECT().CreateFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
						return observation, nil
					})
				mh.FHIR.EXPECT().CreateFHIRDiagnosticReport(mock.Anything, mock.Anything).
					RunAndReturn(func(context1 context.Context, input *domain.FHIRDiagnosticReportInput) (*domain.FHIRDiagnosticReport, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
		{
			name: "Sad case: invalid input(no encounter id)",
			setup: func(mh *usecaseMock.Mocks) args {
				input := dto.PSAInput{
					DiagnosticInput: &dto.DiagnosticReportInput{
						Date:        &scalarutils.Date{Year: 2025, Month: 1, Day: 1},
						EncounterID: "",
						Note:        "Additional information",
						Media: []*dto.Media{
							{
								ID:        gofakeit.UUID(),
								MediaLink: gofakeit.URL(),
								Name:      gofakeit.BeerName(),
							},
						},
						Findings: "ASCUS or greater",
					},
				}

				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get tenant meta tags",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
						return conceptpayload, nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to create observation",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
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
				mh.FHIR.EXPECT().CreateFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			_, err := clinicalUsecase.RecordPSA(args.ctx, &args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.RecordPSA() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCasesClinicalImpl_RecordTestResult(t *testing.T) {
	ctx := usecaseMock.AddTenantIdentifierContext(context.Background())

	input := dto.TestResultInput{
		Entry: dto.DiagnosticReportInput{
			Date:        &scalarutils.Date{Year: 2025, Month: 1, Day: 1},
			EncounterID: gofakeit.UUID(),
			Note:        "Additional information",
			Media: []*dto.Media{
				{
					ID:        gofakeit.UUID(),
					MediaLink: gofakeit.URL(),
					Name:      gofakeit.BeerName(),
				},
			},
			Findings: "ASCUS or greater",
		},
		ServiceRequestID: gofakeit.UUID(),
	}

	ID := uuid.NewString()
	ref := "Reference/12345"
	encounter := &domain.FHIREncounterRelayPayload{
		Resource: &domain.FHIREncounter{
			ID:     &ID,
			Status: domain.EncounterStatusEnumInProgress,
			Subject: &domain.FHIRReference{
				ID:        &ID,
				Reference: &ref,
				Display:   mock.Anything,
			},
		},
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

	status := domain.ObservationStatusEnumRegistered
	valueString := gofakeit.BeerName()
	noteText := scalarutils.Markdown(mock.Anything)
	instant := scalarutils.Instant(time.Now().GoString())
	system := scalarutils.URI(gofakeit.BeerName())
	code := scalarutils.Code("code")
	observation := &domain.FHIRObservation{
		ID: &ID,
		Text: &domain.FHIRNarrative{
			Div: scalarutils.XHTML(gofakeit.BeerName()),
		},
		ValueString: &valueString,
		Status:      &status,
		Code: &domain.FHIRCodeableConcept{
			Coding: []*domain.FHIRCoding{
				{
					System:  &system,
					Display: gofakeit.BeerName(),
					Code:    &code,
				},
			},
		},
		Subject: &domain.FHIRReference{
			ID:      &ID,
			Display: gofakeit.BeerName(),
		},
		Encounter: &domain.FHIRReference{
			Display: gofakeit.BeerName(),
		},
		Note: []*domain.FHIRAnnotation{
			{
				ID:   &ID,
				Text: &noteText,
			},
		},
		Category: []*domain.FHIRCodeableConcept{
			{
				ID: &ID,
				Coding: []*domain.FHIRCoding{
					{
						ID: &ID,
					},
				},
				Text: mock.Anything,
			},
		},
		EffectiveInstant: &instant,
	}

	type args struct {
		ctx   context.Context
		input dto.TestResultInput
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy case: record test result",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return organization, nil
					})
				mh.FHIR.EXPECT().GetFHIRServiceRequest(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRServiceRequestRelayPayload, error) {
						return &domain.FHIRServiceRequestRelayPayload{
							Resource: &domain.FHIRServiceRequest{
								ID:         &ID,
								Identifier: []*domain.FHIRIdentifier{},
								Status:     domain.ServiceRequestStatusActive,
								Intent:     domain.ServiceRequestIntentDirective,
								Subject: &domain.FHIRReference{
									ID:        &ID,
									Reference: &ref,
									Display:   mock.Anything,
								},
								Encounter: &domain.FHIRReference{
									Display:   mock.Anything,
									Reference: &ref,
								},
								Code: &domain.FHIRCodeableReference{
									Concept: &domain.FHIRCodeableConcept{
										Coding: []*domain.FHIRCoding{},
									},
								},
							},
						}, nil
					})
				mh.FHIR.EXPECT().CreateFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
						return observation, nil
					})
				mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
						return conceptpayload, nil
					})
				mh.FHIR.EXPECT().CreateFHIRDiagnosticReport(mock.Anything, mock.Anything).
					RunAndReturn(func(context1 context.Context, input *domain.FHIRDiagnosticReportInput) (*domain.FHIRDiagnosticReport, error) {
						issued := mock.Anything
						diagnosticsReport := &domain.FHIRDiagnosticReport{
							ID: &ID,
							Encounter: &domain.FHIRReference{
								ID: &ID,
							},
							Subject: &domain.FHIRReference{
								ID: &ID,
							},
							Status: domain.DiagnosticReportStatusAmended,
							Issued: &issued,
						}
						return diagnosticsReport, nil
					})
				mh.PubSub.EXPECT().NotifyCreateFollowUpTask(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, data *domain.FHIRTaskInput) error {
						return nil
					})
				mh.FHIR.EXPECT().PatchFHIRServiceRequest(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string, input domain.FHIRServiceRequestInput) (*domain.FHIRServiceRequestRelayPayload, error) {
						return &domain.FHIRServiceRequestRelayPayload{
							Resource: &domain.FHIRServiceRequest{
								Status: domain.ServiceRequestStatusActive,
							},
						}, nil
					})

				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to record test result - no encounter id",
			setup: func(mh *usecaseMock.Mocks) args {
				input := dto.TestResultInput{
					Entry: dto.DiagnosticReportInput{
						Date: &scalarutils.Date{Year: 2025, Month: 1, Day: 1},
						Note: "Additional information",
						Media: []*dto.Media{
							{
								ID:        gofakeit.UUID(),
								MediaLink: gofakeit.URL(),
								Name:      gofakeit.BeerName(),
							},
						},
						Findings: "ASCUS or greater",
					},
					ServiceRequestID: gofakeit.UUID(),
				}

				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to record test result",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return organization, nil
					})
				mh.FHIR.EXPECT().CreateFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
						return observation, nil
					})
				mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
						return conceptpayload, nil
					})
				mh.FHIR.EXPECT().CreateFHIRDiagnosticReport(mock.Anything, mock.Anything).
					RunAndReturn(func(context1 context.Context, input *domain.FHIRDiagnosticReportInput) (*domain.FHIRDiagnosticReport, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to record test result - no test result findings",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return organization, nil
					})
				mh.FHIR.EXPECT().CreateFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
						return observation, nil
					})
				mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
						return conceptpayload, nil
					})
				mh.FHIR.EXPECT().CreateFHIRDiagnosticReport(mock.Anything, mock.Anything).
					RunAndReturn(func(context1 context.Context, input *domain.FHIRDiagnosticReportInput) (*domain.FHIRDiagnosticReport, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			_, err := clinicalUsecase.RecordTestResult(args.ctx, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.RecordTestResult() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCasesClinicalImpl_CreateLabOrder(t *testing.T) {
	ctx := usecaseMock.AddTenantIdentifierContext(context.Background())
	labOrderInput := &dto.IntraLabOrderInput{
		Code:        gofakeit.UUID(),
		Name:        gofakeit.UUID(),
		EncounterID: gofakeit.UUID(),
		Patient: dto.Patient{
			ID:   gofakeit.UUID(),
			Name: gofakeit.BeerName(),
		},
		ObservationID: gofakeit.UUID(),
		UsageContext:  dto.CervicalCancerScreeningTypeEnum,
	}

	tenantIDs := &dto.TenantIdentifiers{
		OrganizationID: uuid.NewString(),
		FacilityID:     uuid.NewString(),
	}

	orgName := mock.Anything
	ID := uuid.NewString()
	organization := &domain.FHIROrganizationRelayPayload{
		Resource: &domain.FHIROrganization{
			ID:   &ID,
			Name: &orgName,
		},
	}

	patient := &domain.FHIRPatientRelayPayload{
		Resource: &domain.FHIRPatient{
			ID: &ID,
		},
	}

	type args struct {
		ctx           context.Context
		labOrderInput *dto.IntraLabOrderInput
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy case: create lab order",
			setup: func(mh *usecaseMock.Mocks) args {
				ref := mock.Anything
				servicerequest := &domain.FHIRServiceRequestRelayPayload{
					Resource: &domain.FHIRServiceRequest{
						ID:         &ID,
						Identifier: []*domain.FHIRIdentifier{},
						Status:     domain.ServiceRequestStatusActive,
						Intent:     domain.ServiceRequestIntentDirective,
						Subject: &domain.FHIRReference{
							ID:        &ID,
							Reference: &ref,
							Display:   mock.Anything,
						},
						Encounter: &domain.FHIRReference{
							ID:        &ID,
							Display:   mock.Anything,
							Reference: &ref,
						},
						Code: &domain.FHIRCodeableReference{
							Concept: &domain.FHIRCodeableConcept{
								Coding: []*domain.FHIRCoding{},
							},
						},
					},
				}

				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return patient, nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return organization, nil
					})
				mh.FHIR.EXPECT().CreateFHIRServiceRequest(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRServiceRequestInput) (*domain.FHIRServiceRequestRelayPayload, error) {
						return servicerequest, nil
					})
				mh.PubSub.EXPECT().NotifyCreatePatientReferralTask(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, data dto.PatientReferralTaskPayload) error {
						return nil
					})

				return args{ctx: ctx, labOrderInput: labOrderInput}
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to get tenant tags",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return patient, nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return nil, fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, labOrderInput: labOrderInput}
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get patient by id",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, labOrderInput: labOrderInput}
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get facility id from context",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return patient, nil
					})

				return args{
					ctx: context.Background(),
					labOrderInput: &dto.IntraLabOrderInput{
						Code:        gofakeit.UUID(),
						Name:        gofakeit.UUID(),
						EncounterID: gofakeit.UUID(),
						Patient: dto.Patient{
							ID:   gofakeit.UUID(),
							Name: gofakeit.BeerName(),
						},
						ObservationID: gofakeit.UUID(),
						UsageContext:  dto.CervicalCancerScreeningTypeEnum,
					},
				}
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to create service request",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return patient, nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return organization, nil
					})
				mh.FHIR.EXPECT().CreateFHIRServiceRequest(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRServiceRequestInput) (*domain.FHIRServiceRequestRelayPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{
					ctx: usecaseMock.AddTenantIdentifierContext(context.Background()),
					labOrderInput: &dto.IntraLabOrderInput{
						Code:        gofakeit.UUID(),
						Name:        gofakeit.UUID(),
						EncounterID: gofakeit.UUID(),
						Patient: dto.Patient{
							ID:   gofakeit.UUID(),
							Name: gofakeit.BeerName(),
						},
						ObservationID: gofakeit.UUID(),
						UsageContext:  dto.CervicalCancerScreeningTypeEnum,
					},
				}
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to create intra-referral task",
			setup: func(mh *usecaseMock.Mocks) args {
				ref := mock.Anything
				servicerequest := &domain.FHIRServiceRequestRelayPayload{
					Resource: &domain.FHIRServiceRequest{
						ID:         &ID,
						Identifier: []*domain.FHIRIdentifier{},
						Status:     domain.ServiceRequestStatusActive,
						Intent:     domain.ServiceRequestIntentDirective,
						Subject: &domain.FHIRReference{
							ID:        &ID,
							Reference: &ref,
							Display:   mock.Anything,
						},
						Encounter: &domain.FHIRReference{
							ID:        &ID,
							Display:   mock.Anything,
							Reference: &ref,
						},
						Code: &domain.FHIRCodeableReference{
							Concept: &domain.FHIRCodeableConcept{
								Coding: []*domain.FHIRCoding{},
							},
						},
					},
				}

				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return patient, nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return organization, nil
					})
				mh.FHIR.EXPECT().CreateFHIRServiceRequest(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRServiceRequestInput) (*domain.FHIRServiceRequestRelayPayload, error) {
						return servicerequest, nil
					})
				mh.PubSub.EXPECT().NotifyCreatePatientReferralTask(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, data dto.PatientReferralTaskPayload) error {
						return fmt.Errorf("an error occurred")
					})

				return args{
					ctx: usecaseMock.AddTenantIdentifierContext(context.Background()),
					labOrderInput: &dto.IntraLabOrderInput{
						Code:        gofakeit.UUID(),
						Name:        gofakeit.UUID(),
						EncounterID: gofakeit.UUID(),
						Patient: dto.Patient{
							ID:   gofakeit.UUID(),
							Name: gofakeit.BeerName(),
						},
						ObservationID: gofakeit.UUID(),
						UsageContext:  dto.CervicalCancerScreeningTypeEnum,
					},
				}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			_, err := clinicalUsecase.CreateLabOrder(args.ctx, args.labOrderInput)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.CreateLabOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCasesClinicalImpl_RecordLabTests(t *testing.T) {
	ctx := usecaseMock.AddTenantIdentifierContext(context.Background())

	ID := uuid.NewString()
	ref := "Reference/1234"
	encounter := &domain.FHIREncounterRelayPayload{
		Resource: &domain.FHIREncounter{
			ID:     &ID,
			Status: domain.EncounterStatusEnumInProgress,
			Subject: &domain.FHIRReference{
				ID:        &ID,
				Reference: &ref,
				Display:   mock.Anything,
			},
		},
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

	status := domain.ObservationStatusEnumRegistered
	valueString := gofakeit.BeerName()
	noteText := scalarutils.Markdown(mock.Anything)
	instant := scalarutils.Instant(time.Now().GoString())
	system := scalarutils.URI(gofakeit.BeerName())
	code := scalarutils.Code("code")
	observation := &domain.FHIRObservation{
		ID: &ID,
		Text: &domain.FHIRNarrative{
			Div: scalarutils.XHTML(gofakeit.BeerName()),
		},
		ValueString: &valueString,
		Status:      &status,
		Code: &domain.FHIRCodeableConcept{
			Coding: []*domain.FHIRCoding{
				{
					System:  &system,
					Display: gofakeit.BeerName(),
					Code:    &code,
				},
			},
		},
		Subject: &domain.FHIRReference{
			ID:      &ID,
			Display: gofakeit.BeerName(),
		},
		Encounter: &domain.FHIRReference{
			Display: gofakeit.BeerName(),
		},
		Note: []*domain.FHIRAnnotation{
			{
				ID:   &ID,
				Text: &noteText,
			},
		},
		Category: []*domain.FHIRCodeableConcept{
			{
				ID: &ID,
				Coding: []*domain.FHIRCoding{
					{
						ID: &ID,
					},
				},
				Text: mock.Anything,
			},
		},
		EffectiveInstant: &instant,
	}

	type args struct {
		ctx   context.Context
		input dto.TestInput
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy case: successfully records mammogram test",
			setup: func(mh *usecaseMock.Mocks) args {
				input := dto.TestInput{
					TestType: dto.LabTestTypeEnum(dto.MammogramTest),
					Input: dto.DiagnosticReportInput{
						Date:        &scalarutils.Date{Year: 2025, Month: 1, Day: 1},
						EncounterID: gofakeit.UUID(),
						Findings:    "False positive",
					},
				}

				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
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
				mh.FHIR.EXPECT().CreateFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
						return observation, nil
					})
				mh.FHIR.EXPECT().CreateFHIRDiagnosticReport(mock.Anything, mock.Anything).
					RunAndReturn(func(context1 context.Context, input *domain.FHIRDiagnosticReportInput) (*domain.FHIRDiagnosticReport, error) {
						issued := mock.Anything
						diagnosticsReport := &domain.FHIRDiagnosticReport{
							ID: &ID,
							Encounter: &domain.FHIRReference{
								ID: &ID,
							},
							Subject: &domain.FHIRReference{
								ID: &ID,
							},
							Status: domain.DiagnosticReportStatusAmended,
							Issued: &issued,
						}
						return diagnosticsReport, nil
					})

				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "Happy case: successfully records MRI test",
			setup: func(mh *usecaseMock.Mocks) args {
				input := dto.TestInput{
					TestType: dto.LabTestTypeEnum(dto.MRITest),
					Input: dto.DiagnosticReportInput{
						Date:        &scalarutils.Date{Year: 2025, Month: 1, Day: 1},
						EncounterID: gofakeit.UUID(),
						Findings:    "False positive",
					},
				}

				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
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
				mh.FHIR.EXPECT().CreateFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
						return observation, nil
					})
				mh.FHIR.EXPECT().CreateFHIRDiagnosticReport(mock.Anything, mock.Anything).
					RunAndReturn(func(context1 context.Context, input *domain.FHIRDiagnosticReportInput) (*domain.FHIRDiagnosticReport, error) {
						issued := mock.Anything
						diagnosticsReport := &domain.FHIRDiagnosticReport{
							ID: &ID,
							Encounter: &domain.FHIRReference{
								ID: &ID,
							},
							Subject: &domain.FHIRReference{
								ID: &ID,
							},
							Status: domain.DiagnosticReportStatusAmended,
							Issued: &issued,
						}
						return diagnosticsReport, nil
					})

				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "Happy case: successfully records CBE test",
			setup: func(mh *usecaseMock.Mocks) args {
				input := dto.TestInput{
					TestType: dto.LabTestTypeEnum(dto.CBETest),
					Input: dto.DiagnosticReportInput{
						Date:        &scalarutils.Date{Year: 2025, Month: 1, Day: 1},
						EncounterID: gofakeit.UUID(),
						Findings:    "False positive",
					},
				}

				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
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
				mh.FHIR.EXPECT().CreateFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
						return observation, nil
					})
				mh.FHIR.EXPECT().CreateFHIRDiagnosticReport(mock.Anything, mock.Anything).
					RunAndReturn(func(context1 context.Context, input *domain.FHIRDiagnosticReportInput) (*domain.FHIRDiagnosticReport, error) {
						issued := mock.Anything
						diagnosticsReport := &domain.FHIRDiagnosticReport{
							ID: &ID,
							Encounter: &domain.FHIRReference{
								ID: &ID,
							},
							Subject: &domain.FHIRReference{
								ID: &ID,
							},
							Status: domain.DiagnosticReportStatusAmended,
							Issued: &issued,
						}
						return diagnosticsReport, nil
					})

				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "Happy case: successfully records Papsmear test",
			setup: func(mh *usecaseMock.Mocks) args {
				input := dto.TestInput{
					TestType: dto.LabTestTypeEnum(dto.PapSmearTest),
					Input: dto.DiagnosticReportInput{
						Date:        &scalarutils.Date{Year: 2025, Month: 1, Day: 1},
						EncounterID: gofakeit.UUID(),
						Findings:    "False positive",
					},
				}

				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
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
				mh.FHIR.EXPECT().CreateFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
						return observation, nil
					})
				mh.FHIR.EXPECT().CreateFHIRDiagnosticReport(mock.Anything, mock.Anything).
					RunAndReturn(func(context1 context.Context, input *domain.FHIRDiagnosticReportInput) (*domain.FHIRDiagnosticReport, error) {
						issued := mock.Anything
						diagnosticsReport := &domain.FHIRDiagnosticReport{
							ID: &ID,
							Encounter: &domain.FHIRReference{
								ID: &ID,
							},
							Subject: &domain.FHIRReference{
								ID: &ID,
							},
							Status: domain.DiagnosticReportStatusAmended,
							Issued: &issued,
						}
						return diagnosticsReport, nil
					})

				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "Happy case: successfully records Biopsy test",
			setup: func(mh *usecaseMock.Mocks) args {
				input := dto.TestInput{
					TestType: dto.LabTestTypeEnum(dto.BiopsyTest),
					Input: dto.DiagnosticReportInput{
						Date:        &scalarutils.Date{Year: 2025, Month: 1, Day: 1},
						EncounterID: gofakeit.UUID(),
						Findings:    "False positive",
					},
				}

				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
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
				mh.FHIR.EXPECT().CreateFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
						return observation, nil
					})
				mh.FHIR.EXPECT().CreateFHIRDiagnosticReport(mock.Anything, mock.Anything).
					RunAndReturn(func(context1 context.Context, input *domain.FHIRDiagnosticReportInput) (*domain.FHIRDiagnosticReport, error) {
						issued := mock.Anything
						diagnosticsReport := &domain.FHIRDiagnosticReport{
							ID: &ID,
							Encounter: &domain.FHIRReference{
								ID: &ID,
							},
							Subject: &domain.FHIRReference{
								ID: &ID,
							},
							Status: domain.DiagnosticReportStatusAmended,
							Issued: &issued,
						}
						return diagnosticsReport, nil
					})

				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "Happy case: successfully records Ultrasound test",
			setup: func(mh *usecaseMock.Mocks) args {
				input := dto.TestInput{
					TestType: dto.LabTestTypeEnum(dto.UltrasoundTest),
					Input: dto.DiagnosticReportInput{
						Date:        &scalarutils.Date{Year: 2025, Month: 1, Day: 1},
						EncounterID: gofakeit.UUID(),
						Findings:    "False positive",
					},
				}

				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
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
				mh.FHIR.EXPECT().CreateFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
						return observation, nil
					})
				mh.FHIR.EXPECT().CreateFHIRDiagnosticReport(mock.Anything, mock.Anything).
					RunAndReturn(func(context1 context.Context, input *domain.FHIRDiagnosticReportInput) (*domain.FHIRDiagnosticReport, error) {
						issued := mock.Anything
						diagnosticsReport := &domain.FHIRDiagnosticReport{
							ID: &ID,
							Encounter: &domain.FHIRReference{
								ID: &ID,
							},
							Subject: &domain.FHIRReference{
								ID: &ID,
							},
							Status: domain.DiagnosticReportStatusAmended,
							Issued: &issued,
						}
						return diagnosticsReport, nil
					})

				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "Happy case: successfully records PSA test ProstaticSerumAntigenTest",
			setup: func(mh *usecaseMock.Mocks) args {
				input := dto.TestInput{
					TestType: dto.LabTestTypeEnum(dto.ProstaticSerumAntigenTest),
					Input: dto.DiagnosticReportInput{
						Date:        &scalarutils.Date{Year: 2025, Month: 1, Day: 1},
						EncounterID: gofakeit.UUID(),
						Findings:    "False positive",
					},
				}

				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
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
				mh.FHIR.EXPECT().CreateFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
						return observation, nil
					})
				mh.FHIR.EXPECT().CreateFHIRDiagnosticReport(mock.Anything, mock.Anything).
					RunAndReturn(func(context1 context.Context, input *domain.FHIRDiagnosticReportInput) (*domain.FHIRDiagnosticReport, error) {
						issued := mock.Anything
						diagnosticsReport := &domain.FHIRDiagnosticReport{
							ID: &ID,
							Encounter: &domain.FHIRReference{
								ID: &ID,
							},
							Subject: &domain.FHIRReference{
								ID: &ID,
							},
							Status: domain.DiagnosticReportStatusAmended,
							Issued: &issued,
						}
						return diagnosticsReport, nil
					})

				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "Happy case: successfully records PSA Wholeblood test",
			setup: func(mh *usecaseMock.Mocks) args {
				input := dto.TestInput{
					TestType: dto.LabTestTypeEnum(dto.ProstaticSerumAntigen),
					Input: dto.DiagnosticReportInput{
						Date:        &scalarutils.Date{Year: 2025, Month: 1, Day: 1},
						EncounterID: gofakeit.UUID(),
						Findings:    "False positive",
					},
				}

				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
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
				mh.FHIR.EXPECT().CreateFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
						return observation, nil
					})
				mh.FHIR.EXPECT().CreateFHIRDiagnosticReport(mock.Anything, mock.Anything).
					RunAndReturn(func(context1 context.Context, input *domain.FHIRDiagnosticReportInput) (*domain.FHIRDiagnosticReport, error) {
						issued := mock.Anything
						diagnosticsReport := &domain.FHIRDiagnosticReport{
							ID: &ID,
							Encounter: &domain.FHIRReference{
								ID: &ID,
							},
							Subject: &domain.FHIRReference{
								ID: &ID,
							},
							Status: domain.DiagnosticReportStatusAmended,
							Issued: &issued,
						}
						return diagnosticsReport, nil
					})

				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "Happy case: successfully records IHC progesterone receptor test",
			setup: func(mh *usecaseMock.Mocks) args {
				input := dto.TestInput{
					TestType: dto.LabTestTypeEnum("IHC_PROGESTERONE_RECEPTOR"),
					Input: dto.DiagnosticReportInput{
						Date:        &scalarutils.Date{Year: 2025, Month: 1, Day: 1},
						EncounterID: gofakeit.UUID(),
						Findings:    "False positive",
					},
				}

				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
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
				mh.FHIR.EXPECT().CreateFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
						return observation, nil
					})
				mh.FHIR.EXPECT().CreateFHIRDiagnosticReport(mock.Anything, mock.Anything).
					RunAndReturn(func(context1 context.Context, input *domain.FHIRDiagnosticReportInput) (*domain.FHIRDiagnosticReport, error) {
						issued := mock.Anything
						diagnosticsReport := &domain.FHIRDiagnosticReport{
							ID: &ID,
							Encounter: &domain.FHIRReference{
								ID: &ID,
							},
							Subject: &domain.FHIRReference{
								ID: &ID,
							},
							Status: domain.DiagnosticReportStatusAmended,
							Issued: &issued,
						}
						return diagnosticsReport, nil
					})

				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "Happy case: successfully records IHC estrogen receptor test",
			setup: func(mh *usecaseMock.Mocks) args {
				input := dto.TestInput{
					TestType: dto.LabTestTypeEnum("IHC_ESTROGEN_RECEPTOR"),
					Input: dto.DiagnosticReportInput{
						Date:        &scalarutils.Date{Year: 2025, Month: 1, Day: 1},
						EncounterID: gofakeit.UUID(),
						Findings:    "False positive",
					},
				}

				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
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
				mh.FHIR.EXPECT().CreateFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
						return observation, nil
					})
				mh.FHIR.EXPECT().CreateFHIRDiagnosticReport(mock.Anything, mock.Anything).
					RunAndReturn(func(context1 context.Context, input *domain.FHIRDiagnosticReportInput) (*domain.FHIRDiagnosticReport, error) {
						issued := mock.Anything
						diagnosticsReport := &domain.FHIRDiagnosticReport{
							ID: &ID,
							Encounter: &domain.FHIRReference{
								ID: &ID,
							},
							Subject: &domain.FHIRReference{
								ID: &ID,
							},
							Status: domain.DiagnosticReportStatusAmended,
							Issued: &issued,
						}
						return diagnosticsReport, nil
					})

				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "Happy case: successfully records IHC HER2 test",
			setup: func(mh *usecaseMock.Mocks) args {
				input := dto.TestInput{
					TestType: dto.LabTestTypeEnum("IHC_HER2"),
					Input: dto.DiagnosticReportInput{
						Date:        &scalarutils.Date{Year: 2025, Month: 1, Day: 1},
						EncounterID: gofakeit.UUID(),
						Findings:    "False positive",
					},
				}

				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
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
				mh.FHIR.EXPECT().CreateFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
						return observation, nil
					})
				mh.FHIR.EXPECT().CreateFHIRDiagnosticReport(mock.Anything, mock.Anything).
					RunAndReturn(func(context1 context.Context, input *domain.FHIRDiagnosticReportInput) (*domain.FHIRDiagnosticReport, error) {
						issued := mock.Anything
						diagnosticsReport := &domain.FHIRDiagnosticReport{
							ID: &ID,
							Encounter: &domain.FHIRReference{
								ID: &ID,
							},
							Subject: &domain.FHIRReference{
								ID: &ID,
							},
							Status: domain.DiagnosticReportStatusAmended,
							Issued: &issued,
						}
						return diagnosticsReport, nil
					})

				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "Happy case: successfully records IHC KI67 test",
			setup: func(mh *usecaseMock.Mocks) args {
				input := dto.TestInput{
					TestType: dto.LabTestTypeEnum("IHC_KI67"),
					Input: dto.DiagnosticReportInput{
						Date:        &scalarutils.Date{Year: 2025, Month: 1, Day: 1},
						EncounterID: gofakeit.UUID(),
						Findings:    "False positive",
					},
				}

				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
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
				mh.FHIR.EXPECT().CreateFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
						return observation, nil
					})
				mh.FHIR.EXPECT().CreateFHIRDiagnosticReport(mock.Anything, mock.Anything).
					RunAndReturn(func(context1 context.Context, input *domain.FHIRDiagnosticReportInput) (*domain.FHIRDiagnosticReport, error) {
						issued := mock.Anything
						diagnosticsReport := &domain.FHIRDiagnosticReport{
							ID: &ID,
							Encounter: &domain.FHIRReference{
								ID: &ID,
							},
							Subject: &domain.FHIRReference{
								ID: &ID,
							},
							Status: domain.DiagnosticReportStatusAmended,
							Issued: &issued,
						}
						return diagnosticsReport, nil
					})

				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "Happy case: successfully records whole blood test",
			setup: func(mh *usecaseMock.Mocks) args {
				input := dto.TestInput{
					TestType: dto.WholeBloodTest,
					Input: dto.DiagnosticReportInput{
						Date:        &scalarutils.Date{Year: 2025, Month: 1, Day: 1},
						EncounterID: gofakeit.UUID(),
						Findings:    "False positive",
					},
				}

				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
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
				mh.FHIR.EXPECT().CreateFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
						return observation, nil
					})
				mh.FHIR.EXPECT().CreateFHIRDiagnosticReport(mock.Anything, mock.Anything).
					RunAndReturn(func(context1 context.Context, input *domain.FHIRDiagnosticReportInput) (*domain.FHIRDiagnosticReport, error) {
						issued := mock.Anything
						diagnosticsReport := &domain.FHIRDiagnosticReport{
							ID: &ID,
							Encounter: &domain.FHIRReference{
								ID: &ID,
							},
							Subject: &domain.FHIRReference{
								ID: &ID,
							},
							Status: domain.DiagnosticReportStatusAmended,
							Issued: &issued,
						}
						return diagnosticsReport, nil
					})

				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to record observation",
			setup: func(mh *usecaseMock.Mocks) args {
				input := dto.TestInput{
					TestType: dto.LabTestTypeEnum(dto.MammogramTest),
					Input: dto.DiagnosticReportInput{
						Date:        &scalarutils.Date{Year: 2025, Month: 1, Day: 1},
						EncounterID: gofakeit.UUID(),
						Findings:    "False positive",
					},
				}

				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
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
				mh.FHIR.EXPECT().CreateFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
						return nil, fmt.Errorf("error")
					})

				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
		{
			name: "Sad case: invalid input",
			setup: func(mh *usecaseMock.Mocks) args {
				input := dto.TestInput{
					Input: dto.DiagnosticReportInput{
						Date:        &scalarutils.Date{Year: 2025, Month: 1, Day: 1},
						EncounterID: gofakeit.UUID(),
						Findings:    "False positive",
					},
				}

				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
		{
			name: "Sad case: fails to record test",
			setup: func(mh *usecaseMock.Mocks) args {
				input := dto.TestInput{
					TestType: dto.LabTestTypeEnum(dto.ProstaticSerumAntigen),
					Input: dto.DiagnosticReportInput{
						Date:        &scalarutils.Date{Year: 2025, Month: 1, Day: 1},
						EncounterID: gofakeit.UUID(),
						Findings:    "False positive",
					},
				}

				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
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
				mh.FHIR.EXPECT().CreateFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
						return observation, nil
					})
				mh.FHIR.EXPECT().CreateFHIRDiagnosticReport(mock.Anything, mock.Anything).
					RunAndReturn(func(context1 context.Context, input *domain.FHIRDiagnosticReportInput) (*domain.FHIRDiagnosticReport, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			clinicalUsecase, mockService := usecaseMock.SetupMocks(t)
			args := tt.setup(&mockService)

			_, err := clinicalUsecase.RecordTests(args.ctx, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.RecordLabTests() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCasesClinicalImpl_DeleteTest(t *testing.T) {
	ctx := usecaseMock.AddTenantIdentifierContext(context.Background())

	ID := uuid.NewString()
	diagnosticReport := &domain.PagedFHIRDiagnosticReport{
		DiagnosticReport: []domain.FHIRDiagnosticReport{
			{
				ID: &ID,
			},
		},
	}

	tenantIDs := &dto.TenantIdentifiers{
		OrganizationID: uuid.NewString(),
		FacilityID:     uuid.NewString(),
	}

	type args struct {
		ctx           context.Context
		observationID string
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy case: delete tests",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchFHIRDiagnosticReport(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRDiagnosticReport, error) {
						return diagnosticReport, nil
					})
				mh.FHIR.EXPECT().DeleteFHIRResource(mock.Anything, "DiagnosticReport", mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, resourceID string) (bool, error) {
						return true, nil
					})
				mh.FHIR.EXPECT().DeleteFHIRResource(mock.Anything, "Observation", mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, resourceID string) (bool, error) {
						return true, nil
					})
				return args{ctx: ctx, observationID: uuid.NewString()}
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad case: unable to get tenant identifiers",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return nil, fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, observationID: uuid.NewString()}
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case: unable to search diagnostic report",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchFHIRDiagnosticReport(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRDiagnosticReport, error) {
						return nil, fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, observationID: uuid.NewString()}
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case: unable to delete diagnostic report",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchFHIRDiagnosticReport(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRDiagnosticReport, error) {
						return diagnosticReport, nil
					})
				mh.FHIR.EXPECT().DeleteFHIRResource(mock.Anything, "DiagnosticReport", mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, resourceID string) (bool, error) {
						return false, fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, observationID: uuid.NewString()}
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case: unable to delete observation",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchFHIRDiagnosticReport(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRDiagnosticReport, error) {
						return diagnosticReport, nil
					})
				mh.FHIR.EXPECT().DeleteFHIRResource(mock.Anything, "DiagnosticReport", mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, resourceID string) (bool, error) {
						return true, nil
					})
				mh.FHIR.EXPECT().DeleteFHIRResource(mock.Anything, "Observation", mock.Anything).
					RunAndReturn(func(ctx context.Context, resourceType, resourceID string) (bool, error) {
						return false, fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, observationID: uuid.NewString()}
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			_, err := clinicalUsecase.DeleteTest(args.ctx, args.observationID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.DeleteTest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestReviewFollowUpTask(t *testing.T) {
	diagnosticResportOutput := &dto.DiagnosticReport{
		ID:          gofakeit.UUID(),
		Status:      dto.ObservationStatusFinal,
		PatientID:   gofakeit.UUID(),
		EncounterID: gofakeit.UUID(),
		Issued:      gofakeit.UUID(),
		Result:      []*dto.Observation{},
		Media:       []*dto.Media{},
		Conclusion:  gofakeit.UUID(),
	}

	obs := &dto.ObservationPayload{
		ObservationInput: dto.ObservationInput{
			Status:             dto.ObservationStatusCancelled,
			Concept:            dto.ObservationConceptEnumBMI,
			EncounterID:        gofakeit.UUID(),
			Note:               gofakeit.UUID(),
			Value:              gofakeit.UUID(),
			ObservationSubType: dto.HPV_PCR_DNA,
			UsageContext:       dto.BreastCancerScreeningTypeEnum,
		},
		VitalSignsConceptID: gofakeit.UUID(),
		ServiceRequestID:    gofakeit.UUID(),
	}

	testValue := gofakeit.UUID()

	type args struct {
		ctx                    context.Context
		diagnosticReportOutput *dto.DiagnosticReport
		observation            *dto.ObservationPayload
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy case create test review followup task",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})

				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return organization, nil
					})

				mh.FHIR.EXPECT().GetFHIRServiceRequest(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRServiceRequestRelayPayload, error) {
						return &domain.FHIRServiceRequestRelayPayload{
							Resource: &domain.FHIRServiceRequest{
								Status: domain.ServiceRequestStatusActive,
								Code: &domain.FHIRCodeableReference{
									Concept: &domain.FHIRCodeableConcept{
										ID: new(string),
										Coding: []*domain.FHIRCoding{
											{
												ID:           new(string),
												System:       (*scalarutils.URI)(&testValue),
												Version:      new(string),
												Code:         (*scalarutils.Code)(&testValue),
												Display:      testValue,
												UserSelected: new(bool),
											},
										},
										Text: "",
									},
								},
								Subject: &domain.FHIRReference{
									Display: gofakeit.Name(),
								},
							},
						}, nil
					})

				mh.PubSub.EXPECT().NotifyCreateFollowUpTask(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, data *domain.FHIRTaskInput) error {
						return nil
					})

				mh.FHIR.EXPECT().PatchFHIRServiceRequest(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string, input domain.FHIRServiceRequestInput) (*domain.FHIRServiceRequestRelayPayload, error) {
						return &domain.FHIRServiceRequestRelayPayload{
							Resource: &domain.FHIRServiceRequest{
								Status: domain.ServiceRequestStatusActive,
							},
						}, nil
					})

				return args{ctx: usecaseMock.AddTenantIdentifierContext(context.Background()), diagnosticReportOutput: diagnosticResportOutput, observation: obs}
			},
			wantErr: false,
		},
		{
			name: "Happy case create test review followup task (no usage context specified)",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})

				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return organization, nil
					})

				mh.FHIR.EXPECT().GetFHIRServiceRequest(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRServiceRequestRelayPayload, error) {
						return &domain.FHIRServiceRequestRelayPayload{
							Resource: &domain.FHIRServiceRequest{
								Status: domain.ServiceRequestStatusActive,
								Code: &domain.FHIRCodeableReference{
									Concept: &domain.FHIRCodeableConcept{
										ID: new(string),
										Coding: []*domain.FHIRCoding{
											{
												ID:           new(string),
												System:       (*scalarutils.URI)(&testValue),
												Version:      new(string),
												Code:         (*scalarutils.Code)(&testValue),
												Display:      testValue,
												UserSelected: new(bool),
											},
										},
										Text: "",
									},
								},
								Subject: &domain.FHIRReference{
									Display: gofakeit.Name(),
								},
							},
						}, nil
					})

				mh.PubSub.EXPECT().NotifyCreateFollowUpTask(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, data *domain.FHIRTaskInput) error {
						return nil
					})

				mh.FHIR.EXPECT().PatchFHIRServiceRequest(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string, input domain.FHIRServiceRequestInput) (*domain.FHIRServiceRequestRelayPayload, error) {
						return &domain.FHIRServiceRequestRelayPayload{
							Resource: &domain.FHIRServiceRequest{
								Status: domain.ServiceRequestStatusActive,
							},
						}, nil
					})

				return args{ctx: usecaseMock.AddTenantIdentifierContext(context.Background()), diagnosticReportOutput: diagnosticResportOutput, observation: &dto.ObservationPayload{
					ObservationInput: dto.ObservationInput{
						UsageContext: dto.ScreeningTypeEnum("test"),
					},
					VitalSignsConceptID: "",
					ServiceRequestID:    "",
				}}
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to get meta tags",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return nil, fmt.Errorf("error")
					})

				return args{ctx: usecaseMock.AddTenantIdentifierContext(context.Background()), diagnosticReportOutput: diagnosticResportOutput, observation: obs}
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get service request",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})

				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return organization, nil
					})

				mh.FHIR.EXPECT().GetFHIRServiceRequest(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRServiceRequestRelayPayload, error) {
						return nil, fmt.Errorf("error")
					})

				return args{ctx: usecaseMock.AddTenantIdentifierContext(context.Background()), diagnosticReportOutput: diagnosticResportOutput, observation: obs}
			},
			wantErr: true,
		},
		{
			name: "Happy case create test review followup task",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})

				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return organization, nil
					})

				mh.FHIR.EXPECT().GetFHIRServiceRequest(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRServiceRequestRelayPayload, error) {
						return &domain.FHIRServiceRequestRelayPayload{
							Resource: &domain.FHIRServiceRequest{
								Status: domain.ServiceRequestStatusActive,
								Code: &domain.FHIRCodeableReference{
									Concept: &domain.FHIRCodeableConcept{
										ID: new(string),
										Coding: []*domain.FHIRCoding{
											{
												ID:           new(string),
												System:       (*scalarutils.URI)(&testValue),
												Version:      new(string),
												Code:         (*scalarutils.Code)(&testValue),
												Display:      testValue,
												UserSelected: new(bool),
											},
										},
										Text: "",
									},
								},
								Subject: &domain.FHIRReference{
									Display: gofakeit.Name(),
								},
							},
						}, nil
					})

				mh.PubSub.EXPECT().NotifyCreateFollowUpTask(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, data *domain.FHIRTaskInput) error {
						return fmt.Errorf("error")
					})

				return args{ctx: usecaseMock.AddTenantIdentifierContext(context.Background()), diagnosticReportOutput: diagnosticResportOutput, observation: obs}
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to patch service request",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})

				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return organization, nil
					})

				mh.FHIR.EXPECT().GetFHIRServiceRequest(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRServiceRequestRelayPayload, error) {
						return &domain.FHIRServiceRequestRelayPayload{
							Resource: &domain.FHIRServiceRequest{
								Status: domain.ServiceRequestStatusActive,
								Code: &domain.FHIRCodeableReference{
									Concept: &domain.FHIRCodeableConcept{
										ID: new(string),
										Coding: []*domain.FHIRCoding{
											{
												ID:           new(string),
												System:       (*scalarutils.URI)(&testValue),
												Version:      new(string),
												Code:         (*scalarutils.Code)(&testValue),
												Display:      testValue,
												UserSelected: new(bool),
											},
										},
										Text: "",
									},
								},
								Subject: &domain.FHIRReference{
									Display: gofakeit.Name(),
								},
							},
						}, nil
					})

				mh.PubSub.EXPECT().NotifyCreateFollowUpTask(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, data *domain.FHIRTaskInput) error {
						return nil
					})

				mh.FHIR.EXPECT().PatchFHIRServiceRequest(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string, input domain.FHIRServiceRequestInput) (*domain.FHIRServiceRequestRelayPayload, error) {
						return nil, fmt.Errorf("error")
					})

				return args{ctx: usecaseMock.AddTenantIdentifierContext(context.Background()), diagnosticReportOutput: diagnosticResportOutput, observation: obs}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			if err := clinicalUsecase.TestReviewFollowUpTask(args.ctx, args.diagnosticReportOutput, args.observation); (err != nil) != tt.wantErr {
				t.Errorf("ClinicalImpl.TestReviewFollowUpTask() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
