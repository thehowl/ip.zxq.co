# build application
FROM golang:1.15-alpine AS build

COPY * /root/

RUN cd /root/ && \
    go get -v && \
    CGO_ENABLED=0 go build -o /dist/main

# create a new image
FROM alpine:latest
COPY --from=build /dist/main /app/main

WORKDIR /app/
ENTRYPOINT [ "/app/main" ]
