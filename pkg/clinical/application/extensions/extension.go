package extensions

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/savannahghi/converterandformatter"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/profileutils"
	"github.com/savannahghi/pubsubtools"
	"github.com/savannahghi/serverutils"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/utils"
)

// ISCClientExtension represents the base ISC client
type ISCClientExtension interface {
	MakeRequest(ctx context.Context, method string, path string, body interface{}) (*http.Response, error)
}

// ISCExtensionImpl ...
type ISCExtensionImpl struct{}

// NewISCExtension initializes an ISC extension
func NewISCExtension() *ISCExtensionImpl {
	return &ISCExtensionImpl{}
}

// MakeRequest performs an inter service http request and returns a response
func (i *ISCExtensionImpl) MakeRequest(ctx context.Context, method string, path string, body interface{}) (*http.Response, error) {
	var isc interserviceclient.InterServiceClient
	return isc.MakeRequest(ctx, method, path, body)
}

// BaseExtensionImpl ...
type BaseExtensionImpl struct{}

// NewBaseExtensionImpl ...
func NewBaseExtensionImpl() *BaseExtensionImpl {
	return &BaseExtensionImpl{}
}

// GetLoggedInUser retrieves logged in user information
func (b *BaseExtensionImpl) GetLoggedInUser(ctx context.Context) (*profileutils.UserInfo, error) {
	return profileutils.GetLoggedInUser(ctx)
}

// GetLoggedInUserUID get the logged in user uid
func (b *BaseExtensionImpl) GetLoggedInUserUID(ctx context.Context) (string, error) {
	return firebasetools.GetLoggedInUserUID(ctx)
}

// NormalizeMSISDN validates the input phone number.
func (b *BaseExtensionImpl) NormalizeMSISDN(msisdn string) (*string, error) {
	return converterandformatter.NormalizeMSISDN(msisdn)
}

// LoadDepsFromYAML ...
func (b *BaseExtensionImpl) LoadDepsFromYAML() (*interserviceclient.DepsConfig, error) {
	return interserviceclient.LoadDepsFromYAML()
}

// SetupISCclient ...
func (b *BaseExtensionImpl) SetupISCclient(config interserviceclient.DepsConfig, serviceName string) (*interserviceclient.InterServiceClient, error) {
	return interserviceclient.SetupISCclient(config, serviceName)
}

// GetEnvVar ...
func (b *BaseExtensionImpl) GetEnvVar(envName string) (string, error) {
	return serverutils.GetEnvVar(envName)
}

// WriteJSONResponse writes the content supplied via the `source` parameter to
// the supplied http ResponseWriter. The response is returned with the indicated
// status.
func (b *BaseExtensionImpl) WriteJSONResponse(
	w http.ResponseWriter,
	source interface{},
	status int,
) {
	serverutils.WriteJSONResponse(w, source, status)
}

// ErrorMap turns the supplied error into a map with "error" as the key
func (b *BaseExtensionImpl) ErrorMap(err error) map[string]string {
	return serverutils.ErrorMap(err)
}

// GetTenantIdentifiers retrieves the tenant identifiers e.g OrganizationID, FacilityID from the context
func (b *BaseExtensionImpl) GetTenantIdentifiers(ctx context.Context) (*dto.TenantIdentifiers, error) {
	organizationID, err := GetOrganizationIDFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve organization ID from context: %w", err)
	}

	facilityID, err := GetFacilityIDFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve facility ID from context: %w", err)
	}

	return &dto.TenantIdentifiers{
		OrganizationID: organizationID,
		FacilityID:     facilityID,
	}, nil
}

// VerifyPubSubJWTAndDecodePayload confirms that there is a valid Google signed
// JWT and decodes the pubsub message payload into a struct.
//
// It's use will simplify & shorten the handler funcs that process Cloud Pubsub
// push notifications.
func (b *BaseExtensionImpl) VerifyPubSubJWTAndDecodePayload(w http.ResponseWriter, r *http.Request) (*pubsubtools.PubSubPayload, error) {
	return pubsubtools.VerifyPubSubJWTAndDecodePayload(w, r)
}

// GetPubSubTopic retrieves a pubsub topic from a pubsub payload.
//
// It follows a convention where the topic is sent as an attribute under the
// `topicID` key.
func (b *BaseExtensionImpl) GetPubSubTopic(m *pubsubtools.PubSubPayload) (string, error) {
	return pubsubtools.GetPubSubTopic(m)
}

// GetOrganizationIDFromContext is a function that retrieves the organization ID from the context.
func GetOrganizationIDFromContext(ctx context.Context) (string, error) {
	organizationID, ok := ctx.Value(utils.OrganizationIDContextKey).(string)
	if !ok {
		return "", errors.New("unable to get organization ID from context")
	}

	return organizationID, nil
}

// GetFacilityIDFromContext is a function that retrieves the facility ID from the context.
func GetFacilityIDFromContext(ctx context.Context) (string, error) {
	facilityID, ok := ctx.Value(utils.FacilityIDContextKey).(string)
	if !ok {
		return "", errors.New("unable to get facility ID from context")
	}

	return facilityID, nil
}
