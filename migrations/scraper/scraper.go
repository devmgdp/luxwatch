package scraper

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
)

type ScrapedProduct struct {
	Name     string
	Price    float64
	ImageURL string
	URL      string
}

func ExtractData(url string) (ScrapedProduct, error) {
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36"),
	)

	var product ScrapedProduct
	product.URL = url

	c.OnHTML("h1, .ui-pdp-title, #productTitle, .product-name", func(e *colly.HTMLElement) {
		if product.Name == "" {
			product.Name = strings.TrimSpace(e.Text)
		}
	})

	c.OnHTML("meta[property='og:image']", func(e *colly.HTMLElement) {
		if product.ImageURL == "" {
			product.ImageURL = e.Attr("content")
		}
	})

	imageSelectors := []string{"img.ui-pdp-image", "#imgBlkFront", "img[data-testid='main-image']", ".poly-component__picture"}
	for _, selector := range imageSelectors {
		c.OnHTML(selector, func(e *colly.HTMLElement) {
			if product.ImageURL == "" {
				src := e.Attr("src")
				if strings.Contains(src, "data:image") || src == "" {
					src = e.Attr("data-src")
				}
				product.ImageURL = src
			}
		})
	}

	// Seletores focados no PREÇO ATUAL (ignora o riscado)
	priceSelectors := []string{
		".ui-pdp-price__second-line .andes-money-amount__fraction",
		".poly-price__current .andes-money-amount__fraction",
		".andes-money-amount__fraction",
	}

	for _, selector := range priceSelectors {
		c.OnHTML(selector, func(e *colly.HTMLElement) {
			if product.Price == 0 {
				rawPrice := e.Text
				clean := strings.ReplaceAll(rawPrice, "R$", "")
				clean = strings.ReplaceAll(clean, ".", "")
				clean = strings.ReplaceAll(clean, ",", ".")
				clean = strings.TrimSpace(clean)
				price, _ := strconv.ParseFloat(clean, 64)
				if price > 0 {
					product.Price = price
				}
			}
		})
	}

	c.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Scraper Error: Status %d on URL %s\n", r.StatusCode, r.Request.URL)
	})

	err := c.Visit(url)

	if product.Name == "" {
		product.Name = "Product Name Not Detected"
	}
	if product.ImageURL == "" {
		product.ImageURL = "https://via.placeholder.com/300?text=No+Image+Found"
	}

	return product, err
}
