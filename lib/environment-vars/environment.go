package environmentvars

import (
	"crypto/tls"
	"github.com/zerospam/check-smtp/lib"
	"github.com/zerospam/check-smtp/lib/mail-sender"
	"github.com/zerospam/check-smtp/lib/tls_parser"
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

		var emailSpoof *mail.Address
		emailFromSpoof := os.Getenv("SMTP_FROM_SPOOF")

		if emailFromSpoof == "" {
			emailFromSpoof = "spoof@amazon.com"
		}

		emailSpoof, err = mail.ParseAddress(emailFromSpoof)
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
			tlsVersionParsed, err = tls_parser.DecodeString(tlsVersion)
			if err != nil {
				panic(err)
			}
		}

		instance = &Env{
			ApplicationPort:       port,
			SharedKey:             os.Getenv("SHARED_KEY"),
			SmtpCN:                commonName,
			SmtpConnectionTimeout: timeoutParsed,
			SmtpOperationTimeout:  timeoutOptParsed,
			SmtpMailSpoof:         emailSpoof,
			TLSMinVersion:         tlsVersionParsed,
		}
	})
	return instance
}

func (e *Env) NewSmtpClient(server *lib.TransportServer) (*mail_sender.Client, *lib.SmtpError) {
	return mail_sender.NewClient(server, e.SmtpCN, e.SmtpConnectionTimeout, e.SmtpOperationTimeout, e.TLSMinVersion)
}
