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
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/empower-clinical/pkg/clinical/domain"
	usecaseMock "github.com/savannahghi/empower-clinical/pkg/clinical/usecases/mock"
)

func resourceBnundle() []map[string]any {
	answer := []map[string]interface{}{
		{
			"resourceType": "ServiceRequest",
			"id":           "0631c31d-ced5-45dd-9b4d-d29b24d84457",
			"status":       "active",
			"intent":       "order",
			"category": []any{
				map[string]any{
					"coding": []any{
						map[string]any{
							"system":  "https://fhir.test.slade360edi.com/fhir/CodeSystem/service-request-cs",
							"code":    "referral",
							"display": "Referral",
						},
					},
					"text": "Referral",
				},
			},
			"priority": "urgent",
			"code": map[string]any{
				"concept": map[string]any{
					"coding": []any{
						map[string]any{
							"system":  "http://loinc.org",
							"code":    "42349-1",
							"display": "Reason for referral (narrative)",
						},
					},
					"text": "Referred to Dr Vimal Gavin specialist for further evaluation",
				},
			},
			"subject": map[string]any{
				"id":        "2194a2e9-c366-488d-a228-db9c30240254",
				"reference": "Patient/2194a2e9-c366-488d-a228-db9c30240254",
				"display":   "Njugush, Njeri",
			},
			"encounter": map[string]any{
				"reference": "Encounter/d7ffd2f6-5735-41a8-bf24-75cdcfcfad56",
				"display":   "d7ffd2f6-5735-41a8-bf24-75cdcfcfad56",
			},
			"authoredOn": "2025-11-05T07:46:20+07:00",
			"requester": map[string]any{
				"reference": "Organization/85e4b0d3-1d69-47ba-b265-579d125f18e5",
				"display":   "85e4b0d3-1d69-47ba-b265-579d125f18e5",
			},
			"performer": []any{
				map[string]any{
					"reference": "Organization/cb9c0147-91ef-428c-a13a-144cabe45b3f",
					"display":   "cb9c0147-91ef-428c-a13a-144cabe45b3f",
				},
			},
			"reason": []any{
				map[string]any{
					"concept": map[string]any{
						"coding": []any{
							map[string]any{
								"system":  "http://loinc.org",
								"code":    "57133-1",
								"display": "Referral note",
							},
						},
						"text": "Referral note",
					},
				},
			},
			"note": []any{
				map[string]any{
					"time": "2025-11-05T07:46:20+07:00",
					"text": "Specialized evaluation",
				},
			},
		},
		{

			"resourceType": "Organization",
			"id":           "cb9c0147-91ef-428c-a13a-144cabe45b3f",
			"active":       true,
			"name":         "Valley View",
			"contact": []map[string]any{
				{
					"purpose": map[string]any{
						"coding": []map[string]any{
							{
								"system":  "https://terminology.hl7.org/5.1.0/CodeSystem-contactentity-type.html",
								"code":    "ADMIN",
								"display": "Administrative",
							},
						},
					},
					"name": []map[string]any{
						{
							"text": "Valley View",
						},
					},
					"telecom": []map[string]any{
						{
							"system": "phone",
							"value":  "+25490360360",
							"use":    "work",
							"rank":   1,
						},
						{
							"system": "email",
							"value":  "nrb@hosi.com",
							"use":    "work",
							"rank":   1,
						},
					},
				},
			},
		},
	}

	return answer
}
func TestUseCasesClinicalImpl_GetPatientReferrals(t *testing.T) {
	ctx := context.Background()
	firstTen := 10
	patientID := uuid.New().String()
	encounterID := uuid.New().String()
	pagination := &dto.Pagination{
		First: &firstTen,
	}

	tenantIDs := &dto.TenantIdentifiers{
		OrganizationID: uuid.NewString(),
		FacilityID:     uuid.NewString(),
	}

	pagedFHIRResource := &domain.PagedFHIRResource{
		Resources: []map[string]interface{}{
			{
				"authoredOn": "2024-06-10T08:09:03+08:00",
				"category": []map[string]interface{}{
					{
						"coding": []map[string]interface{}{
							{
								"code":    "167731",
								"display": "Referral",
							},
						},
						"text": "Referral",
					},
				},
				"code": map[string]interface{}{
					"concept": map[string]interface{}{
						"coding": []map[string]interface{}{
							{
								"code":    "159623",
								"display": "Diagnostics",
							},
							{
								"code":    "TEST",
								"display": "Mammogram",
							},
						},
						"text": "Facility Referral Reason",
					},
				},
				"encounter": map[string]interface{}{
					"id":        "65d22133-2f17-4599-99ab-9f3ee0ef020a",
					"reference": "Encounter/65d22133-2f17-4599-99ab-9f3ee0ef020a",
				},
				"extension": []map[string]interface{}{
					{
						"extension": []map[string]interface{}{
							{
								"url":         "facilityName",
								"valueString": "One pad more",
							},
							{
								"url":         "facilityContact",
								"valueString": "+254727645367",
							},
							{
								"url":         "facilityCounty",
								"valueString": "NAIROBI",
							},
						},
						"url": "http://savannahghi.org/fhir/StructureDefinition/referred-facility",
					},
				},
				"id":       "551ea02e-22cd-4916-bb31-6ffccce9b203",
				"intent":   "order",
				"language": "EN",
				"meta": map[string]interface{}{
					"lastUpdated": "2024-06-10T08:09:03.793059+00:00",
					"tag": []map[string]interface{}{
						{
							"code":         "e303711e-d24a-4fda-87d4-a4c177c7e90d",
							"display":      "Roche",
							"system":       "http://mycarehub/tenant-identification/organisation",
							"userSelected": false,
							"version":      "1.0",
						},
						{
							"code":         "ad031755-9149-44a8-a702-367bafa2ed40",
							"display":      "Main Branch",
							"system":       "http://mycarehub/tenant-identification/facility",
							"userSelected": false,
							"version":      "1.0",
						},
					},
					"versionId": "MTcxODAwNjk0Mzc5MzA1OTAwMA",
				},
				"note": []map[string]interface{}{
					{
						"text": "Testing",
						"time": "2024-06-10T08:09:03+08:00",
					},
				},
				"priority":     "urgent",
				"resourceType": "ServiceRequest",
				"status":       "active",
				"subject": map[string]interface{}{
					"display":   "Talisha, Idah ",
					"id":        "4d808b60-5149-443b-9d03-a6016e5af1b5",
					"reference": "Patient/4d808b60-5149-443b-9d03-a6016e5af1b5",
				},
			},
			{
				"content": []map[string]interface{}{
					{
						"attachment": map[string]interface{}{
							"contentType": "application/pdf",
							"title":       "Doe, Jane 's Referral report",
							"url":         "https://example.invalid/fixtures/report.pdf",
						},
					},
				},
				"context": []map[string]interface{}{
					{
						"reference": "ServiceRequest/551ea02e-22cd-4916-bb31-6ffccce9b203",
					},
				},
				"date":      "2024-06-11T07:57:44Z",
				"docStatus": "final",
				"id":        "9100df95-fbfd-457d-808e-6844cd16ffc2",
				"language":  "EN",
				"meta": map[string]interface{}{
					"lastUpdated": "2024-06-11T07:57:45.480128+00:00",
					"tag": []map[string]interface{}{
						{
							"code":         "e303711e-d24a-4fda-87d4-a4c177c7e90d",
							"display":      "Roche",
							"system":       "http://mycarehub/tenant-identification/organisation",
							"userSelected": false,
							"version":      "1.0",
						},
						{
							"code":         "ad031755-9149-44a8-a702-367bafa2ed40",
							"display":      "Main Branch",
							"system":       "http://mycarehub/tenant-identification/facility",
							"userSelected": false,
							"version":      "1.0",
						},
					},
					"versionId": "MTcxODA5MjY2NTQ4MDEyODAwMA",
				},
				"resourceType": "DocumentReference",
				"status":       "current",
				"subject": map[string]interface{}{
					"id":        "e8131505-b8da-4ad9-b991-f5fc43e5ab48",
					"reference": "Patient/e8131505-b8da-4ad9-b991-f5fc43e5ab48",
				},
				"type": map[string]interface{}{
					"coding": []map[string]interface{}{
						{
							"system": "",
						},
					},
				},
			},
		},
		HasNextPage:     false,
		NextCursor:      "",
		HasPreviousPage: false,
		PreviousCursor:  "",
		TotalCount:      0,
	}
	orgName := "KNH"
	type args struct {
		ctx         context.Context
		searchInput *dto.ReferralSearchInput
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully get patient referrals",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID, resourceType string, params map[string]any, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRResource, error) {
						return pagedFHIRResource, nil
					})
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).RunAndReturn(
					func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return &domain.FHIROrganizationRelayPayload{
							Resource: &domain.FHIROrganization{
								Name: &orgName,
							},
						}, nil
					},
				)
				return args{
					ctx: ctx,
					searchInput: &dto.ReferralSearchInput{
						PatientID:   &patientID,
						EncounterID: &encounterID,
						Pagination:  pagination,
						Status:      domain.ServiceRequestStatusActive,
					},
				}
			},
			wantErr: false,
		},
		{
			name: "Happy Case - general patient search resolves to a subject IN-list",
			setup: func(mh *usecaseMock.Mocks) args {
				p1 := "patient-1"
				p2 := "patient-2"

				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID, resourceType string, params map[string]any, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRResource, error) {
						if resourceType == "Patient" {
							return &domain.PagedFHIRResource{Resources: []map[string]interface{}{
								{"id": p1},
								{"id": p2},
							}}, nil
						}

						if params["subject"] != "Patient/patient-1,Patient/patient-2" {
							return nil, fmt.Errorf("expected subject IN-list, got %v", params["subject"])
						}

						return &domain.PagedFHIRResource{Resources: []map[string]interface{}{}}, nil
					})

				return args{
					ctx: ctx,
					searchInput: &dto.ReferralSearchInput{
						Pagination:    pagination,
						Status:        domain.ServiceRequestStatusActive,
						PatientSearch: "jane",
					},
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
				// No patients match, so the service request search must never run.
				mh.FHIR.EXPECT().SearchFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID, resourceType string, params map[string]any, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRResource, error) {
						if resourceType != "Patient" {
							return nil, fmt.Errorf("service request search should not run when no patients match")
						}

						return &domain.PagedFHIRResource{Resources: []map[string]interface{}{}}, nil
					})

				return args{
					ctx: ctx,
					searchInput: &dto.ReferralSearchInput{
						Pagination:    pagination,
						Status:        domain.ServiceRequestStatusActive,
						PatientSearch: "nobody",
					},
				}
			},
			wantErr: false,
		},
		{
			name: "Sad Case: Fail to get tenant identifiers",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return nil, fmt.Errorf("an error occurred")
					})
				return args{
					ctx: ctx,
					searchInput: &dto.ReferralSearchInput{
						PatientID:   &patientID,
						EncounterID: &encounterID,
						Pagination:  pagination,
						Status:      domain.ServiceRequestStatusActive,
					},
				}
			},
			wantErr: true,
		},
		{
			name: "Sad Case: Fail to get resource",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID, resourceType string, params map[string]any, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRResource, error) {
						return nil, fmt.Errorf("an error occurred")
					})
				return args{
					ctx: ctx,
					searchInput: &dto.ReferralSearchInput{
						PatientID:   &patientID,
						EncounterID: &encounterID,
						Pagination:  pagination,
						Status:      domain.ServiceRequestStatusActive,
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

			_, err := clinicalUsecase.GetPatientReferrals(args.ctx, args.searchInput)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.GetPatientReferrals() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCasesClinicalImpl_ReferPatient(t *testing.T) {
	ctx := context.Background()

	tenantIDs := &dto.TenantIdentifiers{
		OrganizationID: uuid.NewString(),
		FacilityID:     uuid.NewString(),
	}

	ID := uuid.NewString()
	encounter := &domain.FHIREncounterRelayPayload{
		Resource: &domain.FHIREncounter{
			ID:     &ID,
			Status: domain.EncounterStatusEnumInProgress,
			Subject: &domain.FHIRReference{
				ID:      &ID,
				Display: ID,
			},
		},
	}

	orgName := mock.Anything
	organization := &domain.FHIROrganizationRelayPayload{
		Resource: &domain.FHIROrganization{
			ID:   &ID,
			Name: &orgName,
		},
	}

	referralInput := &dto.ReferralInput{
		EncounterID: gofakeit.UUID(),
		Tests:       []string{"VIA"},
		Specialist:  "Oncologist",
		Facility: &dto.FacilityInput{
			Name:               "KNH",
			County:             "Nairobi",
			Contact:            "+254710100100",
			FHIROrganisationID: gofakeit.UUID(),
		},
		ReferralNote: "test",
	}

	ref := mock.Anything
	noteText := scalarutils.Markdown(gofakeit.BeerName())
	serviceCode := scalarutils.Code("159623")
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
				Display:   mock.Anything,
			},
			Encounter: &domain.FHIRReference{
				ID:        &ID,
				Display:   mock.Anything,
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
		},
	}

	type args struct {
		ctx   context.Context
		input *dto.ReferralInput
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy Case: Successfully refer patient for treatment",
			setup: func(mh *usecaseMock.Mocks) args {
				referralInput.ReferralType = dto.TREATMENT

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
				mh.FHIR.EXPECT().CreateFHIRServiceRequest(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRServiceRequestInput) (*domain.FHIRServiceRequestRelayPayload, error) {
						return servicerequest, nil
					})
				mh.PubSub.EXPECT().NotifyCreatePatientReferralTask(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, data dto.PatientReferralTaskPayload) error {
						return nil
					})

				return args{ctx: ctx, input: referralInput}
			},
			wantErr: false,
		},
		{
			name: "Happy Case: Successfully refer patient for specialist",
			setup: func(mh *usecaseMock.Mocks) args {
				referralInput.ReferralType = dto.SPECIALIST

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
				mh.FHIR.EXPECT().CreateFHIRServiceRequest(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRServiceRequestInput) (*domain.FHIRServiceRequestRelayPayload, error) {
						return servicerequest, nil
					})
				mh.PubSub.EXPECT().NotifyCreatePatientReferralTask(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, data dto.PatientReferralTaskPayload) error {
						return nil
					})

				return args{ctx: ctx, input: referralInput}
			},
			wantErr: false,
		},
		{
			name: "Happy Case: Successfully refer patient for diagnostics",
			setup: func(mh *usecaseMock.Mocks) args {
				referralInput.ReferralType = dto.DIAGNOSTICS

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
				mh.FHIR.EXPECT().CreateFHIRServiceRequest(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRServiceRequestInput) (*domain.FHIRServiceRequestRelayPayload, error) {
						return servicerequest, nil
					})
				mh.PubSub.EXPECT().NotifyCreatePatientReferralTask(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, data dto.PatientReferralTaskPayload) error {
						return nil
					})

				return args{ctx: ctx, input: referralInput}
			},
			wantErr: false,
		},
		{
			name: "Sad Case: unable to create diagnostics referral task",
			setup: func(mh *usecaseMock.Mocks) args {
				referralInput.ReferralType = dto.DIAGNOSTICS

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
				mh.FHIR.EXPECT().CreateFHIRServiceRequest(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRServiceRequestInput) (*domain.FHIRServiceRequestRelayPayload, error) {
						return servicerequest, nil
					})
				mh.PubSub.EXPECT().NotifyCreatePatientReferralTask(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, data dto.PatientReferralTaskPayload) error {
						return fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, input: referralInput}
			},
			wantErr: true,
		},
		{
			name: "Sad Case: referral receiving facility unspecified",
			setup: func(mh *usecaseMock.Mocks) args {
				referralInput := &dto.ReferralInput{
					EncounterID:  gofakeit.UUID(),
					ReferralType: dto.TREATMENT,
					Tests:        []string{"VIA"},
					Specialist:   "Oncologist",
					Facility: &dto.FacilityInput{
						Name:    "KNH",
						County:  "Nairobi",
						Contact: "+254710100100",
					},
					ReferralNote: "",
				}

				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})

				return args{ctx: ctx, input: referralInput}
			},
			wantErr: true,
		},
		{
			name: "Sad Case: unable to create treatment referral task",
			setup: func(mh *usecaseMock.Mocks) args {
				referralInput.ReferralType = dto.TREATMENT

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
				mh.FHIR.EXPECT().CreateFHIRServiceRequest(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRServiceRequestInput) (*domain.FHIRServiceRequestRelayPayload, error) {
						return servicerequest, nil
					})
				mh.PubSub.EXPECT().NotifyCreatePatientReferralTask(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, data dto.PatientReferralTaskPayload) error {
						return fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, input: referralInput}
			},
			wantErr: true,
		},
		{
			name: "Sad Case: unable to create specialist referral task",
			setup: func(mh *usecaseMock.Mocks) args {
				referralInput.ReferralType = dto.SPECIALIST

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
				mh.FHIR.EXPECT().CreateFHIRServiceRequest(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRServiceRequestInput) (*domain.FHIRServiceRequestRelayPayload, error) {
						return servicerequest, nil
					})
				mh.PubSub.EXPECT().NotifyCreatePatientReferralTask(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, data dto.PatientReferralTaskPayload) error {
						return fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, input: referralInput}
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

				return args{ctx: ctx, input: referralInput}
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

				return args{ctx: ctx, input: referralInput}
			},
			wantErr: true,
		},
		{
			name: "Sad Case: Fail to get tenant meta tags",
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

				return args{ctx: ctx, input: referralInput}
			},
			wantErr: true,
		},
		{
			name: "Sad Case: Fail to create service request",
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
				mh.FHIR.EXPECT().CreateFHIRServiceRequest(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, input domain.FHIRServiceRequestInput) (*domain.FHIRServiceRequestRelayPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, input: referralInput}
			},
			wantErr: true,
		},
		{
			name: "Sad Case: Input validation - missing encounter ID",
			setup: func(mh *usecaseMock.Mocks) args {
				referralInput := &dto.ReferralInput{
					ReferralType: "DIAGNOSTICS",
					Tests:        []string{"VIA"},
					Facility: &dto.FacilityInput{
						Name:               "KNH",
						County:             "Nairobi",
						Contact:            "+254710100100",
						FHIROrganisationID: gofakeit.UUID(),
					},
				}

				return args{ctx: ctx, input: referralInput}
			},
			wantErr: true,
		},
		{
			name: "Sad Case: Unable to get tenant identifiers",
			setup: func(mh *usecaseMock.Mocks) args {
				referralInput.ReferralType = dto.TREATMENT

				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return nil, fmt.Errorf("error")
					})

				return args{ctx: ctx, input: referralInput}
			},
			wantErr: true,
		},
		{
			name: "Sad Case: referral date in the future is rejected",
			setup: func(mh *usecaseMock.Mocks) args {
				future := time.Now().AddDate(0, 0, 1)

				// A future-dated referral is rejected before any FHIR call, so no mocks are needed.
				input := &dto.ReferralInput{
					EncounterID:  gofakeit.UUID(),
					ReferralType: dto.TREATMENT,
					Tests:        []string{"VIA"},
					Facility: &dto.FacilityInput{
						Name:               "KNH",
						FHIROrganisationID: gofakeit.UUID(),
					},
					ReferralNote: "test",
					ReferralDate: &scalarutils.Date{Year: future.Year(), Month: int(future.Month()), Day: future.Day()},
				}

				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
		{
			name: "Happy Case: a provided past referral date is used as the authored date",
			setup: func(mh *usecaseMock.Mocks) args {
				past := time.Now().AddDate(0, 0, -3)
				wantDate := past.Format("2006-01-02")

				input := &dto.ReferralInput{
					EncounterID:  gofakeit.UUID(),
					ReferralType: dto.TREATMENT,
					Tests:        []string{"VIA"},
					Facility: &dto.FacilityInput{
						Name:               "KNH",
						FHIROrganisationID: gofakeit.UUID(),
					},
					ReferralNote: "test",
					ReferralDate: &scalarutils.Date{Year: past.Year(), Month: int(past.Month()), Day: past.Day()},
				}

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
				mh.FHIR.EXPECT().CreateFHIRServiceRequest(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, in domain.FHIRServiceRequestInput) (*domain.FHIRServiceRequestRelayPayload, error) {
						if in.AuthoredOn == nil {
							t.Errorf("expected AuthoredOn to be set from the provided referral date")
						} else if got := string(*in.AuthoredOn); len(got) < 10 || got[:10] != wantDate {
							t.Errorf("AuthoredOn = %q, want date prefix %q", got, wantDate)
						}

						return servicerequest, nil
					})
				mh.PubSub.EXPECT().NotifyCreatePatientReferralTask(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, data dto.PatientReferralTaskPayload) error {
						return nil
					})

				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		// NB: these subtests are intentionally not run in parallel — several cases mutate the
		// shared referralInput fixture, so parallel execution would race on it.
		t.Run(tt.name, func(t *testing.T) {
			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			got, err := clinicalUsecase.ReferPatient(args.ctx, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.ReferPatient() error = %v, wantErr %v", err, tt.wantErr)
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

func TestUseCasesClinicalImpl_GetPatientReferralDetails(t *testing.T) {
	ctx := context.Background()
	status := dto.ObservationStatusFinal
	valueConcept := "222"
	UUID := gofakeit.UUID()
	instant := "2023-05-28T14:20:00Z"
	currentTime := time.Now().Format(time.RFC3339)
	pagedFHIRObservation := &domain.PagedFHIRObservations{
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
				EffectiveDateTime: (*scalarutils.DateTime)(&currentTime),
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
				EffectiveInstant: (*scalarutils.Instant)(&instant),
			},
		},
		HasNextPage:     false,
		NextPageURL:     "",
		HasPreviousPage: false,
		PreviousPageURL: "",
		TotalCount:      0,
	}

	implicitRules := mock.Anything
	language := mock.Anything
	url := scalarutils.URL(gofakeit.URL())
	contentType := scalarutils.Code("pdf")
	title := gofakeit.BeerName()
	pagedFHIRDocRef := &domain.PagedFHIRDocumentReference{
		DocumentReferences: []domain.FHIRDocumentReference{
			{
				ID:            uuid.NewString(),
				Status:        domain.DocumentReferenceStatusEnumCurrent,
				ImplicitRules: &implicitRules,
				Language:      &language,
				Meta: &domain.FHIRMeta{
					VersionID:   mock.Anything,
					Source:      mock.Anything,
					LastUpdated: gofakeit.Date(),
				},
				Content: []domain.FHIRDocumentReferenceContent{
					{
						Attachment: domain.FHIRAttachment{
							URL:         &url,
							ContentType: &contentType,
							Title:       &title,
						},
					},
				},
			},
		},
		HasNextPage:     false,
		NextCursor:      "",
		HasPreviousPage: false,
		PreviousCursor:  "",
		TotalCount:      0,
	}

	type args struct {
		ctx              context.Context
		serviceRequestID string
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		want    *domain.PatientReferralDetails
		wantErr bool
	}{
		{
			name: "Happy Case: Successfully get patient referral details",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().SearchFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).RunAndReturn(
					func(ctx context.Context, bundleID, resourceType string, params map[string]any, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRResource, error) {
						return &domain.PagedFHIRResource{
							Resources: resourceBnundle(),
						}, nil
					},
				)
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					}).Twice()
				mh.FHIR.EXPECT().SearchPatientObservations(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID string, searchParameters map[string]interface{}, tenant dto.TenantIdentifiers, pagination serverutils.PaginationInput) (*domain.PagedFHIRObservations, error) {
						return pagedFHIRObservation, nil
					})
				mh.FHIR.EXPECT().SearchFHIRDocumentReference(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, searchParams map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRDocumentReference, error) {
						return pagedFHIRDocRef, nil
					})

				return args{ctx: ctx, serviceRequestID: uuid.NewString()}
			},
			wantErr: false,
		},
		{
			name: "Sad Case: missing service request id",
			setup: func(mh *usecaseMock.Mocks) args {
				return args{ctx: ctx}
			},
			wantErr: true,
		},
		{
			name: "Sad Case: invalid service request id",
			setup: func(mh *usecaseMock.Mocks) args {
				return args{ctx: ctx, serviceRequestID: mock.Anything}
			},
			wantErr: true,
		},
		{
			name: "Sad Case: fail to get fhir service request",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})

				mh.FHIR.EXPECT().SearchFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).RunAndReturn(
					func(ctx context.Context, bundleID, resourceType string, params map[string]any, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRResource, error) {
						return nil, fmt.Errorf("an error occurred")
					},
				)

				return args{ctx: ctx, serviceRequestID: uuid.NewString()}
			},
			wantErr: true,
		},
		{
			name: "Sad Case: fail to get tenant identifiers",
			setup: func(mh *usecaseMock.Mocks) args {

				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, serviceRequestID: uuid.NewString()}
			},
			wantErr: true,
		},
		{
			name: "Sad Case: fail to get document reference",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().SearchFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).RunAndReturn(
					func(ctx context.Context, bundleID, resourceType string, params map[string]any, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRResource, error) {
						return &domain.PagedFHIRResource{
							Resources: resourceBnundle(),
						}, nil
					},
				)

				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchPatientObservations(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID string, searchParameters map[string]interface{}, tenant dto.TenantIdentifiers, pagination serverutils.PaginationInput) (*domain.PagedFHIRObservations, error) {
						return pagedFHIRObservation, nil
					})
				mh.FHIR.EXPECT().SearchFHIRDocumentReference(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, searchParams map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRDocumentReference, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, serviceRequestID: uuid.NewString()}
			},
			wantErr: true,
		},
		{
			name: "Sad Case: fail to get patient observations",
			setup: func(mh *usecaseMock.Mocks) args {

				mh.FHIR.EXPECT().SearchFHIRResource(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).RunAndReturn(
					func(ctx context.Context, bundleID, resourceType string, params map[string]any, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRResource, error) {
						return &domain.PagedFHIRResource{
							Resources: resourceBnundle(),
						}, nil
					},
				)
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchPatientObservations(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID string, searchParameters map[string]interface{}, tenant dto.TenantIdentifiers, pagination serverutils.PaginationInput) (*domain.PagedFHIRObservations, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, serviceRequestID: uuid.NewString()}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			_, err := clinicalUsecase.GetPatientReferralDetails(args.ctx, args.serviceRequestID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.GetPatientReferralDetails() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func TestUseCasesClinicalImpl_CreateReferral(t *testing.T) {
	ctx := context.Background()
	input := &dto.CreateReferralInput{
		EncounterID:     gofakeit.UUID(),
		ReferralType:    "OUTBOUND",
		Urgency:         "asap",
		ClinicalHistory: "",
		ReferralDate:    scalarutils.Date{},
		Diagnosis:       "",
		Tests:           []string{"HPV"},
		Specialist:      "Gynaecologist",
		ReferredFrom: &dto.ReferralFacility{
			FHIROrganisationID: gofakeit.UUID(),
			Name:               "KNH",
			PhoneNumber:        "+2547123456",
			Email:              "knk@knh.com",
			Branch:             "Nairobi",
		},
		ReferredTo: dto.ReferralFacility{
			FHIROrganisationID: gofakeit.UUID(),
			Name:               "AKUH",
			PhoneNumber:        "+2547123456",
			Email:              "akuh@akuh.com",
			Branch:             "Nairobi",
		},
		UsageContext: dto.BreastCancerScreeningTypeEnum,
	}

	tenantIDs := &dto.TenantIdentifiers{
		OrganizationID: uuid.NewString(),
		FacilityID:     uuid.NewString(),
	}

	ID := uuid.NewString()
	orgName := mock.Anything
	organization := &domain.FHIROrganizationRelayPayload{
		Resource: &domain.FHIROrganization{
			ID:   &ID,
			Name: &orgName,
		},
	}

	ref := mock.Anything
	encounter := &domain.FHIREncounterRelayPayload{
		Resource: &domain.FHIREncounter{
			ID:     &ID,
			Status: domain.EncounterStatusEnumInProgress,
			Subject: &domain.FHIRReference{
				ID:        &ID,
				Reference: &ref,
				Display:   gofakeit.UUID(),
			},
		},
	}

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
			Category: []*domain.FHIRCodeableConcept{
				{
					Coding: []*domain.FHIRCoding{
						{
							Code:    &serviceCode,
							Display: gofakeit.UUID(),
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
		},
	}

	type args struct {
		ctx   context.Context
		input *dto.CreateReferralInput
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy Case: Successfully create a patient referral",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).RunAndReturn(
					func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
						return &domain.Concept{}, nil
					},
				)
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

				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "Happy Case: Successfully create a patient referral - with referred from retrieved from context",
			setup: func(mh *usecaseMock.Mocks) args {
				input := &dto.CreateReferralInput{
					EncounterID:     gofakeit.UUID(),
					ReferralType:    "OUTBOUND",
					Urgency:         "asap",
					ClinicalHistory: "",
					ReferralDate:    scalarutils.Date{},
					Diagnosis:       "",
					ReferredTo: dto.ReferralFacility{
						FHIROrganisationID: gofakeit.UUID(),
						Name:               "AKUH",
						PhoneNumber:        "+2547123456",
						Email:              "akuh@akuh.com",
						Branch:             "Nairobi",
					},
					UsageContext: dto.BreastCancerScreeningTypeEnum,
				}

				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).RunAndReturn(
					func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
						return &domain.Concept{}, nil
					},
				)
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

				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "Sad Case: unable to create a patient referral - unable to get tenant identifiers",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).RunAndReturn(
					func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
						return &domain.Concept{}, nil
					},
				)
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
		{
			name: "Sad Case: unable to get conditions",
			setup: func(mh *usecaseMock.Mocks) args {
				input := &dto.CreateReferralInput{
					EncounterID:     gofakeit.UUID(),
					ReferralType:    "OUTBOUND",
					Urgency:         "asap",
					ClinicalHistory: "",
					ReferralDate:    scalarutils.Date{},
					Diagnosis:       gofakeit.UUID(),
					ReferredTo: dto.ReferralFacility{
						FHIROrganisationID: gofakeit.UUID(),
						Name:               "AKUH",
						PhoneNumber:        "+2547123456",
						Email:              "akuh@akuh.com",
						Branch:             "Nairobi",
					},
				}

				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).RunAndReturn(
					func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
						return &domain.Concept{}, nil
					},
				)
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().GetFHIRCondition(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIRConditionRelayPayload, error) {
						return nil, fmt.Errorf("an error occurred")
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
			name: "Sad Case: Fail to get tenant meta tags",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).RunAndReturn(
					func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
						return &domain.Concept{}, nil
					},
				)
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
			name: "Sad Case: Fail to create service request",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).RunAndReturn(
					func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
						return &domain.Concept{}, nil
					},
				)
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

				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
		{
			name: "Sad Case: unable to create referral task",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
					})
				mh.OCL.EXPECT().GetConcept(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).RunAndReturn(
					func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
						return &domain.Concept{}, nil
					},
				)
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

				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
		{
			name: "Sad Case: Record referral in finished encounter",
			setup: func(mh *usecaseMock.Mocks) args {
				encounter.Resource.Status = domain.EncounterStatusEnumCompleted
				mh.FHIR.EXPECT().GetFHIREncounter(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
						return encounter, nil
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

			got, err := clinicalUsecase.CreateReferral(args.ctx, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.CreateReferral() error = %v, wantErr %v", err, tt.wantErr)
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
