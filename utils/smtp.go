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

	if err := smtp.SendMail(s.Host+":25", s.GetAuth(), s.Email, []string{to}, []byte(msg)); err != nil {
		return err
	}

	return nil
}

//GetAuth func
func (s *SMTPData) GetAuth() smtp.Auth {

	return smtp.PlainAuth("", s.Email, s.Password, s.Host)
}

//SendEmailWithAuth func
func (s *SMTPData) SendEmailWithAuth(auth smtp.Auth, msg []byte, to string) error {

	if err := smtp.SendMail(s.Host+":25", auth, s.Email, []string{to}, msg); err != nil {
		return err
	}

	return nil
}
