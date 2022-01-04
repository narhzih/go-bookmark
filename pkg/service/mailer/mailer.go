package mailer

import (
	"github.com/jordan-wright/email"
	"net/smtp"
)

type Mailer struct {
	Auth        smtp.Auth
	Transporter *email.Email
	Addr        string
}

type MailConfig struct {
	Username string
	Password string
	SmtpHost string
	SmtpPort string
	MailFrom string
}

func NewMailer(config MailConfig) *Mailer {
	auth := smtp.PlainAuth("", config.Username, config.Password, config.SmtpHost)
	transporter := &email.Email{
		From: "Mypipe Desk <service@mypipe.app>",
	}
	return &Mailer{
		Auth:        auth,
		Addr:        config.SmtpHost + ":" + config.SmtpPort,
		Transporter: transporter,
	}
}

func (m *Mailer) SendVerificationEmail(mailTo []string) error {
	m.Transporter.HTML = []byte("<h1>Your verification code is 555555</h1>")
	m.Transporter.To = mailTo
	err := m.Transporter.Send(m.Addr, m.Auth)
	if err != nil {
		return err
	}
	return nil
}
