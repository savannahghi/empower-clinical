package main

import (
	"context"
	"fmt"
	"time"

	"log"

	"github.com/savannahghi/serverutils"
	silotel "github.com/savannahghi/sil-gotel"
	"github.com/savannahghi/empower-clinical/docs"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/common/helpers"
	"github.com/savannahghi/empower-clinical/pkg/clinical/presentation"
)

const (
	JaegerCollectorEndpoint = "JAEGER_COLLECTOR_ENDPOINT"
)

// @title									Clinical Data Repository API
// @version								1.0
// @description							This is the clinical data repository API.
// @termsOfService							http://swagger.io/terms/
// @license.name							Apache 2.0
// @license.url							http://www.apache.org/licenses/LICENSE-2.0.html
// @query.collection.format				multi
// @accept									application/x-www-form-urlencoded
// @securitydefinitions.oauth2.password	OAuth2Password
// @tokenUrl								https://keycloak.example/realms/example/protocol/openid-connect/token
// @scope.read								Grants read access
// @in										header
// @name									Authorization
func main() {
	ctx := context.Background()
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	log.Println("🔄 initializing Open Telemetry")

	otelClient := &silotel.Client{
		OTLPBaseURL: serverutils.MustGetEnvVar("OTLP_ENDPOINT"),
		ServiceName: "clinical-service-backend",
		Environment: serverutils.GetRunningEnvironment(),
		Version:     "1.0.0",
	}

	otelShutdown, err := silotel.NewOtelSDK(ctx, otelClient)
	if err != nil {
		log.Println("❌ could not init Open Telemetry", "error", err)
		log.Fatalf("could not initialize Open Telemetry: %v", err)
	}

	defer func() {
		// context.Background() ensures the shutdown isn't tied to the parent context's lifecycle.
		// The timeout gives the batch processors enough time to flush any buffered spans,
		// metrics, and logs while still ensuring the application doesn't hang indefinitely.
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := otelShutdown(shutdownCtx); err != nil {
			log.Fatalf("OpenTelemetry shutdown error: %v", err)
		}
	}()

	log.Println("✅ Open Telemetry initialized successfully")

	port := serverutils.MustGetEnvVar(serverutils.PortEnvVarName)

	tokenURL := fmt.Sprintf(
		"%s/realms/%s/protocol/openid-connect/token",
		serverutils.MustGetEnvVar("KEYCLOAK_BASE_URL"),
		serverutils.MustGetEnvVar("KEYCLOAK_REALM"),
	)

	docs.SwaggerInfo.SwaggerTemplate = helpers.SetTokenURL(
		docs.SwaggerInfo.SwaggerTemplate,
		// must match the @tokenUrl annotation above
		"https://keycloak.example/realms/example/protocol/openid-connect/token",
		tokenURL,
	)

	presentation.StartServer(ctx, port)
}
