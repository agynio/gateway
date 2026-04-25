# Gateway

Schema-first HTTP gateway built on Go 1.24.10 with ConnectRPC services generated from agynio/api
protobuf definitions.

Architecture: [Gateway](https://github.com/agynio/architecture/blob/main/architecture/gateway.md)

## Local Development

Full setup: [Local Development](https://github.com/agynio/architecture/blob/main/architecture/operations/local-development.md)

### Prepare environment

```bash
git clone https://github.com/agynio/bootstrap.git
cd bootstrap
chmod +x apply.sh
./apply.sh -y
```

See [bootstrap](https://github.com/agynio/bootstrap) for details.

### Run from sources

```bash
# Deploy once (exit when healthy)
devspace dev

# Watch mode (streams logs, re-syncs on changes)
devspace dev -w
```

## Adding a New API Domain

Every API domain in the gateway must be defined in protobuf and exposed via
ConnectRPC. The standard flow mirrors the existing gateway handlers:

1. Add or update the protobuf definition in `agynio/api` and include the path
   in `buf.gen.yaml` and CI `buf generate` commands.
2. Regenerate stubs with `buf generate` and implement forwarding handlers in
   `internal/gateway/<domain>.go` that satisfy the generated Connect handler
   interfaces.
3. Wire the new handler in `cmd/gateway/main.go` using the Connect mux and
   interceptors.
