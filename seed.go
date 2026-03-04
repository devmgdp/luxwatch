package main

import (
	"fmt"
	"log"
	"luxwatch/internal/database"
)

func main() {
	// 1. Inicializa a conexão com o banco
	database.InitDB()

	// Limpar histórico antigo apenas para esses produtos de teste

	// LISTA DE PRODUTOS PARA TESTE
	// Produto 1: Lavadora (Categoria House)
	p1 := "Lavadora De Alta Pressão Bosch Ghp 200 1700w 2000 Psi Com Acessórios Azul 220v 60"

	// Produto 2: Smartphone (Categoria Electronics)
	p2 := "Smartphone Samsung Galaxy S24 Ultra 512GB Titanium Gray"

	// INSERÇÃO PRODUTO 1 (House)
	insertHistory(p1, 1200.00, "2026-01-15 10:00:00")
	insertHistory(p1, 1150.50, "2026-02-10 14:30:00")
	insertHistory(p1, 1084.00, "2026-03-01 12:00:00")

	// INSERÇÃO PRODUTO 2 (Electronics)
	insertHistory(p2, 6500.00, "2026-01-20 09:00:00")
	insertHistory(p2, 6200.00, "2026-02-15 11:00:00")
	insertHistory(p2, 5890.00, "2026-03-02 16:45:00")

}

// Função auxiliar para evitar repetição de código e garantir os placeholders do Postgres ($1, $2...)
func insertHistory(name string, price float64, date string) {
	query := `INSERT INTO price_history (product_name, price, checked_at) VALUES ($1, $2, $3)`
	_, err := database.DB.Exec(query, name, price, date)
	if err != nil {
		log.Printf("❌ Erro ao inserir %s: %v", name, err)
	} else {
		fmt.Printf("✔ Dados inseridos para: %s (R$ %.2f)\n", name, price)
	}
}
