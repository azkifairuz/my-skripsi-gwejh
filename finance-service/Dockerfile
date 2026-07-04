FROM golang:1.26-alpine AS builder

WORKDIR /src

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/finance-service ./cmd/finance-service

FROM alpine:3.22

WORKDIR /app
COPY --from=builder /out/finance-service /app/finance-service

ENTRYPOINT ["/app/finance-service"]
