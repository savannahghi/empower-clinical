package mock

import (
	"context"
	"net/http"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/profileutils"
	"github.com/savannahghi/pubsubtools"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/dto"
)

// FakeBaseExtension is an mock of the BaseExtension
type FakeBaseExtension struct {
	GetLoggedInUserFn    func(ctx context.Context) (*profileutils.UserInfo, error)
	GetLoggedInUserUIDFn func(ctx context.Context) (string, error)
	NormalizeMSISDNFn    func(msisdn string) (*string, error)
	LoadDepsFromYAMLFn   func() (*interserviceclient.DepsConfig, error)
	SetupISCclientFn     func(config interserviceclient.DepsConfig, serviceName string) (*interserviceclient.InterServiceClient, error)
	GetEnvVarFn          func(envName string) (string, error)
	ErrorMapFn           func(err error) map[string]string
	WriteJSONResponseFn  func(
		w http.ResponseWriter,
		source interface{},
		status int,
	)
	MockGetTenantIdentifiersFn            func(ctx context.Context) (*dto.TenantIdentifiers, error)
	MockVerifyPubSubJWTAndDecodePayloadFn func(w http.ResponseWriter, r *http.Request) (*pubsubtools.PubSubPayload, error)
	MockGetPubSubTopicFn                  func(m *pubsubtools.PubSubPayload) (string, error)
	MockGinContextFn                      func(w http.ResponseWriter)
	MockGetFacilityIDFromContextFn        func(ctx context.Context) (string, error)
}

// NewFakeBaseExtensionMock initializes a new instance of extension mock
func NewFakeBaseExtensionMock() *FakeBaseExtension {
	return &FakeBaseExtension{
		GetLoggedInUserFn: func(ctx context.Context) (*profileutils.UserInfo, error) {
			return &profileutils.UserInfo{
				DisplayName: gofakeit.BeerName(),
				Email:       "test@email.com",
				PhoneNumber: interserviceclient.TestUserPhoneNumber,
				PhotoURL:    "google.com/photo",
				ProviderID:  uuid.NewString(),
				UID:         uuid.NewString(),
			}, nil
		},
		GetLoggedInUserUIDFn: func(ctx context.Context) (string, error) {
			return uuid.NewString(), nil
		},
		NormalizeMSISDNFn: func(msisdn string) (*string, error) {
			p := interserviceclient.TestUserPhoneNumber
			return &p, nil
		},
		LoadDepsFromYAMLFn: func() (*interserviceclient.DepsConfig, error) {
			return &interserviceclient.DepsConfig{
				Staging: []interserviceclient.Dep{
					{
						DepName:       "mycarehub",
						DepRootDomain: "https://mycarehub-staging.savannahghi.org",
					},
				},
				Testing: []interserviceclient.Dep{
					{
						DepName:       "mycarehub",
						DepRootDomain: "https://mycarehub-testing.savannahghi.org",
					},
				},
				Production: []interserviceclient.Dep{
					{
						DepName:       "mycarehub",
						DepRootDomain: "https://mycarehub-prod.savannahghi.org",
					},
				},
			}, nil
		},
		SetupISCclientFn: func(config interserviceclient.DepsConfig, serviceName string) (*interserviceclient.InterServiceClient, error) {
			return &interserviceclient.InterServiceClient{
				Name:              "clinical",
				RequestRootDomain: "clinical.com",
			}, nil
		},
		GetEnvVarFn: func(envName string) (string, error) {
			return "test", nil
		},
		ErrorMapFn: func(err error) map[string]string {
			m := map[string]string{
				"key": "value",
			}

			return m
		},
		WriteJSONResponseFn: func(w http.ResponseWriter, source interface{}, status int) {},
		MockGetTenantIdentifiersFn: func(ctx context.Context) (*dto.TenantIdentifiers, error) {
			return &dto.TenantIdentifiers{
				OrganizationID: uuid.New().String(),
				FacilityID:     uuid.New().String(),
			}, nil
		},
		MockVerifyPubSubJWTAndDecodePayloadFn: func(w http.ResponseWriter, r *http.Request) (*pubsubtools.PubSubPayload, error) {
			return &pubsubtools.PubSubPayload{}, nil
		},
		MockGetPubSubTopicFn: func(m *pubsubtools.PubSubPayload) (string, error) {
			return "", nil
		},
		MockGetFacilityIDFromContextFn: func(ctx context.Context) (string, error) {
			return uuid.NewString(), nil
		},
	}
}

// GetLoggedInUser retrieves logged in user information
func (b *FakeBaseExtension) GetLoggedInUser(ctx context.Context) (*profileutils.UserInfo, error) {
	return b.GetLoggedInUserFn(ctx)
}

// GetLoggedInUserUID is a mock implementation of GetLoggedInUserUID method
func (b *FakeBaseExtension) GetLoggedInUserUID(ctx context.Context) (string, error) {
	return b.GetLoggedInUserUIDFn(ctx)
}

// NormalizeMSISDN is a mock implementation of NormalizeMSISDN method
func (b *FakeBaseExtension) NormalizeMSISDN(msisdn string) (*string, error) {
	return b.NormalizeMSISDNFn(msisdn)
}

// LoadDepsFromYAML is a mock implementation of LoadDepsFromYAML method
func (b *FakeBaseExtension) LoadDepsFromYAML() (*interserviceclient.DepsConfig, error) {
	return b.LoadDepsFromYAMLFn()
}

// SetupISCclient is a mock implementation of SetupISCclient method
func (b *FakeBaseExtension) SetupISCclient(config interserviceclient.DepsConfig, serviceName string) (*interserviceclient.InterServiceClient, error) {
	return b.SetupISCclientFn(config, serviceName)
}

// GetEnvVar is a mock implementation of GetEnvVar method
func (b *FakeBaseExtension) GetEnvVar(envName string) (string, error) {
	return b.GetEnvVarFn(envName)
}

// ErrorMap is a mock implementation of ErrorMapG method
func (b *FakeBaseExtension) ErrorMap(err error) map[string]string {
	return b.ErrorMapFn(err)
}

// WriteJSONResponse is a mock implementation of WriteJSONResponse method
func (b *FakeBaseExtension) WriteJSONResponse(
	_ http.ResponseWriter,
	_ interface{},
	_ int,
) {
}

// GetTenantIdentifiers mocks the implementation of getting the tenant identifiers
func (b *FakeBaseExtension) GetTenantIdentifiers(ctx context.Context) (*dto.TenantIdentifiers, error) {
	return b.MockGetTenantIdentifiersFn(ctx)
}

// VerifyPubSubJWTAndDecodePayload confirms that there is a valid Google signed
// JWT and decodes the pubsub message payload into a struct.
//
// It's use will simplify & shorten the handler funcs that process Cloud Pubsub
// push notifications.
func (b *FakeBaseExtension) VerifyPubSubJWTAndDecodePayload(w http.ResponseWriter, r *http.Request) (*pubsubtools.PubSubPayload, error) {
	return b.MockVerifyPubSubJWTAndDecodePayloadFn(w, r)
}

// GetPubSubTopic retrieves a pubsub topic from a pubsub payload.
//
// It follows a convention where the topic is sent as an attribute under the
// `topicID` key.
func (b *FakeBaseExtension) GetPubSubTopic(m *pubsubtools.PubSubPayload) (string, error) {
	return b.MockGetPubSubTopicFn(m)
}
