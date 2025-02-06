package utils

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
)

// SendEmail sends an email using SMTP.
func SendEmail(to, subject, body string) error {
	from := os.Getenv("SMTP_USERNAME")
	pass := os.Getenv("SMTP_PASSWORD")
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	if from == "" || pass == "" || host == "" || port == "" {
		return fmt.Errorf("SMTP configuration not set")
	}
	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" + body

	log.Printf("[DEBUG] Attempting to send email to %s using server %s:%s", to, host, port)
	err := smtp.SendMail(host+":"+port, smtp.PlainAuth("", from, pass, host), from, []string{to}, []byte(msg))
	if err != nil {
		log.Printf("[ERROR] smtp.SendMail error: %v", err)
		return fmt.Errorf("error sending email: %v", err)
	}
	log.Printf("[DEBUG] Email successfully sent to %s", to)
	return nil
}
