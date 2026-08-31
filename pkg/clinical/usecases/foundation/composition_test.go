package foundation_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/scalarutils"
	"github.com/stretchr/testify/mock"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/common"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/empower-clinical/pkg/clinical/domain"
	usecaseMock "github.com/savannahghi/empower-clinical/pkg/clinical/usecases/mock"
)

func TestUseCasesClinicalImpl_CreateComposition(t *testing.T) {
	ctx := context.Background()
	ID := gofakeit.UUID()
	input := dto.CompositionInput{
		EncounterID: ID,
		Type:        dto.ProgressNote,
		Category:    dto.AssessmentAndPlan,
		Status:      "final",
		Note:        "Patient is deteriorating",
	}

	UUID := uuid.New().String()
	compositionTitle := gofakeit.Name() + "assessment note"
	typeSystem := scalarutils.URI("http://hl7.org/fhir/ValueSet/doc-typecodes")
	categorySystem := scalarutils.URI("http://hl7.org/fhir/ValueSet/referenced-item-category")
	category := "Assessment + plan"
	compositionType := "Progress note"
	treatmentPlan := "Treatment Plan"
	compositionStatus := "active"
	note := scalarutils.Markdown("Fever Fever")
	PatientRef := "Patient/" + uuid.NewString()
	patientType := "Patient"
	organizationRef := "Organization/" + uuid.NewString()
	compositionSectionTextStatus := "generated"
	typeCode := scalarutils.Code(string(common.LOINCProgressNoteCode))
	categoryCode := scalarutils.Code(string(common.LOINCAssessmentPlanCode))
	compositionpayload := &domain.FHIRCompositionRelayPayload{
		Resource: &domain.FHIRComposition{
			ID:         &UUID,
			Text:       &domain.FHIRNarrative{},
			Identifier: []*domain.FHIRIdentifier{},
			Status:     (*domain.CompositionStatusEnum)(&compositionStatus),
			Type: &domain.FHIRCodeableConcept{
				ID: new(string),
				Coding: []*domain.FHIRCoding{
					{
						ID:      &UUID,
						System:  &typeSystem,
						Code:    &typeCode,
						Display: compositionType,
					},
				},
				Text: "Progress note",
			},
			Category: []*domain.FHIRCodeableConcept{
				{
					ID: new(string),
					Coding: []*domain.FHIRCoding{
						{
							ID:      &UUID,
							System:  &categorySystem,
							Version: new(string),
							Code:    &categoryCode,
							Display: category,
						},
					},
					Text: "Assessment + plan",
				},
			},
			Subject: []*domain.FHIRReference{
				{
					ID:        &UUID,
					Reference: &PatientRef,
					Type:      (*scalarutils.URI)(&patientType),
				},
			},
			Encounter: &domain.FHIRReference{
				ID: &UUID,
			},
			Date: &scalarutils.Date{
				Year:  2023,
				Month: 9,
				Day:   25,
			},
			Author: []*domain.FHIRReference{
				{
					Reference: &organizationRef,
				},
			},
			Title: &compositionTitle,
			Section: []*domain.FHIRCompositionSection{
				{
					ID:    &UUID,
					Title: &treatmentPlan,
					Code: &domain.FHIRCodeableConceptInput{
						ID: new(string),
						Coding: []*domain.FHIRCodingInput{
							{
								ID:      UUID,
								System:  &categorySystem,
								Version: new(string),
								Code:    scalarutils.Code(string(common.LOINCAssessmentPlanCode)),
								Display: category,
							},
						},
						Text: "Assessment + plan",
					},
					Author: []*domain.FHIRReference{
						{
							Reference: new(string),
						},
					},
					Text: &domain.FHIRNarrative{
						ID:     &UUID,
						Status: (*domain.NarrativeStatusEnum)(&compositionSectionTextStatus),
						Div:    scalarutils.XHTML(note),
					},
				},
			},
		},
	}

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

	patient := &domain.FHIRPatientRelayPayload{
		Resource: &domain.FHIRPatient{
			ID: &ID,
			Name: []*domain.FHIRHumanName{
				{
					ID:   &ID,
					Text: gofakeit.Name(),
				},
			},
		},
	}

	conceptpayload := &domain.Concept{
		ConceptClass: mock.Anything,
		DataType:     mock.Anything,
		ID:           gofakeit.UUID(),
	}

	orgName := mock.Anything
	organization := &domain.FHIROrganizationRelayPayload{
		Resource: &domain.FHIROrganization{
			ID:   &ID,
			Name: &orgName,
		},
	}

	type args struct {
		ctx   context.Context
		input dto.CompositionInput
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		args    args
		wantErr bool
	}{
		{
			name: "happy case: create composition - AssessmentAndPlan",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().CreateFHIRComposition(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRCompositionInput) (*domain.FHIRCompositionRelayPayload, error) {
						return compositionpayload, nil
					})
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
						return conceptpayload, nil
					})
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return patient, nil
					})
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return organization, nil
					})

				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "happy case: create composition - HistoryOfPresentingIllness",
			setup: func(mh *usecaseMock.Mocks) args {
				input.Category = dto.HistoryOfPresentingIllness

				mh.FHIR.EXPECT().CreateFHIRComposition(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRCompositionInput) (*domain.FHIRCompositionRelayPayload, error) {
						return compositionpayload, nil
					})
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
						return conceptpayload, nil
					})
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return patient, nil
					})
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return organization, nil
					})

				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "happy case: create composition - SocialHistory",
			setup: func(mh *usecaseMock.Mocks) args {
				input.Category = dto.SocialHistory

				mh.FHIR.EXPECT().CreateFHIRComposition(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRCompositionInput) (*domain.FHIRCompositionRelayPayload, error) {
						return compositionpayload, nil
					})
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
						return conceptpayload, nil
					})
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return patient, nil
					})
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return organization, nil
					})

				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "happy case: create composition - FamilyHistory",
			setup: func(mh *usecaseMock.Mocks) args {
				input.Category = dto.FamilyHistory

				mh.FHIR.EXPECT().CreateFHIRComposition(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRCompositionInput) (*domain.FHIRCompositionRelayPayload, error) {
						return compositionpayload, nil
					})
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
						return conceptpayload, nil
					})
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return patient, nil
					})
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return organization, nil
					})

				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "happy case: create composition - Examination",
			setup: func(mh *usecaseMock.Mocks) args {
				input.Category = dto.Examination

				mh.FHIR.EXPECT().CreateFHIRComposition(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRCompositionInput) (*domain.FHIRCompositionRelayPayload, error) {
						return compositionpayload, nil
					})
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
						return conceptpayload, nil
					})
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return patient, nil
					})
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return organization, nil
					})

				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "happy case: create composition - PlanOfCare",
			setup: func(mh *usecaseMock.Mocks) args {
				input.Category = dto.PlanOfCare

				mh.FHIR.EXPECT().CreateFHIRComposition(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRCompositionInput) (*domain.FHIRCompositionRelayPayload, error) {
						return compositionpayload, nil
					})
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
						return conceptpayload, nil
					})
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return patient, nil
					})
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return organization, nil
					})

				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "happy case: create composition - Past medical surgery history",
			setup: func(mh *usecaseMock.Mocks) args {
				input.Category = dto.PastMedicalSurgeryHistory

				mh.FHIR.EXPECT().CreateFHIRComposition(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRCompositionInput) (*domain.FHIRCompositionRelayPayload, error) {
						return compositionpayload, nil
					})
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
						return conceptpayload, nil
					})
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return patient, nil
					})
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return organization, nil
					})

				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "happy case: create composition - ChiefComplaint",
			setup: func(mh *usecaseMock.Mocks) args {
				input.Category = dto.ChiefComplaint

				mh.FHIR.EXPECT().CreateFHIRComposition(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRCompositionInput) (*domain.FHIRCompositionRelayPayload, error) {
						return compositionpayload, nil
					})
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
						return conceptpayload, nil
					})
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return patient, nil
					})
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return organization, nil
					})

				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "happy case: create composition - ProgressNote",
			setup: func(mh *usecaseMock.Mocks) args {
				input.Category = dto.ProviderUnspecifiedProgressNote

				mh.FHIR.EXPECT().CreateFHIRComposition(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRCompositionInput) (*domain.FHIRCompositionRelayPayload, error) {
						return compositionpayload, nil
					})
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
						return conceptpayload, nil
					})
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return patient, nil
					})
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return organization, nil
					})

				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "sad case: create composition - unsupported category",
			setup: func(mh *usecaseMock.Mocks) args {
				input := dto.CompositionInput{
					EncounterID: ID,
					Type:        dto.ProgressNote,
					Category:    dto.CompositionCategory(dto.CompositionStatusEnumFinal),
					Status:      "final",
					Note:        "Patient is deteriorating",
				}

				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return patient, nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})

				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
		{
			name: "sad case: fail to get identifiers",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return nil, fmt.Errorf("an error occurred")
					})
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return patient, nil
					})
				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
		{
			name: "sad case: error fetching concept",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
						return nil, fmt.Errorf("an error occurred")
					})
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return patient, nil
					})
				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
		{
			name: "sad case: fail to get encounter",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
		{
			name: "sad case: fail to get patient",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
		{
			name: "sad case: fail on finished encounter",
			setup: func(mh *usecaseMock.Mocks) args {
				encounter := &domain.FHIREncounterRelayPayload{
					Resource: &domain.FHIREncounter{
						ID:     &ID,
						Status: domain.EncounterStatusEnumCompleted,
					},
				}
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
		{
			name: "sad case: failed to create composition",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().CreateFHIRComposition(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRCompositionInput) (*domain.FHIRCompositionRelayPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					})
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
						return conceptpayload, nil
					})
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return patient, nil
					})
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return organization, nil
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

			_, err := clinicalUsecase.CreateComposition(args.ctx, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.CreateComposition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCasesClinicalImpl_ListPatientCompositions(t *testing.T) {
	first := 3
	EncounterID := uuid.New().String()
	ctx := context.Background()

	ID := gofakeit.UUID()
	tenantIDs := &dto.TenantIdentifiers{
		OrganizationID: ID,
		FacilityID:     ID,
	}

	patient := &domain.FHIRPatientRelayPayload{
		Resource: &domain.FHIRPatient{
			ID: &ID,
		},
	}

	ref := mock.Anything
	status := domain.CompositionStatusEnumAmended
	compositionresponse := &domain.PagedFHIRComposition{
		Compositions: []domain.FHIRComposition{
			{
				ID:     &ID,
				Status: &status,
				Type: &domain.FHIRCodeableConcept{
					ID:   &ID,
					Text: mock.Anything,
				},
				Encounter: &domain.FHIRReference{
					ID:      &ID,
					Display: mock.Anything,
				},
				Subject: []*domain.FHIRReference{
					{
						ID:      &ID,
						Display: mock.Anything,
					},
				},
				Section: []*domain.FHIRCompositionSection{
					{
						ID: &ID,
						Code: &domain.FHIRCodeableConceptInput{
							ID: &ID,
							Coding: []*domain.FHIRCodingInput{
								{
									ID:      ID,
									Display: mock.Anything,
								},
							},
						},
						Text: &domain.FHIRNarrative{},
						Author: []*domain.FHIRReference{
							{
								ID:        &ID,
								Reference: &ref,
							},
						},
						Entry:       []*domain.FHIRReference{},
						EmptyReason: &domain.FHIRCodeableConcept{},
						Section:     []*domain.FHIRCompositionSection{},
					},
				},
			},
		},
		HasNextPage:     false,
		HasPreviousPage: false,
		TotalCount:      1,
	}

	type args struct {
		ctx         context.Context
		patientID   string
		encounterID *string
		date        *scalarutils.Date
		pagination  dto.Pagination
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		args    args
		wantErr bool
	}{
		{
			name: "happy case: list compositions",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return patient, nil
					})
				mh.FHIR.EXPECT().SearchFHIRComposition(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRComposition, error) {
						return compositionresponse, nil
					})

				return args{ctx: ctx, patientID: gofakeit.UUID(), pagination: dto.Pagination{}}
			},
			wantErr: false,
		},
		{
			name: "happy case: list compositions with encounterID",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return patient, nil
					})
				mh.FHIR.EXPECT().SearchFHIRComposition(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRComposition, error) {
						return compositionresponse, nil
					})

				return args{ctx: ctx, patientID: gofakeit.UUID(), pagination: dto.Pagination{}, encounterID: &EncounterID}
			},
			wantErr: false,
		},
		{
			name: "happy case: list compositions with date",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return patient, nil
					})
				mh.FHIR.EXPECT().SearchFHIRComposition(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRComposition, error) {
						return compositionresponse, nil
					})

				return args{ctx: ctx, patientID: gofakeit.UUID(), pagination: dto.Pagination{}, date: &scalarutils.Date{Year: 2023, Month: 12, Day: 11}}
			},
			wantErr: false,
		},
		{
			name: "sad case: invalid patient id",
			setup: func(mh *usecaseMock.Mocks) args {
				return args{ctx: ctx, patientID: "invalid", pagination: dto.Pagination{}}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - invalid pagination",
			setup: func(mh *usecaseMock.Mocks) args {
				return args{ctx: ctx, patientID: gofakeit.UUID(), pagination: dto.Pagination{First: &first, Last: &first}}
			},
			wantErr: true,
		},
		{
			name: "sad case: fail to get identifiers",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, patientID: gofakeit.UUID(), pagination: dto.Pagination{}}
			},
			wantErr: true,
		},
		{
			name: "sad case: fail to get patient",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, patientID: gofakeit.UUID(), pagination: dto.Pagination{}}
			},
			wantErr: true,
		},
		{
			name: "sad case: fail to search composition",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return patient, nil
					})
				mh.FHIR.EXPECT().SearchFHIRComposition(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRComposition, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, patientID: gofakeit.UUID(), pagination: dto.Pagination{}}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			got, err := clinicalUsecase.ListPatientCompositions(args.ctx, args.patientID, args.encounterID, args.date, args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListPatientCompositions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a value to be returned, got: %v", got)
				return
			}
		})
	}
}

func TestUseCasesClinicalImpl_AppendNoteToComposition(t *testing.T) {
	ctx := context.Background()
	ID := uuid.NewString()
	input := dto.PatchCompositionInput{
		Type:     dto.ProgressNote,
		Category: dto.HistoryOfPresentingIllness,
		Status:   dto.CompositionStatusEnumFinal,
		Note:     "Patient condition is deteriorating",
	}

	tenantIDs := &dto.TenantIdentifiers{
		OrganizationID: ID,
		FacilityID:     ID,
	}

	encounter := &domain.FHIREncounterRelayPayload{
		Resource: &domain.FHIREncounter{
			ID: &ID,
			Subject: &domain.FHIRReference{
				ID:      &ID,
				Display: ID,
			},
		},
	}

	conceptpayload := &domain.Concept{
		ConceptClass: mock.Anything,
		DataType:     mock.Anything,
		ID:           gofakeit.UUID(),
	}

	ref := mock.Anything
	composition := &domain.FHIRCompositionRelayPayload{
		Resource: &domain.FHIRComposition{
			ID: &ID,
			Encounter: &domain.FHIRReference{
				ID:      &ID,
				Display: mock.Anything,
			},
			Section: []*domain.FHIRCompositionSection{
				{
					ID: &ID,
					Code: &domain.FHIRCodeableConceptInput{
						ID:     &ID,
						Coding: []*domain.FHIRCodingInput{},
					},
					Text:        &domain.FHIRNarrative{},
					Author:      []*domain.FHIRReference{},
					Entry:       []*domain.FHIRReference{},
					EmptyReason: &domain.FHIRCodeableConcept{},
					Section: []*domain.FHIRCompositionSection{
						{
							ID: &ID,
							Code: &domain.FHIRCodeableConceptInput{
								ID: &ID,
								Coding: []*domain.FHIRCodingInput{
									{
										ID:      ID,
										Display: mock.Anything,
									},
								},
							},
							Text: &domain.FHIRNarrative{},
							Author: []*domain.FHIRReference{
								{
									ID:        &ID,
									Reference: &ref,
								},
							},
							Entry:       []*domain.FHIRReference{},
							EmptyReason: &domain.FHIRCodeableConcept{},
							Section:     []*domain.FHIRCompositionSection{},
						},
					},
				},
			},
		},
	}

	status := domain.CompositionStatusEnumAmended
	compositionoutput := &domain.FHIRComposition{
		ID:     &ID,
		Status: &status,
		Type: &domain.FHIRCodeableConcept{
			ID:   &ID,
			Text: mock.Anything,
		},
		Encounter: &domain.FHIRReference{
			ID:      &ID,
			Display: mock.Anything,
		},
		Subject: []*domain.FHIRReference{
			{
				ID:      &ID,
				Display: mock.Anything,
			},
		},
		Section: []*domain.FHIRCompositionSection{
			{
				ID: &ID,
				Code: &domain.FHIRCodeableConceptInput{
					ID: &ID,
					Coding: []*domain.FHIRCodingInput{
						{
							ID:      ID,
							Display: mock.Anything,
						},
					},
				},
				Text: &domain.FHIRNarrative{},
				Author: []*domain.FHIRReference{
					{
						ID:        &ID,
						Reference: &ref,
					},
				},
				Entry:       []*domain.FHIRReference{},
				EmptyReason: &domain.FHIRCodeableConcept{},
				Section:     []*domain.FHIRCompositionSection{},
			},
		},
	}

	type args struct {
		ctx   context.Context
		id    string
		input dto.PatchCompositionInput
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy Case: Successfully patch a composition",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
						return conceptpayload, nil
					})
				mh.FHIR.EXPECT().GetFHIRComposition(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRCompositionRelayPayload, error) {
						return composition, nil
					})
				mh.FHIR.EXPECT().PatchFHIRComposition(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string, input domain.FHIRCompositionInput) (*domain.FHIRComposition, error) {
						return compositionoutput, nil
					})

				return args{ctx: ctx, id: ID, input: input}
			},
			wantErr: false,
		},
		{
			name: "Happy Case: Successfully patch a composition - with FamilyHistory",
			setup: func(mh *usecaseMock.Mocks) args {
				input.Category = dto.FamilyHistory

				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
						return conceptpayload, nil
					})
				mh.FHIR.EXPECT().GetFHIRComposition(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRCompositionRelayPayload, error) {
						return composition, nil
					})
				mh.FHIR.EXPECT().PatchFHIRComposition(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string, input domain.FHIRCompositionInput) (*domain.FHIRComposition, error) {
						return compositionoutput, nil
					})

				return args{ctx: ctx, id: ID, input: input}
			},
			wantErr: false,
		},
		{
			name: "Happy Case: Successfully patch a composition - with Examination",
			setup: func(mh *usecaseMock.Mocks) args {
				input.Category = dto.Examination

				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
						return conceptpayload, nil
					})
				mh.FHIR.EXPECT().GetFHIRComposition(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRCompositionRelayPayload, error) {
						return composition, nil
					})
				mh.FHIR.EXPECT().PatchFHIRComposition(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string, input domain.FHIRCompositionInput) (*domain.FHIRComposition, error) {
						return compositionoutput, nil
					})

				return args{ctx: ctx, id: ID, input: input}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to get encounter",
			setup: func(mh *usecaseMock.Mocks) args {
				input.Category = dto.Examination

				mh.FHIR.EXPECT().GetFHIRComposition(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRCompositionRelayPayload, error) {
						return composition, nil
					})
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, id: ID, input: input}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - return a finished encounter",
			setup: func(mh *usecaseMock.Mocks) args {
				input.Category = dto.Examination
				encounter := &domain.FHIREncounterRelayPayload{
					Resource: &domain.FHIREncounter{
						ID:     &ID,
						Status: domain.EncounterStatusEnumCompleted,
					},
				}

				mh.FHIR.EXPECT().GetFHIRComposition(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRCompositionRelayPayload, error) {
						return composition, nil
					})
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})

				return args{ctx: ctx, id: ID, input: input}
			},
			wantErr: true,
		},
		{
			name: "Sad Case: Missing composition id",
			setup: func(mh *usecaseMock.Mocks) args {
				return args{ctx: ctx, id: "", input: input}
			},
			wantErr: true,
		},
		{
			name: "Sad Case: Fail to get OCL concept",
			setup: func(mh *usecaseMock.Mocks) args {
				input.Category = dto.SocialHistory

				mh.FHIR.EXPECT().GetFHIRComposition(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRCompositionRelayPayload, error) {
						return composition, nil
					})
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
						return nil, fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, id: ID, input: input}
			},
			wantErr: true,
		},
		{
			name: "Sad Case: Fail to get OCL concept - no type",
			setup: func(mh *usecaseMock.Mocks) args {
				input := dto.PatchCompositionInput{
					Category: dto.PlanOfCare,
					Status:   "final",
					Note:     "Patient is deteriorating",
				}

				mh.FHIR.EXPECT().GetFHIRComposition(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRCompositionRelayPayload, error) {
						return composition, nil
					})
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, id: ID, input: input}
			},
			wantErr: true,
		},
		{
			name: "Sad Case: Fail to get tenant identifiers",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRComposition(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRCompositionRelayPayload, error) {
						return composition, nil
					})
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, id: ID, input: input}
			},
			wantErr: true,
		},
		{
			name: "Sad Case: Wrong/Missing Category code",
			setup: func(mh *usecaseMock.Mocks) args {
				input := dto.PatchCompositionInput{
					Type:     dto.ProgressNote,
					Category: "invalid",
					Status:   dto.CompositionStatusEnumFinal,
					Note:     "Patient condition is deteriorating",
				}
				mh.FHIR.EXPECT().GetFHIRComposition(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRCompositionRelayPayload, error) {
						return composition, nil
					})
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})

				return args{ctx: ctx, id: ID, input: input}
			},
			wantErr: true,
		},
		{
			name: "Sad Case: Fail to get fhir composition",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRComposition(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRCompositionRelayPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, id: ID, input: input}
			},
			wantErr: true,
		},
		{
			name: "Sad Case: Fail to patch composition",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
						return conceptpayload, nil
					})
				mh.FHIR.EXPECT().GetFHIRComposition(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRCompositionRelayPayload, error) {
						return composition, nil
					})
				mh.FHIR.EXPECT().PatchFHIRComposition(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string, input domain.FHIRCompositionInput) (*domain.FHIRComposition, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, id: ID, input: input}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			_, err := clinicalUsecase.AppendNoteToComposition(args.ctx, args.id, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("AppendNoteToComposition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
