# --- ESTÁGIO DE COMPILAÇÃO ---
FROM golang:1.24 AS builder

# ESSA LINHA É A CHAVE: Força o Go a ignorar a discrepância de versão
ENV GOTOOLCHAIN=go1.24.0

RUN apt-get update && apt-get install -y gcc libc6-dev

WORKDIR /app

# Copia tudo de uma vez
COPY . .

# Remove a trava de versão e compila
RUN go mod edit -go=1.24
RUN go mod tidy
RUN go build -o main ./cmd/server/main.go

# --- ESTÁGIO FINAL ---
# Usamos Debian Slim para manter o container leve e compatível
FROM debian:bookworm-slim

# Instala Node.js e dependências do Chromium para o scraper
RUN apt-get update && apt-get install -y \
    curl \
    gnupg \
    nodejs \
    npm \
    chromium \
    ca-certificates \
    fonts-liberation \
    --no-install-recommends && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Traz o binário e as pastas do estágio builder
COPY --from=builder /app/main .
COPY --from=builder /app/suggestions ./suggestions
COPY --from=builder /app/static ./static
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/package*.json ./migrations/scraper/

# Instala as dependências do scraper (Puppeteer/Playwright)
RUN cd migrations/scraper && npm install

EXPOSE 8080

CMD ["./main"]