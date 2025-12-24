package services

import (
	"fmt"
	"net/smtp"
	"os"
)

type EmailService struct {
	smtpHost string
	smtpPort string
	username string
	password string
}

func NewEmailService() *EmailService {
	return &EmailService{
		smtpHost: os.Getenv("SMTP_HOST"),
		smtpPort: os.Getenv("SMTP_PORT"),
		username: os.Getenv("SMTP_USERNAME"),
		password: os.Getenv("SMTP_PASSWORD"),
	}
}

func (s *EmailService) SendOTP(to, code string) error {
	from := s.username
	subject := "Your Verification Code - Khusa Mahal"
	body := fmt.Sprintf("Your verification code is: %s\n\nThis code will expire in 10 minutes.", code)

	message := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: %s\r\n"+
		"\r\n"+
		"%s\r\n", to, subject, body))

	auth := smtp.PlainAuth("", s.username, s.password, s.smtpHost)

	// In development/testing if no creds, just log it
	if s.username == "" || s.password == "" {
		fmt.Printf(" [MOCK EMAIL] To: %s | OTP: %s\n", to, code)
		return nil
	}

	err := smtp.SendMail(s.smtpHost+":"+s.smtpPort, auth, from, []string{to}, message)
	if err != nil {
		// Log detailed error for debugging
		fmt.Printf("Failed to send email: %v\n", err)
		return err
	}

	return nil
}
