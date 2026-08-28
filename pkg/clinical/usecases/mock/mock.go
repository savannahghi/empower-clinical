package mock

import (
	"context"
	"log/slog"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/utils"
	"github.com/savannahghi/empower-clinical/pkg/clinical/infrastructure"
	infraMock "github.com/savannahghi/empower-clinical/pkg/clinical/infrastructure/mock"
	advMock "github.com/savannahghi/empower-clinical/pkg/clinical/infrastructure/services/advantage/mock"
	mailMock "github.com/savannahghi/empower-clinical/pkg/clinical/infrastructure/services/mail/mock"
	mapperMock "github.com/savannahghi/empower-clinical/pkg/clinical/infrastructure/services/mapper/mock"
	pubsubMock "github.com/savannahghi/empower-clinical/pkg/clinical/infrastructure/services/pubsub/mock"
	uploadMock "github.com/savannahghi/empower-clinical/pkg/clinical/infrastructure/services/upload/mock"
	urlMock "github.com/savannahghi/empower-clinical/pkg/clinical/infrastructure/services/urlshortener/mock"
	fakeRepoFHIR "github.com/savannahghi/empower-clinical/pkg/clinical/repository/mock"
	"github.com/savannahghi/empower-clinical/pkg/clinical/usecases/base"
	"github.com/savannahghi/empower-clinical/pkg/clinical/usecases/clinical"
	"github.com/savannahghi/empower-clinical/pkg/clinical/usecases/foundation"
	specializedUsecase "github.com/savannahghi/empower-clinical/pkg/clinical/usecases/specialized"
)

type Mocks struct {
	Ext          *infraMock.MockBaseExtension
	FHIR         *fakeRepoFHIR.MockFHIR
	OCL          *infraMock.MockServiceOCL
	PubSub       *pubsubMock.MockServicePubsub
	Upload       *uploadMock.MockServiceUpload
	Advantage    *advMock.MockAdvantageService
	URLShortener *urlMock.MockIServiceURLShortener
	MailSender   *mailMock.MockMailSender
	Mapper       *mapperMock.MockTimelineMapper
}

func SetupMocks(t *testing.T) (AllUsecaseImpls, Mocks) {
	fakeExt := infraMock.NewMockBaseExtension(t)
	fakeFHIR := fakeRepoFHIR.NewMockFHIR(t)
	fakeOCL := infraMock.NewMockServiceOCL(t)
	fakePubSub := pubsubMock.NewMockServicePubsub(t)
	fakeUpload := uploadMock.NewMockServiceUpload(t)
	fakeAdvantage := advMock.NewMockAdvantageService(t)
	fakeURLShortener := urlMock.NewMockIServiceURLShortener(t)
	fakeMailSender := mailMock.NewMockMailSender(t)
	fakeMapper := mapperMock.NewMockTimelineMapper(t)

	infra := infrastructure.NewInfrastructureInteractor(fakeExt, fakeFHIR, fakeOCL, fakeUpload, fakePubSub, fakeAdvantage, fakeURLShortener, fakeMailSender, fakeMapper)
	foundationUsecase := foundation.NewFoundationImpl(infra, slog.Default())
	specialized := specializedUsecase.NewSpecializedImpl(infra, *foundationUsecase, slog.Default())
	clinical := clinical.NewClinicalImpl(infra, *foundationUsecase, slog.Default())
	baseUsecase := base.NewBaseImpl(infra, *clinical, *foundationUsecase, slog.Default())

	allUseCases := NewUseCasesClinicalImpl(specialized, foundationUsecase, baseUsecase, clinical)

	mocks := Mocks{
		Ext:          fakeExt,
		FHIR:         fakeFHIR,
		OCL:          fakeOCL,
		PubSub:       fakePubSub,
		Upload:       fakeUpload,
		Advantage:    fakeAdvantage,
		URLShortener: fakeURLShortener,
		MailSender:   fakeMailSender,
		Mapper:       fakeMapper,
	}

	return *allUseCases, mocks
}

type AllUsecaseImpls struct {
	specializedUsecase.SpecializedImpl
	foundation.FoundationImpl
	base.BaseImpl
	clinical.ClinicalImpl
}

func NewUseCasesClinicalImpl(
	specialized *specializedUsecase.SpecializedImpl,
	foundation *foundation.FoundationImpl,
	base *base.BaseImpl,
	clinical *clinical.ClinicalImpl,
) *AllUsecaseImpls {
	return &AllUsecaseImpls{
		*specialized,
		*foundation,
		*base,
		*clinical,
	}
}

func AddTenantIdentifierContext(ctx context.Context) context.Context {
	ctx = context.WithValue(ctx, utils.OrganizationIDContextKey, gofakeit.UUID())
	ctx = context.WithValue(ctx, utils.FacilityIDContextKey, gofakeit.UUID())

	return ctx
}
