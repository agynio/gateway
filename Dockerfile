# syntax=docker/dockerfile:1.7
ARG GO_VERSION=1.25
ARG BUF_VERSION=1.66.1

# Stage 1: Download buf binary
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-bookworm AS buf
ARG BUF_VERSION
RUN curl -sSL \
      "https://github.com/bufbuild/buf/releases/download/v${BUF_VERSION}/buf-$(uname -s)-$(uname -m)" \
      -o /usr/local/bin/buf && \
    chmod +x /usr/local/bin/buf

# Stage 2: Generate + compile
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-bookworm AS build
WORKDIR /src
ENV CGO_ENABLED=0 GO111MODULE=on

COPY --from=buf /usr/local/bin/buf /usr/local/bin/buf

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/root/go/pkg/mod \
    go mod download && go mod verify

COPY buf.gen.yaml buf.yaml ./

# Generate protobuf stubs
RUN buf generate https://github.com/agynio/api.git#commit=ec008b1e2dfacec3e4d85776729fe1c3d5f2c42d \
      --include-imports \
      --path proto/agynio/api/agents/v1 \
      --path proto/agynio/api/apps/v1 \
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

COPY . .

ARG TARGETOS TARGETARCH
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/root/go/pkg/mod \
    GOOS=$TARGETOS GOARCH=$TARGETARCH go build \
      -trimpath -ldflags="-s -w" -o /out/gateway ./cmd/gateway

# Stage 3: Runtime
FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app
LABEL org.opencontainers.image.source="https://github.com/agynio/gateway"
COPY --from=build /out/gateway ./gateway
EXPOSE 8080
ENTRYPOINT ["/app/gateway"]
