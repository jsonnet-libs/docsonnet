FROM golang:1.16.4 as base

ENV GO111MODULE=on
WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

FROM base AS builder
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build .

FROM alpine:3.12
COPY --from=builder /app/docsonnet /usr/local/bin

ENTRYPOINT ["docsonnet"]
