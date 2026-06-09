FROM golang:1.24-alpine AS builder
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o weatherserver ./cmd/weatherserver

FROM alpine:latest
RUN apk add --no-cache tzdata
COPY --from=builder /build/weatherserver /usr/local/bin/weatherserver
USER 65532:65532
EXPOSE 8081
CMD ["weatherserver"]
