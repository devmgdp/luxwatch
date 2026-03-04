FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/suggestions/templates ./suggestions/templates
COPY --from=builder /app/static ./static
# Se o script node estiver em migrations/scraper:
COPY --from=builder /app/migrations/scraper ./migrations/scraper

# Instalar Node.js para o scraper funcionar dentro do container
RUN apk add --no-cache nodejs npm
RUN cd migrations/scraper && npm install playwright-extra puppeteer-extra-plugin-stealth

CMD ["./main"]