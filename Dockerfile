# --- ESTÁGIO DE COMPILAÇÃO ---
FROM golang:latest AS builder
RUN apk add --no-cache gcc musl-dev
WORKDIR /app

# Copia tudo de uma vez para não travar no download isolado
COPY . .

# Ignora a verificação rígida de versão do Go
RUN go mod tidy
RUN go build -o main ./cmd/server/main.go

# --- ESTÁGIO FINAL ---
FROM alpine:latest
RUN apk add --no-cache nodejs npm chromium nss freetype harfbuzz ca-certificates ttf-freefont

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/suggestions ./suggestions
COPY --from=builder /app/static ./static
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/package*.json ./migrations/scraper/

RUN cd migrations/scraper && npm install

EXPOSE 8080
CMD ["./main"]