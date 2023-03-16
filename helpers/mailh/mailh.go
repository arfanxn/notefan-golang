package mailh

import (
	"bytes"
	"net/smtp"
	"os"
	"strings"
)

// Send sends a mail
func Send(from string, subject string, body string, receipents ...string) (err error) {
	username := os.Getenv("MAIL_USERNAME")
	password := os.Getenv("MAIL_PASSWORD")
	host := os.Getenv("MAIL_HOST")
	port := os.Getenv("MAIL_PORT")
	addr := host + ":" + port

	messageBuf := bytes.NewBufferString("From: ")
	messageBuf.WriteString(from)
	messageBuf.WriteString("\r\n")
	messageBuf.WriteString("To: ")
	messageBuf.WriteString(strings.Join(receipents, ", "))
	messageBuf.WriteString("\r\n")
	messageBuf.WriteString("Subject: ")
	messageBuf.WriteString(subject)
	messageBuf.WriteString("\r\n\r\n")
	messageBuf.WriteString(body)

	auth := smtp.PlainAuth("", username, password, host)

	err = smtp.SendMail(addr, auth, from, receipents, messageBuf.Bytes())
	if err != nil {
		return
	}

	return nil
}
