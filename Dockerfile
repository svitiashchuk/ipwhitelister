FROM golang:1.18 as builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /app
COPY . .

ENV CGO_ENABLED=1
ENV GOOS=$TARGETOS
ENV GOARCH=$TARGETARCH

# See https://github.com/mattn/go-sqlite3/blob/master/_example/simple/Dockerfile
RUN go build \
    -ldflags='-s -w -extldflags "-static"' \
    -o ipwhitelister \
    ./cmd/ipwhitelister

# Docker multi-stage. App image based on distroless for a smaller final image
FROM gcr.io/distroless/static-debian12 as ipwhitelister

COPY --from=builder /app/ipwhitelister /app/

CMD ["./ipwhitelister"]
