FROM golang:1.15.5 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o application



FROM alpine:latest

WORKDIR /app
VOLUME [ "/app" ]
COPY --from=builder /app/application /application

ENTRYPOINT [ "../application" ]
