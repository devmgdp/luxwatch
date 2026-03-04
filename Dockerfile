# --- ESTÁGIO DE COMPILAÇÃO ---
# Usamos a imagem oficial (Debian) para garantir suporte a versões novas do Go
FROM golang:1.24 AS builder

# No Debian, usamos apt-get em vez de apk
RUN apt-get update && apt-get install -y gcc libc6-dev

WORKDIR /app

# Copia tudo para resolver o problema da sua versão 1.25 local
COPY . .

# Garante que as dependências estejam certas antes de buildar
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