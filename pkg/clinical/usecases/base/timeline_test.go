package base_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/savannahghi/scalarutils"
	"github.com/stretchr/testify/mock"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/empower-clinical/pkg/clinical/domain"
	usecaseMock "github.com/savannahghi/empower-clinical/pkg/clinical/usecases/mock"
)

func TestBaseImpl_GetPatientBanner(t *testing.T) {
	ctx := usecaseMock.AddTenantIdentifierContext(context.Background())

	type args struct {
		ctx       context.Context
		patientID string
		params    *dto.PatientEverythingFilterParams
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy case: get patient banner",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						patientID := gofakeit.UUID()

						return &domain.FHIRPatientRelayPayload{
							Resource: &domain.FHIRPatient{
								ID: &patientID,
							},
						}, nil
					})
				mh.FHIR.EXPECT().GetFHIRPatientEverything(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string, params map[string]interface{}) (*domain.PagedFHIRResource, error) {
						return &domain.PagedFHIRResource{
							Resources: []map[string]interface{}{
								{
									"resourceType": "Condition",
									"id":           "9012",
									"code": map[string]interface{}{
										"coding": []map[string]interface{}{
											{
												"code":    "161826",
												"display": "Biopsy of cervix",
												"system":  "/orgs/CIEL/sources/CIEL/concepts/161826/",
											},
										},
										"text": "Condition",
									},
									"clinicalStatus": map[string]interface{}{
										"text": "active",
									},
									"onsetDateTime": "2024-02-13T10:22:54+03:00",
									"category": []map[string]interface{}{
										{
											"coding": []map[string]interface{}{
												{
													"code":         "laboratory",
													"display":      "Laboratory",
													"system":       "http://terminology.hl7.org/CodeSystem/observation-category",
													"userSelected": false,
												},
											},
											"text": "Laboratory",
										},
									},
								},
								{
									"resourceType": "MedicationStatement",
									"id":           "9012",
									"status":       "active",
									"code": map[string]interface{}{
										"coding": []map[string]interface{}{
											{
												"code":    "161826",
												"display": "Biopsy of cervix",
												"system":  "/orgs/CIEL/sources/CIEL/concepts/161826/",
											},
										},
										"text": "Biopsy of cervix",
									},
									"dateAsserted": "2025-02-13",
									"medication": []map[string]interface{}{
										{
											"concept": map[string]interface{}{
												"text": "medication concept",
											},
										},
									},
									"category": []map[string]interface{}{
										{
											"coding": []map[string]interface{}{
												{
													"code":         "laboratory",
													"display":      "Laboratory",
													"system":       "http://terminology.hl7.org/CodeSystem/observation-category",
													"userSelected": false,
												},
											},
											"text": "Laboratory",
										},
									},
								},
								{
									"resourceType": "AllergyIntolerance",
									"id":           "9012",
									"code": map[string]interface{}{
										"coding": []map[string]interface{}{
											{
												"code":    "161826",
												"display": "Allergy",
												"system":  "/orgs/CIEL/sources/CIEL/concepts/161826/",
											},
										},
										"text": "Allergy",
									},
									"reaction": []map[string]interface{}{
										{
											"manifestation": []map[string]interface{}{
												{
													"concept": "allergy concept",
												},
											},
										},
									},
									"clinicalStatus": map[string]interface{}{
										"text": "active",
									},
									"recordedDate": "2024-02-13T10:22:54+03:00",
								},
							},
							HasNextPage:     false,
							HasPreviousPage: false,
							NextCursor:      "",
							PreviousCursor:  "",
							TotalCount:      1,
						}, nil
					})
				mh.Mapper.EXPECT().ToTimeline(mock.Anything).
					RunAndReturn(func(resource interface{}) (*dto.TimelineResource, error) {
						return &dto.TimelineResource{
							ID:           gofakeit.UUID(),
							ResourceType: "Condition",
							Name:         "Condition",
							Value:        "active",
							Status:       "active",
							Date:         scalarutils.Date{Year: 2024, Month: 02, Day: 13},
							TimeRecorded: time.Now(),
							Category:     "",
						}, nil
					}).Once()
				mh.Mapper.EXPECT().ToTimeline(mock.Anything).
					RunAndReturn(func(resource interface{}) (*dto.TimelineResource, error) {
						return &dto.TimelineResource{
							ID:           gofakeit.UUID(),
							ResourceType: "MedicationStatement",
							Name:         "MedicationStatement",
							Value:        "active",
							Status:       "active",
							Date:         scalarutils.Date{Year: 2025, Month: 02, Day: 13},
							TimeRecorded: time.Now(),
							Category:     "",
						}, nil
					}).Once()
				mh.Mapper.EXPECT().ToTimeline(mock.Anything).
					RunAndReturn(func(resource interface{}) (*dto.TimelineResource, error) {
						return &dto.TimelineResource{
							ID:           gofakeit.UUID(),
							ResourceType: "AllergyIntolerance",
							Name:         "AllergyIntolerance",
							Value:        "active",
							Status:       "active",
							Date:         scalarutils.Date{Year: 2025, Month: 02, Day: 13},
							TimeRecorded: time.Now(),
							Category:     "",
						}, nil
					}).Once()

				return args{
					ctx:       ctx,
					patientID: gofakeit.UUID(),
					params: &dto.PatientEverythingFilterParams{
						Count: "3",
					},
				}
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to get patient banner",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						patientID := gofakeit.UUID()

						return &domain.FHIRPatientRelayPayload{
							Resource: &domain.FHIRPatient{
								ID: &patientID,
							},
						}, nil
					})
				mh.FHIR.EXPECT().GetFHIRPatientEverything(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string, params map[string]interface{}) (*domain.PagedFHIRResource, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{
					ctx:       ctx,
					patientID: gofakeit.UUID(),
					params: &dto.PatientEverythingFilterParams{
						Count: "3",
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

			_, err := clinicalUsecase.GetPatientBanner(args.ctx, args.patientID, args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("BaseImpl.GetPatientBanner() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestBaseImpl_GetPatientTimeline(t *testing.T) {
	type args struct {
		ctx       context.Context
		patientID string
		params    *dto.PatientEverythingFilterParams
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy case: get patient timeline",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						patientID := gofakeit.UUID()

						return &domain.FHIRPatientRelayPayload{
							Resource: &domain.FHIRPatient{
								ID: &patientID,
							},
						}, nil
					})
				mh.FHIR.EXPECT().GetFHIRPatientEverything(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string, params map[string]interface{}) (*domain.PagedFHIRResource, error) {
						return &domain.PagedFHIRResource{
							Resources: []map[string]interface{}{
								{
									"resourceType": "Condition",
									"id":           gofakeit.UUID(),
									"code": map[string]interface{}{
										"coding": []map[string]interface{}{
											{
												"code":    "161826",
												"display": "Biopsy of cervix",
												"system":  "/orgs/CIEL/sources/CIEL/concepts/161826/",
											},
										},
										"text": "Condition",
									},
									"clinicalStatus": map[string]interface{}{
										"text": "active",
									},
									"onsetDateTime": "2024-02-13T10:22:54+03:00",
									"category": []map[string]interface{}{
										{
											"coding": []map[string]interface{}{
												{
													"code":         "laboratory",
													"display":      "Laboratory",
													"system":       "http://terminology.hl7.org/CodeSystem/observation-category",
													"userSelected": false,
												},
											},
											"text": "Laboratory",
										},
									},
								},
								{
									"resourceType": "MedicationStatement",
									"id":           gofakeit.UUID(),
									"status":       "active",
									"code": map[string]interface{}{
										"coding": []map[string]interface{}{
											{
												"code":    "161826",
												"display": "Biopsy of cervix",
												"system":  "/orgs/CIEL/sources/CIEL/concepts/161826/",
											},
										},
										"text": "Biopsy of cervix",
									},
									"dateAsserted": "2025-02-13",
									"medication": []map[string]interface{}{
										{
											"concept": map[string]interface{}{
												"text": "medication concept",
											},
										},
									},
									"category": []map[string]interface{}{
										{
											"coding": []map[string]interface{}{
												{
													"code":         "laboratory",
													"display":      "Laboratory",
													"system":       "http://terminology.hl7.org/CodeSystem/observation-category",
													"userSelected": false,
												},
											},
											"text": "active",
										},
									},
								},
								{
									"resourceType": "AllergyIntolerance",
									"id":           "9012",
									"code": map[string]interface{}{
										"coding": []map[string]interface{}{
											{
												"code":    "161826",
												"display": "Allergy",
												"system":  "/orgs/CIEL/sources/CIEL/concepts/161826/",
											},
										},
										"text": "Allergy",
									},
									"reaction": []map[string]interface{}{
										{
											"manifestation": []map[string]interface{}{
												{
													"concept": "allergy concept",
												},
											},
										},
									},
									"clinicalStatus": map[string]interface{}{
										"text": "active",
									},
									"recordedDate": "2024-02-13T10:22:54+03:00",
								},
							},
							HasNextPage:     false,
							HasPreviousPage: false,
							NextCursor:      "",
							PreviousCursor:  "",
							TotalCount:      1,
						}, nil
					})
				mh.Mapper.EXPECT().ToTimeline(mock.Anything).
					RunAndReturn(func(resource interface{}) (*dto.TimelineResource, error) {
						return &dto.TimelineResource{
							ID:           gofakeit.UUID(),
							ResourceType: "Condition",
							Name:         "Condition",
							Value:        "active",
							Status:       "active",
							Date:         scalarutils.Date{Year: 2024, Month: 02, Day: 13},
							TimeRecorded: time.Now(),
							Category:     "",
						}, nil
					}).Once()
				mh.Mapper.EXPECT().ToTimeline(mock.Anything).
					RunAndReturn(func(resource interface{}) (*dto.TimelineResource, error) {
						return &dto.TimelineResource{
							ID:           gofakeit.UUID(),
							ResourceType: "MedicationStatement",
							Name:         "MedicationStatement",
							Value:        "active",
							Status:       "active",
							Date:         scalarutils.Date{Year: 2025, Month: 02, Day: 13},
							TimeRecorded: time.Now(),
							Category:     "",
						}, nil
					}).Once()
				mh.Mapper.EXPECT().ToTimeline(mock.Anything).
					RunAndReturn(func(resource interface{}) (*dto.TimelineResource, error) {
						return &dto.TimelineResource{
							ID:           gofakeit.UUID(),
							ResourceType: "AllergyIntolerance",
							Name:         "AllergyIntolerance",
							Value:        "active",
							Status:       "active",
							Date:         scalarutils.Date{Year: 2025, Month: 02, Day: 13},
							TimeRecorded: time.Now(),
							Category:     "",
						}, nil
					}).Once()

				return args{
					ctx:       usecaseMock.AddTenantIdentifierContext(context.Background()),
					patientID: gofakeit.UUID(),
					params: &dto.PatientEverythingFilterParams{
						Count: "3",
					},
				}
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to get patient timeline",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						patientID := gofakeit.UUID()

						return &domain.FHIRPatientRelayPayload{
							Resource: &domain.FHIRPatient{
								ID: &patientID,
							},
						}, nil
					})
				mh.FHIR.EXPECT().GetFHIRPatientEverything(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string, params map[string]interface{}) (*domain.PagedFHIRResource, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{
					ctx:       usecaseMock.AddTenantIdentifierContext(context.Background()),
					patientID: gofakeit.UUID(),
					params: &dto.PatientEverythingFilterParams{
						Count: "3",
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

			_, err := clinicalUsecase.GetPatientTimeline(args.ctx, args.patientID, args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("BaseImpl.GetPatientTimeline() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
