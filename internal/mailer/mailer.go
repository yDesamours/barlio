package mailer

import (
	"github.com/go-mail/mail"
)

type Mailer struct {
	sender string
	dialer *mail.Dialer
}

func NewMailer(port int, host, username, password, sender string) *Mailer {
	return &Mailer{
		sender: sender,
		dialer: mail.NewDialer(host, port, username, password),
	}
}

func (m *Mailer) Send(receiver string, infos map[string]string) (err error) {
	msg := mail.NewMessage()

	msg.SetBody("text/html", infos["body"])
	msg.AddAlternative("text/plain", infos["alternative"])
	msg.SetHeader("From", m.sender)
	msg.SetHeader("To", receiver)
	msg.SetHeader("Subject", infos["subject"])

	for i := 0; i < 3; i++ {
		err = m.dialer.DialAndSend(msg)
		if err == nil {
			break
		}
	}
	return
}
