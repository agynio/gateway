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
  --path agynio/api/users/v1 \
  --path agynio/api/ziti_management/v1 \
  --path agynio/api/gateway/v1/agent_state.proto \
  --path agynio/api/gateway/v1/agents.proto \
  --path agynio/api/gateway/v1/chat.proto \
  --path agynio/api/gateway/v1/files.proto \
  --path agynio/api/gateway/v1/llm.proto \
  --path agynio/api/gateway/v1/notifications.proto \
  --path agynio/api/gateway/v1/secrets.proto \
  --path agynio/api/gateway/v1/threads.proto \
  --path agynio/api/gateway/v1/token_counting.proto \
  --path agynio/api/gateway/v1/users.proto

echo "Downloading Go modules..."
go mod download

echo "Starting dev server (air)..."
exec air
