package mailer

import (
	"fmt"
	"net/smtp"
)

type Mailer struct {
	Auth smtp.Auth
}

type MailConfig struct {
	Username string
	Password string
	SmtpHost string
	SmtpPort string
}

func NewMailer(config MailConfig) Mailer {
	auth := smtp.PlainAuth("", config.Username, config.Password, config.SmtpHost)
	return Mailer{
		Auth: auth,
	}
}

func (m *Mailer) SendEmail(mailTo []string) (string, error) {
	from := "2a1aafa047bab7"
	password := "54c8d84430f82f"
	//username := "my username"
	smtpHost := "smtp.mailtrap.io"
	smtpPort := "2525"

	message := []byte("Your verification code is 444444")

	auth := smtp.PlainAuth("", from, password, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, mailTo, message)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return "", err
}
