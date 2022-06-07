.PHONY: build test push push-image docs

IMAGE_NAME ?= docsonnet
IMAGE_PREFIX ?= jsonnetlibs
IMAGE_TAG ?= 0.0.3

build:
	docker buildx build -t $(IMAGE_PREFIX)/$(IMAGE_NAME):$(IMAGE_TAG) .

test: build

push: build test push-image

push-image:
	docker push $(IMAGE_PREFIX)/$(IMAGE_NAME):$(IMAGE_TAG)
	docker push $(IMAGE_PREFIX)/$(IMAGE_NAME):latest

docs:
	jsonnet -J doc-util -S -c -m doc-util/ doc-util/docs.jsonnet

