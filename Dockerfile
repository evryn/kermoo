# syntax=docker/dockerfile:1
FROM golang:1.20-alpine

ARG BUILD_VERSION="dev"
ARG BUILD_DATE="2001-01-01T00:00:00Z"
ARG BUILD_REF="123456"

# Meta variables for annonating the labels
ARG META_TITLE="BuggyBox"
ARG META_DESCRIPTION="An app with the purpose of demonstating real-world malfunctioning applications. Good for testing and learnign container management topics."
ARG META_VENDOR="Evryn"
ARG META_AUTHOR="Amirreza Nasiri <nasiri.amirreza.96@gmail.com>"
ARG META_SOURCE="https://github.com/evryn/buggybox"
ARG META_HOME=$META_SOURCE

# Annotations for generally recognized labels
LABEL maintainer  = $META_AUTHOR
LABEL description = $META_DESCRIPTION

# Annotations for org.opencontainers
# @see https://specs.opencontainers.org/image-spec/annotations/
LABEL org.opencontainers.image.title         = $META_TITLE
LABEL org.opencontainers.image.description   = $META_DESCRIPTION
LABEL org.opencontainers.image.licenses      = "Apache-2.0"
LABEL org.opencontainers.image.vendor        = $META_VENDOR
LABEL org.opencontainers.image.version       = $BUILD_VERSION
LABEL org.opencontainers.image.source        = $META_SOURCE
LABEL org.opencontainers.image.url           = $META_HOME
LABEL org.opencontainers.image.documentation = $META_HOME
LABEL org.opencontainers.image.revision      = $BUILD_REF
LABEL org.opencontainers.image.authors       = $META_AUTHOR
LABEL org.opencontainers.image.created       = $BUILD_DATE

# Annotations for org.label-schema (for backward compatibility) and generally recognized labels
# @see http://label-schema.org/rc1/
LABEL org.label-schema.schema-version    = "1.0"
LABEL org.label-schema.name              = $META_TITLE
LABEL org.label-schema.description       = $META_DESCRIPTION
LABEL org.label-schema.vendor            = $META_VENDOR
LABEL org.label-schema.version           = $BUILD_VERSION
LABEL org.label-schema.vcs-ref           = $BUILD_REF
LABEL org.label-schema.build-date        = $BUILD_DATE
LABEL org.label-schema.vcs-url           = $META_SOURCE
LABEL org.label-schema.url               = $META_HOME
LABEL org.label-schema.usage             = $META_HOME
LABEL org.label-schema.docker.cmd        = "docker run -d -p 8080:8080 -p 5000:5000 buggybox start fixed"
LABEL org.label-schema.docker.cmd.help   = "docker exec -it $CONTAINER buggybox --help"

WORKDIR /usr/src/buggybox

COPY go.* ./
RUN go mod download && go mod verify

COPY . .

RUN go build -ldflags "-X buggybox/config.BuildVersion=\"$BUILD_VERSION\" -X buggybox/config.BuildRef=\"$BUILD_REF\" -X buggybox/config.BuildDate=\"$BUILD_DATE\"" -v -o /usr/local/bin/buggybox .

CMD ["buggybox"]