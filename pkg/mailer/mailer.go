package mailer

import (
	"bytes"
	"embed"
	"text/template"

	"github.com/go-mail/mail/v2"
	"github.com/omekov/superapp-backend/internal/auth/config"
)

//go:embed "templates"
var templateFS embed.FS

// MailerConfig ...

// Mailer ...
type Mailer struct {
	dailer *mail.Dialer
	config config.MailerConfig
	sender string
}

// New ...
func New(cfg config.MailerConfig) Mailer {
	dailer := mail.NewDialer(cfg.Host, cfg.Port, cfg.Username, cfg.Password)
	dailer.Timeout = cfg.Timeout

	return Mailer{
		dailer: dailer,
		sender: cfg.Sender,
		config: cfg,
	}
}

// Send ...
func (m Mailer) Send(to, templateFile, subject string, data interface{}) error {

	if m.config.TemplatePath == "" {
		m.config.TemplatePath = "templates/"
	}

	tmpl, err := template.ParseFS(templateFS, m.config.TemplatePath+templateFile)
	if err != nil {
		return err
	}

	htmlBody := new(bytes.Buffer)
	err = tmpl.Execute(htmlBody, data)
	if err != nil {
		return err
	}

	msg := mail.NewMessage()
	msg.SetHeader("To", to)
	msg.SetHeader("From", m.sender)
	msg.SetHeader("Subject", subject)
	msg.AddAlternative("text/html", htmlBody.String())

	return m.dailer.DialAndSend(msg)
}
