.PHONY: build help

COMMIT_HASH = $(shell git rev-parse --short=7 HEAD)

build: ## build image
	docker buildx build --platform linux/amd64,linux/arm64,linux/arm/v7 -t registry.cn-hangzhou.aliyuncs.com/breawang/cloud-backup:${COMMIT_HASH} -o type=registry .
	docker buildx build --platform linux/amd64,linux/arm64,linux/arm/v7 -t oldwang6/cloud-backup:${COMMIT_HASH} -o type=registry .

help: ## help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf " \033[36m%-20s\033[0m  %s\n", $$1, $$2}' $(MAKEFILE_LIST)