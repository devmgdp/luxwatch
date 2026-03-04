# --- ESTÁGIO DE COMPILAÇÃO ---
FROM golang:1.24 AS builder

RUN apt-get update && apt-get install -y gcc libc6-dev

WORKDIR /app

# Copia os arquivos de módulo
COPY go.mod go.sum ./
# Agora o download vai funcionar direto!
RUN go mod download

# Copia o resto do código
COPY . .

# Compila
RUN go build -o main ./cmd/server/main.go

# --- ESTÁGIO FINAL ---
FROM debian:bookworm-slim

# Instala Node e Chromium para o scraper
RUN apt-get update && apt-get install -y \
    curl \
    gnupg \
    nodejs \
    npm \
    chromium \
    ca-certificates \
    --no-install-recommends && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/suggestions ./suggestions
COPY --from=builder /app/static ./static
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/package*.json ./migrations/scraper/

RUN cd migrations/scraper && npm install
RUN npx playwright install chromium

EXPOSE 8080
CMD ["./main"]