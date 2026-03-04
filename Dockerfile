# --- ESTÁGIO DE COMPILAÇÃO ---
FROM golang:1.22-alpine AS builder
RUN apk add --no-cache gcc musl-dev
WORKDIR /app

# Copia dependências primeiro (cache)
COPY go.mod go.sum ./
RUN go mod download

# Copia o restante do código
COPY . .

# CORREÇÃO DO CAMINHO: Aponta para onde o main.go realmente está
RUN go build -o main ./cmd/server/main.go

# --- ESTÁGIO FINAL ---
FROM alpine:latest
RUN apk add --no-cache nodejs npm chromium nss freetype harfbuzz ca-certificates ttf-freefont

WORKDIR /app

# Copia o binário compilado
COPY --from=builder /app/main .

# Copia as pastas de assets e templates
COPY --from=builder /app/suggestions ./suggestions
COPY --from=builder /app/static ./static
COPY --from=builder /app/migrations ./migrations

# CORREÇÃO DO NPM: Copia os arquivos de configuração do Node para a pasta certa
COPY package*.json ./migrations/scraper/
RUN cd migrations/scraper && npm install

EXPOSE 8080

CMD ["./main"]