# Базовый образ на основании которого мы создаем свой
FROM golang

RUN go version 

WORKDIR /mkc-backend/
COPY . .
RUN go mod download
RUN go build -o ./bin/mkc ./cmd/mkc/main.go
RUN mkdir storage
ENV TZ=Europe/Moscow

EXPOSE 8080

CMD ["./bin/mkc"]