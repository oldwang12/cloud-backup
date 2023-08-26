.PHONY: docker-build-amd64 docker-build-arm-v7 manifest help

COMMIT_HASH = $(shell git rev-parse --short=7 HEAD)

docker-build-amd64:
	docker build --platform linux/amd64 -t oldwang6/cloud-backup:amd64 .
	docker push oldwang6/cloud-backup:amd64

docker-build-armv7:
	docker build --platform linux/arm/v7 -t oldwang6/cloud-backup:armv7 .
	docker push oldwang6/cloud-backup:armv7

manifest:
	docker manifest create oldwang6/cloud-backup:${COMMIT_HASH} \
           oldwang6/cloud-backup:amd64 \
           oldwang6/cloud-backup:armv7

	docker manifest create oldwang6/cloud-backup:latest \
           oldwang6/cloud-backup:amd64 \
           oldwang6/cloud-backup:armv7

	docker manifest push oldwang6/cloud-backup:${COMMIT_HASH}
	docker manifest push oldwang6/cloud-backup:latest

test:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GOFLAGS=-mod=vendor go build -o cloud-backup-amd64 main.go
	GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=0 GOFLAGS=-mod=vendor go build -o cloud-backup-armv7 main.go
	docker build --platform linux/amd64 -t oldwang6/cloud-backup:amd64 -f Dockerfile.amd64 .
	docker build --platform linux/arm/v7 -t oldwang6/cloud-backup:armv7 -f Dockerfile.armv7 .
	docker push oldwang6/cloud-backup:amd64
	docker push oldwang6/cloud-backup:armv7

	docker manifest create oldwang6/cloud-backup:test \
           oldwang6/cloud-backup:amd64 \
           oldwang6/cloud-backup:armv7
	docker manifest push oldwang6/cloud-backup:test

help: ## help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf " \033[36m%-20s\033[0m  %s\n", $$1, $$2}' $(MAKEFILE_LIST)