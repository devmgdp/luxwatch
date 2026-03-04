package cron

import (
	"luxwatch/internal/database"
	"luxwatch/internal/mailer"
	"luxwatch/migrations/scraper"
	"time"
)

func StartPriceMonitor() {
	go func() {
		for {

			query := `
				SELECT u.email, w.product_url, w.target_price 
				FROM watchlist w
				JOIN users u ON w.user_id = u.id`

			rows, err := database.DB.Query(query)
			if err != nil {
				time.Sleep(10 * time.Minute) // Espera um pouco antes de tentar o banco de novo
				continue
			}

			for rows.Next() {
				var email, url string
				var targetPrice float64
				rows.Scan(&email, &url, &targetPrice)

				// Faz o scrape de UM produto
				currentProduct, err := scraper.ExtractData(url)

				if err == nil && currentProduct.Price > 0 && currentProduct.Price < targetPrice {
					mailer.SendPriceAlert(email, currentProduct.Name, targetPrice, currentProduct.Price, url)

					database.DB.Exec(`
						UPDATE watchlist SET target_price = $1 
						WHERE product_url = $2 AND user_id = (SELECT id FROM users WHERE email = $3)`,
						currentProduct.Price, url, email)
				}

				// Aguarda 5 segundos entre cada produto da lista para não ser bloqueado
				time.Sleep(5 * time.Second)
			}
			rows.Close()

			time.Sleep(2 * time.Hour)
		}
	}()
}
