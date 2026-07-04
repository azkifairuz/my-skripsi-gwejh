FROM golang:1.26-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/whatsapp-service ./cmd/whatsapp-service

FROM alpine:3.22

WORKDIR /app
COPY --from=builder /out/whatsapp-service /app/whatsapp-service

ENTRYPOINT ["/app/whatsapp-service"]
