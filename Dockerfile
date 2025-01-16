FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o app ./cmd/main.go

FROM alpine:latest
WORKDIR /root/
RUN apk add --no-cache postgresql-client
COPY --from=builder /app/app .
COPY --from=builder /app/config.hcl .
CMD ["./app"]