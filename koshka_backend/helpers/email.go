package helpers

import (
	"fmt"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"gopkg.in/gomail.v2"
	"log"
	"os"
	"strconv"
)

// SendEmail - Sends an email using SMTP
func MailGmail(to string, subject string, body string) error {
	smtpHost := os.Getenv("SMTP_HOST") // Example: "smtp.gmail.com"
	smtpPort := os.Getenv("SMTP_PORT") // Example: "587"
	smtpUser := os.Getenv("SMTP_USER") // Example: "your-email@gmail.com"
	smtpPass := os.Getenv("SMTP_PASS") // Example: "your-email-password"

	// Convert SMTP port to int
	port, err := strconv.Atoi(smtpPort)
	if err != nil {
		log.Println("Invalid SMTP_PORT:", smtpPort)
		return err
	}

	// Configure email message
	m := gomail.NewMessage()
	m.SetHeader("From", smtpUser)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	// Configure SMTP sender
	d := gomail.NewDialer(smtpHost, port, smtpUser, smtpPass)

	// Send email
	if err := d.DialAndSend(m); err != nil {
		log.Println("Failed to send email:", err)
		return err
	}

	log.Println("Email sent successfully to", to)
	return nil
}

// SendEmail - Sends an email using SendGrid API
func MailSendgrid(to string, subject string, body string) error {
	sendgridAPIKey := os.Getenv("SENDGRID_API_KEY")
	fromEmail := os.Getenv("SENDGRID_FROM_EMAIL")

	if sendgridAPIKey == "" || fromEmail == "" {
		log.Println("Missing SendGrid API credentials")
		return fmt.Errorf("missing SendGrid API credentials")
	}

	from := mail.NewEmail("Koshka Support", fromEmail)
	toRecipient := mail.NewEmail("", to)
	message := mail.NewSingleEmail(from, subject, toRecipient, "", body)

	client := sendgrid.NewSendClient(sendgridAPIKey)
	response, err := client.Send(message)

	if err != nil {
		log.Println("Failed to send email:", err)
		return err
	}

	log.Printf("Email sent! Status: %d\n", response.StatusCode)
	return nil
}

// SendEmailWithTemplate - Sends an email using SendGrid Templates
func SendEmailWithTemplate(to string, subject string, templateID string, dynamicData map[string]string) error {
	sendgridAPIKey := os.Getenv("SENDGRID_API_KEY")
	fromEmail := os.Getenv("SENDGRID_FROM_EMAIL")

	if sendgridAPIKey == "" || fromEmail == "" {
		log.Println("Missing SendGrid API credentials")
		return fmt.Errorf("missing SendGrid API credentials")
	}

	from := mail.NewEmail("Koshka Support", fromEmail)
	toRecipient := mail.NewEmail("", to)
	message := mail.NewV3MailInit(from, subject, toRecipient)

	// Attach template ID
	message.SetTemplateID(templateID)

	// Attach dynamic template data
	personalization := mail.NewPersonalization()
	personalization.AddTos(toRecipient)
	for key, value := range dynamicData {
		personalization.SetDynamicTemplateData(key, value)
	}
	message.AddPersonalizations(personalization)

	client := sendgrid.NewSendClient(sendgridAPIKey)
	response, err := client.Send(message)

	if err != nil {
		log.Println("Failed to send email:", err)
		return err
	}

	log.Printf("Email sent! Status: %d\n", response.StatusCode)
	return nil
}
