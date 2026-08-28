package base_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/scalarutils"
	"github.com/stretchr/testify/mock"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/empower-clinical/pkg/clinical/domain"
	usecaseMock "github.com/savannahghi/empower-clinical/pkg/clinical/usecases/mock"
)

func TestClinicalUseCaseImpl_GetMedicalData(t *testing.T) {
	tenantIDs := &dto.TenantIdentifiers{
		OrganizationID: uuid.NewString(),
		FacilityID:     uuid.NewString(),
	}

	type args struct {
		ctx       context.Context
		patientID string
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy case: patient timeline",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchFHIRMedicationStatement(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIRMedicationStatementRelayConnection, error) {
						status := dto.MedicationStatementStatusEnumActive
						code := "123"
						system := gofakeit.URL()

						return &domain.FHIRMedicationStatementRelayConnection{
							Edges: []*domain.FHIRMedicationStatementRelayEdge{
								{
									Cursor: new(string),
									Node: &domain.FHIRMedicationStatement{
										ID:     new(string),
										Status: (*domain.MedicationStatementStatusEnum)(&status),
										MedicationCodeableConcept: &domain.FHIRCodeableConcept{
											ID: new(string),
											Coding: []*domain.FHIRCoding{{
												ID:           new(string),
												System:       (*scalarutils.URI)(&system),
												Version:      new(string),
												Code:         (*scalarutils.Code)(&code),
												Display:      gofakeit.BS(),
												UserSelected: new(bool),
											}},
										},
										Subject: &domain.FHIRReference{
											ID: new(string),
										},
										Text:         &domain.FHIRNarrative{},
										Identifier:   []*domain.FHIRIdentifier{},
										BasedOn:      []*domain.FHIRReference{},
										PartOf:       []*domain.FHIRReference{},
										StatusReason: []*domain.FHIRCodeableConcept{},
										Category: []*domain.FHIRCodeableConcept{
											{
												ID: new(string),
												Coding: []*domain.FHIRCoding{
													{
														Code:    (*scalarutils.Code)(&code),
														Display: gofakeit.BS(),
													},
												},
												Text: "",
											},
										},
										Context:           &domain.FHIRReference{},
										EffectiveDateTime: &scalarutils.Date{},
										EffectivePeriod:   &domain.FHIRPeriod{},
										DateAsserted:      &scalarutils.Date{},
										InformationSource: &domain.FHIRReference{},
										DerivedFrom:       []*domain.FHIRReference{},
										ReasonCode:        []*domain.FHIRCodeableConcept{},
										ReasonReference:   []*domain.FHIRReference{},
										Note:              []*domain.FHIRAnnotation{},
										Dosage:            []*domain.FHIRDosage{},
										Meta:              &domain.FHIRMeta{},
										Extension:         []*domain.FHIRExtension{},
									},
								},
							},
							PageInfo: &firebasetools.PageInfo{},
						}, nil
					})

				mh.FHIR.EXPECT().SearchFHIRAllergyIntolerance(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRAllergy, error) {
						code := "123"
						system := gofakeit.URL()
						return &domain.PagedFHIRAllergy{
							Allergies: []domain.FHIRAllergyIntolerance{
								{
									ID: new(string),
									Code: &domain.FHIRCodeableConcept{
										ID: new(string),
										Coding: []*domain.FHIRCoding{
											{
												Code:    (*scalarutils.Code)(&code),
												Display: gofakeit.BS(),
												System:  (*scalarutils.URI)(&system),
											},
										},
										Text: gofakeit.BeerName(),
									},
									Patient: &domain.FHIRReference{
										ID: new(string),
									},
									Encounter: &domain.FHIRReference{
										ID: new(string),
									},
									OnsetPeriod: &domain.FHIRPeriod{
										ID:    new(string),
										Start: "2020-09-24T18:02:38.661033Z",
										End:   "2020-09-24T18:02:38.661033Z",
									},
									Reaction: []*domain.FHIRAllergyintoleranceReaction{},
								},
							},
							HasNextPage:     false,
							NextCursor:      "",
							HasPreviousPage: false,
							PreviousCursor:  "",
							TotalCount:      0,
						}, nil
					})
				mh.FHIR.EXPECT().SearchFHIRObservation(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRObservations, error) {
						status := dto.ObservationStatusFinal
						instant := gofakeit.TimeZone()
						valueConcept := "222"
						UUID := gofakeit.UUID()
						currentTime := time.Now().Format(time.RFC3339)
						return &domain.PagedFHIRObservations{
							Observations: []domain.FHIRObservation{
								{
									ID:     new(string),
									Status: (*domain.ObservationStatusEnum)(&status),
									Code: &domain.FHIRCodeableConcept{
										ID: new(string),
										Coding: []*domain.FHIRCoding{{
											Display: gofakeit.BS(),
											Code:    (*scalarutils.Code)(&valueConcept),
										}},
									},
									EffectiveDateTime: (*scalarutils.DateTime)(&currentTime),
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
						}, nil
					}).Times(4)

				return args{
					ctx:       context.Background(),
					patientID: gofakeit.UUID(),
				}
			},
			wantErr: false,
		},
		{
			name: "Sad case: failed to get tenant identifiers from context",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{
					ctx:       context.Background(),
					patientID: gofakeit.UUID(),
				}
			},
			wantErr: true,
		},
		{
			name: "Happy Case - Successfully search medication statement",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchFHIRMedicationStatement(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIRMedicationStatementRelayConnection, error) {
						status := dto.MedicationStatementStatusEnumActive
						code := "123"
						system := gofakeit.URL()

						return &domain.FHIRMedicationStatementRelayConnection{
							Edges: []*domain.FHIRMedicationStatementRelayEdge{
								{
									Cursor: new(string),
									Node: &domain.FHIRMedicationStatement{
										ID:     new(string),
										Status: (*domain.MedicationStatementStatusEnum)(&status),
										MedicationCodeableConcept: &domain.FHIRCodeableConcept{
											ID: new(string),
											Coding: []*domain.FHIRCoding{{
												ID:           new(string),
												System:       (*scalarutils.URI)(&system),
												Version:      new(string),
												Code:         (*scalarutils.Code)(&code),
												Display:      gofakeit.BS(),
												UserSelected: new(bool),
											}},
										},
										Subject: &domain.FHIRReference{
											ID: new(string),
										},
										Text:         &domain.FHIRNarrative{},
										Identifier:   []*domain.FHIRIdentifier{},
										BasedOn:      []*domain.FHIRReference{},
										PartOf:       []*domain.FHIRReference{},
										StatusReason: []*domain.FHIRCodeableConcept{},
										Category: []*domain.FHIRCodeableConcept{
											{
												ID: new(string),
												Coding: []*domain.FHIRCoding{
													{
														Code:    (*scalarutils.Code)(&code),
														Display: gofakeit.BS(),
													},
												},
												Text: "",
											},
										},
										Context:           &domain.FHIRReference{},
										EffectiveDateTime: &scalarutils.Date{},
										EffectivePeriod:   &domain.FHIRPeriod{},
										DateAsserted:      &scalarutils.Date{},
										InformationSource: &domain.FHIRReference{},
										DerivedFrom:       []*domain.FHIRReference{},
										ReasonCode:        []*domain.FHIRCodeableConcept{},
										ReasonReference:   []*domain.FHIRReference{},
										Note:              []*domain.FHIRAnnotation{},
										Dosage:            []*domain.FHIRDosage{},
										Meta:              &domain.FHIRMeta{},
										Extension:         []*domain.FHIRExtension{},
									},
								},
							},
							PageInfo: &firebasetools.PageInfo{},
						}, nil
					})

				mh.FHIR.EXPECT().SearchFHIRAllergyIntolerance(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRAllergy, error) {
						return &domain.PagedFHIRAllergy{
							Allergies: []domain.FHIRAllergyIntolerance{},
						}, nil
					})
				mh.FHIR.EXPECT().SearchFHIRObservation(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRObservations, error) {
						return &domain.PagedFHIRObservations{
							Observations: []domain.FHIRObservation{},
						}, nil
					})

				return args{
					ctx:       context.Background(),
					patientID: gofakeit.UUID(),
				}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to search medication statement",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchFHIRMedicationStatement(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIRMedicationStatementRelayConnection, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{
					ctx:       context.Background(),
					patientID: gofakeit.UUID(),
				}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to search medication statement - nil node",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchFHIRMedicationStatement(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIRMedicationStatementRelayConnection, error) {
						return &domain.FHIRMedicationStatementRelayConnection{
							Edges: []*domain.FHIRMedicationStatementRelayEdge{
								{
									Cursor: new(string),
								},
							},
							PageInfo: &firebasetools.PageInfo{},
						}, nil
					})

				mh.FHIR.EXPECT().SearchFHIRAllergyIntolerance(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRAllergy, error) {
						code := "123"
						system := gofakeit.URL()
						return &domain.PagedFHIRAllergy{
							Allergies: []domain.FHIRAllergyIntolerance{
								{
									ID: new(string),
									Code: &domain.FHIRCodeableConcept{
										ID: new(string),
										Coding: []*domain.FHIRCoding{
											{
												Code:    (*scalarutils.Code)(&code),
												Display: gofakeit.BS(),
												System:  (*scalarutils.URI)(&system),
											},
										},
										Text: gofakeit.BeerName(),
									},
									Patient: &domain.FHIRReference{
										ID: new(string),
									},
									Encounter: &domain.FHIRReference{
										ID: new(string),
									},
									OnsetPeriod: &domain.FHIRPeriod{
										ID:    new(string),
										Start: "2020-09-24T18:02:38.661033Z",
										End:   "2020-09-24T18:02:38.661033Z",
									},
									Reaction: []*domain.FHIRAllergyintoleranceReaction{},
								},
							},
							HasNextPage:     false,
							NextCursor:      "",
							HasPreviousPage: false,
							PreviousCursor:  "",
							TotalCount:      0,
						}, nil
					})
				mh.FHIR.EXPECT().SearchFHIRObservation(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRObservations, error) {
						status := dto.ObservationStatusFinal
						instant := gofakeit.TimeZone()
						valueConcept := "222"
						UUID := gofakeit.UUID()
						currentTime := time.Now().Format(time.RFC3339)
						return &domain.PagedFHIRObservations{
							Observations: []domain.FHIRObservation{
								{
									ID:     new(string),
									Status: (*domain.ObservationStatusEnum)(&status),
									Code: &domain.FHIRCodeableConcept{
										ID: new(string),
										Coding: []*domain.FHIRCoding{{
											Display: gofakeit.BS(),
											Code:    (*scalarutils.Code)(&valueConcept),
										}},
									},
									EffectiveDateTime: (*scalarutils.DateTime)(&currentTime),
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
						}, nil
					})

				return args{
					ctx:       context.Background(),
					patientID: gofakeit.UUID(),
				}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to search medication statement - nil node id",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchFHIRMedicationStatement(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIRMedicationStatementRelayConnection, error) {
						status := dto.MedicationStatementStatusEnumActive
						code := "123"

						return &domain.FHIRMedicationStatementRelayConnection{
							Edges: []*domain.FHIRMedicationStatementRelayEdge{
								{
									Cursor: new(string),
									Node: &domain.FHIRMedicationStatement{
										Status: (*domain.MedicationStatementStatusEnum)(&status),
										MedicationCodeableConcept: &domain.FHIRCodeableConcept{
											ID: new(string),
											Coding: []*domain.FHIRCoding{{
												Code:    (*scalarutils.Code)(&code),
												Display: gofakeit.BS(),
											}},
										},
										Subject: &domain.FHIRReference{
											ID: new(string),
										},
									},
								},
							},
							PageInfo: &firebasetools.PageInfo{},
						}, nil
					})

				mh.FHIR.EXPECT().SearchFHIRAllergyIntolerance(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRAllergy, error) {
						code := "123"
						system := gofakeit.URL()
						return &domain.PagedFHIRAllergy{
							Allergies: []domain.FHIRAllergyIntolerance{
								{
									ID: new(string),
									Code: &domain.FHIRCodeableConcept{
										ID: new(string),
										Coding: []*domain.FHIRCoding{
											{
												Code:    (*scalarutils.Code)(&code),
												Display: gofakeit.BS(),
												System:  (*scalarutils.URI)(&system),
											},
										},
										Text: gofakeit.BeerName(),
									},
									Patient: &domain.FHIRReference{
										ID: new(string),
									},
									Encounter: &domain.FHIRReference{
										ID: new(string),
									},
									OnsetPeriod: &domain.FHIRPeriod{
										ID:    new(string),
										Start: "2020-09-24T18:02:38.661033Z",
										End:   "2020-09-24T18:02:38.661033Z",
									},
									Reaction: []*domain.FHIRAllergyintoleranceReaction{},
								},
							},
							HasNextPage:     false,
							NextCursor:      "",
							HasPreviousPage: false,
							PreviousCursor:  "",
							TotalCount:      0,
						}, nil
					})
				mh.FHIR.EXPECT().SearchFHIRObservation(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRObservations, error) {
						status := dto.ObservationStatusFinal
						instant := gofakeit.TimeZone()
						valueConcept := "222"
						UUID := gofakeit.UUID()
						currentTime := time.Now().Format(time.RFC3339)
						return &domain.PagedFHIRObservations{
							Observations: []domain.FHIRObservation{
								{
									ID:     new(string),
									Status: (*domain.ObservationStatusEnum)(&status),
									Code: &domain.FHIRCodeableConcept{
										ID: new(string),
										Coding: []*domain.FHIRCoding{{
											Display: gofakeit.BS(),
											Code:    (*scalarutils.Code)(&valueConcept),
										}},
									},
									EffectiveDateTime: (*scalarutils.DateTime)(&currentTime),
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
						}, nil
					})

				return args{
					ctx:       context.Background(),
					patientID: gofakeit.UUID(),
				}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to search medication statement - nil status",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchFHIRMedicationStatement(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIRMedicationStatementRelayConnection, error) {
						code := "123"

						return &domain.FHIRMedicationStatementRelayConnection{
							Edges: []*domain.FHIRMedicationStatementRelayEdge{
								{
									Cursor: new(string),
									Node: &domain.FHIRMedicationStatement{
										ID: new(string),
										MedicationCodeableConcept: &domain.FHIRCodeableConcept{
											ID:     new(string),
											Coding: []*domain.FHIRCoding{{Code: (*scalarutils.Code)(&code), Display: gofakeit.BS()}},
											Text:   "",
										},
										Subject: &domain.FHIRReference{
											ID: new(string),
										},
									},
								},
							},
							PageInfo: &firebasetools.PageInfo{},
						}, nil
					})

				mh.FHIR.EXPECT().SearchFHIRAllergyIntolerance(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRAllergy, error) {
						code := "123"
						system := gofakeit.URL()
						return &domain.PagedFHIRAllergy{
							Allergies: []domain.FHIRAllergyIntolerance{
								{
									ID: new(string),
									Code: &domain.FHIRCodeableConcept{
										ID: new(string),
										Coding: []*domain.FHIRCoding{
											{
												Code:    (*scalarutils.Code)(&code),
												Display: gofakeit.BS(),
												System:  (*scalarutils.URI)(&system),
											},
										},
										Text: gofakeit.BeerName(),
									},
									Patient: &domain.FHIRReference{
										ID: new(string),
									},
									Encounter: &domain.FHIRReference{
										ID: new(string),
									},
									OnsetPeriod: &domain.FHIRPeriod{
										ID:    new(string),
										Start: "2020-09-24T18:02:38.661033Z",
										End:   "2020-09-24T18:02:38.661033Z",
									},
									Reaction: []*domain.FHIRAllergyintoleranceReaction{},
								},
							},
							HasNextPage:     false,
							NextCursor:      "",
							HasPreviousPage: false,
							PreviousCursor:  "",
							TotalCount:      0,
						}, nil
					})
				mh.FHIR.EXPECT().SearchFHIRObservation(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRObservations, error) {
						status := dto.ObservationStatusFinal
						instant := gofakeit.TimeZone()
						valueConcept := "222"
						UUID := gofakeit.UUID()
						currentTime := time.Now().Format(time.RFC3339)
						return &domain.PagedFHIRObservations{
							Observations: []domain.FHIRObservation{
								{
									ID:     new(string),
									Status: (*domain.ObservationStatusEnum)(&status),
									Code: &domain.FHIRCodeableConcept{
										ID: new(string),
										Coding: []*domain.FHIRCoding{{
											Display: gofakeit.BS(),
											Code:    (*scalarutils.Code)(&valueConcept),
										}},
									},
									EffectiveDateTime: (*scalarutils.DateTime)(&currentTime),
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
						}, nil
					})

				return args{
					ctx:       context.Background(),
					patientID: gofakeit.UUID(),
				}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to search medication statement - nil coding",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchFHIRMedicationStatement(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIRMedicationStatementRelayConnection, error) {
						code := "123"

						return &domain.FHIRMedicationStatementRelayConnection{
							Edges: []*domain.FHIRMedicationStatementRelayEdge{
								{
									Cursor: new(string),
									Node: &domain.FHIRMedicationStatement{
										ID: new(string),
										MedicationCodeableConcept: &domain.FHIRCodeableConcept{
											ID:     new(string),
											Coding: []*domain.FHIRCoding{{Code: (*scalarutils.Code)(&code), Display: gofakeit.BS()}},
											Text:   "",
										},
										Subject: &domain.FHIRReference{
											ID: new(string),
										},
									},
								},
							},
							PageInfo: &firebasetools.PageInfo{},
						}, nil
					})

				mh.FHIR.EXPECT().SearchFHIRAllergyIntolerance(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRAllergy, error) {
						code := "123"
						system := gofakeit.URL()
						return &domain.PagedFHIRAllergy{
							Allergies: []domain.FHIRAllergyIntolerance{
								{
									ID: new(string),
									Code: &domain.FHIRCodeableConcept{
										ID: new(string),
										Coding: []*domain.FHIRCoding{
											{
												Code:    (*scalarutils.Code)(&code),
												Display: gofakeit.BS(),
												System:  (*scalarutils.URI)(&system),
											},
										},
										Text: gofakeit.BeerName(),
									},
									Patient: &domain.FHIRReference{
										ID: new(string),
									},
									Encounter: &domain.FHIRReference{
										ID: new(string),
									},
									OnsetPeriod: &domain.FHIRPeriod{
										ID:    new(string),
										Start: "2020-09-24T18:02:38.661033Z",
										End:   "2020-09-24T18:02:38.661033Z",
									},
									Reaction: []*domain.FHIRAllergyintoleranceReaction{},
								},
							},
							HasNextPage:     false,
							NextCursor:      "",
							HasPreviousPage: false,
							PreviousCursor:  "",
							TotalCount:      0,
						}, nil
					})
				mh.FHIR.EXPECT().SearchFHIRObservation(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRObservations, error) {
						status := dto.ObservationStatusFinal
						instant := gofakeit.TimeZone()
						valueConcept := "222"
						UUID := gofakeit.UUID()
						currentTime := time.Now().Format(time.RFC3339)
						return &domain.PagedFHIRObservations{
							Observations: []domain.FHIRObservation{
								{
									ID:     new(string),
									Status: (*domain.ObservationStatusEnum)(&status),
									Code: &domain.FHIRCodeableConcept{
										ID: new(string),
										Coding: []*domain.FHIRCoding{{
											Display: gofakeit.BS(),
											Code:    (*scalarutils.Code)(&valueConcept),
										}},
									},
									EffectiveDateTime: (*scalarutils.DateTime)(&currentTime),
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
						}, nil
					})

				return args{
					ctx:       context.Background(),
					patientID: gofakeit.UUID(),
				}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to search medication statement - empty coding",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchFHIRMedicationStatement(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIRMedicationStatementRelayConnection, error) {
						status := dto.MedicationStatementStatusEnumActive

						return &domain.FHIRMedicationStatementRelayConnection{
							Edges: []*domain.FHIRMedicationStatementRelayEdge{
								{
									Cursor: new(string),
									Node: &domain.FHIRMedicationStatement{
										ID:     new(string),
										Status: (*domain.MedicationStatementStatusEnum)(&status),
										MedicationCodeableConcept: &domain.FHIRCodeableConcept{
											ID:     new(string),
											Coding: []*domain.FHIRCoding{},
										},
										Subject: &domain.FHIRReference{
											ID: new(string),
										},
									},
								},
							},
							PageInfo: &firebasetools.PageInfo{},
						}, nil
					})

				mh.FHIR.EXPECT().SearchFHIRAllergyIntolerance(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRAllergy, error) {
						code := "123"
						system := gofakeit.URL()
						return &domain.PagedFHIRAllergy{
							Allergies: []domain.FHIRAllergyIntolerance{
								{
									ID: new(string),
									Code: &domain.FHIRCodeableConcept{
										ID: new(string),
										Coding: []*domain.FHIRCoding{
											{
												Code:    (*scalarutils.Code)(&code),
												Display: gofakeit.BS(),
												System:  (*scalarutils.URI)(&system),
											},
										},
										Text: gofakeit.BeerName(),
									},
									Patient: &domain.FHIRReference{
										ID: new(string),
									},
									Encounter: &domain.FHIRReference{
										ID: new(string),
									},
									OnsetPeriod: &domain.FHIRPeriod{
										ID:    new(string),
										Start: "2020-09-24T18:02:38.661033Z",
										End:   "2020-09-24T18:02:38.661033Z",
									},
									Reaction: []*domain.FHIRAllergyintoleranceReaction{},
								},
							},
							HasNextPage:     false,
							NextCursor:      "",
							HasPreviousPage: false,
							PreviousCursor:  "",
							TotalCount:      0,
						}, nil
					})
				mh.FHIR.EXPECT().SearchFHIRObservation(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRObservations, error) {
						status := dto.ObservationStatusFinal
						instant := gofakeit.TimeZone()
						valueConcept := "222"
						UUID := gofakeit.UUID()
						currentTime := time.Now().Format(time.RFC3339)
						return &domain.PagedFHIRObservations{
							Observations: []domain.FHIRObservation{
								{
									ID:     new(string),
									Status: (*domain.ObservationStatusEnum)(&status),
									Code: &domain.FHIRCodeableConcept{
										ID: new(string),
										Coding: []*domain.FHIRCoding{{
											Display: gofakeit.BS(),
											Code:    (*scalarutils.Code)(&valueConcept),
										}},
									},
									EffectiveDateTime: (*scalarutils.DateTime)(&currentTime),
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
						}, nil
					})

				return args{
					ctx:       context.Background(),
					patientID: gofakeit.UUID(),
				}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to search medication statement - nil subject id",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchFHIRMedicationStatement(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIRMedicationStatementRelayConnection, error) {
						code := "123"
						status := dto.MedicationStatementStatusEnumActive

						return &domain.FHIRMedicationStatementRelayConnection{
							Edges: []*domain.FHIRMedicationStatementRelayEdge{
								{
									Cursor: new(string),
									Node: &domain.FHIRMedicationStatement{
										ID:     new(string),
										Status: (*domain.MedicationStatementStatusEnum)(&status),
										MedicationCodeableConcept: &domain.FHIRCodeableConcept{
											ID: new(string),
											Coding: []*domain.FHIRCoding{{
												Code:    (*scalarutils.Code)(&code),
												Display: gofakeit.BS(),
											}},
										},
										Subject: &domain.FHIRReference{},
									},
								},
							},
							PageInfo: &firebasetools.PageInfo{},
						}, nil
					})

				mh.FHIR.EXPECT().SearchFHIRAllergyIntolerance(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRAllergy, error) {
						code := "123"
						system := gofakeit.URL()
						return &domain.PagedFHIRAllergy{
							Allergies: []domain.FHIRAllergyIntolerance{
								{
									ID: new(string),
									Code: &domain.FHIRCodeableConcept{
										ID: new(string),
										Coding: []*domain.FHIRCoding{
											{
												Code:    (*scalarutils.Code)(&code),
												Display: gofakeit.BS(),
												System:  (*scalarutils.URI)(&system),
											},
										},
										Text: gofakeit.BeerName(),
									},
									Patient: &domain.FHIRReference{
										ID: new(string),
									},
									Encounter: &domain.FHIRReference{
										ID: new(string),
									},
									OnsetPeriod: &domain.FHIRPeriod{
										ID:    new(string),
										Start: "2020-09-24T18:02:38.661033Z",
										End:   "2020-09-24T18:02:38.661033Z",
									},
									Reaction: []*domain.FHIRAllergyintoleranceReaction{},
								},
							},
							HasNextPage:     false,
							NextCursor:      "",
							HasPreviousPage: false,
							PreviousCursor:  "",
							TotalCount:      0,
						}, nil
					})
				mh.FHIR.EXPECT().SearchFHIRObservation(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRObservations, error) {
						status := dto.ObservationStatusFinal
						instant := gofakeit.TimeZone()
						valueConcept := "222"
						UUID := gofakeit.UUID()
						currentTime := time.Now().Format(time.RFC3339)
						return &domain.PagedFHIRObservations{
							Observations: []domain.FHIRObservation{
								{
									ID:     new(string),
									Status: (*domain.ObservationStatusEnum)(&status),
									Code: &domain.FHIRCodeableConcept{
										ID: new(string),
										Coding: []*domain.FHIRCoding{{
											Display: gofakeit.BS(),
											Code:    (*scalarutils.Code)(&valueConcept),
										}},
									},
									EffectiveDateTime: (*scalarutils.DateTime)(&currentTime),
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
						}, nil
					})

				return args{
					ctx:       context.Background(),
					patientID: gofakeit.UUID(),
				}
			},
			wantErr: false,
		},
		{
			name: "Happy Case - Successfully search allergy intolerance",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchFHIRMedicationStatement(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIRMedicationStatementRelayConnection, error) {
						return &domain.FHIRMedicationStatementRelayConnection{
							Edges: []*domain.FHIRMedicationStatementRelayEdge{},
						}, nil
					})

				mh.FHIR.EXPECT().SearchFHIRAllergyIntolerance(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRAllergy, error) {
						code := "123"
						system := gofakeit.URL()
						return &domain.PagedFHIRAllergy{
							Allergies: []domain.FHIRAllergyIntolerance{
								{
									ID: new(string),
									Code: &domain.FHIRCodeableConcept{
										ID: new(string),
										Coding: []*domain.FHIRCoding{
											{
												Code:    (*scalarutils.Code)(&code),
												Display: gofakeit.BS(),
												System:  (*scalarutils.URI)(&system),
											},
										},
										Text: gofakeit.BeerName(),
									},
									Patient: &domain.FHIRReference{
										ID: new(string),
									},
									Encounter: &domain.FHIRReference{
										ID: new(string),
									},
									OnsetPeriod: &domain.FHIRPeriod{
										ID:    new(string),
										Start: "2020-09-24T18:02:38.661033Z",
										End:   "2020-09-24T18:02:38.661033Z",
									},
									Reaction: []*domain.FHIRAllergyintoleranceReaction{},
								},
							},
							HasNextPage:     false,
							NextCursor:      "",
							HasPreviousPage: false,
							PreviousCursor:  "",
							TotalCount:      0,
						}, nil
					})
				mh.FHIR.EXPECT().SearchFHIRObservation(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRObservations, error) {
						return &domain.PagedFHIRObservations{
							Observations: []domain.FHIRObservation{},
						}, nil
					})

				return args{
					ctx:       context.Background(),
					patientID: gofakeit.UUID(),
				}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to search allergy intolerance",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchFHIRMedicationStatement(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIRMedicationStatementRelayConnection, error) {
						return &domain.FHIRMedicationStatementRelayConnection{
							Edges: []*domain.FHIRMedicationStatementRelayEdge{},
						}, nil
					})

				mh.FHIR.EXPECT().SearchFHIRAllergyIntolerance(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRAllergy, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{
					ctx:       context.Background(),
					patientID: gofakeit.UUID(),
				}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to search allergy intolerance - nil node",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchFHIRMedicationStatement(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIRMedicationStatementRelayConnection, error) {
						status := dto.MedicationStatementStatusEnumActive
						code := "123"
						system := gofakeit.URL()

						return &domain.FHIRMedicationStatementRelayConnection{
							Edges: []*domain.FHIRMedicationStatementRelayEdge{
								{
									Cursor: new(string),
									Node: &domain.FHIRMedicationStatement{
										ID:     new(string),
										Status: (*domain.MedicationStatementStatusEnum)(&status),
										MedicationCodeableConcept: &domain.FHIRCodeableConcept{
											ID: new(string),
											Coding: []*domain.FHIRCoding{{
												ID:           new(string),
												System:       (*scalarutils.URI)(&system),
												Version:      new(string),
												Code:         (*scalarutils.Code)(&code),
												Display:      gofakeit.BS(),
												UserSelected: new(bool),
											}},
										},
										Subject: &domain.FHIRReference{
											ID: new(string),
										},
										Text:         &domain.FHIRNarrative{},
										Identifier:   []*domain.FHIRIdentifier{},
										BasedOn:      []*domain.FHIRReference{},
										PartOf:       []*domain.FHIRReference{},
										StatusReason: []*domain.FHIRCodeableConcept{},
										Category: []*domain.FHIRCodeableConcept{
											{
												ID: new(string),
												Coding: []*domain.FHIRCoding{
													{
														Code:    (*scalarutils.Code)(&code),
														Display: gofakeit.BS(),
													},
												},
												Text: "",
											},
										},
										Context:           &domain.FHIRReference{},
										EffectiveDateTime: &scalarutils.Date{},
										EffectivePeriod:   &domain.FHIRPeriod{},
										DateAsserted:      &scalarutils.Date{},
										InformationSource: &domain.FHIRReference{},
										DerivedFrom:       []*domain.FHIRReference{},
										ReasonCode:        []*domain.FHIRCodeableConcept{},
										ReasonReference:   []*domain.FHIRReference{},
										Note:              []*domain.FHIRAnnotation{},
										Dosage:            []*domain.FHIRDosage{},
										Meta:              &domain.FHIRMeta{},
										Extension:         []*domain.FHIRExtension{},
									},
								},
							},
							PageInfo: &firebasetools.PageInfo{},
						}, nil
					})

				mh.FHIR.EXPECT().SearchFHIRAllergyIntolerance(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRAllergy, error) {
						return &domain.PagedFHIRAllergy{
							Allergies: []domain.FHIRAllergyIntolerance{
								{
									ID: new(string),
								},
							},
							HasNextPage:     false,
							NextCursor:      "",
							HasPreviousPage: false,
							PreviousCursor:  "",
							TotalCount:      0,
						}, nil
					})
				mh.FHIR.EXPECT().SearchFHIRObservation(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRObservations, error) {
						status := dto.ObservationStatusFinal
						instant := gofakeit.TimeZone()
						valueConcept := "222"
						UUID := gofakeit.UUID()
						currentTime := time.Now().Format(time.RFC3339)
						return &domain.PagedFHIRObservations{
							Observations: []domain.FHIRObservation{
								{
									ID:     new(string),
									Status: (*domain.ObservationStatusEnum)(&status),
									Code: &domain.FHIRCodeableConcept{
										ID: new(string),
										Coding: []*domain.FHIRCoding{{
											Display: gofakeit.BS(),
											Code:    (*scalarutils.Code)(&valueConcept),
										}},
									},
									EffectiveDateTime: (*scalarutils.DateTime)(&currentTime),
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
						}, nil
					})

				return args{
					ctx:       context.Background(),
					patientID: gofakeit.UUID(),
				}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to search allergy intolerance - nil node id",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchFHIRMedicationStatement(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIRMedicationStatementRelayConnection, error) {
						status := dto.MedicationStatementStatusEnumActive
						code := "123"
						system := gofakeit.URL()

						return &domain.FHIRMedicationStatementRelayConnection{
							Edges: []*domain.FHIRMedicationStatementRelayEdge{
								{
									Cursor: new(string),
									Node: &domain.FHIRMedicationStatement{
										ID:     new(string),
										Status: (*domain.MedicationStatementStatusEnum)(&status),
										MedicationCodeableConcept: &domain.FHIRCodeableConcept{
											ID: new(string),
											Coding: []*domain.FHIRCoding{{
												ID:           new(string),
												System:       (*scalarutils.URI)(&system),
												Version:      new(string),
												Code:         (*scalarutils.Code)(&code),
												Display:      gofakeit.BS(),
												UserSelected: new(bool),
											}},
										},
										Subject: &domain.FHIRReference{
											ID: new(string),
										},
										Text:         &domain.FHIRNarrative{},
										Identifier:   []*domain.FHIRIdentifier{},
										BasedOn:      []*domain.FHIRReference{},
										PartOf:       []*domain.FHIRReference{},
										StatusReason: []*domain.FHIRCodeableConcept{},
										Category: []*domain.FHIRCodeableConcept{
											{
												ID: new(string),
												Coding: []*domain.FHIRCoding{
													{
														Code:    (*scalarutils.Code)(&code),
														Display: gofakeit.BS(),
													},
												},
												Text: "",
											},
										},
										Context:           &domain.FHIRReference{},
										EffectiveDateTime: &scalarutils.Date{},
										EffectivePeriod:   &domain.FHIRPeriod{},
										DateAsserted:      &scalarutils.Date{},
										InformationSource: &domain.FHIRReference{},
										DerivedFrom:       []*domain.FHIRReference{},
										ReasonCode:        []*domain.FHIRCodeableConcept{},
										ReasonReference:   []*domain.FHIRReference{},
										Note:              []*domain.FHIRAnnotation{},
										Dosage:            []*domain.FHIRDosage{},
										Meta:              &domain.FHIRMeta{},
										Extension:         []*domain.FHIRExtension{},
									},
								},
							},
							PageInfo: &firebasetools.PageInfo{},
						}, nil
					})

				mh.FHIR.EXPECT().SearchFHIRAllergyIntolerance(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRAllergy, error) {
						code := "123"
						system := gofakeit.URL()

						return &domain.PagedFHIRAllergy{
							Allergies: []domain.FHIRAllergyIntolerance{
								{
									Code: &domain.FHIRCodeableConcept{
										ID: new(string),
										Coding: []*domain.FHIRCoding{
											{
												Code:    (*scalarutils.Code)(&code),
												Display: gofakeit.BS(),
												System:  (*scalarutils.URI)(&system),
											},
										},
									},
									Patient: &domain.FHIRReference{
										ID: new(string),
									},
									Encounter: &domain.FHIRReference{
										ID: new(string),
									},
									OnsetPeriod: &domain.FHIRPeriod{
										Start: "2000-01-01",
									},
								},
							},
							HasNextPage:     false,
							NextCursor:      "",
							HasPreviousPage: false,
							PreviousCursor:  "",
							TotalCount:      0,
						}, nil
					})
				mh.FHIR.EXPECT().SearchFHIRObservation(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRObservations, error) {
						status := dto.ObservationStatusFinal
						instant := gofakeit.TimeZone()
						valueConcept := "222"
						UUID := gofakeit.UUID()
						currentTime := time.Now().Format(time.RFC3339)
						return &domain.PagedFHIRObservations{
							Observations: []domain.FHIRObservation{
								{
									ID:     new(string),
									Status: (*domain.ObservationStatusEnum)(&status),
									Code: &domain.FHIRCodeableConcept{
										ID: new(string),
										Coding: []*domain.FHIRCoding{{
											Display: gofakeit.BS(),
											Code:    (*scalarutils.Code)(&valueConcept),
										}},
									},
									EffectiveDateTime: (*scalarutils.DateTime)(&currentTime),
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
						}, nil
					})

				return args{
					ctx:       context.Background(),
					patientID: gofakeit.UUID(),
				}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to search allergy intolerance - nil patient",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchFHIRMedicationStatement(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIRMedicationStatementRelayConnection, error) {
						status := dto.MedicationStatementStatusEnumActive
						code := "123"
						system := gofakeit.URL()

						return &domain.FHIRMedicationStatementRelayConnection{
							Edges: []*domain.FHIRMedicationStatementRelayEdge{
								{
									Cursor: new(string),
									Node: &domain.FHIRMedicationStatement{
										ID:     new(string),
										Status: (*domain.MedicationStatementStatusEnum)(&status),
										MedicationCodeableConcept: &domain.FHIRCodeableConcept{
											ID: new(string),
											Coding: []*domain.FHIRCoding{{
												ID:           new(string),
												System:       (*scalarutils.URI)(&system),
												Version:      new(string),
												Code:         (*scalarutils.Code)(&code),
												Display:      gofakeit.BS(),
												UserSelected: new(bool),
											}},
										},
										Subject: &domain.FHIRReference{
											ID: new(string),
										},
										Text:         &domain.FHIRNarrative{},
										Identifier:   []*domain.FHIRIdentifier{},
										BasedOn:      []*domain.FHIRReference{},
										PartOf:       []*domain.FHIRReference{},
										StatusReason: []*domain.FHIRCodeableConcept{},
										Category: []*domain.FHIRCodeableConcept{
											{
												ID: new(string),
												Coding: []*domain.FHIRCoding{
													{
														Code:    (*scalarutils.Code)(&code),
														Display: gofakeit.BS(),
													},
												},
												Text: "",
											},
										},
										Context:           &domain.FHIRReference{},
										EffectiveDateTime: &scalarutils.Date{},
										EffectivePeriod:   &domain.FHIRPeriod{},
										DateAsserted:      &scalarutils.Date{},
										InformationSource: &domain.FHIRReference{},
										DerivedFrom:       []*domain.FHIRReference{},
										ReasonCode:        []*domain.FHIRCodeableConcept{},
										ReasonReference:   []*domain.FHIRReference{},
										Note:              []*domain.FHIRAnnotation{},
										Dosage:            []*domain.FHIRDosage{},
										Meta:              &domain.FHIRMeta{},
										Extension:         []*domain.FHIRExtension{},
									},
								},
							},
							PageInfo: &firebasetools.PageInfo{},
						}, nil
					})

				mh.FHIR.EXPECT().SearchFHIRAllergyIntolerance(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRAllergy, error) {
						code := "123"
						system := gofakeit.URL()

						return &domain.PagedFHIRAllergy{
							Allergies: []domain.FHIRAllergyIntolerance{
								{
									Code: &domain.FHIRCodeableConcept{
										ID: new(string),
										Coding: []*domain.FHIRCoding{
											{
												Code:    (*scalarutils.Code)(&code),
												Display: gofakeit.BS(),
												System:  (*scalarutils.URI)(&system),
											},
										},
									},
									Encounter: &domain.FHIRReference{
										ID: new(string),
									},
									OnsetPeriod: &domain.FHIRPeriod{
										Start: "2000-01-01",
									},
								},
							},
							HasNextPage:     false,
							NextCursor:      "",
							HasPreviousPage: false,
							PreviousCursor:  "",
							TotalCount:      0,
						}, nil
					})
				mh.FHIR.EXPECT().SearchFHIRObservation(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRObservations, error) {
						status := dto.ObservationStatusFinal
						instant := gofakeit.TimeZone()
						valueConcept := "222"
						UUID := gofakeit.UUID()
						currentTime := time.Now().Format(time.RFC3339)
						return &domain.PagedFHIRObservations{
							Observations: []domain.FHIRObservation{
								{
									ID:     new(string),
									Status: (*domain.ObservationStatusEnum)(&status),
									Code: &domain.FHIRCodeableConcept{
										ID: new(string),
										Coding: []*domain.FHIRCoding{{
											Display: gofakeit.BS(),
											Code:    (*scalarutils.Code)(&valueConcept),
										}},
									},
									EffectiveDateTime: (*scalarutils.DateTime)(&currentTime),
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
						}, nil
					})

				return args{
					ctx:       context.Background(),
					patientID: gofakeit.UUID(),
				}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to search allergy intolerance - nil patient id",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchFHIRMedicationStatement(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIRMedicationStatementRelayConnection, error) {
						status := dto.MedicationStatementStatusEnumActive
						code := "123"
						system := gofakeit.URL()

						return &domain.FHIRMedicationStatementRelayConnection{
							Edges: []*domain.FHIRMedicationStatementRelayEdge{
								{
									Cursor: new(string),
									Node: &domain.FHIRMedicationStatement{
										ID:     new(string),
										Status: (*domain.MedicationStatementStatusEnum)(&status),
										MedicationCodeableConcept: &domain.FHIRCodeableConcept{
											ID: new(string),
											Coding: []*domain.FHIRCoding{{
												ID:           new(string),
												System:       (*scalarutils.URI)(&system),
												Version:      new(string),
												Code:         (*scalarutils.Code)(&code),
												Display:      gofakeit.BS(),
												UserSelected: new(bool),
											}},
										},
										Subject: &domain.FHIRReference{
											ID: new(string),
										},
										Text:         &domain.FHIRNarrative{},
										Identifier:   []*domain.FHIRIdentifier{},
										BasedOn:      []*domain.FHIRReference{},
										PartOf:       []*domain.FHIRReference{},
										StatusReason: []*domain.FHIRCodeableConcept{},
										Category: []*domain.FHIRCodeableConcept{
											{
												ID: new(string),
												Coding: []*domain.FHIRCoding{
													{
														Code:    (*scalarutils.Code)(&code),
														Display: gofakeit.BS(),
													},
												},
												Text: "",
											},
										},
										Context:           &domain.FHIRReference{},
										EffectiveDateTime: &scalarutils.Date{},
										EffectivePeriod:   &domain.FHIRPeriod{},
										DateAsserted:      &scalarutils.Date{},
										InformationSource: &domain.FHIRReference{},
										DerivedFrom:       []*domain.FHIRReference{},
										ReasonCode:        []*domain.FHIRCodeableConcept{},
										ReasonReference:   []*domain.FHIRReference{},
										Note:              []*domain.FHIRAnnotation{},
										Dosage:            []*domain.FHIRDosage{},
										Meta:              &domain.FHIRMeta{},
										Extension:         []*domain.FHIRExtension{},
									},
								},
							},
							PageInfo: &firebasetools.PageInfo{},
						}, nil
					})

				mh.FHIR.EXPECT().SearchFHIRAllergyIntolerance(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRAllergy, error) {
						code := "123"
						system := gofakeit.URL()

						return &domain.PagedFHIRAllergy{
							Allergies: []domain.FHIRAllergyIntolerance{
								{
									Code: &domain.FHIRCodeableConcept{
										ID: new(string),
										Coding: []*domain.FHIRCoding{
											{
												Code:    (*scalarutils.Code)(&code),
												Display: gofakeit.BS(),
												System:  (*scalarutils.URI)(&system),
											},
										},
									},
									Patient: &domain.FHIRReference{
										ID: nil,
									},
									Encounter: &domain.FHIRReference{
										ID: new(string),
									},
									OnsetPeriod: &domain.FHIRPeriod{
										Start: "2000-01-01",
									},
								},
							},
							HasNextPage:     false,
							NextCursor:      "",
							HasPreviousPage: false,
							PreviousCursor:  "",
							TotalCount:      0,
						}, nil
					})
				mh.FHIR.EXPECT().SearchFHIRObservation(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRObservations, error) {
						status := dto.ObservationStatusFinal
						instant := gofakeit.TimeZone()
						valueConcept := "222"
						UUID := gofakeit.UUID()
						currentTime := time.Now().Format(time.RFC3339)
						return &domain.PagedFHIRObservations{
							Observations: []domain.FHIRObservation{
								{
									ID:     new(string),
									Status: (*domain.ObservationStatusEnum)(&status),
									Code: &domain.FHIRCodeableConcept{
										ID: new(string),
										Coding: []*domain.FHIRCoding{{
											Display: gofakeit.BS(),
											Code:    (*scalarutils.Code)(&valueConcept),
										}},
									},
									EffectiveDateTime: (*scalarutils.DateTime)(&currentTime),
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
						}, nil
					})

				return args{
					ctx:       context.Background(),
					patientID: gofakeit.UUID(),
				}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to search allergy intolerance - nil encounter",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchFHIRMedicationStatement(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIRMedicationStatementRelayConnection, error) {
						status := dto.MedicationStatementStatusEnumActive
						code := "123"
						system := gofakeit.URL()

						return &domain.FHIRMedicationStatementRelayConnection{
							Edges: []*domain.FHIRMedicationStatementRelayEdge{
								{
									Cursor: new(string),
									Node: &domain.FHIRMedicationStatement{
										ID:     new(string),
										Status: (*domain.MedicationStatementStatusEnum)(&status),
										MedicationCodeableConcept: &domain.FHIRCodeableConcept{
											ID: new(string),
											Coding: []*domain.FHIRCoding{{
												ID:           new(string),
												System:       (*scalarutils.URI)(&system),
												Version:      new(string),
												Code:         (*scalarutils.Code)(&code),
												Display:      gofakeit.BS(),
												UserSelected: new(bool),
											}},
										},
										Subject: &domain.FHIRReference{
											ID: new(string),
										},
										Text:         &domain.FHIRNarrative{},
										Identifier:   []*domain.FHIRIdentifier{},
										BasedOn:      []*domain.FHIRReference{},
										PartOf:       []*domain.FHIRReference{},
										StatusReason: []*domain.FHIRCodeableConcept{},
										Category: []*domain.FHIRCodeableConcept{
											{
												ID: new(string),
												Coding: []*domain.FHIRCoding{
													{
														Code:    (*scalarutils.Code)(&code),
														Display: gofakeit.BS(),
													},
												},
												Text: "",
											},
										},
										Context:           &domain.FHIRReference{},
										EffectiveDateTime: &scalarutils.Date{},
										EffectivePeriod:   &domain.FHIRPeriod{},
										DateAsserted:      &scalarutils.Date{},
										InformationSource: &domain.FHIRReference{},
										DerivedFrom:       []*domain.FHIRReference{},
										ReasonCode:        []*domain.FHIRCodeableConcept{},
										ReasonReference:   []*domain.FHIRReference{},
										Note:              []*domain.FHIRAnnotation{},
										Dosage:            []*domain.FHIRDosage{},
										Meta:              &domain.FHIRMeta{},
										Extension:         []*domain.FHIRExtension{},
									},
								},
							},
							PageInfo: &firebasetools.PageInfo{},
						}, nil
					})

				mh.FHIR.EXPECT().SearchFHIRAllergyIntolerance(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRAllergy, error) {
						code := "123"
						system := gofakeit.URL()

						return &domain.PagedFHIRAllergy{
							Allergies: []domain.FHIRAllergyIntolerance{
								{
									Code: &domain.FHIRCodeableConcept{
										ID: new(string),
										Coding: []*domain.FHIRCoding{
											{
												Code:    (*scalarutils.Code)(&code),
												Display: gofakeit.BS(),
												System:  (*scalarutils.URI)(&system),
											},
										},
									},
									Patient: &domain.FHIRReference{
										ID: new(string),
									},
									OnsetPeriod: &domain.FHIRPeriod{
										Start: "2000-01-01",
									},
								},
							},
							HasNextPage:     false,
							NextCursor:      "",
							HasPreviousPage: false,
							PreviousCursor:  "",
							TotalCount:      0,
						}, nil
					})
				mh.FHIR.EXPECT().SearchFHIRObservation(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRObservations, error) {
						status := dto.ObservationStatusFinal
						instant := gofakeit.TimeZone()
						valueConcept := "222"
						UUID := gofakeit.UUID()
						currentTime := time.Now().Format(time.RFC3339)
						return &domain.PagedFHIRObservations{
							Observations: []domain.FHIRObservation{
								{
									ID:     new(string),
									Status: (*domain.ObservationStatusEnum)(&status),
									Code: &domain.FHIRCodeableConcept{
										ID: new(string),
										Coding: []*domain.FHIRCoding{{
											Display: gofakeit.BS(),
											Code:    (*scalarutils.Code)(&valueConcept),
										}},
									},
									EffectiveDateTime: (*scalarutils.DateTime)(&currentTime),
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
						}, nil
					})

				return args{
					ctx:       context.Background(),
					patientID: gofakeit.UUID(),
				}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to search allergy intolerance - nil encounter id",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchFHIRMedicationStatement(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIRMedicationStatementRelayConnection, error) {
						status := dto.MedicationStatementStatusEnumActive
						code := "123"
						system := gofakeit.URL()

						return &domain.FHIRMedicationStatementRelayConnection{
							Edges: []*domain.FHIRMedicationStatementRelayEdge{
								{
									Cursor: new(string),
									Node: &domain.FHIRMedicationStatement{
										ID:     new(string),
										Status: (*domain.MedicationStatementStatusEnum)(&status),
										MedicationCodeableConcept: &domain.FHIRCodeableConcept{
											ID: new(string),
											Coding: []*domain.FHIRCoding{{
												ID:           new(string),
												System:       (*scalarutils.URI)(&system),
												Version:      new(string),
												Code:         (*scalarutils.Code)(&code),
												Display:      gofakeit.BS(),
												UserSelected: new(bool),
											}},
										},
										Subject: &domain.FHIRReference{
											ID: new(string),
										},
										Text:         &domain.FHIRNarrative{},
										Identifier:   []*domain.FHIRIdentifier{},
										BasedOn:      []*domain.FHIRReference{},
										PartOf:       []*domain.FHIRReference{},
										StatusReason: []*domain.FHIRCodeableConcept{},
										Category: []*domain.FHIRCodeableConcept{
											{
												ID: new(string),
												Coding: []*domain.FHIRCoding{
													{
														Code:    (*scalarutils.Code)(&code),
														Display: gofakeit.BS(),
													},
												},
												Text: "",
											},
										},
										Context:           &domain.FHIRReference{},
										EffectiveDateTime: &scalarutils.Date{},
										EffectivePeriod:   &domain.FHIRPeriod{},
										DateAsserted:      &scalarutils.Date{},
										InformationSource: &domain.FHIRReference{},
										DerivedFrom:       []*domain.FHIRReference{},
										ReasonCode:        []*domain.FHIRCodeableConcept{},
										ReasonReference:   []*domain.FHIRReference{},
										Note:              []*domain.FHIRAnnotation{},
										Dosage:            []*domain.FHIRDosage{},
										Meta:              &domain.FHIRMeta{},
										Extension:         []*domain.FHIRExtension{},
									},
								},
							},
							PageInfo: &firebasetools.PageInfo{},
						}, nil
					})

				mh.FHIR.EXPECT().SearchFHIRAllergyIntolerance(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRAllergy, error) {
						code := "123"
						system := gofakeit.URL()

						return &domain.PagedFHIRAllergy{
							Allergies: []domain.FHIRAllergyIntolerance{
								{
									Code: &domain.FHIRCodeableConcept{
										ID: new(string),
										Coding: []*domain.FHIRCoding{
											{
												Code:    (*scalarutils.Code)(&code),
												Display: gofakeit.BS(),
												System:  (*scalarutils.URI)(&system),
											},
										},
									},
									Patient: &domain.FHIRReference{
										ID: new(string),
									},
									Encounter: &domain.FHIRReference{
										ID: new(string),
									},
									OnsetPeriod: &domain.FHIRPeriod{
										Start: "2000-01-01",
									},
								},
							},
							HasNextPage:     false,
							NextCursor:      "",
							HasPreviousPage: false,
							PreviousCursor:  "",
							TotalCount:      0,
						}, nil
					})
				mh.FHIR.EXPECT().SearchFHIRObservation(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRObservations, error) {
						status := dto.ObservationStatusFinal
						instant := gofakeit.TimeZone()
						valueConcept := "222"
						UUID := gofakeit.UUID()
						currentTime := time.Now().Format(time.RFC3339)
						return &domain.PagedFHIRObservations{
							Observations: []domain.FHIRObservation{
								{
									ID:     new(string),
									Status: (*domain.ObservationStatusEnum)(&status),
									Code: &domain.FHIRCodeableConcept{
										ID: new(string),
										Coding: []*domain.FHIRCoding{{
											Display: gofakeit.BS(),
											Code:    (*scalarutils.Code)(&valueConcept),
										}},
									},
									EffectiveDateTime: (*scalarutils.DateTime)(&currentTime),
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
						}, nil
					})

				return args{
					ctx:       context.Background(),
					patientID: gofakeit.UUID(),
				}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to search allergy intolerance - nil code",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchFHIRMedicationStatement(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIRMedicationStatementRelayConnection, error) {
						status := dto.MedicationStatementStatusEnumActive
						code := "123"
						system := gofakeit.URL()

						return &domain.FHIRMedicationStatementRelayConnection{
							Edges: []*domain.FHIRMedicationStatementRelayEdge{
								{
									Cursor: new(string),
									Node: &domain.FHIRMedicationStatement{
										ID:     new(string),
										Status: (*domain.MedicationStatementStatusEnum)(&status),
										MedicationCodeableConcept: &domain.FHIRCodeableConcept{
											ID: new(string),
											Coding: []*domain.FHIRCoding{{
												ID:           new(string),
												System:       (*scalarutils.URI)(&system),
												Version:      new(string),
												Code:         (*scalarutils.Code)(&code),
												Display:      gofakeit.BS(),
												UserSelected: new(bool),
											}},
										},
										Subject: &domain.FHIRReference{
											ID: new(string),
										},
										Text:         &domain.FHIRNarrative{},
										Identifier:   []*domain.FHIRIdentifier{},
										BasedOn:      []*domain.FHIRReference{},
										PartOf:       []*domain.FHIRReference{},
										StatusReason: []*domain.FHIRCodeableConcept{},
										Category: []*domain.FHIRCodeableConcept{
											{
												ID: new(string),
												Coding: []*domain.FHIRCoding{
													{
														Code:    (*scalarutils.Code)(&code),
														Display: gofakeit.BS(),
													},
												},
												Text: "",
											},
										},
										Context:           &domain.FHIRReference{},
										EffectiveDateTime: &scalarutils.Date{},
										EffectivePeriod:   &domain.FHIRPeriod{},
										DateAsserted:      &scalarutils.Date{},
										InformationSource: &domain.FHIRReference{},
										DerivedFrom:       []*domain.FHIRReference{},
										ReasonCode:        []*domain.FHIRCodeableConcept{},
										ReasonReference:   []*domain.FHIRReference{},
										Note:              []*domain.FHIRAnnotation{},
										Dosage:            []*domain.FHIRDosage{},
										Meta:              &domain.FHIRMeta{},
										Extension:         []*domain.FHIRExtension{},
									},
								},
							},
							PageInfo: &firebasetools.PageInfo{},
						}, nil
					})

				mh.FHIR.EXPECT().SearchFHIRAllergyIntolerance(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRAllergy, error) {
						return &domain.PagedFHIRAllergy{
							Allergies: []domain.FHIRAllergyIntolerance{
								{
									Patient: &domain.FHIRReference{
										ID: new(string),
									},
									Encounter: &domain.FHIRReference{
										ID: new(string),
									},
									OnsetPeriod: &domain.FHIRPeriod{
										Start: "2000-01-01",
									},
								},
							},
							HasNextPage:     false,
							NextCursor:      "",
							HasPreviousPage: false,
							PreviousCursor:  "",
							TotalCount:      0,
						}, nil
					})
				mh.FHIR.EXPECT().SearchFHIRObservation(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRObservations, error) {
						status := dto.ObservationStatusFinal
						instant := gofakeit.TimeZone()
						valueConcept := "222"
						UUID := gofakeit.UUID()
						currentTime := time.Now().Format(time.RFC3339)
						return &domain.PagedFHIRObservations{
							Observations: []domain.FHIRObservation{
								{
									ID:     new(string),
									Status: (*domain.ObservationStatusEnum)(&status),
									Code: &domain.FHIRCodeableConcept{
										ID: new(string),
										Coding: []*domain.FHIRCoding{{
											Display: gofakeit.BS(),
											Code:    (*scalarutils.Code)(&valueConcept),
										}},
									},
									EffectiveDateTime: (*scalarutils.DateTime)(&currentTime),
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
						}, nil
					})

				return args{
					ctx:       context.Background(),
					patientID: gofakeit.UUID(),
				}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to search allergy intolerance - nil coding",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchFHIRMedicationStatement(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIRMedicationStatementRelayConnection, error) {
						status := dto.MedicationStatementStatusEnumActive
						code := "123"
						system := gofakeit.URL()

						return &domain.FHIRMedicationStatementRelayConnection{
							Edges: []*domain.FHIRMedicationStatementRelayEdge{
								{
									Cursor: new(string),
									Node: &domain.FHIRMedicationStatement{
										ID:     new(string),
										Status: (*domain.MedicationStatementStatusEnum)(&status),
										MedicationCodeableConcept: &domain.FHIRCodeableConcept{
											ID: new(string),
											Coding: []*domain.FHIRCoding{{
												ID:           new(string),
												System:       (*scalarutils.URI)(&system),
												Version:      new(string),
												Code:         (*scalarutils.Code)(&code),
												Display:      gofakeit.BS(),
												UserSelected: new(bool),
											}},
										},
										Subject: &domain.FHIRReference{
											ID: new(string),
										},
										Text:         &domain.FHIRNarrative{},
										Identifier:   []*domain.FHIRIdentifier{},
										BasedOn:      []*domain.FHIRReference{},
										PartOf:       []*domain.FHIRReference{},
										StatusReason: []*domain.FHIRCodeableConcept{},
										Category: []*domain.FHIRCodeableConcept{
											{
												ID: new(string),
												Coding: []*domain.FHIRCoding{
													{
														Code:    (*scalarutils.Code)(&code),
														Display: gofakeit.BS(),
													},
												},
												Text: "",
											},
										},
										Context:           &domain.FHIRReference{},
										EffectiveDateTime: &scalarutils.Date{},
										EffectivePeriod:   &domain.FHIRPeriod{},
										DateAsserted:      &scalarutils.Date{},
										InformationSource: &domain.FHIRReference{},
										DerivedFrom:       []*domain.FHIRReference{},
										ReasonCode:        []*domain.FHIRCodeableConcept{},
										ReasonReference:   []*domain.FHIRReference{},
										Note:              []*domain.FHIRAnnotation{},
										Dosage:            []*domain.FHIRDosage{},
										Meta:              &domain.FHIRMeta{},
										Extension:         []*domain.FHIRExtension{},
									},
								},
							},
							PageInfo: &firebasetools.PageInfo{},
						}, nil
					})

				mh.FHIR.EXPECT().SearchFHIRAllergyIntolerance(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRAllergy, error) {
						return &domain.PagedFHIRAllergy{
							Allergies: []domain.FHIRAllergyIntolerance{
								{
									Code: &domain.FHIRCodeableConcept{
										ID: new(string),
									},
									Patient: &domain.FHIRReference{
										ID: new(string),
									},
									Encounter: &domain.FHIRReference{
										ID: new(string),
									},
									OnsetPeriod: &domain.FHIRPeriod{
										Start: "2000-01-01",
									},
								},
							},
							HasNextPage:     false,
							NextCursor:      "",
							HasPreviousPage: false,
							PreviousCursor:  "",
							TotalCount:      0,
						}, nil
					})
				mh.FHIR.EXPECT().SearchFHIRObservation(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRObservations, error) {
						status := dto.ObservationStatusFinal
						instant := gofakeit.TimeZone()
						valueConcept := "222"
						UUID := gofakeit.UUID()
						currentTime := time.Now().Format(time.RFC3339)
						return &domain.PagedFHIRObservations{
							Observations: []domain.FHIRObservation{
								{
									ID:     new(string),
									Status: (*domain.ObservationStatusEnum)(&status),
									Code: &domain.FHIRCodeableConcept{
										ID: new(string),
										Coding: []*domain.FHIRCoding{{
											Display: gofakeit.BS(),
											Code:    (*scalarutils.Code)(&valueConcept),
										}},
									},
									EffectiveDateTime: (*scalarutils.DateTime)(&currentTime),
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
						}, nil
					})

				return args{
					ctx:       context.Background(),
					patientID: gofakeit.UUID(),
				}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to search allergy intolerance - coding length < 1",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchFHIRMedicationStatement(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIRMedicationStatementRelayConnection, error) {
						status := dto.MedicationStatementStatusEnumActive
						code := "123"
						system := gofakeit.URL()

						return &domain.FHIRMedicationStatementRelayConnection{
							Edges: []*domain.FHIRMedicationStatementRelayEdge{
								{
									Cursor: new(string),
									Node: &domain.FHIRMedicationStatement{
										ID:     new(string),
										Status: (*domain.MedicationStatementStatusEnum)(&status),
										MedicationCodeableConcept: &domain.FHIRCodeableConcept{
											ID: new(string),
											Coding: []*domain.FHIRCoding{{
												ID:           new(string),
												System:       (*scalarutils.URI)(&system),
												Version:      new(string),
												Code:         (*scalarutils.Code)(&code),
												Display:      gofakeit.BS(),
												UserSelected: new(bool),
											}},
										},
										Subject: &domain.FHIRReference{
											ID: new(string),
										},
										Text:         &domain.FHIRNarrative{},
										Identifier:   []*domain.FHIRIdentifier{},
										BasedOn:      []*domain.FHIRReference{},
										PartOf:       []*domain.FHIRReference{},
										StatusReason: []*domain.FHIRCodeableConcept{},
										Category: []*domain.FHIRCodeableConcept{
											{
												ID: new(string),
												Coding: []*domain.FHIRCoding{
													{
														Code:    (*scalarutils.Code)(&code),
														Display: gofakeit.BS(),
													},
												},
												Text: "",
											},
										},
										Context:           &domain.FHIRReference{},
										EffectiveDateTime: &scalarutils.Date{},
										EffectivePeriod:   &domain.FHIRPeriod{},
										DateAsserted:      &scalarutils.Date{},
										InformationSource: &domain.FHIRReference{},
										DerivedFrom:       []*domain.FHIRReference{},
										ReasonCode:        []*domain.FHIRCodeableConcept{},
										ReasonReference:   []*domain.FHIRReference{},
										Note:              []*domain.FHIRAnnotation{},
										Dosage:            []*domain.FHIRDosage{},
										Meta:              &domain.FHIRMeta{},
										Extension:         []*domain.FHIRExtension{},
									},
								},
							},
							PageInfo: &firebasetools.PageInfo{},
						}, nil
					})

				mh.FHIR.EXPECT().SearchFHIRAllergyIntolerance(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRAllergy, error) {
						return &domain.PagedFHIRAllergy{
							Allergies: []domain.FHIRAllergyIntolerance{
								{
									Code: &domain.FHIRCodeableConcept{
										ID:     new(string),
										Coding: []*domain.FHIRCoding{},
									},
									Patient: &domain.FHIRReference{
										ID: new(string),
									},
									Encounter: &domain.FHIRReference{
										ID: new(string),
									},
									OnsetPeriod: &domain.FHIRPeriod{
										Start: "2000-01-01",
									},
								},
							},
							HasNextPage:     false,
							NextCursor:      "",
							HasPreviousPage: false,
							PreviousCursor:  "",
							TotalCount:      0,
						}, nil
					})
				mh.FHIR.EXPECT().SearchFHIRObservation(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRObservations, error) {
						status := dto.ObservationStatusFinal
						instant := gofakeit.TimeZone()
						valueConcept := "222"
						UUID := gofakeit.UUID()
						currentTime := time.Now().Format(time.RFC3339)
						return &domain.PagedFHIRObservations{
							Observations: []domain.FHIRObservation{
								{
									ID:     new(string),
									Status: (*domain.ObservationStatusEnum)(&status),
									Code: &domain.FHIRCodeableConcept{
										ID: new(string),
										Coding: []*domain.FHIRCoding{{
											Display: gofakeit.BS(),
											Code:    (*scalarutils.Code)(&valueConcept),
										}},
									},
									EffectiveDateTime: (*scalarutils.DateTime)(&currentTime),
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
						}, nil
					})

				return args{
					ctx:       context.Background(),
					patientID: gofakeit.UUID(),
				}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to search allergy intolerance - nil reaction",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchFHIRMedicationStatement(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIRMedicationStatementRelayConnection, error) {
						status := dto.MedicationStatementStatusEnumActive
						code := "123"
						system := gofakeit.URL()

						return &domain.FHIRMedicationStatementRelayConnection{
							Edges: []*domain.FHIRMedicationStatementRelayEdge{
								{
									Cursor: new(string),
									Node: &domain.FHIRMedicationStatement{
										ID:     new(string),
										Status: (*domain.MedicationStatementStatusEnum)(&status),
										MedicationCodeableConcept: &domain.FHIRCodeableConcept{
											ID: new(string),
											Coding: []*domain.FHIRCoding{{
												ID:           new(string),
												System:       (*scalarutils.URI)(&system),
												Version:      new(string),
												Code:         (*scalarutils.Code)(&code),
												Display:      gofakeit.BS(),
												UserSelected: new(bool),
											}},
										},
										Subject: &domain.FHIRReference{
											ID: new(string),
										},
										Text:         &domain.FHIRNarrative{},
										Identifier:   []*domain.FHIRIdentifier{},
										BasedOn:      []*domain.FHIRReference{},
										PartOf:       []*domain.FHIRReference{},
										StatusReason: []*domain.FHIRCodeableConcept{},
										Category: []*domain.FHIRCodeableConcept{
											{
												ID: new(string),
												Coding: []*domain.FHIRCoding{
													{
														Code:    (*scalarutils.Code)(&code),
														Display: gofakeit.BS(),
													},
												},
												Text: "",
											},
										},
										Context:           &domain.FHIRReference{},
										EffectiveDateTime: &scalarutils.Date{},
										EffectivePeriod:   &domain.FHIRPeriod{},
										DateAsserted:      &scalarutils.Date{},
										InformationSource: &domain.FHIRReference{},
										DerivedFrom:       []*domain.FHIRReference{},
										ReasonCode:        []*domain.FHIRCodeableConcept{},
										ReasonReference:   []*domain.FHIRReference{},
										Note:              []*domain.FHIRAnnotation{},
										Dosage:            []*domain.FHIRDosage{},
										Meta:              &domain.FHIRMeta{},
										Extension:         []*domain.FHIRExtension{},
									},
								},
							},
							PageInfo: &firebasetools.PageInfo{},
						}, nil
					})

				mh.FHIR.EXPECT().SearchFHIRAllergyIntolerance(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRAllergy, error) {
						code := "123"
						system := gofakeit.URL()

						return &domain.PagedFHIRAllergy{
							Allergies: []domain.FHIRAllergyIntolerance{
								{
									Code: &domain.FHIRCodeableConcept{
										ID: new(string),
										Coding: []*domain.FHIRCoding{
											{
												Code:    (*scalarutils.Code)(&code),
												Display: gofakeit.BS(),
												System:  (*scalarutils.URI)(&system),
											},
										},
									},
									Patient: &domain.FHIRReference{
										ID: new(string),
									},
									Encounter: &domain.FHIRReference{
										ID: new(string),
									},
									OnsetPeriod: &domain.FHIRPeriod{
										Start: "2000-01-01",
									},
								},
							},
							HasNextPage:     false,
							NextCursor:      "",
							HasPreviousPage: false,
							PreviousCursor:  "",
							TotalCount:      0,
						}, nil
					})
				mh.FHIR.EXPECT().SearchFHIRObservation(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRObservations, error) {
						status := dto.ObservationStatusFinal
						instant := gofakeit.TimeZone()
						valueConcept := "222"
						UUID := gofakeit.UUID()
						currentTime := time.Now().Format(time.RFC3339)
						return &domain.PagedFHIRObservations{
							Observations: []domain.FHIRObservation{
								{
									ID:     new(string),
									Status: (*domain.ObservationStatusEnum)(&status),
									Code: &domain.FHIRCodeableConcept{
										ID: new(string),
										Coding: []*domain.FHIRCoding{{
											Display: gofakeit.BS(),
											Code:    (*scalarutils.Code)(&valueConcept),
										}},
									},
									EffectiveDateTime: (*scalarutils.DateTime)(&currentTime),
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
						}, nil
					})

				return args{
					ctx:       context.Background(),
					patientID: gofakeit.UUID(),
				}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to search allergy intolerance - reaction length < 1",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchFHIRMedicationStatement(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIRMedicationStatementRelayConnection, error) {
						status := dto.MedicationStatementStatusEnumActive
						code := "123"
						system := gofakeit.URL()

						return &domain.FHIRMedicationStatementRelayConnection{
							Edges: []*domain.FHIRMedicationStatementRelayEdge{
								{
									Cursor: new(string),
									Node: &domain.FHIRMedicationStatement{
										ID:     new(string),
										Status: (*domain.MedicationStatementStatusEnum)(&status),
										MedicationCodeableConcept: &domain.FHIRCodeableConcept{
											ID: new(string),
											Coding: []*domain.FHIRCoding{{
												ID:           new(string),
												System:       (*scalarutils.URI)(&system),
												Version:      new(string),
												Code:         (*scalarutils.Code)(&code),
												Display:      gofakeit.BS(),
												UserSelected: new(bool),
											}},
										},
										Subject: &domain.FHIRReference{
											ID: new(string),
										},
										Text:         &domain.FHIRNarrative{},
										Identifier:   []*domain.FHIRIdentifier{},
										BasedOn:      []*domain.FHIRReference{},
										PartOf:       []*domain.FHIRReference{},
										StatusReason: []*domain.FHIRCodeableConcept{},
										Category: []*domain.FHIRCodeableConcept{
											{
												ID: new(string),
												Coding: []*domain.FHIRCoding{
													{
														Code:    (*scalarutils.Code)(&code),
														Display: gofakeit.BS(),
													},
												},
												Text: "",
											},
										},
										Context:           &domain.FHIRReference{},
										EffectiveDateTime: &scalarutils.Date{},
										EffectivePeriod:   &domain.FHIRPeriod{},
										DateAsserted:      &scalarutils.Date{},
										InformationSource: &domain.FHIRReference{},
										DerivedFrom:       []*domain.FHIRReference{},
										ReasonCode:        []*domain.FHIRCodeableConcept{},
										ReasonReference:   []*domain.FHIRReference{},
										Note:              []*domain.FHIRAnnotation{},
										Dosage:            []*domain.FHIRDosage{},
										Meta:              &domain.FHIRMeta{},
										Extension:         []*domain.FHIRExtension{},
									},
								},
							},
							PageInfo: &firebasetools.PageInfo{},
						}, nil
					})

				mh.FHIR.EXPECT().SearchFHIRAllergyIntolerance(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRAllergy, error) {
						code := "123"
						system := gofakeit.URL()

						return &domain.PagedFHIRAllergy{
							Allergies: []domain.FHIRAllergyIntolerance{
								{
									Code: &domain.FHIRCodeableConcept{
										ID: new(string),
										Coding: []*domain.FHIRCoding{
											{
												Code:    (*scalarutils.Code)(&code),
												Display: gofakeit.BS(),
												System:  (*scalarutils.URI)(&system),
											},
										},
									},
									Patient: &domain.FHIRReference{
										ID: new(string),
									},
									Encounter: &domain.FHIRReference{
										ID: new(string),
									},
									Reaction: []*domain.FHIRAllergyintoleranceReaction{},
									OnsetPeriod: &domain.FHIRPeriod{
										Start: "2000-01-01",
									},
								},
							},
							HasNextPage:     false,
							NextCursor:      "",
							HasPreviousPage: false,
							PreviousCursor:  "",
							TotalCount:      0,
						}, nil
					})
				mh.FHIR.EXPECT().SearchFHIRObservation(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRObservations, error) {
						status := dto.ObservationStatusFinal
						instant := gofakeit.TimeZone()
						valueConcept := "222"
						UUID := gofakeit.UUID()
						currentTime := time.Now().Format(time.RFC3339)
						return &domain.PagedFHIRObservations{
							Observations: []domain.FHIRObservation{
								{
									ID:     new(string),
									Status: (*domain.ObservationStatusEnum)(&status),
									Code: &domain.FHIRCodeableConcept{
										ID: new(string),
										Coding: []*domain.FHIRCoding{{
											Display: gofakeit.BS(),
											Code:    (*scalarutils.Code)(&valueConcept),
										}},
									},
									EffectiveDateTime: (*scalarutils.DateTime)(&currentTime),
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
						}, nil
					})

				return args{
					ctx:       context.Background(),
					patientID: gofakeit.UUID(),
				}
			},
			wantErr: false,
		},
		{
			name: "Happy Case - Successfully search observation",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchFHIRMedicationStatement(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIRMedicationStatementRelayConnection, error) {
						return &domain.FHIRMedicationStatementRelayConnection{
							Edges: []*domain.FHIRMedicationStatementRelayEdge{},
						}, nil
					})

				mh.FHIR.EXPECT().SearchFHIRAllergyIntolerance(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRAllergy, error) {
						return &domain.PagedFHIRAllergy{
							Allergies: []domain.FHIRAllergyIntolerance{},
						}, nil
					})
				mh.FHIR.EXPECT().SearchFHIRObservation(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRObservations, error) {
						status := dto.ObservationStatusFinal
						instant := gofakeit.TimeZone()
						valueConcept := "222"
						UUID := gofakeit.UUID()
						currentTime := time.Now().Format(time.RFC3339)
						return &domain.PagedFHIRObservations{
							Observations: []domain.FHIRObservation{
								{
									ID:     new(string),
									Status: (*domain.ObservationStatusEnum)(&status),
									Code: &domain.FHIRCodeableConcept{
										ID: new(string),
										Coding: []*domain.FHIRCoding{{
											Display: gofakeit.BS(),
											Code:    (*scalarutils.Code)(&valueConcept),
										}},
									},
									EffectiveDateTime: (*scalarutils.DateTime)(&currentTime),
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
						}, nil
					})

				return args{
					ctx:       context.Background(),
					patientID: gofakeit.UUID(),
				}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to search observation - nil id",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchFHIRMedicationStatement(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIRMedicationStatementRelayConnection, error) {
						status := dto.MedicationStatementStatusEnumActive
						code := "123"
						system := gofakeit.URL()

						return &domain.FHIRMedicationStatementRelayConnection{
							Edges: []*domain.FHIRMedicationStatementRelayEdge{
								{
									Cursor: new(string),
									Node: &domain.FHIRMedicationStatement{
										ID:     new(string),
										Status: (*domain.MedicationStatementStatusEnum)(&status),
										MedicationCodeableConcept: &domain.FHIRCodeableConcept{
											ID: new(string),
											Coding: []*domain.FHIRCoding{{
												ID:           new(string),
												System:       (*scalarutils.URI)(&system),
												Version:      new(string),
												Code:         (*scalarutils.Code)(&code),
												Display:      gofakeit.BS(),
												UserSelected: new(bool),
											}},
										},
										Subject: &domain.FHIRReference{
											ID: new(string),
										},
										Text:         &domain.FHIRNarrative{},
										Identifier:   []*domain.FHIRIdentifier{},
										BasedOn:      []*domain.FHIRReference{},
										PartOf:       []*domain.FHIRReference{},
										StatusReason: []*domain.FHIRCodeableConcept{},
										Category: []*domain.FHIRCodeableConcept{
											{
												ID: new(string),
												Coding: []*domain.FHIRCoding{
													{
														Code:    (*scalarutils.Code)(&code),
														Display: gofakeit.BS(),
													},
												},
												Text: "",
											},
										},
										Context:           &domain.FHIRReference{},
										EffectiveDateTime: &scalarutils.Date{},
										EffectivePeriod:   &domain.FHIRPeriod{},
										DateAsserted:      &scalarutils.Date{},
										InformationSource: &domain.FHIRReference{},
										DerivedFrom:       []*domain.FHIRReference{},
										ReasonCode:        []*domain.FHIRCodeableConcept{},
										ReasonReference:   []*domain.FHIRReference{},
										Note:              []*domain.FHIRAnnotation{},
										Dosage:            []*domain.FHIRDosage{},
										Meta:              &domain.FHIRMeta{},
										Extension:         []*domain.FHIRExtension{},
									},
								},
							},
							PageInfo: &firebasetools.PageInfo{},
						}, nil
					})

				mh.FHIR.EXPECT().SearchFHIRAllergyIntolerance(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRAllergy, error) {
						code := "123"
						system := gofakeit.URL()

						return &domain.PagedFHIRAllergy{
							Allergies: []domain.FHIRAllergyIntolerance{
								{
									Code: &domain.FHIRCodeableConcept{
										ID: new(string),
										Coding: []*domain.FHIRCoding{
											{
												Code:    (*scalarutils.Code)(&code),
												Display: gofakeit.BS(),
												System:  (*scalarutils.URI)(&system),
											},
										},
									},
									Patient: &domain.FHIRReference{
										ID: new(string),
									},
									Encounter: &domain.FHIRReference{
										ID: new(string),
									},
									OnsetPeriod: &domain.FHIRPeriod{
										Start: "2000-01-01",
									},
								},
							},
							HasNextPage:     false,
							NextCursor:      "",
							HasPreviousPage: false,
							PreviousCursor:  "",
							TotalCount:      0,
						}, nil
					})
				mh.FHIR.EXPECT().SearchFHIRObservation(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRObservations, error) {
						status := dto.ObservationStatusFinal
						valueConcept := "222"
						UUID := gofakeit.UUID()
						currentTime := time.Now().Format(time.RFC3339)
						return &domain.PagedFHIRObservations{
							Observations: []domain.FHIRObservation{
								{
									Status: (*domain.ObservationStatusEnum)(&status),
									Code: &domain.FHIRCodeableConcept{
										ID: new(string),
									},
									EffectiveDateTime: (*scalarutils.DateTime)(&currentTime),
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
						}, nil
					})

				return args{
					ctx:       context.Background(),
					patientID: gofakeit.UUID(),
				}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to search observation - nil coding",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchFHIRMedicationStatement(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIRMedicationStatementRelayConnection, error) {
						status := dto.MedicationStatementStatusEnumActive
						code := "123"
						system := gofakeit.URL()

						return &domain.FHIRMedicationStatementRelayConnection{
							Edges: []*domain.FHIRMedicationStatementRelayEdge{
								{
									Cursor: new(string),
									Node: &domain.FHIRMedicationStatement{
										ID:     new(string),
										Status: (*domain.MedicationStatementStatusEnum)(&status),
										MedicationCodeableConcept: &domain.FHIRCodeableConcept{
											ID: new(string),
											Coding: []*domain.FHIRCoding{{
												ID:           new(string),
												System:       (*scalarutils.URI)(&system),
												Version:      new(string),
												Code:         (*scalarutils.Code)(&code),
												Display:      gofakeit.BS(),
												UserSelected: new(bool),
											}},
										},
										Subject: &domain.FHIRReference{
											ID: new(string),
										},
										Text:         &domain.FHIRNarrative{},
										Identifier:   []*domain.FHIRIdentifier{},
										BasedOn:      []*domain.FHIRReference{},
										PartOf:       []*domain.FHIRReference{},
										StatusReason: []*domain.FHIRCodeableConcept{},
										Category: []*domain.FHIRCodeableConcept{
											{
												ID: new(string),
												Coding: []*domain.FHIRCoding{
													{
														Code:    (*scalarutils.Code)(&code),
														Display: gofakeit.BS(),
													},
												},
												Text: "",
											},
										},
										Context:           &domain.FHIRReference{},
										EffectiveDateTime: &scalarutils.Date{},
										EffectivePeriod:   &domain.FHIRPeriod{},
										DateAsserted:      &scalarutils.Date{},
										InformationSource: &domain.FHIRReference{},
										DerivedFrom:       []*domain.FHIRReference{},
										ReasonCode:        []*domain.FHIRCodeableConcept{},
										ReasonReference:   []*domain.FHIRReference{},
										Note:              []*domain.FHIRAnnotation{},
										Dosage:            []*domain.FHIRDosage{},
										Meta:              &domain.FHIRMeta{},
										Extension:         []*domain.FHIRExtension{},
									},
								},
							},
							PageInfo: &firebasetools.PageInfo{},
						}, nil
					})

				mh.FHIR.EXPECT().SearchFHIRAllergyIntolerance(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRAllergy, error) {
						code := "123"
						system := gofakeit.URL()

						return &domain.PagedFHIRAllergy{
							Allergies: []domain.FHIRAllergyIntolerance{
								{
									Code: &domain.FHIRCodeableConcept{
										ID: new(string),
										Coding: []*domain.FHIRCoding{
											{
												Code:    (*scalarutils.Code)(&code),
												Display: gofakeit.BS(),
												System:  (*scalarutils.URI)(&system),
											},
										},
									},
									Patient: &domain.FHIRReference{
										ID: new(string),
									},
									Encounter: &domain.FHIRReference{
										ID: new(string),
									},
									OnsetPeriod: &domain.FHIRPeriod{
										Start: "2000-01-01",
									},
								},
							},
							HasNextPage:     false,
							NextCursor:      "",
							HasPreviousPage: false,
							PreviousCursor:  "",
							TotalCount:      0,
						}, nil
					})
				mh.FHIR.EXPECT().SearchFHIRObservation(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRObservations, error) {
						status := dto.ObservationStatusFinal
						valueConcept := "222"
						UUID := gofakeit.UUID()
						currentTime := time.Now().Format(time.RFC3339)
						return &domain.PagedFHIRObservations{
							Observations: []domain.FHIRObservation{
								{
									Status: (*domain.ObservationStatusEnum)(&status),
									Code: &domain.FHIRCodeableConcept{
										ID: new(string),
									},
									EffectiveDateTime: (*scalarutils.DateTime)(&currentTime),
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
						}, nil
					})

				return args{
					ctx:       context.Background(),
					patientID: gofakeit.UUID(),
				}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to search observation - empty coding",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchFHIRMedicationStatement(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIRMedicationStatementRelayConnection, error) {
						status := dto.MedicationStatementStatusEnumActive
						code := "123"
						system := gofakeit.URL()

						return &domain.FHIRMedicationStatementRelayConnection{
							Edges: []*domain.FHIRMedicationStatementRelayEdge{
								{
									Cursor: new(string),
									Node: &domain.FHIRMedicationStatement{
										ID:     new(string),
										Status: (*domain.MedicationStatementStatusEnum)(&status),
										MedicationCodeableConcept: &domain.FHIRCodeableConcept{
											ID: new(string),
											Coding: []*domain.FHIRCoding{{
												ID:           new(string),
												System:       (*scalarutils.URI)(&system),
												Version:      new(string),
												Code:         (*scalarutils.Code)(&code),
												Display:      gofakeit.BS(),
												UserSelected: new(bool),
											}},
										},
										Subject: &domain.FHIRReference{
											ID: new(string),
										},
										Text:         &domain.FHIRNarrative{},
										Identifier:   []*domain.FHIRIdentifier{},
										BasedOn:      []*domain.FHIRReference{},
										PartOf:       []*domain.FHIRReference{},
										StatusReason: []*domain.FHIRCodeableConcept{},
										Category: []*domain.FHIRCodeableConcept{
											{
												ID: new(string),
												Coding: []*domain.FHIRCoding{
													{
														Code:    (*scalarutils.Code)(&code),
														Display: gofakeit.BS(),
													},
												},
												Text: "",
											},
										},
										Context:           &domain.FHIRReference{},
										EffectiveDateTime: &scalarutils.Date{},
										EffectivePeriod:   &domain.FHIRPeriod{},
										DateAsserted:      &scalarutils.Date{},
										InformationSource: &domain.FHIRReference{},
										DerivedFrom:       []*domain.FHIRReference{},
										ReasonCode:        []*domain.FHIRCodeableConcept{},
										ReasonReference:   []*domain.FHIRReference{},
										Note:              []*domain.FHIRAnnotation{},
										Dosage:            []*domain.FHIRDosage{},
										Meta:              &domain.FHIRMeta{},
										Extension:         []*domain.FHIRExtension{},
									},
								},
							},
							PageInfo: &firebasetools.PageInfo{},
						}, nil
					})

				mh.FHIR.EXPECT().SearchFHIRAllergyIntolerance(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRAllergy, error) {
						code := "123"
						system := gofakeit.URL()

						return &domain.PagedFHIRAllergy{
							Allergies: []domain.FHIRAllergyIntolerance{
								{
									Code: &domain.FHIRCodeableConcept{
										ID: new(string),
										Coding: []*domain.FHIRCoding{
											{
												Code:    (*scalarutils.Code)(&code),
												Display: gofakeit.BS(),
												System:  (*scalarutils.URI)(&system),
											},
										},
									},
									Patient: &domain.FHIRReference{
										ID: new(string),
									},
									Encounter: &domain.FHIRReference{
										ID: new(string),
									},
									OnsetPeriod: &domain.FHIRPeriod{
										Start: "2000-01-01",
									},
								},
							},
							HasNextPage:     false,
							NextCursor:      "",
							HasPreviousPage: false,
							PreviousCursor:  "",
							TotalCount:      0,
						}, nil
					})
				mh.FHIR.EXPECT().SearchFHIRObservation(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRObservations, error) {
						status := dto.ObservationStatusFinal
						valueConcept := "222"
						UUID := gofakeit.UUID()
						currentTime := time.Now().Format(time.RFC3339)
						return &domain.PagedFHIRObservations{
							Observations: []domain.FHIRObservation{
								{
									Status: (*domain.ObservationStatusEnum)(&status),
									Code: &domain.FHIRCodeableConcept{
										ID:     new(string),
										Coding: []*domain.FHIRCoding{},
									},
									EffectiveDateTime: (*scalarutils.DateTime)(&currentTime),
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
						}, nil
					})

				return args{
					ctx:       context.Background(),
					patientID: gofakeit.UUID(),
				}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to search observation - nil status",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchFHIRMedicationStatement(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIRMedicationStatementRelayConnection, error) {
						status := dto.MedicationStatementStatusEnumActive
						code := "123"
						system := gofakeit.URL()

						return &domain.FHIRMedicationStatementRelayConnection{
							Edges: []*domain.FHIRMedicationStatementRelayEdge{
								{
									Cursor: new(string),
									Node: &domain.FHIRMedicationStatement{
										ID:     new(string),
										Status: (*domain.MedicationStatementStatusEnum)(&status),
										MedicationCodeableConcept: &domain.FHIRCodeableConcept{
											ID: new(string),
											Coding: []*domain.FHIRCoding{{
												ID:           new(string),
												System:       (*scalarutils.URI)(&system),
												Version:      new(string),
												Code:         (*scalarutils.Code)(&code),
												Display:      gofakeit.BS(),
												UserSelected: new(bool),
											}},
										},
										Subject: &domain.FHIRReference{
											ID: new(string),
										},
										Text:         &domain.FHIRNarrative{},
										Identifier:   []*domain.FHIRIdentifier{},
										BasedOn:      []*domain.FHIRReference{},
										PartOf:       []*domain.FHIRReference{},
										StatusReason: []*domain.FHIRCodeableConcept{},
										Category: []*domain.FHIRCodeableConcept{
											{
												ID: new(string),
												Coding: []*domain.FHIRCoding{
													{
														Code:    (*scalarutils.Code)(&code),
														Display: gofakeit.BS(),
													},
												},
												Text: "",
											},
										},
										Context:           &domain.FHIRReference{},
										EffectiveDateTime: &scalarutils.Date{},
										EffectivePeriod:   &domain.FHIRPeriod{},
										DateAsserted:      &scalarutils.Date{},
										InformationSource: &domain.FHIRReference{},
										DerivedFrom:       []*domain.FHIRReference{},
										ReasonCode:        []*domain.FHIRCodeableConcept{},
										ReasonReference:   []*domain.FHIRReference{},
										Note:              []*domain.FHIRAnnotation{},
										Dosage:            []*domain.FHIRDosage{},
										Meta:              &domain.FHIRMeta{},
										Extension:         []*domain.FHIRExtension{},
									},
								},
							},
							PageInfo: &firebasetools.PageInfo{},
						}, nil
					})

				mh.FHIR.EXPECT().SearchFHIRAllergyIntolerance(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRAllergy, error) {
						code := "123"
						system := gofakeit.URL()

						return &domain.PagedFHIRAllergy{
							Allergies: []domain.FHIRAllergyIntolerance{
								{
									Code: &domain.FHIRCodeableConcept{
										ID: new(string),
										Coding: []*domain.FHIRCoding{
											{
												Code:    (*scalarutils.Code)(&code),
												Display: gofakeit.BS(),
												System:  (*scalarutils.URI)(&system),
											},
										},
									},
									Patient: &domain.FHIRReference{
										ID: new(string),
									},
									Encounter: &domain.FHIRReference{
										ID: new(string),
									},
									OnsetPeriod: &domain.FHIRPeriod{
										Start: "2000-01-01",
									},
								},
							},
							HasNextPage:     false,
							NextCursor:      "",
							HasPreviousPage: false,
							PreviousCursor:  "",
							TotalCount:      0,
						}, nil
					})
				mh.FHIR.EXPECT().SearchFHIRObservation(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRObservations, error) {
						valueConcept := "222"
						UUID := gofakeit.UUID()
						currentTime := time.Now().Format(time.RFC3339)
						return &domain.PagedFHIRObservations{
							Observations: []domain.FHIRObservation{
								{
									ID: new(string),
									Code: &domain.FHIRCodeableConcept{
										ID: new(string),
										Coding: []*domain.FHIRCoding{{
											Display: gofakeit.BS(),
											Code:    (*scalarutils.Code)(&valueConcept),
										}},
									},
									EffectiveDateTime: (*scalarutils.DateTime)(&currentTime),
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
						}, nil
					})

				return args{
					ctx:       context.Background(),
					patientID: gofakeit.UUID(),
				}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to search observation - nil subject",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchFHIRMedicationStatement(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIRMedicationStatementRelayConnection, error) {
						status := dto.MedicationStatementStatusEnumActive
						code := "123"
						system := gofakeit.URL()

						return &domain.FHIRMedicationStatementRelayConnection{
							Edges: []*domain.FHIRMedicationStatementRelayEdge{
								{
									Cursor: new(string),
									Node: &domain.FHIRMedicationStatement{
										ID:     new(string),
										Status: (*domain.MedicationStatementStatusEnum)(&status),
										MedicationCodeableConcept: &domain.FHIRCodeableConcept{
											ID: new(string),
											Coding: []*domain.FHIRCoding{{
												ID:           new(string),
												System:       (*scalarutils.URI)(&system),
												Version:      new(string),
												Code:         (*scalarutils.Code)(&code),
												Display:      gofakeit.BS(),
												UserSelected: new(bool),
											}},
										},
										Subject: &domain.FHIRReference{
											ID: new(string),
										},
										Text:         &domain.FHIRNarrative{},
										Identifier:   []*domain.FHIRIdentifier{},
										BasedOn:      []*domain.FHIRReference{},
										PartOf:       []*domain.FHIRReference{},
										StatusReason: []*domain.FHIRCodeableConcept{},
										Category: []*domain.FHIRCodeableConcept{
											{
												ID: new(string),
												Coding: []*domain.FHIRCoding{
													{
														Code:    (*scalarutils.Code)(&code),
														Display: gofakeit.BS(),
													},
												},
												Text: "",
											},
										},
										Context:           &domain.FHIRReference{},
										EffectiveDateTime: &scalarutils.Date{},
										EffectivePeriod:   &domain.FHIRPeriod{},
										DateAsserted:      &scalarutils.Date{},
										InformationSource: &domain.FHIRReference{},
										DerivedFrom:       []*domain.FHIRReference{},
										ReasonCode:        []*domain.FHIRCodeableConcept{},
										ReasonReference:   []*domain.FHIRReference{},
										Note:              []*domain.FHIRAnnotation{},
										Dosage:            []*domain.FHIRDosage{},
										Meta:              &domain.FHIRMeta{},
										Extension:         []*domain.FHIRExtension{},
									},
								},
							},
							PageInfo: &firebasetools.PageInfo{},
						}, nil
					})

				mh.FHIR.EXPECT().SearchFHIRAllergyIntolerance(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRAllergy, error) {
						code := "123"
						system := gofakeit.URL()

						return &domain.PagedFHIRAllergy{
							Allergies: []domain.FHIRAllergyIntolerance{
								{
									Code: &domain.FHIRCodeableConcept{
										ID: new(string),
										Coding: []*domain.FHIRCoding{
											{
												Code:    (*scalarutils.Code)(&code),
												Display: gofakeit.BS(),
												System:  (*scalarutils.URI)(&system),
											},
										},
									},
									Patient: &domain.FHIRReference{
										ID: new(string),
									},
									Encounter: &domain.FHIRReference{
										ID: new(string),
									},
									OnsetPeriod: &domain.FHIRPeriod{
										Start: "2000-01-01",
									},
								},
							},
							HasNextPage:     false,
							NextCursor:      "",
							HasPreviousPage: false,
							PreviousCursor:  "",
							TotalCount:      0,
						}, nil
					})
				mh.FHIR.EXPECT().SearchFHIRObservation(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRObservations, error) {
						instant := gofakeit.TimeZone()
						valueConcept := "222"
						UUID := gofakeit.UUID()
						currentTime := time.Now().Format(time.RFC3339)
						return &domain.PagedFHIRObservations{
							Observations: []domain.FHIRObservation{
								{
									ID: new(string),
									Code: &domain.FHIRCodeableConcept{
										ID: new(string),
										Coding: []*domain.FHIRCoding{{
											Display: gofakeit.BS(),
											Code:    (*scalarutils.Code)(&valueConcept),
										}},
									},
									EffectiveDateTime: (*scalarutils.DateTime)(&currentTime),
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
						}, nil
					})

				return args{
					ctx:       context.Background(),
					patientID: gofakeit.UUID(),
				}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to search observation - nil subject id",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchFHIRMedicationStatement(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIRMedicationStatementRelayConnection, error) {
						status := dto.MedicationStatementStatusEnumActive
						code := "123"
						system := gofakeit.URL()

						return &domain.FHIRMedicationStatementRelayConnection{
							Edges: []*domain.FHIRMedicationStatementRelayEdge{
								{
									Cursor: new(string),
									Node: &domain.FHIRMedicationStatement{
										ID:     new(string),
										Status: (*domain.MedicationStatementStatusEnum)(&status),
										MedicationCodeableConcept: &domain.FHIRCodeableConcept{
											ID: new(string),
											Coding: []*domain.FHIRCoding{{
												ID:           new(string),
												System:       (*scalarutils.URI)(&system),
												Version:      new(string),
												Code:         (*scalarutils.Code)(&code),
												Display:      gofakeit.BS(),
												UserSelected: new(bool),
											}},
										},
										Subject: &domain.FHIRReference{
											ID: new(string),
										},
										Text:         &domain.FHIRNarrative{},
										Identifier:   []*domain.FHIRIdentifier{},
										BasedOn:      []*domain.FHIRReference{},
										PartOf:       []*domain.FHIRReference{},
										StatusReason: []*domain.FHIRCodeableConcept{},
										Category: []*domain.FHIRCodeableConcept{
											{
												ID: new(string),
												Coding: []*domain.FHIRCoding{
													{
														Code:    (*scalarutils.Code)(&code),
														Display: gofakeit.BS(),
													},
												},
												Text: "",
											},
										},
										Context:           &domain.FHIRReference{},
										EffectiveDateTime: &scalarutils.Date{},
										EffectivePeriod:   &domain.FHIRPeriod{},
										DateAsserted:      &scalarutils.Date{},
										InformationSource: &domain.FHIRReference{},
										DerivedFrom:       []*domain.FHIRReference{},
										ReasonCode:        []*domain.FHIRCodeableConcept{},
										ReasonReference:   []*domain.FHIRReference{},
										Note:              []*domain.FHIRAnnotation{},
										Dosage:            []*domain.FHIRDosage{},
										Meta:              &domain.FHIRMeta{},
										Extension:         []*domain.FHIRExtension{},
									},
								},
							},
							PageInfo: &firebasetools.PageInfo{},
						}, nil
					})

				mh.FHIR.EXPECT().SearchFHIRAllergyIntolerance(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRAllergy, error) {
						code := "123"
						system := gofakeit.URL()

						return &domain.PagedFHIRAllergy{
							Allergies: []domain.FHIRAllergyIntolerance{
								{
									Code: &domain.FHIRCodeableConcept{
										ID: new(string),
										Coding: []*domain.FHIRCoding{
											{
												Code:    (*scalarutils.Code)(&code),
												Display: gofakeit.BS(),
												System:  (*scalarutils.URI)(&system),
											},
										},
									},
									Patient: &domain.FHIRReference{
										ID: new(string),
									},
									Encounter: &domain.FHIRReference{
										ID: new(string),
									},
									OnsetPeriod: &domain.FHIRPeriod{
										Start: "2000-01-01",
									},
								},
							},
							HasNextPage:     false,
							NextCursor:      "",
							HasPreviousPage: false,
							PreviousCursor:  "",
							TotalCount:      0,
						}, nil
					})
				mh.FHIR.EXPECT().SearchFHIRObservation(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRObservations, error) {
						status := dto.ObservationStatusFinal
						valueConcept := "222"
						UUID := gofakeit.UUID()
						currentTime := time.Now().Format(time.RFC3339)
						return &domain.PagedFHIRObservations{
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
									EffectiveDateTime: (*scalarutils.DateTime)(&currentTime),
									Subject:           &domain.FHIRReference{},
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
						}, nil
					})

				return args{
					ctx:       context.Background(),
					patientID: gofakeit.UUID(),
				}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to search observation - nil encounter",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchFHIRMedicationStatement(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIRMedicationStatementRelayConnection, error) {
						status := dto.MedicationStatementStatusEnumActive
						code := "123"
						system := gofakeit.URL()

						return &domain.FHIRMedicationStatementRelayConnection{
							Edges: []*domain.FHIRMedicationStatementRelayEdge{
								{
									Cursor: new(string),
									Node: &domain.FHIRMedicationStatement{
										ID:     new(string),
										Status: (*domain.MedicationStatementStatusEnum)(&status),
										MedicationCodeableConcept: &domain.FHIRCodeableConcept{
											ID: new(string),
											Coding: []*domain.FHIRCoding{{
												ID:           new(string),
												System:       (*scalarutils.URI)(&system),
												Version:      new(string),
												Code:         (*scalarutils.Code)(&code),
												Display:      gofakeit.BS(),
												UserSelected: new(bool),
											}},
										},
										Subject: &domain.FHIRReference{
											ID: new(string),
										},
										Text:         &domain.FHIRNarrative{},
										Identifier:   []*domain.FHIRIdentifier{},
										BasedOn:      []*domain.FHIRReference{},
										PartOf:       []*domain.FHIRReference{},
										StatusReason: []*domain.FHIRCodeableConcept{},
										Category: []*domain.FHIRCodeableConcept{
											{
												ID: new(string),
												Coding: []*domain.FHIRCoding{
													{
														Code:    (*scalarutils.Code)(&code),
														Display: gofakeit.BS(),
													},
												},
												Text: "",
											},
										},
										Context:           &domain.FHIRReference{},
										EffectiveDateTime: &scalarutils.Date{},
										EffectivePeriod:   &domain.FHIRPeriod{},
										DateAsserted:      &scalarutils.Date{},
										InformationSource: &domain.FHIRReference{},
										DerivedFrom:       []*domain.FHIRReference{},
										ReasonCode:        []*domain.FHIRCodeableConcept{},
										ReasonReference:   []*domain.FHIRReference{},
										Note:              []*domain.FHIRAnnotation{},
										Dosage:            []*domain.FHIRDosage{},
										Meta:              &domain.FHIRMeta{},
										Extension:         []*domain.FHIRExtension{},
									},
								},
							},
							PageInfo: &firebasetools.PageInfo{},
						}, nil
					})

				mh.FHIR.EXPECT().SearchFHIRAllergyIntolerance(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRAllergy, error) {
						code := "123"
						system := gofakeit.URL()

						return &domain.PagedFHIRAllergy{
							Allergies: []domain.FHIRAllergyIntolerance{
								{
									Code: &domain.FHIRCodeableConcept{
										ID: new(string),
										Coding: []*domain.FHIRCoding{
											{
												Code:    (*scalarutils.Code)(&code),
												Display: gofakeit.BS(),
												System:  (*scalarutils.URI)(&system),
											},
										},
									},
									Patient: &domain.FHIRReference{
										ID: new(string),
									},
									Encounter: &domain.FHIRReference{
										ID: new(string),
									},
									OnsetPeriod: &domain.FHIRPeriod{
										Start: "2000-01-01",
									},
								},
							},
							HasNextPage:     false,
							NextCursor:      "",
							HasPreviousPage: false,
							PreviousCursor:  "",
							TotalCount:      0,
						}, nil
					})
				mh.FHIR.EXPECT().SearchFHIRObservation(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRObservations, error) {
						status := dto.ObservationStatusFinal
						valueConcept := "222"
						UUID := gofakeit.UUID()
						currentTime := time.Now().Format(time.RFC3339)
						return &domain.PagedFHIRObservations{
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
									EffectiveDateTime: (*scalarutils.DateTime)(&currentTime),
									Subject:           &domain.FHIRReference{},
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
						}, nil
					})

				return args{
					ctx:       context.Background(),
					patientID: gofakeit.UUID(),
				}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to search observation - nil encounter id",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchFHIRMedicationStatement(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIRMedicationStatementRelayConnection, error) {
						status := dto.MedicationStatementStatusEnumActive
						code := "123"
						system := gofakeit.URL()

						return &domain.FHIRMedicationStatementRelayConnection{
							Edges: []*domain.FHIRMedicationStatementRelayEdge{
								{
									Cursor: new(string),
									Node: &domain.FHIRMedicationStatement{
										ID:     new(string),
										Status: (*domain.MedicationStatementStatusEnum)(&status),
										MedicationCodeableConcept: &domain.FHIRCodeableConcept{
											ID: new(string),
											Coding: []*domain.FHIRCoding{{
												ID:           new(string),
												System:       (*scalarutils.URI)(&system),
												Version:      new(string),
												Code:         (*scalarutils.Code)(&code),
												Display:      gofakeit.BS(),
												UserSelected: new(bool),
											}},
										},
										Subject: &domain.FHIRReference{
											ID: new(string),
										},
										Text:         &domain.FHIRNarrative{},
										Identifier:   []*domain.FHIRIdentifier{},
										BasedOn:      []*domain.FHIRReference{},
										PartOf:       []*domain.FHIRReference{},
										StatusReason: []*domain.FHIRCodeableConcept{},
										Category: []*domain.FHIRCodeableConcept{
											{
												ID: new(string),
												Coding: []*domain.FHIRCoding{
													{
														Code:    (*scalarutils.Code)(&code),
														Display: gofakeit.BS(),
													},
												},
												Text: "",
											},
										},
										Context:           &domain.FHIRReference{},
										EffectiveDateTime: &scalarutils.Date{},
										EffectivePeriod:   &domain.FHIRPeriod{},
										DateAsserted:      &scalarutils.Date{},
										InformationSource: &domain.FHIRReference{},
										DerivedFrom:       []*domain.FHIRReference{},
										ReasonCode:        []*domain.FHIRCodeableConcept{},
										ReasonReference:   []*domain.FHIRReference{},
										Note:              []*domain.FHIRAnnotation{},
										Dosage:            []*domain.FHIRDosage{},
										Meta:              &domain.FHIRMeta{},
										Extension:         []*domain.FHIRExtension{},
									},
								},
							},
							PageInfo: &firebasetools.PageInfo{},
						}, nil
					})

				mh.FHIR.EXPECT().SearchFHIRAllergyIntolerance(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRAllergy, error) {
						code := "123"
						system := gofakeit.URL()

						return &domain.PagedFHIRAllergy{
							Allergies: []domain.FHIRAllergyIntolerance{
								{
									Code: &domain.FHIRCodeableConcept{
										ID: new(string),
										Coding: []*domain.FHIRCoding{
											{
												Code:    (*scalarutils.Code)(&code),
												Display: gofakeit.BS(),
												System:  (*scalarutils.URI)(&system),
											},
										},
									},
									Patient: &domain.FHIRReference{
										ID: new(string),
									},
									Encounter: &domain.FHIRReference{
										ID: new(string),
									},
									OnsetPeriod: &domain.FHIRPeriod{
										Start: "2000-01-01",
									},
								},
							},
							HasNextPage:     false,
							NextCursor:      "",
							HasPreviousPage: false,
							PreviousCursor:  "",
							TotalCount:      0,
						}, nil
					})
				mh.FHIR.EXPECT().SearchFHIRObservation(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRObservations, error) {
						status := dto.ObservationStatusFinal
						valueConcept := "222"
						UUID := gofakeit.UUID()
						currentTime := time.Now().Format(time.RFC3339)
						return &domain.PagedFHIRObservations{
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
									EffectiveDateTime: (*scalarutils.DateTime)(&currentTime),
									Subject:           &domain.FHIRReference{},
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
						}, nil
					})

				return args{
					ctx:       context.Background(),
					patientID: gofakeit.UUID(),
				}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to search weight",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchFHIRMedicationStatement(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIRMedicationStatementRelayConnection, error) {
						return &domain.FHIRMedicationStatementRelayConnection{
							Edges: []*domain.FHIRMedicationStatementRelayEdge{},
						}, nil
					})

				mh.FHIR.EXPECT().SearchFHIRAllergyIntolerance(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRAllergy, error) {
						return &domain.PagedFHIRAllergy{
							Allergies: []domain.FHIRAllergyIntolerance{},
						}, nil
					})
				mh.FHIR.EXPECT().SearchFHIRObservation(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRObservations, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{
					ctx:       context.Background(),
					patientID: gofakeit.UUID(),
				}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to search BMI",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchFHIRMedicationStatement(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIRMedicationStatementRelayConnection, error) {
						return &domain.FHIRMedicationStatementRelayConnection{
							Edges: []*domain.FHIRMedicationStatementRelayEdge{},
						}, nil
					})

				mh.FHIR.EXPECT().SearchFHIRAllergyIntolerance(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRAllergy, error) {
						return &domain.PagedFHIRAllergy{
							Allergies: []domain.FHIRAllergyIntolerance{},
						}, nil
					})
				mh.FHIR.EXPECT().SearchFHIRObservation(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRObservations, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{
					ctx:       context.Background(),
					patientID: gofakeit.UUID(),
				}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to search viralLoad",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchFHIRMedicationStatement(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIRMedicationStatementRelayConnection, error) {
						return &domain.FHIRMedicationStatementRelayConnection{
							Edges: []*domain.FHIRMedicationStatementRelayEdge{},
						}, nil
					})

				mh.FHIR.EXPECT().SearchFHIRAllergyIntolerance(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRAllergy, error) {
						return &domain.PagedFHIRAllergy{
							Allergies: []domain.FHIRAllergyIntolerance{},
						}, nil
					})
				mh.FHIR.EXPECT().SearchFHIRObservation(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRObservations, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{
					ctx:       context.Background(),
					patientID: gofakeit.UUID(),
				}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to search cd4Count",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchFHIRMedicationStatement(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIRMedicationStatementRelayConnection, error) {
						return &domain.FHIRMedicationStatementRelayConnection{
							Edges: []*domain.FHIRMedicationStatementRelayEdge{},
						}, nil
					})

				mh.FHIR.EXPECT().SearchFHIRAllergyIntolerance(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRAllergy, error) {
						return &domain.PagedFHIRAllergy{
							Allergies: []domain.FHIRAllergyIntolerance{},
						}, nil
					})
				mh.FHIR.EXPECT().SearchFHIRObservation(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRObservations, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{
					ctx:       context.Background(),
					patientID: gofakeit.UUID(),
				}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			got, err := clinicalUsecase.GetMedicalData(args.ctx, args.patientID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClinicalUseCaseImpl.GetMedicalData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && got != nil {
				t.Errorf("expected patient medical data to be nil for %v", tt.name)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("expected patient medical data not to be nil for %v", tt.name)
				return
			}
		})
	}

}

func TestUseCasesClinicalImpl_CreatePatient(t *testing.T) {

	type args struct {
		ctx   context.Context
		input dto.PatientInput
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "happy case: register a patient",
			setup: func(mh *usecaseMock.Mocks) args {
				input := dto.PatientInput{
					FirstName: gofakeit.Name(),
					LastName:  gofakeit.Name(),
					BirthDate: &scalarutils.Date{
						Year:  1997,
						Month: 12,
						Day:   12,
					},
					Gender: dto.GenderFemale,
					Identifiers: []dto.IdentifierInput{
						{
							Type:  dto.IdentifierTypeNationalID,
							Value: "12345678",
						},
					},
					Contacts: []dto.ContactInput{
						{
							Type:  dto.ContactTypePhoneNumber,
							Value: "0700000000",
						},
					},
				}

				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						orgID := uuid.NewString()
						name := gofakeit.Name()
						return &domain.FHIROrganizationRelayPayload{
							Resource: &domain.FHIROrganization{
								ID:         &orgID,
								Identifier: []*domain.FHIRIdentifier{},
								Name:       &name,
							},
						}, nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return &dto.TenantIdentifiers{
							OrganizationID: uuid.NewString(),
							FacilityID:     uuid.NewString(),
						}, nil
					})
				mh.FHIR.EXPECT().CreateFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRPatientInput) (*domain.PatientPayload, error) {
						patientID := uuid.NewString()
						system := domain.ContactPointSystemEnumPhone
						phone := gofakeit.Phone()
						active := false
						name := gofakeit.Name()
						gender := domain.PatientGenderEnumMale
						return &domain.PatientPayload{
							PatientRecord: &domain.FHIRPatient{
								ID: &patientID,
								Telecom: []*domain.FHIRContactPoint{
									{
										System: &system,
										Value:  &phone,
									},
								},
								Active: &active,
								Name: []*domain.FHIRHumanName{
									{
										Given: []*string{
											&name,
										},
									},
								},
								Gender: &gender,
							},
						}, nil
					})

				return args{ctx: usecaseMock.AddTenantIdentifierContext(context.Background()), input: input}
			},
			wantErr: false,
		},
		{
			name: "sad case: invalid phone number",
			setup: func(mh *usecaseMock.Mocks) args {
				input := dto.PatientInput{
					FirstName: gofakeit.Name(),
					LastName:  gofakeit.Name(),
					BirthDate: &scalarutils.Date{
						Year:  1997,
						Month: 12,
						Day:   12,
					},
					Gender: dto.GenderFemale,
					Identifiers: []dto.IdentifierInput{
						{
							Type:  dto.IdentifierTypeNationalID,
							Value: "12345678",
						},
					},
					Contacts: []dto.ContactInput{
						{
							Type:  dto.ContactTypePhoneNumber,
							Value: "070000",
						},
					},
				}

				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						orgID := uuid.NewString()
						name := gofakeit.Name()
						return &domain.FHIROrganizationRelayPayload{
							Resource: &domain.FHIROrganization{
								ID:         &orgID,
								Identifier: []*domain.FHIRIdentifier{},
								Name:       &name,
							},
						}, nil
					})

				return args{ctx: usecaseMock.AddTenantIdentifierContext(context.Background()), input: input}
			},
			wantErr: true,
		},
		{
			name: "sad case: fail to get tenant tags",
			setup: func(mh *usecaseMock.Mocks) args {
				input := dto.PatientInput{
					FirstName: gofakeit.Name(),
					LastName:  gofakeit.Name(),
					BirthDate: &scalarutils.Date{
						Year:  1997,
						Month: 12,
						Day:   12,
					},
					Gender: dto.GenderFemale,
					Identifiers: []dto.IdentifierInput{
						{
							Type:  dto.IdentifierTypeNationalID,
							Value: "12345678",
						},
					},
					Contacts: []dto.ContactInput{
						{
							Type:  dto.ContactTypePhoneNumber,
							Value: "0700000000",
						},
					},
				}

				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						orgID := uuid.NewString()
						name := gofakeit.Name()
						return &domain.FHIROrganizationRelayPayload{
							Resource: &domain.FHIROrganization{
								ID:         &orgID,
								Identifier: []*domain.FHIRIdentifier{},
								Name:       &name,
							},
						}, nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: usecaseMock.AddTenantIdentifierContext(context.Background()), input: input}
			},
			wantErr: true,
		},
		{
			name: "sad case: fail to create patient",
			setup: func(mh *usecaseMock.Mocks) args {
				input := dto.PatientInput{
					FirstName: gofakeit.Name(),
					LastName:  gofakeit.Name(),
					BirthDate: &scalarutils.Date{
						Year:  1997,
						Month: 12,
						Day:   12,
					},
					Gender: dto.GenderFemale,
					Identifiers: []dto.IdentifierInput{
						{
							Type:  dto.IdentifierTypeNationalID,
							Value: "12345678",
						},
					},
					Contacts: []dto.ContactInput{
						{
							Type:  dto.ContactTypePhoneNumber,
							Value: "0700000000",
						},
					},
				}

				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						orgID := uuid.NewString()
						name := gofakeit.Name()
						return &domain.FHIROrganizationRelayPayload{
							Resource: &domain.FHIROrganization{
								ID:         &orgID,
								Identifier: []*domain.FHIRIdentifier{},
								Name:       &name,
							},
						}, nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return &dto.TenantIdentifiers{
							OrganizationID: uuid.NewString(),
							FacilityID:     uuid.NewString(),
						}, nil
					})
				mh.FHIR.EXPECT().CreateFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRPatientInput) (*domain.PatientPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: usecaseMock.AddTenantIdentifierContext(context.Background()), input: input}
			},
			wantErr: true,
		},
		{
			name: "sad case: no facility id in context",
			setup: func(mh *usecaseMock.Mocks) args {
				input := dto.PatientInput{
					FirstName: gofakeit.Name(),
					LastName:  gofakeit.Name(),
					BirthDate: &scalarutils.Date{
						Year:  1997,
						Month: 12,
						Day:   12,
					},
					Gender: dto.GenderFemale,
					Identifiers: []dto.IdentifierInput{
						{
							Type:  dto.IdentifierTypeNationalID,
							Value: "12345678",
						},
					},
					Contacts: []dto.ContactInput{
						{
							Type:  dto.ContactTypePhoneNumber,
							Value: "0700000000",
						},
					},
				}

				return args{ctx: context.Background(), input: input}
			},
			wantErr: true,
		},
		{
			name: "sad case: fail to find facility",
			setup: func(mh *usecaseMock.Mocks) args {
				input := dto.PatientInput{
					FirstName: gofakeit.Name(),
					LastName:  gofakeit.Name(),
					BirthDate: &scalarutils.Date{
						Year:  1997,
						Month: 12,
						Day:   12,
					},
					Gender: dto.GenderFemale,
					Identifiers: []dto.IdentifierInput{
						{
							Type:  dto.IdentifierTypeNationalID,
							Value: "12345678",
						},
					},
					Contacts: []dto.ContactInput{
						{
							Type:  dto.ContactTypePhoneNumber,
							Value: "0700000000",
						},
					},
				}

				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: usecaseMock.AddTenantIdentifierContext(context.Background()), input: input}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			got, err := clinicalUsecase.CreatePatient(args.ctx, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreatePatient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && got != nil {
				t.Errorf("expected patient to be nil for %v", tt.name)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("expected patients not to be nil for %v", tt.name)
				return
			}
		})
	}
}

func TestUseCasesClinicalImpl_PatchPatient(t *testing.T) {
	type args struct {
		ctx   context.Context
		id    string
		input dto.PatientInput
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully patch a patient (single field)",
			setup: func(mh *usecaseMock.Mocks) args {
				input := dto.PatientInput{
					BirthDate: &scalarutils.Date{
						Year:  2000,
						Month: 6,
						Day:   14,
					},
				}

				mh.FHIR.EXPECT().PatchFHIRPatient(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string, input domain.FHIRPatientInput) (*domain.FHIRPatient, error) {
						patientID := uuid.NewString()
						active := true
						system := domain.ContactPointSystemEnumPhone
						value := "0712312345"
						gender := domain.PatientGenderEnumMale

						return &domain.FHIRPatient{
							ID:     &patientID,
							Active: &active,
							Telecom: []*domain.FHIRContactPoint{
								{
									System: &system,
									Value:  &value,
								},
							},
							Name: []*domain.FHIRHumanName{
								{
									Text: gofakeit.Name(),
								},
							},
							Gender: &gender,
							BirthDate: &scalarutils.Date{
								Year:  2000,
								Month: 6,
								Day:   14,
							},
						}, nil
					})

				return args{ctx: usecaseMock.AddTenantIdentifierContext(context.Background()), id: uuid.NewString(), input: input}
			},
			wantErr: false,
		},
		{
			name: "Happy Case - Successfully patch a patient (multiple fields)",
			setup: func(mh *usecaseMock.Mocks) args {
				input := dto.PatientInput{
					BirthDate: &scalarutils.Date{
						Year:  2000,
						Month: 6,
						Day:   14,
					},
					FirstName: gofakeit.Name(),
					Identifiers: []dto.IdentifierInput{
						{
							Type:  dto.IdentifierTypeNationalID,
							Value: "12345678",
						},
					},
					Gender: dto.GenderFemale,
				}

				mh.FHIR.EXPECT().PatchFHIRPatient(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string, input domain.FHIRPatientInput) (*domain.FHIRPatient, error) {
						patientID := uuid.NewString()
						active := true
						system := domain.ContactPointSystemEnumPhone
						value := "0712312345"
						gender := domain.PatientGenderEnumMale

						return &domain.FHIRPatient{
							ID:     &patientID,
							Active: &active,
							Telecom: []*domain.FHIRContactPoint{
								{
									System: &system,
									Value:  &value,
								},
							},
							Name: []*domain.FHIRHumanName{
								{
									Text: gofakeit.Name(),
								},
							},
							Gender: &gender,
							BirthDate: &scalarutils.Date{
								Year:  2000,
								Month: 6,
								Day:   14,
							},
						}, nil
					})

				return args{ctx: usecaseMock.AddTenantIdentifierContext(context.Background()), id: uuid.NewString(), input: input}
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Missing patient ID",
			setup: func(mh *usecaseMock.Mocks) args {
				input := dto.PatientInput{
					BirthDate: &scalarutils.Date{
						Year:  2000,
						Month: 6,
						Day:   14,
					},
					FirstName: gofakeit.Name(),
					Identifiers: []dto.IdentifierInput{
						{
							Type:  dto.IdentifierTypeNationalID,
							Value: "12345678",
						},
					},
					Gender: dto.GenderFemale,
				}
				return args{ctx: usecaseMock.AddTenantIdentifierContext(context.Background()), input: input}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Missing facility ID in context",
			setup: func(mh *usecaseMock.Mocks) args {
				input := dto.PatientInput{
					BirthDate: &scalarutils.Date{
						Year:  2000,
						Month: 6,
						Day:   14,
					},
					FirstName: gofakeit.Name(),
					Identifiers: []dto.IdentifierInput{
						{
							Type:  dto.IdentifierTypeNationalID,
							Value: "12345678",
						},
					},
					Gender: dto.GenderFemale,
				}
				return args{ctx: context.Background(), id: gofakeit.UUID(), input: input}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Invalid phone number",
			setup: func(mh *usecaseMock.Mocks) args {
				input := dto.PatientInput{
					BirthDate: &scalarutils.Date{
						Year:  2000,
						Month: 6,
						Day:   14,
					},
					FirstName: gofakeit.Name(),
					Identifiers: []dto.IdentifierInput{
						{
							Type:  dto.IdentifierTypeNationalID,
							Value: "12345678",
						},
					},
					Gender: dto.GenderFemale,
					Contacts: []dto.ContactInput{
						{
							Type:  dto.ContactTypePhoneNumber,
							Value: "070000",
						},
					},
				}
				return args{ctx: usecaseMock.AddTenantIdentifierContext(context.Background()), id: uuid.NewString(), input: input}
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to patch patient",
			setup: func(mh *usecaseMock.Mocks) args {
				input := dto.PatientInput{
					BirthDate: &scalarutils.Date{
						Year:  2000,
						Month: 6,
						Day:   14,
					},
					FirstName: gofakeit.Name(),
					Identifiers: []dto.IdentifierInput{
						{
							Type:  dto.IdentifierTypeNationalID,
							Value: "12345678",
						},
					},
					Gender: dto.GenderFemale,
					Contacts: []dto.ContactInput{
						{
							Type:  dto.ContactTypePhoneNumber,
							Value: "0712345634",
						},
					},
				}

				mh.FHIR.EXPECT().PatchFHIRPatient(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string, input domain.FHIRPatientInput) (*domain.FHIRPatient, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: usecaseMock.AddTenantIdentifierContext(context.Background()), id: uuid.NewString(), input: input}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			got, err := clinicalUsecase.PatchPatient(args.ctx, args.id, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("PatchPatient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && got != nil {
				t.Errorf("Expected patient to be nil for %v", tt.name)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("Expected patient not to be nil for %v", tt.name)
				return
			}
		})
	}
}

func TestUseCasesClinicalImpl_DeletePatient(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully delete patient",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().DeleteFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (bool, error) {
						return true, nil
					})

				return args{ctx: ctx, id: uuid.NewString()}
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad Case - Missing patient ID",
			setup: func(mh *usecaseMock.Mocks) args {
				return args{ctx: ctx}
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to delete patient",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().DeleteFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (bool, error) {
						return false, fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, id: uuid.NewString()}
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

			got, err := clinicalUsecase.DeletePatient(args.ctx, args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.DeletePatient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UseCasesClinicalImpl.DeletePatient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCasesClinicalImpl_GetPatientEverything(t *testing.T) {
	type args struct {
		ctx          context.Context
		patientID    string
		filterParams *dto.PatientEverythingFilterParams
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy case: get patient everything",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						patientID := uuid.NewString()
						return &domain.FHIRPatientRelayPayload{
							Resource: &domain.FHIRPatient{
								ID: &patientID,
							},
						}, nil
					})
				mh.FHIR.EXPECT().GetFHIRPatientEverything(mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string, params map[string]interface{}) (*domain.PagedFHIRResource, error) {
						return &domain.PagedFHIRResource{
							TotalCount: 2,
							Resources:  []map[string]interface{}{},
						}, nil
					})

				return args{
					ctx:       usecaseMock.AddTenantIdentifierContext(context.Background()),
					patientID: gofakeit.UUID(),
					filterParams: &dto.PatientEverythingFilterParams{
						Count:     "10",
						PageToken: "",
						Since:     "",
						Type:      "Observation,Condition",
						End:       "",
						Start:     "",
					},
				}
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to get patient everything",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						patientID := uuid.NewString()
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
					filterParams: &dto.PatientEverythingFilterParams{
						Count:     "10",
						PageToken: "",
						Since:     "",
						Type:      "Observation,Condition",
						End:       "",
						Start:     "",
					},
				}
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get patient profile",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{
					ctx:       usecaseMock.AddTenantIdentifierContext(context.Background()),
					patientID: gofakeit.UUID(),
					filterParams: &dto.PatientEverythingFilterParams{
						Count:     "10",
						PageToken: "",
						Since:     "",
						Type:      "Observation,Condition",
						End:       "",
						Start:     "",
					},
				}
			},
			wantErr: true,
		},
		{
			name: "Sad case: invalid id",
			setup: func(mh *usecaseMock.Mocks) args {
				return args{
					ctx:       usecaseMock.AddTenantIdentifierContext(context.Background()),
					patientID: "10",
					filterParams: &dto.PatientEverythingFilterParams{
						Count:     "10",
						PageToken: "",
						Since:     "",
						Type:      "Observation,Condition",
						End:       "",
						Start:     "",
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

			_, err := clinicalUsecase.GetPatientEverything(args.ctx, args.patientID, args.filterParams)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.GetPatientEverything() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
