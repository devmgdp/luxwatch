# Estágio de Compilação
FROM golang:1.22-alpine AS builder

# Instalar dependências necessárias para compilar pacotes CGO
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

# Isso acelera o build pois o Docker faz cache das bibliotecas
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Compila o binário
RUN go build -o main .

# Estágio Final
FROM alpine:latest

# Instalar Node.js e dependências para o Playwright funcionar no Linux
RUN apk add --no-cache \
    nodejs \
    npm \
    chromium \
    nss \
    freetype \
    harfbuzz \
    ca-certificates \
    ttf-freefont

WORKDIR /app

# Copia o binário e as pastas necessárias do estágio builder
COPY --from=builder /app/main .
COPY --from=builder /app/suggestions/templates ./suggestions/templates
COPY --from=builder /app/static ./static
COPY --from=builder /app/migrations/scraper ./migrations/scraper

# Instalar as dependências do Node.js dentro do container
RUN cd migrations/scraper && npm install

# Expõe a porta que definimos no main.go
EXPOSE 8080

CMD ["./main"]