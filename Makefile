.PHONY: build release docs

build:
	goreleaser build --rm-dist --snapshot

release:
	goreleaser release --rm-dist

docs:
	cd doc-util && \
	jsonnet -S -c -m . \
		-e "(import './main.libsonnet').render(import './main.libsonnet')"

