# go-concepts

A hands-on reference repo for Go patterns used in production API and platform engineering. Built while preparing for backend engineering interviews focused on distributed systems, gRPC, and resilient service design.

---

## Contents

### Pattern examples

| File | What it covers |
|---|---|
| `goroutine.go` | Fan-out concurrency, worker pools, buffered channels, goroutine leak prevention |
| `defer.go` | LIFO execution order, named return values, commit/rollback pattern, defer in loops |
| `errors.go` | Sentinel errors, typed errors, wrapping with `%w`, `errors.Is` vs `errors.As` |
| `midle.go` | HTTP middleware chaining, auth, logging, rate limiting with `atomic.Int64` |
| `circuit_breaker.go` | Circuit breaker state machine — closed, open, half-open — with mutex-safe state transitions |

### mini-platform

A minimal two-service platform demonstrating HTTP → gRPC communication.

```
mini-platform/
├── proto/
│   ├── store.proto          # service contract definition
│   ├── store.pb.go          # generated message structs
│   └── store_grpc.pb.go     # generated gRPC client and server interfaces
├── publish-api/
│   └── main.go              # HTTP REST service (port 8080)
└── store-service/
    └── main.go              # gRPC server (port 50051)
```

**publish-api** — accepts `POST /publish` with a JSON body and forwards it to store-service over gRPC. Demonstrates connection reuse, request context propagation, and handler injection via struct.

**store-service** — gRPC server that receives data and stores it in an in-memory map. Demonstrates mutex-safe concurrent writes, proto-generated interface implementation, and `UnimplementedStoreServiceServer` embedding for forward compatibility.

---

## Running the platform

Start store-service first, then publish-api in a separate terminal:

```bash
# terminal 1
go run mini-platform/store-service/main.go

# terminal 2
go run mini-platform/publish-api/main.go
```

Send a request:

```bash
curl -X POST http://localhost:8080/publish \
  -H "Content-Type: application/json" \
  -d '{"id": "1", "data": "hello"}'
```

---

## Regenerating proto code

If you modify `store.proto`, regenerate the Go code from the repo root:

```bash
protoc \
  --go_out=. \
  --go_opt=paths=source_relative \
  --go-grpc_out=. \
  --go-grpc_opt=paths=source_relative \
  mini-platform/proto/store.proto
```

Requires `protoc`, `protoc-gen-go`, and `protoc-gen-go-grpc` to be installed and on your PATH.