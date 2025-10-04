# Dockerfile
ARG APP_NAME=go-module
FROM golang:1.24 AS builder

WORKDIR /app

COPY . .

ARG APP_NAME=go-module
ENV APP_NAME=${APP_NAME}

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ${APP_NAME} .

FROM alpine

WORKDIR /app

ARG APP_NAME=go-module
ENV APP_NAME=${APP_NAME}

COPY --from=builder /app/${APP_NAME} .

# Fix permissions
RUN chmod +x ${APP_NAME}

EXPOSE 8080
ENTRYPOINT ["sh", "-c", "./${APP_NAME}"]
