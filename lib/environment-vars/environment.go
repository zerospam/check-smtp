package environmentvars

import (
	"crypto/tls"
	"fmt"
	"net/mail"
	"os"
	"sync"
	"time"
)

type Env struct {
	ApplicationPort       string
	SharedKey             string
	SmtpCN                string
	SmtpMailFrom          *mail.Address
	SmtpConnectionTimeout time.Duration
	SmtpOperationTimeout  time.Duration
	SmtpMailSpoof         *mail.Address
	TLSMinVersion         uint16
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

		var emailSpoof *mail.Address
		emailFromSpoof := os.Getenv("SMTP_FROM_SPOOF")

		if emailFromSpoof == "" {
			emailFromSpoof = "spoof@amazon.com"
		}

		emailSpoof, err = mail.ParseAddress(emailFrom)
		if err != nil {
			panic(err)
		}

		port := os.Getenv("PORT")
		if port == "" {
			port = "80"
		}

		timeoutParsed := 30 * time.Second
		timeout := os.Getenv("SMTP_CONN_TIMEOUT")
		if timeout != "" {
			timeoutParsed, err = time.ParseDuration(timeout)
			if err != nil {
				panic(err)
			}
		}

		timeoutOptParsed := 30 * time.Second
		timeoutOpt := os.Getenv("SMTP_OPT_TIMEOUT")
		if timeoutOpt != "" {
			timeoutOptParsed, err = time.ParseDuration(timeoutOpt)
			if err != nil {
				panic(err)
			}
		}

		tlsVersionParsed := uint16(tls.VersionTLS12)
		tlsVersion := os.Getenv("TLS_MIN_VERSION")
		if tlsVersion != "" {
			tlsVersionParsed, err = tlsDecodeString(tlsVersion)
			if err != nil {
				panic(err)
			}
		}

		instance = &Env{
			ApplicationPort:       port,
			SharedKey:             os.Getenv("SHARED_KEY"),
			SmtpCN:                commonName,
			SmtpMailFrom:          emailFromAddress,
			SmtpConnectionTimeout: timeoutParsed,
			SmtpOperationTimeout:  timeoutOptParsed,
			SmtpMailSpoof:         emailSpoof,
			TLSMinVersion:         tlsVersionParsed,
		}
	})
	return instance
}
