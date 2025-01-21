package mailer

import (
	"bytes"
	"embed"
	"html/template"

	"github.com/go-mail/mail/v2"
)

type Mailer struct {
	sender string
	dailer *mail.Dialer
}

//go:embed "templates"
var templateFS embed.FS

func New(host string, port int, username, password, sender string) Mailer {
	dailer := mail.NewDialer(host, port, username, password)
	return Mailer{sender, dailer}
}

func (m *Mailer) Send(recipient, templateFile string, data any) error {
	tmpl, err := template.New("email").ParseFS(templateFS, "templates/"+templateFile)

	if err != nil {
		return err
	}

	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}

	plainBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(plainBody, "plainBody", data)
	if err != nil {
		return err
	}

	htmlBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(htmlBody, "htmlBody", data)
	if err != nil {
		return err
	}

	msg := mail.NewMessage()
	msg.SetHeader("To", recipient)
	msg.SetHeader("From", m.sender)
	msg.SetHeader("Subject", subject.String())
	msg.SetBody("text/plain", plainBody.String())
	msg.AddAlternative("text/html", htmlBody.String())

	err = m.dailer.DialAndSend(msg)
	if err != nil {
		return err
	}

	return nil
}
