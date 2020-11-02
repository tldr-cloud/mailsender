#!/bin/zsh

CONFIG_ID="tldr-subscription-config-$(date +%s)"
API_ID="tldr-mailing-apis"
API_YAML_FILE="subscription-api.yaml"
PROJECT_ID="mailsender-288100"
API_SERVICE_ACCOUNT="mailsender@mailsender-288100.iam.gserviceaccount.com"
GATEWAY_ID="tldr-subscriptions-gateway"
GCP_REGION="us-central1"

gcloud beta api-gateway api-configs create "${CONFIG_ID}" \
  --api=${API_ID} --openapi-spec="${API_YAML_FILE}" \
  --project="${PROJECT_ID}" --backend-auth-service-account="${API_SERVICE_ACCOUNT}"

gcloud beta api-gateway gateways update "${GATEWAY_ID}" \
  --api="${API_ID}" --api-config="${CONFIG_ID}" \
  --location="${GCP_REGION}" --project="${PROJECT_ID}"
