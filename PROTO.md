# Protobuf Generation

This project keeps shared service contracts in the root `proto/` directory.

Required tools:

- `protoc`
- `protoc-gen-go`
- `protoc-gen-go-grpc`

Generate all Go protobuf files from the project root:

```bash
make proto
```

Generated files are written into each service:

- `ocr-service/internal/generated/proto/ocr/v1/`
- `whatsapp-service/internal/generated/proto/whatsapp/v1/`
- `finance-service/internal/generated/proto/finance/v1/`
- `nlp-service/internal/generated/proto/nlp/v1/`

After regenerating, validate each service:

```bash
go build ./...
```
