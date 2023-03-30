FROM alpine:latest

RUN mkdir /app

ENV ADDRESS server:8080

COPY ./cmd/agent/agentApp /app

CMD [ "/app/agentApp"]
#RUN /app/agentApp -a server:8080