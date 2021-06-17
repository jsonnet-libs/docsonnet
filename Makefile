.PHONY: build test push push-image

IMAGE_NAME ?= docsonnet
IMAGE_PREFIX ?= jsonnetlibs
IMAGE_TAG ?= 0.0.2

build:
	docker buildx build -t $(IMAGE_PREFIX)/$(IMAGE_NAME):$(IMAGE_TAG) .

test: build

push: build test push-image

push-image:
	docker push $(IMAGE_PREFIX)/$(IMAGE_NAME):$(IMAGE_TAG)
	docker push $(IMAGE_PREFIX)/$(IMAGE_NAME):latest
