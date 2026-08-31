package clinical_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/scalarutils"
	"github.com/savannahghi/serverutils"
	"github.com/stretchr/testify/mock"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/common"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/empower-clinical/pkg/clinical/domain"
	"github.com/savannahghi/empower-clinical/pkg/clinical/usecases/clinical"
	usecaseMock "github.com/savannahghi/empower-clinical/pkg/clinical/usecases/mock"
)

var testObservationCategorySystem = "http://terminology.hl7.org/CodeSystem/observation-category"

// addLabCategory is used to add laboratory categories for various observations records.
var addLabCategory = func(ctx context.Context, observation *domain.FHIRObservationInput) error {
	userSelected := false
	category := []*domain.FHIRCodeableConceptInput{
		{
			Coding: []*domain.FHIRCodingInput{
				{
					System:       (*scalarutils.URI)(&testObservationCategorySystem),
					Code:         "laboratory",
					Display:      "Laboratory",
					UserSelected: &userSelected,
				},
			},
			Text: "Laboratory",
		},
	}

	observation.Category = append(observation.Category, category...)
	return nil
}

// Shared variables throughout the file
var ID = uuid.NewString()

var ref = "Reference/12345"
var encounter = &domain.FHIREncounterRelayPayload{
	Resource: &domain.FHIREncounter{
		ID:     &ID,
		Status: domain.EncounterStatusEnumInProgress,
		Subject: &domain.FHIRReference{
			ID: &ID,
		},
		ServiceProvider: &domain.FHIRReference{
			Display:   gofakeit.UUID(),
			Reference: &ref,
		},
	},
}

var conceptpayload = &domain.Concept{
	ConceptClass: mock.Anything,
	DataType:     mock.Anything,
	ID:           gofakeit.UUID(),
}

var tenantIDs = &dto.TenantIdentifiers{
	OrganizationID: uuid.NewString(),
	FacilityID:     uuid.NewString(),
}

var orgName = mock.Anything
var organization = &domain.FHIROrganizationRelayPayload{
	Resource: &domain.FHIROrganization{
		ID:   &ID,
		Name: &orgName,
	},
}

var valueString = gofakeit.BeerName()
var noteText = scalarutils.Markdown(mock.Anything)
var instant = scalarutils.Instant(time.Now().GoString())
var status = domain.ObservationStatusEnumAmended
var system = scalarutils.URI(gofakeit.BeerName())
var code = scalarutils.Code("code")
var observation = &domain.FHIRObservation{
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

func TestUseCasesClinicalImpl_RecordObservation(t *testing.T) {
	ctx := usecaseMock.AddTenantIdentifierContext(context.Background())
	obsPayload := &dto.ObservationPayload{
		ObservationInput: dto.ObservationInput{
			Status:       dto.ObservationStatusFinal,
			EncounterID:  uuid.New().String(),
			Value:        "positive",
			UsageContext: dto.BreastCancerScreeningTypeEnum,
			Note:         gofakeit.Paragraph(3, 2, 20, "."),
		},
		VitalSignsConceptID: common.BMILOINCTerminologyCode,
		ServiceRequestID:    gofakeit.UUID(),
	}
	mutators := []clinical.ObservationInputMutatorFunc{addLabCategory}
	screeningType := dto.BreastCancerScreeningTypeEnum

	type args struct {
		ctx           context.Context
		obsPayload    *dto.ObservationPayload
		mutators      []clinical.ObservationInputMutatorFunc
		screeningType dto.ScreeningTypeEnum
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully record observation",
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

				return args{ctx: ctx, obsPayload: obsPayload, mutators: mutators, screeningType: screeningType}
			},
			wantErr: false,
		},
		{
			name: "Happy Case - Successfully record observation with a note",
			setup: func(mh *usecaseMock.Mocks) args {
				obsPayload.ObservationInput.Note = mock.Anything

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

				return args{ctx: ctx, obsPayload: obsPayload, mutators: mutators, screeningType: screeningType}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail validation",
			setup: func(mh *usecaseMock.Mocks) args {
				obsPayload := &dto.ObservationPayload{
					ObservationInput: dto.ObservationInput{
						EncounterID: uuid.New().String(),
						Value:       "BIRADS 1",
					},
					VitalSignsConceptID: common.BMILOINCTerminologyCode,
					ServiceRequestID:    "",
				}

				return args{ctx: ctx, obsPayload: obsPayload, mutators: mutators, screeningType: screeningType}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail validation",
			setup: func(mh *usecaseMock.Mocks) args {
				obsPayload := &dto.ObservationPayload{
					ObservationInput: dto.ObservationInput{
						Status: dto.ObservationStatusFinal,
						Value:  "1234",
					},
					VitalSignsConceptID: common.BMILOINCTerminologyCode,
					ServiceRequestID:    "",
				}

				return args{ctx: ctx, obsPayload: obsPayload, mutators: mutators, screeningType: screeningType}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get encounter",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, obsPayload: obsPayload, mutators: mutators, screeningType: screeningType}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - return a finished encounter",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						encounter := &domain.FHIREncounterRelayPayload{
							Resource: &domain.FHIREncounter{
								ID:     &ID,
								Status: domain.EncounterStatusEnumCompleted,
							},
						}
						return encounter, nil
					})

				return args{ctx: ctx, obsPayload: obsPayload, mutators: mutators, screeningType: screeningType}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get CIEL concept",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, obsPayload: obsPayload, mutators: mutators, screeningType: screeningType}
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
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, obsPayload: obsPayload, mutators: mutators, screeningType: screeningType}
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

				return args{ctx: ctx, obsPayload: obsPayload, mutators: mutators, screeningType: screeningType}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - no observation category specified",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
						return conceptpayload, nil
					})

				return args{ctx: ctx, obsPayload: obsPayload, screeningType: screeningType}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - unable to get facility id from context",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})

				return args{ctx: context.Background(), obsPayload: obsPayload, mutators: mutators, screeningType: screeningType}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			got, err := clinicalUsecase.RecordObservation(args.ctx, *args.obsPayload, args.mutators)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.RecordObservation() error = %v, wantErr %v", err, tt.wantErr)
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

func TestUseCasesClinicalImpl_RecordTemperature(t *testing.T) {
	ctx := usecaseMock.AddTenantIdentifierContext(context.Background())
	obsPayload := dto.ObservationInput{
		Status:       dto.ObservationStatusFinal,
		EncounterID:  uuid.NewString(),
		Value:        "12",
		Note:         "Test note.",
		UsageContext: dto.BreastCancerScreeningTypeEnum,
	}

	type args struct {
		ctx   context.Context
		input dto.ObservationInput
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully record temperature",
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

				return args{ctx: ctx, input: obsPayload}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to record temperature",
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

				return args{ctx: ctx, input: obsPayload}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			got, err := clinicalUsecase.RecordTemperature(args.ctx, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.RecordTemperature() error = %v, wantErr %v", err, tt.wantErr)
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

func TestUseCasesClinicalImpl_RecordMuac(t *testing.T) {
	ctx := usecaseMock.AddTenantIdentifierContext(context.Background())
	input := dto.ObservationInput{
		Status:       dto.ObservationStatusFinal,
		EncounterID:  uuid.New().String(),
		Value:        "12",
		Note:         "Test note.",
		UsageContext: dto.BreastCancerScreeningTypeEnum,
	}

	type args struct {
		ctx   context.Context
		input dto.ObservationInput
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully record muac",
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

				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to record muac",
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

			got, err := clinicalUsecase.RecordMuac(args.ctx, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.RecordMuac() error = %v, wantErr %v", err, tt.wantErr)
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

func TestUseCasesClinicalImpl_RecordOxygenSaturation(t *testing.T) {
	ctx := usecaseMock.AddTenantIdentifierContext(context.Background())
	input := dto.ObservationInput{
		Status:       dto.ObservationStatusFinal,
		EncounterID:  uuid.New().String(),
		Value:        "12",
		Note:         "Test note.",
		UsageContext: dto.BreastCancerScreeningTypeEnum,
	}
	type args struct {
		ctx   context.Context
		input dto.ObservationInput
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully record oxygen saturation",
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

				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to record oxygen saturation",
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

			got, err := clinicalUsecase.RecordOxygenSaturation(args.ctx, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.RecordOxygenSaturation() error = %v, wantErr %v", err, tt.wantErr)
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

func TestUseCasesClinicalImpl_RecordHeight(t *testing.T) {
	ctx := usecaseMock.AddTenantIdentifierContext(context.Background())
	input := dto.ObservationInput{
		Status:       dto.ObservationStatusFinal,
		EncounterID:  uuid.New().String(),
		Value:        "185.21",
		Note:         "Test note.",
		UsageContext: dto.BreastCancerScreeningTypeEnum,
	}

	type args struct {
		ctx   context.Context
		input dto.ObservationInput
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully record height",
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

				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to record height",
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

			got, err := clinicalUsecase.RecordHeight(args.ctx, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.RecordHeight() error = %v, wantErr %v", err, tt.wantErr)
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

func TestUseCasesClinicalImpl_RecordWeight(t *testing.T) {
	ctx := usecaseMock.AddTenantIdentifierContext(context.Background())
	input := dto.ObservationInput{
		Status:       dto.ObservationStatusFinal,
		EncounterID:  uuid.New().String(),
		Value:        "185.21",
		Note:         "Test note.",
		UsageContext: dto.BreastCancerScreeningTypeEnum,
	}

	type args struct {
		ctx   context.Context
		input dto.ObservationInput
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully record weight",
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

				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to record weight",
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

			got, err := clinicalUsecase.RecordWeight(args.ctx, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.RecordWeight() error = %v, wantErr %v", err, tt.wantErr)
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

func TestUseCasesClinicalImpl_RecordRespiratoryRate(t *testing.T) {
	ctx := usecaseMock.AddTenantIdentifierContext(context.Background())
	input := dto.ObservationInput{
		Status:       dto.ObservationStatusFinal,
		EncounterID:  uuid.New().String(),
		Value:        "185.21",
		Note:         "Test note.",
		UsageContext: dto.BreastCancerScreeningTypeEnum,
	}

	type args struct {
		ctx   context.Context
		input dto.ObservationInput
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully record respiratory rate",
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

				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to record respiratory rate",
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

			got, err := clinicalUsecase.RecordRespiratoryRate(args.ctx, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.RecordRespiratoryRate() error = %v, wantErr %v", err, tt.wantErr)
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

func TestUseCasesClinicalImpl_RecordPulseRate(t *testing.T) {
	ctx := usecaseMock.AddTenantIdentifierContext(context.Background())
	input := dto.ObservationInput{
		Status:      dto.ObservationStatusFinal,
		EncounterID: uuid.New().String(),
		Value:       "185.21",
	}

	type args struct {
		ctx   context.Context
		input dto.ObservationInput
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully record pulse rate",
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

				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to record pulse rate",
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

			got, err := clinicalUsecase.RecordPulseRate(args.ctx, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.RecordPulseRate() error = %v, wantErr %v", err, tt.wantErr)
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

func TestUseCasesClinicalImpl_RecordBloodPressure(t *testing.T) {
	ctx := usecaseMock.AddTenantIdentifierContext(context.Background())
	input := dto.ObservationInput{
		Status:       dto.ObservationStatusFinal,
		EncounterID:  uuid.New().String(),
		Value:        "185.21",
		Note:         "Test note.",
		UsageContext: dto.BreastCancerScreeningTypeEnum,
	}

	type args struct {
		ctx   context.Context
		input dto.ObservationInput
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully record blood pressure",
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

				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to record blood pressure",
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

			got, err := clinicalUsecase.RecordBloodPressure(args.ctx, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.RecordBloodPressure() error = %v, wantErr %v", err, tt.wantErr)
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

func TestUseCasesClinicalImpl_RecordDiastolicBloodPressure(t *testing.T) {
	ctx := usecaseMock.AddTenantIdentifierContext(context.Background())
	input := dto.ObservationInput{
		Status:       dto.ObservationStatusFinal,
		EncounterID:  uuid.New().String(),
		Value:        "185.21",
		Note:         "Test note.",
		UsageContext: dto.BreastCancerScreeningTypeEnum,
	}

	type args struct {
		ctx   context.Context
		input dto.ObservationInput
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully record diastolic blood pressure",
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

				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to record diastolic blood pressure",
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

			got, err := clinicalUsecase.RecordDiastolicBloodPressure(args.ctx, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.RecordDiastolicBloodPressure() error = %v, wantErr %v", err, tt.wantErr)
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

func TestUseCasesClinicalImpl_RecordColposcopy(t *testing.T) {
	ctx := usecaseMock.AddTenantIdentifierContext(context.Background())
	input := dto.ObservationInput{
		Status:      dto.ObservationStatusFinal,
		EncounterID: uuid.New().String(),
		Value:       "Normal",
		Note:        "Test note.",
	}

	type args struct {
		ctx   context.Context
		input dto.ObservationInput
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully record colposcopy",
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

				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to record colposcopy",
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

			got, err := clinicalUsecase.RecordColposcopy(args.ctx, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.RecordColposcopy() error = %v, wantErr %v", err, tt.wantErr)
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

func TestUseCasesClinicalImpl_RecordBMI(t *testing.T) {
	ctx := usecaseMock.AddTenantIdentifierContext(context.Background())
	input := dto.ObservationInput{
		Status:      dto.ObservationStatusFinal,
		EncounterID: uuid.New().String(),
		Value:       "19.7",
		Note:        "Test note.",
	}

	type args struct {
		ctx   context.Context
		input dto.ObservationInput
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully record BMI",
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

				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to record BMI",
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

			got, err := clinicalUsecase.RecordBMI(args.ctx, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.RecordBMI() error = %v, wantErr %v", err, tt.wantErr)
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

func TestUseCasesClinicalImpl_RecordViralLoad(t *testing.T) {
	ctx := usecaseMock.AddTenantIdentifierContext(context.Background())
	input := dto.ObservationInput{
		Status:      dto.ObservationStatusFinal,
		EncounterID: uuid.New().String(),
		Value:       "110.7",
		Note:        "Test note.",
	}

	type args struct {
		ctx   context.Context
		input dto.ObservationInput
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully record viral load",
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

				return args{ctx: ctx, input: input}
			},
			wantErr: false,
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

			_, err := clinicalUsecase.RecordViralLoad(args.ctx, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.RecordViralLoad() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCasesClinicalImpl_RecordBloodSugar(t *testing.T) {
	ctx := usecaseMock.AddTenantIdentifierContext(context.Background())
	input := dto.ObservationInput{
		Status:      dto.ObservationStatusFinal,
		EncounterID: uuid.New().String(),
		Value:       "11.7",
		Note:        "Test note.",
	}

	type args struct {
		ctx   context.Context
		input dto.ObservationInput
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully record blood sugar",
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

				return args{ctx: ctx, input: input}
			},
			wantErr: false,
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

			got, err := clinicalUsecase.RecordBloodSugar(args.ctx, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.RecordBloodSugar() error = %v, wantErr %v", err, tt.wantErr)
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

func TestUseCasesClinicalImpl_RecordLastMenstrualPeriod(t *testing.T) {
	ctx := usecaseMock.AddTenantIdentifierContext(context.Background())
	input := dto.ObservationInput{
		Status:      dto.ObservationStatusFinal,
		EncounterID: uuid.New().String(),
		Value:       "11.7",
		Note:        "Test note.",
	}

	type args struct {
		ctx   context.Context
		input dto.ObservationInput
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully record last menstrual period",
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

				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to record last menstrual period",
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

			got, err := clinicalUsecase.RecordLastMenstrualPeriod(args.ctx, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.RecordLastMenstrualPeriod() error = %v, wantErr %v", err, tt.wantErr)
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

func TestUseCasesClinicalImpl_RecordVIA(t *testing.T) {
	ctx := usecaseMock.AddTenantIdentifierContext(context.Background())

	type args struct {
		ctx   context.Context
		input dto.ObservationInput
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully record positive VIA",
			setup: func(mh *usecaseMock.Mocks) args {
				input := dto.ObservationInput{
					Status:      dto.ObservationStatusFinal,
					EncounterID: uuid.New().String(),
					Value:       "positive",
					Note:        "Test note.",
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

				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "Happy Case - Successfully record negative VIA",
			setup: func(mh *usecaseMock.Mocks) args {
				input := dto.ObservationInput{
					Status:      dto.ObservationStatusFinal,
					EncounterID: uuid.New().String(),
					Value:       "negative",
					Note:        "Test note.",
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

				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "Happy Case - Successfully record suspicious VIA",
			setup: func(mh *usecaseMock.Mocks) args {
				input := dto.ObservationInput{
					Status:      dto.ObservationStatusFinal,
					EncounterID: uuid.New().String(),
					Value:       "suspicious_for_cancer",
					Note:        "Test note.",
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

				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Invalid VIA value",
			setup: func(mh *usecaseMock.Mocks) args {
				input := dto.ObservationInput{
					Status:      dto.ObservationStatusFinal,
					EncounterID: uuid.New().String(),
					Value:       mock.Anything,
					Note:        "Test note.",
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

			got, err := clinicalUsecase.RecordVIA(args.ctx, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.RecordVIA() error = %v, wantErr %v", err, tt.wantErr)
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

func TestUseCasesClinicalImpl_RecordImmunoHistoChemistry(t *testing.T) {
	ctx := usecaseMock.AddTenantIdentifierContext(context.Background())
	input := dto.ObservationInput{
		Status:      dto.ObservationStatusFinal,
		EncounterID: uuid.New().String(),
		Value:       "110.7",
		Note:        "Test note.",
	}

	type args struct {
		ctx   context.Context
		input dto.ObservationInput
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully record immuno histo chemistry",
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

				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to record immuno histo chemistry",
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

			got, err := clinicalUsecase.RecordImmunoHistoChemistry(args.ctx, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.RecordImmunoHistoChemistry() error = %v, wantErr %v", err, tt.wantErr)
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

func TestUseCasesClinicalImpl_RecordPostCoitalBleeding(t *testing.T) {
	ctx := usecaseMock.AddTenantIdentifierContext(context.Background())
	input := dto.ObservationInput{
		Status:      dto.ObservationStatusFinal,
		EncounterID: uuid.New().String(),
		Value:       "110.7",
		Note:        "Test note.",
	}

	type args struct {
		ctx   context.Context
		input dto.ObservationInput
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully record pcb",
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

				return args{ctx: ctx, input: input}
			},
			wantErr: false,
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

			got, err := clinicalUsecase.RecordPostCoitalBleeding(args.ctx, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.RecordPostCoitalBleeding() error = %v, wantErr %v", err, tt.wantErr)
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

func TestUseCasesClinical_RecoredHistoryOfPresentIlness(t *testing.T) {
	ctx := usecaseMock.AddTenantIdentifierContext(context.Background())
	input := dto.ObservationInput{
		Status:      dto.ObservationStatusFinal,
		EncounterID: uuid.New().String(),
		Note:        "Test note.",
	}

	type args struct {
		ctx   context.Context
		input dto.ObservationInput
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy case - Successfully records history of patient illness",
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

				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "Sad case - Fail to record history of patient illness",
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

			got, err := clinicalUsecase.RecordHistoryOfPresentIllness(args.ctx, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.RecordHistoryOfPresentIllness() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if got == nil {
					t.Errorf("expected a response but go %v", got)
				}
			}
		})
	}
}
func TestUseCasesClinical_RecordPastMedicalAndSurgicalHistory(t *testing.T) {
	ctx := usecaseMock.AddTenantIdentifierContext(context.Background())
	input := dto.ObservationInput{
		Status:      dto.ObservationStatusFinal,
		EncounterID: uuid.New().String(),
		Note:        "Test note.",
	}

	type args struct {
		ctx   context.Context
		input dto.ObservationInput
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy case - Successfully records past medical and surgical history of a patient",
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

				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "Sad case - Fail to record past medical and surgical history of a patient",
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

			got, err := clinicalUsecase.RecordPastMedicalAndSurgicalHistory(args.ctx, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.RecordPastMedicalAndSurgicalHistory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if got == nil {
					t.Errorf("expected a response but go %v", got)
				}
			}
		})
	}
}

func TestUseCasesClinical_RecordFamilyAndSocialHistory(t *testing.T) {
	ctx := usecaseMock.AddTenantIdentifierContext(context.Background())
	input := dto.ObservationInput{
		Status:      dto.ObservationStatusFinal,
		EncounterID: uuid.New().String(),
		Note:        "Test note.",
	}

	type args struct {
		ctx   context.Context
		input dto.ObservationInput
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy case - Successfully records family and social history of a patient",
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

				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "Sad case - Fail to record family and social history of a patient",
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

			got, err := clinicalUsecase.RecordFamilyAndSocialHistory(args.ctx, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.RecordFamilyAndSocialHistory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if got == nil {
					t.Errorf("expected a response but go %v", got)
				}
			}
		})
	}
}

func TestUseCasesClinicalImpl_RecordHPV(t *testing.T) {
	ctx := usecaseMock.AddTenantIdentifierContext(context.Background())

	type args struct {
		ctx   context.Context
		input dto.ObservationInput
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully record HPV (Oncoprotein)",
			setup: func(mh *usecaseMock.Mocks) args {
				input := dto.ObservationInput{
					Status:             dto.ObservationStatusFinal,
					EncounterID:        uuid.New().String(),
					ObservationSubType: dto.HPV_ONCOPROTEIN,
					Value:              "positive",
					Note:               "This is a note",
					UsageContext:       dto.CervicalCancerScreeningTypeEnum,
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

				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "Happy Case - Successfully record HPV (PCT DNA)",
			setup: func(mh *usecaseMock.Mocks) args {
				input := dto.ObservationInput{
					Status:             dto.ObservationStatusFinal,
					EncounterID:        uuid.New().String(),
					ObservationSubType: dto.HPV_PCR_DNA,
					Value:              "positive",
					Note:               "This is a note",
					UsageContext:       dto.CervicalCancerScreeningTypeEnum,
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

				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "Happy Case - Successfully record HPV (Unspecified)",
			setup: func(mh *usecaseMock.Mocks) args {
				input := dto.ObservationInput{
					Status:       dto.ObservationStatusFinal,
					EncounterID:  uuid.New().String(),
					Value:        "positive",
					Note:         "This is a note",
					UsageContext: dto.CervicalCancerScreeningTypeEnum,
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

				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			got, err := clinicalUsecase.RecordHPV(args.ctx, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.RecordHPV() error = %v, wantErr %v", err, tt.wantErr)
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

func TestUseCasesClinicalImpl_RecordObservationV2(t *testing.T) {
	ctx := usecaseMock.AddTenantIdentifierContext(context.Background())
	input := dto.ObservationInput{
		Status:             dto.ObservationStatusFinal,
		UsageContext:       dto.BreastCancerScreeningTypeEnum,
		EncounterID:        gofakeit.UUID(),
		ObservationSubType: dto.HPV_PCR_DNA,
		Value:              "23",
		Note:               "Normal",
		Concept:            dto.ObservationConceptEnumBMI,
	}

	type args struct {
		ctx   context.Context
		input dto.ObservationInput
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy case: Successfully recorded observation",
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

				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "Sad case: Fail to record observation",
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
			name: "Sad case: Fail validation - missing concept",
			setup: func(mh *usecaseMock.Mocks) args {
				input := dto.ObservationInput{
					Status:             dto.ObservationStatusFinal,
					UsageContext:       dto.BreastCancerScreeningTypeEnum,
					EncounterID:        gofakeit.UUID(),
					ObservationSubType: dto.HPV_PCR_DNA,
					Value:              "23",
					Note:               "Normal",
				}

				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			_, err := c.RecordObservationV2(args.ctx, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.RecordObservationHelper() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
func TestUseCasesClinicalImpl_GetPatientObservations(t *testing.T) {
	ctx := usecaseMock.AddTenantIdentifierContext(context.Background())

	first := 10
	encounterId := uuid.New().String()
	exam := dto.ObservationCategoryExam

	status := dto.ObservationStatusFinal
	valueConcept := "222"
	UUID := gofakeit.UUID()
	pagedFHIRObservations := &domain.PagedFHIRObservations{
		Observations: []domain.FHIRObservation{
			{
				ID:     new(string),
				Status: (*domain.ObservationStatusEnum)(&status),
				Code: &domain.FHIRCodeableConcept{
					ID: new(string),
					Coding: []*domain.FHIRCoding{{
						Display: gofakeit.BS(),
					}},
				},
				Subject: &domain.FHIRReference{
					ID: new(string),
				},
				Encounter: &domain.FHIRReference{
					ID: new(string),
				},
				ValueQuantity: &domain.FHIRQuantity{
					Value: 100,
					Unit:  "cm",
				},
				ValueCodeableConcept: &domain.FHIRCodeableConcept{
					Text: valueConcept,
				},
				ValueString:  new(string),
				ValueBoolean: new(bool),
				ValueInteger: new(string),
				ValueRange: &domain.FHIRRange{
					Low: domain.FHIRQuantity{
						Value: 100,
						Unit:  "cm",
					},
					High: domain.FHIRQuantity{
						Value: 100,
						Unit:  "cm",
					},
				},
				ValueRatio: &domain.FHIRRatio{
					Numerator: domain.FHIRQuantity{
						Value: 100,
						Unit:  "cm",
					},
					Denominator: domain.FHIRQuantity{
						Value: 100,
						Unit:  "cm",
					},
				},
				ValueSampledData: &domain.FHIRSampledData{
					ID: &UUID,
				},
				ValueTime: &time.Time{},
				ValueDateTime: &scalarutils.Date{
					Year:  2000,
					Month: 1,
					Day:   1,
				},
				ValuePeriod: &domain.FHIRPeriod{
					Start: scalarutils.DateTime(time.Wednesday.String()),
					End:   scalarutils.DateTime(time.Thursday.String()),
				},
			},
		},
		HasNextPage:     false,
		NextPageURL:     "",
		HasPreviousPage: false,
		PreviousPageURL: "",
		TotalCount:      0,
	}

	type args struct {
		ctx     context.Context
		payload *dto.FetchObservationPayload
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy Case - successfully get patient observations",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchPatientObservations(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID string, searchParameters map[string]interface{}, tenant dto.TenantIdentifiers, pagination serverutils.PaginationInput) (*domain.PagedFHIRObservations, error) {
						return pagedFHIRObservations, nil
					})

				return args{
					ctx: ctx,
					payload: &dto.FetchObservationPayload{
						PatientID:       gofakeit.UUID(),
						EncounterID:     &encounterId,
						ObservationCode: "1234",
						Category:        &exam,
						Pagination:      &dto.Pagination{First: &first},
					},
				}
			},
			wantErr: false,
		},
		{
			name: "Happy Case - successfully get patient observations - searchID",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchPatientObservations(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID string, searchParameters map[string]interface{}, tenant dto.TenantIdentifiers, pagination serverutils.PaginationInput) (*domain.PagedFHIRObservations, error) {
						return pagedFHIRObservations, nil
					})

				return args{
					ctx: ctx,
					payload: &dto.FetchObservationPayload{
						PatientID:       gofakeit.UUID(),
						EncounterID:     &encounterId,
						ObservationCode: "1234",
						Category:        &exam,
						PaginationV2:    &serverutils.PaginationInput{First: &first},
						SearchID:        gofakeit.UUID(),
						Usage:           "SCREENING_EXAMINATIONS",
					},
				}
			},
			wantErr: false,
		},
		{
			name: "Happy Case - successfully get patient observations with vital signs category",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchPatientObservations(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID string, searchParameters map[string]interface{}, tenant dto.TenantIdentifiers, pagination serverutils.PaginationInput) (*domain.PagedFHIRObservations, error) {
						return pagedFHIRObservations, nil
					})

				return args{
					ctx: ctx,
					payload: &dto.FetchObservationPayload{
						PatientID:       gofakeit.UUID(),
						EncounterID:     &encounterId,
						ObservationCode: "1234",
						Category:        &exam,
						Pagination:      &dto.Pagination{First: &first},
					},
				}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Invalid patient ID",
			setup: func(mh *usecaseMock.Mocks) args {
				return args{
					ctx: ctx,
					payload: &dto.FetchObservationPayload{
						PatientID:       "invalid",
						EncounterID:     &encounterId,
						ObservationCode: "1234",
						Category:        &exam,
						Pagination:      &dto.Pagination{First: &first},
					},
				}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - fail to get tenant identifiers",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{
					ctx: ctx,
					payload: &dto.FetchObservationPayload{
						PatientID:       gofakeit.UUID(),
						EncounterID:     &encounterId,
						ObservationCode: "1234",
						Category:        &exam,
						Pagination:      &dto.Pagination{First: &first},
					},
				}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - fail to search patient observations",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchPatientObservations(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID string, searchParameters map[string]interface{}, tenant dto.TenantIdentifiers, pagination serverutils.PaginationInput) (*domain.PagedFHIRObservations, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{
					ctx: ctx,
					payload: &dto.FetchObservationPayload{
						PatientID:       gofakeit.UUID(),
						EncounterID:     &encounterId,
						ObservationCode: "1234",
						Category:        &exam,
						Pagination:      &dto.Pagination{First: &first},
					},
				}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to search observation - nil subject",
			setup: func(mh *usecaseMock.Mocks) args {
				status := dto.ObservationStatusFinal
				valueConcept := "222"
				UUID := gofakeit.UUID()
				pagedFHIRObservations := &domain.PagedFHIRObservations{
					Observations: []domain.FHIRObservation{
						{
							ID:     &UUID,
							Status: (*domain.ObservationStatusEnum)(&status),
							Code: &domain.FHIRCodeableConcept{
								ID: new(string),
								Coding: []*domain.FHIRCoding{{
									Display: gofakeit.BS(),
								}},
							},
							Encounter: &domain.FHIRReference{
								ID: &UUID,
							},
							ValueQuantity: &domain.FHIRQuantity{
								Value: 100,
								Unit:  "cm",
							},
							ValueCodeableConcept: &domain.FHIRCodeableConcept{
								Text: valueConcept,
							},
							ValueString:  new(string),
							ValueBoolean: new(bool),
							ValueInteger: new(string),
							ValueRange: &domain.FHIRRange{
								Low: domain.FHIRQuantity{
									Value: 100,
									Unit:  "cm",
								},
								High: domain.FHIRQuantity{
									Value: 100,
									Unit:  "cm",
								},
							},
							ValueRatio: &domain.FHIRRatio{
								Numerator: domain.FHIRQuantity{
									Value: 100,
									Unit:  "cm",
								},
								Denominator: domain.FHIRQuantity{
									Value: 100,
									Unit:  "cm",
								},
							},
							ValueSampledData: &domain.FHIRSampledData{
								ID: &UUID,
							},
							ValueTime: &time.Time{},
							ValueDateTime: &scalarutils.Date{
								Year:  2000,
								Month: 1,
								Day:   1,
							},
							ValuePeriod: &domain.FHIRPeriod{
								Start: scalarutils.DateTime(time.Wednesday.String()),
								End:   scalarutils.DateTime(time.Thursday.String()),
							},
						},
					},
					HasNextPage:     false,
					NextPageURL:     "",
					HasPreviousPage: false,
					PreviousPageURL: "",
					TotalCount:      0,
				}

				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchPatientObservations(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID string, searchParameters map[string]interface{}, tenant dto.TenantIdentifiers, pagination serverutils.PaginationInput) (*domain.PagedFHIRObservations, error) {
						return pagedFHIRObservations, nil
					})

				return args{
					ctx: ctx,
					payload: &dto.FetchObservationPayload{
						PatientID:       gofakeit.UUID(),
						EncounterID:     &encounterId,
						ObservationCode: "1234",
						Category:        &exam,
						Pagination:      &dto.Pagination{First: &first},
					},
				}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to search observation - nil subject id",
			setup: func(mh *usecaseMock.Mocks) args {
				status := dto.ObservationStatusFinal
				valueConcept := "222"
				UUID := gofakeit.UUID()
				pagedFHIRObservations := &domain.PagedFHIRObservations{
					Observations: []domain.FHIRObservation{
						{
							ID:     new(string),
							Status: (*domain.ObservationStatusEnum)(&status),
							Code: &domain.FHIRCodeableConcept{
								ID: new(string),
								Coding: []*domain.FHIRCoding{{
									Display: gofakeit.BS(),
								}},
							},
							Subject: &domain.FHIRReference{},
							Encounter: &domain.FHIRReference{
								ID: new(string),
							},
							ValueQuantity: &domain.FHIRQuantity{
								Value: 100,
								Unit:  "cm",
							},
							ValueCodeableConcept: &domain.FHIRCodeableConcept{
								Text: valueConcept,
							},
							ValueString:  new(string),
							ValueBoolean: new(bool),
							ValueInteger: new(string),
							ValueRange: &domain.FHIRRange{
								Low: domain.FHIRQuantity{
									Value: 100,
									Unit:  "cm",
								},
								High: domain.FHIRQuantity{
									Value: 100,
									Unit:  "cm",
								},
							},
							ValueRatio: &domain.FHIRRatio{
								Numerator: domain.FHIRQuantity{
									Value: 100,
									Unit:  "cm",
								},
								Denominator: domain.FHIRQuantity{
									Value: 100,
									Unit:  "cm",
								},
							},
							ValueSampledData: &domain.FHIRSampledData{
								ID: &UUID,
							},
							ValueTime: &time.Time{},
							ValueDateTime: &scalarutils.Date{
								Year:  2000,
								Month: 1,
								Day:   1,
							},
							ValuePeriod: &domain.FHIRPeriod{
								Start: scalarutils.DateTime(time.Wednesday.String()),
								End:   scalarutils.DateTime(time.Thursday.String()),
							},
						},
					},
					HasNextPage:     false,
					NextPageURL:     "",
					HasPreviousPage: false,
					PreviousPageURL: "",
					TotalCount:      0,
				}

				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchPatientObservations(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID string, searchParameters map[string]interface{}, tenant dto.TenantIdentifiers, pagination serverutils.PaginationInput) (*domain.PagedFHIRObservations, error) {
						return pagedFHIRObservations, nil
					})

				return args{
					ctx: ctx,
					payload: &dto.FetchObservationPayload{
						PatientID:       gofakeit.UUID(),
						EncounterID:     &encounterId,
						ObservationCode: "1234",
						Category:        &exam,
						Pagination:      &dto.Pagination{First: &first},
					},
				}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to search observation - nil encounter",
			setup: func(mh *usecaseMock.Mocks) args {
				status := dto.ObservationStatusFinal
				instant := gofakeit.TimeZone()
				valueConcept := "222"
				UUID := gofakeit.UUID()
				currentTime := scalarutils.DateTime(time.Now().Format(time.RFC3339))
				pagedFHIRObservations := &domain.PagedFHIRObservations{
					Observations: []domain.FHIRObservation{
						{
							ID:     new(string),
							Status: (*domain.ObservationStatusEnum)(&status),
							Code: &domain.FHIRCodeableConcept{
								ID: new(string),
								Coding: []*domain.FHIRCoding{
									{
										ID:           new(string),
										Version:      new(string),
										Code:         (*scalarutils.Code)(&valueConcept),
										Display:      "Vital",
										UserSelected: new(bool),
									},
								},
								Text: "",
							},
							Subject: &domain.FHIRReference{
								ID: &UUID,
							},
							EffectiveDateTime: &currentTime,
							ValueQuantity: &domain.FHIRQuantity{
								Value: 100,
								Unit:  "cm",
							},
							ValueCodeableConcept: &domain.FHIRCodeableConcept{
								Text: valueConcept,
							},
							ValueString:      new(string),
							ValueBoolean:     new(bool),
							ValueInteger:     new(string),
							EffectiveInstant: (*scalarutils.Instant)(&instant),
							ValueRange: &domain.FHIRRange{
								Low: domain.FHIRQuantity{
									Value: 100,
									Unit:  "cm",
								},
								High: domain.FHIRQuantity{
									Value: 100,
									Unit:  "cm",
								},
							},
							ValueRatio: &domain.FHIRRatio{
								Numerator: domain.FHIRQuantity{
									Value: 100,
									Unit:  "cm",
								},
								Denominator: domain.FHIRQuantity{
									Value: 100,
									Unit:  "cm",
								},
							},
							ValueSampledData: &domain.FHIRSampledData{
								ID: &UUID,
							},
							ValueTime: &time.Time{},
							ValueDateTime: &scalarutils.Date{
								Year:  2000,
								Month: 1,
								Day:   1,
							},
							ValuePeriod: &domain.FHIRPeriod{
								Start: scalarutils.DateTime(time.Wednesday.String()),
								End:   scalarutils.DateTime(time.Thursday.String()),
							},
						},
					},
					HasNextPage:     false,
					NextPageURL:     "",
					HasPreviousPage: false,
					PreviousPageURL: "",
					TotalCount:      0,
				}

				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchPatientObservations(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID string, searchParameters map[string]interface{}, tenant dto.TenantIdentifiers, pagination serverutils.PaginationInput) (*domain.PagedFHIRObservations, error) {
						return pagedFHIRObservations, nil
					})

				return args{
					ctx: ctx,
					payload: &dto.FetchObservationPayload{
						PatientID:       gofakeit.UUID(),
						EncounterID:     &encounterId,
						ObservationCode: "1234",
						Category:        &exam,
						Pagination:      &dto.Pagination{First: &first},
					},
				}
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			got, err := clinicalUsecase.GetPatientObservations(args.ctx, args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.GetPatientObservations() error = %v, wantErr %v", err, tt.wantErr)
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

func TestUseCasesClinicalImpl_GetPatientDiastolicBloodPressureEntries(t *testing.T) {
	ctx := usecaseMock.AddTenantIdentifierContext(context.Background())
	date := &scalarutils.Date{
		Year:  1997,
		Month: 12,
		Day:   12,
	}

	status := dto.ObservationStatusFinal
	valueConcept := "222"
	UUID := gofakeit.UUID()
	pagedFHIRObservations := &domain.PagedFHIRObservations{
		Observations: []domain.FHIRObservation{
			{
				ID:     new(string),
				Status: (*domain.ObservationStatusEnum)(&status),
				Code: &domain.FHIRCodeableConcept{
					ID: new(string),
					Coding: []*domain.FHIRCoding{{
						Display: gofakeit.BS(),
					}},
				},
				Subject: &domain.FHIRReference{
					ID: new(string),
				},
				Encounter: &domain.FHIRReference{
					ID: new(string),
				},
				ValueQuantity: &domain.FHIRQuantity{
					Value: 100,
					Unit:  "cm",
				},
				ValueCodeableConcept: &domain.FHIRCodeableConcept{
					Text: valueConcept,
				},
				ValueString:  new(string),
				ValueBoolean: new(bool),
				ValueInteger: new(string),
				ValueRange: &domain.FHIRRange{
					Low: domain.FHIRQuantity{
						Value: 100,
						Unit:  "cm",
					},
					High: domain.FHIRQuantity{
						Value: 100,
						Unit:  "cm",
					},
				},
				ValueRatio: &domain.FHIRRatio{
					Numerator: domain.FHIRQuantity{
						Value: 100,
						Unit:  "cm",
					},
					Denominator: domain.FHIRQuantity{
						Value: 100,
						Unit:  "cm",
					},
				},
				ValueSampledData: &domain.FHIRSampledData{
					ID: &UUID,
				},
				ValueTime: &time.Time{},
				ValueDateTime: &scalarutils.Date{
					Year:  2000,
					Month: 1,
					Day:   1,
				},
				ValuePeriod: &domain.FHIRPeriod{
					Start: scalarutils.DateTime(time.Wednesday.String()),
					End:   scalarutils.DateTime(time.Thursday.String()),
				},
			},
		},
		HasNextPage:     false,
		NextPageURL:     "",
		HasPreviousPage: false,
		PreviousPageURL: "",
		TotalCount:      0,
	}

	first := 10
	encounterId := uuid.New().String()
	type args struct {
		ctx         context.Context
		patientID   string
		encounterID *string
		date        *scalarutils.Date
		pagination  *dto.Pagination
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully get patient diastolic blood pressure entries",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchPatientObservations(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID string, searchParameters map[string]interface{}, tenant dto.TenantIdentifiers, pagination serverutils.PaginationInput) (*domain.PagedFHIRObservations, error) {
						return pagedFHIRObservations, nil
					})

				return args{ctx: ctx, patientID: uuid.NewString(), encounterID: &encounterId, date: date, pagination: &dto.Pagination{First: &first}}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - fail to get patient observations",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchPatientObservations(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID string, searchParameters map[string]interface{}, tenant dto.TenantIdentifiers, pagination serverutils.PaginationInput) (*domain.PagedFHIRObservations, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, patientID: uuid.NewString(), encounterID: &encounterId, date: date, pagination: &dto.Pagination{First: &first}}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			got, err := clinicalUsecase.GetPatientDiastolicBloodPressureEntries(args.ctx, args.patientID, args.encounterID, args.date, args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.GetPatientDiastolicBloodPressureEntries() error = %v, wantErr %v", err, tt.wantErr)
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

func TestUseCasesClinicalImpl_GetPatientTemperatureEntries(t *testing.T) {
	ctx := usecaseMock.AddTenantIdentifierContext(context.Background())
	date := &scalarutils.Date{
		Year:  1997,
		Month: 12,
		Day:   12,
	}

	status := dto.ObservationStatusFinal
	valueConcept := "222"
	UUID := gofakeit.UUID()
	pagedFHIRObservations := &domain.PagedFHIRObservations{
		Observations: []domain.FHIRObservation{
			{
				ID:     new(string),
				Status: (*domain.ObservationStatusEnum)(&status),
				Code: &domain.FHIRCodeableConcept{
					ID: new(string),
					Coding: []*domain.FHIRCoding{{
						Display: gofakeit.BS(),
					}},
				},
				Subject: &domain.FHIRReference{
					ID: new(string),
				},
				Encounter: &domain.FHIRReference{
					ID: new(string),
				},
				ValueQuantity: &domain.FHIRQuantity{
					Value: 100,
					Unit:  "cm",
				},
				ValueCodeableConcept: &domain.FHIRCodeableConcept{
					Text: valueConcept,
				},
				ValueString:  new(string),
				ValueBoolean: new(bool),
				ValueInteger: new(string),
				ValueRange: &domain.FHIRRange{
					Low: domain.FHIRQuantity{
						Value: 100,
						Unit:  "cm",
					},
					High: domain.FHIRQuantity{
						Value: 100,
						Unit:  "cm",
					},
				},
				ValueRatio: &domain.FHIRRatio{
					Numerator: domain.FHIRQuantity{
						Value: 100,
						Unit:  "cm",
					},
					Denominator: domain.FHIRQuantity{
						Value: 100,
						Unit:  "cm",
					},
				},
				ValueSampledData: &domain.FHIRSampledData{
					ID: &UUID,
				},
				ValueTime: &time.Time{},
				ValueDateTime: &scalarutils.Date{
					Year:  2000,
					Month: 1,
					Day:   1,
				},
				ValuePeriod: &domain.FHIRPeriod{
					Start: scalarutils.DateTime(time.Wednesday.String()),
					End:   scalarutils.DateTime(time.Thursday.String()),
				},
			},
		},
		HasNextPage:     false,
		NextPageURL:     "",
		HasPreviousPage: false,
		PreviousPageURL: "",
		TotalCount:      0,
	}

	first := 10
	encounterId := uuid.New().String()

	type args struct {
		ctx         context.Context
		patientID   string
		encounterID *string
		date        *scalarutils.Date
		pagination  *dto.Pagination
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully get patient temperature entries",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchPatientObservations(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID string, searchParameters map[string]interface{}, tenant dto.TenantIdentifiers, pagination serverutils.PaginationInput) (*domain.PagedFHIRObservations, error) {
						return pagedFHIRObservations, nil
					})

				return args{ctx: ctx, patientID: uuid.NewString(), encounterID: &encounterId, date: date, pagination: &dto.Pagination{First: &first}}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - fail to get patient observations",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchPatientObservations(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID string, searchParameters map[string]interface{}, tenant dto.TenantIdentifiers, pagination serverutils.PaginationInput) (*domain.PagedFHIRObservations, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, patientID: uuid.NewString(), encounterID: &encounterId, date: date, pagination: &dto.Pagination{First: &first}}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			got, err := clinicalUsecase.GetPatientTemperatureEntries(args.ctx, args.patientID, args.encounterID, args.date, args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.GetPatientTemperatureEntries() error = %v, wantErr %v", err, tt.wantErr)
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

func TestUseCasesClinicalImpl_GetPatientBloodPressureEntries(t *testing.T) {
	ctx := usecaseMock.AddTenantIdentifierContext(context.Background())
	date := &scalarutils.Date{
		Year:  1997,
		Month: 12,
		Day:   12,
	}

	status := dto.ObservationStatusFinal
	valueConcept := "222"
	UUID := gofakeit.UUID()
	pagedFHIRObservations := &domain.PagedFHIRObservations{
		Observations: []domain.FHIRObservation{
			{
				ID:     new(string),
				Status: (*domain.ObservationStatusEnum)(&status),
				Code: &domain.FHIRCodeableConcept{
					ID: new(string),
					Coding: []*domain.FHIRCoding{{
						Display: gofakeit.BS(),
					}},
				},
				Subject: &domain.FHIRReference{
					ID: new(string),
				},
				Encounter: &domain.FHIRReference{
					ID: new(string),
				},
				ValueQuantity: &domain.FHIRQuantity{
					Value: 100,
					Unit:  "cm",
				},
				ValueCodeableConcept: &domain.FHIRCodeableConcept{
					Text: valueConcept,
				},
				ValueString:  new(string),
				ValueBoolean: new(bool),
				ValueInteger: new(string),
				ValueRange: &domain.FHIRRange{
					Low: domain.FHIRQuantity{
						Value: 100,
						Unit:  "cm",
					},
					High: domain.FHIRQuantity{
						Value: 100,
						Unit:  "cm",
					},
				},
				ValueRatio: &domain.FHIRRatio{
					Numerator: domain.FHIRQuantity{
						Value: 100,
						Unit:  "cm",
					},
					Denominator: domain.FHIRQuantity{
						Value: 100,
						Unit:  "cm",
					},
				},
				ValueSampledData: &domain.FHIRSampledData{
					ID: &UUID,
				},
				ValueTime: &time.Time{},
				ValueDateTime: &scalarutils.Date{
					Year:  2000,
					Month: 1,
					Day:   1,
				},
				ValuePeriod: &domain.FHIRPeriod{
					Start: scalarutils.DateTime(time.Wednesday.String()),
					End:   scalarutils.DateTime(time.Thursday.String()),
				},
			},
		},
		HasNextPage:     false,
		NextPageURL:     "",
		HasPreviousPage: false,
		PreviousPageURL: "",
		TotalCount:      0,
	}

	first := 10
	encounterId := uuid.New().String()
	type args struct {
		ctx         context.Context
		patientID   string
		encounterID *string
		date        *scalarutils.Date
		pagination  *dto.Pagination
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully get patient blood pressure entries",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchPatientObservations(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID string, searchParameters map[string]interface{}, tenant dto.TenantIdentifiers, pagination serverutils.PaginationInput) (*domain.PagedFHIRObservations, error) {
						return pagedFHIRObservations, nil
					})

				return args{ctx: ctx, patientID: uuid.NewString(), encounterID: &encounterId, date: date, pagination: &dto.Pagination{First: &first}}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - fail to get patient observations",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchPatientObservations(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID string, searchParameters map[string]interface{}, tenant dto.TenantIdentifiers, pagination serverutils.PaginationInput) (*domain.PagedFHIRObservations, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, patientID: uuid.NewString(), encounterID: &encounterId, date: date, pagination: &dto.Pagination{First: &first}}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			got, err := clinicalUsecase.GetPatientBloodPressureEntries(args.ctx, args.patientID, args.encounterID, args.date, args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.GetPatientBloodPressureEntries() error = %v, wantErr %v", err, tt.wantErr)
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

func TestUseCasesClinicalImpl_GetHeight(t *testing.T) {
	ctx := usecaseMock.AddTenantIdentifierContext(context.Background())
	date := &scalarutils.Date{
		Year:  1997,
		Month: 12,
		Day:   12,
	}

	status := dto.ObservationStatusFinal
	valueConcept := "222"
	UUID := gofakeit.UUID()
	pagedFHIRObservations := &domain.PagedFHIRObservations{
		Observations: []domain.FHIRObservation{
			{
				ID:     new(string),
				Status: (*domain.ObservationStatusEnum)(&status),
				Code: &domain.FHIRCodeableConcept{
					ID: new(string),
					Coding: []*domain.FHIRCoding{{
						Display: gofakeit.BS(),
					}},
				},
				Subject: &domain.FHIRReference{
					ID: new(string),
				},
				Encounter: &domain.FHIRReference{
					ID: new(string),
				},
				ValueQuantity: &domain.FHIRQuantity{
					Value: 100,
					Unit:  "cm",
				},
				ValueCodeableConcept: &domain.FHIRCodeableConcept{
					Text: valueConcept,
				},
				ValueString:  new(string),
				ValueBoolean: new(bool),
				ValueInteger: new(string),
				ValueRange: &domain.FHIRRange{
					Low: domain.FHIRQuantity{
						Value: 100,
						Unit:  "cm",
					},
					High: domain.FHIRQuantity{
						Value: 100,
						Unit:  "cm",
					},
				},
				ValueRatio: &domain.FHIRRatio{
					Numerator: domain.FHIRQuantity{
						Value: 100,
						Unit:  "cm",
					},
					Denominator: domain.FHIRQuantity{
						Value: 100,
						Unit:  "cm",
					},
				},
				ValueSampledData: &domain.FHIRSampledData{
					ID: &UUID,
				},
				ValueTime: &time.Time{},
				ValueDateTime: &scalarutils.Date{
					Year:  2000,
					Month: 1,
					Day:   1,
				},
				ValuePeriod: &domain.FHIRPeriod{
					Start: scalarutils.DateTime(time.Wednesday.String()),
					End:   scalarutils.DateTime(time.Thursday.String()),
				},
			},
		},
		HasNextPage:     false,
		NextPageURL:     "",
		HasPreviousPage: false,
		PreviousPageURL: "",
		TotalCount:      0,
	}

	first := 10
	encounterId := uuid.New().String()

	type args struct {
		ctx         context.Context
		patientID   string
		encounterID *string
		date        *scalarutils.Date
		pagination  *dto.Pagination
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy case: get patient height",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchPatientObservations(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID string, searchParameters map[string]interface{}, tenant dto.TenantIdentifiers, pagination serverutils.PaginationInput) (*domain.PagedFHIRObservations, error) {
						return pagedFHIRObservations, nil
					})

				return args{ctx: ctx, patientID: uuid.NewString(), encounterID: &encounterId, date: date, pagination: &dto.Pagination{First: &first}}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - fail to get patient observations",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchPatientObservations(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID string, searchParameters map[string]interface{}, tenant dto.TenantIdentifiers, pagination serverutils.PaginationInput) (*domain.PagedFHIRObservations, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, patientID: uuid.NewString(), encounterID: &encounterId, date: date, pagination: &dto.Pagination{First: &first}}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			_, err := clinicalUsecase.GetPatientHeightEntries(args.ctx, args.patientID, args.encounterID, args.date, args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.GetPatientHeightEntries() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCasesClinicalImpl_GetPatientPulseRateEntries(t *testing.T) {
	ctx := usecaseMock.AddTenantIdentifierContext(context.Background())
	date := &scalarutils.Date{
		Year:  1997,
		Month: 12,
		Day:   12,
	}

	status := dto.ObservationStatusFinal
	valueConcept := "222"
	UUID := gofakeit.UUID()
	pagedFHIRObservations := &domain.PagedFHIRObservations{
		Observations: []domain.FHIRObservation{
			{
				ID:     new(string),
				Status: (*domain.ObservationStatusEnum)(&status),
				Code: &domain.FHIRCodeableConcept{
					ID: new(string),
					Coding: []*domain.FHIRCoding{{
						Display: gofakeit.BS(),
					}},
				},
				Subject: &domain.FHIRReference{
					ID: new(string),
				},
				Encounter: &domain.FHIRReference{
					ID: new(string),
				},
				ValueQuantity: &domain.FHIRQuantity{
					Value: 100,
					Unit:  "cm",
				},
				ValueCodeableConcept: &domain.FHIRCodeableConcept{
					Text: valueConcept,
				},
				ValueString:  new(string),
				ValueBoolean: new(bool),
				ValueInteger: new(string),
				ValueRange: &domain.FHIRRange{
					Low: domain.FHIRQuantity{
						Value: 100,
						Unit:  "cm",
					},
					High: domain.FHIRQuantity{
						Value: 100,
						Unit:  "cm",
					},
				},
				ValueRatio: &domain.FHIRRatio{
					Numerator: domain.FHIRQuantity{
						Value: 100,
						Unit:  "cm",
					},
					Denominator: domain.FHIRQuantity{
						Value: 100,
						Unit:  "cm",
					},
				},
				ValueSampledData: &domain.FHIRSampledData{
					ID: &UUID,
				},
				ValueTime: &time.Time{},
				ValueDateTime: &scalarutils.Date{
					Year:  2000,
					Month: 1,
					Day:   1,
				},
				ValuePeriod: &domain.FHIRPeriod{
					Start: scalarutils.DateTime(time.Wednesday.String()),
					End:   scalarutils.DateTime(time.Thursday.String()),
				},
			},
		},
		HasNextPage:     false,
		NextPageURL:     "",
		HasPreviousPage: false,
		PreviousPageURL: "",
		TotalCount:      0,
	}

	first := 10
	encounterId := uuid.New().String()

	type args struct {
		ctx         context.Context
		patientID   string
		encounterID *string
		date        *scalarutils.Date
		pagination  *dto.Pagination
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy case: get patient pulse rate",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchPatientObservations(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID string, searchParameters map[string]interface{}, tenant dto.TenantIdentifiers, pagination serverutils.PaginationInput) (*domain.PagedFHIRObservations, error) {
						return pagedFHIRObservations, nil
					})

				return args{ctx: ctx, patientID: uuid.NewString(), encounterID: &encounterId, date: date, pagination: &dto.Pagination{First: &first}}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - fail to get patient observations",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchPatientObservations(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID string, searchParameters map[string]interface{}, tenant dto.TenantIdentifiers, pagination serverutils.PaginationInput) (*domain.PagedFHIRObservations, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, patientID: uuid.NewString(), encounterID: &encounterId, date: date, pagination: &dto.Pagination{First: &first}}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			_, err := clinicalUsecase.GetPatientPulseRateEntries(args.ctx, args.patientID, args.encounterID, args.date, args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.GetPatientPulseRateEntries() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCasesClinicalImpl_GetPatientRespiratoryRateEntries(t *testing.T) {
	ctx := usecaseMock.AddTenantIdentifierContext(context.Background())
	date := &scalarutils.Date{
		Year:  1997,
		Month: 12,
		Day:   12,
	}

	status := dto.ObservationStatusFinal
	valueConcept := "222"
	UUID := gofakeit.UUID()
	pagedFHIRObservations := &domain.PagedFHIRObservations{
		Observations: []domain.FHIRObservation{
			{
				ID:     new(string),
				Status: (*domain.ObservationStatusEnum)(&status),
				Code: &domain.FHIRCodeableConcept{
					ID: new(string),
					Coding: []*domain.FHIRCoding{{
						Display: gofakeit.BS(),
					}},
				},
				Subject: &domain.FHIRReference{
					ID: new(string),
				},
				Encounter: &domain.FHIRReference{
					ID: new(string),
				},
				ValueQuantity: &domain.FHIRQuantity{
					Value: 100,
					Unit:  "cm",
				},
				ValueCodeableConcept: &domain.FHIRCodeableConcept{
					Text: valueConcept,
				},
				ValueString:  new(string),
				ValueBoolean: new(bool),
				ValueInteger: new(string),
				ValueRange: &domain.FHIRRange{
					Low: domain.FHIRQuantity{
						Value: 100,
						Unit:  "cm",
					},
					High: domain.FHIRQuantity{
						Value: 100,
						Unit:  "cm",
					},
				},
				ValueRatio: &domain.FHIRRatio{
					Numerator: domain.FHIRQuantity{
						Value: 100,
						Unit:  "cm",
					},
					Denominator: domain.FHIRQuantity{
						Value: 100,
						Unit:  "cm",
					},
				},
				ValueSampledData: &domain.FHIRSampledData{
					ID: &UUID,
				},
				ValueTime: &time.Time{},
				ValueDateTime: &scalarutils.Date{
					Year:  2000,
					Month: 1,
					Day:   1,
				},
				ValuePeriod: &domain.FHIRPeriod{
					Start: scalarutils.DateTime(time.Wednesday.String()),
					End:   scalarutils.DateTime(time.Thursday.String()),
				},
			},
		},
		HasNextPage:     false,
		NextPageURL:     "",
		HasPreviousPage: false,
		PreviousPageURL: "",
		TotalCount:      0,
	}

	first := 10
	encounterId := uuid.New().String()

	type args struct {
		ctx         context.Context
		patientID   string
		encounterID *string
		date        *scalarutils.Date
		pagination  *dto.Pagination
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully get patient respiratory rate entries",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchPatientObservations(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID string, searchParameters map[string]interface{}, tenant dto.TenantIdentifiers, pagination serverutils.PaginationInput) (*domain.PagedFHIRObservations, error) {
						return pagedFHIRObservations, nil
					})

				return args{ctx: ctx, patientID: uuid.NewString(), encounterID: &encounterId, date: date, pagination: &dto.Pagination{First: &first}}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - fail to get patient observations",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchPatientObservations(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID string, searchParameters map[string]interface{}, tenant dto.TenantIdentifiers, pagination serverutils.PaginationInput) (*domain.PagedFHIRObservations, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, patientID: uuid.NewString(), encounterID: &encounterId, date: date, pagination: &dto.Pagination{First: &first}}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			got, err := clinicalUsecase.GetPatientRespiratoryRateEntries(args.ctx, args.patientID, args.encounterID, args.date, args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.GetPatientRespiratoryRateEntries() error = %v, wantErr %v", err, tt.wantErr)
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

func TestUseCasesClinicalImpl_GetPatientBMIEntries(t *testing.T) {
	ctx := usecaseMock.AddTenantIdentifierContext(context.Background())
	date := &scalarutils.Date{
		Year:  1997,
		Month: 12,
		Day:   12,
	}

	status := dto.ObservationStatusFinal
	valueConcept := "222"
	UUID := gofakeit.UUID()
	pagedFHIRObservations := &domain.PagedFHIRObservations{
		Observations: []domain.FHIRObservation{
			{
				ID:     new(string),
				Status: (*domain.ObservationStatusEnum)(&status),
				Code: &domain.FHIRCodeableConcept{
					ID: new(string),
					Coding: []*domain.FHIRCoding{{
						Display: gofakeit.BS(),
					}},
				},
				Subject: &domain.FHIRReference{
					ID: new(string),
				},
				Encounter: &domain.FHIRReference{
					ID: new(string),
				},
				ValueQuantity: &domain.FHIRQuantity{
					Value: 100,
					Unit:  "cm",
				},
				ValueCodeableConcept: &domain.FHIRCodeableConcept{
					Text: valueConcept,
				},
				ValueString:  new(string),
				ValueBoolean: new(bool),
				ValueInteger: new(string),
				ValueRange: &domain.FHIRRange{
					Low: domain.FHIRQuantity{
						Value: 100,
						Unit:  "cm",
					},
					High: domain.FHIRQuantity{
						Value: 100,
						Unit:  "cm",
					},
				},
				ValueRatio: &domain.FHIRRatio{
					Numerator: domain.FHIRQuantity{
						Value: 100,
						Unit:  "cm",
					},
					Denominator: domain.FHIRQuantity{
						Value: 100,
						Unit:  "cm",
					},
				},
				ValueSampledData: &domain.FHIRSampledData{
					ID: &UUID,
				},
				ValueTime: &time.Time{},
				ValueDateTime: &scalarutils.Date{
					Year:  2000,
					Month: 1,
					Day:   1,
				},
				ValuePeriod: &domain.FHIRPeriod{
					Start: scalarutils.DateTime(time.Wednesday.String()),
					End:   scalarutils.DateTime(time.Thursday.String()),
				},
			},
		},
		HasNextPage:     false,
		NextPageURL:     "",
		HasPreviousPage: false,
		PreviousPageURL: "",
		TotalCount:      0,
	}

	first := 10
	encounterId := uuid.New().String()

	type args struct {
		ctx         context.Context
		patientID   string
		encounterID *string
		date        *scalarutils.Date
		pagination  *dto.Pagination
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully get patient bmi entries",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchPatientObservations(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID string, searchParameters map[string]interface{}, tenant dto.TenantIdentifiers, pagination serverutils.PaginationInput) (*domain.PagedFHIRObservations, error) {
						return pagedFHIRObservations, nil
					})

				return args{ctx: ctx, patientID: uuid.NewString(), encounterID: &encounterId, date: date, pagination: &dto.Pagination{First: &first}}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - fail to get patient observations",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchPatientObservations(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID string, searchParameters map[string]interface{}, tenant dto.TenantIdentifiers, pagination serverutils.PaginationInput) (*domain.PagedFHIRObservations, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, patientID: uuid.NewString(), encounterID: &encounterId, date: date, pagination: &dto.Pagination{First: &first}}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			got, err := clinicalUsecase.GetPatientBMIEntries(args.ctx, args.patientID, args.encounterID, args.date, args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.GetPatientBMIEntries() error = %v, wantErr %v", err, tt.wantErr)
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

func TestUseCasesClinicalImpl_GetPatientWeightEntries(t *testing.T) {
	ctx := usecaseMock.AddTenantIdentifierContext(context.Background())
	date := &scalarutils.Date{
		Year:  1997,
		Month: 12,
		Day:   12,
	}

	status := dto.ObservationStatusFinal
	valueConcept := "222"
	UUID := gofakeit.UUID()
	pagedFHIRObservations := &domain.PagedFHIRObservations{
		Observations: []domain.FHIRObservation{
			{
				ID:     new(string),
				Status: (*domain.ObservationStatusEnum)(&status),
				Code: &domain.FHIRCodeableConcept{
					ID: new(string),
					Coding: []*domain.FHIRCoding{{
						Display: gofakeit.BS(),
					}},
				},
				Subject: &domain.FHIRReference{
					ID: new(string),
				},
				Encounter: &domain.FHIRReference{
					ID: new(string),
				},
				ValueQuantity: &domain.FHIRQuantity{
					Value: 100,
					Unit:  "cm",
				},
				ValueCodeableConcept: &domain.FHIRCodeableConcept{
					Text: valueConcept,
				},
				ValueString:  new(string),
				ValueBoolean: new(bool),
				ValueInteger: new(string),
				ValueRange: &domain.FHIRRange{
					Low: domain.FHIRQuantity{
						Value: 100,
						Unit:  "cm",
					},
					High: domain.FHIRQuantity{
						Value: 100,
						Unit:  "cm",
					},
				},
				ValueRatio: &domain.FHIRRatio{
					Numerator: domain.FHIRQuantity{
						Value: 100,
						Unit:  "cm",
					},
					Denominator: domain.FHIRQuantity{
						Value: 100,
						Unit:  "cm",
					},
				},
				ValueSampledData: &domain.FHIRSampledData{
					ID: &UUID,
				},
				ValueTime: &time.Time{},
				ValueDateTime: &scalarutils.Date{
					Year:  2000,
					Month: 1,
					Day:   1,
				},
				ValuePeriod: &domain.FHIRPeriod{
					Start: scalarutils.DateTime(time.Wednesday.String()),
					End:   scalarutils.DateTime(time.Thursday.String()),
				},
			},
		},
		HasNextPage:     false,
		NextPageURL:     "",
		HasPreviousPage: false,
		PreviousPageURL: "",
		TotalCount:      0,
	}

	first := 10
	encounterId := uuid.New().String()

	type args struct {
		ctx         context.Context
		patientID   string
		encounterID *string
		date        *scalarutils.Date
		pagination  *dto.Pagination
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy case: get patient pulse rate",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchPatientObservations(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID string, searchParameters map[string]interface{}, tenant dto.TenantIdentifiers, pagination serverutils.PaginationInput) (*domain.PagedFHIRObservations, error) {
						return pagedFHIRObservations, nil
					})

				return args{ctx: ctx, patientID: uuid.NewString(), encounterID: &encounterId, date: date, pagination: &dto.Pagination{First: &first}}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - fail to get patient observations",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchPatientObservations(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID string, searchParameters map[string]interface{}, tenant dto.TenantIdentifiers, pagination serverutils.PaginationInput) (*domain.PagedFHIRObservations, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, patientID: uuid.NewString(), encounterID: &encounterId, date: date, pagination: &dto.Pagination{First: &first}}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		clinicalUsecase, mock := usecaseMock.SetupMocks(t)
		args := tt.setup(&mock)

		_, err := clinicalUsecase.GetPatientWeightEntries(args.ctx, args.patientID, args.encounterID, args.date, args.pagination)
		if (err != nil) != tt.wantErr {
			t.Errorf("UseCasesClinicalImpl.PatchPatientWeight() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
	}
}

func TestUseCasesClinicalImpl_GetPatientMuacEntries(t *testing.T) {
	ctx := usecaseMock.AddTenantIdentifierContext(context.Background())
	date := &scalarutils.Date{
		Year:  1997,
		Month: 12,
		Day:   12,
	}

	status := dto.ObservationStatusFinal
	valueConcept := "222"
	UUID := gofakeit.UUID()
	pagedFHIRObservations := &domain.PagedFHIRObservations{
		Observations: []domain.FHIRObservation{
			{
				ID:     new(string),
				Status: (*domain.ObservationStatusEnum)(&status),
				Code: &domain.FHIRCodeableConcept{
					ID: new(string),
					Coding: []*domain.FHIRCoding{{
						Display: gofakeit.BS(),
					}},
				},
				Subject: &domain.FHIRReference{
					ID: new(string),
				},
				Encounter: &domain.FHIRReference{
					ID: new(string),
				},
				ValueQuantity: &domain.FHIRQuantity{
					Value: 100,
					Unit:  "cm",
				},
				ValueCodeableConcept: &domain.FHIRCodeableConcept{
					Text: valueConcept,
				},
				ValueString:  new(string),
				ValueBoolean: new(bool),
				ValueInteger: new(string),
				ValueRange: &domain.FHIRRange{
					Low: domain.FHIRQuantity{
						Value: 100,
						Unit:  "cm",
					},
					High: domain.FHIRQuantity{
						Value: 100,
						Unit:  "cm",
					},
				},
				ValueRatio: &domain.FHIRRatio{
					Numerator: domain.FHIRQuantity{
						Value: 100,
						Unit:  "cm",
					},
					Denominator: domain.FHIRQuantity{
						Value: 100,
						Unit:  "cm",
					},
				},
				ValueSampledData: &domain.FHIRSampledData{
					ID: &UUID,
				},
				ValueTime: &time.Time{},
				ValueDateTime: &scalarutils.Date{
					Year:  2000,
					Month: 1,
					Day:   1,
				},
				ValuePeriod: &domain.FHIRPeriod{
					Start: scalarutils.DateTime(time.Wednesday.String()),
					End:   scalarutils.DateTime(time.Thursday.String()),
				},
			},
		},
		HasNextPage:     false,
		NextPageURL:     "",
		HasPreviousPage: false,
		PreviousPageURL: "",
		TotalCount:      0,
	}

	first := 10
	encounterId := uuid.New().String()

	type args struct {
		ctx         context.Context
		patientID   string
		encounterID *string
		date        *scalarutils.Date
		pagination  *dto.Pagination
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy case: get patient muac",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchPatientObservations(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID string, searchParameters map[string]interface{}, tenant dto.TenantIdentifiers, pagination serverutils.PaginationInput) (*domain.PagedFHIRObservations, error) {
						return pagedFHIRObservations, nil
					})

				return args{ctx: ctx, patientID: uuid.NewString(), encounterID: &encounterId, date: date, pagination: &dto.Pagination{First: &first}}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - fail to get patient observations",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchPatientObservations(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID string, searchParameters map[string]interface{}, tenant dto.TenantIdentifiers, pagination serverutils.PaginationInput) (*domain.PagedFHIRObservations, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, patientID: uuid.NewString(), encounterID: &encounterId, date: date, pagination: &dto.Pagination{First: &first}}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			_, err := clinicalUsecase.GetPatientMuacEntries(args.ctx, args.patientID, args.encounterID, args.date, args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.GetPatientMuacEntries() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCasesClinicalImpl_GetPatientOxygenSaturationEntries(t *testing.T) {
	ctx := usecaseMock.AddTenantIdentifierContext(context.Background())
	date := &scalarutils.Date{
		Year:  1997,
		Month: 12,
		Day:   12,
	}

	status := dto.ObservationStatusFinal
	valueConcept := "222"
	UUID := gofakeit.UUID()
	pagedFHIRObservations := &domain.PagedFHIRObservations{
		Observations: []domain.FHIRObservation{
			{
				ID:     new(string),
				Status: (*domain.ObservationStatusEnum)(&status),
				Code: &domain.FHIRCodeableConcept{
					ID: new(string),
					Coding: []*domain.FHIRCoding{{
						Display: gofakeit.BS(),
					}},
				},
				Subject: &domain.FHIRReference{
					ID: new(string),
				},
				Encounter: &domain.FHIRReference{
					ID: new(string),
				},
				ValueQuantity: &domain.FHIRQuantity{
					Value: 100,
					Unit:  "cm",
				},
				ValueCodeableConcept: &domain.FHIRCodeableConcept{
					Text: valueConcept,
				},
				ValueString:  new(string),
				ValueBoolean: new(bool),
				ValueInteger: new(string),
				ValueRange: &domain.FHIRRange{
					Low: domain.FHIRQuantity{
						Value: 100,
						Unit:  "cm",
					},
					High: domain.FHIRQuantity{
						Value: 100,
						Unit:  "cm",
					},
				},
				ValueRatio: &domain.FHIRRatio{
					Numerator: domain.FHIRQuantity{
						Value: 100,
						Unit:  "cm",
					},
					Denominator: domain.FHIRQuantity{
						Value: 100,
						Unit:  "cm",
					},
				},
				ValueSampledData: &domain.FHIRSampledData{
					ID: &UUID,
				},
				ValueTime: &time.Time{},
				ValueDateTime: &scalarutils.Date{
					Year:  2000,
					Month: 1,
					Day:   1,
				},
				ValuePeriod: &domain.FHIRPeriod{
					Start: scalarutils.DateTime(time.Wednesday.String()),
					End:   scalarutils.DateTime(time.Thursday.String()),
				},
			},
		},
		HasNextPage:     false,
		NextPageURL:     "",
		HasPreviousPage: false,
		PreviousPageURL: "",
		TotalCount:      0,
	}

	first := 10
	encounterId := uuid.New().String()

	type args struct {
		ctx         context.Context
		patientID   string
		encounterID *string
		date        *scalarutils.Date
		pagination  *dto.Pagination
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy case: get patient oxygen saturation",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchPatientObservations(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID string, searchParameters map[string]interface{}, tenant dto.TenantIdentifiers, pagination serverutils.PaginationInput) (*domain.PagedFHIRObservations, error) {
						return pagedFHIRObservations, nil
					})

				return args{ctx: ctx, patientID: uuid.NewString(), encounterID: &encounterId, date: date, pagination: &dto.Pagination{First: &first}}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - fail to get patient observations",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchPatientObservations(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID string, searchParameters map[string]interface{}, tenant dto.TenantIdentifiers, pagination serverutils.PaginationInput) (*domain.PagedFHIRObservations, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, patientID: uuid.NewString(), encounterID: &encounterId, date: date, pagination: &dto.Pagination{First: &first}}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			_, err := clinicalUsecase.GetPatientOxygenSaturationEntries(args.ctx, args.patientID, args.encounterID, args.date, args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.GetPatientOxygenSaturationEntries() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCasesClinicalImpl_GetPatientViralLoad(t *testing.T) {
	ctx := usecaseMock.AddTenantIdentifierContext(context.Background())
	date := &scalarutils.Date{
		Year:  1997,
		Month: 12,
		Day:   12,
	}

	status := dto.ObservationStatusFinal
	valueConcept := "222"
	UUID := gofakeit.UUID()
	pagedFHIRObservations := &domain.PagedFHIRObservations{
		Observations: []domain.FHIRObservation{
			{
				ID:     new(string),
				Status: (*domain.ObservationStatusEnum)(&status),
				Code: &domain.FHIRCodeableConcept{
					ID: new(string),
					Coding: []*domain.FHIRCoding{{
						Display: gofakeit.BS(),
					}},
				},
				Subject: &domain.FHIRReference{
					ID: new(string),
				},
				Encounter: &domain.FHIRReference{
					ID: new(string),
				},
				ValueQuantity: &domain.FHIRQuantity{
					Value: 100,
					Unit:  "cm",
				},
				ValueCodeableConcept: &domain.FHIRCodeableConcept{
					Text: valueConcept,
				},
				ValueString:  new(string),
				ValueBoolean: new(bool),
				ValueInteger: new(string),
				ValueRange: &domain.FHIRRange{
					Low: domain.FHIRQuantity{
						Value: 100,
						Unit:  "cm",
					},
					High: domain.FHIRQuantity{
						Value: 100,
						Unit:  "cm",
					},
				},
				ValueRatio: &domain.FHIRRatio{
					Numerator: domain.FHIRQuantity{
						Value: 100,
						Unit:  "cm",
					},
					Denominator: domain.FHIRQuantity{
						Value: 100,
						Unit:  "cm",
					},
				},
				ValueSampledData: &domain.FHIRSampledData{
					ID: &UUID,
				},
				ValueTime: &time.Time{},
				ValueDateTime: &scalarutils.Date{
					Year:  2000,
					Month: 1,
					Day:   1,
				},
				ValuePeriod: &domain.FHIRPeriod{
					Start: scalarutils.DateTime(time.Wednesday.String()),
					End:   scalarutils.DateTime(time.Thursday.String()),
				},
			},
		},
		HasNextPage:     false,
		NextPageURL:     "",
		HasPreviousPage: false,
		PreviousPageURL: "",
		TotalCount:      0,
	}

	first := 10
	encounterId := uuid.New().String()

	type args struct {
		ctx         context.Context
		patientID   string
		encounterID *string
		date        *scalarutils.Date
		pagination  *dto.Pagination
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy case: get patient viral load",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchPatientObservations(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID string, searchParameters map[string]interface{}, tenant dto.TenantIdentifiers, pagination serverutils.PaginationInput) (*domain.PagedFHIRObservations, error) {
						return pagedFHIRObservations, nil
					})

				return args{ctx: ctx, patientID: uuid.NewString(), encounterID: &encounterId, date: date, pagination: &dto.Pagination{First: &first}}
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to get patient viral load",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchPatientObservations(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID string, searchParameters map[string]interface{}, tenant dto.TenantIdentifiers, pagination serverutils.PaginationInput) (*domain.PagedFHIRObservations, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, patientID: uuid.NewString(), encounterID: &encounterId, date: date, pagination: &dto.Pagination{First: &first}}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			_, err := clinicalUsecase.GetPatientViralLoad(args.ctx, args.patientID, args.encounterID, args.date, args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.GetPatientViralLoad() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCasesClinicalImpl_GetPatientBloodSugarEntries(t *testing.T) {
	ctx := usecaseMock.AddTenantIdentifierContext(context.Background())
	date := &scalarutils.Date{
		Year:  1997,
		Month: 12,
		Day:   12,
	}

	status := dto.ObservationStatusFinal
	valueConcept := "222"
	UUID := gofakeit.UUID()
	pagedFHIRObservations := &domain.PagedFHIRObservations{
		Observations: []domain.FHIRObservation{
			{
				ID:     new(string),
				Status: (*domain.ObservationStatusEnum)(&status),
				Code: &domain.FHIRCodeableConcept{
					ID: new(string),
					Coding: []*domain.FHIRCoding{{
						Display: gofakeit.BS(),
					}},
				},
				Subject: &domain.FHIRReference{
					ID: new(string),
				},
				Encounter: &domain.FHIRReference{
					ID: new(string),
				},
				ValueQuantity: &domain.FHIRQuantity{
					Value: 100,
					Unit:  "cm",
				},
				ValueCodeableConcept: &domain.FHIRCodeableConcept{
					Text: valueConcept,
				},
				ValueString:  new(string),
				ValueBoolean: new(bool),
				ValueInteger: new(string),
				ValueRange: &domain.FHIRRange{
					Low: domain.FHIRQuantity{
						Value: 100,
						Unit:  "cm",
					},
					High: domain.FHIRQuantity{
						Value: 100,
						Unit:  "cm",
					},
				},
				ValueRatio: &domain.FHIRRatio{
					Numerator: domain.FHIRQuantity{
						Value: 100,
						Unit:  "cm",
					},
					Denominator: domain.FHIRQuantity{
						Value: 100,
						Unit:  "cm",
					},
				},
				ValueSampledData: &domain.FHIRSampledData{
					ID: &UUID,
				},
				ValueTime: &time.Time{},
				ValueDateTime: &scalarutils.Date{
					Year:  2000,
					Month: 1,
					Day:   1,
				},
				ValuePeriod: &domain.FHIRPeriod{
					Start: scalarutils.DateTime(time.Wednesday.String()),
					End:   scalarutils.DateTime(time.Thursday.String()),
				},
			},
		},
		HasNextPage:     false,
		NextPageURL:     "",
		HasPreviousPage: false,
		PreviousPageURL: "",
		TotalCount:      0,
	}

	first := 10
	encounterId := uuid.New().String()

	type args struct {
		ctx         context.Context
		patientID   string
		encounterID *string
		date        *scalarutils.Date
		pagination  *dto.Pagination
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy case: get patient blood sugar",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchPatientObservations(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID string, searchParameters map[string]interface{}, tenant dto.TenantIdentifiers, pagination serverutils.PaginationInput) (*domain.PagedFHIRObservations, error) {
						return pagedFHIRObservations, nil
					})

				return args{ctx: ctx, patientID: uuid.NewString(), encounterID: &encounterId, date: date, pagination: &dto.Pagination{First: &first}}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - fail to get patient observations",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchPatientObservations(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID string, searchParameters map[string]interface{}, tenant dto.TenantIdentifiers, pagination serverutils.PaginationInput) (*domain.PagedFHIRObservations, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, patientID: uuid.NewString(), encounterID: &encounterId, date: date, pagination: &dto.Pagination{First: &first}}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			_, err := clinicalUsecase.GetPatientBloodSugarEntries(args.ctx, args.patientID, args.encounterID, args.date, args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.GetPatientBloodSugarEntries() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCasesClinicalImpl_GetPatientLastMenstrualPeriodEntries(t *testing.T) {
	ctx := usecaseMock.AddTenantIdentifierContext(context.Background())
	date := &scalarutils.Date{
		Year:  1997,
		Month: 12,
		Day:   12,
	}

	status := dto.ObservationStatusFinal
	valueConcept := "222"
	UUID := gofakeit.UUID()
	pagedFHIRObservations := &domain.PagedFHIRObservations{
		Observations: []domain.FHIRObservation{
			{
				ID:     new(string),
				Status: (*domain.ObservationStatusEnum)(&status),
				Code: &domain.FHIRCodeableConcept{
					ID: new(string),
					Coding: []*domain.FHIRCoding{{
						Display: gofakeit.BS(),
					}},
				},
				Subject: &domain.FHIRReference{
					ID: new(string),
				},
				Encounter: &domain.FHIRReference{
					ID: new(string),
				},
				ValueQuantity: &domain.FHIRQuantity{
					Value: 100,
					Unit:  "cm",
				},
				ValueCodeableConcept: &domain.FHIRCodeableConcept{
					Text: valueConcept,
				},
				ValueString:  new(string),
				ValueBoolean: new(bool),
				ValueInteger: new(string),
				ValueRange: &domain.FHIRRange{
					Low: domain.FHIRQuantity{
						Value: 100,
						Unit:  "cm",
					},
					High: domain.FHIRQuantity{
						Value: 100,
						Unit:  "cm",
					},
				},
				ValueRatio: &domain.FHIRRatio{
					Numerator: domain.FHIRQuantity{
						Value: 100,
						Unit:  "cm",
					},
					Denominator: domain.FHIRQuantity{
						Value: 100,
						Unit:  "cm",
					},
				},
				ValueSampledData: &domain.FHIRSampledData{
					ID: &UUID,
				},
				ValueTime: &time.Time{},
				ValueDateTime: &scalarutils.Date{
					Year:  2000,
					Month: 1,
					Day:   1,
				},
				ValuePeriod: &domain.FHIRPeriod{
					Start: scalarutils.DateTime(time.Wednesday.String()),
					End:   scalarutils.DateTime(time.Thursday.String()),
				},
			},
		},
		HasNextPage:     false,
		NextPageURL:     "",
		HasPreviousPage: false,
		PreviousPageURL: "",
		TotalCount:      0,
	}

	first := 10
	encounterId := uuid.New().String()

	type args struct {
		ctx         context.Context
		patientID   string
		encounterID *string
		date        *scalarutils.Date
		pagination  *dto.Pagination
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy case: get patient last menstrual period",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchPatientObservations(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID string, searchParameters map[string]interface{}, tenant dto.TenantIdentifiers, pagination serverutils.PaginationInput) (*domain.PagedFHIRObservations, error) {
						return pagedFHIRObservations, nil
					})

				return args{ctx: ctx, patientID: uuid.NewString(), encounterID: &encounterId, date: date, pagination: &dto.Pagination{First: &first}}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - fail to get patient observations",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchPatientObservations(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID string, searchParameters map[string]interface{}, tenant dto.TenantIdentifiers, pagination serverutils.PaginationInput) (*domain.PagedFHIRObservations, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, patientID: uuid.NewString(), encounterID: &encounterId, date: date, pagination: &dto.Pagination{First: &first}}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			_, err := clinicalUsecase.GetPatientLastMenstrualPeriodEntries(args.ctx, args.patientID, args.encounterID, args.date, args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.GetPatientLastMenstrualPeriodEntries() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCasesClinicalImpl_GetPatientImmunoHistoChemistryRecords(t *testing.T) {
	ctx := usecaseMock.AddTenantIdentifierContext(context.Background())
	date := &scalarutils.Date{
		Year:  1997,
		Month: 12,
		Day:   12,
	}

	status := dto.ObservationStatusFinal
	valueConcept := "222"
	UUID := gofakeit.UUID()
	pagedFHIRObservations := &domain.PagedFHIRObservations{
		Observations: []domain.FHIRObservation{
			{
				ID:     new(string),
				Status: (*domain.ObservationStatusEnum)(&status),
				Code: &domain.FHIRCodeableConcept{
					ID: new(string),
					Coding: []*domain.FHIRCoding{{
						Display: gofakeit.BS(),
					}},
				},
				Subject: &domain.FHIRReference{
					ID: new(string),
				},
				Encounter: &domain.FHIRReference{
					ID: new(string),
				},
				ValueQuantity: &domain.FHIRQuantity{
					Value: 100,
					Unit:  "cm",
				},
				ValueCodeableConcept: &domain.FHIRCodeableConcept{
					Text: valueConcept,
				},
				ValueString:  new(string),
				ValueBoolean: new(bool),
				ValueInteger: new(string),
				ValueRange: &domain.FHIRRange{
					Low: domain.FHIRQuantity{
						Value: 100,
						Unit:  "cm",
					},
					High: domain.FHIRQuantity{
						Value: 100,
						Unit:  "cm",
					},
				},
				ValueRatio: &domain.FHIRRatio{
					Numerator: domain.FHIRQuantity{
						Value: 100,
						Unit:  "cm",
					},
					Denominator: domain.FHIRQuantity{
						Value: 100,
						Unit:  "cm",
					},
				},
				ValueSampledData: &domain.FHIRSampledData{
					ID: &UUID,
				},
				ValueTime: &time.Time{},
				ValueDateTime: &scalarutils.Date{
					Year:  2000,
					Month: 1,
					Day:   1,
				},
				ValuePeriod: &domain.FHIRPeriod{
					Start: scalarutils.DateTime(time.Wednesday.String()),
					End:   scalarutils.DateTime(time.Thursday.String()),
				},
			},
		},
		HasNextPage:     false,
		NextPageURL:     "",
		HasPreviousPage: false,
		PreviousPageURL: "",
		TotalCount:      0,
	}

	first := 10
	encounterId := uuid.New().String()

	type args struct {
		ctx         context.Context
		patientID   string
		encounterID *string
		date        *scalarutils.Date
		pagination  *dto.Pagination
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully get patient ihc entries",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchPatientObservations(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID string, searchParameters map[string]interface{}, tenant dto.TenantIdentifiers, pagination serverutils.PaginationInput) (*domain.PagedFHIRObservations, error) {
						return pagedFHIRObservations, nil
					})

				return args{ctx: ctx, patientID: uuid.NewString(), encounterID: &encounterId, date: date, pagination: &dto.Pagination{First: &first}}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - fail to get patient observations",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchPatientObservations(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID string, searchParameters map[string]interface{}, tenant dto.TenantIdentifiers, pagination serverutils.PaginationInput) (*domain.PagedFHIRObservations, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, patientID: uuid.NewString(), encounterID: &encounterId, date: date, pagination: &dto.Pagination{First: &first}}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			got, err := clinicalUsecase.GetPatientImmunoHistoChemistryRecords(args.ctx, args.patientID, args.encounterID, args.date, args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.GetPatientImmunoHistoChemistryRecords() error = %v, wantErr %v", err, tt.wantErr)
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

func TestUseCasesClinicalImpl_GetPatientPostCoitalBleedingRecords(t *testing.T) {
	ctx := usecaseMock.AddTenantIdentifierContext(context.Background())
	date := &scalarutils.Date{
		Year:  1997,
		Month: 12,
		Day:   12,
	}

	status := dto.ObservationStatusFinal
	valueConcept := "222"
	UUID := gofakeit.UUID()
	pagedFHIRObservations := &domain.PagedFHIRObservations{
		Observations: []domain.FHIRObservation{
			{
				ID:     new(string),
				Status: (*domain.ObservationStatusEnum)(&status),
				Code: &domain.FHIRCodeableConcept{
					ID: new(string),
					Coding: []*domain.FHIRCoding{{
						Display: gofakeit.BS(),
					}},
				},
				Subject: &domain.FHIRReference{
					ID: new(string),
				},
				Encounter: &domain.FHIRReference{
					ID: new(string),
				},
				ValueQuantity: &domain.FHIRQuantity{
					Value: 100,
					Unit:  "cm",
				},
				ValueCodeableConcept: &domain.FHIRCodeableConcept{
					Text: valueConcept,
				},
				ValueString:  new(string),
				ValueBoolean: new(bool),
				ValueInteger: new(string),
				ValueRange: &domain.FHIRRange{
					Low: domain.FHIRQuantity{
						Value: 100,
						Unit:  "cm",
					},
					High: domain.FHIRQuantity{
						Value: 100,
						Unit:  "cm",
					},
				},
				ValueRatio: &domain.FHIRRatio{
					Numerator: domain.FHIRQuantity{
						Value: 100,
						Unit:  "cm",
					},
					Denominator: domain.FHIRQuantity{
						Value: 100,
						Unit:  "cm",
					},
				},
				ValueSampledData: &domain.FHIRSampledData{
					ID: &UUID,
				},
				ValueTime: &time.Time{},
				ValueDateTime: &scalarutils.Date{
					Year:  2000,
					Month: 1,
					Day:   1,
				},
				ValuePeriod: &domain.FHIRPeriod{
					Start: scalarutils.DateTime(time.Wednesday.String()),
					End:   scalarutils.DateTime(time.Thursday.String()),
				},
			},
		},
		HasNextPage:     false,
		NextPageURL:     "",
		HasPreviousPage: false,
		PreviousPageURL: "",
		TotalCount:      0,
	}

	first := 10
	encounterId := uuid.New().String()

	type args struct {
		ctx         context.Context
		patientID   string
		encounterID *string
		date        *scalarutils.Date
		pagination  *dto.Pagination
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully get patient pcb entries",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchPatientObservations(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID string, searchParameters map[string]interface{}, tenant dto.TenantIdentifiers, pagination serverutils.PaginationInput) (*domain.PagedFHIRObservations, error) {
						return pagedFHIRObservations, nil
					})

				return args{ctx: ctx, patientID: uuid.NewString(), encounterID: &encounterId, date: date, pagination: &dto.Pagination{First: &first}}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - fail to get patient observations",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchPatientObservations(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID string, searchParameters map[string]interface{}, tenant dto.TenantIdentifiers, pagination serverutils.PaginationInput) (*domain.PagedFHIRObservations, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, patientID: uuid.NewString(), encounterID: &encounterId, date: date, pagination: &dto.Pagination{First: &first}}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			got, err := clinicalUsecase.GetPatientPostCoitalBleedingRecords(args.ctx, args.patientID, args.encounterID, args.date, args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.GetPatientPostCoitalBleedingRecords() error = %v, wantErr %v", err, tt.wantErr)
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

func TestUseCasesClinicalImpl_GetHistoryOfPresentIllness(t *testing.T) {
	ctx := usecaseMock.AddTenantIdentifierContext(context.Background())
	date := &scalarutils.Date{
		Year:  1997,
		Month: 12,
		Day:   12,
	}

	status := dto.ObservationStatusFinal
	valueConcept := "222"
	UUID := gofakeit.UUID()
	pagedFHIRObservations := &domain.PagedFHIRObservations{
		Observations: []domain.FHIRObservation{
			{
				ID:     new(string),
				Status: (*domain.ObservationStatusEnum)(&status),
				Code: &domain.FHIRCodeableConcept{
					ID: new(string),
					Coding: []*domain.FHIRCoding{{
						Display: gofakeit.BS(),
					}},
				},
				Subject: &domain.FHIRReference{
					ID: new(string),
				},
				Encounter: &domain.FHIRReference{
					ID: new(string),
				},
				ValueQuantity: &domain.FHIRQuantity{
					Value: 100,
					Unit:  "cm",
				},
				ValueCodeableConcept: &domain.FHIRCodeableConcept{
					Text: valueConcept,
				},
				ValueString:  new(string),
				ValueBoolean: new(bool),
				ValueInteger: new(string),
				ValueRange: &domain.FHIRRange{
					Low: domain.FHIRQuantity{
						Value: 100,
						Unit:  "cm",
					},
					High: domain.FHIRQuantity{
						Value: 100,
						Unit:  "cm",
					},
				},
				ValueRatio: &domain.FHIRRatio{
					Numerator: domain.FHIRQuantity{
						Value: 100,
						Unit:  "cm",
					},
					Denominator: domain.FHIRQuantity{
						Value: 100,
						Unit:  "cm",
					},
				},
				ValueSampledData: &domain.FHIRSampledData{
					ID: &UUID,
				},
				ValueTime: &time.Time{},
				ValueDateTime: &scalarutils.Date{
					Year:  2000,
					Month: 1,
					Day:   1,
				},
				ValuePeriod: &domain.FHIRPeriod{
					Start: scalarutils.DateTime(time.Wednesday.String()),
					End:   scalarutils.DateTime(time.Thursday.String()),
				},
			},
		},
		HasNextPage:     false,
		NextPageURL:     "",
		HasPreviousPage: false,
		PreviousPageURL: "",
		TotalCount:      0,
	}

	first := 10
	encounterId := uuid.New().String()

	type args struct {
		ctx         context.Context
		patientID   string
		encounterID *string
		date        *scalarutils.Date
		pagination  *dto.Pagination
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy case: get history of present illness of a patient",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchPatientObservations(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID string, searchParameters map[string]interface{}, tenant dto.TenantIdentifiers, pagination serverutils.PaginationInput) (*domain.PagedFHIRObservations, error) {
						return pagedFHIRObservations, nil
					})

				return args{ctx: ctx, patientID: uuid.NewString(), encounterID: &encounterId, date: date, pagination: &dto.Pagination{First: &first}}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - fail to get patient observations",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchPatientObservations(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID string, searchParameters map[string]interface{}, tenant dto.TenantIdentifiers, pagination serverutils.PaginationInput) (*domain.PagedFHIRObservations, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, patientID: uuid.NewString(), encounterID: &encounterId, date: date, pagination: &dto.Pagination{First: &first}}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			_, err := clinicalUsecase.GetHistoryOfPresentIllness(args.ctx, args.patientID, args.encounterID, args.date, args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.GetHistoryOfPresentIllness() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCasesClinicalImpl_GetPastMedicalAndSurgicalHistory(t *testing.T) {
	ctx := usecaseMock.AddTenantIdentifierContext(context.Background())
	date := &scalarutils.Date{
		Year:  1997,
		Month: 12,
		Day:   12,
	}

	status := dto.ObservationStatusFinal
	valueConcept := "222"
	UUID := gofakeit.UUID()
	pagedFHIRObservations := &domain.PagedFHIRObservations{
		Observations: []domain.FHIRObservation{
			{
				ID:     new(string),
				Status: (*domain.ObservationStatusEnum)(&status),
				Code: &domain.FHIRCodeableConcept{
					ID: new(string),
					Coding: []*domain.FHIRCoding{{
						Display: gofakeit.BS(),
					}},
				},
				Subject: &domain.FHIRReference{
					ID: new(string),
				},
				Encounter: &domain.FHIRReference{
					ID: new(string),
				},
				ValueQuantity: &domain.FHIRQuantity{
					Value: 100,
					Unit:  "cm",
				},
				ValueCodeableConcept: &domain.FHIRCodeableConcept{
					Text: valueConcept,
				},
				ValueString:  new(string),
				ValueBoolean: new(bool),
				ValueInteger: new(string),
				ValueRange: &domain.FHIRRange{
					Low: domain.FHIRQuantity{
						Value: 100,
						Unit:  "cm",
					},
					High: domain.FHIRQuantity{
						Value: 100,
						Unit:  "cm",
					},
				},
				ValueRatio: &domain.FHIRRatio{
					Numerator: domain.FHIRQuantity{
						Value: 100,
						Unit:  "cm",
					},
					Denominator: domain.FHIRQuantity{
						Value: 100,
						Unit:  "cm",
					},
				},
				ValueSampledData: &domain.FHIRSampledData{
					ID: &UUID,
				},
				ValueTime: &time.Time{},
				ValueDateTime: &scalarutils.Date{
					Year:  2000,
					Month: 1,
					Day:   1,
				},
				ValuePeriod: &domain.FHIRPeriod{
					Start: scalarutils.DateTime(time.Wednesday.String()),
					End:   scalarutils.DateTime(time.Thursday.String()),
				},
			},
		},
		HasNextPage:     false,
		NextPageURL:     "",
		HasPreviousPage: false,
		PreviousPageURL: "",
		TotalCount:      0,
	}

	first := 10
	encounterId := uuid.New().String()

	type args struct {
		ctx         context.Context
		patientID   string
		encounterID *string
		date        *scalarutils.Date
		pagination  *dto.Pagination
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy case: get patient past medical and surgical history",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchPatientObservations(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID string, searchParameters map[string]interface{}, tenant dto.TenantIdentifiers, pagination serverutils.PaginationInput) (*domain.PagedFHIRObservations, error) {
						return pagedFHIRObservations, nil
					})

				return args{ctx: ctx, patientID: uuid.NewString(), encounterID: &encounterId, date: date, pagination: &dto.Pagination{First: &first}}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - fail to get patient observations",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchPatientObservations(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID string, searchParameters map[string]interface{}, tenant dto.TenantIdentifiers, pagination serverutils.PaginationInput) (*domain.PagedFHIRObservations, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, patientID: uuid.NewString(), encounterID: &encounterId, date: date, pagination: &dto.Pagination{First: &first}}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			_, err := clinicalUsecase.GetPastMedicalAndSurgicalHistory(args.ctx, args.patientID, args.encounterID, args.date, args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.GetPastMedicalAndSurgicalHistory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCasesClinicalImpl_GetFamilyAndSocialHistory(t *testing.T) {
	ctx := usecaseMock.AddTenantIdentifierContext(context.Background())
	date := &scalarutils.Date{
		Year:  1997,
		Month: 12,
		Day:   12,
	}

	status := dto.ObservationStatusFinal
	valueConcept := "222"
	UUID := gofakeit.UUID()
	pagedFHIRObservations := &domain.PagedFHIRObservations{
		Observations: []domain.FHIRObservation{
			{
				ID:     new(string),
				Status: (*domain.ObservationStatusEnum)(&status),
				Code: &domain.FHIRCodeableConcept{
					ID: new(string),
					Coding: []*domain.FHIRCoding{{
						Display: gofakeit.BS(),
					}},
				},
				Subject: &domain.FHIRReference{
					ID: new(string),
				},
				Encounter: &domain.FHIRReference{
					ID: new(string),
				},
				ValueQuantity: &domain.FHIRQuantity{
					Value: 100,
					Unit:  "cm",
				},
				ValueCodeableConcept: &domain.FHIRCodeableConcept{
					Text: valueConcept,
				},
				ValueString:  new(string),
				ValueBoolean: new(bool),
				ValueInteger: new(string),
				ValueRange: &domain.FHIRRange{
					Low: domain.FHIRQuantity{
						Value: 100,
						Unit:  "cm",
					},
					High: domain.FHIRQuantity{
						Value: 100,
						Unit:  "cm",
					},
				},
				ValueRatio: &domain.FHIRRatio{
					Numerator: domain.FHIRQuantity{
						Value: 100,
						Unit:  "cm",
					},
					Denominator: domain.FHIRQuantity{
						Value: 100,
						Unit:  "cm",
					},
				},
				ValueSampledData: &domain.FHIRSampledData{
					ID: &UUID,
				},
				ValueTime: &time.Time{},
				ValueDateTime: &scalarutils.Date{
					Year:  2000,
					Month: 1,
					Day:   1,
				},
				ValuePeriod: &domain.FHIRPeriod{
					Start: scalarutils.DateTime(time.Wednesday.String()),
					End:   scalarutils.DateTime(time.Thursday.String()),
				},
			},
		},
		HasNextPage:     false,
		NextPageURL:     "",
		HasPreviousPage: false,
		PreviousPageURL: "",
		TotalCount:      0,
	}

	first := 10
	encounterId := uuid.New().String()

	type args struct {
		ctx         context.Context
		patientID   string
		encounterID *string
		date        *scalarutils.Date
		pagination  *dto.Pagination
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy case: get patient's family and social history",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchPatientObservations(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID string, searchParameters map[string]interface{}, tenant dto.TenantIdentifiers, pagination serverutils.PaginationInput) (*domain.PagedFHIRObservations, error) {
						return pagedFHIRObservations, nil
					})

				return args{ctx: ctx, patientID: uuid.NewString(), encounterID: &encounterId, date: date, pagination: &dto.Pagination{First: &first}}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - fail to get patient observations",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchPatientObservations(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID string, searchParameters map[string]interface{}, tenant dto.TenantIdentifiers, pagination serverutils.PaginationInput) (*domain.PagedFHIRObservations, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, patientID: uuid.NewString(), encounterID: &encounterId, date: date, pagination: &dto.Pagination{First: &first}}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			_, err := clinicalUsecase.GetFamilyAndSocialHistory(args.ctx, args.patientID, args.encounterID, args.date, args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.GetFamilyAndSocialHistory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCasesClinicalImpl_PatchPatientObservations(t *testing.T) {
	ctx := context.Background()

	obsRelayPayload := &domain.FHIRObservationRelayPayload{
		Resource: &domain.FHIRObservation{
			ID: &ID,
			Encounter: &domain.FHIRReference{
				ID: &ID,
			},
		},
	}

	encounter := &domain.FHIREncounterRelayPayload{
		Resource: &domain.FHIREncounter{
			ID:     &ID,
			Status: domain.EncounterStatusEnumInProgress,
		},
	}

	type args struct {
		ctx   context.Context
		id    string
		value string
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy Case - successfully patch patient observations",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRObservationRelayPayload, error) {
						return obsRelayPayload, nil
					})
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.FHIR.EXPECT().PatchFHIRObservation(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
						return observation, nil
					})

				return args{ctx: ctx, id: ID, value: "150"}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail validation nil value",
			setup: func(mh *usecaseMock.Mocks) args {
				return args{ctx: ctx, id: ID, value: ""}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - missing observation ID",
			setup: func(mh *usecaseMock.Mocks) args {
				return args{ctx: ctx, value: "150"}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - fail to get observation",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRObservationRelayPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, id: ID, value: "150"}
			},
			wantErr: true,
		},
		{
			name: "Sad Case: fail to get encounter",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRObservationRelayPayload, error) {
						return obsRelayPayload, nil
					})
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, id: ID, value: "150"}
			},
			wantErr: true,
		},
		{
			name: "Sad Case: fail on finished encounter",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRObservationRelayPayload, error) {
						return obsRelayPayload, nil
					})
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						encounter := &domain.FHIREncounterRelayPayload{
							Resource: &domain.FHIREncounter{
								ID:     &ID,
								Status: domain.EncounterStatusEnumCompleted,
							},
						}

						return encounter, nil
					})

				return args{ctx: ctx, id: ID, value: "150"}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - fail to patch patient observations",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRObservationRelayPayload, error) {
						return obsRelayPayload, nil
					})
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.FHIR.EXPECT().PatchFHIRObservation(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, id: ID, value: "150"}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		clinicalUsecase, mock := usecaseMock.SetupMocks(t)
		args := tt.setup(&mock)

		_, err := clinicalUsecase.PatchPatientObservations(args.ctx, args.id, args.value)
		if (err != nil) != tt.wantErr {
			t.Errorf("UseCasesClinicalImpl.PatchPatientObservations() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
	}
}

func TestUseCasesClinicalImpl_PatchPatientRespirationRate(t *testing.T) {
	ctx := context.Background()

	obsRelayPayload := &domain.FHIRObservationRelayPayload{
		Resource: &domain.FHIRObservation{
			ID: &ID,
			Encounter: &domain.FHIRReference{
				ID: &ID,
			},
		},
	}

	encounter := &domain.FHIREncounterRelayPayload{
		Resource: &domain.FHIREncounter{
			ID:     &ID,
			Status: domain.EncounterStatusEnumInProgress,
		},
	}

	type args struct {
		ctx   context.Context
		id    string
		value string
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy Case - successfully patch patient observation",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRObservationRelayPayload, error) {
						return obsRelayPayload, nil
					})
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.FHIR.EXPECT().PatchFHIRObservation(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
						return observation, nil
					})

				return args{ctx: ctx, id: ID, value: "160"}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - fail to patch patient respiratory rate observation",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRObservationRelayPayload, error) {
						return obsRelayPayload, nil
					})
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.FHIR.EXPECT().PatchFHIRObservation(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
						return nil, fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, id: ID, value: "160"}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		clinicalUsecase, mock := usecaseMock.SetupMocks(t)
		args := tt.setup(&mock)

		_, err := clinicalUsecase.PatchPatientRespiratoryRate(args.ctx, args.id, args.value)
		if (err != nil) != tt.wantErr {
			t.Errorf("UseCasesClinicalImpl.PatchPatientRespiratoryRate() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
	}
}

func TestUseCasesClinicalImpl_PatchPatientPulseRate(t *testing.T) {
	ctx := context.Background()

	obsRelayPayload := &domain.FHIRObservationRelayPayload{
		Resource: &domain.FHIRObservation{
			ID: &ID,
			Encounter: &domain.FHIRReference{
				ID: &ID,
			},
		},
	}

	encounter := &domain.FHIREncounterRelayPayload{
		Resource: &domain.FHIREncounter{
			ID:     &ID,
			Status: domain.EncounterStatusEnumInProgress,
		},
	}

	type args struct {
		ctx   context.Context
		id    string
		value string
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy Case - successfully patch patient observation",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRObservationRelayPayload, error) {
						return obsRelayPayload, nil
					})
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.FHIR.EXPECT().PatchFHIRObservation(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
						return observation, nil
					})

				return args{ctx: ctx, id: ID, value: "80"}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - fail to patch patient pulse rate observation",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRObservationRelayPayload, error) {
						return obsRelayPayload, nil
					})
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.FHIR.EXPECT().PatchFHIRObservation(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, id: ID, value: "80"}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		clinicalUsecase, mock := usecaseMock.SetupMocks(t)
		args := tt.setup(&mock)

		_, err := clinicalUsecase.PatchPatientPulseRate(args.ctx, args.id, args.value)
		if (err != nil) != tt.wantErr {
			t.Errorf("UseCasesClinicalImpl.PatchPatientPulseRate() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
	}
}

func TestUseCasesClinicalImpl_PatchPatientDiastolicBloodPressure(t *testing.T) {
	ctx := context.Background()

	obsRelayPayload := &domain.FHIRObservationRelayPayload{
		Resource: &domain.FHIRObservation{
			ID: &ID,
			Encounter: &domain.FHIRReference{
				ID: &ID,
			},
		},
	}

	encounter := &domain.FHIREncounterRelayPayload{
		Resource: &domain.FHIREncounter{
			ID:     &ID,
			Status: domain.EncounterStatusEnumInProgress,
		},
	}

	type args struct {
		ctx   context.Context
		id    string
		value string
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy Case - successfully patch patient observation",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRObservationRelayPayload, error) {
						return obsRelayPayload, nil
					})
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.FHIR.EXPECT().PatchFHIRObservation(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
						return observation, nil
					})

				return args{ctx: ctx, id: ID, value: "120"}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - fail to patch patient diastolic blood pressure observation",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRObservationRelayPayload, error) {
						return obsRelayPayload, nil
					})
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.FHIR.EXPECT().PatchFHIRObservation(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, id: ID, value: "120"}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		clinicalUsecase, mock := usecaseMock.SetupMocks(t)
		args := tt.setup(&mock)

		_, err := clinicalUsecase.PatchPatientDiastolicBloodPressure(args.ctx, args.id, args.value)
		if (err != nil) != tt.wantErr {
			t.Errorf("UseCasesClinicalImpl.PatchPatientBloodPressure() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
	}
}
func TestUseCasesClinicalImpl_PatchPatientTemperature(t *testing.T) {
	ctx := context.Background()

	obsRelayPayload := &domain.FHIRObservationRelayPayload{
		Resource: &domain.FHIRObservation{
			ID: &ID,
			Encounter: &domain.FHIRReference{
				ID: &ID,
			},
		},
	}

	encounter := &domain.FHIREncounterRelayPayload{
		Resource: &domain.FHIREncounter{
			ID:     &ID,
			Status: domain.EncounterStatusEnumInProgress,
		},
	}

	type args struct {
		ctx   context.Context
		id    string
		value string
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy Case - successfully patch patient observation",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRObservationRelayPayload, error) {
						return obsRelayPayload, nil
					})
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.FHIR.EXPECT().PatchFHIRObservation(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
						return observation, nil
					})

				return args{ctx: ctx, id: ID, value: "120"}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - fail to patch patient temperature observation",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRObservationRelayPayload, error) {
						return obsRelayPayload, nil
					})
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.FHIR.EXPECT().PatchFHIRObservation(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, id: ID, value: "120"}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		clinicalUsecase, mock := usecaseMock.SetupMocks(t)
		args := tt.setup(&mock)

		_, err := clinicalUsecase.PatchPatientTemperature(args.ctx, args.id, args.value)
		if (err != nil) != tt.wantErr {
			t.Errorf("UseCasesClinicalImpl.PatchPatientTemperature() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
	}
}

func TestUseCasesClinicalImpl_PatchPatientSystolicBloodPressure(t *testing.T) {
	ctx := context.Background()

	obsRelayPayload := &domain.FHIRObservationRelayPayload{
		Resource: &domain.FHIRObservation{
			ID: &ID,
			Encounter: &domain.FHIRReference{
				ID: &ID,
			},
		},
	}

	encounter := &domain.FHIREncounterRelayPayload{
		Resource: &domain.FHIREncounter{
			ID:     &ID,
			Status: domain.EncounterStatusEnumInProgress,
		},
	}

	type args struct {
		ctx   context.Context
		id    string
		value string
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy Case - successfully patch patient observation",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRObservationRelayPayload, error) {
						return obsRelayPayload, nil
					})
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.FHIR.EXPECT().PatchFHIRObservation(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
						return observation, nil
					})

				return args{ctx: ctx, id: ID, value: "120"}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - fail to patch patient systolic blood pressure observation",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRObservationRelayPayload, error) {
						return obsRelayPayload, nil
					})
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.FHIR.EXPECT().PatchFHIRObservation(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, id: ID, value: "120"}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		clinicalUsecase, mock := usecaseMock.SetupMocks(t)
		args := tt.setup(&mock)

		_, err := clinicalUsecase.PatchPatientSystolicBloodPressure(args.ctx, args.id, args.value)
		if (err != nil) != tt.wantErr {
			t.Errorf("UseCasesClinicalImpl.PatchPatientSystolicBloodPressure() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
	}
}

func TestUseCasesClinicalImpl_PatchPatientBMI(t *testing.T) {
	ctx := context.Background()

	obsRelayPayload := &domain.FHIRObservationRelayPayload{
		Resource: &domain.FHIRObservation{
			ID: &ID,
			Encounter: &domain.FHIRReference{
				ID: &ID,
			},
		},
	}

	encounter := &domain.FHIREncounterRelayPayload{
		Resource: &domain.FHIREncounter{
			ID:     &ID,
			Status: domain.EncounterStatusEnumInProgress,
		},
	}

	type args struct {
		ctx   context.Context
		id    string
		value string
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy Case - successfully patch patient observation",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRObservationRelayPayload, error) {
						return obsRelayPayload, nil
					})
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.FHIR.EXPECT().PatchFHIRObservation(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
						return observation, nil
					})

				return args{ctx: ctx, id: ID, value: "20"}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - fail to patch patient bmi observation",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRObservationRelayPayload, error) {
						return obsRelayPayload, nil
					})
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.FHIR.EXPECT().PatchFHIRObservation(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, id: ID, value: "20"}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		clinicalUsecase, mock := usecaseMock.SetupMocks(t)
		args := tt.setup(&mock)

		_, err := clinicalUsecase.PatchPatientBMI(args.ctx, args.id, args.value)
		if (err != nil) != tt.wantErr {
			t.Errorf("UseCasesClinicalImpl.PatchPatientBMI() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
	}
}

func TestUseCasesClinicalImpl_PatchPatientWeight(t *testing.T) {
	ctx := context.Background()

	obsRelayPayload := &domain.FHIRObservationRelayPayload{
		Resource: &domain.FHIRObservation{
			ID: &ID,
			Encounter: &domain.FHIRReference{
				ID: &ID,
			},
		},
	}

	encounter := &domain.FHIREncounterRelayPayload{
		Resource: &domain.FHIREncounter{
			ID:     &ID,
			Status: domain.EncounterStatusEnumInProgress,
		},
	}

	type args struct {
		ctx   context.Context
		id    string
		value string
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy Case - successfully patch patient observation",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRObservationRelayPayload, error) {
						return obsRelayPayload, nil
					})
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.FHIR.EXPECT().PatchFHIRObservation(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
						return observation, nil
					})

				return args{ctx: ctx, id: ID, value: "90"}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - fail to patch patient weight observation",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRObservationRelayPayload, error) {
						return obsRelayPayload, nil
					})
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.FHIR.EXPECT().PatchFHIRObservation(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, id: ID, value: "90"}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		clinicalUsecase, mock := usecaseMock.SetupMocks(t)
		args := tt.setup(&mock)

		_, err := clinicalUsecase.PatchPatientWeight(args.ctx, args.id, args.value)
		if (err != nil) != tt.wantErr {
			t.Errorf("UseCasesClinicalImpl.PatchPatientWeight() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
	}
}

func TestUseCasesClinicalImpl_UpdateTestResults(t *testing.T) {
	ctx := context.Background()

	obsRelayPayload := &domain.FHIRObservationRelayPayload{
		Resource: &domain.FHIRObservation{
			ID: &ID,
			Encounter: &domain.FHIRReference{
				ID: &ID,
			},
		},
	}

	encounter := &domain.FHIREncounterRelayPayload{
		Resource: &domain.FHIREncounter{
			ID:     &ID,
			Status: domain.EncounterStatusEnumInProgress,
		},
	}

	type args struct {
		ctx   context.Context
		id    string
		value string
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy Case - successfully patch patient observation",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRObservationRelayPayload, error) {
						return obsRelayPayload, nil
					})
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.FHIR.EXPECT().PatchFHIRObservation(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
						return observation, nil
					})

				return args{ctx: ctx, id: ID, value: "120"}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - fail to update patient test results",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRObservationRelayPayload, error) {
						return obsRelayPayload, nil
					})
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.FHIR.EXPECT().PatchFHIRObservation(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, id: ID, value: "120"}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		clinicalUsecase, mock := usecaseMock.SetupMocks(t)
		args := tt.setup(&mock)

		_, err := clinicalUsecase.UpdateTestResults(args.ctx, args.id, args.value)
		if (err != nil) != tt.wantErr {
			t.Errorf("UseCasesClinicalImpl.UpdateTestResults() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
	}
}

func TestUseCasesClinicalImpl_PatchPatientMuac(t *testing.T) {
	ctx := context.Background()

	obsRelayPayload := &domain.FHIRObservationRelayPayload{
		Resource: &domain.FHIRObservation{
			ID: &ID,
			Encounter: &domain.FHIRReference{
				ID: &ID,
			},
		},
	}

	encounter := &domain.FHIREncounterRelayPayload{
		Resource: &domain.FHIREncounter{
			ID:     &ID,
			Status: domain.EncounterStatusEnumInProgress,
		},
	}

	type args struct {
		ctx   context.Context
		id    string
		value string
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy Case - successfully patch patient observation",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRObservationRelayPayload, error) {
						return obsRelayPayload, nil
					})
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.FHIR.EXPECT().PatchFHIRObservation(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
						return observation, nil
					})

				return args{ctx: ctx, id: ID, value: "90"}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - fail to patch patient muac observation",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRObservationRelayPayload, error) {
						return obsRelayPayload, nil
					})
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.FHIR.EXPECT().PatchFHIRObservation(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, id: ID, value: "90"}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		clinicalUsecase, mock := usecaseMock.SetupMocks(t)
		args := tt.setup(&mock)

		_, err := clinicalUsecase.PatchPatientMuac(args.ctx, args.id, args.value)
		if (err != nil) != tt.wantErr {
			t.Errorf("UseCasesClinicalImpl.PatchPatientMuac() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
	}
}

func TestUseCasesClinicalImpl_PatchPatientOxygenSaturation(t *testing.T) {
	ctx := context.Background()

	obsRelayPayload := &domain.FHIRObservationRelayPayload{
		Resource: &domain.FHIRObservation{
			ID: &ID,
			Encounter: &domain.FHIRReference{
				ID: &ID,
			},
		},
	}

	encounter := &domain.FHIREncounterRelayPayload{
		Resource: &domain.FHIREncounter{
			ID:     &ID,
			Status: domain.EncounterStatusEnumInProgress,
		},
	}

	type args struct {
		ctx   context.Context
		id    string
		value string
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy Case - successfully patch patient observation",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRObservationRelayPayload, error) {
						return obsRelayPayload, nil
					})
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.FHIR.EXPECT().PatchFHIRObservation(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
						return observation, nil
					})

				return args{ctx: ctx, id: ID, value: "80"}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - fail to patch patient oxygen saturation observation",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRObservationRelayPayload, error) {
						return obsRelayPayload, nil
					})
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.FHIR.EXPECT().PatchFHIRObservation(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, id: ID, value: "80"}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		clinicalUsecase, mock := usecaseMock.SetupMocks(t)
		args := tt.setup(&mock)

		_, err := clinicalUsecase.PatchPatientOxygenSaturation(args.ctx, args.id, args.value)
		if (err != nil) != tt.wantErr {
			t.Errorf("UseCasesClinicalImpl.PatchPatientOxygenSaturation() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
	}
}

func TestUseCasesClinicalImpl_PatchPatientViralLoad(t *testing.T) {
	ctx := context.Background()

	obsRelayPayload := &domain.FHIRObservationRelayPayload{
		Resource: &domain.FHIRObservation{
			ID: &ID,
			Encounter: &domain.FHIRReference{
				ID: &ID,
			},
		},
	}

	encounter := &domain.FHIREncounterRelayPayload{
		Resource: &domain.FHIREncounter{
			ID:     &ID,
			Status: domain.EncounterStatusEnumInProgress,
		},
	}

	type args struct {
		ctx   context.Context
		id    string
		value string
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy Case - successfully patch patient observation",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRObservationRelayPayload, error) {
						return obsRelayPayload, nil
					})
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.FHIR.EXPECT().PatchFHIRObservation(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
						return observation, nil
					})

				return args{ctx: ctx, id: ID, value: "80"}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - fail to patch patient viral load observation",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRObservationRelayPayload, error) {
						return obsRelayPayload, nil
					})
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.FHIR.EXPECT().PatchFHIRObservation(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, id: ID, value: "80"}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		clinicalUsecase, mock := usecaseMock.SetupMocks(t)
		args := tt.setup(&mock)

		_, err := clinicalUsecase.PatchPatientViralLoad(args.ctx, args.id, args.value)
		if (err != nil) != tt.wantErr {
			t.Errorf("UseCasesClinicalImpl.PatchPatientViralLoad() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
	}
}

func TestUseCasesClinicalImpl_PatchPatientBloodSugar(t *testing.T) {
	ctx := context.Background()
	ID := uuid.NewString()

	obsRelayPayload := &domain.FHIRObservationRelayPayload{
		Resource: &domain.FHIRObservation{
			ID:   &ID,
			Text: &domain.FHIRNarrative{},
			Encounter: &domain.FHIRReference{
				ID: &ID,
			},
		},
	}

	encounter := &domain.FHIREncounterRelayPayload{
		Resource: &domain.FHIREncounter{
			ID:     &ID,
			Status: domain.EncounterStatusEnumInProgress,
		},
	}

	type args struct {
		ctx   context.Context
		id    string
		value string
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy Case - successfully patch patient observation",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRObservationRelayPayload, error) {
						return obsRelayPayload, nil
					})
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.FHIR.EXPECT().PatchFHIRObservation(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
						return observation, nil
					})

				return args{ctx: ctx, id: ID, value: "100"}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - fail to patch patient blood sugar observation",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRObservationRelayPayload, error) {
						return obsRelayPayload, nil
					})
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.FHIR.EXPECT().PatchFHIRObservation(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, id: ID, value: "100"}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		clinicalUsecase, mock := usecaseMock.SetupMocks(t)
		args := tt.setup(&mock)

		_, err := clinicalUsecase.PatchPatientBloodSugar(args.ctx, args.id, args.value)
		if (err != nil) != tt.wantErr {
			t.Errorf("UseCasesClinicalImpl.PatchPatientBloodSugar() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
	}
}

func TestUseCasesClinicalImpl_PatchPatientLastMenstrualPeriod(t *testing.T) {
	ctx := context.Background()
	ID := uuid.NewString()

	obsRelayPayload := &domain.FHIRObservationRelayPayload{
		Resource: &domain.FHIRObservation{
			ID:   &ID,
			Text: &domain.FHIRNarrative{},
			Encounter: &domain.FHIRReference{
				ID: &ID,
			},
		},
	}

	encounter := &domain.FHIREncounterRelayPayload{
		Resource: &domain.FHIREncounter{
			ID:     &ID,
			Status: domain.EncounterStatusEnumInProgress,
		},
	}

	type args struct {
		ctx   context.Context
		id    string
		value string
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy Case - successfully patch patient observation",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRObservationRelayPayload, error) {
						return obsRelayPayload, nil
					})
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.FHIR.EXPECT().PatchFHIRObservation(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
						return observation, nil
					})

				return args{ctx: ctx, id: ID, value: "20"}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - fail to patch patient last menstrual period observation",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRObservationRelayPayload, error) {
						return obsRelayPayload, nil
					})
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.FHIR.EXPECT().PatchFHIRObservation(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, id: ID, value: "20"}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		clinicalUsecase, mock := usecaseMock.SetupMocks(t)
		args := tt.setup(&mock)

		_, err := clinicalUsecase.PatchPatientLastMenstrualPeriod(args.ctx, args.id, args.value)
		if (err != nil) != tt.wantErr {
			t.Errorf("UseCasesClinicalImpl.PatchPatientLastMenstrualPeriod() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
	}
}

func TestUseCasesClinicalImpl_PatchPatientHeight(t *testing.T) {
	ctx := context.Background()
	ID := uuid.NewString()

	obsRelayPayload := &domain.FHIRObservationRelayPayload{
		Resource: &domain.FHIRObservation{
			ID:   &ID,
			Text: &domain.FHIRNarrative{},
			Encounter: &domain.FHIRReference{
				ID: &ID,
			},
		},
	}

	encounter := &domain.FHIREncounterRelayPayload{
		Resource: &domain.FHIREncounter{
			ID:     &ID,
			Status: domain.EncounterStatusEnumInProgress,
		},
	}

	type args struct {
		ctx   context.Context
		id    string
		value string
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy Case - successfully patch patient observations",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRObservationRelayPayload, error) {
						return obsRelayPayload, nil
					})
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.FHIR.EXPECT().PatchFHIRObservation(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
						return observation, nil
					})

				return args{ctx: ctx, id: ID, value: "20"}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - fail to patch patient height observation",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRObservationRelayPayload, error) {
						return obsRelayPayload, nil
					})
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.FHIR.EXPECT().PatchFHIRObservation(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, id: ID, value: "20"}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		clinicalUsecase, mock := usecaseMock.SetupMocks(t)
		args := tt.setup(&mock)

		_, err := clinicalUsecase.PatchPatientHeight(args.ctx, args.id, args.value)
		if (err != nil) != tt.wantErr {
			t.Errorf("UseCasesClinicalImpl.PatchPatientHeight() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
	}
}

func TestUseCasesClinicalImpl_ListObservations(t *testing.T) {
	ctx := usecaseMock.AddTenantIdentifierContext(context.Background())

	status := dto.ObservationStatusFinal
	valueConcept := "222"
	UUID := gofakeit.UUID()
	pagedFHIRObservations := &domain.PagedFHIRObservations{
		Observations: []domain.FHIRObservation{
			{
				ID:     new(string),
				Status: (*domain.ObservationStatusEnum)(&status),
				Code: &domain.FHIRCodeableConcept{
					ID: new(string),
					Coding: []*domain.FHIRCoding{{
						Display: gofakeit.BS(),
					}},
				},
				Subject: &domain.FHIRReference{
					ID: new(string),
				},
				Encounter: &domain.FHIRReference{
					ID: new(string),
				},
				ValueQuantity: &domain.FHIRQuantity{
					Value: 100,
					Unit:  "cm",
				},
				ValueCodeableConcept: &domain.FHIRCodeableConcept{
					Text: valueConcept,
				},
				ValueString:  new(string),
				ValueBoolean: new(bool),
				ValueInteger: new(string),
				ValueRange: &domain.FHIRRange{
					Low: domain.FHIRQuantity{
						Value: 100,
						Unit:  "cm",
					},
					High: domain.FHIRQuantity{
						Value: 100,
						Unit:  "cm",
					},
				},
				ValueRatio: &domain.FHIRRatio{
					Numerator: domain.FHIRQuantity{
						Value: 100,
						Unit:  "cm",
					},
					Denominator: domain.FHIRQuantity{
						Value: 100,
						Unit:  "cm",
					},
				},
				ValueSampledData: &domain.FHIRSampledData{
					ID: &UUID,
				},
				ValueTime: &time.Time{},
				ValueDateTime: &scalarutils.Date{
					Year:  2000,
					Month: 1,
					Day:   1,
				},
				ValuePeriod: &domain.FHIRPeriod{
					Start: scalarutils.DateTime(time.Wednesday.String()),
					End:   scalarutils.DateTime(time.Thursday.String()),
				},
			},
		},
		HasNextPage:     false,
		NextPageURL:     "",
		HasPreviousPage: false,
		PreviousPageURL: "",
		TotalCount:      0,
	}

	type args struct {
		ctx         context.Context
		patientID   string
		encounterID *string
		date        *scalarutils.Date
		category    *dto.ObservationCategory
		pagination  dto.Pagination
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy case: list observatiions",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchPatientObservations(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID string, searchParameters map[string]interface{}, tenant dto.TenantIdentifiers, pagination serverutils.PaginationInput) (*domain.PagedFHIRObservations, error) {
						return pagedFHIRObservations, nil
					})

				return args{ctx: ctx, patientID: ID, pagination: dto.Pagination{}}
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to list observatiions",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchPatientObservations(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID string, searchParameters map[string]interface{}, tenant dto.TenantIdentifiers, pagination serverutils.PaginationInput) (*domain.PagedFHIRObservations, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, patientID: ID, pagination: dto.Pagination{}}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			_, err := clinicalUsecase.ListObservations(args.ctx, args.patientID, args.encounterID, args.date, args.category, args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.ListObservations() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestClinicalImpl_GetObservationByID(t *testing.T) {
	ctx := usecaseMock.AddTenantIdentifierContext(context.Background())
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy case: Successfully get observation by ID",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRObservationRelayPayload, error) {
						return &domain.FHIRObservationRelayPayload{
							Resource: observation,
						}, nil
					})
				return args{ctx: ctx, id: gofakeit.UUID()}
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to get observation by ID",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRObservationRelayPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, id: gofakeit.UUID()}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			_, err := c.GetObservationByID(args.ctx, args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClinicalImpl.GetObservationByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestClinicalImpl_PatchPatientObservation(t *testing.T) {
	ctx := usecaseMock.AddTenantIdentifierContext(context.Background())
	uuid := gofakeit.UUID()
	status := domain.ObservationStatusEnumFinal

	type args struct {
		ctx           context.Context
		observationID string
		input         *dto.PatchObservationInput
	}
	tests := []struct {
		name    string
		args    args
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy case: patch observation",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRObservation(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRObservationRelayPayload, error) {
						return &domain.FHIRObservationRelayPayload{
							Resource: &domain.FHIRObservation{
								ID: &uuid,
								Encounter: &domain.FHIRReference{
									ID: &uuid,
								},
								Status: &status,
							},
						}, nil
					})
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return &domain.FHIREncounterRelayPayload{
							Resource: &domain.FHIREncounter{
								Status: domain.EncounterStatusEnumInProgress,
							},
						}, nil
					})
				mh.FHIR.EXPECT().PatchFHIRObservation(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
						return &domain.FHIRObservation{
							ID: &uuid,
							Encounter: &domain.FHIRReference{
								ID: &uuid,
							},
							Status: &status,
						}, nil
					})
				return args{
					ctx:           ctx,
					observationID: "123",
					input: &dto.PatchObservationInput{
						Value: gofakeit.Name(),
					},
				}
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to patch observation",
			setup: func(mh *usecaseMock.Mocks) args {
				return args{
					ctx:           ctx,
					observationID: "123",
					input: &dto.PatchObservationInput{
						Value:           gofakeit.Name(),
						ObservationType: "jdjkd",
					},
				}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			_, err := c.PatchPatientObservation(tt.args.ctx, args.observationID, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClinicalImpl.PatchPatientObservation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
