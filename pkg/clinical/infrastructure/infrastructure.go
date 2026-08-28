package infrastructure

import (
	"context"
	"io"
	"net/http"
	"net/url"

	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/profileutils"
	"github.com/savannahghi/pubsubtools"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/empower-clinical/pkg/clinical/domain"
	"github.com/savannahghi/empower-clinical/pkg/clinical/infrastructure/services/advantage"
	"github.com/savannahghi/empower-clinical/pkg/clinical/infrastructure/services/mail"
	"github.com/savannahghi/empower-clinical/pkg/clinical/infrastructure/services/mapper"
	pubsubmessaging "github.com/savannahghi/empower-clinical/pkg/clinical/infrastructure/services/pubsub"
	"github.com/savannahghi/empower-clinical/pkg/clinical/infrastructure/services/upload"
	"github.com/savannahghi/empower-clinical/pkg/clinical/infrastructure/services/urlshortener"
	"github.com/savannahghi/empower-clinical/pkg/clinical/repository"
)

// ServiceOCL ...
type ServiceOCL interface {
	MakeRequest(method string, path string, params url.Values, body io.Reader) (*http.Response, error)
	ListConcepts(
		ctx context.Context, org []string, source []string, verbose bool, q *string,
		sortAsc *string, sortDesc *string, conceptClass *string, dataType *string,
		locale *string, includeRetired *bool,
		includeMappings *bool, includeInverseMappings *bool, paginationInput *dto.Pagination) (*domain.ConceptPage, error)
	GetConcept(
		ctx context.Context, org string, source string, concept string,
		includeMappings bool, includeInverseMappings bool) (*domain.Concept, error)
}

// BaseExtension is an interface that represents some methods in base
// The `onboarding` service has a dependency on `base` library.
// Our first step to making some functions are testable is to remove the base dependency.
// This can be achieved with the below interface.
type BaseExtension interface {
	GetLoggedInUser(ctx context.Context) (*profileutils.UserInfo, error)
	GetLoggedInUserUID(ctx context.Context) (string, error)
	GetTenantIdentifiers(ctx context.Context) (*dto.TenantIdentifiers, error)
	NormalizeMSISDN(msisdn string) (*string, error)
	LoadDepsFromYAML() (*interserviceclient.DepsConfig, error)
	SetupISCclient(config interserviceclient.DepsConfig, serviceName string) (*interserviceclient.InterServiceClient, error)
	GetEnvVar(envName string) (string, error)
	ErrorMap(err error) map[string]string
	WriteJSONResponse(
		w http.ResponseWriter,
		source interface{},
		status int,
	)
	VerifyPubSubJWTAndDecodePayload(w http.ResponseWriter, r *http.Request) (*pubsubtools.PubSubPayload, error)
	GetPubSubTopic(m *pubsubtools.PubSubPayload) (string, error)
}

// Infrastructure ...
type Infrastructure struct {
	FHIR                 repository.FHIR
	OpenConceptLab       ServiceOCL
	BaseExtension        BaseExtension
	Upload               upload.ServiceUpload
	Pubsub               pubsubmessaging.ServicePubsub
	AdvantageService     advantage.AdvantageService
	URLShorteningService urlshortener.IServiceURLShortener
	Mail                 mail.MailSender
	Mapper               mapper.TimelineMapper
}

// NewInfrastructureInteractor initializes a new Infrastructure
func NewInfrastructureInteractor(
	ext BaseExtension,
	fhir repository.FHIR,
	openconceptlab ServiceOCL,
	upload upload.ServiceUpload,
	pubsub pubsubmessaging.ServicePubsub,
	advantage advantage.AdvantageService,
	urlshortener urlshortener.IServiceURLShortener,
	mail mail.MailSender,
	mapper mapper.TimelineMapper,
) Infrastructure {
	return Infrastructure{
		fhir,
		openconceptlab,
		ext,
		upload,
		pubsub,
		advantage,
		urlshortener,
		mail,
		mapper,
	}
}
