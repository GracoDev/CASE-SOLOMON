FROM golang:1.21-alpine

WORKDIR /app

# Copiar código fonte
COPY main.go .

# Inicializar módulo Go e instalar dependências
RUN go mod init pipeline && \
    go get github.com/lib/pq && \
    go mod tidy && \
    go build -o main .

CMD ["./main"]



