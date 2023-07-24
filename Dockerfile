FROM golang:1.20-alpine

RUN apk add --no-cache git

WORKDIR /usr/src/buggybox

COPY go.* ./
RUN go mod download && go mod verify

COPY . .
RUN hash=$(git rev-parse --short HEAD) \
 && build=$(git log -s --pretty='format:%cd' --date=format:'%Y-%m-%d' $hash | head) \
 && go build -ldflags "-X buggybox/config.AppVersion=1.0.0 -X buggybox/config.AppBuildHash=$hash -X buggybox/config.AppBuildTime=$build" -v -o /usr/local/bin/buggybox .

CMD ["buggybox"]