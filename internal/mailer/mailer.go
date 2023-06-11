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

func (m *Mailer) Send(infos map[string]string) (err error) {
	msg := mail.Message{}

	msg.SetBody("text/html", infos["body"])
	msg.AddAlternative("text/plain", infos["alternative"])
	msg.SetHeader("From", m.sender)
	msg.SetHeader("To", infos["to"])
	msg.SetHeader("Object", infos["object"])

	for i := 0; i < 3; i++ {
		err = m.dialer.DialAndSend(&msg)
		if err != nil {
			break
		}
	}
	return
}
