# FROM golang:1.20 as builder
# WORKDIR /root/
# COPY . .
# RUN CGO_ENABLED=0 GOFLAGS=-mod=vendor go build -o cloud-backup main.go

# =================================================================================
# FROM oldwang6/alpine
# WORKDIR /root/
# COPY --from=builder /root/cloud-backup .
# RUN chmod +x /root/cloud-backup
# CMD ["/root/cloud-backup"]

FROM oldwang6/alpine
WORKDIR /root/
COPY /app/cloud-backup .
RUN chmod +x /root/cloud-backup
CMD ["/root/cloud-backup"]