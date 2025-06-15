package service

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/eskokado/startup-auth-go/backend/pkg/domain/vo"
	"gopkg.in/gomail.v2"
)

type MailSender interface {
	DialAndSend(...*gomail.Message) error
}

type EmailService struct {
	sender      MailSender
	from        string
	frontendURL string
}

func NewEmailService(sender MailSender) *EmailService {
	if sender == nil {
		host := os.Getenv("SMTP_HOST")
		portStr := os.Getenv("SMTP_PORT")
		username := os.Getenv("SMTP_USERNAME")
		password := os.Getenv("SMTP_PASSWORD")

		port := parsePort(portStr)
		if port == 0 {
			port = 587
		}

		sender = gomail.NewDialer(host, port, username, password)
	}

	return &EmailService{
		sender:      sender,
		from:        os.Getenv("FROM_EMAIL"),
		frontendURL: os.Getenv("FRONTEND_RESET_URL"),
	}
}

func (s *EmailService) SendResetPasswordEmail(email vo.Email, token string) error {
	if s.from == "" {
		return errors.New("FROM_EMAIL não está definido")
	}

	if s.frontendURL == "" {
		return errors.New("FRONTEND_RESET_URL não está definido")
	}

	m := gomail.NewMessage()
	m.SetHeader("From", s.from)
	m.SetHeader("To", email.String())
	m.SetHeader("Subject", "Redefinição de Senha")

	resetLink := fmt.Sprintf("%s?reset_password_token=%s", s.frontendURL, token)

	htmlBody := fmt.Sprintf(`
		<html>
		<body>
			<h2>Redefinição de Senha</h2>
			<p>Clique no link abaixo para redefinir sua senha:</p>
			<a href="%s">%s</a>
			<p>Este link expira em 1 hora.</p>
		</body>
		</html>
	`, resetLink, resetLink)

	m.SetBody("text/html", htmlBody)

	return s.sender.DialAndSend(m)
}

func ParsePort(port string) (int, error) {
	if port == "" {
		return 0, errors.New("porta não fornecida")
	}
	p, err := strconv.Atoi(port)
	if err != nil {
		return 0, fmt.Errorf("porta inválida: %s", port)
	}
	return p, nil
}

func parsePort(port string) int {
	p, _ := ParsePort(port)
	return p
}
