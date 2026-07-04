.PHONY: proto proto-ocr proto-whatsapp proto-finance proto-nlp proto-whatsapp-clients

PROTOC ?= protoc
PROTOC_FLAGS = -I proto

proto: proto-ocr proto-whatsapp proto-finance proto-nlp proto-whatsapp-clients

proto-ocr:
	mkdir -p ocr-service/internal/generated/proto
	$(PROTOC) $(PROTOC_FLAGS) \
		--go_out=ocr-service/internal/generated/proto \
		--go_opt=paths=source_relative \
		--go-grpc_out=ocr-service/internal/generated/proto \
		--go-grpc_opt=paths=source_relative \
		ocr/v1/ocr.proto

proto-whatsapp:
	mkdir -p whatsapp-service/internal/generated/proto
	$(PROTOC) $(PROTOC_FLAGS) \
		--go_out=whatsapp-service/internal/generated/proto \
		--go_opt=paths=source_relative \
		--go-grpc_out=whatsapp-service/internal/generated/proto \
		--go-grpc_opt=paths=source_relative \
		whatsapp/v1/whatsapp.proto

proto-finance:
	mkdir -p finance-service/internal/generated/proto
	$(PROTOC) $(PROTOC_FLAGS) \
		--go_out=finance-service/internal/generated/proto \
		--go_opt=paths=source_relative \
		--go-grpc_out=finance-service/internal/generated/proto \
		--go-grpc_opt=paths=source_relative \
		finance/v1/finance.proto

proto-nlp:
	mkdir -p nlp-service/internal/generated/proto
	$(PROTOC) $(PROTOC_FLAGS) \
		--go_out=nlp-service/internal/generated/proto \
		--go_opt=paths=source_relative \
		--go-grpc_out=nlp-service/internal/generated/proto \
		--go-grpc_opt=paths=source_relative \
		nlp/v1/nlp.proto

proto-whatsapp-clients:
	mkdir -p whatsapp-service/internal/generated/proto
	$(PROTOC) $(PROTOC_FLAGS) \
		--go_out=whatsapp-service/internal/generated/proto \
		--go_opt=paths=source_relative \
		--go-grpc_out=whatsapp-service/internal/generated/proto \
		--go-grpc_opt=paths=source_relative \
		finance/v1/finance.proto \
		nlp/v1/nlp.proto \
		ocr/v1/ocr.proto
