FROM golang:1.21-alpine

WORKDIR /app

COPY main.go .

RUN go mod init backend2-api && \
    go mod tidy && \
    go build -o main .

CMD ["./main"]



