package mail_sender

import (
	"crypto/tls"
	"github.com/zerospam/check-firewall/lib/tls-generator"
	"github.com/zerospam/check-smtp/lib"
	"github.com/zerospam/check-smtp/test-email"
	"log"
	"net"
	"net/smtp"
	"time"
)

type Client struct {
	*smtp.Client
	server        *lib.TransportServer
	localName     string
	tlsGenerator  *tlsgenerator.CertificateGenerator
	optTimeout    time.Duration
	connection    net.Conn
	lastError     *lib.SmtpError
	lastOperation *Operation
}

type SmtpOperation func() error

func NewClient(server lib.TransportServer, localName string, connTimeout time.Duration, optTimeout time.Duration) (*Client, *lib.SmtpError) {
	conn, err := server.Connect(connTimeout)
	if err != nil {
		return nil, lib.NewSmtpError(Timeout, err)
	}

	client, err := smtp.NewClient(conn, server.Server)

	if err != nil {
		return nil, lib.NewSmtpError(Connection, err)
	}

	return &Client{
		Client:     client,
		localName:  localName,
		server:     &server,
		optTimeout: optTimeout,
		connection: conn,
	}, nil
}

func (c *Client) getLastOperation() (*Operation, *lib.SmtpError) {
	return c.lastOperation, c.lastError
}

func (c *Client) getClientTLSConfig(commonName string) *tls.Config {
	if c.tlsGenerator == nil {
		c.tlsGenerator = tlsgenerator.NewClient(time.Now(), 30*365*24*time.Hour)
	}

	return c.tlsGenerator.GetTlsClientConfig(commonName)
}

func (c *Client) doOperation(operation Operation, optCallback SmtpOperation) {

	if c.lastError != nil {
		return
	}
	c.lastOperation = &operation

	err := c.connection.SetDeadline(time.Now().Add(c.optTimeout))
	if err != nil {
		c.lastError = lib.NewSmtpError(Timeout, err)
	}

	if err := optCallback(); err != nil {
		c.lastError = lib.NewSmtpError(operation, err)
	}

}

func (c *Client) setTls() error {
	if tlsSupport, _ := c.Client.Extension("STARTTLS"); !tlsSupport {
		return nil
	}
	tlsConfig := c.getClientTLSConfig(c.localName)
	tlsConfig.ServerName = c.server.Server
	tlsConfig.MinVersion = tls.VersionTLS11
	err := c.Client.StartTLS(tlsConfig)
	if err != nil {
		log.Printf("Couldn't start TLS transaction: %s", err)
		return err
	}
	state, _ := c.Client.TLSConnectionState()
	tlsVersion := tlsgenerator.TlsVersion(state)
	log.Printf("[%s] TLS: %s", c.server.Server, tlsVersion)
	return nil
}

//Try to send the test email
func (c *Client) SendTestEmail(email test_email.Email) *lib.SmtpError {

	defer c.Client.Close()

	c.doOperation(Ehlo, func() error {
		return c.Client.Hello(c.localName)
	})

	c.doOperation(StartTls, func() error {
		return c.setTls()
	})

	c.doOperation(MailFrom, func() error {
		return c.Client.Mail(email.From)
	})
	c.doOperation(RcptTo, func() error {
		return c.Client.Rcpt(c.server.TestEmail)
	})

	c.doOperation(Data, func() error {
		w, err := c.Data()

		if err != nil {
			return err
		}

		email.PrepareHeaders(c.localName)

		_, err = w.Write([]byte(email.String()))
		if err != nil {
			return err
		}

		err = w.Close()

		return err
	})

	c.doOperation(Quit, func() error {
		return c.Client.Quit()
	})

	return c.lastError
}

//Try to send the test email
func (c *Client) SpoofingTest(from string) *lib.SmtpError {

	c.lastError = nil

	defer c.Client.Quit()
	defer c.Client.Close()

	c.doOperation(Ehlo, func() error {
		return c.Client.Hello(c.localName)
	})

	c.doOperation(StartTls, func() error {
		return c.setTls()
	})

	c.doOperation(SpfFail, func() error {
		return c.Client.Mail(from)
	})
	c.doOperation(SpfFail, func() error {
		return c.Client.Rcpt(c.server.TestEmail)
	})

	return c.lastError
}
