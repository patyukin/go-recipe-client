FROM golang:1.22.3-alpine3.19 AS builder

ENV config=docker

WORKDIR /app

COPY . /app

ENV YAML_CONFIG_FILE_PATH=config/config.yaml
ENV ENV_CONFIG_FILE_PATH=config/.env

RUN go mod tidy && \
    go mod download && \
    go get github.com/githubnemo/CompileDaemon && \
    go install github.com/githubnemo/CompileDaemon

ENTRYPOINT CompileDaemon --build="go build -o bin/client cmd/client/main.go" --command=./bin/client
