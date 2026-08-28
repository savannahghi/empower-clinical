package clinical_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/scalarutils"
	"github.com/stretchr/testify/mock"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/empower-clinical/pkg/clinical/domain"
	usecaseMock "github.com/savannahghi/empower-clinical/pkg/clinical/usecases/mock"
)

func TestUseCasesClinicalImpl_CreateAllergyIntolerance(t *testing.T) {
	id := uuid.NewString()
	patientRef := "Patient/" + uuid.NewString()
	encounter := &domain.FHIREncounterRelayPayload{
		Resource: &domain.FHIREncounter{
			ID:            &id,
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
				ID:        &id,
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
	ID := uuid.NewString()
	organization := &domain.FHIROrganizationRelayPayload{
		Resource: &domain.FHIROrganization{
			ID:   &ID,
			Name: &orgName,
		},
	}

	type args struct {
		ctx   context.Context
		input dto.AllergyInput
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy case: create allergy intolerance",
			setup: func(mh *usecaseMock.Mocks) args {
				input := dto.AllergyInput{
					PatientID:         gofakeit.UUID(),
					Code:              "C12345",
					TerminologySource: domain.TerminologySourceCIEL,
					EncounterID:       gofakeit.UUID(),
					Reaction: &dto.ReactionInput{
						Code:     "2000",
						System:   gofakeit.BS(),
						Severity: "fatal",
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
				mh.FHIR.EXPECT().CreateFHIRAllergyIntolerance(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRAllergyIntoleranceInput) (*domain.FHIRAllergyIntoleranceRelayPayload, error) {
						codingCode := "20"
						manifestationCodingCode := scalarutils.Code(gofakeit.BS())
						system := gofakeit.URL()
						UUID := gofakeit.UUID()
						mildSeverity := domain.AllergyIntoleranceReactionSeverityEnumMild
						allergyIntolerance := &domain.FHIRAllergyIntoleranceRelayPayload{
							Resource: &domain.FHIRAllergyIntolerance{
								ID:          &UUID,
								Criticality: "fatal",
								Code: &domain.FHIRCodeableConcept{
									Coding: []*domain.FHIRCoding{{
										System:  (*scalarutils.URI)(&system),
										Code:    (*scalarutils.Code)(&codingCode),
										Display: "Example display",
									}},
								},
								OnsetPeriod: &domain.FHIRPeriod{
									Start: scalarutils.DateTime("2000-01-01T00:00:00"),
								},
								Patient: &domain.FHIRReference{
									ID: &UUID,
								},
								Encounter: &domain.FHIRReference{
									ID: &UUID,
								},
								Reaction: []*domain.FHIRAllergyintoleranceReaction{
									{
										Severity: &mildSeverity,
										Manifestation: []*domain.FHIRCodeableReference{
											{
												Concept: &domain.FHIRCodeableConcept{
													Coding: []*domain.FHIRCoding{
														{
															System: (*scalarutils.URI)(&system),
															Code:   &manifestationCodingCode,
														},
													},
												},
											},
										},
									},
								},
							},
						}

						return allergyIntolerance, nil
					})

				return args{ctx: context.Background(), input: input}
			},
			wantErr: false,
		},
		{
			name: "Happy case: create allergy intolerance, no reaction",
			setup: func(mh *usecaseMock.Mocks) args {
				input := dto.AllergyInput{
					PatientID:         gofakeit.UUID(),
					Code:              "100",
					TerminologySource: domain.TerminologySourceLOINC,
					EncounterID:       gofakeit.UUID(),
				}

				system := gofakeit.URL()
				display := gofakeit.FirstName()
				UUID := gofakeit.UUID()
				codingCode := "20"
				allergyIntolerance := &domain.FHIRAllergyIntoleranceRelayPayload{
					Resource: &domain.FHIRAllergyIntolerance{
						ID:          &UUID,
						Criticality: "let",
						Code: &domain.FHIRCodeableConcept{
							Coding: []*domain.FHIRCoding{{
								System:  (*scalarutils.URI)(&system),
								Code:    (*scalarutils.Code)(&codingCode),
								Display: display,
							}},
						},
						OnsetPeriod: &domain.FHIRPeriod{
							Start: scalarutils.DateTime("2000-01-01T00:00:00"),
						},
						Patient: &domain.FHIRReference{
							ID: &UUID,
						},
						Encounter: &domain.FHIRReference{
							ID: &UUID,
						},
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
				mh.FHIR.EXPECT().CreateFHIRAllergyIntolerance(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRAllergyIntoleranceInput) (*domain.FHIRAllergyIntoleranceRelayPayload, error) {
						return allergyIntolerance, nil
					})

				return args{ctx: context.Background(), input: input}
			},
			wantErr: false,
		},

		{
			name: "Sad case: fail to get encounter",
			setup: func(mh *usecaseMock.Mocks) args {
				input := dto.AllergyInput{
					PatientID:         gofakeit.UUID(),
					Code:              "100",
					TerminologySource: domain.TerminologySource("invalid"),
					EncounterID:       gofakeit.UUID(),
				}

				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: context.Background(), input: input}
			},
			wantErr: true,
		},

		{
			name: "Sad case: unsupported concept source",
			setup: func(mh *usecaseMock.Mocks) args {
				input := dto.AllergyInput{
					PatientID:         gofakeit.UUID(),
					Code:              "100",
					TerminologySource: domain.TerminologySource("invalid"),
					EncounterID:       gofakeit.UUID(),
				}

				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})

				return args{ctx: context.Background(), input: input}
			},
			wantErr: true,
		},

		{
			name: "Sad case: failed to create allergy intolerance",
			setup: func(mh *usecaseMock.Mocks) args {
				input := dto.AllergyInput{
					PatientID:         gofakeit.UUID(),
					Code:              "100",
					TerminologySource: domain.TerminologySourceCIEL,
					EncounterID:       gofakeit.UUID(),
					Reaction: &dto.ReactionInput{
						Code:     "2000",
						System:   gofakeit.BS(),
						Severity: "fatal",
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
				mh.FHIR.EXPECT().CreateFHIRAllergyIntolerance(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRAllergyIntoleranceInput) (*domain.FHIRAllergyIntoleranceRelayPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: context.Background(), input: input}
			},
			wantErr: true,
		},

		{
			name: "Sad case: no encounter id passed",
			setup: func(mh *usecaseMock.Mocks) args {
				input := dto.AllergyInput{
					PatientID:         gofakeit.UUID(),
					Code:              "100",
					TerminologySource: domain.TerminologySourceCIEL,
					Reaction: &dto.ReactionInput{
						Code:     "2000",
						System:   gofakeit.BS(),
						Severity: "fatal",
					},
				}

				return args{ctx: context.Background(), input: input}
			},
			wantErr: true,
		},

		{
			name: "Sad case: failed to get fhir encounter",
			setup: func(mh *usecaseMock.Mocks) args {
				input := dto.AllergyInput{
					PatientID:         gofakeit.UUID(),
					Code:              "100",
					TerminologySource: domain.TerminologySourceCIEL,
					EncounterID:       gofakeit.UUID(),
					Reaction: &dto.ReactionInput{
						Code:     "2000",
						System:   gofakeit.BS(),
						Severity: "fatal",
					},
				}

				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: context.Background(), input: input}
			},
			wantErr: true,
		},

		{
			name: "Sad case - fail on finished encounter",
			setup: func(mh *usecaseMock.Mocks) args {
				input := dto.AllergyInput{
					PatientID:         gofakeit.UUID(),
					Code:              "100",
					TerminologySource: domain.TerminologySourceCIEL,
					EncounterID:       gofakeit.UUID(),
					Reaction: &dto.ReactionInput{
						Code:     "2000",
						System:   gofakeit.BS(),
						Severity: "fatal",
					},
				}

				encounter := &domain.FHIREncounterRelayPayload{
					Resource: &domain.FHIREncounter{
						ID:            &id,
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
							ID:        &id,
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

				return args{ctx: context.Background(), input: input}
			},
			wantErr: true,
		},
		{
			name: "Sad case: fail to get tags",
			setup: func(mh *usecaseMock.Mocks) args {
				input := dto.AllergyInput{
					PatientID:         gofakeit.UUID(),
					Code:              "100",
					TerminologySource: domain.TerminologySourceCIEL,
					EncounterID:       gofakeit.UUID(),
					Reaction: &dto.ReactionInput{
						Code:     "2000",
						System:   gofakeit.BS(),
						Severity: "fatal",
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
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: context.Background(), input: input}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - fail to get ciel concept",
			setup: func(mh *usecaseMock.Mocks) args {
				input := dto.AllergyInput{
					PatientID:         gofakeit.UUID(),
					Code:              "100",
					TerminologySource: domain.TerminologySourceCIEL,
					EncounterID:       gofakeit.UUID(),
					Reaction: &dto.ReactionInput{
						Code:     "2000",
						System:   gofakeit.BS(),
						Severity: "fatal",
					},
				}

				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: context.Background(), input: input}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			_, err := clinicalUsecase.CreateAllergyIntolerance(args.ctx, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.CreateAllergyIntolerance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCasesClinicalImpl_SearchAllergy(t *testing.T) {
	first := 5
	type args struct {
		ctx        context.Context
		name       string
		pagination dto.Pagination
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy case: search for an allergy",
			setup: func(mh *usecaseMock.Mocks) args {
				next := "100"
				previous := "80"
				conceptPage := &domain.ConceptPage{
					Results: []*domain.Concept{
						{
							ID:          uuid.NewString(),
							Source:      "ICD10",
							DisplayName: mock.Anything,
						},
					},
					Next:     &next,
					Previous: &previous,
				}

				mh.OCL.EXPECT().ListConcepts(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything,
					mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, org, source []string, verbose bool, q, sortAsc, sortDesc, conceptClass, dataType, locale *string, includeRetired, includeMappings, includeInverseMappings *bool, paginationInput *dto.Pagination) (*domain.ConceptPage, error) {
						return conceptPage, nil
					})

				return args{ctx: context.Background(), name: "Peanuts", pagination: dto.Pagination{First: &first}}
			},
			wantErr: false,
		},
		{
			name: "Sad case: invalid pagination",
			setup: func(mh *usecaseMock.Mocks) args {
				return args{ctx: context.Background(), name: "Peanuts", pagination: dto.Pagination{First: &first, Last: &first}}
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to search for an allergy",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.OCL.EXPECT().ListConcepts(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything,
					mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, org, source []string, verbose bool, q, sortAsc, sortDesc, conceptClass, dataType, locale *string, includeRetired, includeMappings, includeInverseMappings *bool, paginationInput *dto.Pagination) (*domain.ConceptPage, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: context.Background(), name: "Peanuts"}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			_, err := clinicalUsecase.SearchAllergy(args.ctx, args.name, args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.SearchAllergy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCasesClinicalImpl_GetAllergyIntolerance(t *testing.T) {
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
			name: "Happy case: get allergy intolerance",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRAllergyIntolerance(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRAllergyIntoleranceRelayPayload, error) {
						codingCode := "20"
						display := gofakeit.FirstName()
						manifestationCodingCode := scalarutils.Code(gofakeit.BS())
						system := gofakeit.URL()
						UUID := gofakeit.UUID()
						mildSeverity := domain.AllergyIntoleranceReactionSeverityEnumMild
						allergyIntolerance := &domain.FHIRAllergyIntoleranceRelayPayload{
							Resource: &domain.FHIRAllergyIntolerance{
								ID:          &UUID,
								Criticality: "fatal",
								Code: &domain.FHIRCodeableConcept{
									Coding: []*domain.FHIRCoding{{
										System:  (*scalarutils.URI)(&system),
										Code:    (*scalarutils.Code)(&codingCode),
										Display: display,
									}},
								},
								OnsetPeriod: &domain.FHIRPeriod{
									Start: scalarutils.DateTime("2000-01-01T00:00:00"),
								},
								Patient: &domain.FHIRReference{
									ID: &UUID,
								},
								Encounter: &domain.FHIRReference{
									ID: &UUID,
								},
								Reaction: []*domain.FHIRAllergyintoleranceReaction{
									{
										Severity: &mildSeverity,
										Manifestation: []*domain.FHIRCodeableReference{
											{
												Concept: &domain.FHIRCodeableConcept{
													Coding: []*domain.FHIRCoding{
														{
															System:  (*scalarutils.URI)(&system),
															Code:    &manifestationCodingCode,
															Display: display,
														},
													},
												},
											},
										},
									},
								},
							},
						}

						return allergyIntolerance, nil
					})

				return args{ctx: context.Background(), id: gofakeit.UUID()}
			},
			wantErr: false,
		},
		{
			name: "Sad case: invalid uuid",
			setup: func(mh *usecaseMock.Mocks) args {
				return args{ctx: context.Background(), id: "12"}
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get allergy intolerance",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRAllergyIntolerance(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRAllergyIntoleranceRelayPayload, error) {

						return nil, fmt.Errorf("an error occurred")
					})
				return args{ctx: context.Background(), id: gofakeit.UUID()}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			_, err := clinicalUsecase.GetAllergyIntolerance(args.ctx, args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.GetAllergyIntolerance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCasesClinicalImpl_GetPatientAllergies(t *testing.T) {
	tenantIDs := &dto.TenantIdentifiers{
		OrganizationID: uuid.NewString(),
		FacilityID:     uuid.NewString(),
	}

	first := 3
	type args struct {
		ctx        context.Context
		patientID  string
		pagination dto.Pagination
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy case: get allergy patient intolerances",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchPatientAllergyIntolerance(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, patientReference string, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRAllergy, error) {
						id := uuid.NewString()
						ref := "Reference/123345"
						code := scalarutils.Code("ABS")
						severity := domain.AllergyIntoleranceReactionSeverityEnumMild
						pagedAllergyIntolerance := &domain.PagedFHIRAllergy{
							Allergies: []domain.FHIRAllergyIntolerance{
								{
									ID:         &id,
									Text:       &domain.FHIRNarrative{},
									Identifier: []*domain.FHIRIdentifier{},
									Patient: &domain.FHIRReference{
										ID:        &id,
										Reference: &ref,
									},
									Code: &domain.FHIRCodeableConcept{
										ID:   &id,
										Text: gofakeit.Paragraph(1, 2, 100, ","),
										Coding: []*domain.FHIRCoding{
											{
												Display: mock.Anything,
												Code:    &code,
												System:  &system,
											},
										},
									},
									Encounter: &domain.FHIRReference{
										ID: &id,
									},
									OnsetPeriod: &domain.FHIRPeriod{
										Start: scalarutils.DateTime("2025-05-05"),
									},
									Reaction: []*domain.FHIRAllergyintoleranceReaction{
										{
											Severity: &severity,
											Manifestation: []*domain.FHIRCodeableReference{
												{
													ID: &id,
													Concept: &domain.FHIRCodeableConcept{
														ID:   &id,
														Text: gofakeit.Paragraph(1, 2, 100, ","),
														Coding: []*domain.FHIRCoding{
															{
																Display: mock.Anything,
																Code:    &code,
																System:  &system,
															},
														},
													},
													Reference: &domain.FHIRReference{
														ID:        &id,
														Reference: &ref,
													},
												},
											},
										},
									},
								},
							},
							HasNextPage:     false,
							HasPreviousPage: false,
							NextCursor:      "",
							PreviousCursor:  "",
							TotalCount:      1,
						}
						return pagedAllergyIntolerance, nil
					})

				return args{ctx: context.Background(), patientID: uuid.NewString(), pagination: dto.Pagination{First: &first}}
			},
			wantErr: false,
		},
		{
			name: "Sad case: invalid uuid",
			setup: func(mh *usecaseMock.Mocks) args {
				return args{ctx: context.Background(), patientID: "1", pagination: dto.Pagination{First: &first}}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - invalid pagination",
			setup: func(mh *usecaseMock.Mocks) args {
				return args{ctx: context.Background(), patientID: uuid.NewString(), pagination: dto.Pagination{First: &first, Last: &first}}
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get patient allergy intolerances",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchPatientAllergyIntolerance(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, patientReference string, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRAllergy, error) {

						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: context.Background(), patientID: uuid.NewString(), pagination: dto.Pagination{First: &first}}
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get meta tags",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return nil, fmt.Errorf("an error occurred")
					})
				return args{ctx: context.Background(), patientID: uuid.NewString(), pagination: dto.Pagination{First: &first}}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			_, err := clinicalUsecase.ListPatientAllergies(args.ctx, args.patientID, args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.GetPatientAllergies() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
