# Name this stage as 'builder'
FROM golang:1.21.0 as builder
WORKDIR /usr/src/app
COPY . .
RUN go mod download


CMD ["go", "run", "/usr/src/app/cmd", "."]

