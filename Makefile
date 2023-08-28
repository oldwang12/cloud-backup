.PHONY: docker-build-amd64 docker-build-arm64 docker-build-armv7 docker-build-armv8 manifest help

COMMIT_HASH = $(shell git rev-parse --short=7 HEAD)
TAG := $(shell git describe --exact-match --abbrev=0 --tags 2>/dev/null)

ifdef TAG
    IMAGE_TAG := $(TAG)
else
    IMAGE_TAG := $(COMMIT_HASH)
endif

docker-build-amd64:
	docker build --platform linux/amd64 -t oldwang6/cloud-backup:amd64 .
	docker push oldwang6/cloud-backup:amd64

docker-build-arm64:
	docker build --platform linux/arm64 -t oldwang6/cloud-backup:arm64 .
	docker push oldwang6/cloud-backup:arm64

docker-build-armv7:
	docker build --platform linux/arm/v7 -t oldwang6/cloud-backup:armv7 .
	docker push oldwang6/cloud-backup:armv7

docker-build-armv8:
	docker build --platform linux/arm/v8 -t oldwang6/cloud-backup:armv8 .
	docker push oldwang6/cloud-backup:armv8

manifest:
	docker manifest create oldwang6/cloud-backup:${IMAGE_TAG} \
           oldwang6/cloud-backup:amd64 \
		   oldwang6/cloud-backup:arm64 \
           oldwang6/cloud-backup:armv7 \
		   oldwang6/cloud-backup:armv8

	docker manifest create oldwang6/cloud-backup:latest \
           oldwang6/cloud-backup:amd64 \
           oldwang6/cloud-backup:arm64 \
           oldwang6/cloud-backup:armv7 \
           oldwang6/cloud-backup:armv8

	docker manifest push oldwang6/cloud-backup:${IMAGE_TAG}
	docker manifest push oldwang6/cloud-backup:latest

help: ## help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf " \033[36m%-20s\033[0m  %s\n", $$1, $$2}' $(MAKEFILE_LIST)