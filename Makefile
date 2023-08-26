.PHONY: docker-build-amd64 docker-build-arm64 docker-build-arm-v7 docker-push manifest help

COMMIT_HASH = $(shell git rev-parse --short=7 HEAD)

docker-build-amd64:
	docker build -t oldwang6/cloud-backup:amd64 .

docker-build-arm64:
	docker build -t oldwang6/cloud-backup:arm64 .

docker-build-arm-v7:
	docker build -t oldwang6/cloud-backup:arm-v7 .

docker-push:
	docker push oldwang6/cloud-backup:amd64 
	docker push oldwang6/cloud-backup:arm64
	docker push oldwang6/cloud-backup:arm-v7

manifest:
	docker manifest create oldwang6/cloud-backup:${COMMIT_HASH} \
           oldwang6/cloud-backup:amd64 \
           oldwang6/cloud-backup:arm64 \
           oldwang6/cloud-backup:arm-v7
	docker manifest push oldwang6/cloud-backup:${COMMIT_HASH}


# docker buildx create --name all --node local --driver docker-container --platform linux/amd64,linux/arm64,linux/arm/v7 --use
# docker buildx use all
# docker buildx build --platform linux/amd64 -t oldwang6/cloud-backup-amd64:${COMMIT_HASH} --push .
# docker buildx build --platform linux/arm/v7 -t oldwang6/cloud-backup-armv7:${COMMIT_HASH} --push .
# docker buildx build --platform linux/amd64,linux/arm/v7 -t oldwang6/cloud-backup:latest --push .

# docker buildx build --platform linux/arm64 -t oldwang6/cloud-backup-arm64:${COMMIT_HASH} --push .

# build-amd64: ## build linux/amd64 image
# 	docker buildx create --name test --node local --driver docker-container --platform linux/amd64 --use
# 	docker buildx use test
# 	docker buildx build --platform linux/amd64 -t oldwang6/cloud-backup:latest -o type=registry .

# build-arm64: ## build linux/arm64 image
# 	docker buildx create --name test --node local --driver docker-container --platform linux/arm64 --use
# 	docker buildx use test
# 	docker buildx build --platform linux/arm64 -t oldwang6/cloud-backup:latest -o type=registry .

# build-arm-v7: ## build linux/arm/v7 image
# 	docker buildx create --name test --node local --driver docker-container --platform linux/arm/v7 --use
# 	docker buildx use test
# 	docker buildx build --platform linux/arm/v7 -t oldwang6/cloud-backup:latest -o type=registry .

help: ## help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf " \033[36m%-20s\033[0m  %s\n", $$1, $$2}' $(MAKEFILE_LIST)