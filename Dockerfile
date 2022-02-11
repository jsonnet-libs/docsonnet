FROM --platform=$BUILDPLATFORM golang:1.17.6 as base

ENV GO111MODULE=on
WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

FROM base AS builder

ENV GOARCH=$TARGETARCH
RUN CGO_ENABLED=0 go build -ldflags='-s -w -extldflags "-static"' .

FROM alpine:3.12
COPY --from=builder /app/docsonnet /usr/local/bin

ENTRYPOINT ["docsonnet"]
