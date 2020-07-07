package utils

import (
	"net/smtp"
)

//SMTPData struct for emails
type SMTPData struct {
	Host      string `json:"host"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	MockEmail string `json:"mock_email"`
}

//SendEmail func
func (s *SMTPData) SendEmail(to string, msg string) error {

	auth := smtp.PlainAuth("", s.Email, s.Password, s.Host)

	if err := smtp.SendMail(s.Host+":25", auth, s.Email, []string{to}, []byte(msg)); err != nil {
		return err
	}

	return nil
}
