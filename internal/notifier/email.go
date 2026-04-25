package notifier

import (
	"crypto/tls"
	"fmt"
	"strconv"

	"gopkg.in/gomail.v2"
)

type EmailNotifier struct {
	host     string
	port     int
	username string
	password string
	from     string
}

func NewEmailNotifier(host, portStr, username, password, from string) (*EmailNotifier, error) {
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("Invalid SMTP port: %s", portStr)
	}
	return &EmailNotifier{
		host:     host,
		port:     port,
		username: username,
		password: password,
		from:     from,
	}, nil
}

func (e *EmailNotifier) SendAdoptionConfirmation(toEmail, userName, animalName string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", e.from)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", "🐾 Заявка на усыновление принята — Happytail")
	m.SetBody("text/html", fmt.Sprintf(`
		<h2>Привет, %s!</h2>
		<p>Твоя заявка на усыновление <strong>%s</strong> успешно принята.</p>
		<p>Приют свяжется с тобой в ближайшее время для подтверждения.</p>
		<br>
		<p>С заботой, команда Happytail 🐾</p>
	`, userName, animalName))

	d := gomail.NewDialer(e.host, e.port, e.username, e.password)

	d.TLSConfig = &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         e.host,
	}
	return d.DialAndSend(m)
}
