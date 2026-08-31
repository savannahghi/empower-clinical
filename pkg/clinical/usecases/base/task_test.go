package base_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/scalarutils"
	"github.com/stretchr/testify/mock"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/common"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/empower-clinical/pkg/clinical/domain"
	usecaseMock "github.com/savannahghi/empower-clinical/pkg/clinical/usecases/mock"
)

func TestUseCasesClinicalImpl_AddTestResultsLater(t *testing.T) {
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
		ctx  context.Context
		task *dto.TaskInput
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy case: add test results for later",
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
				mh.FHIR.EXPECT().CreateFHIRTask(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input *domain.FHIRTaskInput) (*domain.FHIRTask, error) {
						taskID := uuid.NewString()
						status := scalarutils.Code(mock.Anything)
						return &domain.FHIRTask{
							ID: &taskID,
							BusinessStatus: &domain.FHIRCodeableConcept{
								ID:   &ID,
								Text: mock.Anything,
							},
							Description: mock.Anything,
							Status:      &status,
							StatusReason: &domain.FHIRCodeableReference{
								ID: &id,
								Concept: &domain.FHIRCodeableConcept{
									ID:   &id,
									Text: mock.Anything,
								},
							},
						}, nil
					})

				return args{
					ctx: usecaseMock.AddTenantIdentifierContext(context.Background()),
					task: &dto.TaskInput{
						EncounterID: gofakeit.UUID(),
						Task:        "VIA",
						Workflow:    dto.BreastCancerScreeningTypeEnum,
						Description: "A VIA test",
					},
				}
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to get encounter",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{
					ctx: usecaseMock.AddTenantIdentifierContext(context.Background()),
					task: &dto.TaskInput{
						EncounterID: gofakeit.UUID(),
						Task:        "VIA",
						Workflow:    dto.BreastCancerScreeningTypeEnum,
						Description: "A VIA test",
					},
				}
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get meta tags",
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

				return args{
					ctx: usecaseMock.AddTenantIdentifierContext(context.Background()),
					task: &dto.TaskInput{
						EncounterID: gofakeit.UUID(),
						Task:        "VIA",
						Workflow:    dto.BreastCancerScreeningTypeEnum,
						Description: "A VIA test",
					},
				}
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to create task",
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
				mh.FHIR.EXPECT().CreateFHIRTask(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input *domain.FHIRTaskInput) (*domain.FHIRTask, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{
					ctx: usecaseMock.AddTenantIdentifierContext(context.Background()),
					task: &dto.TaskInput{
						EncounterID: gofakeit.UUID(),
						Task:        "VIA",
						Workflow:    dto.BreastCancerScreeningTypeEnum,
						Description: "A VIA test",
					},
				}
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get facilityID from context",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})

				return args{
					ctx: context.Background(),
					task: &dto.TaskInput{
						EncounterID: gofakeit.UUID(),
						Task:        "VIA",
						Workflow:    dto.BreastCancerScreeningTypeEnum,
						Description: "A VIA test",
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

			_, err := clinicalUsecase.AddTestResultsLater(args.ctx, args.task)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.AddTestResultsLater() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCasesClinicalImpl_CreateTask(t *testing.T) {
	status := "requested"
	id := gofakeit.UUID()

	type args struct {
		ctx  context.Context
		task *domain.FHIRTaskInput
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy case: create fhir task",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().CreateFHIRTask(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input *domain.FHIRTaskInput) (*domain.FHIRTask, error) {
						taskID := uuid.NewString()
						status := scalarutils.Code(mock.Anything)
						return &domain.FHIRTask{
							ID: &taskID,
							BusinessStatus: &domain.FHIRCodeableConcept{
								ID:   &id,
								Text: mock.Anything,
							},
							Description: mock.Anything,
							Status:      &status,
							StatusReason: &domain.FHIRCodeableReference{
								ID: &id,
								Concept: &domain.FHIRCodeableConcept{
									ID:   &id,
									Text: mock.Anything,
								},
							},
						}, nil
					})

				return args{
					ctx: usecaseMock.AddTenantIdentifierContext(context.Background()),
					task: &domain.FHIRTaskInput{
						Status: (*scalarutils.Code)(&status),
						For: &domain.FHIRReferenceInput{
							ID: &id,
						},
					},
				}
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to create fhir task",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().CreateFHIRTask(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input *domain.FHIRTaskInput) (*domain.FHIRTask, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{
					ctx: usecaseMock.AddTenantIdentifierContext(context.Background()),
					task: &domain.FHIRTaskInput{
						Status: (*scalarutils.Code)(&status),
						For: &domain.FHIRReferenceInput{
							ID: &id,
						},
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

			_, err := clinicalUsecase.CreateTask(args.ctx, args.task)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.CreateTask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCasesClinicalImpl_UpdateTask(t *testing.T) {
	id := uuid.NewString()
	date := scalarutils.DateTime(time.Thursday.String())

	type args struct {
		ctx        context.Context
		taskID     string
		updateData *dto.PatchTaskInput
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy case: update task",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().PatchFHIRTask(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRTaskInput) (*domain.FHIRTask, error) {
						taskID := uuid.NewString()
						status := scalarutils.Code(mock.Anything)

						return &domain.FHIRTask{
							ID: &taskID,
							BusinessStatus: &domain.FHIRCodeableConcept{
								ID:   &id,
								Text: mock.Anything,
							},
							Description: mock.Anything,
							Status:      &status,
							StatusReason: &domain.FHIRCodeableReference{
								ID: &id,
								Concept: &domain.FHIRCodeableConcept{
									ID:   &id,
									Text: mock.Anything,
								},
							},
						}, nil
					})

				return args{
					ctx:    usecaseMock.AddTenantIdentifierContext(context.Background()),
					taskID: gofakeit.UUID(),
					updateData: &dto.PatchTaskInput{
						Status:       dto.CompletedTasksStatus,
						DueDate:      date,
						Author:       gofakeit.BeerName(),
						Notes:        gofakeit.HipsterSentence(10),
						UpdateReason: "Test returned",
					},
				}
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to update task",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().PatchFHIRTask(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRTaskInput) (*domain.FHIRTask, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{
					ctx:    usecaseMock.AddTenantIdentifierContext(context.Background()),
					taskID: gofakeit.UUID(),
					updateData: &dto.PatchTaskInput{
						Status:       dto.CompletedTasksStatus,
						DueDate:      date,
						Author:       gofakeit.BeerName(),
						Notes:        gofakeit.HipsterSentence(10),
						UpdateReason: "Test returned",
					},
				}
			},
			wantErr: true,
		},
		{
			name: "Sad case: invalid task ID",
			setup: func(mh *usecaseMock.Mocks) args {
				return args{
					ctx:    usecaseMock.AddTenantIdentifierContext(context.Background()),
					taskID: "invalid",
					updateData: &dto.PatchTaskInput{
						Status:       dto.CompletedTasksStatus,
						DueDate:      date,
						Author:       gofakeit.BeerName(),
						Notes:        gofakeit.HipsterSentence(10),
						UpdateReason: "Test returned",
					},
				}
			},
			wantErr: true,
		},
		{
			name: "Sad case: fail if update reason is not provided",
			setup: func(mh *usecaseMock.Mocks) args {
				return args{
					ctx:    usecaseMock.AddTenantIdentifierContext(context.Background()),
					taskID: gofakeit.UUID(),
					updateData: &dto.PatchTaskInput{
						Status:  dto.CompletedTasksStatus,
						Author:  gofakeit.BeerName(),
						Notes:   gofakeit.HipsterSentence(10),
						DueDate: date,
					},
				}
			},
			wantErr: true,
		},
		{
			name: "Sad case: fail if update reason provided but no status",
			setup: func(mh *usecaseMock.Mocks) args {
				return args{
					ctx:    usecaseMock.AddTenantIdentifierContext(context.Background()),
					taskID: gofakeit.UUID(),
					updateData: &dto.PatchTaskInput{
						Author:       gofakeit.BeerName(),
						Notes:        gofakeit.HipsterSentence(10),
						DueDate:      date,
						UpdateReason: "test",
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

			_, err := clinicalUsecase.UpdateTask(args.ctx, args.taskID, args.updateData)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.UpdateTask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCasesClinicalImpl_ListTask(t *testing.T) {
	tenantIDs := &dto.TenantIdentifiers{
		OrganizationID: uuid.NewString(),
		FacilityID:     uuid.NewString(),
	}
	firstTen := 10
	businessStatus := "DIAGNOSTICS"

	type args struct {
		ctx        context.Context
		filters    *dto.TaskFilterInput
		pagination dto.Pagination
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully get task",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchFHIRTask(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRTask, error) {
						taskID := uuid.NewString()
						status := scalarutils.Code(mock.Anything)
						id := uuid.NewString()
						ref := "Reference/123456"
						authoredOn := "2025-01-12"
						note := scalarutils.Markdown(mock.Anything)
						priotity := scalarutils.Code("asap")

						return &domain.PagedFHIRTask{
							Tasks: []domain.FHIRTask{
								{
									ID: &taskID,
									BusinessStatus: &domain.FHIRCodeableConcept{
										ID:   &id,
										Text: mock.Anything,
									},
									Description: mock.Anything,
									Status:      &status,
									StatusReason: &domain.FHIRCodeableReference{
										ID: &id,
										Concept: &domain.FHIRCodeableConcept{
											ID:   &id,
											Text: mock.Anything,
										},
									},
									Priority: &priotity,
									For: &domain.FHIRReference{
										Reference: &ref,
										Display:   mock.Anything,
									},
									AuthoredOn: &authoredOn,
									Meta: &domain.FHIRMeta{
										LastUpdated: gofakeit.Date(),
									},
									Encounter: &domain.FHIRReference{
										Reference: &ref,
									},
									Reason: []*domain.FHIRCodeableReference{
										{
											Concept: &domain.FHIRCodeableConcept{
												Text: mock.Anything,
											},
										},
									},
									ExecutionPeriod: &domain.FHIRPeriod{
										End: scalarutils.DateTime(gofakeit.Date().UTC().GoString()),
									},
									Note: []*domain.FHIRAnnotation{
										{
											Text: &note,
										},
									},
								},
							},
						}, nil
					})

				return args{
					ctx: usecaseMock.AddTenantIdentifierContext(context.Background()),
					filters: &dto.TaskFilterInput{
						FilterInput: dto.FilterInput{
							PatientID:   gofakeit.UUID(),
							EncounterID: gofakeit.UUID(),
							Date: &scalarutils.Date{
								Year:  2023,
								Month: 10,
								Day:   4,
							},
						},
						Type: businessStatus,
					},
					pagination: dto.Pagination{
						First: &firstTen,
					},
				}
			},
			wantErr: false,
		},
		{
			name: "Happy Case - Successfully get task - with business status",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchFHIRTask(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRTask, error) {
						taskID := uuid.NewString()
						status := scalarutils.Code(mock.Anything)
						id := uuid.NewString()
						ref := "Reference/123456"
						authoredOn := "2025-01-12"
						note := scalarutils.Markdown(mock.Anything)
						priotity := scalarutils.Code("asap")

						return &domain.PagedFHIRTask{
							Tasks: []domain.FHIRTask{
								{
									ID: &taskID,
									BusinessStatus: &domain.FHIRCodeableConcept{
										ID:   &id,
										Text: mock.Anything,
									},
									Status: &status,
									StatusReason: &domain.FHIRCodeableReference{
										ID: &id,
										Concept: &domain.FHIRCodeableConcept{
											ID:   &id,
											Text: mock.Anything,
										},
									},
									Priority: &priotity,
									For: &domain.FHIRReference{
										Reference: &ref,
										Display:   mock.Anything,
									},
									AuthoredOn: &authoredOn,
									Meta: &domain.FHIRMeta{
										LastUpdated: gofakeit.Date(),
									},
									Encounter: &domain.FHIRReference{
										Reference: &ref,
									},
									Reason: []*domain.FHIRCodeableReference{
										{
											Concept: &domain.FHIRCodeableConcept{
												Text: mock.Anything,
											},
										},
									},
									ExecutionPeriod: &domain.FHIRPeriod{
										End: scalarutils.DateTime(gofakeit.Date().UTC().GoString()),
									},
									Note: []*domain.FHIRAnnotation{
										{
											Text: &note,
										},
									},
								},
							},
						}, nil
					})

				return args{
					ctx: usecaseMock.AddTenantIdentifierContext(context.Background()),
					filters: &dto.TaskFilterInput{
						FilterInput: dto.FilterInput{
							PatientID:   gofakeit.UUID(),
							EncounterID: gofakeit.UUID(),
							Date: &scalarutils.Date{
								Year:  2023,
								Month: 10,
								Day:   4,
							},
						},
						Type: businessStatus,
					},
					pagination: dto.Pagination{
						First: &firstTen,
					},
				}
			},
			wantErr: false,
		},
		{
			name: "Happy Case - general patient search resolves to a patient IN-list",
			setup: func(mh *usecaseMock.Mocks) args {
				p1 := "patient-1"
				p2 := "patient-2"

				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID, resourceType string, params map[string]any, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRResource, error) {
						return &domain.PagedFHIRResource{Resources: []map[string]interface{}{
							{"id": p1},
							{"id": p2},
						}}, nil
					})
				mh.FHIR.EXPECT().SearchFHIRTask(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRTask, error) {
						if params["patient"] != "Patient/patient-1,Patient/patient-2" {
							return nil, fmt.Errorf("expected patient IN-list, got %v", params["patient"])
						}

						return &domain.PagedFHIRTask{Tasks: []domain.FHIRTask{}}, nil
					})

				return args{
					ctx:        usecaseMock.AddTenantIdentifierContext(context.Background()),
					filters:    &dto.TaskFilterInput{PatientSearch: "jane"},
					pagination: dto.Pagination{First: &firstTen},
				}
			},
			wantErr: false,
		},
		{
			name: "Happy Case - general patient search with no matches short-circuits to empty",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				// No patients match, so the task search must never run.
				mh.FHIR.EXPECT().SearchFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID, resourceType string, params map[string]any, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRResource, error) {
						return &domain.PagedFHIRResource{Resources: []map[string]interface{}{}}, nil
					})

				return args{
					ctx:        usecaseMock.AddTenantIdentifierContext(context.Background()),
					filters:    &dto.TaskFilterInput{PatientSearch: "nobody"},
					pagination: dto.Pagination{First: &firstTen},
				}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - invalid patient id",
			setup: func(mh *usecaseMock.Mocks) args {
				return args{
					ctx: usecaseMock.AddTenantIdentifierContext(context.Background()),
					filters: &dto.TaskFilterInput{
						FilterInput: dto.FilterInput{
							PatientID:   "id",
							EncounterID: gofakeit.UUID(),
							Date: &scalarutils.Date{
								Year:  2023,
								Month: 10,
								Day:   4,
							},
						},
						Type: businessStatus,
					},
					pagination: dto.Pagination{
						First: &firstTen,
					},
				}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - invalid encounter id",
			setup: func(mh *usecaseMock.Mocks) args {
				return args{
					ctx: usecaseMock.AddTenantIdentifierContext(context.Background()),
					filters: &dto.TaskFilterInput{
						FilterInput: dto.FilterInput{
							PatientID:   gofakeit.UUID(),
							EncounterID: "id",
							Date: &scalarutils.Date{
								Year:  2023,
								Month: 10,
								Day:   4,
							},
						},
						Type: businessStatus,
					},
					pagination: dto.Pagination{
						First: &firstTen,
					},
				}
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

				return args{
					ctx: usecaseMock.AddTenantIdentifierContext(context.Background()),
					filters: &dto.TaskFilterInput{
						FilterInput: dto.FilterInput{
							PatientID:   gofakeit.UUID(),
							EncounterID: gofakeit.UUID(),
							Date: &scalarutils.Date{
								Year:  2023,
								Month: 10,
								Day:   4,
							},
						},
						Type: businessStatus,
					},
					pagination: dto.Pagination{
						First: &firstTen,
					},
				}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to search for task",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchFHIRTask(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRTask, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{
					ctx: usecaseMock.AddTenantIdentifierContext(context.Background()),
					filters: &dto.TaskFilterInput{
						FilterInput: dto.FilterInput{
							PatientID:   gofakeit.UUID(),
							EncounterID: gofakeit.UUID(),
							Date: &scalarutils.Date{
								Year:  2023,
								Month: 10,
								Day:   4,
							},
						},
						Type: businessStatus,
					},
					pagination: dto.Pagination{
						First: &firstTen,
					},
				}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			_, err := clinicalUsecase.ListTasks(args.ctx, args.filters, args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.ListTasks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCasesClinicalImpl_CreateAppointmentTask(t *testing.T) {
	id := gofakeit.UUID()
	ref := fmt.Sprintf("Reference/%s", id)
	tagsInput := &domain.FHIRMetaInput{
		Tag: []domain.FHIRCodingInput{
			{
				Code: scalarutils.Code(id),
			},
		},
	}

	tags := &domain.FHIRMeta{
		Tag: []domain.FHIRCoding{
			{
				Code: (*scalarutils.Code)(&id),
			},
		},
	}

	participant := []*domain.FHIRAppointmentParticipant{
		{
			Actor: &domain.FHIRReference{
				ID:        &id,
				Reference: &ref,
			},
		},
	}

	instant := scalarutils.Instant(time.Now().Format(time.RFC3339))

	type args struct {
		ctx         context.Context
		encounterID string
		tags        *domain.FHIRMetaInput
		appointment *domain.FHIRAppointment
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy case: create appointment task",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().CreateFHIRTask(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input *domain.FHIRTaskInput) (*domain.FHIRTask, error) {
						return nil, nil
					})

				return args{
					ctx:         usecaseMock.AddTenantIdentifierContext(context.Background()),
					encounterID: id,
					tags:        tagsInput,
					appointment: &domain.FHIRAppointment{
						ID:          &id,
						Meta:        tags,
						Participant: participant,
						End:         &instant,
						Reason: []*domain.FHIRCodeableReference{
							{
								Concept: &domain.FHIRCodeableConcept{
									Text: mock.Anything,
								},
							},
						},
					},
				}
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to create appointment task",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().CreateFHIRTask(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input *domain.FHIRTaskInput) (*domain.FHIRTask, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{
					ctx:         usecaseMock.AddTenantIdentifierContext(context.Background()),
					encounterID: id,
					tags:        tagsInput,
					appointment: &domain.FHIRAppointment{
						ID:          &id,
						Meta:        tags,
						Participant: participant,
						End:         &instant,
					},
				}
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get facility ID from context",
			setup: func(mh *usecaseMock.Mocks) args {
				return args{
					ctx:         context.Background(),
					encounterID: id,
					tags:        tagsInput,
					appointment: &domain.FHIRAppointment{
						ID:          &id,
						Meta:        tags,
						Participant: participant,
						End:         &instant,
					},
				}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			_, err := clinicalUsecase.CreateAppointmentTask(args.ctx, args.encounterID, args.tags, args.appointment)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.CreateAppointmentTask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCasesClinicalImpl_CreateReferralTask(t *testing.T) {
	ID := gofakeit.UUID()
	ref := fmt.Sprintf("Reference/%s", ID)

	noteText := scalarutils.Markdown(gofakeit.BeerName())
	serviceCode := scalarutils.Code("159623")
	authoredOn := scalarutils.DateTime(gofakeit.Date().GoString())
	system := scalarutils.URI("http://mycarehub/tenant-identification/facility")
	code := scalarutils.Code(gofakeit.UUID())
	status := domain.NarrativeStatusEnumAdditional
	servicerequest := &domain.FHIRServiceRequestRelayPayload{
		Resource: &domain.FHIRServiceRequest{
			ID: &ID,
			Text: &domain.FHIRNarrative{
				Status: &status,
				Div:    scalarutils.XHTML(gofakeit.BeerName()),
			},
			Priority:   domain.ServiceRequestPriorityAsap,
			Identifier: []*domain.FHIRIdentifier{},
			Status:     domain.ServiceRequestStatusActive,
			Intent:     domain.ServiceRequestIntentDirective,
			Subject: &domain.FHIRReference{
				ID:        &ID,
				Reference: &ref,
				Display:   gofakeit.UUID(),
			},
			AuthoredOn: &authoredOn,
			Encounter: &domain.FHIRReference{
				ID:        &ID,
				Display:   gofakeit.UUID(),
				Reference: &ref,
			},
			Note: []*domain.FHIRAnnotation{
				{
					Text: &noteText,
				},
			},
			Code: &domain.FHIRCodeableReference{
				Concept: &domain.FHIRCodeableConcept{
					Coding: []*domain.FHIRCoding{
						{
							Display: gofakeit.UUID(),
							Code:    &serviceCode,
						},
					},
				},
			},
			Extension: []*domain.FHIRExtension{
				{
					URL: "http://savannahghi.org/fhir/StructureDefinition/referred-facility",
					Extension: []domain.Extension{
						{
							URL:         "facilityName",
							ValueString: gofakeit.BeerName(),
						},
						{
							URL:         "facilityCounty",
							ValueString: gofakeit.BeerName(),
						},
						{
							URL:         "facilityContact",
							ValueString: gofakeit.Contact().Phone,
						},
						{
							URL:         "facilityEmail",
							ValueString: gofakeit.Contact().Email,
						},
					},
				},
			},
			Meta: &domain.FHIRMeta{
				Tag: []domain.FHIRCoding{
					{
						System: &system,
						Code:   &code,
					},
				},
			},
			Performer: []*domain.FHIRReference{},
		},
	}

	patient := &domain.FHIRPatientRelayPayload{
		Resource: &domain.FHIRPatient{
			ID: &ID,
		},
	}

	type args struct {
		ctx            context.Context
		tags           *dto.MetaInput
		serviceRequest *dto.ServiceRequest
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy case: Create referral test",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRServiceRequest(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRServiceRequestRelayPayload, error) {
						return servicerequest, nil
					})
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return patient, nil
					})
				mh.FHIR.EXPECT().CreateFHIRServiceRequest(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRServiceRequestInput) (*domain.FHIRServiceRequestRelayPayload, error) {
						return servicerequest, nil
					})
				mh.FHIR.EXPECT().CreateFHIRTask(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input *domain.FHIRTaskInput) (*domain.FHIRTask, error) {
						return nil, nil
					})

				return args{
					ctx: usecaseMock.AddTenantIdentifierContext(context.Background()),
					tags: &dto.MetaInput{
						Tag: []dto.Coding{
							{
								Code:         (*scalarutils.Code)(&ID),
								Display:      "Facility",
								UserSelected: false,
							},
						},
					},
					serviceRequest: &dto.ServiceRequest{
						ID: ID,
						Subject: dto.Reference{
							ID:        ID,
							Reference: ref,
						},
						Encounter: &dto.Reference{
							ID:        ID,
							Reference: ref,
						},
					},
				}
			},
			wantErr: false,
		},
		{
			name: "Sad case: Unable to create referral test",
			setup: func(mh *usecaseMock.Mocks) args {

				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return patient, nil
					})
				mh.FHIR.EXPECT().GetFHIRServiceRequest(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRServiceRequestRelayPayload, error) {
						return servicerequest, nil
					})
				mh.FHIR.EXPECT().CreateFHIRServiceRequest(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRServiceRequestInput) (*domain.FHIRServiceRequestRelayPayload, error) {
						return servicerequest, nil
					})
				mh.FHIR.EXPECT().CreateFHIRTask(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input *domain.FHIRTaskInput) (*domain.FHIRTask, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{
					ctx: usecaseMock.AddTenantIdentifierContext(context.Background()),
					tags: &dto.MetaInput{
						Tag: []dto.Coding{
							{
								Code:         (*scalarutils.Code)(&ID),
								Display:      "Facility",
								UserSelected: false,
							},
						},
					},
					serviceRequest: &dto.ServiceRequest{
						ID: ID,
						Subject: dto.Reference{
							ID:        ID,
							Reference: ref,
						},
						Encounter: &dto.Reference{
							ID:        ID,
							Reference: ref,
						},
					},
				}
			},
			wantErr: true,
		},
		{
			name: "Sad case: Unable to create service request",
			setup: func(mh *usecaseMock.Mocks) args {

				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return patient, nil
					})
				mh.FHIR.EXPECT().GetFHIRServiceRequest(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRServiceRequestRelayPayload, error) {
						return servicerequest, nil
					})
				mh.FHIR.EXPECT().CreateFHIRServiceRequest(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRServiceRequestInput) (*domain.FHIRServiceRequestRelayPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{
					ctx: usecaseMock.AddTenantIdentifierContext(context.Background()),
					tags: &dto.MetaInput{
						Tag: []dto.Coding{
							{
								Code:         (*scalarutils.Code)(&ID),
								Display:      "Facility",
								UserSelected: false,
							},
						},
					},
					serviceRequest: &dto.ServiceRequest{
						ID: ID,
						Subject: dto.Reference{
							ID:        ID,
							Reference: ref,
						},
						Encounter: &dto.Reference{
							ID:        ID,
							Reference: ref,
						},
					},
				}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			_, err := clinicalUsecase.CreateReferralTask(args.ctx, args.tags, args.serviceRequest)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.CreateReferralTask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCasesClinicalImpl_GetTaskByID(t *testing.T) {
	tenantIDs := &dto.TenantIdentifiers{
		OrganizationID: uuid.NewString(),
		FacilityID:     uuid.NewString(),
	}

	taskID := uuid.NewString()
	status := scalarutils.Code(mock.Anything)
	ref := "Reference/123456"
	authoredOn := "2025-01-12"
	note := scalarutils.Markdown(mock.Anything)
	priotity := scalarutils.Code("asap")
	basedOnType := scalarutils.URI("Patient Referral (Service Request)")
	id := uuid.NewString()

	taskOutput := &domain.FHIRTaskRelayPayload{
		Resource: &domain.FHIRTask{
			ID: &taskID,
			BusinessStatus: &domain.FHIRCodeableConcept{
				ID:   &id,
				Text: mock.Anything,
			},
			Description: mock.Anything,
			Status:      &status,
			StatusReason: &domain.FHIRCodeableReference{
				ID: &id,
				Concept: &domain.FHIRCodeableConcept{
					ID:   &id,
					Text: mock.Anything,
				},
			},
			Priority: &priotity,
			For: &domain.FHIRReference{
				Reference: &ref,
				Display:   mock.Anything,
			},
			AuthoredOn: &authoredOn,
			Meta: &domain.FHIRMeta{
				LastUpdated: gofakeit.Date(),
			},
			Encounter: &domain.FHIRReference{
				Reference: &ref,
			},
			Reason: []*domain.FHIRCodeableReference{
				{
					Concept: &domain.FHIRCodeableConcept{
						Text: mock.Anything,
					},
				},
			},
			ExecutionPeriod: &domain.FHIRPeriod{
				End: scalarutils.DateTime(gofakeit.Date().UTC().GoString()),
			},
			Note: []*domain.FHIRAnnotation{
				{
					Text: &note,
				},
			},
			BasedOn: []*domain.FHIRReference{
				{
					ID:   &id,
					Type: &basedOnType,
				},
			},
		},
	}

	phoneSystem := domain.ContactPointSystemEnumPhone
	phoneNumber := gofakeit.Phone()
	patientPayload := &domain.FHIRPatientRelayPayload{
		Resource: &domain.FHIRPatient{
			ID: &id,
			Telecom: []*domain.FHIRContactPoint{
				{
					System: &phoneSystem,
					Value:  &phoneNumber,
				},
			},
		},
	}

	type args struct {
		ctx    context.Context
		taskID string
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy case: get task by id",
			setup: func(mh *usecaseMock.Mocks) args {
				resourceID := gofakeit.UUID()
				URL := gofakeit.URL()
				contentType := "application/pdf"
				title := "Test title"
				code := common.IntraReferralLOINCCode
				pagedDocRef := &domain.PagedFHIRDocumentReference{
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
									ID:                resourceID,
									Extension:         []domain.Extension{},
									ModifierExtension: []domain.Extension{},
									Attachment: domain.FHIRAttachment{
										ContentType: (*scalarutils.Code)(&contentType),
										URL:         (*scalarutils.URL)(&URL),
										Title:       &title,
									},
									Format: &domain.FHIRCoding{},
								},
							},
							Type: &domain.FHIRCodeableConcept{
								ID: new(string),
								Coding: []*domain.FHIRCoding{
									{
										Code:    (*scalarutils.Code)(&code),
										Display: "",
									},
								},
								Text: "",
							},
						},
					},
					HasNextPage:     false,
					NextCursor:      "",
					HasPreviousPage: false,
					PreviousCursor:  "",
					TotalCount:      0,
				}

				mh.FHIR.EXPECT().GetFHIRTask(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRTaskRelayPayload, error) {
						return taskOutput, nil
					})
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return patientPayload, nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchFHIRDocumentReference(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, searchParams map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRDocumentReference, error) {
						return pagedDocRef, nil
					})

				return args{
					ctx:    context.Background(),
					taskID: gofakeit.UUID(),
				}
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to get tenant identifiers",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRTask(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRTaskRelayPayload, error) {
						return taskOutput, nil
					})
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return patientPayload, nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{
					ctx:    context.Background(),
					taskID: gofakeit.UUID(),
				}
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to search document reference",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRTask(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRTaskRelayPayload, error) {
						return taskOutput, nil
					})
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return patientPayload, nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchFHIRDocumentReference(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, searchParams map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRDocumentReference, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{
					ctx:    context.Background(),
					taskID: gofakeit.UUID(),
				}
			},
			wantErr: true,
		},
		{
			name: "Sad case: invalid task ID",
			setup: func(mh *usecaseMock.Mocks) args {
				return args{
					ctx:    context.Background(),
					taskID: "invalid",
				}
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get attachment link",
			setup: func(mh *usecaseMock.Mocks) args {
				contentType := "application/pdf"
				title := "Test title"
				code := common.IntraReferralLOINCCode
				pagedDocRef := &domain.PagedFHIRDocumentReference{
					DocumentReferences: []domain.FHIRDocumentReference{
						{
							Content: []domain.FHIRDocumentReferenceContent{
								{
									Attachment: domain.FHIRAttachment{
										ContentType: (*scalarutils.Code)(&contentType),
										Title:       &title,
									},
									Format: &domain.FHIRCoding{},
								},
							},
							Type: &domain.FHIRCodeableConcept{
								ID: new(string),
								Coding: []*domain.FHIRCoding{
									{
										Code:    (*scalarutils.Code)(&code),
										Display: "",
									},
								},
								Text: "",
							},
						},
					},
					HasNextPage:     false,
					NextCursor:      "",
					HasPreviousPage: false,
					PreviousCursor:  "",
					TotalCount:      0,
				}

				mh.FHIR.EXPECT().GetFHIRTask(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRTaskRelayPayload, error) {
						return taskOutput, nil
					})
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return patientPayload, nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchFHIRDocumentReference(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, searchParams map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRDocumentReference, error) {
						return pagedDocRef, nil
					})

				return args{
					ctx:    context.Background(),
					taskID: gofakeit.UUID(),
				}
			},
			wantErr: true,
		},
		{
			name: "Happy case: show attachment as external referral",
			setup: func(mh *usecaseMock.Mocks) args {
				resourceID := gofakeit.UUID()
				URL := gofakeit.URL()
				contentType := "application/pdf"
				title := "Test title"
				code := common.ReferralLOINCTerminologySystem
				pagedDocRef := &domain.PagedFHIRDocumentReference{
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
									ID:                resourceID,
									Extension:         []domain.Extension{},
									ModifierExtension: []domain.Extension{},
									Attachment: domain.FHIRAttachment{
										ContentType: (*scalarutils.Code)(&contentType),
										URL:         (*scalarutils.URL)(&URL),
										Title:       &title,
									},
									Format: &domain.FHIRCoding{},
								},
							},
							Type: &domain.FHIRCodeableConcept{
								ID: new(string),
								Coding: []*domain.FHIRCoding{
									{
										Code:    (*scalarutils.Code)(&code),
										Display: "",
									},
								},
								Text: "",
							},
						},
					},
					HasNextPage:     false,
					NextCursor:      "",
					HasPreviousPage: false,
					PreviousCursor:  "",
					TotalCount:      0,
				}

				mh.FHIR.EXPECT().GetFHIRTask(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRTaskRelayPayload, error) {
						return taskOutput, nil
					})
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return patientPayload, nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchFHIRDocumentReference(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, searchParams map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRDocumentReference, error) {
						return pagedDocRef, nil
					})

				return args{
					ctx:    context.Background(),
					taskID: gofakeit.UUID(),
				}
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to get patient for subject telecom",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRTask(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRTaskRelayPayload, error) {
						return taskOutput, nil
					})
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{
					ctx:    context.Background(),
					taskID: gofakeit.UUID(),
				}
			},
			wantErr: true,
		},
		{
			name: "Happy case: subject identifier present, skips telecom lookup",
			setup: func(mh *usecaseMock.Mocks) args {
				taskWithIdentifier := &domain.FHIRTaskRelayPayload{
					Resource: &domain.FHIRTask{
						ID:       &taskID,
						Status:   &status,
						Priority: &priotity,
						BusinessStatus: &domain.FHIRCodeableConcept{
							ID:   &id,
							Text: mock.Anything,
						},
						For: &domain.FHIRReference{
							Reference: &ref,
							Display:   mock.Anything,
							Identifier: &domain.FHIRIdentifier{
								Value: phoneNumber,
							},
						},
					},
				}

				mh.FHIR.EXPECT().GetFHIRTask(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRTaskRelayPayload, error) {
						return taskWithIdentifier, nil
					})

				return args{
					ctx:    context.Background(),
					taskID: gofakeit.UUID(),
				}
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			_, err := clinicalUsecase.GetTaskByID(args.ctx, args.taskID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.GetTaskByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestBaseImpl_FetchPatientCarePlan(t *testing.T) {
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

	type args struct {
		ctx         context.Context
		encounterID string
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy case: Successfully fetch patient plan",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.FHIR.EXPECT().SearchFHIRCarePlan(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, pagination dto.Pagination) (*domain.PagedFHIRCarePlan, error) {
						title := mock.Anything
						description := mock.Anything
						ref := "Resource/12345"
						careplan := &domain.PagedFHIRCarePlan{
							CarePlan: []domain.FHIRCarePlan{
								{
									ID:           &id,
									ResourceType: "Service Request",
									Title:        &title,
									Description:  &description,
									Subject: domain.FHIRReference{
										Display:   mock.Anything,
										Reference: &ref,
									},
									Activity: []domain.CarePlanActivity{
										{
											ID: &id,
											PerformedActivity: []domain.FHIRCodeableReference{
												{
													ID: &id,
													Concept: &domain.FHIRCodeableConcept{
														ID:   &id,
														Text: mock.Anything,
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
						}
						return careplan, nil
					})
				mh.FHIR.EXPECT().GetFHIRTask(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRTaskRelayPayload, error) {
						taskID := uuid.NewString()
						status := scalarutils.Code(mock.Anything)

						return &domain.FHIRTaskRelayPayload{
							Resource: &domain.FHIRTask{
								ID: &taskID,
								BusinessStatus: &domain.FHIRCodeableConcept{
									ID:   &id,
									Text: mock.Anything,
								},
								Description: mock.Anything,
								Status:      &status,
								StatusReason: &domain.FHIRCodeableReference{
									ID: &id,
									Concept: &domain.FHIRCodeableConcept{
										ID:   &id,
										Text: mock.Anything,
									},
								},
							},
						}, nil
					})
				mh.FHIR.EXPECT().SearchFHIRTask(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRTask, error) {
						taskID := uuid.NewString()
						status := scalarutils.Code(mock.Anything)

						return &domain.PagedFHIRTask{
							Tasks: []domain.FHIRTask{
								{
									ID: &taskID,
									BusinessStatus: &domain.FHIRCodeableConcept{
										ID:   &id,
										Text: mock.Anything,
									},
									Description: mock.Anything,
									Status:      &status,
									StatusReason: &domain.FHIRCodeableReference{
										ID: &id,
										Concept: &domain.FHIRCodeableConcept{
											ID:   &id,
											Text: mock.Anything,
										},
									},
								},
							},
						}, nil
					})

				return args{
					ctx:         usecaseMock.AddTenantIdentifierContext(context.Background()),
					encounterID: gofakeit.UUID(),
				}
			},
			wantErr: false,
		},
		{
			name: "Sad case: Unable to search care plan",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.FHIR.EXPECT().SearchFHIRCarePlan(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, pagination dto.Pagination) (*domain.PagedFHIRCarePlan, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{
					ctx:         usecaseMock.AddTenantIdentifierContext(context.Background()),
					encounterID: gofakeit.UUID(),
				}
			},
			wantErr: true,
		},
		{
			name: "Sad case: Unable to get task",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.FHIR.EXPECT().SearchFHIRCarePlan(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, pagination dto.Pagination) (*domain.PagedFHIRCarePlan, error) {
						title := mock.Anything
						description := mock.Anything
						ref := "Resource/12345"
						careplan := &domain.PagedFHIRCarePlan{
							CarePlan: []domain.FHIRCarePlan{
								{
									ID:           &id,
									ResourceType: "Service Request",
									Title:        &title,
									Description:  &description,
									Subject: domain.FHIRReference{
										Display:   mock.Anything,
										Reference: &ref,
									},
									Activity: []domain.CarePlanActivity{
										{
											ID: &id,
											PerformedActivity: []domain.FHIRCodeableReference{
												{
													ID: &id,
													Concept: &domain.FHIRCodeableConcept{
														ID:   &id,
														Text: mock.Anything,
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
						}
						return careplan, nil
					})
				mh.FHIR.EXPECT().GetFHIRTask(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRTaskRelayPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{
					ctx:         usecaseMock.AddTenantIdentifierContext(context.Background()),
					encounterID: gofakeit.UUID(),
				}
			},
			wantErr: true,
		},
		{
			name: "Sad case: Unable to get encounter",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{
					ctx:         usecaseMock.AddTenantIdentifierContext(context.Background()),
					encounterID: gofakeit.UUID(),
				}
			},
			wantErr: true,
		},
		{
			name: "Sad case: Unable to search task",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.FHIR.EXPECT().SearchFHIRCarePlan(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, pagination dto.Pagination) (*domain.PagedFHIRCarePlan, error) {
						title := mock.Anything
						description := mock.Anything
						ref := "Resource/12345"
						careplan := &domain.PagedFHIRCarePlan{
							CarePlan: []domain.FHIRCarePlan{
								{
									ID:           &id,
									ResourceType: "Service Request",
									Title:        &title,
									Description:  &description,
									Subject: domain.FHIRReference{
										Display:   mock.Anything,
										Reference: &ref,
									},
									Activity: []domain.CarePlanActivity{
										{
											ID: &id,
											PerformedActivity: []domain.FHIRCodeableReference{
												{
													ID: &id,
													Concept: &domain.FHIRCodeableConcept{
														ID:   &id,
														Text: mock.Anything,
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
						}
						return careplan, nil
					})
				mh.FHIR.EXPECT().GetFHIRTask(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRTaskRelayPayload, error) {
						taskID := uuid.NewString()
						status := scalarutils.Code(mock.Anything)

						return &domain.FHIRTaskRelayPayload{
							Resource: &domain.FHIRTask{
								ID: &taskID,
								BusinessStatus: &domain.FHIRCodeableConcept{
									ID:   &id,
									Text: mock.Anything,
								},
								Description: mock.Anything,
								Status:      &status,
								StatusReason: &domain.FHIRCodeableReference{
									ID: &id,
									Concept: &domain.FHIRCodeableConcept{
										ID:   &id,
										Text: mock.Anything,
									},
								},
							},
						}, nil
					})
				mh.FHIR.EXPECT().SearchFHIRTask(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRTask, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{
					ctx:         usecaseMock.AddTenantIdentifierContext(context.Background()),
					encounterID: gofakeit.UUID(),
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

			_, err := clinicalUsecase.FetchPatientCarePlan(args.ctx, args.encounterID)
			if (err != nil) != tt.wantErr {
				t.Errorf("BaseImpl.FetchPatientCarePlan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
