FROM alpine:latest

RUN mkdir /app

ENV ADDRESS server:8080

COPY ./cmd/server/serverApp /app
COPY ../../app/listAllMetrics.html /app


CMD [ "/app/serverApp"]