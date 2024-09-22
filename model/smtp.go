package model

import (
	"fmt"
	"net/smtp"
)

type SmtpInfo struct {
	Host string
	Port int16
	User string
	Password string

	FromEmail string
	MockUserEmail string
}

func (info *SmtpInfo) SendPlainAuth(subject string, body string, toEmail string) error {
	auth := smtp.PlainAuth("", info.User, info.Password, info.Host)

	addr := fmt.Sprintf("%s:%d", info.Host, info.Port)

	format := ("From: %s\r\n" +
			"To: %s\r\n" +
			"Subject: %s\r\n" +
			"\r\n" +
			"%s")

	message := []byte(fmt.Sprintf(format, info.FromEmail, toEmail, subject, body))
	err := smtp.SendMail(addr, auth, info.FromEmail, []string { toEmail }, message)

	return err
}