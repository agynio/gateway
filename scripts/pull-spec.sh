#!/usr/bin/env bash
set -euo pipefail

REGISTRY="${OPENAPI_REGISTRY:-ghcr.io}"
OUTPUT_DIR=".openapi"

if ! command -v oras &>/dev/null; then
  echo "ERROR: oras CLI not found" >&2; exit 1
fi

mkdir -p "${OUTPUT_DIR}"

pull_spec() {
  local image="$1" filename="$2"
  local output_file="${OUTPUT_DIR}/${filename}" dist_file="${OUTPUT_DIR}/dist/${filename}"
  echo "Pulling ${REGISTRY}/${image} -> ${output_file}"
  oras pull "${REGISTRY}/${image}" --output "${OUTPUT_DIR}"
  [[ -f "${dist_file}" ]] && cp "${dist_file}" "${output_file}"
  if [[ ! -f "${output_file}" ]]; then
    echo "ERROR: expected ${output_file} not found" >&2; ls -la "${OUTPUT_DIR}" >&2; exit 1
  fi
  echo "Spec downloaded: ${output_file} ($(wc -c < "${output_file}") bytes)"
}

pull_spec "${TEAM_OPENAPI_IMAGE:-agynio/openapi/team:1}" "team-v1.yaml"
pull_spec "${LLM_OPENAPI_IMAGE:-agynio/openapi/llm:1}" "llm-v1.yaml"
