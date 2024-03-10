FROM golang:1.18 as builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -v -o ipwhitelister

# Docker multi-stage. App image based on distroless for a smaller final image
FROM gcr.io/distroless/static-debian11
RUN apk --no-cache add ca-certificates

WORKDIR /app/

COPY --from=builder /app/ipwhitelister .

CMD ["./ipwhitelister"]
