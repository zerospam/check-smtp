package smtp_commands

import "bytes"

type Commands int

const (
	Timeout Commands = iota
	Connection
	Ehlo
	StartTls
	RcptTo
	MailFrom
	Data
	Quit
	SpfFail
)

func (c Commands) String() string {
	names := []string{
		"TIMEOUT",
		"CONNECTION",
		"EHLO",
		"STARTTLS",
		"RCPT_TO",
		"MAIL_FROM",
		"DATA",
		"QUIT",
		"SPF_FAIL"}

	if c < Timeout || c > SpfFail {
		return "Unknown"
	}

	return names[c]
}

func (c *Commands) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(c.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}
