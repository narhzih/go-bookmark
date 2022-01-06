package mailer

import (
	"fmt"
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

func (m *Mailer) SendVerificationEmail(mailTo []string, token string) error {
	m.Transporter.HTML = []byte(fmt.Sprintf("<h2>YOur account verification token is %v</h2>", token))
	m.Transporter.Subject = "MyPipe account verification"
	m.Transporter.To = mailTo
	err := m.Transporter.Send(m.Addr, m.Auth)
	if err != nil {
		return err
	}
	return nil
}

func (m *Mailer) SendPasswordResetToken(mailTo []string, token string) error {
	m.Transporter.HTML = []byte(fmt.Sprintf("<h2>Your password reset token is %v</h2>", token))
	m.Transporter.Subject = "MyPipe account password reset"
	m.Transporter.To = mailTo
	err := m.Transporter.Send(m.Addr, m.Auth)
	if err != nil {
		return err
	}
	return nil
}
