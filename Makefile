.PHONY: build-amd64 build-arm64 build-arm-v7 help

COMMIT_HASH = $(shell git rev-parse --short=7 HEAD)

# build: ## build image
# 	docker buildx create --name all --node local --driver docker-container --platform linux/amd64,linux/arm64,linux/arm/v7 --use
# 	docker buildx use all
# 	docker buildx build --platform linux/amd64,linux/arm64,linux/arm/v7 -t oldwang6/cloud-backup:${COMMIT_HASH} -o type=registry .
# 	docker buildx build --platform linux/amd64,linux/arm64,linux/arm/v7 -t oldwang6/cloud-backup:latest -o type=registry .

build-amd64: ## build linux/amd64 image
	docker buildx create --name test --node local --driver docker-container --platform linux/amd64 --use
	docker buildx use test
	docker buildx build --platform linux/amd64 -t oldwang6/cloud-backup:latest -o type=registry .

build-arm64: ## build linux/arm64 image
	docker buildx create --name test --node local --driver docker-container --platform linux/arm64 --use
	docker buildx use test
	docker buildx build --platform linux/arm64 -t oldwang6/cloud-backup:latest -o type=registry .

build-arm-v7: ## build linux/arm/v7 image
	docker buildx create --name test --node local --driver docker-container --platform linux/arm/v7 --use
	docker buildx use test
	docker buildx build --platform linux/arm/v7 -t oldwang6/cloud-backup:latest -o type=registry .

help: ## help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf " \033[36m%-20s\033[0m  %s\n", $$1, $$2}' $(MAKEFILE_LIST)