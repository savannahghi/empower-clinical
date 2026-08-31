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

func TestUseCasesClinicalImpl_CreateCondition(t *testing.T) {
	ctx := context.Background()
	input := dto.ConditionInput{
		Code:        "386661006",
		System:      domain.TerminologySourceICD11WHO,
		Name:        "Malaria",
		Status:      dto.ConditionStatusActive,
		Category:    dto.ConditionCategoryProblemList,
		EncounterID: gofakeit.UUID(),
		Note:        "Fever Fever",
		OnsetDate: &scalarutils.Date{
			Year:  2022,
			Month: 12,
			Day:   12,
		},
	}

	UUID := uuid.New().String()
	statusSystem := scalarutils.URI("http://terminology.hl7.org/CodeSystem/condition-clinical")
	status := "active"
	note := scalarutils.Markdown("Fever Fever")
	today := time.Now().Format(time.RFC3339)
	currentTime := (*scalarutils.DateTime)(&today)
	uri := scalarutils.URI("1234567")
	codingCode := "1234"
	categoryCode := "ENCOUNTER_DIAGNOSIS"
	condition := &domain.FHIRConditionRelayPayload{
		Resource: &domain.FHIRCondition{
			ID:         &UUID,
			Text:       &domain.FHIRNarrative{},
			Identifier: []*domain.FHIRIdentifier{},
			ClinicalStatus: &domain.FHIRCodeableConcept{
				Coding: []*domain.FHIRCoding{
					{
						System:  &statusSystem,
						Code:    (*scalarutils.Code)(&status),
						Display: string(status),
					},
				},
				Text: string(status),
			},
			Code: &domain.FHIRCodeableConcept{
				Coding: []*domain.FHIRCoding{
					{
						System:  &uri,
						Code:    (*scalarutils.Code)(&codingCode),
						Display: "1234",
					},
				},
				Text: "1234",
			},
			OnsetDateTime: &scalarutils.Date{},
			RecordedDate:  &scalarutils.Date{},
			Note: []*domain.FHIRAnnotation{
				{
					Time: currentTime,
					Text: &note,
				},
			},
			Subject: &domain.FHIRReference{
				ID: &UUID,
			},
			Encounter: &domain.FHIRReference{
				ID: &UUID,
			},
			Category: []*domain.FHIRCodeableConcept{
				{
					ID: &UUID,
					Coding: []*domain.FHIRCoding{
						{
							ID:           &UUID,
							System:       (*scalarutils.URI)(&UUID),
							Version:      &UUID,
							Code:         (*scalarutils.Code)(&categoryCode),
							Display:      gofakeit.BeerAlcohol(),
							UserSelected: new(bool),
						},
					},
					Text: "ENCOUNTER_DIAGNOSIS",
				},
			},
		},
	}

	invalidCategoryCode := "INVALID"
	invalidCondition := &domain.FHIRConditionRelayPayload{
		Resource: &domain.FHIRCondition{
			ID:         &UUID,
			Text:       &domain.FHIRNarrative{},
			Identifier: []*domain.FHIRIdentifier{},
			ClinicalStatus: &domain.FHIRCodeableConcept{
				Coding: []*domain.FHIRCoding{
					{
						System:  &statusSystem,
						Code:    (*scalarutils.Code)(&status),
						Display: string(status),
					},
				},
				Text: string(status),
			},
			Code: &domain.FHIRCodeableConcept{
				Coding: []*domain.FHIRCoding{
					{
						System:  &uri,
						Code:    (*scalarutils.Code)(&codingCode),
						Display: "1234",
					},
				},
				Text: "1234",
			},
			OnsetDateTime: &scalarutils.Date{},
			RecordedDate:  &scalarutils.Date{},
			Note: []*domain.FHIRAnnotation{
				{
					Time: currentTime,
					Text: &note,
				},
			},
			Subject: &domain.FHIRReference{
				ID: &UUID,
			},
			Encounter: &domain.FHIRReference{
				ID: &UUID,
			},
			Category: []*domain.FHIRCodeableConcept{
				{
					ID: &UUID,
					Coding: []*domain.FHIRCoding{
						{
							ID:           &UUID,
							System:       (*scalarutils.URI)(&UUID),
							Version:      &UUID,
							Code:         (*scalarutils.Code)(&invalidCategoryCode),
							Display:      gofakeit.BeerAlcohol(),
							UserSelected: new(bool),
						},
					},
					Text: "INVALID",
				},
			},
		},
	}

	patientRef := "Patient/" + uuid.NewString()
	encounter := &domain.FHIREncounterRelayPayload{
		Resource: &domain.FHIREncounter{
			ID:            &UUID,
			Text:          &domain.FHIRNarrative{},
			Identifier:    []*domain.FHIRIdentifier{},
			Status:        domain.EncounterStatusEnum(domain.EncounterStatusEnumDischarged),
			StatusHistory: []*domain.FHIREncounterStatushistory{},
			Class:         []domain.FHIRCodeableConcept{},
			ClassHistory:  []*domain.FHIREncounterClasshistory{},
			Type:          []*domain.FHIRCodeableConcept{},
			ServiceType:   &domain.FHIRCodeableConcept{},
			Priority:      &domain.FHIRCodeableConcept{},
			Subject: &domain.FHIRReference{
				ID:        &UUID,
				Reference: &patientRef,
			},
			EpisodeOfCare:   []*domain.FHIRReference{},
			BasedOn:         []*domain.FHIRReference{},
			Participant:     []*domain.FHIREncounterParticipant{},
			Appointment:     []*domain.FHIRReference{},
			ActualPeriod:    &domain.FHIRPeriod{},
			Length:          &domain.FHIRDuration{},
			ReasonReference: []*domain.FHIRReference{},
			Diagnosis:       []*domain.FHIREncounterDiagnosis{},
			Account:         []*domain.FHIRReference{},
			Hospitalization: &domain.FHIREncounterHospitalization{},
			Location:        []*domain.FHIREncounterLocation{},
			ServiceProvider: &domain.FHIRReference{},
			PartOf:          &domain.FHIRReference{},
		},
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

	type args struct {
		ctx   context.Context
		input dto.ConditionInput
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "happy case: create condition - encounter diagnosis",
			setup: func(mh *usecaseMock.Mocks) args {
				input.Category = dto.ConditionCategoryDiagnosis

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
				mh.FHIR.EXPECT().CreateFHIRCondition(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRConditionInput) (*domain.FHIRConditionRelayPayload, error) {
						return condition, nil
					})

				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "happy case: create condition - problem list",
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
				mh.FHIR.EXPECT().CreateFHIRCondition(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRConditionInput) (*domain.FHIRConditionRelayPayload, error) {
						condition.Resource.Category[0].Text = "PROBLEM_LIST_ITEM"
						return condition, nil
					})

				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "happy case: create condition -  invalid category",
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
				mh.FHIR.EXPECT().CreateFHIRCondition(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRConditionInput) (*domain.FHIRConditionRelayPayload, error) {
						return invalidCondition, nil
					})

				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "sad case: fail to get tags",
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
						return nil, fmt.Errorf("an error occurred")
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
			name: "Sad Case - fail on finished encounter",
			setup: func(mh *usecaseMock.Mocks) args {
				encounter := &domain.FHIREncounterRelayPayload{
					Resource: &domain.FHIREncounter{
						ID:            &UUID,
						Text:          &domain.FHIRNarrative{},
						Identifier:    []*domain.FHIRIdentifier{},
						Status:        domain.EncounterStatusEnum(domain.EncounterStatusEnumCompleted),
						StatusHistory: []*domain.FHIREncounterStatushistory{},
						Class:         []domain.FHIRCodeableConcept{},
						ClassHistory:  []*domain.FHIREncounterClasshistory{},
						Type:          []*domain.FHIRCodeableConcept{},
						ServiceType:   &domain.FHIRCodeableConcept{},
						Priority:      &domain.FHIRCodeableConcept{},
						Subject: &domain.FHIRReference{
							ID:        &UUID,
							Reference: &patientRef,
						},
						EpisodeOfCare:   []*domain.FHIRReference{},
						BasedOn:         []*domain.FHIRReference{},
						Participant:     []*domain.FHIREncounterParticipant{},
						Appointment:     []*domain.FHIRReference{},
						ActualPeriod:    &domain.FHIRPeriod{},
						Length:          &domain.FHIRDuration{},
						ReasonReference: []*domain.FHIRReference{},
						Diagnosis:       []*domain.FHIREncounterDiagnosis{},
						Account:         []*domain.FHIRReference{},
						Hospitalization: &domain.FHIREncounterHospitalization{},
						Location:        []*domain.FHIREncounterLocation{},
						ServiceProvider: &domain.FHIRReference{},
						PartOf:          &domain.FHIRReference{},
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
			name: "sad case: failed to create condition",
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
				mh.FHIR.EXPECT().CreateFHIRCondition(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRConditionInput) (*domain.FHIRConditionRelayPayload, error) {
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

			_, err := clinicalUsecase.CreateCondition(args.ctx, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.CreateCondition() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUseCasesClinicalImpl_RecordTreatmentEnrollment(t *testing.T) {
	ctx := context.Background()
	input := &dto.TreatmentEnrollmentInput{
		EncounterID: gofakeit.UUID(),
		Condition: dto.ValueSetData{
			Code:    "2A21",
			Display: "Malaria",
		},
		Date: &scalarutils.Date{
			Year:  2026,
			Month: 6,
			Day:   30,
		},
		LinkedToTreatment: true,
		TreatmentFacility: "Nairobi Cancer Centre",
		TreatmentProgram:  "Oncology Treatment Program",
		EnrollmentDate: &scalarutils.Date{
			Year:  2026,
			Month: 6,
			Day:   30,
		},
		Severity: dto.ConditionSeverityMild,
	}

	UUID := uuid.New().String()
	patientRef := "Patient/" + uuid.NewString()

	encounter := &domain.FHIREncounterRelayPayload{
		Resource: &domain.FHIREncounter{
			ID:     &UUID,
			Status: domain.EncounterStatusEnum(domain.EncounterStatusEnumInProgress),
			Subject: &domain.FHIRReference{
				ID:        &UUID,
				Reference: &patientRef,
				Display:   patientRef,
			},
		},
	}

	tenantIDs := &dto.TenantIdentifiers{
		OrganizationID: uuid.NewString(),
		FacilityID:     uuid.NewString(),
	}

	orgName := "Test Facility"
	orgID := uuid.NewString()
	organization := &domain.FHIROrganizationRelayPayload{
		Resource: &domain.FHIROrganization{
			ID:   &orgID,
			Name: &orgName,
		},
	}

	linkedToTreatment := true
	condition := &domain.FHIRConditionRelayPayload{
		Resource: &domain.FHIRCondition{
			ID: &UUID,
			ClinicalStatus: &domain.FHIRCodeableConcept{
				Text: "Active",
			},
			Code: &domain.FHIRCodeableConcept{
				Text: "Malaria",
			},
			RecordedDate: input.Date,
			Subject: &domain.FHIRReference{
				Display: patientRef,
			},
			Encounter: &domain.FHIRReference{
				Display: UUID,
			},
			Extension: []*domain.FHIRExtension{
				{
					URL: "http://savannahghi.org/fhir/StructureDefinition/diagnosis-treatment-linkage",
					Extension: []domain.Extension{
						{URL: "linkedToTreatment", ValueBoolean: &linkedToTreatment},
						{URL: "treatmentFacility", ValueString: "Nairobi Cancer Centre"},
						{URL: "treatmentProgram", ValueString: "Oncology Treatment Program"},
						{URL: "enrollmentDate", ValueDate: "2026-06-30"},
					},
				},
			},
		},
	}

	type args struct {
		ctx   context.Context
		input *dto.TreatmentEnrollmentInput
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "happy case: link diagnosis to treatment",
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
				mh.FHIR.EXPECT().CreateFHIRCondition(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRConditionInput) (*domain.FHIRConditionRelayPayload, error) {
						return condition, nil
					})

				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "sad case: invalid input",
			setup: func(mh *usecaseMock.Mocks) args {
				return args{ctx: ctx, input: &dto.TreatmentEnrollmentInput{}}
			},
			wantErr: true,
		},
		{
			name: "sad case: linked to treatment without a treatment facility",
			setup: func(mh *usecaseMock.Mocks) args {
				return args{ctx: ctx, input: &dto.TreatmentEnrollmentInput{
					EncounterID: gofakeit.UUID(),
					Condition: dto.ValueSetData{
						Code:    "2A21",
						Display: "Malaria",
					},
					Date: &scalarutils.Date{
						Year:  2026,
						Month: 6,
						Day:   30,
					},
					LinkedToTreatment: true,
				}}
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
			name: "sad case: cannot record in a completed encounter",
			setup: func(mh *usecaseMock.Mocks) args {
				completedEncounter := &domain.FHIREncounterRelayPayload{
					Resource: &domain.FHIREncounter{
						ID:     &UUID,
						Status: domain.EncounterStatusEnum(domain.EncounterStatusEnumCompleted),
						Subject: &domain.FHIRReference{
							ID:        &UUID,
							Reference: &patientRef,
						},
					},
				}
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return completedEncounter, nil
					})

				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
		{
			name: "sad case: fail to get tenant tags",
			setup: func(mh *usecaseMock.Mocks) args {
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
			name: "sad case: fail to create condition",
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
				mh.FHIR.EXPECT().CreateFHIRCondition(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRConditionInput) (*domain.FHIRConditionRelayPayload, error) {
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

			got, err := clinicalUsecase.RecordTreatmentEnrollment(args.ctx, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClinicalImpl.RecordTreatmentEnrollment() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				if got == nil || got.TreatmentLinkage == nil {
					t.Fatalf("expected treatment linkage in the response, got %+v", got)
				}

				if !got.TreatmentLinkage.LinkedToTreatment {
					t.Errorf("expected LinkedToTreatment to be true")
				}

				if got.TreatmentLinkage.TreatmentFacility != "Nairobi Cancer Centre" {
					t.Errorf("expected treatment facility %q, got %q", "Nairobi Cancer Centre", got.TreatmentLinkage.TreatmentFacility)
				}

				if got.TreatmentLinkage.TreatmentProgram != "Oncology Treatment Program" {
					t.Errorf("expected treatment program %q, got %q", "Oncology Treatment Program", got.TreatmentLinkage.TreatmentProgram)
				}

				if got.TreatmentLinkage.EnrollmentDate == nil {
					t.Errorf("expected enrollment date to be set")
				}
			}
		})
	}
}

func TestUseCasesClinicalImpl_UpdateTreatmentEnrollment(t *testing.T) {
	ctx := context.Background()
	id := uuid.NewString()
	patientRef := "Patient/" + uuid.NewString()
	linkedToTreatment := true

	// enrollment is an existing treatment enrollment carrying the diagnosis-treatment linkage.
	enrollment := func() *domain.FHIRConditionRelayPayload {
		return &domain.FHIRConditionRelayPayload{
			Resource: &domain.FHIRCondition{
				ID:             &id,
				ClinicalStatus: &domain.FHIRCodeableConcept{Text: "Active"},
				Code:           &domain.FHIRCodeableConcept{Text: "Malaria"},
				RecordedDate:   &scalarutils.Date{Year: 2026, Month: 6, Day: 30},
				Subject:        &domain.FHIRReference{Display: patientRef},
				Encounter:      &domain.FHIRReference{Display: id},
				Extension: []*domain.FHIRExtension{
					{
						URL: "http://savannahghi.org/fhir/StructureDefinition/diagnosis-treatment-linkage",
						Extension: []domain.Extension{
							{URL: "linkedToTreatment", ValueBoolean: &linkedToTreatment},
							{URL: "treatmentFacility", ValueString: "Nairobi Cancer Centre"},
							{URL: "treatmentProgram", ValueString: "Oncology Treatment Program"},
							{URL: "enrollmentDate", ValueDate: "2026-06-30"},
						},
					},
				},
			},
		}
	}

	// enrollmentDateOf returns the linkage's enrollmentDate sub-extension value.
	enrollmentDateOf := func(extensions []*domain.FHIRExtension) string {
		for _, extension := range extensions {
			if extension == nil {
				continue
			}
			for _, field := range extension.Extension {
				if field.URL == "enrollmentDate" {
					return field.ValueDate
				}
			}
		}

		return ""
	}

	tenantIDs := &dto.TenantIdentifiers{
		OrganizationID: uuid.NewString(),
		FacilityID:     uuid.NewString(),
	}

	orgName := "Test Facility"
	orgID := uuid.NewString()
	organization := &domain.FHIROrganizationRelayPayload{
		Resource: &domain.FHIROrganization{ID: &orgID, Name: &orgName},
	}

	type args struct {
		ctx   context.Context
		id    string
		input *dto.UpdateTreatmentEnrollmentInput
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "happy case: update condition, date and enrollment date",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRCondition(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRConditionRelayPayload, error) {
						return enrollment(), nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return organization, nil
					})
				mh.FHIR.EXPECT().UpdateFHIRCondition(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, in domain.FHIRConditionInput) (*domain.FHIRConditionRelayPayload, error) {
						if in.ID == nil || *in.ID != id {
							t.Errorf("expected condition id %q to be set on the update", id)
						}
						if in.Code == nil {
							t.Errorf("expected code to be updated")
						}
						if in.RecordedDate == nil {
							t.Errorf("expected recorded date to be updated")
						}
						if len(in.Meta.Tag) == 0 {
							t.Errorf("expected meta tags to be preserved to avoid wiping tenant tags")
						}
						if got := enrollmentDateOf(in.Extension); got != "2027-02-20" {
							t.Errorf("enrollment date = %q, want %q", got, "2027-02-20")
						}

						return enrollment(), nil
					})

				return args{ctx: ctx, id: id, input: &dto.UpdateTreatmentEnrollmentInput{
					Condition:      &dto.ValueSetData{Code: "2A22", Display: "Updated Condition"},
					Date:           &scalarutils.Date{Year: 2027, Month: 1, Day: 15},
					EnrollmentDate: &scalarutils.Date{Year: 2027, Month: 2, Day: 20},
				}}
			},
			wantErr: false,
		},
		{
			name: "happy case: update only enrollment date preserves the rest of the linkage",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRCondition(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRConditionRelayPayload, error) {
						return enrollment(), nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return organization, nil
					})
				mh.FHIR.EXPECT().UpdateFHIRCondition(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, in domain.FHIRConditionInput) (*domain.FHIRConditionRelayPayload, error) {
						if in.Code != nil {
							t.Errorf("expected code to be left unchanged")
						}
						if in.RecordedDate != nil {
							t.Errorf("expected recorded date to be left unchanged")
						}
						if got := enrollmentDateOf(in.Extension); got != "2027-03-10" {
							t.Errorf("enrollment date = %q, want %q", got, "2027-03-10")
						}
						// The other linkage sub-extensions must survive the enrollment-date change.
						linkage := treatmentLinkageForTest(in.Extension)
						if linkage == nil || linkage["treatmentFacility"] != "Nairobi Cancer Centre" || linkage["treatmentProgram"] != "Oncology Treatment Program" {
							t.Errorf("expected treatment facility/program to be preserved, got %+v", linkage)
						}

						return enrollment(), nil
					})

				return args{ctx: ctx, id: id, input: &dto.UpdateTreatmentEnrollmentInput{
					EnrollmentDate: &scalarutils.Date{Year: 2027, Month: 3, Day: 10},
				}}
			},
			wantErr: false,
		},
		{
			name: "sad case: missing id",
			setup: func(mh *usecaseMock.Mocks) args {
				return args{ctx: ctx, id: "", input: &dto.UpdateTreatmentEnrollmentInput{
					Date: &scalarutils.Date{Year: 2027, Month: 1, Day: 15},
				}}
			},
			wantErr: true,
		},
		{
			name: "sad case: no fields provided",
			setup: func(mh *usecaseMock.Mocks) args {
				return args{ctx: ctx, id: id, input: &dto.UpdateTreatmentEnrollmentInput{}}
			},
			wantErr: true,
		},
		{
			name: "sad case: condition not found",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRCondition(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRConditionRelayPayload, error) {
						return nil, fmt.Errorf("not found")
					})

				return args{ctx: ctx, id: id, input: &dto.UpdateTreatmentEnrollmentInput{
					Date: &scalarutils.Date{Year: 2027, Month: 1, Day: 15},
				}}
			},
			wantErr: true,
		},
		{
			name: "sad case: condition is not a treatment enrollment",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRCondition(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRConditionRelayPayload, error) {
						// A plain condition with no diagnosis-treatment linkage extension.
						return &domain.FHIRConditionRelayPayload{
							Resource: &domain.FHIRCondition{ID: &id, Code: &domain.FHIRCodeableConcept{Text: "Malaria"}},
						}, nil
					})

				return args{ctx: ctx, id: id, input: &dto.UpdateTreatmentEnrollmentInput{
					Date: &scalarutils.Date{Year: 2027, Month: 1, Day: 15},
				}}
			},
			wantErr: true,
		},
		{
			name: "sad case: update fails",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRCondition(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRConditionRelayPayload, error) {
						return enrollment(), nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return organization, nil
					})
				mh.FHIR.EXPECT().UpdateFHIRCondition(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, in domain.FHIRConditionInput) (*domain.FHIRConditionRelayPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, id: id, input: &dto.UpdateTreatmentEnrollmentInput{
					Date: &scalarutils.Date{Year: 2027, Month: 1, Day: 15},
				}}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			got, err := clinicalUsecase.UpdateTreatmentEnrollment(args.ctx, args.id, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClinicalImpl.UpdateTreatmentEnrollment() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got nil")
			}
		})
	}
}

// treatmentLinkageForTest flattens the diagnosis-treatment linkage sub-extensions' string values
// into a map keyed by sub-extension URL, for convenient assertions in tests.
func treatmentLinkageForTest(extensions []*domain.FHIRExtension) map[string]string {
	for _, extension := range extensions {
		if extension == nil || extension.URL != "http://savannahghi.org/fhir/StructureDefinition/diagnosis-treatment-linkage" {
			continue
		}

		values := map[string]string{}
		for _, field := range extension.Extension {
			values[field.URL] = field.ValueString
		}

		return values
	}

	return nil
}

func TestUseCasesClinicalImpl_ListPatientConditions(t *testing.T) {
	tenantIDs := &dto.TenantIdentifiers{
		OrganizationID: uuid.NewString(),
		FacilityID:     uuid.NewString(),
	}

	id := uuid.NewString()
	patient := &domain.FHIRPatientRelayPayload{
		Resource: &domain.FHIRPatient{
			ID: &id,
		},
	}

	UUID := uuid.New().String()
	statusSystem := scalarutils.URI("http://terminology.hl7.org/CodeSystem/condition-clinical")
	status := "active"
	note := scalarutils.Markdown("Fever Fever")
	today := time.Now().Format(time.RFC3339)
	currentTime := (*scalarutils.DateTime)(&today)
	uri := scalarutils.URI("1234567")
	codingCode := "1234"
	categoryCode := "ENCOUNTER_DIAGNOSIS"
	pagedConditions := &domain.PagedFHIRCondition{
		Conditions: []domain.FHIRCondition{
			{
				ID:         &UUID,
				Text:       &domain.FHIRNarrative{},
				Identifier: []*domain.FHIRIdentifier{},
				ClinicalStatus: &domain.FHIRCodeableConcept{
					Coding: []*domain.FHIRCoding{
						{
							System:  &statusSystem,
							Code:    (*scalarutils.Code)(&status),
							Display: string(status),
						},
					},
					Text: string(status),
				},
				Code: &domain.FHIRCodeableConcept{
					Coding: []*domain.FHIRCoding{
						{
							System:  &uri,
							Code:    (*scalarutils.Code)(&codingCode),
							Display: "1234",
						},
					},
					Text: "1234",
				},
				OnsetDateTime: &scalarutils.Date{},
				RecordedDate:  &scalarutils.Date{},
				Note: []*domain.FHIRAnnotation{
					{
						Time: currentTime,
						Text: &note,
					},
				},
				Subject: &domain.FHIRReference{
					ID: &UUID,
				},
				Encounter: &domain.FHIRReference{
					ID: &UUID,
				},
				Category: []*domain.FHIRCodeableConcept{
					{
						ID: &UUID,
						Coding: []*domain.FHIRCoding{
							{
								ID:           &UUID,
								System:       (*scalarutils.URI)(&UUID),
								Version:      &UUID,
								Code:         (*scalarutils.Code)(&categoryCode),
								Display:      gofakeit.BeerAlcohol(),
								UserSelected: new(bool),
							},
						},
						Text: dto.ConditionCategoryDiagnosis.ToString(),
					},
				},
			},
		},
	}

	pagedConditions.Conditions = append(pagedConditions.Conditions, domain.FHIRCondition{
		ID: &UUID,
		ClinicalStatus: &domain.FHIRCodeableConcept{
			Text: status,
		},
		Code: &domain.FHIRCodeableConcept{
			Text: "Breast Cancer",
		},
		RecordedDate: &scalarutils.Date{},
		Subject: &domain.FHIRReference{
			ID: &UUID,
		},
		Encounter: &domain.FHIRReference{
			ID: &UUID,
		},
		Category: []*domain.FHIRCodeableConcept{
			{
				Text: dto.ConditionCategoryProblemList.ToString(),
			},
			{
				Text: "Primary Tumor Code - C50.9",
			},
			{
				Text: "Morphology Code - 8500/3",
			},
		},
		Stage: []*domain.FHIRConditionStage{
			{
				Summary: &domain.FHIRCodeableConcept{
					Text: "Stage 3",
				},
			},
		},
		// An unrelated extension should be ignored by the treatment-linkage extraction.
		Extension: []*domain.FHIRExtension{
			{
				URL: "http://savannahghi.org/fhir/StructureDefinition/some-other-extension",
			},
		},
	})

	first := 3
	encounterId := uuid.New().String()
	type args struct {
		ctx         context.Context
		patientID   string
		encounterID *string
		date        *scalarutils.Date
		strategy    string
		pagination  dto.Pagination
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "happy case: list conditions",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return patient, nil
					})
				mh.FHIR.EXPECT().SearchFHIRCondition(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRCondition, error) {
						return pagedConditions, nil
					})

				return args{ctx: context.Background(), patientID: gofakeit.UUID(), pagination: dto.Pagination{}}
			},
			wantErr: false,
		},
		{
			name: "happy case: list only retrospective conditions",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return patient, nil
					})
				mh.FHIR.EXPECT().SearchFHIRCondition(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRCondition, error) {
						// The retrospective filter must translate into a FHIR `_source` search parameter.
						source := "http://savannahghi.org/fhir/record-source/linkage"
						if params["_source"] != source {
							return nil, fmt.Errorf("expected _source filter %q, got %v", source, params["_source"])
						}

						return pagedConditions, nil
					})

				return args{ctx: context.Background(), patientID: gofakeit.UUID(), strategy: "linkage", pagination: dto.Pagination{}}
			},
			wantErr: false,
		},
		{
			name: "happy case: list conditions with encounterID",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return patient, nil
					})
				mh.FHIR.EXPECT().SearchFHIRCondition(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRCondition, error) {
						return pagedConditions, nil
					})

				return args{ctx: context.Background(), patientID: gofakeit.UUID(), encounterID: &encounterId, pagination: dto.Pagination{}}
			},
			wantErr: false,
		},
		{
			name: "happy case: list conditions with date",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return patient, nil
					})
				mh.FHIR.EXPECT().SearchFHIRCondition(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRCondition, error) {
						return pagedConditions, nil
					})

				date := &scalarutils.Date{
					Year:  2013,
					Month: 12,
					Day:   8,
				}
				return args{ctx: context.Background(), patientID: gofakeit.UUID(), pagination: dto.Pagination{}, date: date}
			},
			wantErr: false,
		},
		{
			name: "sad case: invalid patient id",
			setup: func(mh *usecaseMock.Mocks) args {
				return args{ctx: context.Background(), patientID: "invalid", pagination: dto.Pagination{}}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - invalid pagination",
			setup: func(mh *usecaseMock.Mocks) args {
				return args{ctx: context.Background(), patientID: gofakeit.UUID(), pagination: dto.Pagination{First: &first, Last: &first}}
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

				return args{ctx: context.Background(), patientID: gofakeit.UUID(), pagination: dto.Pagination{}}
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

				return args{ctx: context.Background(), patientID: gofakeit.UUID(), pagination: dto.Pagination{}}
			},
			wantErr: true,
		},
		{
			name: "sad case: fail to search condition",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return patient, nil
					})
				mh.FHIR.EXPECT().SearchFHIRCondition(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRCondition, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: context.Background(), patientID: gofakeit.UUID(), pagination: dto.Pagination{}}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			_, err := clinicalUsecase.ListPatientConditions(args.ctx, args.patientID, args.encounterID, args.date, args.strategy, args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListPatientConditions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

// TODO: Get clarification on how the test cases for missing grade and stage should be evaluated

// func TestClinicalImpl_RecordOncologicalDiagnosis(t *testing.T) {
// 	tenantIDs := &dto.TenantIdentifiers{
// 		OrganizationID: uuid.NewString(),
// 		FacilityID:     uuid.NewString(),
// 	}

// 	orgName := mock.Anything
// 	ID := uuid.NewString()
// 	organization := &domain.FHIROrganizationRelayPayload{
// 		Resource: &domain.FHIROrganization{
// 			ID:   &ID,
// 			Name: &orgName,
// 		},
// 	}

// 	patientRef := "Patient/" + uuid.NewString()
// 	encounter := &domain.FHIREncounterRelayPayload{
// 		Resource: &domain.FHIREncounter{
// 			ID:            &ID,
// 			Text:          &domain.FHIRNarrative{},
// 			Identifier:    []*domain.FHIRIdentifier{},
// 			Status:        domain.EncounterStatusEnum(domain.EncounterStatusEnumDischarged),
// 			StatusHistory: []*domain.FHIREncounterStatushistory{},
// 			Class:         []domain.FHIRCodeableConcept{},
// 			ClassHistory:  []*domain.FHIREncounterClasshistory{},
// 			Type:          []*domain.FHIRCodeableConcept{},
// 			ServiceType:   &domain.FHIRCodeableConcept{},
// 			Priority:      &domain.FHIRCodeableConcept{},
// 			Subject: &domain.FHIRReference{
// 				ID:        &ID,
// 				Reference: &patientRef,
// 			},
// 			EpisodeOfCare:   []*domain.FHIRReference{},
// 			BasedOn:         []*domain.FHIRReference{},
// 			Participant:     []*domain.FHIREncounterParticipant{},
// 			Appointment:     []*domain.FHIRReference{},
// 			ActualPeriod:    &domain.FHIRPeriod{},
// 			Length:          &domain.FHIRDuration{},
// 			ReasonReference: []*domain.FHIRReference{},
// 			Diagnosis:       []*domain.FHIREncounterDiagnosis{},
// 			Account:         []*domain.FHIRReference{},
// 			Hospitalization: &domain.FHIREncounterHospitalization{},
// 			Location:        []*domain.FHIREncounterLocation{},
// 			ServiceProvider: &domain.FHIRReference{},
// 			PartOf:          &domain.FHIRReference{},
// 		},
// 	}

// 	observation := &domain.FHIRObservation{
// 		ID:         &ID,
// 		Text:       &domain.FHIRNarrative{},
// 		Identifier: []*domain.FHIRIdentifier{},
// 	}

// 	statusSystem := scalarutils.URI("http://terminology.hl7.org/CodeSystem/condition-clinical")
// 	status := "active"
// 	note := scalarutils.Markdown("Fever Fever")
// 	today := time.Now().Format(time.RFC3339)
// 	currentTime := (*scalarutils.DateTime)(&today)
// 	uri := scalarutils.URI("1234567")
// 	codingCode := "1234"
// 	categoryCode := "ENCOUNTER_DIAGNOSIS"
// 	UUID := gofakeit.UUID()
// 	condition := &domain.FHIRConditionRelayPayload{
// 		Resource: &domain.FHIRCondition{
// 			ID:         &UUID,
// 			Text:       &domain.FHIRNarrative{},
// 			Identifier: []*domain.FHIRIdentifier{},
// 			ClinicalStatus: &domain.FHIRCodeableConcept{
// 				Coding: []*domain.FHIRCoding{
// 					{
// 						System:  &statusSystem,
// 						Code:    (*scalarutils.Code)(&status),
// 						Display: string(status),
// 					},
// 				},
// 				Text: string(status),
// 			},
// 			Code: &domain.FHIRCodeableConcept{
// 				Coding: []*domain.FHIRCoding{
// 					{
// 						System:  &uri,
// 						Code:    (*scalarutils.Code)(&codingCode),
// 						Display: "1234",
// 					},
// 				},
// 				Text: "1234",
// 			},
// 			OnsetDateTime: &scalarutils.Date{},
// 			RecordedDate:  &scalarutils.Date{},
// 			Note: []*domain.FHIRAnnotation{
// 				{
// 					Time: currentTime,
// 					Text: &note,
// 				},
// 			},
// 			Subject: &domain.FHIRReference{
// 				ID: &UUID,
// 			},
// 			Encounter: &domain.FHIRReference{
// 				ID: &UUID,
// 			},
// 			Category: []*domain.FHIRCodeableConcept{
// 				{
// 					ID: &UUID,
// 					Coding: []*domain.FHIRCoding{
// 						{
// 							ID:           &UUID,
// 							System:       (*scalarutils.URI)(&UUID),
// 							Version:      &UUID,
// 							Code:         (*scalarutils.Code)(&categoryCode),
// 							Display:      gofakeit.BeerAlcohol(),
// 							UserSelected: new(bool),
// 						},
// 					},
// 					Text: "ENCOUNTER_DIAGNOSIS",
// 				},
// 			},
// 		},
// 	}

// 	type args struct {
// 		ctx   context.Context
// 		input *dto.OncologyDiagnosisInput
// 	}
// 	tests := []struct {
// 		name    string
// 		setup   func(mh *usecaseMock.Mocks) args
// 		wantErr bool
// 	}{
// 		{
// 			name: "Happy case: create oncological diagnosis",
// 			setup: func(mh *usecaseMock.Mocks) args {
// 				ctx := usecaseMock.AddTenantIdentifierContext(context.Background())
// 				input := &dto.OncologyDiagnosisInput{
// 					EncounterID: gofakeit.UUID(),
// 					Condition: dto.ValueSetData{
// 						Code:    gofakeit.UUID(),
// 						Display: gofakeit.UUID(),
// 					},
// 					ICDO3PrimaryTumorCode: gofakeit.UUID(),
// 					ICDO3MorphologyCode:   gofakeit.UUID(),
// 					Behavior: dto.ValueSetData{
// 						Code:    gofakeit.UUID(),
// 						Display: gofakeit.UUID(),
// 					},
// 					Grade: dto.ValueSetData{
// 						Code:    gofakeit.UUID(),
// 						Display: gofakeit.UUID(),
// 					},
// 					Stage: dto.ValueSetData{
// 						Code:    gofakeit.UUID(),
// 						Display: gofakeit.UUID(),
// 					},
// 					Notes: gofakeit.UUID(),
// 				}

// 				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
// 					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
// 						return encounter, nil
// 					})
// 				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
// 					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
// 						return tenantIDs, nil
// 					})
// 				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
// 					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
// 						return organization, nil
// 					})
// 				mh.FHIR.EXPECT().CreateFHIRObservation(mock.Anything, mock.Anything).
// 					RunAndReturn(func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
// 						return observation, nil
// 					})
// 				mh.FHIR.EXPECT().CreateFHIRDiagnosticReport(mock.Anything, mock.Anything).
// 					RunAndReturn(func(context1 context.Context, input *domain.FHIRDiagnosticReportInput) (*domain.FHIRDiagnosticReport, error) {
// 						implicitRules := mock.Anything
// 						language := mock.Anything
// 						diagnosticsReport := &domain.FHIRDiagnosticReport{
// 							ID:            &ID,
// 							ImplicitRules: &implicitRules,
// 							Language:      &language,
// 						}
// 						return diagnosticsReport, nil
// 					})
// 				mh.FHIR.EXPECT().CreateFHIRCondition(mock.Anything, mock.Anything).
// 					RunAndReturn(func(ctx context.Context, input domain.FHIRConditionInput) (*domain.FHIRConditionRelayPayload, error) {
// 						return condition, nil
// 					})

// 				return args{ctx: ctx, input: input}
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "Happy case: unbound condition",
// 			setup: func(mh *usecaseMock.Mocks) args {
// 				ctx := usecaseMock.AddTenantIdentifierContext(context.Background())
// 				input := &dto.OncologyDiagnosisInput{
// 					EncounterID: gofakeit.UUID(),
// 					Condition: dto.ValueSetData{
// 						Display: gofakeit.UUID(),
// 					},
// 					ICDO3PrimaryTumorCode: gofakeit.UUID(),
// 					ICDO3MorphologyCode:   gofakeit.UUID(),
// 					Behavior: dto.ValueSetData{
// 						Code:    gofakeit.UUID(),
// 						Display: gofakeit.UUID(),
// 					},
// 					Grade: dto.ValueSetData{
// 						Code:    gofakeit.UUID(),
// 						Display: gofakeit.UUID(),
// 					},
// 					Stage: dto.ValueSetData{
// 						Code:    gofakeit.UUID(),
// 						Display: gofakeit.UUID(),
// 					},
// 					Notes: gofakeit.UUID(),
// 				}

// 				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
// 					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
// 						return encounter, nil
// 					})
// 				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
// 					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
// 						return tenantIDs, nil
// 					})
// 				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
// 					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
// 						return organization, nil
// 					})
// 				mh.FHIR.EXPECT().CreateFHIRObservation(mock.Anything, mock.Anything).
// 					RunAndReturn(func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
// 						return observation, nil
// 					})
// 				mh.FHIR.EXPECT().CreateFHIRDiagnosticReport(mock.Anything, mock.Anything).
// 					RunAndReturn(func(context1 context.Context, input *domain.FHIRDiagnosticReportInput) (*domain.FHIRDiagnosticReport, error) {
// 						implicitRules := mock.Anything
// 						language := mock.Anything
// 						diagnosticsReport := &domain.FHIRDiagnosticReport{
// 							ID:            &ID,
// 							ImplicitRules: &implicitRules,
// 							Language:      &language,
// 						}
// 						return diagnosticsReport, nil
// 					})
// 				mh.FHIR.EXPECT().CreateFHIRCondition(mock.Anything, mock.Anything).
// 					RunAndReturn(func(ctx context.Context, input domain.FHIRConditionInput) (*domain.FHIRConditionRelayPayload, error) {
// 						return condition, nil
// 					})

// 				return args{ctx: ctx, input: input}
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "Sad case: unable to get encounter",
// 			setup: func(mh *usecaseMock.Mocks) args {
// 				ctx := usecaseMock.AddTenantIdentifierContext(context.Background())
// 				input := &dto.OncologyDiagnosisInput{
// 					EncounterID: gofakeit.UUID(),
// 					Condition: dto.ValueSetData{
// 						Code:    gofakeit.UUID(),
// 						Display: gofakeit.UUID(),
// 					},
// 					ICDO3PrimaryTumorCode: gofakeit.UUID(),
// 					ICDO3MorphologyCode:   gofakeit.UUID(),
// 					Behavior: dto.ValueSetData{
// 						Code:    gofakeit.UUID(),
// 						Display: gofakeit.UUID(),
// 					},
// 					Grade: dto.ValueSetData{
// 						Code:    gofakeit.UUID(),
// 						Display: gofakeit.UUID(),
// 					},
// 					Stage: dto.ValueSetData{
// 						Code:    gofakeit.UUID(),
// 						Display: gofakeit.UUID(),
// 					},
// 					Notes: gofakeit.UUID(),
// 				}

// 				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
// 					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
// 						return nil, fmt.Errorf("an error occurred")
// 					})

// 				return args{ctx: ctx, input: input}
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "Sad case: missing grade",
// 			setup: func(mh *usecaseMock.Mocks) args {
// 				ctx := usecaseMock.AddTenantIdentifierContext(context.Background())
// 				input := &dto.OncologyDiagnosisInput{
// 					EncounterID: gofakeit.UUID(),
// 					Condition: dto.ValueSetData{
// 						Code:    gofakeit.UUID(),
// 						Display: gofakeit.UUID(),
// 					},
// 					ICDO3PrimaryTumorCode: gofakeit.UUID(),
// 					ICDO3MorphologyCode:   gofakeit.UUID(),
// 					Behavior: dto.ValueSetData{
// 						Code:    gofakeit.UUID(),
// 						Display: gofakeit.UUID(),
// 					},
// 					Stage: dto.ValueSetData{
// 						Code:    gofakeit.UUID(),
// 						Display: gofakeit.UUID(),
// 					},
// 					Notes: gofakeit.UUID(),
// 				}

// 				return args{ctx: ctx, input: input}
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "Sad case: missing stage",
// 			setup: func(mh *usecaseMock.Mocks) args {
// 				ctx := usecaseMock.AddTenantIdentifierContext(context.Background())
// 				input := &dto.OncologyDiagnosisInput{
// 					EncounterID: gofakeit.UUID(),
// 					Condition: dto.ValueSetData{
// 						Code:    gofakeit.UUID(),
// 						Display: gofakeit.UUID(),
// 					},
// 					ICDO3PrimaryTumorCode: gofakeit.UUID(),
// 					ICDO3MorphologyCode:   gofakeit.UUID(),
// 					Behavior: dto.ValueSetData{
// 						Code:    gofakeit.UUID(),
// 						Display: gofakeit.UUID(),
// 					},
// 					Grade: dto.ValueSetData{
// 						Code:    gofakeit.UUID(),
// 						Display: gofakeit.UUID(),
// 					},
// 					Notes: gofakeit.UUID(),
// 				}

// 				return args{ctx: ctx, input: input}
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "Sad case: unable to get meta tags",
// 			setup: func(mh *usecaseMock.Mocks) args {
// 				ctx := usecaseMock.AddTenantIdentifierContext(context.Background())
// 				input := &dto.OncologyDiagnosisInput{
// 					EncounterID: gofakeit.UUID(),
// 					Condition: dto.ValueSetData{
// 						Code:    gofakeit.UUID(),
// 						Display: gofakeit.UUID(),
// 					},
// 					ICDO3PrimaryTumorCode: gofakeit.UUID(),
// 					ICDO3MorphologyCode:   gofakeit.UUID(),
// 					Behavior: dto.ValueSetData{
// 						Code:    gofakeit.UUID(),
// 						Display: gofakeit.UUID(),
// 					},
// 					Grade: dto.ValueSetData{
// 						Code:    gofakeit.UUID(),
// 						Display: gofakeit.UUID(),
// 					},
// 					Stage: dto.ValueSetData{
// 						Code:    gofakeit.UUID(),
// 						Display: gofakeit.UUID(),
// 					},
// 					Notes: gofakeit.UUID(),
// 				}

// 				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
// 					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
// 						return encounter, nil
// 					})
// 				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
// 					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
// 						return tenantIDs, nil
// 					})
// 				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
// 					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
// 						return nil, fmt.Errorf("an error occurred")
// 					})

// 				return args{ctx: ctx, input: input}
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "Sad case: unable to get facility from context",
// 			setup: func(mh *usecaseMock.Mocks) args {
// 				ctx := context.Background()
// 				input := &dto.OncologyDiagnosisInput{
// 					EncounterID: gofakeit.UUID(),
// 					Condition: dto.ValueSetData{
// 						Code:    gofakeit.UUID(),
// 						Display: gofakeit.UUID(),
// 					},
// 					ICDO3PrimaryTumorCode: gofakeit.UUID(),
// 					ICDO3MorphologyCode:   gofakeit.UUID(),
// 					Behavior: dto.ValueSetData{
// 						Code:    gofakeit.UUID(),
// 						Display: gofakeit.UUID(),
// 					},
// 					Grade: dto.ValueSetData{
// 						Code:    gofakeit.UUID(),
// 						Display: gofakeit.UUID(),
// 					},
// 					Stage: dto.ValueSetData{
// 						Code:    gofakeit.UUID(),
// 						Display: gofakeit.UUID(),
// 					},
// 					Notes: gofakeit.UUID(),
// 				}

// 				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
// 					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
// 						return encounter, nil
// 					})
// 				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
// 					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
// 						return tenantIDs, nil
// 					})
// 				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
// 					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
// 						return organization, nil
// 					})

// 				return args{ctx: ctx, input: input}
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "Sad case: unable to create observation",
// 			setup: func(mh *usecaseMock.Mocks) args {
// 				ctx := usecaseMock.AddTenantIdentifierContext(context.Background())
// 				input := &dto.OncologyDiagnosisInput{
// 					EncounterID: gofakeit.UUID(),
// 					Condition: dto.ValueSetData{
// 						Code:    gofakeit.UUID(),
// 						Display: gofakeit.UUID(),
// 					},
// 					ICDO3PrimaryTumorCode: gofakeit.UUID(),
// 					ICDO3MorphologyCode:   gofakeit.UUID(),
// 					Behavior: dto.ValueSetData{
// 						Code:    gofakeit.UUID(),
// 						Display: gofakeit.UUID(),
// 					},
// 					Grade: dto.ValueSetData{
// 						Code:    gofakeit.UUID(),
// 						Display: gofakeit.UUID(),
// 					},
// 					Stage: dto.ValueSetData{
// 						Code:    gofakeit.UUID(),
// 						Display: gofakeit.UUID(),
// 					},
// 					Notes: gofakeit.UUID(),
// 				}

// 				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
// 					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
// 						return encounter, nil
// 					})
// 				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
// 					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
// 						return tenantIDs, nil
// 					})
// 				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
// 					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
// 						return organization, nil
// 					})
// 				mh.FHIR.EXPECT().CreateFHIRObservation(mock.Anything, mock.Anything).
// 					RunAndReturn(func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
// 						return nil, fmt.Errorf("an error occurred")
// 					})

// 				return args{ctx: ctx, input: input}
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "Sad case: unable to create diagnostic report",
// 			setup: func(mh *usecaseMock.Mocks) args {
// 				ctx := usecaseMock.AddTenantIdentifierContext(context.Background())
// 				input := &dto.OncologyDiagnosisInput{
// 					EncounterID: gofakeit.UUID(),
// 					Condition: dto.ValueSetData{
// 						Code:    gofakeit.UUID(),
// 						Display: gofakeit.UUID(),
// 					},
// 					ICDO3PrimaryTumorCode: gofakeit.UUID(),
// 					ICDO3MorphologyCode:   gofakeit.UUID(),
// 					Behavior: dto.ValueSetData{
// 						Code:    gofakeit.UUID(),
// 						Display: gofakeit.UUID(),
// 					},
// 					Grade: dto.ValueSetData{
// 						Code:    gofakeit.UUID(),
// 						Display: gofakeit.UUID(),
// 					},
// 					Stage: dto.ValueSetData{
// 						Code:    gofakeit.UUID(),
// 						Display: gofakeit.UUID(),
// 					},
// 					Notes: gofakeit.UUID(),
// 				}

// 				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
// 					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
// 						return encounter, nil
// 					})
// 				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
// 					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
// 						return tenantIDs, nil
// 					})
// 				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
// 					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
// 						return organization, nil
// 					})
// 				mh.FHIR.EXPECT().CreateFHIRObservation(mock.Anything, mock.Anything).
// 					RunAndReturn(func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
// 						return observation, nil
// 					})
// 				mh.FHIR.EXPECT().CreateFHIRDiagnosticReport(mock.Anything, mock.Anything).
// 					RunAndReturn(func(context1 context.Context, input *domain.FHIRDiagnosticReportInput) (*domain.FHIRDiagnosticReport, error) {
// 						return nil, fmt.Errorf("an error occurred")
// 					})

// 				return args{ctx: ctx, input: input}
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "Sad case: unable to create condition",
// 			setup: func(mh *usecaseMock.Mocks) args {
// 				ctx := usecaseMock.AddTenantIdentifierContext(context.Background())
// 				input := &dto.OncologyDiagnosisInput{
// 					EncounterID: gofakeit.UUID(),
// 					Condition: dto.ValueSetData{
// 						Code:    gofakeit.UUID(),
// 						Display: gofakeit.UUID(),
// 					},
// 					ICDO3PrimaryTumorCode: gofakeit.UUID(),
// 					ICDO3MorphologyCode:   gofakeit.UUID(),
// 					Behavior: dto.ValueSetData{
// 						Code:    gofakeit.UUID(),
// 						Display: gofakeit.UUID(),
// 					},
// 					Grade: dto.ValueSetData{
// 						Code:    gofakeit.UUID(),
// 						Display: gofakeit.UUID(),
// 					},
// 					Stage: dto.ValueSetData{
// 						Code:    gofakeit.UUID(),
// 						Display: gofakeit.UUID(),
// 					},
// 					Notes: gofakeit.UUID(),
// 				}

// 				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
// 					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
// 						return encounter, nil
// 					})
// 				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
// 					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
// 						return tenantIDs, nil
// 					})
// 				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
// 					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
// 						return organization, nil
// 					})
// 				mh.FHIR.EXPECT().CreateFHIRObservation(mock.Anything, mock.Anything).
// 					RunAndReturn(func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
// 						return observation, nil
// 					})
// 				mh.FHIR.EXPECT().CreateFHIRDiagnosticReport(mock.Anything, mock.Anything).
// 					RunAndReturn(func(context1 context.Context, input *domain.FHIRDiagnosticReportInput) (*domain.FHIRDiagnosticReport, error) {
// 						implicitRules := mock.Anything
// 						language := mock.Anything
// 						diagnosticsReport := &domain.FHIRDiagnosticReport{
// 							ID:            &ID,
// 							ImplicitRules: &implicitRules,
// 							Language:      &language,
// 						}
// 						return diagnosticsReport, nil
// 					})
// 				mh.FHIR.EXPECT().CreateFHIRCondition(mock.Anything, mock.Anything).
// 					RunAndReturn(func(ctx context.Context, input domain.FHIRConditionInput) (*domain.FHIRConditionRelayPayload, error) {
// 						return nil, fmt.Errorf("an error occurred")
// 					})

// 				return args{ctx: ctx, input: input}
// 			},
// 			wantErr: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Parallel()

// 			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
// 			args := tt.setup(&mock)

// 			_, err := clinicalUsecase.RecordOncologicalDiagnosis(args.ctx, args.input)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("ClinicalImpl.RecordOncologicalDiagnosis() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 		})
// 	}
// }
