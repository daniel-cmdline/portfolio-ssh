# Estágio 1: Compilação do binário Go
FROM golang:1.25-alpine AS builder
WORKDIR /workspace

# 1. Copia as dependências primeiro para usar o cache do Docker
COPY go.mod go.sum ./
RUN go mod download

# 2. Copia somente os pacotes necessários para compilar o servidor SSH
COPY cmd ./cmd
COPY input ./input
COPY ui ./ui
COPY utils ./utils

# 3. Compila especificamente o seu servidor SSH
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server ./cmd/ssh-server/main.go


# Estágio 2: Ambiente de execução final ultra-leve
FROM alpine:latest
WORKDIR /root/

RUN apk add --no-cache ca-certificates

# 4. Copia apenas o binário compilado. A Host Key deve ser montada em runtime
# via SSH_HOST_KEY_PATH, ou o servidor gera uma chave efêmera.
COPY --from=builder /workspace/server .

EXPOSE 2222
CMD ["./server"]
