#!/usr/bin/env bash
set -euo pipefail

SOURCE=".openapi/team-v1.yaml"
TARGET="internal/apischema/teamv1/team-v1.yaml"

if [[ ! -f "${SOURCE}" ]]; then
  echo "ERROR: ${SOURCE} not found. Run scripts/pull-spec.sh first." >&2
  exit 1
fi

cp "${SOURCE}" "${TARGET}"
echo "Synced ${SOURCE} -> ${TARGET}"
