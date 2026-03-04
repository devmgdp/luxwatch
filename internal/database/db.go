package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq" // Driver do PostgreSQL
)

var DB *sql.DB

func InitDB() {
	var err error

	// Formato: postgres://usuario:senha@localhost:5432/nome_do_banco?sslmode=disable
	// SUBSTITUA 'suasenha' pela senha real do seu pgAdmin
	connStr := os.Getenv("DB_URL")

	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Erro ao abrir conexão:", err)
	}

	// Verifica se a conexão com o servidor do Postgres está ativa
	err = DB.Ping()
	if err != nil {
		log.Fatal("❌ Erro real: Não foi possível conectar ao PostgreSQL:", err)
	}

	// Criar tabelas no Postgres (Sintaxe SERIAL em vez de AUTOINCREMENT)
	createTables()
}

func createTables() {
	// Tabela de Ofertas Atuais
	_, err := DB.Exec(`CREATE TABLE IF NOT EXISTS daily_deals (
		id SERIAL PRIMARY KEY,
		name TEXT,
		price DECIMAL(10,2),
		image_url TEXT,
		product_url TEXT UNIQUE,
		category TEXT
	)`)
	if err != nil {
		log.Println("Erro ao criar daily_deals:", err)
	}

	// Tabela de Histórico
	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS price_history (
		id SERIAL PRIMARY KEY,
		product_name TEXT,
		price DECIMAL(10,2),
		checked_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		log.Println("Erro ao criar price_history:", err)
	}

	// Tabela de Usuários
	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`)

	// Tabela Watchlist atualizada (vincular ao ID do usuário)
	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS watchlist (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    product_url TEXT NOT NULL,
    target_price DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, product_url)
	)`)
}
