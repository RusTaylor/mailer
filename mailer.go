package mailer

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/mail"
	"net/smtp"
)

type MailSender struct {
	From     mail.Address
	To       mail.Address
	Password string
	HostName string
	HostPort string
	Message  string
	Subject  string
}

// SendMail {
//		From: mail.Address{Name: "Name", Address: "your_email@inbox.ru"},
//		To: mail.Address{Name: "Name2", Address: "to@gmail.com"},
//		Password: "password",
//		HostName: "smtp.mail.ru",
//		HostPort: "465",
//		Subject: "message theme",
//		Message: "message body",
//	}
func (m MailSender) SendMail() {

	headers := make(map[string]string)
	headers["From"] = m.From.String()
	headers["To"] = m.To.String()
	headers["Subject"] = m.Subject

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + m.Message

	//host, _, _ := net.SplitHostPort(servername)

	auth := smtp.PlainAuth("", m.From.Address, m.Password, m.HostName)

	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         m.HostName,
	}

	conn, err := tls.Dial("tcp", m.HostName+":"+m.HostPort, tlsconfig)
	if err != nil {
		log.Panic(err)
	}

	a, err := smtp.NewClient(conn, m.HostName)
	if err != nil {
		log.Panic(err)
	}

	// Auth
	if err = a.Auth(auth); err != nil {
		log.Panic(err)
	}

	// To && From
	if err = a.Mail(m.From.Address); err != nil {
		log.Panic(err)
	}

	if err = a.Rcpt(m.To.Address); err != nil {
		log.Panic(err)
	}

	// Data
	w, err := a.Data()
	if err != nil {
		log.Panic(err)
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		log.Panic(err)
	}

	err = w.Close()
	if err != nil {
		log.Panic(err)
	}

	err = a.Quit()
	if err != nil {
		log.Panic(err)
	}
}
