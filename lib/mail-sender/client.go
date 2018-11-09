package mail_sender

import (
	"crypto/tls"
	"fmt"
	"github.com/zerospam/check-firewall/lib/tls-generator"
	"github.com/zerospam/check-smtp/lib"
	"github.com/zerospam/check-smtp/lib/mail-sender/smtp-commands"
	"github.com/zerospam/check-smtp/lib/tls_parser"
	"github.com/zerospam/check-smtp/test-email"
	"net/smtp"
	"strings"
	"time"
)

type Client struct {
	*smtp.Client
	server         *lib.TransportServer
	localName      string
	tlsGenerator   *tlsgenerator.CertificateGenerator
	optTimeout     time.Duration
	lastError      *lib.SmtpError
	lastCommand    *smtp_commands.Commands
	sentTestEmail  bool
	helloBanner    string
	tlsMinVersion  uint16
	tlsVersionUsed uint16
}

type SmtpOperation func() error

//Create new client to send the test email and test the SMTP server
func NewClient(server *lib.TransportServer, localName string, connTimeout time.Duration, optTimeout time.Duration, tlsMinVersion uint16) (*Client, *lib.SmtpError) {
	conn, err := server.Connect(connTimeout)
	if err != nil {
		return nil, lib.NewSmtpError(smtp_commands.Connection, err)
	}

	connection := NewConnection(conn)
	client, err := smtp.NewClient(connection, server.Server)

	if err != nil {
		return nil, lib.NewSmtpError(smtp_commands.Connection, err)
	}

	banner := strings.Trim(string(connection.firstRead), "\u0000\r\n")

	return &Client{
		Client:         client,
		localName:      localName,
		server:         server,
		optTimeout:     optTimeout,
		helloBanner:    banner,
		tlsMinVersion:  tlsMinVersion,
		tlsVersionUsed: 0,
	}, nil
}

func (c *Client) GetLastCommand() (*smtp_commands.Commands, *lib.SmtpError) {
	return c.lastCommand, c.lastError
}

func (c *Client) GetHelloBanner() (banner string, tlsVersion string) {

	if c.tlsVersionUsed == 0 {
		return c.helloBanner, "None"
	}

	tlsVersion, err := tls_parser.ToString(c.tlsVersionUsed)
	if err != nil {
		tlsVersion = err.Error()
	}

	return c.helloBanner, tlsVersion
}

func (c *Client) getClientTLSConfig(commonName string) *tls.Config {
	if c.tlsGenerator == nil {
		c.tlsGenerator = tlsgenerator.NewClient(time.Now(), 30*365*24*time.Hour)
	}

	return c.tlsGenerator.GetTlsClientConfig(commonName)
}

func (c *Client) doCommand(command smtp_commands.Commands, optCallback SmtpOperation) {

	if c.lastError != nil {
		return
	}

	c.lastCommand = &command
	//second parameter to not wait for a receiver.
	//This happen in the case the timeout returns before the command
	ch := make(chan error, 1)

	go func() {
		ch <- optCallback()
	}()

	timer := time.NewTimer(c.optTimeout)
	select {
	case err := <-ch:
		//Stop the timer as the command returned a result
		//@see https://golang.org/pkg/time/#After
		//Avoiding having a hanging timers
		timer.Stop()
		if err != nil {
			c.lastError = lib.NewSmtpError(command, err)
		}
		break

	case <-timer.C:
		c.lastError = lib.NewSmtpError(smtp_commands.Timeout, fmt.Errorf("CMD [%s] Timed out after %s", command, c.optTimeout))
		break
	}

}

func (c *Client) setTls() error {
	if tlsSupport, _ := c.Client.Extension("STARTTLS"); !tlsSupport {
		return nil
	}
	tlsConfig := c.getClientTLSConfig(c.localName)
	tlsConfig.ServerName = c.server.Server
	tlsConfig.MinVersion = c.tlsMinVersion
	//It's impossible to verify correctly the server in the case of a SMTP transaction
	//Better be permissive
	tlsConfig.InsecureSkipVerify = true
	err := c.Client.StartTLS(tlsConfig)
	if err != nil {
		return err
	}
	state, _ := c.Client.TLSConnectionState()
	c.tlsVersionUsed = state.Version

	return nil
}

//Try to send the test email
func (c *Client) SendTestEmail(email *test_email.Email) *lib.SmtpError {

	defer c.Client.Close()

	c.doCommand(smtp_commands.Ehlo, func() error {
		return c.Client.Hello(c.localName)
	})

	c.doCommand(smtp_commands.StartTls, func() error {
		return c.setTls()
	})

	c.doCommand(smtp_commands.MailFrom, func() error {
		return c.Client.Mail(email.From)
	})
	c.doCommand(smtp_commands.RcptTo, func() error {
		return c.Client.Rcpt(c.server.TestEmail)
	})

	c.doCommand(smtp_commands.Data, func() error {
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

	c.doCommand(smtp_commands.Quit, func() error {
		return c.Client.Quit()
	})

	return c.lastError
}

//Try to send the test email
func (c *Client) SpoofingTest(from string) *lib.SmtpError {

	defer c.Client.Quit()
	defer c.Client.Close()

	c.doCommand(smtp_commands.Ehlo, func() error {
		return c.Client.Hello(c.localName)
	})

	c.doCommand(smtp_commands.StartTls, func() error {
		return c.setTls()
	})

	c.doCommand(smtp_commands.SpfFail, func() error {
		return c.Client.Mail(from)
	})
	c.doCommand(smtp_commands.SpfFail, func() error {
		return c.Client.Rcpt(c.server.TestEmail)
	})

	return c.lastError
}
