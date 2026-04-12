#!/usr/bin/env bash
set -eu

echo "=== DevSpace startup ==="

echo "Generating protobuf types..."
buf generate https://github.com/agynio/api.git#commit=ec008b1e2dfacec3e4d85776729fe1c3d5f2c42d \
  --include-imports \
  --path proto/agynio/api/agents/v1 \
  --path proto/agynio/api/threads/v1 \
  --path proto/agynio/api/chat/v1 \
  --path proto/agynio/api/notifications/v1 \
  --path proto/agynio/api/files/v1 \
  --path proto/agynio/api/agent_state/v1 \
  --path proto/agynio/api/token_counting/v1 \
  --path proto/agynio/api/llm/v1 \
  --path proto/agynio/api/secrets/v1 \
  --path proto/agynio/api/identity/v1 \
  --path proto/agynio/api/tracing/v1 \
  --path proto/agynio/api/users/v1 \
  --path proto/agynio/api/organizations/v1 \
  --path proto/agynio/api/runners/v1 \
  --path proto/agynio/api/ziti_management/v1 \
  --path proto/agynio/api/gateway/v1

echo "Downloading Go modules..."
go mod download

echo "Starting dev server (air)..."
exec air
