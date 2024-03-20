# Builder stage
FROM golang:1.21.0 as builder
WORKDIR /usr/src/app
COPY . .
RUN go mod download && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd

# Final stage
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /usr/src/app/main .
COPY --from=builder /usr/src/app/.env .
EXPOSE 8080
CMD ["./main"]
