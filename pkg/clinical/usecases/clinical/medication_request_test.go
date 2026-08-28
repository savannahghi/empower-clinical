package clinical_test

import (
	"context"
	"encoding/json"
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

func TestUseCasesClinicalImpl_ListMedicationRequests(t *testing.T) {
	ctx := usecaseMock.AddTenantIdentifierContext(context.Background())
	filter := &dto.MedicationRequestFilterInput{
		FilterInput: dto.FilterInput{
			PatientID:   gofakeit.UUID(),
			EncounterID: gofakeit.UUID(),
			Date: &scalarutils.Date{
				Year:  2023,
				Month: 10,
				Day:   4,
			},
		},
	}

	firstTen := 10
	pagination := dto.Pagination{
		First: &firstTen,
	}

	tenantIDs := &dto.TenantIdentifiers{
		OrganizationID: uuid.NewString(),
		FacilityID:     uuid.NewString(),
	}

	type args struct {
		ctx        context.Context
		filter     *dto.MedicationRequestFilterInput
		pagination dto.Pagination
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully get medication request",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchFHIRMedicationRequest(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRMedicationRequest, error) {
						status := domain.ActiveMedicationStatus
						id := uuid.NewString()
						authoredOn := scalarutils.DateTime(time.Now().GoString())
						priority := scalarutils.Code("asap")
						noteText := scalarutils.Markdown(mock.Anything)
						periodUnit := domain.UnitsOfTimeEnumD
						frequency := 3
						patientInstruction := gofakeit.BeerStyle()
						asneeded := true
						number := json.Number("12")
						system := scalarutils.URI("http://mycarehub/tenant-identification/facility")

						medicationRequestOutput := &domain.FHIRMedicationRequest{
							ID:         &id,
							Status:     &status,
							AuthoredOn: &authoredOn,
							Priority:   &priority,
							Subject: &domain.FHIRReference{
								Display: uuid.NewString(),
							},
							Encounter: &domain.FHIRReference{
								Display: uuid.NewString(),
							},
							Note: []*domain.FHIRAnnotation{
								{
									ID:   &id,
									Text: &noteText,
								},
							},
							Meta: &domain.FHIRMeta{
								VersionID:   gofakeit.UUID(),
								LastUpdated: time.Now(),
								Tag: []domain.FHIRCoding{
									{
										System:  &system,
										Display: gofakeit.BeerName(),
									},
								},
							},
							DosageInstruction: []*domain.FHIRDosage{
								{
									ID: &id,
									Route: &domain.FHIRCodeableConcept{
										Text: gofakeit.BeerName(),
									},
									DoseAndRate: []*domain.FHIRDosageDoseandrate{
										{
											DoseQuantity: &domain.FHIRQuantity{
												Value: 12.3,
												Unit:  "ml",
											},
										},
									},
									Timing: &domain.FHIRTiming{
										Repeat: &domain.FHIRTimingRepeat{
											Period:       &number,
											PeriodUnit:   &periodUnit,
											Frequency:    &frequency,
											Duration:     &number,
											DurationUnit: &periodUnit,
											BoundsPeriod: &domain.FHIRPeriod{
												Start: scalarutils.DateTime(time.Now().GoString()),
												End:   scalarutils.DateTime(time.Now().GoString()),
											},
										},
									},
									PatientInstruction: &patientInstruction,
									AsNeeded:           &asneeded,
								},
							},
						}
						pagedMedicationRequest := &domain.PagedFHIRMedicationRequest{
							MedicationRequests: []domain.FHIRMedicationRequest{
								*medicationRequestOutput,
							},
							HasNextPage:     false,
							HasPreviousPage: false,
							NextCursor:      "",
							PreviousCursor:  "",
							TotalCount:      1,
						}
						return pagedMedicationRequest, nil
					})

				return args{ctx: ctx, filter: filter, pagination: pagination}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - invalid patient id",
			setup: func(mh *usecaseMock.Mocks) args {
				filter := &dto.MedicationRequestFilterInput{
					FilterInput: dto.FilterInput{
						PatientID:   "id",
						EncounterID: gofakeit.UUID(),
						Date: &scalarutils.Date{
							Year:  2023,
							Month: 10,
							Day:   4,
						},
					},
				}
				return args{ctx: ctx, filter: filter, pagination: pagination}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - invalid encounter id",
			setup: func(mh *usecaseMock.Mocks) args {
				filter := &dto.MedicationRequestFilterInput{
					FilterInput: dto.FilterInput{
						PatientID:   gofakeit.UUID(),
						EncounterID: "id",
						Date: &scalarutils.Date{
							Year:  2023,
							Month: 10,
							Day:   4,
						},
					},
				}
				return args{ctx: ctx, filter: filter, pagination: pagination}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get tenant identifiers",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, filter: filter, pagination: pagination}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to search for medication request",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchFHIRMedicationRequest(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRMedicationRequest, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, filter: filter, pagination: pagination}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			_, err := clinicalUsecase.ListMedicationRequests(args.ctx, args.filter, args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.ListMedicationRequests() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCasesClinicalImpl_CreatePrescription(t *testing.T) {
	ctx := context.Background()

	periodUnit := domain.UnitsOfTimeEnumD
	durationUnit := domain.UnitsOfTimeEnumD
	date := scalarutils.DateTime(time.Wednesday.String())
	dosageInstruction := dto.DosageInstruction{
		Route: dto.ValueSetData{
			Code:    "ro",
			Display: "Oral",
		},
		DoseQuantity:          2,
		DoseUnit:              "Capsules",
		Period:                "8",
		PeriodUnit:            &periodUnit,
		Frequency:             1,
		Duration:              "5",
		DurationUnit:          &durationUnit,
		StartDate:             &date,
		EndDate:               &date,
		Condition:             "After meals..",
		PatientInstruction:    "",
		AdditionalInstruction: []string{},
		AsNeeded:              false,
	}

	input := dto.PrescriptionInput{
		EncounterID: uuid.New().String(),
		Medications: []dto.PrescriptionMedicationInput{
			{
				MedicationID: gofakeit.Name(),
				DosageInstructions: []dto.DosageInstruction{
					dosageInstruction,
				},
			},
		},
	}

	medicationStatus := domain.MedicationStatusEnumActive
	id := uuid.NewString()
	medication := &domain.FHIRMedication{
		ID:     &id,
		Status: &medicationStatus,
	}

	ref := "Reference/12345"
	encounter := &domain.FHIREncounterRelayPayload{
		Resource: &domain.FHIREncounter{
			ID:     &id,
			Status: domain.EncounterStatusEnumInProgress,
			Subject: &domain.FHIRReference{
				ID: &id,
			},
			ServiceProvider: &domain.FHIRReference{
				Display:   gofakeit.UUID(),
				Reference: &ref,
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
			ID:   &id,
			Name: &orgName,
		},
	}

	status := domain.ActiveMedicationStatus
	authoredOn := scalarutils.DateTime(time.Now().GoString())
	priority := scalarutils.Code("asap")
	noteText := scalarutils.Markdown(mock.Anything)
	frequency := 3
	patientInstruction := gofakeit.BeerStyle()
	asneeded := true
	number := json.Number("12")
	system := scalarutils.URI("http://mycarehub/tenant-identification/facility")

	medicationRequestOutput := &domain.FHIRMedicationRequestRelayPayload{
		Resource: &domain.FHIRMedicationRequest{
			ID:         &id,
			Status:     &status,
			AuthoredOn: &authoredOn,
			Priority:   &priority,
			Subject: &domain.FHIRReference{
				Display: uuid.NewString(),
			},
			Encounter: &domain.FHIRReference{
				Display: uuid.NewString(),
			},
			Note: []*domain.FHIRAnnotation{
				{
					ID:   &id,
					Text: &noteText,
				},
			},
			Meta: &domain.FHIRMeta{
				VersionID:   gofakeit.UUID(),
				LastUpdated: time.Now(),
				Tag: []domain.FHIRCoding{
					{
						System:  &system,
						Display: gofakeit.BeerName(),
					},
				},
			},
			DosageInstruction: []*domain.FHIRDosage{
				{
					ID: &id,
					Route: &domain.FHIRCodeableConcept{
						Text: gofakeit.BeerName(),
					},
					DoseAndRate: []*domain.FHIRDosageDoseandrate{
						{
							DoseQuantity: &domain.FHIRQuantity{
								Value: 12.3,
								Unit:  "ml",
							},
						},
					},
					Timing: &domain.FHIRTiming{
						Repeat: &domain.FHIRTimingRepeat{
							Period:       &number,
							PeriodUnit:   &periodUnit,
							Frequency:    &frequency,
							Duration:     &number,
							DurationUnit: &periodUnit,
							BoundsPeriod: &domain.FHIRPeriod{
								Start: scalarutils.DateTime(time.Now().GoString()),
								End:   scalarutils.DateTime(time.Now().GoString()),
							},
						},
					},
					PatientInstruction: &patientInstruction,
					AsNeeded:           &asneeded,
				},
			},
		},
	}

	type args struct {
		ctx   context.Context
		input dto.PrescriptionInput
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy Case: Successfully create a prescription",
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
				mh.FHIR.EXPECT().FetchMedicationByID(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRMedication, error) {
						return medication, nil
					})
				mh.FHIR.EXPECT().CreateFHIRMedicationRequest(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRMedicationRequestInput) (*domain.FHIRMedicationRequestRelayPayload, error) {
						return medicationRequestOutput, nil
					})

				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "Sad Case: Empty encounter ID",
			setup: func(mh *usecaseMock.Mocks) args {
				input := dto.PrescriptionInput{
					EncounterID: "",
					Medications: []dto.PrescriptionMedicationInput{
						{
							MedicationID: gofakeit.Name(),
							DosageInstructions: []dto.DosageInstruction{
								dosageInstruction,
							},
						},
					},
				}
				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
		{
			name: "Sad Case: unable to fetch medication by ID",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.FHIR.EXPECT().FetchMedicationByID(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRMedication, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
		{
			name: "Sad Case: cannot create prescription on completed encounter",
			setup: func(mh *usecaseMock.Mocks) args {
				encounter := &domain.FHIREncounterRelayPayload{
					Resource: &domain.FHIREncounter{
						ID:     &id,
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
			name: "Sad Case: Fail to get encounter",
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
			name: "Sad Case: Fail to create medication request",
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
				mh.FHIR.EXPECT().FetchMedicationByID(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRMedication, error) {
						return medication, nil
					})
				mh.FHIR.EXPECT().CreateFHIRMedicationRequest(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRMedicationRequestInput) (*domain.FHIRMedicationRequestRelayPayload, error) {
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

			got, err := clinicalUsecase.CreatePrescription(args.ctx, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.CreatePrescription() error = %v, wantErr %v", err, tt.wantErr)
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

func TestUseCasesClinicalImpl_PatchMedicationRequests(t *testing.T) {
	ctx := usecaseMock.AddTenantIdentifierContext(context.Background())
	ID := uuid.NewString()
	value := domain.CompletedMedicationStatus

	type args struct {
		ctx   context.Context
		id    string
		value domain.MedicationRequestStatus
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy case: patch medication requests",
			setup: func(mh *usecaseMock.Mocks) args {
				status := domain.ActiveMedicationStatus
				id := uuid.NewString()
				authoredOn := scalarutils.DateTime(time.Now().GoString())
				priority := scalarutils.Code("asap")
				noteText := scalarutils.Markdown(mock.Anything)
				periodUnit := domain.UnitsOfTimeEnumD
				frequency := 3
				patientInstruction := gofakeit.BeerStyle()
				asneeded := true
				number := json.Number("12")
				system := scalarutils.URI("http://mycarehub/tenant-identification/facility")

				medicationRequest := &domain.FHIRMedicationRequest{
					ID:         &id,
					Status:     &status,
					AuthoredOn: &authoredOn,
					Priority:   &priority,
					Subject: &domain.FHIRReference{
						Display: uuid.NewString(),
					},
					Encounter: &domain.FHIRReference{
						Display: uuid.NewString(),
					},
					Note: []*domain.FHIRAnnotation{
						{
							ID:   &id,
							Text: &noteText,
						},
					},
					Meta: &domain.FHIRMeta{
						VersionID:   gofakeit.UUID(),
						LastUpdated: time.Now(),
						Tag: []domain.FHIRCoding{
							{
								System:  &system,
								Display: gofakeit.BeerName(),
							},
						},
					},
					DosageInstruction: []*domain.FHIRDosage{
						{
							ID: &id,
							Route: &domain.FHIRCodeableConcept{
								Text: gofakeit.BeerName(),
							},
							DoseAndRate: []*domain.FHIRDosageDoseandrate{
								{
									DoseQuantity: &domain.FHIRQuantity{
										Value: 12.3,
										Unit:  "ml",
									},
								},
							},
							Timing: &domain.FHIRTiming{
								Repeat: &domain.FHIRTimingRepeat{
									Period:       &number,
									PeriodUnit:   &periodUnit,
									Frequency:    &frequency,
									Duration:     &number,
									DurationUnit: &periodUnit,
									BoundsPeriod: &domain.FHIRPeriod{
										Start: scalarutils.DateTime(time.Now().GoString()),
										End:   scalarutils.DateTime(time.Now().GoString()),
									},
								},
							},
							PatientInstruction: &patientInstruction,
							AsNeeded:           &asneeded,
						},
					},
				}

				mh.FHIR.EXPECT().PatchFHIRMedicationRequest(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string, input domain.FHIRMedicationRequestInput) (*domain.FHIRMedicationRequest, error) {
						return medicationRequest, nil
					})
				return args{ctx: ctx, id: ID, value: value}
			},
			wantErr: false,
		},
		{
			name: "Sad case: invalid status",
			setup: func(mh *usecaseMock.Mocks) args {
				return args{ctx: ctx, id: ID, value: "invalid"}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail patch medication request",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().PatchFHIRMedicationRequest(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string, input domain.FHIRMedicationRequestInput) (*domain.FHIRMedicationRequest, error) {
						return nil, fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, id: ID, value: value}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - No medication request ID provided",
			setup: func(mh *usecaseMock.Mocks) args {
				return args{ctx: ctx, value: value}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Invalid medication request ID provided",
			setup: func(mh *usecaseMock.Mocks) args {
				return args{ctx: ctx, id: "invalid", value: value}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			_, err := clinicalUsecase.PatchMedicationRequests(args.ctx, args.id, args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.PatchMedicationRequests() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCasesClinicalImpl_FetchMedicationRequestByID(t *testing.T) {
	ctx := usecaseMock.AddTenantIdentifierContext(context.Background())
	medicationRequestID := uuid.NewString()

	type args struct {
		ctx                 context.Context
		medicationRequestID string
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy case: fetch medication request by ID",
			setup: func(mh *usecaseMock.Mocks) args {
				status := domain.ActiveMedicationStatus
				id := uuid.NewString()
				authoredOn := scalarutils.DateTime(time.Now().GoString())
				priority := scalarutils.Code("asap")
				noteText := scalarutils.Markdown(mock.Anything)
				periodUnit := domain.UnitsOfTimeEnumD
				frequency := 3
				patientInstruction := gofakeit.BeerStyle()
				asneeded := true
				number := json.Number("12")
				system := scalarutils.URI("http://mycarehub/tenant-identification/facility")

				medicationRequestOutput := &domain.FHIRMedicationRequestRelayPayload{
					Resource: &domain.FHIRMedicationRequest{
						ID:         &id,
						Status:     &status,
						AuthoredOn: &authoredOn,
						Priority:   &priority,
						Subject: &domain.FHIRReference{
							Display: uuid.NewString(),
						},
						Encounter: &domain.FHIRReference{
							Display: uuid.NewString(),
						},
						Note: []*domain.FHIRAnnotation{
							{
								ID:   &id,
								Text: &noteText,
							},
						},
						Meta: &domain.FHIRMeta{
							VersionID:   gofakeit.UUID(),
							LastUpdated: time.Now(),
							Tag: []domain.FHIRCoding{
								{
									System:  &system,
									Display: gofakeit.BeerName(),
								},
							},
						},
						DosageInstruction: []*domain.FHIRDosage{
							{
								ID: &id,
								Route: &domain.FHIRCodeableConcept{
									Text: gofakeit.BeerName(),
								},
								DoseAndRate: []*domain.FHIRDosageDoseandrate{
									{
										DoseQuantity: &domain.FHIRQuantity{
											Value: 12.3,
											Unit:  "ml",
										},
									},
								},
								Timing: &domain.FHIRTiming{
									Repeat: &domain.FHIRTimingRepeat{
										Period:       &number,
										PeriodUnit:   &periodUnit,
										Frequency:    &frequency,
										Duration:     &number,
										DurationUnit: &periodUnit,
										BoundsPeriod: &domain.FHIRPeriod{
											Start: scalarutils.DateTime(time.Now().GoString()),
											End:   scalarutils.DateTime(time.Now().GoString()),
										},
									},
								},
								PatientInstruction: &patientInstruction,
								AsNeeded:           &asneeded,
							},
						},
					},
				}

				mh.FHIR.EXPECT().GetFHIRMedicationRequest(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRMedicationRequestRelayPayload, error) {
						return medicationRequestOutput, nil
					})

				return args{ctx: ctx, medicationRequestID: medicationRequestID}
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to fetch medication request by ID",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRMedicationRequest(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRMedicationRequestRelayPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, medicationRequestID: medicationRequestID}
			},
			wantErr: true,
		},
		{
			name: "Sad case: medication request ID not provided",
			setup: func(mh *usecaseMock.Mocks) args {
				return args{ctx: ctx}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			_, err := clinicalUsecase.FetchMedicationRequestByID(args.ctx, args.medicationRequestID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.FetchMedicationRequestByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCasesClinicalImpl_FetchMedicationByID(t *testing.T) {
	ctx := usecaseMock.AddTenantIdentifierContext(context.Background())
	medicationID := uuid.NewString()

	type args struct {
		ctx          context.Context
		medicationID string
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy case: fetch medication by ID",
			setup: func(mh *usecaseMock.Mocks) args {
				status := domain.MedicationStatusEnumActive
				ID := uuid.NewString()
				code := scalarutils.Code(gofakeit.BeerName())
				lotNumber := "12"
				expirationDate := time.Now().GoString()
				medication := &domain.FHIRMedication{
					ID:     &ID,
					Status: &status,
					DoseForm: &domain.FHIRCodeableConcept{
						Coding: []*domain.FHIRCoding{
							{
								Code:    &code,
								Display: gofakeit.BeerName(),
							},
						},
					},
					Extension: []domain.Extension{
						{
							ValueString: gofakeit.BeerName(),
						},
					},
					Batch: &domain.FHIRMedicationBatch{
						LotNumber:      &lotNumber,
						ExpirationDate: &expirationDate,
					},
				}

				mh.FHIR.EXPECT().FetchMedicationByID(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRMedication, error) {
						return medication, nil
					})

				return args{ctx: ctx, medicationID: medicationID}
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to fetch medication by ID",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().FetchMedicationByID(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRMedication, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, medicationID: medicationID}
			},
			wantErr: true,
		},
		{
			name: "Sad case: medication ID not provided",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().FetchMedicationByID(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRMedication, error) {
						return nil, fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			_, err := clinicalUsecase.FetchMedicationByID(args.ctx, args.medicationID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.FetchMedicationByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
