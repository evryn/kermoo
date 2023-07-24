FROM golang:1.20-alpine

RUN apk add --no-cache git

WORKDIR /usr/src/app

COPY go.* ./
RUN go mod download && go mod verify

COPY . .
RUN go build -ldflags "-X main.Version=1.0.0 -X 'main.Build=$(date)'" -v -o /usr/local/bin/buggybox ./...

CMD ["buggybox"]