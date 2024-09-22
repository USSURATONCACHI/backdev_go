package smtp_io

type SmtpClient interface {
	SendEmail(subject string, body string, toEmail string) error
}