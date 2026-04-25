#!/usr/bin/env bash
set -eu

echo "=== DevSpace startup ==="

echo "Generating protobuf types..."
buf generate buf.build/agynio/api

echo "Downloading Go modules..."
go mod download

echo "Starting dev server (air)..."
exec air -c /opt/app/data/scripts/air.toml
