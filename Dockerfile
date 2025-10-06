# Dockerfile
ARG APP_NAME=Vajean
FROM golang:1.25 AS builder

WORKDIR /app

COPY . .

ARG APP_NAME=Vajean
ENV APP_NAME=vajean

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o vajean .

FROM alpine

WORKDIR /app

ARG APP_NAME=Vajean
ENV APP_NAME=vajean

COPY --from=builder /app/vajean .

# Fix permissions
RUN chmod +x vajean

EXPOSE 8080
ENTRYPOINT ["sh", "-c", "./vajean"]
