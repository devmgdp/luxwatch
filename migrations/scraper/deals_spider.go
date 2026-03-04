package scraper

import (
	"encoding/json"
	"fmt"
	"luxwatch/internal/database"
	"os/exec"
	"strconv"
	"strings"
)

type Item struct {
	Name  string `json:"name"`
	Price string `json:"price"`
	Image string `json:"image"`
	Link  string `json:"link"`
}

func RunDealFinder() {
	// 1. Limpa os dados antigos antes de novos scans
	database.DB.Exec("DELETE FROM daily_deals")

	url := "https://www.mercadolivre.com.br/ofertas"

	// IMPORTANTE: Verifique se o caminho "migrations/scraper/stealth_browser.js" está correto
	cmd := exec.Command("node", "migrations/scraper/stealth_browser.js", url)
	out, err := cmd.Output()

	if err != nil {
		fmt.Println("Erro ao executar o Node.js:", err)
		return
	}

	var items []Item
	err = json.Unmarshal(out, &items)
	if err != nil {
		fmt.Println("Erro ao converter JSON do Scraper:", err)
		return
	}

	// Contador para sabermos quantos foram salvos
	count := 0

	for _, it := range items {
		p := cleanPrice(it.Price)
		if p <= 0 {
			continue
		}
		cat := categorize(it.Name)

		// Usamos Exec para inserir no banco
		_, err := database.DB.Exec("INSERT INTO daily_deals (name, price, image_url, product_url, category) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (product_url) DO UPDATE SET price = $2",
			it.Name, p, it.Image, it.Link, cat)

		if err == nil {
			// Salva no histórico também
			var exists bool
			checkQuery := `
        SELECT EXISTS(
            SELECT 1 FROM price_history 
            WHERE product_name = $1 
            AND checked_at >= CURRENT_DATE
        )`
			database.DB.QueryRow(checkQuery, it.Name).Scan(&exists)

			if !exists {
				// Só insere se não houver registro hoje
				database.DB.Exec("INSERT INTO price_history (product_name, price) VALUES ($1, $2)", it.Name, p)
			}

			count++
		}
	}
}

func cleanPrice(p string) float64 {
	// Remove tudo que não é número ou vírgula
	p = strings.ReplaceAll(p, "R$", "")
	p = strings.ReplaceAll(p, ".", "")
	p = strings.ReplaceAll(p, ",", ".")
	p = strings.TrimSpace(p)
	val, _ := strconv.ParseFloat(p, 64)
	return val
}

func categorize(n string) string {
	n = strings.ToLower(n)

	electronicsKeywords := []string{
		"celular", "iphone", "samsung", "motorola", "xiaomi", "redmi",
		"camera", "câmera", "intelbras", "wi-fi", "monitor", "tv", "televisao",
		"fone", "headset", "notebook", "laptop", "pc", "gamer", "teclado", "mouse",
		"tablet", "ipad", "alexa", "echo", "smartwatch", "relogio inteligente",
	}

	for _, key := range electronicsKeywords {
		if strings.Contains(n, key) {
			return "Electronics"
		}
	}

	if strings.Contains(n, "whey") || strings.Contains(n, "creatina") || strings.Contains(n, "suplemento") {
		return "Suplements"
	}

	// Se não for nenhum acima, cai em House
	return "House"
}
