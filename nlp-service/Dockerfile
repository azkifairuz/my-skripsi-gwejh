FROM golang:1.26-alpine AS builder

WORKDIR /src

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/nlp-service ./cmd/nlp-service

FROM alpine:3.22

WORKDIR /app
COPY --from=builder /out/nlp-service /app/nlp-service

ENTRYPOINT ["/app/nlp-service"]
