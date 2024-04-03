# Базовый образ на основании которого мы создаем свой
FROM golang

RUN go version 

WORKDIR /home/mkc-backend/
COPY . .
RUN go mod download
RUN go build -o ./bin/main ./cmd/mkc/main.go

EXPOSE 8080

CMD ["./bin/mkc"]