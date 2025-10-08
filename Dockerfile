# ---------- build stage ----------
FROM golang:1.25 AS build
WORKDIR /src

# Cache deps
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build the cmd/bot target
# (strip symbols, static binary)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o /out/valjean ./cmd/bot

# ---------- runtime stage ----------
FROM gcr.io/distroless/base-debian12
WORKDIR /
COPY --from=build /out/valjean /valjean
ENV ADDR=:8080
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/valjean"]
