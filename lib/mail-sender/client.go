package mail_sender

import (
	"crypto/tls"
	"github.com/zerospam/check-firewall/lib/tls-generator"
	"github.com/zerospam/check-smtp/lib"
	"github.com/zerospam/check-smtp/lib/environment-vars"
	"github.com/zerospam/check-smtp/test-email"
	"log"
	"net/smtp"
	"time"
)

type Client struct {
	*smtp.Client
	server       *lib.TransportServer
	localName    string
	tlsGenerator *tlsgenerator.CertificateGenerator
}

func NewClient(server lib.TransportServer, localName string, timeout time.Duration) (*Client, *lib.SmtpError) {
	conn, err := server.Connect(timeout)
	if err != nil {
		return nil, lib.NewSmtpError(TIMEOUT, err)
	}

	client, err := smtp.NewClient(conn, server.Server)

	if err != nil {
		return nil, lib.NewSmtpError(CONN, err)
	}
	return &Client{
		Client:    client,
		localName: localName,
		server:    &server,
	}, nil
}

func (c *Client) getClientTLSConfig(commonName string) *tls.Config {
	if c.tlsGenerator == nil {
		c.tlsGenerator = tlsgenerator.NewClient(time.Now(), 30*365*24*time.Hour)
	}

	return c.tlsGenerator.GetTlsClientConfig(commonName)
}

//Try to send the test email
func (c *Client) SendTestEmail(email test_email.TestEmail) *lib.SmtpError {

	defer c.Client.Quit()
	defer c.Client.Close()

	var err error

	if err = c.Client.Hello(environmentvars.GetVars().SmtpCN); err != nil {
		return lib.NewSmtpError(EHLO, err)
	}

	if tlsSupport, _ := c.Client.Extension("STARTTLS"); tlsSupport {
		tlsConfig := c.getClientTLSConfig(environmentvars.GetVars().SmtpCN)
		tlsConfig.ServerName = c.server.Server
		tlsConfig.MinVersion = tls.VersionTLS11
		err = c.Client.StartTLS(tlsConfig)
		if err != nil {
			log.Printf("Couldn't start TLS transaction: %s", err)
			return lib.NewSmtpError(STARTTLS, err)
		}
		state, _ := c.Client.TLSConnectionState()
		tlsVersion := tlsgenerator.TlsVersion(state)
		log.Printf("[%s] TLS: %s", c.server.Server, tlsVersion)
	}

	if err = c.Client.Mail(environmentvars.GetVars().SmtpMailFrom.Address); err != nil {
		return lib.NewSmtpError(MAIL_FROM, err)
	}

	if err = c.Client.Rcpt(c.server.TestEmail); err != nil {
		return lib.NewSmtpError(RCPT_TO, err)
	}

	w, err := c.Data()

	if err != nil {
		return lib.NewSmtpError(DATA, err)
	}

	_, err = w.Write([]byte(email.String()))
	if err != nil {
		return lib.NewSmtpError(DATA, err)
	}

	err = w.Close()

	if err != nil {
		return lib.NewSmtpError(DATA, err)
	}

	if err = c.Client.Quit(); err != nil {
		return lib.NewSmtpError(QUIT, err)
	}

	return nil
}
