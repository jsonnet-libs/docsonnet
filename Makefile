.PHONY: build release docs

build:
	goreleaser build --rm-dist --snapshot

release:
	goreleaser release --rm-dist

docs:
	jsonnet -S -c -m doc-util/ \
		-e "(import 'doc-util/main.libsonnet').render(import 'doc-util/main.libsonnet')"

