# build stage
FROM golang:1.22 AS build
WORKDIR /go/src/app

ARG TAG
ENV GO111MODULE=on
ENV GOPROXY=http://build.evixo.ru:9050

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-X cmd/version.appVersion=${TAG}-$(date -u +%Y%m%d%H%M)" -o "./bin" main.go

# final stage
FROM alpine:3.11

WORKDIR /app

RUN apk add  --update --no-cache tzdata ca-certificates nano curl \
    && rm -rf /tmp/* \
    && rm -rf /var/cache/apk/*

COPY ./deployments/build/docker/app/files /
COPY ./addr_cache.json /app/
COPY ./init.sql /app/
COPY ./id_rsa.key /app/
COPY --from=build /go/src/app/bin bin

RUN chmod +x /docker/bin/* \
    && chmod +x bin

ENTRYPOINT ["/docker/bin/entrypoint.sh"]