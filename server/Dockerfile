FROM golang:1.23.0-alpine AS build-env

LABEL maintainer="bondhan novandy<bondhan.novandy@gmail.com>"
ARG CONDUCTOR_SERVER_URL
ARG WORKFLOW_VERSION

ENV APP_NAME=server
ENV GO111MODULE=on
ENV GOPRIVATE=github.com/bondhan
ENV TZ=Asia/Jakarta
ENV GIT_TERMINAL_PROMPT=0
ENV CGO_ENABLED=0


RUN apk update && apk upgrade
RUN apk add --no-cache --virtual .build-deps --no-progress -q \
    bash \
    curl \
    busybox-extras \
    make \
    git \
    tzdata && \
    cp /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone
RUN apk update && apk add --no-cache coreutils

WORKDIR /src

RUN mkdir -p /src/server
COPY . /src/server

WORKDIR /src/server

RUN echo "CONDUCTOR_SERVER_URL=${CONDUCTOR_SERVER_URL}" >> .env
RUN echo "WORKFLOW_VERSION=${WORKFLOW_VERSION}" >> .env

RUN ls -lah

RUN go build -o ./bin/$APP_NAME main.go

# clean container
FROM alpine:latest

RUN apk add --no-cache --virtual .build-deps --no-progress -q \
    bash \
    tzdata

ENV TZ=Asia/Jakarta
ENV APP_NAME=server

RUN mkdir /app
WORKDIR /app

COPY --from=build-env /src/server/bin/$APP_NAME /app/$APP_NAME
COPY --from=build-env /src/server/.env /app/.env

EXPOSE 3000

CMD ["/app/server"]