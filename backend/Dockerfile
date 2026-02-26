# Dockerfile para a aplicação Go
# Multi-stage build para imagem final pequena

# Etapa 1: Build
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copia os arquivos de dependências
COPY go.mod go.sum ./
RUN go mod download

# Copia o código fonte
COPY . .

# Compila a aplicação
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api

# Etapa 2: Imagem final
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copia o binário compilado
COPY --from=builder /app/main .

# Expõe a porta
EXPOSE 8080

# Comando de execução
CMD ["./main"]
