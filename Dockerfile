FROM golang:1.19.3-alpine AS builder
LABEL maintainer="imlonghao <dockerfile@esd.cc>"
WORKDIR /builder
COPY . /builder
RUN apk add upx && \
    go build -ldflags="-s -w" -o /app && \
    upx --lzma --best /app

FROM alpine:latest
COPY --from=builder /app /
RUN apk add --no-cache tini
ENTRYPOINT ["/sbin/tini", "--"]
CMD ["/app"]
