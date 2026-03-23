#!/usr/bin/env bash
set -eu

echo "=== DevSpace startup ==="

echo "Generating protobuf types..."
buf generate buf.build/agynio/api \
  --path agynio/api/agents/v1 \
  --path agynio/api/threads/v1 \
  --path agynio/api/chat/v1 \
  --path agynio/api/notifications/v1 \
  --path agynio/api/files/v1 \
  --path agynio/api/agent_state/v1 \
  --path agynio/api/token_counting/v1 \
  --path agynio/api/llm/v1 \
  --path agynio/api/secrets/v1 \
  --path agynio/api/tracing/v1 \
  --path agynio/api/users/v1 \
  --path agynio/api/ziti_management/v1 \
  --path agynio/api/gateway/v1

echo "Downloading Go modules..."
go mod download

echo "Starting dev server (air)..."
exec air
