package foundation_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/scalarutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/empower-clinical/pkg/clinical/domain"
	usecaseMock "github.com/savannahghi/empower-clinical/pkg/clinical/usecases/mock"
)

func TestUseCasesClinicalImpl_GetTenantMetaTags(t *testing.T) {
	ctx := context.Background()
	ID := uuid.NewString()
	tenantIDs := &dto.TenantIdentifiers{
		OrganizationID: ID,
		FacilityID:     ID,
	}
	orgName := mock.Anything
	organization := &domain.FHIROrganizationRelayPayload{
		Resource: &domain.FHIROrganization{
			ID:   &ID,
			Name: &orgName,
		},
	}

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		args    args
		wantErr bool
	}{
		{
			name: "happy case: get tenant org from context",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return organization, nil
					}).Twice()
				return args{ctx: ctx}
			},
			wantErr: false,
		},
		{
			name: "sad case: missing tenant org in context",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return nil, fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx}
			},
			wantErr: true,
		},
		{
			name: "sad case: error retrieving organisation",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx}
			},
			wantErr: true,
		},
		{
			name: "sad case: error retrieving facility",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return organization, nil
					}).Once()
				mh.FHIR.EXPECT().GetFHIROrganization(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error) {
						return nil, fmt.Errorf("an error occurred")
					}).Once()
				return args{ctx: ctx}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			got, err := clinicalUsecase.GetTenantMetaTags(args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.GetTenantMetaTags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && got != nil {
				t.Errorf("expected result to be nil for %v", tt.name)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("expected result not to be nil for %v", tt.name)
				return
			}
		})
	}
}

func TestUseCasesClinicalImpl_SearchPatientReferences(t *testing.T) {
	ctx := context.Background()
	tenant := dto.TenantIdentifiers{
		OrganizationID: uuid.NewString(),
		FacilityID:     uuid.NewString(),
	}

	type args struct {
		ctx    context.Context
		term   string
		tenant dto.TenantIdentifiers
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		want    []string
		wantErr bool
	}{
		{
			name: "happy case: resolves and de-duplicates matching patients",
			setup: func(mh *usecaseMock.Mocks) args {
				// The name and identifier searches both return the same two patients; the result
				// must be de-duplicated.
				mh.FHIR.EXPECT().SearchFHIRResource(mock.Anything, mock.Anything, "Patient", mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID, resourceType string, params map[string]any, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRResource, error) {
						return &domain.PagedFHIRResource{Resources: []map[string]interface{}{
							{"id": "patient-1"},
							{"id": "patient-2"},
						}}, nil
					})

				return args{ctx: ctx, term: "jane", tenant: tenant}
			},
			want:    []string{"Patient/patient-1", "Patient/patient-2"},
			wantErr: false,
		},
		{
			name: "happy case: returns an empty slice when nothing matches",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().SearchFHIRResource(mock.Anything, mock.Anything, "Patient", mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID, resourceType string, params map[string]any, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRResource, error) {
						return &domain.PagedFHIRResource{Resources: []map[string]interface{}{}}, nil
					})

				return args{ctx: ctx, term: "ghost", tenant: tenant}
			},
			want:    []string{},
			wantErr: false,
		},
		{
			name: "happy case: a phone-like term also runs a normalized phone search",
			setup: func(mh *usecaseMock.Mocks) args {
				// The term is a valid phone number, so name, identifier AND phone searches run. Each
				// field returns a distinct patient, proving all three searches were issued.
				mh.FHIR.EXPECT().SearchFHIRResource(mock.Anything, mock.Anything, "Patient", mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID, resourceType string, params map[string]any, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRResource, error) {
						switch {
						case params["name"] != nil:
							return &domain.PagedFHIRResource{Resources: []map[string]interface{}{{"id": "name-match"}}}, nil
						case params["identifier"] != nil:
							return &domain.PagedFHIRResource{Resources: []map[string]interface{}{{"id": "identifier-match"}}}, nil
						case params["phone"] != nil:
							return &domain.PagedFHIRResource{Resources: []map[string]interface{}{{"id": "phone-match"}}}, nil
						default:
							return &domain.PagedFHIRResource{}, nil
						}
					})

				return args{ctx: ctx, term: "+254712345678", tenant: tenant}
			},
			want:    []string{"Patient/name-match", "Patient/identifier-match", "Patient/phone-match"},
			wantErr: false,
		},
		{
			name: "sad case: propagates the patient search error",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.FHIR.EXPECT().SearchFHIRResource(mock.Anything, mock.Anything, "Patient", mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, bundleID, resourceType string, params map[string]any, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRResource, error) {
						return nil, fmt.Errorf("an error occurred")
					})

				return args{ctx: ctx, term: "jane", tenant: tenant}
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clinicalUsecase, mocks := usecaseMock.SetupMocks(t)
			args := tt.setup(&mocks)

			got, err := clinicalUsecase.SearchPatientReferences(args.ctx, args.term, args.tenant)
			if (err != nil) != tt.wantErr {
				t.Errorf("SearchPatientReferences() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

// TestUseCasesClinicalImpl_SearchPatientReferences_capsResults verifies the match set is bounded so
// a common term cannot produce an unbounded reference IN-list.
func TestUseCasesClinicalImpl_SearchPatientReferences_capsResults(t *testing.T) {
	clinicalUsecase, mocks := usecaseMock.SetupMocks(t)

	// A single field search returns more patients than the cap; the helper must stop at the cap.
	resources := make([]map[string]interface{}, 0, 150)
	for i := 0; i < 150; i++ {
		resources = append(resources, map[string]interface{}{"id": fmt.Sprintf("patient-%d", i)})
	}

	mocks.FHIR.EXPECT().SearchFHIRResource(mock.Anything, mock.Anything, "Patient", mock.Anything, mock.Anything, mock.Anything).
		RunAndReturn(func(ctx context.Context, bundleID, resourceType string, params map[string]any, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRResource, error) {
			return &domain.PagedFHIRResource{Resources: resources}, nil
		})

	tenant := dto.TenantIdentifiers{
		OrganizationID: uuid.NewString(),
		FacilityID:     uuid.NewString(),
	}

	got, err := clinicalUsecase.SearchPatientReferences(context.Background(), "jane", tenant)
	assert.NoError(t, err)
	assert.Len(t, got, 100)
}

func TestUseCasesClinicalImpl_CheckPatientExistenceUsingPhoneNumber(t *testing.T) {
	ctx := context.Background()
	input := domain.SimplePatientRegistrationInput{
		PhoneNumbers: []*domain.PhoneNumberInput{
			{
				Msisdn:             gofakeit.Phone(),
				VerificationCode:   "1234",
				IsUssd:             false,
				CommunicationOptIn: false,
			},
		},
	}
	ID := uuid.NewString()
	tenantIDs := &dto.TenantIdentifiers{
		OrganizationID: ID,
		FacilityID:     ID,
	}

	patientsconn := domain.PatientConnection{
		Edges: []*domain.PatientEdge{
			{
				Cursor: mock.Anything,
				Node: &domain.FHIRPatient{
					ID: &ID,
				},
			},
		},
	}

	type args struct {
		ctx   context.Context
		input domain.SimplePatientRegistrationInput
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "happy case: nil inputs",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchFHIRPatient(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, searchParams string, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PatientConnection, error) {
						return &patientsconn, nil
					})
				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "happy case: contacts to contact point",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchFHIRPatient(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, searchParams string, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PatientConnection, error) {
						return &patientsconn, nil
					})
				return args{ctx: ctx, input: input}
			},
			wantErr: false,
		},
		{
			name: "sad case: missing FHIR patient",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})
				mh.FHIR.EXPECT().SearchFHIRPatient(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, searchParams string, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PatientConnection, error) {
						return nil, fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
		{
			name: "sad case: missing tenant org in context",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return nil, fmt.Errorf("an error occurred")
					})
				return args{ctx: ctx, input: input}
			},
			wantErr: true,
		},
		{
			name: "sad case: invalid phone",
			setup: func(mh *usecaseMock.Mocks) args {
				mh.Ext.EXPECT().GetTenantIdentifiers(mock.Anything).
					RunAndReturn(func(ctx context.Context) (*dto.TenantIdentifiers, error) {
						return tenantIDs, nil
					})

				input := domain.SimplePatientRegistrationInput{
					PhoneNumbers: []*domain.PhoneNumberInput{
						{
							Msisdn:             "0722",
							VerificationCode:   "1234",
							IsUssd:             false,
							CommunicationOptIn: false,
						},
					},
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

			got, err := clinicalUsecase.CheckPatientExistenceUsingPhoneNumber(args.ctx, args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.CheckPatientExistenceUsingPhoneNumber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && got != false {
				t.Errorf("expected result to be nil for %v", tt.name)
				return
			}
		})
	}
}

func TestUseCasesClinicalImpl_ContactsToContactPointInput(t *testing.T) {
	ctx := context.Background()
	email := gofakeit.Email()
	invalidEmail := "gofakeit.Email()"
	phones := []*domain.PhoneNumberInput{
		{
			Msisdn:             gofakeit.Phone(),
			VerificationCode:   "1234",
			IsUssd:             false,
			CommunicationOptIn: false,
		},
	}
	emails := []*domain.EmailInput{
		{
			Email:              &email,
			CommunicationOptIn: false,
		},
	}

	type args struct {
		ctx    context.Context
		phones []*domain.PhoneNumberInput
		emails []*domain.EmailInput
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "happy case: nil inputs",
			setup: func(mh *usecaseMock.Mocks) args {
				return args{ctx: ctx, phones: nil, emails: nil}
			},
			wantErr: false,
		},
		{
			name: "happy case: contacts to contact point",
			setup: func(mh *usecaseMock.Mocks) args {
				return args{ctx: ctx, phones: phones, emails: emails}
			},
			wantErr: false,
		},
		{
			name: "sad case: invalid phone",
			setup: func(mh *usecaseMock.Mocks) args {
				phones := []*domain.PhoneNumberInput{
					{
						Msisdn:             "0722",
						VerificationCode:   "1234",
						IsUssd:             false,
						CommunicationOptIn: false,
					},
				}
				return args{ctx: ctx, phones: phones, emails: emails}
			},
			wantErr: true,
		},
		{
			name: "sad case: invalid email",
			setup: func(mh *usecaseMock.Mocks) args {
				emails := []*domain.EmailInput{
					{
						Email:              &invalidEmail,
						CommunicationOptIn: false,
					},
				}
				return args{ctx: ctx, phones: phones, emails: emails}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			got, err := clinicalUsecase.ContactsToContactPointInput(args.ctx, args.phones, args.emails)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.ContactsToContactPointInput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && got != nil {
				t.Errorf("expected result to be nil for %v", tt.name)
				return
			}
		})
	}
}

func TestUseCasesClinicalImpl_SimplePatientRegistrationInputToPatientInput(t *testing.T) {
	ctx := context.Background()
	email := gofakeit.Email()
	invalidEmail := "invalid"
	address := gofakeit.BS()
	input := domain.SimplePatientRegistrationInput{
		ID: gofakeit.UUID(),
		Names: []*domain.NameInput{
			{
				FirstName: gofakeit.Name(),
				LastName:  gofakeit.Name(),
			},
		},
		IdentificationDocuments: []*domain.IdentificationDocument{
			{
				DocumentType:   domain.IDDocumentTypePassport,
				DocumentNumber: gofakeit.SSN(),
			},
		},
		BirthDate: &scalarutils.Date{
			Year:  2000,
			Month: 10,
			Day:   10,
		},
		PhoneNumbers: []*domain.PhoneNumberInput{},
		Photos: []*domain.PhotoInput{
			{
				PhotoContentType: enumutils.ContentTypeJpg,
				PhotoBase64data:  "qweqwdwedwed",
				PhotoFilename:    "test",
			},
		},
		Emails: []*domain.EmailInput{
			{
				Email:              &email,
				CommunicationOptIn: false,
			},
		},
		PhysicalAddresses: []*domain.PhysicalAddress{
			{
				MapsCode:        "123",
				PhysicalAddress: &address,
			},
		},
		PostalAddresses: []*domain.PostalAddress{
			{
				PostalAddress: &address,
				PostalCode:    "1234",
			},
		},
		Gender:        "",
		Active:        true,
		MaritalStatus: "",
		Languages:     []enumutils.Language{"en"},
		ReplicateUSSD: false,
	}

	type args struct {
		ctx            context.Context
		input          domain.SimplePatientRegistrationInput
		organizationID string
	}
	tests := []struct {
		name    string
		setup   func(mh *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "happy case: fhir patient input",
			setup: func(mh *usecaseMock.Mocks) args {
				return args{ctx: ctx, input: input, organizationID: gofakeit.UUID()}
			},
			wantErr: false,
		},
		{
			name: "sad case: invalid email",
			setup: func(mh *usecaseMock.Mocks) args {
				input.Emails = []*domain.EmailInput{
					{
						Email:              &invalidEmail,
						CommunicationOptIn: false,
					},
				}
				return args{ctx: ctx, input: input, organizationID: gofakeit.UUID()}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clinicalUsecase, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			got, err := clinicalUsecase.SimplePatientRegistrationInputToPatientInput(args.ctx, args.input, args.organizationID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.SimplePatientRegistrationInputToPatientInput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && got != nil {
				t.Errorf("expected result to be nil for %v", tt.name)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("expected result not to be nil for %v", tt.name)
				return
			}
		})
	}
}

func TestUseCasesClinicalImpl_ConceptMapper(t *testing.T) {
	type args struct {
		concept dto.ObservationConceptEnum
	}
	tests := []struct {
		name    string
		setup   func(_ *usecaseMock.Mocks) args
		wantErr bool
	}{
		{
			name: "Happy case: Valid concept BMI",
			setup: func(_ *usecaseMock.Mocks) args {
				return args{concept: "BMI"}
			},
			wantErr: false,
		},
		{
			name: "Happy case: Valid concept TEMPERATURE",
			setup: func(_ *usecaseMock.Mocks) args {
				return args{concept: "TEMPERATURE"}
			},
			wantErr: false,
		},
		{
			name: "Happy case: Valid concept HEIGHT",
			setup: func(_ *usecaseMock.Mocks) args {
				return args{concept: "HEIGHT"}
			},
			wantErr: false,
		},
		{
			name: "Happy case: Valid concept WEIGHT",
			setup: func(_ *usecaseMock.Mocks) args {
				return args{concept: "WEIGHT"}
			},
			wantErr: false,
		},
		{
			name: "Happy case: Valid concept RESPIRATORY_RATE",
			setup: func(_ *usecaseMock.Mocks) args {
				return args{concept: "RESPIRATORY_RATE"}
			},
			wantErr: false,
		},
		{
			name: "Happy case: Valid concept PULSE_RATE",
			setup: func(_ *usecaseMock.Mocks) args {
				return args{concept: "PULSE_RATE"}
			},
			wantErr: false,
		},
		{
			name: "Happy case: Valid concept BLOOD_PRESSURE",
			setup: func(_ *usecaseMock.Mocks) args {
				return args{concept: "BLOOD_PRESSURE"}
			},
			wantErr: false,
		},
		{
			name: "Happy case: Valid concept VIRAL_LOAD",
			setup: func(_ *usecaseMock.Mocks) args {
				return args{concept: "VIRAL_LOAD"}
			},
			wantErr: false,
		},
		{
			name: "Happy case: Valid concept MUAC",
			setup: func(_ *usecaseMock.Mocks) args {
				return args{concept: "MUAC"}
			},
			wantErr: false,
		},
		{
			name: "Happy case: Valid concept OXYGEN_SATURATION",
			setup: func(_ *usecaseMock.Mocks) args {
				return args{concept: "OXYGEN_SATURATION"}
			},
			wantErr: false,
		},
		{
			name: "Happy case: Valid concept BLOOD_SUGAR",
			setup: func(_ *usecaseMock.Mocks) args {
				return args{concept: "BLOOD_SUGAR"}
			},
			wantErr: false,
		},
		{
			name: "Happy case: Valid concept LAST_MENSTRUAL_PERIOD",
			setup: func(_ *usecaseMock.Mocks) args {
				return args{concept: "LAST_MENSTRUAL_PERIOD"}
			},
			wantErr: false,
		},
		{
			name: "Happy case: Valid concept DIASTOLIC_BLOOD_PRESSURE",
			setup: func(_ *usecaseMock.Mocks) args {
				return args{concept: "DIASTOLIC_BLOOD_PRESSURE"}
			},
			wantErr: false,
		},
		{
			name: "Happy case: Valid concept COLPOSCOPY",
			setup: func(_ *usecaseMock.Mocks) args {
				return args{concept: "COLPOSCOPY"}
			},
			wantErr: false,
		},
		{
			name: "Happy case: Valid concept VIA",
			setup: func(_ *usecaseMock.Mocks) args {
				return args{concept: "VIA"}
			},
			wantErr: false,
		},
		{
			name: "Happy case: Valid concept IMMUNO_HISTO_CHEMISTRY",
			setup: func(_ *usecaseMock.Mocks) args {
				return args{concept: "IMMUNO_HISTO_CHEMISTRY"}
			},
			wantErr: false,
		},
		{
			name: "Happy case: Valid concept POST_COITAL_BLEEDING",
			setup: func(_ *usecaseMock.Mocks) args {
				return args{concept: "POST_COITAL_BLEEDING"}
			},
			wantErr: false,
		},
		{
			name: "Happy case: Valid concept HISTORY_OF_PRESENTING_ILLNESS",
			setup: func(_ *usecaseMock.Mocks) args {
				return args{concept: "HISTORY_OF_PRESENTING_ILLNESS"}
			},
			wantErr: false,
		},
		{
			name: "Happy case: Valid concept PAST_MEDICAL_AND_SURGICAL_HISTORY",
			setup: func(_ *usecaseMock.Mocks) args {
				return args{concept: "PAST_MEDICAL_AND_SURGICAL_HISTORY"}
			},
			wantErr: false,
		},
		{
			name: "Happy case: Valid concept CHIEF_COMPLAINT",
			setup: func(_ *usecaseMock.Mocks) args {
				return args{concept: "CHIEF_COMPLAINT"}
			},
			wantErr: false,
		},
		{
			name: "Happy case: Valid concept CBE",
			setup: func(_ *usecaseMock.Mocks) args {
				return args{concept: "CBE"}
			},
			wantErr: false,
		},
		{
			name: "Happy case: Valid concept GENERAL_EXAMINATION",
			setup: func(_ *usecaseMock.Mocks) args {
				return args{concept: "GENERAL_EXAMINATION"}
			},
			wantErr: false,
		},
		{
			name: "Happy case: Valid concept FAMILY_AND_SOCIAL_HISTORY",
			setup: func(_ *usecaseMock.Mocks) args {
				return args{concept: "FAMILY_AND_SOCIAL_HISTORY"}
			},
			wantErr: false,
		},
		{
			name: "Happy case: Valid concept HPV",
			setup: func(_ *usecaseMock.Mocks) args {
				return args{concept: "HPV"}
			},
			wantErr: false,
		},
		{
			name: "Sad case: Invalid concept",
			setup: func(_ *usecaseMock.Mocks) args {
				return args{concept: mock.Anything}
			},
			wantErr: true,
		},
		{
			name: "Sad case: Missing concept",
			setup: func(_ *usecaseMock.Mocks) args {
				return args{concept: ""}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clinicalUsecases, mock := usecaseMock.SetupMocks(t)
			args := tt.setup(&mock)

			_, _, err := clinicalUsecases.ConceptMapper(args.concept)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.ConceptMapper() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
