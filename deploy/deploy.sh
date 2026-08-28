#!/usr/bin/env sh

set -eux

# Create the namespace
kubectl create namespace $DEPLOY_ENV_NAMESPACE || true

# Delete Kubernetes secret if exists
kubectl delete secret clinical-service-account --namespace $DEPLOY_ENV_NAMESPACE || true

# Create GCP service account file
cp $GOOGLE_APPLICATION_CREDENTIALS ./service-account.json

# Recreate service account file as Kubernetes secret
kubectl create secret generic clinical-service-account \
    --namespace $DEPLOY_ENV_NAMESPACE \
    --from-file=key.json=./service-account.json

helm repo add chartmuseum $CHARTMUSEUM_DOMAIN --username $CHARTMUSEUM_USERNAME --password $CHARTMUSEUM_PASSWORD

helm pull chartmuseum/clinical --untar

helm upgrade \
    --install \
    --debug \
    --create-namespace \
    --namespace "${DEPLOY_ENV_NAMESPACE}" \
    --set app.replicaCount="${DEPLOY_ENV_APP_REPLICA_COUNT}" \
    --set service.port="${DEPLOY_ENV_PORT}"\
    --set app.container.image="${DEPLOY_ENV_CONTAINER_REGISTRY_PATH}:${CI_COMMIT_SHORT_SHA}"\
    --set app.container.env.googleCloudProject="${DEPLOY_ENV_GOOGLE_CLOUD_PROJECT}"\
    --set app.container.env.environment="${DEPLOY_ENV_ENVIRONMENT}"\
    --set app.container.env.googleProjectNumber="${DEPLOY_ENV_GOOGLE_PROJECT_NUMBER}"\
    --set app.container.env.sentryDSN="${DEPLOY_ENV_SENTRY_DSN}"\
    --set app.container.env.cloudHealthPubsubTopic="${DEPLOY_ENV_CLOUD_HEALTH_PUBSUB_TOPIC}"\
    --set app.container.env.cloudHealthDatasetID="${DEPLOY_ENV_CLOUD_HEALTH_DATASET_ID}"\
    --set app.container.env.cloudHealthDatasetLocation="${DEPLOY_ENV_CLOUD_HEALTH_DATASET_LOCATION}"\
    --set app.container.env.cloudHealthFHIRStoreID="${DEPLOY_ENV_CLOUD_HEALTH_FHIRSTORE_ID}"\
    --set app.container.env.openConceptLabToken="${DEPLOY_ENV_OPENCONCEPTLAB_TOKEN}"\
    --set app.container.env.serviceHost="${DEPLOY_ENV_SERVICE_HOST}"\
    --set app.container.env.openConceptAPIUrl="${DEPLOY_ENV_OPENCONCEPTLAB_API_URL}"\
    --set app.container.env.savannahAdminEmail="${DEPLOY_ENV_SAVANNAH_ADMIN_EMAIL}"\
    --set app.container.env.authserverEndpoint="${DEPLOY_ENV_AUTHSERVER_ENDPOINT}"\
    --set app.container.env.clientID="${DEPLOY_ENV_CLIENT_ID}"\
    --set app.container.env.clientSecret="${DEPLOY_ENV_CLIENT_SECRET}"\
    --set app.container.env.authUsername="${DEPLOY_ENV_AUTH_USERNAME}"\
    --set app.container.env.redisAddress="${DEPLOY_ENV_REDIS_ADDRESS}" \
    --set app.container.env.authPassword="${DEPLOY_ENV_AUTH_PASSWORD}"\
    --set app.container.env.grantType="${DEPLOY_ENV_GRANT_TYPE}"\
    --set app.container.env.clinicalBucketName="${DEPLOY_ENV_CLINICAL_BUCKET_NAME}"\
    --set app.container.env.urlShortenerDomain="${DEPLOY_ENV_URL_SHORTENER_DOMAIN}"\
    --set app.container.env.urlShortenerApiKey="${DEPLOY_ENV_URL_SHORTENER_API_KEY}"\
    --set app.container.env.awsSecretAccessKey="${DEPLOY_ENV_AWS_SECRET_ACCESS_KEY}"\
    --set app.container.env.awsAccessKeyID="${DEPLOY_ENV_AWS_ACCESS_KEY_ID}"\
    --set app.container.env.jaegerCollectorEndpoint="${DEPLOY_ENV_JAEGER_COLLECTOR_ENDPOINT}"\
    --set app.container.env.sesawsRegion="${DEPLOY_ENV_SES_AWS_REGION}"\
    --set app.container.env.defaultFromEmail="${DEPLOY_ENV_DEFAULT_FROM_EMAIL}"\
    --set app.container.env.hapiFHIRBaseURL="${DEPLOY_ENV_HAPI_FHIR_BASE_URL}"\
    --set networking.issuer.name="letsencrypt-prod"\
    --set app.container.env.jwtKey="${DEPLOY_ENV_JWT_KEY}"\
    --set networking.issuer.privateKeySecretRef="letsencrypt-prod"\
    --set networking.ingress.host="${DEPLOY_ENV_APP_DOMAIN}"\
    --set networking.useHttps="${DEPLOY_ENV_APP_USE_HTTPS:-true}"\
    --set app.container.env.defaultSentryTraceSampleRate="${DEPLOY_ENV_SENTRY_TRACE_SAMPLE_RATE}"\
    --set app.container.env.advantageBaseURL="${DEPLOY_ENV_ADVANTAGE_BASE_URL}"\
    --version="${HELM_CHART_VERSION}" \
    --wait \
    --timeout 300s \
    -f ./clinical/values.yaml \
    $DEPLOY_ENV_APP_NAME \
    chartmuseum/clinical