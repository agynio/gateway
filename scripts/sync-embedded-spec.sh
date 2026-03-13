#!/usr/bin/env bash
set -euo pipefail

sync_spec() {
  local source="$1" target="$2"
  if [[ ! -f "${source}" ]]; then
    echo "ERROR: ${source} not found. Run scripts/pull-spec.sh first." >&2; exit 1
  fi
  cp "${source}" "${target}"; echo "Synced ${source} -> ${target}"
}

sync_spec ".openapi/team-v1.yaml" "internal/apischema/teamv1/team-v1.yaml"
sync_spec ".openapi/llm-v1.yaml" "internal/apischema/llmv1/llm-v1.yaml"
