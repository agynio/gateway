# syntax=docker/dockerfile:1.7

ARG GO_VERSION=1.22

FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-bookworm AS build

WORKDIR /src

ENV CGO_ENABLED=0 \
    GO111MODULE=on

COPY go.mod go.sum ./

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/root/go/pkg/mod \
    go mod download && \
    go mod verify

COPY . .

ARG TARGETOS
ARG TARGETARCH

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/root/go/pkg/mod \
    GOOS=$TARGETOS GOARCH=$TARGETARCH go build \
      -trimpath \
      -ldflags="-s -w" \
      -o /out/gateway ./cmd/gateway

FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app

LABEL org.opencontainers.image.source="https://github.com/agynio/gateway"

COPY --from=build /out/gateway ./gateway
COPY --from=build /src/spec ./spec

EXPOSE 8080

ENTRYPOINT ["/app/gateway"]
