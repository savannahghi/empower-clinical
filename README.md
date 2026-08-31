# clinical

![Linting and Tests](https://github.com/savannahghi/empower-clinical/actions/workflows/ci.yml/badge.svg)
[![Coverage Status](https://coveralls.io/repos/github/savannahghi/clinical/badge.svg?branch=develop-v2)](https://coveralls.io/github/savannahghi/clinical?branch=develop-v2)

APIs to bridge to a FHIR clinical repository

## Environment variables

For local development, you need to *export* the following env vars:

```bash
# Google Cloud Settings
export GOOGLE_APPLICATION_CREDENTIALS="<a path to a Google service account JSON file>"
export GOOGLE_CLOUD_PROJECT="<the name of the project that the service account above belongs to>"
export FIREBASE_WEB_API_KEY="<an API key from the Firebase console for the project mentioned above>"
```

The server deploys to Google Cloud Run. For Cloud Run, the necessary environment
variables are:

- `GOOGLE_CLOUD_PROJECT`
- `FIREBASE_WEB_API_KEY`

## Deployment
This application is deployed via Google Cloud Build ( https://cloud.google.com/build ) to Google Cloud Run ( https://cloud.google.com/run ). There's a cloudbuild.yaml file in the home folder. Secrets (e.g production settings) are managed with Google Secret Manager ( https://cloud.google.com/secret-manager ).
