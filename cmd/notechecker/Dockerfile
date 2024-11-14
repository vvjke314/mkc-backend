FROM golang

WORKDIR /app
COPY ./ ./
COPY ./../internal ./internal
COPY ./../config.yaml ./config.yaml

RUN go mod download
RUN go build -o ./notechecker ./cmd/notechecker/notechecker.go
ENV TZ=Europe/Moscow

ENTRYPOINT ["./notechecker"]