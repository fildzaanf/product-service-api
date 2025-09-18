package mailer

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
	"product-service-api/infrastructure/config"

	"text/template"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"gopkg.in/mail.v2"
)

func EmailNotification(to []string, templateContent string, data interface{}) (bool, error) {
	config, err := config.LoadConfig()
	if err != nil {
		logrus.Fatalf("failed to load smtp configuration: %v", err)
	}

	m := mail.NewMessage()
	m.SetHeader("From", config.SMTP.SMTP_USER)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", "Go Commerce Payment Notification")

	tmpl, err := template.New("emailTemplate").Parse(templateContent)
	if err != nil {
		return false, fmt.Errorf("failed to parse email template: %v", err)
	}

	var emailContent bytes.Buffer
	if err := tmpl.Execute(&emailContent, data); err != nil {
		return false, fmt.Errorf("failed to execute template: %v", err)
	}

	m.SetBody("text/html", emailContent.String())

	SMTP_PORT, err := strconv.Atoi(config.SMTP.SMTP_PORT)
	if err != nil {
		return false, fmt.Errorf("invalid smtp port: %v", err)
	}

	d := mail.NewDialer(
		config.SMTP.SMTP_HOST,
		SMTP_PORT,
		config.SMTP.SMTP_USER,
		config.SMTP.SMTP_PASS,
	)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	log.Printf("sending email to: %v", to)
	if err := d.DialAndSend(m); err != nil {
		return false, fmt.Errorf("failed to send email: %v", err)
	}
	log.Println("email successfully sent!")
	return true, nil
}

func SendEmailNotificationPayment(name, email, paymentCode, productName string, price decimal.Decimal, quantity int, totalAmount decimal.Decimal, status string, updatedAt time.Time) {
	go func() {
		rootPath, err := filepath.Abs("../") 
		if err != nil {
			log.Fatalf("failed to get root directory: %v", err)
		}

		filePath := filepath.Join(rootPath, "pkg/email/template/payment-notification.html")
		log.Printf("Looking for email template at: %s", filePath)
		log.Printf("Looking for email template at: %s", filePath)
		emailTemplate, err := os.ReadFile(filePath)
		if err != nil {
			log.Printf("failed to load email template: %v", err)
			return
		}

		data := map[string]string{
			"Name":        name,
			"Email":       email,
			"PaymentCode": paymentCode,
			"ProductName": productName,
			"Price":       price.String(),
			"Quantity":    strconv.Itoa(quantity),
			"TotalAmount": totalAmount.String(),
			"Status":      status,
			"UpdatedAt":   updatedAt.Format("2006-01-02 15:04:05"),
		}

		log.Printf("Preparing email for: %s", email)
		success, errEmail := EmailNotification([]string{email}, string(emailTemplate), data)
		if !success || errEmail != nil {
			log.Printf("Failed to send notification email to %s: %v", email, errEmail)
		}
	}()
}
