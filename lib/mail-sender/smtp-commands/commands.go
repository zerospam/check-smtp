package smtp_commands

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type Commands int

const (
	Timeout Commands = iota
	Connection
	Ehlo
	StartTls
	MailFrom
	RcptTo
	Data
	Quit
	SpfFail
)

type CommandLog map[Commands]string

func (c Commands) String() string {
	names := []string{
		"TIMEOUT",
		"CONNECTION",
		"EHLO",
		"STARTTLS",
		"MAIL FROM",
		"RCPT TO",
		"DATA",
		"QUIT",
		"SPF-FAIL"}

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

func (i CommandLog) MarshalJSON() ([]byte, error) {
	x := make(map[string]string)
	for k, v := range i {
		x[fmt.Sprintf("%d/%s", k, k.String())] = v
	}
	marshal, e := json.Marshal(x)
	return marshal, e
}
