package mailer

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendPriceAlert(userEmail string, productName string, oldPrice float64, newPrice float64, url string) error {
	from := os.Getenv("SMTP_EMAIL")
	password := os.Getenv("SMTP_PASS")

	// Configuração do servidor SMTP do Gmail
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	auth := smtp.PlainAuth("", from, password, smtpHost)

	body := fmt.Sprintf(
		"Subject: 💎 BAIXOU O PREÇO: %s\r\n\r\n"+
			"Boas notícias! O produto que você está rastreando caiu de preço.\n\n"+
			"Preço anterior: R$ %.2f\n"+
			"Preço ATUAL: R$ %.2f\n\n"+
			"Link para aproveitar: %s",
		productName, oldPrice, newPrice, url,
	)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{userEmail}, []byte(body))
	return err
}
