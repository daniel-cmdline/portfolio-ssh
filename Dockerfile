# Estágio 1: Compilação do binário Go
FROM golang:1.25-alpine AS builder
WORKDIR /workspace

# 1. Copia as dependências primeiro para usar o cache do Docker
COPY go.mod go.sum ./
RUN go mod download

# 2. Copia TODO o seu projeto (api, cmd, input, ui, id_rsa, etc.)
COPY . .

# 3. Compila especificamente o seu servidor SSH
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server ./cmd/ssh-server/main.go


# Estágio 2: Ambiente de execução final ultra-leve
FROM alpine:latest
WORKDIR /root/

# 4. Copia o binário compilado do estágio anterior
COPY --from=builder /workspace/server .

# 5. Copia as suas pastas de módulos e layouts que o servidor SSH precisa ler em tempo de execução
COPY --from=builder /workspace/ui ./ui
COPY --from=builder /workspace/api ./api
COPY --from=builder /workspace/input ./input

# 6. Copia a chave privada que o seu código usa para assinar o SSH
COPY --from=builder /workspace/id_rsa .

EXPOSE 22
CMD ["./server"]