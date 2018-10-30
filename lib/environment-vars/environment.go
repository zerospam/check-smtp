package environmentvars

import (
	"fmt"
	"net/mail"
	"os"
	"sync"
)

type Env struct {
	ApplicationPort string
	SharedKey       string
	SmtpCN          string
	SmtpCheck       bool
	SmtpMailFrom    *mail.Address
}

var instance *Env
var once sync.Once

func GetVars() *Env {
	once.Do(func() {
		hostname, err := os.Hostname()
		if err != nil {
			panic(err)
		}

		commonName := os.Getenv("SMTP_CN")
		if commonName == "" {
			commonName = hostname
		}

		var emailFromAddress *mail.Address
		emailFrom := os.Getenv("SMTP_FROM")

		if emailFrom == "" {
			emailFrom = fmt.Sprintf("%s@%s", "local", hostname)
		}

		emailFromAddress, err = mail.ParseAddress(emailFrom)
		if err != nil {
			panic(err)
		}

		port := os.Getenv("PORT")
		if port == "" {
			port = "80"
		}

		instance = &Env{
			ApplicationPort: port,
			SharedKey:       os.Getenv("SHARED_KEY"),
			SmtpCN:          commonName,
			SmtpMailFrom:    emailFromAddress,
		}
	})
	return instance
}
