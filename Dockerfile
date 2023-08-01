FROM golang:1.20 as builder
WORKDIR /root/
COPY . .
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GOFLAGS=-mod=vendor go build -o sync main.go

# =================================== 分层编译 ==============================================
FROM alpine:3.9 AS final

# 国内使用的goproxy
ENV GOPROXY=https://goproxy.cn

# 设置时区
ENV TZ=Asia/Shanghai
RUN apk add --update tzdata \
    && cp /usr/share/zoneinfo/$TZ /etc/localtime \
    && echo $TZ > /etc/timezone \
    && rm -rf /var/cache/apk/*

WORKDIR /root/
COPY --from=builder /root/sync .
RUN chmod +x sync
ENTRYPOINT ["/root/sync"]