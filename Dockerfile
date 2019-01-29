# build application
FROM golang:1.11-alpine AS build
ARG BUILDPATH=github.com/thehowl/ip.zxq.co

COPY * /go/src/${BUILDPATH}/

RUN apk -U add git && \
    cd /go/src/${BUILDPATH}/ && \
    go get -v && \
    go build -o /dist/main

# create a new image
FROM alpine:latest
COPY --from=build /dist/main /app/main
COPY GeoLite2-City.mmdb /app/GeoLite2-City.mmdb

WORKDIR /app/
ENTRYPOINT [ "/app/main" ]