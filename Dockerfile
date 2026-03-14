# syntax=docker/dockerfile:1.7
ARG GO_VERSION=1.22
ARG BUF_VERSION=1.64.0
ARG ORAS_VERSION=1.1.0
ARG OAPI_CODEGEN_VERSION=v2.4.0

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
ARG ORAS_VERSION
ARG OAPI_CODEGEN_VERSION

COPY --from=buf /usr/local/bin/buf /usr/local/bin/buf

# Install ORAS
RUN curl -fsSL -o /tmp/oras.tar.gz \
      "https://github.com/oras-project/oras/releases/download/v${ORAS_VERSION}/oras_${ORAS_VERSION}_linux_$(dpkg --print-architecture).tar.gz" && \
    tar -xzf /tmp/oras.tar.gz -C /usr/local/bin oras && \
    rm /tmp/oras.tar.gz

# Install oapi-codegen
RUN go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@${OAPI_CODEGEN_VERSION}

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/root/go/pkg/mod \
    go mod download && go mod verify

COPY buf.gen.yaml buf.yaml ./
COPY scripts/ scripts/

# Generate protobuf stubs
RUN buf generate buf.build/agynio/api \
      --path agynio/api/files/v1 \
      --path agynio/api/llm/v1 \
      --path agynio/api/secrets/v1 \
      --path agynio/api/teams/v1

# Generate OpenAPI code (pull specs + oapi-codegen + embed wrappers)
RUN bash scripts/generate-openapi.sh

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
