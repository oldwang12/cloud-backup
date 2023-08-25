FROM golang:1.20 as builder
WORKDIR /root/
COPY . .
RUN CGO_ENABLED=0 GOFLAGS=-mod=vendor go build -o cloud-backup main.go

# =================================================================================
FROM alpine
WORKDIR /root/

ENV TZ=Asia/Shanghai

COPY --from=builder /root/cloud-backup .
RUN chmod +x /root/cloud-backup \
    apk add --update tzdata \
    && cp /usr/share/zoneinfo/$TZ /etc/localtime \
    && echo $TZ > /etc/timezone \
    && rm -rf /var/cache/apk/*

CMD ["/root/cloud-backup"]