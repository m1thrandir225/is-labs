package mail

type MailService interface {
	SendMail(from string, to string, subject, content string) error
}
