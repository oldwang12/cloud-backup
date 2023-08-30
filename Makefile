.PHONY: docker-build-amd64 docker-build-arm64 docker-build-armv7 docker-build-armv8 manifest help

COMMIT_HASH = $(shell git rev-parse --short=7 HEAD)
TAG := $(shell git describe --exact-match --abbrev=0 --tags 2>/dev/null)

ifdef TAG
    IMAGE_TAG := $(TAG)
else
    IMAGE_TAG := $(COMMIT_HASH)
endif

docker-build-amd64: ## 编译 amd64 镜像
	docker build --platform linux/amd64 -t oldwang6/cloud-backup:amd64 -f build/Dockerfile .
	docker push oldwang6/cloud-backup:amd64

docker-build-arm64: ## 编译 arm64 镜像
	docker build --platform linux/arm64 -t oldwang6/cloud-backup:arm64 -f build/Dockerfile .
	docker push oldwang6/cloud-backup:arm64

docker-build-armv7: ## 编译 armv7 镜像
	docker build --platform linux/arm/v7 -t oldwang6/cloud-backup:armv7 -f build/Dockerfile .
	docker push oldwang6/cloud-backup:armv7

docker-build-armv8: ## 编译 armv8 镜像
	docker build --platform linux/arm/v8 -t oldwang6/cloud-backup:armv8 -f build/Dockerfile .
	docker push oldwang6/cloud-backup:armv8

# ================================= 本地测试 =================================
docker-build-amd64-local: ## 编译 amd64 镜像
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GOFLAGS=-mod=vendor go build -o cloud-backup-amd64-local main.go
	docker build --platform linux/amd64 -t oldwang6/cloud-backup:amd64-local -f build/Dockerfile.local.amd64 .
	docker push oldwang6/cloud-backup:amd64-local
	rm -f cloud-backup-amd64-local

docker-build-arm64-local: ## 编译 arm64 镜像
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 GOFLAGS=-mod=vendor go build -o cloud-backup-arm64-local main.go
	docker build --platform linux/arm64 -t oldwang6/cloud-backup:arm64-local -f build/Dockerfile.local.arm64 .
	docker push oldwang6/cloud-backup:arm64-local
	rm -f cloud-backup-arm64-local

manifest: ## 合并镜像
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

help: ## 查看帮助
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf " \033[36m%-20s\033[0m  %s\n", $$1, $$2}' $(MAKEFILE_LIST)