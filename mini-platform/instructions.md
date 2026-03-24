# Initialization

```bash
go get google.golang.org/grpc
go get google.golang.org/protobuf
brew install protobuf
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

```bash
protoc \
  --go_out=. \
  --go_opt=paths=source_relative \
  --go-grpc_out=. \
  --go-grpc_opt=paths=source_relative \
  proto/store.proto
```