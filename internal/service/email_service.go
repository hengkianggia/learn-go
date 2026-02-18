package service

import (
	"fmt"
	"learn/internal/config"
	"learn/internal/pkg/email"
	"log/slog"

	"gopkg.in/gomail.v2"
)

type EmailService interface {
	SendOTP(to string, otp string) error
}

type emailService struct {
	logger *slog.Logger
}

func NewEmailService(logger *slog.Logger) EmailService {
	return &emailService{
		logger: logger,
	}
}

func (s *emailService) SendOTP(to string, otp string) error {
	from := config.AppConfig.SMTPFromEmail
	password := config.AppConfig.SMTPPassword
	smtpHost := config.AppConfig.SMTPHost
	smtpPort := config.AppConfig.SMTPPort

	// If SMTP credentials are not set (mock/dev), just log
	if smtpHost == "" || password == "" {
		s.logger.Warn("SMTP credentials not set, logging OTP instead", slog.String("otp", otp), slog.String("to", to))
		return nil
	}

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Your OTP Code")
	m.SetBody("text/plain", fmt.Sprintf("Your OTP code is: %s\n\nThis code will expire in 5 minutes.", otp))
	m.AddAlternative("text/html", email.GetOTPTemplate(otp))

	d := gomail.NewDialer(smtpHost, smtpPort, config.AppConfig.SMTPUser, password)

	// In case of self-signed certs or similar issues in dev environments (optional security relaxation)
	// d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		s.logger.Error("failed to send email", slog.String("error", err.Error()))
		return err
	}

	s.logger.Info("OTP sent successfully", slog.String("to", to))
	return nil
}
