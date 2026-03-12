#!/usr/bin/env bash
set -euo pipefail

REGISTRY="${OPENAPI_REGISTRY:-ghcr.io}"
IMAGE="${OPENAPI_IMAGE:-agynio/openapi/team:1}"
OUTPUT_DIR=".openapi"
OUTPUT_FILE="${OUTPUT_DIR}/team-v1.yaml"

if ! command -v oras &>/dev/null; then
  echo "ERROR: oras CLI not found. Install from https://oras.land" >&2
  exit 1
fi

mkdir -p "${OUTPUT_DIR}"

echo "Pulling ${REGISTRY}/${IMAGE} -> ${OUTPUT_FILE}"
oras pull "${REGISTRY}/${IMAGE}" \
  --output "${OUTPUT_DIR}"

if [[ ! -f "${OUTPUT_FILE}" ]]; then
  if [[ -f "${OUTPUT_DIR}/dist/team-v1.yaml" ]]; then
    cp "${OUTPUT_DIR}/dist/team-v1.yaml" "${OUTPUT_FILE}"
  else
    echo "ERROR: expected ${OUTPUT_FILE} not found after pull" >&2
    ls -la "${OUTPUT_DIR}" >&2
    exit 1
  fi
fi

echo "Spec downloaded: ${OUTPUT_FILE} ($(wc -c < "${OUTPUT_FILE}") bytes)"
