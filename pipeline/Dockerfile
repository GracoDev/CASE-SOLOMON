FROM golang:1.21-alpine

WORKDIR /app

COPY main.go .

RUN go mod init pipeline && \
    go mod tidy && \
    go build -o main .

CMD ["./main"]



