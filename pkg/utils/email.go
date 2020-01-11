package utils

import (
	"net-alert/pkg/logging"

	"github.com/go-gomail/gomail"
)

var serverURL string
var portNumber int
var emailAccount string
var emailPassword string

//InitSMTP init the smtp client
func InitSMTP(server string, port int, email string, pass string) {
	serverURL = server
	portNumber = port
	emailAccount = email
	emailPassword = pass
}

//SendEmail send an email
func SendEmail(to string, title string, body string) {
	m := gomail.NewMessage()
	m.SetHeader("From", emailAccount)
	m.SetHeader("To", to)
	m.SetHeader("Subject", title)
	m.SetBody("text/html", body)
	d := gomail.NewPlainDialer(serverURL, portNumber, emailAccount, emailPassword)
	if err := d.DialAndSend(m); err != nil {
		logging.LogError(err)
	}
}
