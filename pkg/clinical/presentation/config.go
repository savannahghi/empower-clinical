package presentation

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"regexp"
	"time"

	"cloud.google.com/go/pubsub/v2"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/Nerzal/gocloak/v13"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/chenyahui/gin-cache/persist"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	sloggin "github.com/samber/slog-gin"
	"github.com/savannahghi/authutils"
	hapifhirgo "github.com/savannahghi/hapi-fhir-go"
	"github.com/savannahghi/serverutils"
	silgotel "github.com/savannahghi/sil-gotel"
	"github.com/savannahghi/silurlshortener"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/savannahghi/empower-clinical/docs"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/extensions"
	"github.com/savannahghi/empower-clinical/pkg/clinical/infrastructure"
	fhir "github.com/savannahghi/empower-clinical/pkg/clinical/infrastructure/datastore/cloudhealthcare"
	hapifhir "github.com/savannahghi/empower-clinical/pkg/clinical/infrastructure/datastore/hapi_fhir"
	"github.com/savannahghi/empower-clinical/pkg/clinical/infrastructure/services/advantage"
	auth "github.com/savannahghi/empower-clinical/pkg/clinical/infrastructure/services/authutils"
	"github.com/savannahghi/empower-clinical/pkg/clinical/infrastructure/services/keycloak"
	"github.com/savannahghi/empower-clinical/pkg/clinical/infrastructure/services/mail"
	"github.com/savannahghi/empower-clinical/pkg/clinical/infrastructure/services/mapper"
	"github.com/savannahghi/empower-clinical/pkg/clinical/infrastructure/services/openconceptlab"
	pubsubmessaging "github.com/savannahghi/empower-clinical/pkg/clinical/infrastructure/services/pubsub"
	"github.com/savannahghi/empower-clinical/pkg/clinical/infrastructure/services/upload"
	"github.com/savannahghi/empower-clinical/pkg/clinical/infrastructure/services/urlshortener"
	resolver "github.com/savannahghi/empower-clinical/pkg/clinical/presentation/graph"
	"github.com/savannahghi/empower-clinical/pkg/clinical/presentation/graph/generated"
	"github.com/savannahghi/empower-clinical/pkg/clinical/presentation/rest"
	"github.com/savannahghi/empower-clinical/pkg/clinical/usecases"
	"github.com/savannahghi/empower-clinical/pkg/clinical/usecases/base"
	"github.com/savannahghi/empower-clinical/pkg/clinical/usecases/clinical"
	"github.com/savannahghi/empower-clinical/pkg/clinical/usecases/foundation"
	"github.com/savannahghi/empower-clinical/pkg/clinical/usecases/specialized"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

var allowedOriginPatterns = []string{
	`^(https?://)?(.+)-?.slade360\.co.ke$`,
	`^(https?://)?(.+)-?.slade360\.com$`,
	`^(https?://)?(.+)-?.healthcloud\.co.ke$`,
	`^(https?://)?(.+)-?.ncikenya\.go.ke$`,
	`^(https?://)?(.+)-?.slade360edi\.com$`,
	`^(https?://)?(.+)-?.savannahghi\.org$`,
	`^(https?://)?(.+)-?.bewell\.co.ke$`,
	`^(https?://)?(.+)-?.tiberbu\.health$`,
	`^(https?://)?(.+)-?.dha\.go.ke$`,
	`^https://.+\.web\.app$`,
}

// Compile the regex patterns into a slice of *regexp.Regexp
func compilePatterns(patterns []string) []*regexp.Regexp {
	var compiledPatterns []*regexp.Regexp

	for _, pattern := range patterns {
		compiledPattern := regexp.MustCompile(pattern)
		compiledPatterns = append(compiledPatterns, compiledPattern)
	}

	return compiledPatterns
}

// Check if the origin is allowed by matching it against the compiled regex patterns
func isAllowedOrigin(origin string, compiledPatterns []*regexp.Regexp) bool {
	for _, pattern := range compiledPatterns {
		if pattern.MatchString(origin) {
			return true
		}
	}

	return false
}

var (
	authServerEndpoint = serverutils.MustGetEnvVar("AUTHSERVER_ENDPOINT")
	clientID           = serverutils.MustGetEnvVar("CLIENT_ID")
	clientSecret       = serverutils.MustGetEnvVar("CLIENT_SECRET")
	username           = serverutils.MustGetEnvVar("AUTH_USERNAME")
	password           = serverutils.MustGetEnvVar("AUTH_PASSWORD")
	grantType          = serverutils.MustGetEnvVar("GRANT_TYPE")
)

// resolveHAPIAuthOption selects how the HAPI FHIR client authenticates, driven by HAPI_AUTH_MODE:
//   - "keycloak": a client-credentials bearer from a dedicated Keycloak client, validated with an
//     eager token mint so a wrong secret/realm or an unreachable Keycloak fails startup here
//     instead of 401ing on the first FHIR call.
//   - "basic" or unset: legacy HTTP Basic auth (HAPI_USERNAME/HAPI_PASSWORD).
//
// Unset defaults to basic so environments whose FHIR server does not yet accept Keycloak bearer
// tokens keep working unchanged; an environment opts in by setting HAPI_AUTH_MODE=keycloak.
func resolveHAPIAuthOption(ctx context.Context, logger *slog.Logger) (hapifhirgo.ClientOption, error) {
	mode := os.Getenv("HAPI_AUTH_MODE")

	switch mode {
	case "keycloak":
		hapiTokenURL := fmt.Sprintf(
			"%s/realms/%s/protocol/openid-connect/token",
			serverutils.MustGetEnvVar("KEYCLOAK_BASE_URL"),
			serverutils.MustGetEnvVar("KEYCLOAK_REALM"),
		)

		// Bound every token round-trip (initial mint + background refreshes) so a slow or
		// unreachable Keycloak can neither hang startup nor stall in-flight FHIR calls.
		hapiTokenHTTPClient := &http.Client{Timeout: 15 * time.Second}

		// Derive the token-source context from ctx (inherits its values) but drop cancellation:
		// the source lives for the whole process, so a shutdown cancel of ctx must not sever token
		// refreshes while in-flight FHIR calls are still draining.
		hapiTokenCtx := context.WithValue(context.WithoutCancel(ctx), oauth2.HTTPClient, hapiTokenHTTPClient)

		hapiToken := (&clientcredentials.Config{
			ClientID:     serverutils.MustGetEnvVar("HAPI_KEYCLOAK_CLIENT_ID"),
			ClientSecret: serverutils.MustGetEnvVar("HAPI_KEYCLOAK_CLIENT_SECRET"),
			TokenURL:     hapiTokenURL,
		}).TokenSource(hapiTokenCtx)

		if _, err := hapiToken.Token(); err != nil {
			return nil, fmt.Errorf("failed to mint initial hapi fhir keycloak token: %w", err)
		}

		logger.Info("hapi fhir auth mode: keycloak - initial bearer token minted successfully")

		return hapifhirgo.WithTokenProvider(func(context.Context) (string, error) {
			token, err := hapiToken.Token()
			if err != nil {
				return "", fmt.Errorf("failed to obtain hapi fhir keycloak token: %w", err)
			}

			return token.AccessToken, nil
		}), nil
	case "basic", "":
		logger.Info("hapi fhir auth mode: basic")

		return hapifhirgo.WithBasicAuth(
			serverutils.MustGetEnvVar("HAPI_USERNAME"),
			serverutils.MustGetEnvVar("HAPI_PASSWORD"),
		), nil
	default:
		return nil, fmt.Errorf("invalid HAPI_AUTH_MODE %q: expected \"keycloak\" or \"basic\"", mode)
	}
}

func StartServer(
	ctx context.Context,
	port string,
) {
	err := serverutils.Sentry()
	if err != nil {
		serverutils.LogStartupError(ctx, err)
	}

	baseExtension := extensions.NewBaseExtensionImpl()

	authServerConfig := authutils.Config{
		AuthServerEndpoint: authServerEndpoint,
		ClientID:           clientID,
		ClientSecret:       clientSecret,
		GrantType:          grantType,
		Username:           username,
		Password:           password,
	}

	authclient, err := authutils.NewClient(authServerConfig)
	if err != nil {
		serverutils.LogStartupError(ctx, err)
	}

	redisClient, err := initRedisClient()
	if err != nil {
		// fail fast if redis is not available
		log.Panicf("failed to initialize Redis client: %v", err)
	}

	logger := initializeLogger()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		serverutils.LogStartupError(ctx, fmt.Errorf("failed to initialize AWS sdk: %w", err))
	}

	amazonSESClient := ses.NewFromConfig(cfg, func(o *ses.Options) {
		o.Region = serverutils.MustGetEnvVar("SES_AWS_REGION")
	})

	hapiAuthOption, err := resolveHAPIAuthOption(ctx, logger)
	if err != nil {
		serverutils.LogStartupError(ctx, fmt.Errorf("failed to configure hapi fhir auth: %w", err))
		os.Exit(1)
	}

	hapiFHIR, err := hapifhirgo.NewClient(
		serverutils.MustGetEnvVar("HAPI_FHIR_BASE_URL"),
		hapiAuthOption,
		// The library sends JSON Patch bodies with a FHIR content type, which
		// every FHIR server rejects. Correct it in transit.
		hapifhirgo.WithTransport(hapifhir.NewPatchContentTypeTransport(nil)),
	)
	if err != nil {
		serverutils.LogStartupError(ctx, fmt.Errorf("failed to initialize hapi fhir: %w", err))
	}

	hapiRepo := hapifhir.NewHAPIFHIRImplementation(hapiFHIR)

	mailSvc := mail.NewMailSender(amazonSESClient)

	fhir := fhir.NewFHIRStoreImpl(redisClient, hapiRepo)
	ocl := openconceptlab.NewServiceOCL()

	upload := upload.NewServiceUpload(ctx)

	// Integration events go over NATS by default, which needs no cloud account.
	// Set PUBSUB_BACKEND=gcp for Google Cloud Pub/Sub.
	//
	// Either backend is a hard failure: these events carry referral tasks, care
	// plans, follow-ups and FHIR ID writebacks.
	var (
		pubsubSvc  pubsubmessaging.ServicePubsub
		natsPubSub *pubsubmessaging.NATSPubSub
	)

	if strings.EqualFold(os.Getenv("PUBSUB_BACKEND"), "gcp") {
		projectID := serverutils.MustGetEnvVar(serverutils.GoogleCloudProjectIDEnvVarName)

		pubSubClient, clientErr := pubsub.NewClient(ctx, projectID)
		if clientErr != nil {
			serverutils.LogStartupError(ctx, fmt.Errorf("unable to initialize pubsub client: %w", clientErr))
			os.Exit(1)
		}

		pubsubSvc, err = pubsubmessaging.NewServicePubSubMessaging(ctx, pubSubClient)
		if err != nil {
			serverutils.LogStartupError(ctx, fmt.Errorf("failed to initialize pubsub messaging service: %w", err))
			os.Exit(1)
		}
	} else {
		natsPubSub, err = pubsubmessaging.NewNATSPubSub(ctx, os.Getenv("NATS_URL"), "")
		if err != nil {
			serverutils.LogStartupError(ctx, fmt.Errorf("failed to initialize NATS messaging: %w", err))
			os.Exit(1)
		}

		pubsubSvc = natsPubSub

		defer natsPubSub.Close()
	}

	advantageSvc := advantage.NewServiceAdvantage(authclient)

	silURLShortnerClient, err := silurlshortener.NewURLShortener()
	if err != nil {
		serverutils.LogStartupError(ctx, err)
	}

	shortenerService := urlshortener.NewURLShortenerService(silURLShortnerClient)

	timelineMapper := mapper.NewTimelineMapper()

	infrastructure := infrastructure.NewInfrastructureInteractor(baseExtension, fhir, ocl, upload, pubsubSvc, advantageSvc, shortenerService, mailSvc, timelineMapper)

	foundation := foundation.NewFoundationImpl(infrastructure, logger)
	clinical := clinical.NewClinicalImpl(infrastructure, *foundation, logger)

	base := base.NewBaseImpl(infrastructure, *clinical, *foundation, logger)
	specialized := specialized.NewSpecializedImpl(infrastructure, *foundation, logger)

	allUsecases, err := usecases.NewUsecasesImpl(*foundation, *base, specialized, *clinical)
	if err != nil {
		serverutils.LogStartupError(ctx, err)
	}

	// Five of the topics published above are consumed by this service. On Cloud
	// Pub/Sub they return as a push to /pubsub; over NATS nothing delivers them
	// unless the service subscribes.
	if natsPubSub != nil {
		handlers := rest.NewPresentationHandlers(allUsecases, baseExtension, advantageSvc)

		subscriber, err := natsPubSub.Subscribe(ctx, handlers.HandlePubSubEvent)
		if err != nil {
			serverutils.LogStartupError(ctx, fmt.Errorf("failed to subscribe to clinical events: %w", err))
			os.Exit(1)
		}

		defer subscriber.Stop()
	}

	r := gin.Default()

	memoryStore := persist.NewMemoryStore(60 * time.Minute)

	config := keycloak.Config{
		ClientID:     serverutils.MustGetEnvVar("KEYCLOAK_CLIENT_ID"),
		ClientSecret: serverutils.MustGetEnvVar("KEYCLOAK_CLIENT_SECRET"),
		Realm:        serverutils.MustGetEnvVar("KEYCLOAK_REALM"),
		Timeout:      30 * time.Second,
		Logger:       logger,
	}

	gocloak := gocloak.NewClient(serverutils.MustGetEnvVar("KEYCLOAK_BASE_URL"))

	kcClient, err := keycloak.New(config, gocloak)
	if err != nil {
		serverutils.LogStartupError(ctx, fmt.Errorf("failed to create keycloak client: %w", err))
	}

	SetupRoutes(r, logger, memoryStore, authclient, kcClient, allUsecases, infrastructure)

	addr := fmt.Sprintf(":%v", port)

	if err := r.Run(addr); err != nil {
		serverutils.LogStartupError(ctx, err)
	}
}

func initRedisClient() (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: serverutils.MustGetEnvVar("REDIS_ADDRESS"),
	})

	ctx := context.Background()

	pingResponse, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	log.Printf("redis ping response: %v", pingResponse)

	return rdb, nil
}

func initializeLogger() *slog.Logger {
	handler := slog.Handler(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))

	serviceName := fmt.Sprintf("clinical-%v", serverutils.GetRunningEnvironment())

	attrs := []slog.Attr{
		{Key: "service", Value: slog.StringValue(serviceName)},
	}

	handler.WithAttrs(attrs)

	log := slog.New(handler)
	slog.SetDefault(log)

	return log
}

func SetupRoutes(
	r *gin.Engine,
	log *slog.Logger,
	cacheStore persist.CacheStore,
	authclient auth.OAuthClientService,
	keycloakClient *keycloak.Client,
	allUsecases *usecases.Usecases,
	infra infrastructure.Infrastructure,
) {
	compiledPatterns := compilePatterns(allowedOriginPatterns)

	r.Use(cors.New(cors.Config{
		AllowWildcard: true,
		AllowMethods:  []string{http.MethodPut, http.MethodPatch, http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodOptions},
		AllowHeaders: []string{
			"Accept",
			"Accept-Charset",
			"Accept-Language",
			"Accept-Encoding",
			"Origin",
			"Host",
			"User-Agent",
			"Content-Length",
			"Content-Type",
			"Authorization",
			"X-Authorization",
			"Clinical-Organization-ID",
			"Clinical-Facility-ID",
			"traceparent",
			"tracestate",
			"X-Active-Tenant",
			"Access-Control-Allow-Headers",
			"X-Variant",
		},
		ExposeHeaders:    []string{"Content-Length", "Link"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			// Specific localhost origins
			docs.SwaggerInfo.Host = origin

			docs.SwaggerInfo.BasePath = origin
			if origin == "http://localhost:5000" || origin == "http://localhost:4200" ||
				origin == "http://localhost:8000" || origin == "http://localhost:8080" {
				return true
			}

			allowed := isAllowedOrigin(origin, compiledPatterns)

			return allowed
		},
		MaxAge:          12 * time.Hour,
		AllowWebSockets: true,
	}))

	r.Use(sloggin.NewWithConfig(log, sloggin.Config{
		WithRequestBody:  true,
		WithResponseBody: true,
	}))
	r.Use(gin.Recovery())

	handlers := rest.NewPresentationHandlers(allUsecases, infra.BaseExtension, infra.AdvantageService)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	graphQL := r.Group("/graphql")
	graphQL.Use(rest.AuthenticationGinMiddleware(cacheStore, authclient, keycloakClient))
	graphQL.Use(rest.TenantIdentifierExtractionMiddleware(infra.FHIR))
	graphQL.Any("", GQLHandler(allUsecases))

	// Unauthenticated routes
	ide := r.Group("/ide")
	ide.Any("", playgroundHandler())

	r.POST("/pubsub", handlers.ReceivePubSubPushMessage)

	// Authenticated routes
	apis := r.Group("/api")

	questionnaireGroup := apis.Group("/questionnaire")
	{
		questionnaireGroup.GET("", handlers.FetchQuestionnaire)
		questionnaireGroup.POST("", handlers.Assessment)
	}

	apis.Use(rest.AuthenticationGinMiddleware(cacheStore, authclient, keycloakClient))
	apis.Use(silgotel.GinRequestMetrics("clinical-service-backend"))

	v1 := apis.Group("/v1")

	tenants := v1.Group("/tenants")
	tenants.POST("", handlers.RegisterTenant)

	facilities := v1.Group("/facilities")
	facilities.POST("", handlers.RegisterFacility)

	provision := v1.Group("/identity/v1/provision")
	{
		provision.POST("", handlers.ProvisionTenant)
		provision.GET("/:tenant-id", handlers.GetTenantProvisioningStatus)
	}

	// Multitenant routes
	v1.Use(rest.TenantIdentifierExtractionMiddleware(infra.FHIR))

	planDefinitionGroup := v1.Group("/plan-definition")
	{
		planDefinitionGroup.POST("", handlers.CreatePlanDefinition)
		planDefinitionGroup.GET("", handlers.ListPlanDefinition)
	}

	carePlanGroup := v1.Group("/care-plan")
	{
		carePlanGroup.POST("", handlers.CreatePatientCarePlan)
		carePlanGroup.GET("", handlers.FetchCarePlan)
	}

	questionnaire := v1.Group("/questionnaire")
	{
		questionnaire.POST("", handlers.LoadQuestionnaire)
	}

	questionnaireResponse := v1.Group("/questionnaire-response")
	{
		questionnaireResponse.GET("/:id", handlers.SimpleQuestionnaireResponse)
		questionnaireResponse.POST("", handlers.CreateQuestionnaireResponses)
	}

	upload := v1.Group("/media")
	{
		upload.POST("", handlers.UploadMedia)
		upload.GET("", handlers.ListPatientMedia)
	}

	questionnaireList := v1.Group("/questionnaires")
	questionnaireList.GET("", handlers.ListQuestionnaire)

	referralReport := v1.Group("/referral-report")
	referralReport.GET("", handlers.GenerateReferralReport)

	riskAssessment := v1.Group("/risk-assessment")
	{
		riskAssessment.GET("", handlers.ListRiskAssessment)
		riskAssessment.GET("/:id", handlers.RiskAssessment)
	}

	taskGroup := v1.Group("/task")
	{
		taskGroup.GET("", handlers.ListTasks)
		taskGroup.PATCH("/:id", handlers.PatchTask)
		taskGroup.GET("/:id", handlers.GetTaskByID)
	}

	patientReferral := v1.Group("")
	{
		patientReferral.POST("/refer", handlers.ReferPatient)
		patientReferral.POST("/referV2", handlers.CreateReferral)
		patientReferral.GET("/referral-details", handlers.ReferralDetails)
		patientReferral.GET("/referral", handlers.PatientReferralDetailsByID)
		patientReferral.POST("/referral-form", handlers.ShareReferralForm)
	}

	listLabOrders := v1.Group("/lab-orders")
	{
		listLabOrders.GET("", handlers.ListLabOrders)
		listLabOrders.GET("/:id", handlers.GetLabOrderByID)
		listLabOrders.POST("/test", handlers.GetLabOrderByID)
	}

	condition := v1.Group("/condition")
	{
		condition.POST("", handlers.CreateCondition)
		condition.POST("/oncology-diagnosis", handlers.OncologyDiagnosis)
		condition.GET("", handlers.ListPatientConditions)
		condition.POST("/treatment-enrollment", handlers.RecordTreatmentEnrollment)
		condition.PATCH("/treatment-enrollment/:id", handlers.UpdateTreatmentEnrollment)
	}

	episodesOfCare := v1.Group("/episode-of-care")
	{
		episodesOfCare.POST("", handlers.CreateEpisodeOfCare)
		episodesOfCare.PATCH("/:id", handlers.PatchEpisodeOfCare)
		episodesOfCare.PATCH("/:id/end", handlers.EndEpisodeOfCare)
		episodesOfCare.GET("/:id", handlers.GetEpisodeOfCare)
	}

	allergyIntolerance := v1.Group("/allergyintolerance")
	{
		allergyIntolerance.POST("", handlers.CreateAllergyIntolerance)
		allergyIntolerance.GET("/search", handlers.SearchAllergy)
		allergyIntolerance.GET("/:id", handlers.GetAllergyIntolerance)
		allergyIntolerance.GET("", handlers.ListPatientAllergies)
	}

	patient := v1.Group("/patient")
	{
		patient.POST("", handlers.CreatePatient)
		patient.DELETE("/:id", handlers.DeletePatient)
		patient.GET("/medical-data/:id", handlers.GetMedicalData)
		patient.GET("/:id/timeline", handlers.GetPatientTimeline)
		patient.GET("/:id/banner", handlers.PatientBanner)
		patient.PATCH("/:id", handlers.PatchPatient)
	}

	encounters := v1.Group("/encounter")
	{
		encounters.POST("", handlers.StartEncounter)
		encounters.PATCH("/:id", handlers.PatchEncounter)
		encounters.PATCH("/:id/end", handlers.EndEncounter)
		encounters.PATCH("/end-screening/:id", handlers.EndScreening)
		encounters.GET("/:id/associated-resources", handlers.EncounterAssociatedResources)
		encounters.GET("", handlers.ListPatientEncounters)
	}

	screeningReport := v1.Group("")
	{
		screeningReport.GET("/screening-report", handlers.ScreeningReport)
	}

	composition := v1.Group("/compositions")
	{
		composition.POST("", handlers.CreateComposition)
		composition.PATCH("/:id", handlers.PatchComposition)
		composition.GET("", handlers.ListPatientCompositions)
	}

	consent := v1.Group("/consent")
	{
		consent.POST("", handlers.RecordConsent)
	}

	observation := v1.Group("/observations")
	{
		observation.POST("", handlers.RecordObservation)
		observation.GET("", handlers.GetObservations)
		observation.GET("/:id", handlers.Observation)
		observation.PATCH("/:id", handlers.PatchPatientObservations)
	}

	diagnostics := v1.Group("")
	{
		diagnostics.POST("/tests", handlers.RecordTests)
		diagnostics.POST("/lab-order", handlers.CreateLabOrderResult)
		diagnostics.POST("/test-results", handlers.RecordTestResult)
	}

	appointmentGroup := v1.Group("/appointment")
	{
		appointmentGroup.POST("", handlers.ScheduleAppointment)
	}

	medicationGroup := v1.Group("/medication")
	{
		medicationGroup.POST("/prescription", handlers.CreatePrescription)
		medicationGroup.PATCH("/:id", handlers.UpdateMedication)
		medicationGroup.POST("", handlers.CreateMedication)
		medicationGroup.GET("/:id", handlers.FetchMedicationByID)
		medicationGroup.GET("/medication-request/:id", handlers.FetchMedicationRequestByID)
		medicationGroup.GET("/medication-request", handlers.ListMedicationRequests)
	}

	diagnosticReport := v1.Group("/diagnostic-report")
	diagnosticReport.DELETE("", handlers.DeleteTests)
}

// GQLHandler sets up a GraphQL resolver.
func GQLHandler(usecases *usecases.Usecases) gin.HandlerFunc {
	resolver, err := resolver.NewResolver(usecases)
	if err != nil {
		log.Panicf("failed to start graph resolver: %s", err)
	}

	server := handler.New(
		generated.NewExecutableSchema(
			generated.Config{
				Resolvers: resolver,
			},
		),
	)

	server.AddTransport(transport.Options{})
	server.AddTransport(transport.GET{})
	server.AddTransport(transport.POST{})
	server.AddTransport(transport.MultipartForm{})
	server.Use(extension.Introspection{})

	return func(ctx *gin.Context) {
		server.ServeHTTP(ctx.Writer, ctx.Request)
	}
}

func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL IDE", "/graphql")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
