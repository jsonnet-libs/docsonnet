.PHONY: docsonnet cparse

VERSION := $(shell git describe --tags --dirty --always)

docsonnet: pkged.go
	go build -o docsonnet -ldflags "-X main.Version=dev-$(VERSION)" .

cparse: pkged.go
	go build -o cparse -ldflags "-X main.Version=dev-$(VERSION)" ./cmd/cparse

pkged.go: doc-util/main.libsonnet pkg/comments/dsl.libsonnet
	rm -f cmd/cparse/pkged.go
	pkger
	ln -s $(PWD)/pkged.go cmd/cparse/pkged.go
