FROM golang:1.20 as builder
WORKDIR /root/
COPY . .
RUN GOOS=linux GOARCH=arm CGO_ENABLED=0 GOFLAGS=-mod=vendor go build -o cloud-backup main.go

# =================================== 分层编译 ==============================================
FROM alpine AS final
WORKDIR /root/

# 国内使用的goproxy
ENV GOPROXY=https://goproxy.cn

# 设置时区
ENV TZ=Asia/Shanghai

COPY tar.sh tar.sh
# COPY config.yaml configs.yaml
COPY --from=builder /root/cloud-backup .

RUN apk add --update tzdata \
    && apk add --no-cache curl \
    && cp /usr/share/zoneinfo/$TZ /etc/localtime \
    && echo $TZ > /etc/timezone \
    && rm -rf /var/cache/apk/* \
    && chmod +x tar.sh \
    && chmod +x cloud-backup

ENTRYPOINT ["/root/cloud-backup"]