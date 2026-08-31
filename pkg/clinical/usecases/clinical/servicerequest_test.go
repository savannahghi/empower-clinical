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
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/common"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/empower-clinical/pkg/clinical/domain"
	usecaseMock "github.com/savannahghi/empower-clinical/pkg/clinical/usecases/mock"
)

func TestUseCasesClinicalImpl_CreateServiceRequest(t *testing.T) {
	ctx := usecaseMock.AddTenantIdentifierContext(context.Background())
	ID := gofakeit.UUID()
	input := &domain.FHIRServiceRequestInput{
		ID: &ID,
	}

	type args struct {
		ctx   context.Context
		input *domain.FHIRServiceRequestInput
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy case: create service request",
			setup: func(mh *usecaseMock.Mocks) args {
				ref := mock.Anything
				noteText := scalarutils.Markdown(gofakeit.BeerName())
				serviceCode := scalarutils.Code("159623")
				authoredOn := scalarutils.DateTime(gofakeit.Date().GoString())
				servicerequest := &domain.FHIRServiceRequestRelayPayload{
					Resource: &domain.FHIRServiceRequest{
						ID: &ID,
						Text: &domain.FHIRNarrative{
							Div: scalarutils.XHTML(gofakeit.BeerName()),
						},
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
					},
				}
				mh.FHIR.EXPECT().CreateFHIRServiceRequest(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRServiceRequestInput) (*domain.FHIRServiceRequestRelayPayload, error) {
						return servicerequest, nil
					})
				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to create service request",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().CreateFHIRServiceRequest(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRServiceRequestInput) (*domain.FHIRServiceRequestRelayPayload, error) {
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

			_, err := clinicalUsecase.CreateServiceRequest(args.ctx, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.CreateServiceRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCasesClinicalImpl_CreateLaboratoryOrder(t *testing.T) {
	ctx := usecaseMock.AddTenantIdentifierContext(context.Background())
	serviceRequestID := uuid.NewString()

	ref := mock.Anything
	ID := uuid.NewString()
	noteText := scalarutils.Markdown(gofakeit.BeerName())
	serviceCode := scalarutils.Code("159623")
	authoredOn := scalarutils.DateTime(gofakeit.Date().GoString())
	system := scalarutils.URI("http://mycarehub/tenant-identification/facility")
	code := scalarutils.Code(gofakeit.UUID())
	idCode := scalarutils.Code("slade-health-id")
	status := domain.NarrativeStatusEnumAdditional
	servicerequest := &domain.FHIRServiceRequestRelayPayload{
		Resource: &domain.FHIRServiceRequest{
			ID: &ID,
			Text: &domain.FHIRNarrative{
				Status: &status,
				Div:    scalarutils.XHTML(gofakeit.BeerName()),
			},
			Identifier: []*domain.FHIRIdentifier{
				{
					Type: domain.FHIRCodeableConcept{
						Coding: []*domain.FHIRCoding{
							{
								Code: &idCode,
							},
						},
					},
					Value: gofakeit.BeerName(),
				},
			},
			Status: domain.ServiceRequestStatusActive,
			Intent: domain.ServiceRequestIntentDirective,
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

	idSystem := scalarutils.URI("HEALTH_ID")
	patient := domain.FHIRPatientRelayPayload{
		Resource: &domain.FHIRPatient{
			ID: &ID,
			Identifier: []*domain.FHIRIdentifier{
				{
					Type: domain.FHIRCodeableConcept{
						Coding: []*domain.FHIRCoding{
							{
								System:  &idSystem,
								Display: gofakeit.UUID(),
								Code:    &idCode,
							},
						},
					},
				},
			},
		},
	}

	type args struct {
		ctx              context.Context
		serviceRequestID string
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy case: create service request",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRServiceRequest(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRServiceRequestRelayPayload, error) {
						return servicerequest, nil
					})
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return &patient, nil
					})
				mh.FHIR.EXPECT().CreateFHIRServiceRequest(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRServiceRequestInput) (*domain.FHIRServiceRequestRelayPayload, error) {
						return servicerequest, nil
					})

				return args{ctx: ctx, serviceRequestID: serviceRequestID}
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable get service request",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRServiceRequest(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRServiceRequestRelayPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, serviceRequestID: serviceRequestID}
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to create service request",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRServiceRequest(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRServiceRequestRelayPayload, error) {
						return servicerequest, nil
					})
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return &patient, nil
					})
				mh.FHIR.EXPECT().CreateFHIRServiceRequest(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRServiceRequestInput) (*domain.FHIRServiceRequestRelayPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, serviceRequestID: serviceRequestID}
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get facility ID",
			setup: func(mh *usecaseMock.Mocks) args {
				ref := mock.Anything
				ID := uuid.NewString()
				noteText := scalarutils.Markdown(gofakeit.BeerName())
				serviceCode := scalarutils.Code("159623")
				authoredOn := scalarutils.DateTime(gofakeit.Date().GoString())
				system := scalarutils.URI("http://mycarehub/tenant-identification/nofacility")
				code := scalarutils.Code(gofakeit.UUID())
				status := domain.NarrativeStatusEnumAdditional
				servicerequest := &domain.FHIRServiceRequestRelayPayload{
					Resource: &domain.FHIRServiceRequest{
						ID: &ID,
						Text: &domain.FHIRNarrative{
							Status: &status,
							Div:    scalarutils.XHTML(gofakeit.BeerName()),
						},
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
				mh.FHIR.EXPECT().GetFHIRServiceRequest(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRServiceRequestRelayPayload, error) {
						return servicerequest, nil
					})
				return args{ctx: ctx, serviceRequestID: serviceRequestID}
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get patient",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRServiceRequest(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRServiceRequestRelayPayload, error) {
						return servicerequest, nil
					})
				mh.FHIR.EXPECT().GetFHIRPatient(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, serviceRequestID: serviceRequestID}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			_, err := clinicalUsecase.CreateLaboratoryOrder(args.ctx, args.serviceRequestID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.CreateLaboratoryOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCasesClinicalImpl_ListLabOrders(t *testing.T) {
	ctx := usecaseMock.AddTenantIdentifierContext(context.Background())

	filter := &dto.ServiceRequestFilterInput{
		FilterInput: dto.FilterInput{},
		Type:        common.LaboratoryOrderLOINCCode,
		Status:      dto.ServiceRequestStatusActive,
	}

	count := 10
	pagination := dto.Pagination{
		First: &count,
	}

	ID := uuid.NewString()
	tenantIDs := &dto.TenantIdentifiers{
		OrganizationID: ID,
		FacilityID:     ID,
	}

	ref := mock.Anything
	noteText := scalarutils.Markdown(gofakeit.BeerName())
	serviceCode := scalarutils.Code("159623")
	authoredOn := scalarutils.DateTime(gofakeit.Date().GoString())
	system := scalarutils.URI("http://mycarehub/tenant-identification/facility")
	code := scalarutils.Code(gofakeit.UUID())
	labCategoryCode := scalarutils.Code(domain.LaboratoryProcedureCategoryType)
	status := domain.NarrativeStatusEnumAdditional
	pagedfhirservicerequest := &domain.PagedFHIRServiceRequest{
		ServiceRequests: []domain.FHIRServiceRequest{
			{
				ID: &ID,
				Category: []*domain.FHIRCodeableConcept{
					{
						Coding: []*domain.FHIRCoding{
							{Code: &labCategoryCode},
						},
						Text: domain.LaboratoryProcedureCategoryType.Display(),
					},
				},
				Text: &domain.FHIRNarrative{
					Status: &status,
					Div:    scalarutils.XHTML(gofakeit.BeerName()),
				},
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
		},
	}

	type args struct {
		ctx        context.Context
		filter     *dto.ServiceRequestFilterInput
		pagination dto.Pagination
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy case: list lab orders",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchFHIRServiceRequest(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(
						func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRServiceRequest, error) {
							if params["category"] != "laboratory-procedure" {
								return nil, fmt.Errorf("expected laboratory-procedure category, got %v", params["category"])
							}

							if params["code"] != common.LaboratoryOrderLOINCCode {
								return nil, fmt.Errorf("expected code filter %q, got %v", common.LaboratoryOrderLOINCCode, params["code"])
							}

							return pagedfhirservicerequest, nil
						})

				return args{ctx: ctx, filter: filter, pagination: pagination}
			},
			wantErr: false,
		},
		{
			name: "Happy case: general patient search resolves to a subject IN-list",
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
				mh.FHIR.EXPECT().SearchFHIRServiceRequest(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRServiceRequest, error) {
						if params["subject"] != "Patient/patient-1,Patient/patient-2" {
							return nil, fmt.Errorf("expected subject IN-list, got %v", params["subject"])
						}

						return pagedfhirservicerequest, nil
					})

				searchFilter := &dto.ServiceRequestFilterInput{
					Type:          common.LaboratoryOrderLOINCCode,
					Status:        dto.ServiceRequestStatusActive,
					PatientSearch: "jane",
				}

				return args{ctx: ctx, filter: searchFilter, pagination: pagination}
			},
			wantErr: false,
		},
		{
			name: "Happy case: general patient search with no matches short-circuits to empty",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				// No patients match, so the service request search must never run.
				mh.FHIR.EXPECT().SearchFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID, resourceType string, params map[string]any, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRResource, error) {
						return &domain.PagedFHIRResource{Resources: []map[string]interface{}{}}, nil
					})

				searchFilter := &dto.ServiceRequestFilterInput{
					Type:          common.LaboratoryOrderLOINCCode,
					Status:        dto.ServiceRequestStatusActive,
					PatientSearch: "nobody",
				}

				return args{ctx: ctx, filter: searchFilter, pagination: pagination}
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to list lab orders",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchFHIRServiceRequest(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRServiceRequest, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, filter: filter, pagination: pagination}
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

				return args{ctx: ctx, filter: filter, pagination: pagination}
			},
			wantErr: true,
		},
		{
			name: "Sad case: invalid filters",
			setup: func(mh *usecaseMock.Mocks) args {
				filter := &dto.ServiceRequestFilterInput{
					FilterInput: dto.FilterInput{
						PatientID: "123",
					},
					Type:   common.LaboratoryOrderLOINCCode,
					Status: dto.ServiceRequestStatusActive,
				}

				return args{ctx: ctx, filter: filter, pagination: pagination}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			_, err := clinicalUsecase.ListLabOrders(args.ctx, args.filter, args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.ListLabOrders() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCasesClinicalImpl_GetLabOrder(t *testing.T) {
	ctx := context.Background()
	servicerequestID := uuid.NewString()

	ID := uuid.NewString()
	tenantIDs := &dto.TenantIdentifiers{
		OrganizationID: ID,
		FacilityID:     ID,
	}

	ref := mock.Anything
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

	observationstatus := dto.ObservationStatusFinal
	obsTime := scalarutils.Instant(time.Now().Format(time.RFC3339))
	valueString := gofakeit.BeerName()
	pagedfhirobservation := &domain.PagedFHIRObservations{
		Observations: []domain.FHIRObservation{
			{
				ID:     &ID,
				Status: (*domain.ObservationStatusEnum)(&observationstatus),
				Code: &domain.FHIRCodeableConcept{
					Coding: []*domain.FHIRCoding{{
						Display: gofakeit.BS(),
					}},
					Text: gofakeit.BS(),
				},
				EffectiveInstant: &obsTime,
				ValueString:      &valueString,
				Subject: &domain.FHIRReference{
					ID:      &ID,
					Display: ID,
				},
				Encounter: &domain.FHIRReference{
					ID:      &ID,
					Display: ID,
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
		ctx              context.Context
		serviceRequestID string
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		want    *dto.ServiceRequest
		wantErr bool
	}{
		{
			name: "Happy Case: Successfully get a lab order",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().GetFHIRServiceRequest(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRServiceRequestRelayPayload, error) {
						return servicerequest, nil
					})
				mh.FHIR.EXPECT().SearchFHIRObservation(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRObservations, error) {
						return pagedfhirobservation, nil
					})

				return args{ctx: ctx, serviceRequestID: servicerequestID}
			},
			wantErr: false,
		},
		{
			name: "Sad Case: Missing service request ID",
			setup: func(mh *usecaseMock.Mocks) args {
				return args{ctx: ctx, serviceRequestID: ""}
			},
			wantErr: true,
		},
		{
			name: "Sad Case: Fail to get service request (Lab order)",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIRServiceRequest(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRServiceRequestRelayPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, serviceRequestID: servicerequestID}
			},
			wantErr: true,
		},
		{
			name: "Sad Case: unable to search for service request",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().GetFHIRServiceRequest(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRServiceRequestRelayPayload, error) {
						return servicerequest, nil
					})
				mh.FHIR.EXPECT().SearchFHIRObservation(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRObservations, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, serviceRequestID: servicerequestID}
			},
			wantErr: true,
		},
		{
			name: "Sad Case: unable to get tenant identifier",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return nil, fmt.Errorf("an error occurred")
					})
				mh.FHIR.EXPECT().GetFHIRServiceRequest(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRServiceRequestRelayPayload, error) {
						return servicerequest, nil
					})
				return args{ctx: ctx, serviceRequestID: servicerequestID}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			_, err := clinicalUsecase.GetLabOrder(args.ctx, args.serviceRequestID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.GetLabOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCasesClinicalImpl_CreateTestOrder(t *testing.T) {
	ctx := context.Background()
	ID := gofakeit.UUID()

	tenantIDs := &dto.TenantIdentifiers{
		OrganizationID: ID,
		FacilityID:     ID,
	}

	input := &dto.TestOrder{
		EncounterID: gofakeit.UUID(),
		Name:        gofakeit.Name(),
		LoincCode:   gofakeit.Digit(),
		Status:      gofakeit.Name(),
		Facility: dto.ReferralFacility{
			FHIROrganisationID: gofakeit.UUID(),
			Name:               gofakeit.Name(),
			PhoneNumber:        gofakeit.Phone(),
			Email:              gofakeit.Email(),
			Branch:             "",
		},
		Diagnosis:    "Allergy",
		ClinicalNote: gofakeit.Sentence(5),
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

	orgName := mock.Anything
	organization := &domain.FHIROrganizationRelayPayload{
		Resource: &domain.FHIROrganization{
			ID:   &ID,
			Name: &orgName,
		},
	}

	ref := mock.Anything
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
							System:  &system,
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

	type args struct {
		ctx   context.Context
		input *dto.TestOrder
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		want    *dto.ServiceRequest
		wantErr bool
	}{
		{
			name: "Happy Case: Successfully create a test order",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return organization, nil
					})
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
						return conceptpayload, nil
					})
				mh.FHIR.EXPECT().GetFHIRCondition(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRConditionRelayPayload, error) {
						return &domain.FHIRConditionRelayPayload{}, nil
					})

				mh.FHIR.EXPECT().CreateFHIRServiceRequest(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRServiceRequestInput) (*domain.FHIRServiceRequestRelayPayload, error) {
						return servicerequest, nil
					})

				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "Sad Case: fail to create a test order",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return organization, nil
					})
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
						return conceptpayload, nil
					})
				mh.FHIR.EXPECT().GetFHIRCondition(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRConditionRelayPayload, error) {
						return &domain.FHIRConditionRelayPayload{}, nil
					})
				mh.FHIR.EXPECT().CreateFHIRServiceRequest(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRServiceRequestInput) (*domain.FHIRServiceRequestRelayPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
		{
			name: "Sad Case: fail to get encounter",
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
			name: "Sad Case: Record referral in finished encounter",
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
			name: "Sad Case: Unable to get condition",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
						return conceptpayload, nil
					})
				mh.FHIR.EXPECT().GetFHIRCondition(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRConditionRelayPayload, error) {
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

			got, err := clinicalUsecase.CreateTestOrder(args.ctx, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.CreateTestOrder() error = %v, wantErr %v", err, tt.wantErr)
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

// coding builds a *domain.FHIRCoding from a code and display for use in table tests.
func coding(code, display string) *domain.FHIRCoding {
	c := scalarutils.Code(code)

	return &domain.FHIRCoding{Code: &c, Display: display}
}

// serviceRequestWithCodings wraps the supplied codings in a service request's Code block.
func serviceRequestWithCodings(codings ...*domain.FHIRCoding) domain.FHIRServiceRequest {
	return domain.FHIRServiceRequest{
		Code: &domain.FHIRCodeableReference{
			Concept: &domain.FHIRCodeableConcept{
				Coding: codings,
			},
		},
	}
}
