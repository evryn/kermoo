# syntax=docker/dockerfile:1

# Build Stage
FROM golang:1.20-alpine3.18 as builder

ARG BUILD_VERSION
ARG BUILD_DATE
ARG BUILD_REF

WORKDIR /go/src/app
COPY . .
RUN go get -d -v ./... \
 && go build -ldflags "-X buggybox/config.BuildVersion=\"$BUILD_VERSION\" -X buggybox/config.BuildRef=\"$BUILD_REF\" -X buggybox/config.BuildDate=\"$BUILD_DATE\"" -o /go/bin/app -v .

# A layer for running automated tests
FROM builder as test
RUN go test ./...

# A layer with final executable
FROM alpine:3.18
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/bin/app /usr/local/bin/buggybox
EXPOSE 80
ENTRYPOINT ["buggybox"]
