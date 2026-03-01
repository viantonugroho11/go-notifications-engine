package config

import (
	"crypto/tls"

	"gopkg.in/gomail.v2"
)

type Email struct {
	Host     string `json:"host"`
	Port     int `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
}


func InitializeEmail(email Email) (*gomail.Dialer, error) {
	
	dialer := gomail.NewDialer(
		email.Host,
		email.Port,
		email.User,
		email.Password,
	)

	dialer.TLSConfig = &tls.Config{
		InsecureSkipVerify: true,
	}
	return dialer, nil
}