# Name this stage as 'builder'
FROM golang:1.21.0 as builder
WORKDIR /usr/src/app
COPY . .
RUN go mod download

EXPOSE 4000

CMD ["go", "run", "/usr/src/app/cmd", "."]

## Этап сборки
 #FROM golang:1.21.0 as builder
 #WORKDIR /usr/src/app
 #COPY . .
 #RUN go mod download && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .
 #
 ## Этап финального образа
 #FROM alpine:latest
 #WORKDIR /root/
 ## Копируем скомпилированный файл из предыдущего этапа
 #COPY --from=builder /usr/src/app/main .
 ## Запускаем скомпилированное приложение
 #CMD ["./main"]