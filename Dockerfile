FROM golang

RUN go version 
ENV GOPATH=/

COPY ./ ./
RUN go mod download
RUN go build -o ./bin/main ./cmd/mkc/main.go

EXPOSE 8080

CMD ["./bin/main"]