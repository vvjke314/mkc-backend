FROM golang AS builder

WORKDIR /app
COPY ./ ./
COPY ./../internal ./internal
COPY ./../config.yaml ./config.yaml

RUN go mod download
RUN go build -o ./backup_app ./cmd/db-backup/db-backup.go

FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y postgresql-client
RUN mkdir -p /backups

COPY --from=builder /app/backup_app /backup_app
COPY --from=builder /app/config.yaml /config.yaml
ENTRYPOINT ["/backup_app"]