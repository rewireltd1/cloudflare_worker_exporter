FROM golang:1.15-alpine as builder
RUN apk --update add ca-certificates

RUN mkdir /cloudflare_worker_exporter
WORKDIR /cloudflare_worker_exporter
COPY go.mod go.sum ./

RUN go mod download
COPY . .

RUN CGO_ENABLED=0  GOOS=linux GOARCH=amd64 go build -a -o cloudflare_worker_exporter

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /cloudflare_worker_exporter/cloudflare_worker_exporter /cloudflare_worker_exporter
ENTRYPOINT ["/cloudflare_worker_exporter"]
