# Dockerfile
FROM golang:1.25 AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o vajean .

FROM alpine:3.22.1
WORKDIR /app
COPY --from=builder /app/vajean ./vajean
RUN chmod +x vajean

EXPOSE 8080
ENTRYPOINT ["./vajean"]
