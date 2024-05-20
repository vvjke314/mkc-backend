FROM golang

WORKDIR /app
COPY . .
COPY ../internal .

RUN go mod download
RUN go build -o ./backup_app ./db-backup.go

FROM debian:buster

RUN apt-get update && apt-get install -y postgresql-client

COPY --from=builder /backup_app /backup_app
ENTRYPOINT ["/backup_app"]