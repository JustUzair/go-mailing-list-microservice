# Setup

This project requires a `gcc` compiler installed and the `protobuf` code generation tools.
## Install missing dependencies
```bash
    go mod tidy
```
## Start the server 

```bash
    go run ./server
```

## Start the client
```bash
    go run ./client
```

## Install protobuf compiler

Install the `protoc` tool using the instructions available at [https://grpc.io/docs/protoc-installation/](https://grpc.io/docs/protoc-installation/).

Alternatively you can download a pre-built binary from [https://github.com/protocolbuffers/protobuf/releases](https://github.com/protocolbuffers/protobuf/releases) and placing the extracted binary somewhere in your `$PATH`.

## Install Go protobuf codegen tools

`go install google.golang.org/protobuf/cmd/protoc-gen-go@latest`

`go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest`

## Generate Go code from .proto files

```bash
protoc --go_out=. --go_opt=paths=source_relative \ --go-grpc_out=. --go-grpc_opt=paths=source_relative \ Proto/mail.proto
```

### NOTE : If above command generates an error try the following command

```bash
protoc --go_out=. --go_opt=paths=source_relative ./proto/mail.proto --go-grpc_out=. --go-grpc_opt=paths=source_relative ./proto/mail.proto
```


