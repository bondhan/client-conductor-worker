FROM golang:1.23.0 AS builder

ARG POLLING_TIME
ARG BATCH_SIZE
ARG CONDUCTOR_SERVER_URL

RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app
RUN echo "POLLING_TIME=${POLLING_TIME}" >> .env
RUN echo "POLLING_TIME=${BATCH_SIZE}" >> .env
RUN echo "POLLING_TIME=${CONDUCTOR_SERVER_URL}" >> .env

COPY . .
RUN CGO_ENABLED=0 go build -o worker main.go

FROM alpine:latest

RUN mkdir -p /app

WORKDIR /app

COPY --from=builder /usr/src/app/worker .
COPY --from=builder /usr/src/app/.env .

EXPOSE 8080

ENTRYPOINT ["/app/worker"]