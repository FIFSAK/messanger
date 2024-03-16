# Name this stage as 'builder'
FROM golang:1.21.0 as builder

WORKDIR /usr/src/app

# Copy the entire project from the current directory to the Working Directory inside the container
COPY . .

# Download all dependencies. Dependencies will be cached if the go.mod and the go.sum files are not changed
RUN go mod download

# Build the Go app inside the cmd directory
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd

# Start a new stage from scratch
FROM alpine:latest

WORKDIR /root/

# Copy the Pre-built binary file from the 'builder' stage
COPY --from=builder /usr/src/app/main .

# Copy .env file into the final stage to avoid it being skipped
#COPY --from=builder /usr/src/app/.env .
#
## Now COPY can find the '.env' file since it has been explicitly copied
#RUN ls -la /root/ && cat /root/.env

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
