FROM golang:1.22.3-alpine3.19 AS builder

COPY . /app
WORKDIR /app

RUN go mod download
RUN go mod tidy
RUN go build -o ./bin/recipe_client cmd/client/main.go

FROM alpine:3.19

WORKDIR /app
COPY --from=builder /app/bin/recipe_client .
COPY config/.env .env
ENV ENV_CONFIG_FILE_PATH=.env

CMD ["./recipe_client"]
