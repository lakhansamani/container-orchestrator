FROM golang:1.22-alpine as builder
WORKDIR /app
COPY . .

RUN go build && \
    chmod 777 container-orchestrator

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/container-orchestrator .
EXPOSE 5600
ENTRYPOINT [ "./container-orchestrator" ]