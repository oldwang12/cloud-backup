FROM golang:1.20 as builder
WORKDIR /root/
COPY . .
RUN go build -o cloud-backup main.go

# =================================================================================
FROM alpine AS final
WORKDIR /root/

ENV TZ=Asia/Shanghai

COPY --from=builder /root/cloud-backup .

RUN apk add --update tzdata \
    && apk add --no-cache curl \
    && cp /usr/share/zoneinfo/$TZ /etc/localtime \
    && echo $TZ > /etc/timezone \
    && rm -rf /var/cache/apk/* \
    && chmod +x cloud-backup

ENTRYPOINT ["/root/cloud-backup"]